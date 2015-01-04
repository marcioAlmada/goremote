GoRemote
========

A quickly hacked (POG alert) **Virtual TV remote controller** because sometimes it's much easier
to use a notebook to control your "smart" TV through Wi-Fi network:

<img
    src="https://github.com/marcioAlmada/goremote/blob/master/screenshots/cover.png"
    title="GTK gui VS curses interface" />

> This was tested with a Sony Bravia KDL-50W805B and alike. PRs supporting other devices are welcome.

## Install

```
go get github.com/marcioAlmada/goremote
cd $GOPATH/src/github.com/marcioAlmada/goremote
make
make install
```

## Usage

### GUI

To use the GTK GUI interface:

```bash
$ goremote <tv-ip> # click the buttons :)
```

> First time TV access can't be done with the GUI yet because I was too lazy to implement
the auth screen, [here](/gtkApplication.go#L34).

### Terminal

To use the curses (terminal) interface:

```bash
$ goremote <tv-ip> --repl # use the computer keyboard :)
```

> Ex: my TV has a fixed IP so I always: `goremote 10.0.0.101 --repl`.

## Is My Smart TV Compatible?

With fingers crossed, try to run:

```bash
curl "http://<tv-ip>/sony/accessControl/actRegister"
# OR
curl "http://<tv-ip>/accessControl/actRegister"
```

If you get `{"error":[501,"Not Implemented"]}` response your TV is probably compatible :)

## To Do

- Add the pin code promt screen to GTK app
- Add device autodiscover so it won't be necessary to type the IP
- Add support to other smart TVs

## Copyright

Copyright (c) 2014 - 2015 Márcio Almada. Distributed under the terms of an MIT-style license.
See LICENSE for details.
