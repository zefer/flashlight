package main

import (
	"flag"
	"net/http"
	"runtime"
	"strings"

	"github.com/zefer/mothership/mpd"
	"github.com/zefer/mpdlcd/lcd"
)

var (
	client  *mpd.Client
	mpdAddr = flag.String("mpdaddr", "127.0.0.1:6600", "MPD address")
)

func banana(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	flag.Parse()
	lcd.Start()
	defer lcd.Stop()

	// This client connection provides an API to MPD's commands.
	client = mpd.NewClient(*mpdAddr)
	defer client.Close()

	// This watcher notifies us when MPD's state changes, without polling.
	watch := mpd.NewWatcher(*mpdAddr)
	defer watch.Close()
	watch.OnStateChange(func(s string) {
		// http://www.musicpd.org/doc/protocol/command_reference.html#command_idle
		if s == "player" {
			displayStatus()
		}
	})

	displayStatus()

	for {
		runtime.Gosched()
	}
}

func trim(s string) string {
	if len(s) > lcd.Cols {
		return s[0:lcd.Cols]
	}
	return s
}

func displayStatus() {
	d, err := mpdStatus()
	if err != nil {
		return
	}

	var msg string

	if d["Artist"] != nil && d["Title"] != nil {
		msg = trim(d["Artist"].(string)) + "\n" + trim(d["Title"].(string))
	} else if d["file"] != nil {
		parts := strings.Split(d["file"].(string), "/")
		file := parts[len(parts)-1]
		msg = trim(file)
		if msg != file {
			msg += "\n" + trim(file[lcd.Cols:])
		}
	} else {
		msg = `¯\_(ツ)_/¯`
	}

	lcd.Display(msg)
}

func mpdStatus() (map[string]interface{}, error) {
	out := map[string]interface{}{}
	s, err := client.C.Status()
	if err != nil {
		return nil, err
	}
	for k, v := range s {
		out[k] = v
	}

	s, err = client.C.CurrentSong()
	if err != nil {
		return nil, err
	}
	for k, v := range s {
		out[k] = v
	}

	return out, err
}
