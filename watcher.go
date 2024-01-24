package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/fsnotify/fsnotify"
)

type watcherOptions struct {
	path            string
	commands        [][]string
	registredEvents []fsnotify.Op
	recursive       bool
}

func (opt *watcherOptions) print() {
	//TODO: colorize the option with emoji
	fmt.Println("watcher version 0.1.0")
	fmt.Println("path: ", opt.path)
	fmt.Println("commands: ", opt.commands)
	fmt.Println("events: ", opt.registredEvents)
	fmt.Println("recursive: ", opt.recursive)
}

func watchEvents(watcher *fsnotify.Watcher, options watcherOptions) {
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
			for _, op := range options.registredEvents {
				if event.Has(op) && (time.Since(eventTime) > (time.Millisecond*400) || lastEvent != event.Op) {
					for _, s := range options.commands {
						cmd := exec.Command(s[0], s[1:]...)
						cmd.Stdout = os.Stdout
						cmd.Stderr = os.Stderr
						err := cmd.Start()
						if err != nil {
							fmt.Fprintf(os.Stderr, "error: can't start command: %s\n", err.Error())
							continue
						}
						eventTime = time.Now()
						lastEvent = event.Op
						// TODO: a way to organize the output of the commands in a better way
					}
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		}
	}
}
