package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"os"
	"os/signal"
	"strings"

	"github.com/errnoh/gocurse/curses"
)

var controlsTable = make(map[string]string)
var client = http.Client{Jar: oneArg(cookiejar.New(nil)).(http.CookieJar)}
var reader = bufio.NewReader(os.Stdin)
var authRequestBody = []byte(`{
    "id": 1,
    "version": "1.0",
    "method": "actRegister",
    "params": [
        {
            "clientid": "GoRemoteController",
            "nickname": "go-remote",
            "level": "private"
        },
        [{
            "value": "yes",
            "function": "WOL"
        }]
    ]
}`)
var constrolsRequestBody = []byte(`{
    "id": 2,
    "version": "1.0",
    "method": "getRemoteControllerInfo",
    "params": []
}`)
var controlRequestEnvelope = `<?xml version="1.0"?>
<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
    <s:Body>
        <u:X_SendIRCC xmlns:u="urn:schemas-sony-com:service:IRCC:1">
            <IRCCCode>{signal}</IRCCCode>
        </u:X_SendIRCC>
    </s:Body>
</s:Envelope>`

type Config struct {
	Ip  string
	Url string
}
var config Config

type Envelope struct {
	Id     uint64
	Result []json.RawMessage
}
type Control struct {
	Name  string
	Value string
}
type Controls []Control

type Key struct {
	Command string
	Help    string
}
var keyboard = map[int]Key{ // TV functions
	curses.KEY_HOME:      {Command: "Home", Help: "home"},
	curses.KEY_UP:        {Command: "Up", Help: "Arrow Up"},
	curses.KEY_DOWN:      {Command: "Down", Help: "Arrow Down"},
	curses.KEY_LEFT:      {Command: "Left", Help: "Arrow Left"},
	curses.KEY_RIGHT:     {Command: "Right", Help: "Arrow Right"},
	curses.KEY_BACKSPACE: {Command: "Return", Help: "Backspace"},
	curses.KEY_ENTER:     {Command: "Confirm", Help: "Enter"},
	10:                   {Command: "Confirm", Help: "Enter"},
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
	119:                  {Command: "wide", Help: "W"},
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

	config.Ip = os.Args[1]
	config.Url = "http://" + config.Ip + "/sony"

	handshakeRequest, _ := http.NewRequest("POST", config.Url+"/accessControl", bytes.NewBuffer(authRequestBody))
	handshakeRequest.Header.Set("content-type", "application/json")
	response, err := client.Do(handshakeRequest)
	nuke(err, "Could not find device.")

	if 401 == response.StatusCode { // authorize in case client is not authorized yet
		pin := prompt("Enter provided pin code: ")
		authRequest, _ := http.NewRequest("POST", config.Url+"/accessControl", bytes.NewBuffer(authRequestBody))
		authRequest.Header.Set("content-type", "application/json")
		authRequest.Header.Set("authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(":"+pin)))
		response, err = client.Do(authRequest) // notice we are overwriting response from handshakeRequest
		nuke(err, "Authentication faile. Is the pin code right?")
	}

	if 200 == response.StatusCode { // let's get UPnP control list from device
		controlsResponse, err := http.Post(config.Url+"/system", "application/json", bytes.NewBuffer(constrolsRequestBody))
		nuke(err, "Could not retrieve UPnP control list from device.")
		envelope := new(Envelope)
		json.NewDecoder(controlsResponse.Body).Decode(envelope)
		if 0 == len(envelope.Result) {
			fmt.Fprintln(os.Stderr, "Could not retrieve UPnP control list from device.")
			os.Exit(1)
		}
		var controls Controls
		json.Unmarshal(envelope.Result[1], &controls)
		for _, control := range controls {
			controlsTable[string(control.Name)] = control.Value
		}
		runDaemon()
	}
}

func runDaemon() {
	// handle program termination
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		curses.Endwin()
		os.Exit(0)
	}()
	defer curses.Endwin()
	
	curses.Noecho()
	screen, _ := curses.Initscr()
	screen.Keypad(true) // interpret escape sequences
	screen.Addstr(0, 0, fmt.Sprintf("Use the keyboard:"), 0)
	for {
		keyCode := screen.Getch()
		key, ok := dispatch(keyCode)
		if ok {
			screen.Move(0, 1)
			screen.Clrtoeol()
			screen.Addstr(0, 1, fmt.Sprintf("%s >> %s", key.Command, config.Ip), 0)
		}
	}
}

func dispatch(keyCode int) (key Key, ok bool) {
	key, ok = keyboard[keyCode]
	if ok {
		signal, ok := controlsTable[key.Command]
		if ok {
			request, _ := http.NewRequest(
				"POST",
				config.Url+"/IRCC",
				bytes.NewBuffer(
					[]byte(
						strings.Replace(controlRequestEnvelope, "{signal}", signal, -1))))
			request.Header.Set("content-type", "text/xml; charset=UTF-8")
			request.Header.Set("soapaction", "urn:schemas-sony-com:service:IRCC:1#X_SendIRCC")
			client.Do(request)
		}
	}
	return
}

func prompt(message string) (str string) {
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

func oneArg(x interface{}, _ ...interface{}) interface{} {
	return x
}
