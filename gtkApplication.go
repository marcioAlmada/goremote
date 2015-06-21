package main

import (
	//#cgo pkg-config: glib-2.0 gobject-2.0 gtk+-3.0
	//#include "bindings.go.h"
	//#include "resources.h"
	"C"
	"os"
	"unsafe"

	"github.com/conformal/gotk3/glib"
	"github.com/conformal/gotk3/gtk"
	"github.com/marcioAlmada/goremote/upnp"
)

// connectable is any object capable of respond to gtk signals
type connectable interface {
	Connect(detailedSignal string, f interface{}, userData ...interface{}) (glib.SignalHandle, error)
}

type gtkApplication struct {
	preferDarkTheme uint
	builder         *gtk.Builder
	goSettings      *C.struct__GtkSettings

	application
}

func newGtkApplication() gtkApplication {
	return gtkApplication{}
}

func (gtkApplication) PromptPinCode() (str string) {
	return
}

// Run executes the graphical application
func (app gtkApplication) Run(client upnp.Client, keyMap, _ keyMap) (e error) {
	gtk.Init(&os.Args)
	app.preferDarkTheme = 0
	app.goSettings = C.gtk_settings_get_default()
	app.builder, _ = gtk.BuilderNew()
	if e = app.builder.AddFromResource("/marcio/gocontroller/resources/main.ui"); e != nil {
		return
	}
	if obj, e := app.builder.GetObject("Window"); e == nil {
		// try to load ui main window object
		if window, ok := obj.(*gtk.Window); ok {
			window.SetTitle("GoRemote")
			// link ui buttons manually
			// wrapper for builder.ConnectSignals(nil) is not ready in gotk3 yet
			// see https://github.com/conformal/gotk3/issues/50
			for _, key := range keyMap {
				command := key.Command
				if button, ok := app.getObject(key.Command); ok {
					button.Connect("clicked", func() {
						go client.SendCommand(command)
					})
				}
			}
			// theme toggler
			if toggler, ok := app.getObject("ThemeToggler"); ok {
				toggler.Connect("clicked", func() { app.toggleDarkTheme() })
			}
			// make app prefer dark theme by default
			app.toggleDarkTheme()
			// start gui
			window.Connect("destroy", func() { gtk.MainQuit() })
			window.ShowAll()
			gtk.Main()
		}
	}
	return
}

// pools gtk ui and returns a gtk widget that implements connectable interface
func (app gtkApplication) getObject(id string) (widget connectable, ok bool) {
	if obj, e := app.builder.GetObject(id); e == nil {
		if widget, ok = obj.(*gtk.Button); ok {
			return widget, ok
		} else if widget, ok = obj.(*gtk.ToggleButton); ok {
			return widget, ok
		} else if widget, ok = obj.(*gtk.Window); ok {
			return widget, ok
		}
	}
	return
}

// this executes C glib and gtk API while gotk3 has no g_object_set access
// see https://github.com/conformal/gotk3/issues/95
func (app *gtkApplication) toggleDarkTheme() {
	alternateTable := []uint{1, 0}
	app.preferDarkTheme = alternateTable[app.preferDarkTheme]
	cstr := C.CString("gtk-application-prefer-dark-theme")
	defer C.free(unsafe.Pointer(cstr))
	c := C.guint(app.preferDarkTheme)
	p := unsafe.Pointer(&c)
	C._g_object_set_one(C.gpointer(app.goSettings), (*C.gchar)(cstr), p)
}
