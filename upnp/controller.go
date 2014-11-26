package upnp

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

// Controller is a client to comunicate with remote device
type Controller struct {
	IP            string
	url           string
	client        http.Client
	controlsTable map[string]string
}

// NewController instantiates new device remote controller client
func NewController(ip string) Controller {
	cookieJar, _ := cookiejar.New(nil)
	return Controller{
		IP:            ip,
		url:           strings.Replace("http://{ip}/sony/{endpoint}", "{ip}", ip, -1),
		client:        http.Client{Jar: cookieJar},
		controlsTable: make(map[string]string),
	}
}

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

// Handshake with remote device
func (c Controller) Handshake() (response *http.Response, e error) {
	request, _ := http.NewRequest("POST", c.getURL("accessControl"), bytes.NewBuffer(authRequestBody))
	request.Header.Set("content-type", "application/json")
	response, e = c.client.Do(request)
	return
}

// Authorize with remote device providing an "Authorization: Basic *" header
func (c Controller) Authorize(pin string) (response *http.Response, e error) {
	request, _ := http.NewRequest("POST", c.getURL("accessControl"), bytes.NewBuffer(authRequestBody))
	request.Header.Set("content-type", "application/json")
	request.Header.Set("authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(":"+pin)))
	response, e = c.client.Do(request)
	return
}

type envelope struct {
	ID     uint64            `json:"id"`
	Result []json.RawMessage `json:"result"`
}
type control struct {
	Name  string
	Value string
}
type controls []control

var constrolsRequestBody = []byte(`{
    "id": 2,
    "version": "1.0",
    "method": "getRemoteControllerInfo",
    "params": []
}`)

// RequestControlsList gets UPnP controller description from remote device
func (c Controller) RequestControlsList() (response *http.Response, e error) {
	request, _ := http.NewRequest("POST", c.getURL("system"), bytes.NewBuffer(constrolsRequestBody))
	request.Header.Set("content-type", "application/json")
	response, e = c.client.Do(request)

	if e == nil {
		envelope := new(envelope)
		json.NewDecoder(response.Body).Decode(envelope)
		if 0 == len(envelope.Result) {
			e = errors.New("Could not retrieve UPnP control list from device.")
		} else {
			var controls controls
			json.Unmarshal(envelope.Result[1], &controls)
			for _, control := range controls {
				c.controlsTable[string(control.Name)] = control.Value
			}
		}
	}
	return
}

var commandRequestBody = `<?xml version="1.0"?>
<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
    <s:Body>
        <u:X_SendIRCC xmlns:u="urn:schemas-sony-com:service:IRCC:1">
            <IRCCCode>{signal}</IRCCCode>
        </u:X_SendIRCC>
    </s:Body>
</s:Envelope>`

// SendCommand sends a signal to activate a remote device function
func (c Controller) SendCommand(command string) (ok bool) {
	signal, ok := c.controlsTable[command]
	if ok {
		request, _ := http.NewRequest(
			"POST",
			c.getURL("IRCC"),
			bytes.NewBuffer(
				[]byte(
					strings.Replace(commandRequestBody, "{signal}", signal, -1))))
		request.Header.Set("content-type", "text/xml; charset=UTF-8")
		request.Header.Set("soapaction", "urn:schemas-sony-com:service:IRCC:1#X_SendIRCC")
		c.client.Do(request)
	}
	return
}

func (c Controller) getURL(endpoint string) string {
	return strings.Replace(c.url, "{endpoint}", endpoint, -1)
}
