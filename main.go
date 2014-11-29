package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"./upnp"
	"github.com/errnoh/gocurse/curses"
)

type key struct {
	Command string
	Help    string
}

var keyboard = map[int]key{ // keyboard mapping
	curses.KEY_HOME:      {Command: "Home", Help: "home"},
	curses.KEY_UP:        {Command: "Up", Help: "Arrow Up"},
	curses.KEY_DOWN:      {Command: "Down", Help: "Arrow Down"},
	curses.KEY_LEFT:      {Command: "Left", Help: "Arrow Left"},
	curses.KEY_RIGHT:     {Command: "Right", Help: "Arrow Right"},
	curses.KEY_BACKSPACE: {Command: "Return", Help: "Backspace"},
	curses.KEY_ENTER:     {Command: "Confirm", Help: "Enter"},
	10:                   {Command: "Confirm", Help: "Enter"}, // fallback when KEY_ENTER fails
	27:                   {Command: "Exit", Help: "Esc"},
	32:                   {Command: "Play", Help: "Space"},
	111:                  {Command: "Options", Help: "O"},
	339:                  {Command: "ChannelUp", Help: "PageUp"},
	338:                  {Command: "ChannelDown", Help: "PageDown"},
	65:                   {Command: "VolumeUp", Help: "+"},
	66:                   {Command: "VolumeDown", Help: "-"},
	109:                  {Command: "Mute", Help: "M"},
	110:                  {Command: "Netflix", Help: "N"},
	106:                  {Command: "Jump", Help: "J"},
	119:                  {Command: "Wide", Help: "W"},
	112:                  {Command: "PAP", Help: "P"},
	100:                  {Command: "Display", Help: "D"},
	99:                   {Command: "SceneSelect", Help: "C"},
	115:                  {Command: "ClosedCaption", Help: "S"},
	104:                  {Command: "iManual", Help: "H"},
	105:                  {Command: "Input", Help: "I"},
	267:                  {Command: "Mode3D", Help: "F3"},
	107:                  {Command: "KeyPad", Help: "K"},
	102:                  {Command: "FootballMode", Help: "F"},
	276:                  {Command: "PowerOff", Help: "F12"},
	114:                  {Command: "Red", Help: "R"},
	103:                  {Command: "Green", Help: "G"},
	121:                  {Command: "Yellow", Help: "Y"},
	98:                   {Command: "Blue", Help: "B"},
	46:                   {Command: "DOT", Help: "."},
	48:                   {Command: "Num0", Help: "0"},
	49:                   {Command: "Num1", Help: "1"},
	50:                   {Command: "Num2", Help: "2"},
	51:                   {Command: "Num3", Help: "3"},
	52:                   {Command: "Num4", Help: "4"},
	53:                   {Command: "Num5", Help: "5"},
	54:                   {Command: "Num6", Help: "6"},
	55:                   {Command: "Num7", Help: "7"},
	56:                   {Command: "Num8", Help: "8"},
	57:                   {Command: "Num9", Help: "9"},
}

func main() {
	if len(os.Args) != 2 {
		nuke(errors.New("Missing argument 1"), "Please inform device address")
	}
	client := upnp.NewController(os.Args[1])
	response, e := client.Handshake()
	nuke(e, "Could not find device.")

	if 401 == response.StatusCode { // authorize in case client is not authorized yet
		pin := prompt("Enter provided pin code: ") // get pin code from devie
		response, e = client.Authorize(pin)
		nuke(e, "Authentication failed. Is the pin code right?")
	}

	if 200 == response.StatusCode { // let's get UPnP control list from device
		_, e := client.RequestControlsList()
		nuke(e, "Could not retrieve UPnP control list from device.")
		runDaemon(client)
	}
}

func runDaemon(client upnp.Controller) {
	handleProcTermination()
	curses.Noecho()
	screen, _ := curses.Initscr()
	screen.Keypad(true) // interpret escape sequences
	screen.Addstr(0, 0, fmt.Sprintf("Use the keyboard:"), 0)
	screen.Move(0, 1)
	for {
		keyCode := screen.Getch()
		// fmt.Println(keyCode)
		key, ok := keyboard[keyCode]
		if ok { // is the key mapped? Otherwise ignore it
			ok := client.SendCommand(key.Command)
			if ok { // show status on screen when request is made
				screen.Move(0, 1)
				screen.Clrtoeol()
				screen.Addstr(0, 1, fmt.Sprintf("%s >> %s", key.Command, client.IP), 0)
			}
		}
	}
}

func handleProcTermination() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		curses.Endwin()
		os.Exit(0)
	}()
	defer curses.Endwin()
}

func prompt(message string) (str string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(message)
	str, err := reader.ReadString('\n')
	nuke(err, "Could not read input.")
	str = strings.TrimRight(str, "\n\r")
	return
}

func nuke(err error, msg string) {
	if err != nil {
		fmt.Fprintln(os.Stderr, msg)
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
