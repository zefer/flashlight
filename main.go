package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/zefer/mothership/mpd"
	"github.com/zefer/mpdlcd/lcd"
	"gopkg.in/airbrake/glog.v1"
	"gopkg.in/airbrake/gobrake.v1"
)

var (
	client      *mpd.Client
	mpdAddr     = flag.String("mpdaddr", "127.0.0.1:6600", "MPD address")
	abProjectID = flag.Int64("abprojectid", 0, "Airbrake project ID")
	abApiKey    = flag.String("abapikey", "", "Airbrake API key")
	abEnv       = flag.String("abenv", "development", "Airbrake environment name")
)

func main() {
	flag.Parse()

	// Free-up GPIO resources on exit.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		lcd.Stop()
		os.Exit(1)
	}()

	if *abProjectID > int64(0) && *abApiKey != "" {
		airbrake := gobrake.NewNotifier(*abProjectID, *abApiKey)
		airbrake.SetContext("environment", *abEnv)
		glog.Gobrake = airbrake
	}

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

	// Prevent program exit.
	<-make(chan int)
}

func displayStatus() {
	d, err := mpdStatus()
	if err != nil {
		return
	}

	if d["state"] == "play" {
		displayPlaying(d)
	} else {
		lcd.Clear()
	}
}

func trim(s string) string {
	if len(s) > lcd.Cols {
		return s[0:lcd.Cols]
	}
	return s
}

func displayPlaying(state map[string]interface{}) {
	var msg string
	if state["Artist"] != nil && state["Title"] != nil {
		msg = trim(state["Artist"].(string)) + "\n" + trim(state["Title"].(string))
	} else if state["file"] != nil {
		parts := strings.Split(state["file"].(string), "/")
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
