package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/marcioAlmada/goremote/upnp"
)

type application interface {
	PromptPinCode() string
	Run(upnp.Client, keyMap, keyMap) error
}

type options struct {
	Args struct{ IP string } `positional-args:"yes" required:"yes"`
	REPL bool                `short:"r" long:"repl" description:"Use a command line REPL session instead of gui"`
}

func main() {
	storagePath := os.Getenv("HOME") + "/.gocontroller/"

	os.MkdirAll(storagePath, 0777) // setup

	var options options
	_, e := flags.Parse(&options)
	nuke(e)

	var app application
	if options.REPL {
		app = newCursesApplication()
	} else {
		app = newGtkApplication()
	}

	client := upnp.NewClient(options.Args.IP, fileGetContents(storagePath+options.Args.IP))

	response, e := client.Handshake()
	nuke(e, "Could not find device.")

	if 401 == response.StatusCode { // authorize in case client is not authorized yet
		pin := app.PromptPinCode() // get pin code from device
		response, e = client.Authorize(pin)
		nuke(e, "Authentication failed. Is the pin code right?")
	}

	if 200 == response.StatusCode { // authenticated
		_, e := client.RequestControlsList() // let's get UPnP control list from device
		nuke(e, "Maybe device is off?")
		client.RequestSystemInformation()
		filePutContents(storagePath+client.IP, client.Mac)    // cache device info
		e = app.Run(client, defaultKeyMap, alternativeKeyMap) // run!
		nuke(e, "Failed to launch application")
	}
}

func nuke(err error, msgs ...interface{}) {
	if err != nil {
		if len(msgs) > 0 {
			fmt.Fprintln(os.Stderr, msgs...)
		}
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func getStoragePath() string {
	return os.Getenv("HOME") + "/.gocontroller/"
}

func fileGetContents(filename string) string {
	fp, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return ""
	}
	defer fp.Close()
	reader := bufio.NewReader(fp)
	bytes, _ := ioutil.ReadAll(reader)
	contents := strings.TrimRight(string(bytes), "\n\r")
	return contents
}

func filePutContents(filename string, content string) error {
	fp, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer fp.Close()
	_, err = fp.Write([]byte(content))
	return err
}
