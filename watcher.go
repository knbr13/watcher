package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

// fmt.Println("👀  Watcher v0.1.0")
// fmt.Printf("📂  Path: %s\n", opt.path)
// fmt.Printf("🔍  Events: %v\n", opt.registredEvents)
// fmt.Printf("🔄  Recursive: %v\n", opt.recursive)

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
				for _, v := range cf.Write {
					if cmd := wrapCmd(parseCommand(v)); cmd != nil {
						err := cmd.Run()
						if err != nil {
							fmt.Fprintf(os.Stderr, "error running command %q: %s\n", v, err)
							continue
						}
					}
				}
			case fsnotify.Create.String():
				for _, v := range cf.Create {
					if cmd := wrapCmd(parseCommand(v)); cmd != nil {
						err := cmd.Run()
						if err != nil {
							fmt.Fprintf(os.Stderr, "error running command %q: %s\n", v, err)
							continue
						}
					}
				}
			case fsnotify.Remove.String():
				for _, v := range cf.Remove {
					if cmd := wrapCmd(parseCommand(v)); cmd != nil {
						err := cmd.Run()
						if err != nil {
							fmt.Fprintf(os.Stderr, "error running command %q: %s\n", v, err)
							continue
						}
					}
				}
			case fsnotify.Rename.String():
				for _, v := range cf.Rename {
					if cmd := wrapCmd(parseCommand(v)); cmd != nil {
						err := cmd.Run()
						if err != nil {
							fmt.Fprintf(os.Stderr, "error running command %q: %s\n", v, err)
							continue
						}
					}
				}
			case fsnotify.Chmod.String():
				for _, v := range cf.Chmod {
					if cmd := wrapCmd(parseCommand(v)); cmd != nil {
						err := cmd.Run()
						if err != nil {
							fmt.Fprintf(os.Stderr, "error running command %q: %s\n", v, err)
							continue
						}
					}
				}
			}
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
