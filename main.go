package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"marcioAlmada/tv/upnp"
	"github.com/jessevdk/go-flags"
	"github.com/errnoh/gocurse/curses"
)

}

type options struct {
	Args struct { IP string } `positional-args:"yes" required:"yes"`
	REPL bool `short:"r" long:"repl" description:"Use a command line REPL session instead of gui"`
}

func main() {
	var options options;
	_, e := flags.Parse(&options)
	nuke(e)

	client := upnp.NewClient(options.Args.IP)
	response, e := client.Handshake()
	nuke(e, "Could not find device.")

	if 401 == response.StatusCode { // authorize in case client is not authorized yet
		pin := prompt("Enter provided pin code: ") // get pin code from devie
		response, e = client.Authorize(pin)
		nuke(e, "Authentication failed. Is the pin code right?")
	}

	if 200 == response.StatusCode { // let's get UPnP control list from device
		_, e := client.RequestControlsList()
		runREPL(client)
		nuke(e)
	}
}

func runREPL(client upnp.Client) {
	// handle process termination
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		curses.Endwin()
		os.Exit(0)
	}()
	defer curses.Endwin()
	// curses gui
	curses.Noecho()
	screen, _ := curses.Initscr()
	screen.Keypad(true) // interpret escape sequences
	screen.Addstr(0, 0, fmt.Sprintf("Use the keyboard:"), 0)
	screen.Move(0, 1)
	// add alternative key bidings to actionsMap
	actionsMap.Merge(alternativeMap)
	// run the REPL!
	for {
		keyCode := screen.Getch()
		key, ok := actionsMap[keyCode]
		if ok { // is the key mapped? Otherwise ignore it
			signal, ok, _, _ := client.SendCommand(key.Command)
			if ok { // show status on screen when request is made
				screen.Move(0, 1)
				screen.Clrtoeol()
				screen.Addstr(0, 1, fmt.Sprintf("%s (%s) >> %s", key.Command, signal, client.IP), 0)
			}
		}
	}
}

func prompt(message string) (str string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(message)
	str, err := reader.ReadString('\n')
	nuke(err, "Could not read input.")
	str = strings.TrimRight(str, "\n\r")
	return
}

func nuke(err error, msg ...string) {
	if err != nil {
		if(len(msg) > 1){
			fmt.Fprintln(os.Stderr, msg[1])
		}
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
