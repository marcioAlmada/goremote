package main

import "github.com/tncardoso/gocurses"

type key struct {
	Command string
	Help    string
}

type keyMap map[int]key

func (k keyMap) Merge(m keyMap) keyMap {
	for code, key := range m {
		k[code] = key
	}
	return k
}

// action maps, map keys are system key codes
var defaultKeyMap = keyMap{
	gocurses.KEY_UP:    {Command: "Up", Help: "Arrow Up"},
	gocurses.KEY_DOWN:  {Command: "Down", Help: "Arrow Down"},
	gocurses.KEY_LEFT:  {Command: "Left", Help: "Arrow Left"},
	gocurses.KEY_RIGHT: {Command: "Right", Help: "Arrow Right"},
	//
	262: {Command: "Home", Help: "home"},
	194: {Command: "Confirm", Help: "Enter"},
	27:  {Command: "Return", Help: "Esc"},
	32:  {Command: "Play", Help: "Space"},
	111: {Command: "Options", Help: "O"},
	339: {Command: "ChannelUp", Help: "PageUp"},
	338: {Command: "ChannelDown", Help: "PageDown"},
	43:  {Command: "VolumeUp", Help: "+"},
	45:  {Command: "VolumeDown", Help: "-"},
	109: {Command: "Mute", Help: "M"},
	116: {Command: "Theater", Help: "T"},
	110: {Command: "Netflix", Help: "N"},
	106: {Command: "Jump", Help: "J"},
	119: {Command: "Wide", Help: "W"},
	112: {Command: "PAP", Help: "P"},
	100: {Command: "Display", Help: "D"},
	99:  {Command: "SceneSelect", Help: "C"},
	115: {Command: "ClosedCaption", Help: "S"},
	104: {Command: "iManual", Help: "H"},
	105: {Command: "Input", Help: "I"},
	267: {Command: "Mode3D", Help: "F3"},
	107: {Command: "KeyPad", Help: "K"},
	102: {Command: "FootballMode", Help: "F"},
	276: {Command: "PowerOff", Help: "F12"},
	114: {Command: "Red", Help: "R"},
	103: {Command: "Green", Help: "G"},
	121: {Command: "Yellow", Help: "Y"},
	98:  {Command: "Blue", Help: "B"},
	330: {Command: "PicOff", Help: "Delete"},
	46:  {Command: "DOT", Help: "."},
	48:  {Command: "Num0", Help: "0"},
	49:  {Command: "Num1", Help: "1"},
	50:  {Command: "Num2", Help: "2"},
	51:  {Command: "Num3", Help: "3"},
	52:  {Command: "Num4", Help: "4"},
	53:  {Command: "Num5", Help: "5"},
	54:  {Command: "Num6", Help: "6"},
	55:  {Command: "Num7", Help: "7"},
	56:  {Command: "Num8", Help: "8"},
	57:  {Command: "Num9", Help: "9"},
	// no key bind yet
	// 900: {Command: "Social", Help: ""},
	// 901: {Command: "Forward", Help: ""},
	// 902: {Command: "Rewind", Help: ""},
	// 903: {Command: "Prev", Help: ""},
	// 904: {Command: "Stop", Help: ""},
	// 905: {Command: "Next", Help: ""},
	// 906: {Command: "GGuide", Help: ""},
	// 907: {Command: "EPG", Help: ""},
	// 908: {Command: "Favorites", Help: ""},
	// 910: {Command: "Analog", Help: ""},
	// 911: {Command: "Teletext", Help: ""},
	// 912: {Command: "Exit", Help: ""},
	// 913: {Command: "Analog2", Help: ""},
	// 914: {Command: "AD", Help: ""},
	// 915: {Command: "Digital", Help: ""},
	// 916: {Command: "Analog?", Help: ""},
	// 917: {Command: "BS", Help: ""},
	// 918: {Command: "CS", Help: ""},
	// 919: {Command: "BSCS", Help: ""},
	// 920: {Command: "Ddata", Help: ""},
	// 921: {Command: "Tv_Radio", Help: ""},
	// 922: {Command: "SEN", Help: ""},
	// 923: {Command: "InternetWidgets", Help: ""},
	// 924: {Command: "InternetVideo", Help: ""},
	// 925: {Command: "PAP", Help: ""},
	// 926: {Command: "MyEPG", Help: ""},
	// 927: {Command: "ProgramDescription", Help: ""},
	// 928: {Command: "WriteChapter", Help: ""},
	// 929: {Command: "TrackID", Help: ""},
	// 930: {Command: "TenKey", Help: ""},
	// 931: {Command: "AppliCast", Help: ""},
	// 932: {Command: "acTVila", Help: ""},
	// 933: {Command: "DeleteVideo", Help: ""},
	// 934: {Command: "PhotoFrame", Help: ""},
	// 935: {Command: "TvPause", Help: ""},
	// 936: {Command: "KeyPad", Help: ""},
	// 937: {Command: "Media", Help: ""},
	// 938: {Command: "SyncMenu", Help: ""},
	// 939: {Command: "Rec", Help: ""},
	// 940: {Command: "Eject", Help: ""},
	// 941: {Command: "FlashPlus", Help: ""},
	// 942: {Command: "FlashMinus", Help: ""},
	// 943: {Command: "TopMenu", Help: ""},
	// 944: {Command: "PopUpMenu", Help: ""},
	// 945: {Command: "RakurakuStart", Help: ""},
	// 946: {Command: "OneTouchTimeRec", Help: ""},
	// 947: {Command: "OneTouchView", Help: ""},
	// 948: {Command: "OneTouchRec", Help: ""},
	// 949: {Command: "OneTouchStop", Help: ""},
	// 950: {Command: "DUX", Help: ""},
}

// alternative key bindigs
var alternativeKeyMap = keyMap{
	263: {Command: "Return", Help: "Backspace"},
	10:  {Command: "Confirm", Help: "Enter"}, // fallback when curses KEY_ENTER fails
}
