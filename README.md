GoRemote
========

I got fed up of searching for my TV remote controller below the sofa and pillows only to find it later inside the refrigerator.

So I quickly hacked this amazing (POG) **Virtual TV remote controller**:

<img
    src="https://github.com/marcioAlmada/goremote/blob/master/screenshots/cover.png"
    title="GTK gui VS curses interface" />

Besides the "physical indexing" issue:

- Sometimes it's much easier to use a notebook to control your TV
- The remote control needs battery to be changed once in a while, virtual one doesn't
- I needed a Go pet project to maintain and keep the gopher on my desk alive

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
the auth screen, [here](blob/master/gtkApplication.go#L34).

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
