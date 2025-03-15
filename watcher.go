package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

func watchEvents(watcher *fsnotify.Watcher, cf CommandsFile) {
	if watcher == nil {
		panic("watcher is nil!")
	}
	eventTime := time.Now()
	var lastEvent fsnotify.Op

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if !(time.Since(eventTime) > (time.Millisecond*400) || lastEvent != event.Op) {
				continue
			}
			switch event.Op.String() {
			case fsnotify.Write.String():
				handleEvent(cf.Write)
			case fsnotify.Create.String():
				handleEvent(cf.Create)
			case fsnotify.Remove.String():
				handleEvent(cf.Remove)
			case fsnotify.Rename.String():
				handleEvent(cf.Rename)
			case fsnotify.Chmod.String():
				handleEvent(cf.Chmod)
			}
			handleEvent(cf.Common)
			eventTime = time.Now()
			lastEvent = event.Op
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		}
	}
}

func handleEvent(cmds []string) {
	for _, v := range cmds {
		go func(cmd string) {
			if cmd := wrapCmd(parseCommand(v)); cmd != nil {
				err := cmd.Run()
				if err != nil {
					fmt.Fprintf(os.Stderr, "error running command %q: %s\n", v, err)
					return
				}
			}
		}(v)
	}
}
