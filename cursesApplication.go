package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/marcioAlmada/goremote/upnp"
	"github.com/tncardoso/gocurses"
)

type cursesApplication struct {
	screen *gocurses.Window
	application
}

func newCursesApplication() cursesApplication {
	return cursesApplication{}
}

func (cursesApplication) PromptPinCode() (str string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter provided pin code: ")
	str, err := reader.ReadString('\n')
	nuke(err, "Could not read input.")
	str = strings.TrimRight(str, "\n\r")
	return
}

// Run executes the command line curses application
func (app cursesApplication) Run(client upnp.Client, keyMap, altKeyMap keyMap) (e error) {
	// handle process termination
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		gocurses.End()
		os.Exit(1)
	}()
	// add alternative key bidings
	keyMap.Merge(altKeyMap)
	// env config
	os.Setenv("ESCDELAY", "0")
	// curses config
	gocurses.CursSet(0)
	app.screen = gocurses.Initscr()
	defer gocurses.End()
	app.screen.Keypad(true)   // enable keypad support
	app.screen.Scrollok(true) // infinite screen
	gocurses.Noecho()         // avoid char leak of unmapped keys
	// run the REPL!
	app.screen.Addstr("> Ready to rumble!")
	for {
		ch := app.screen.Getch()
		if 4 == ch { // handles CTRL+D
			break
		}
		go func() {
			if key, ok := keyMap[ch]; ok {
				if _, _, e := client.SendCommand(key.Command); e == nil {
					app.screen.Addstr("\n", client.IP, "> ", key.Command)
				} else {
					app.screen.Addstr("\n", "ERROR: ", e)
				}
				app.screen.Refresh()
			}
		}()
	}
	return
}
