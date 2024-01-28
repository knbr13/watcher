package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
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
	fmt.Println("👀  Watcher v0.1.0")
	fmt.Printf("📂  Path: %s\n", opt.path)
	fmt.Printf("🔍  Events: %v\n", opt.registredEvents)
	fmt.Printf("🔄  Recursive: %v\n", opt.recursive)

	if len(opt.commands) > 0 {
		fmt.Println("🚀  Commands to run:")
		for _, command := range opt.commands {
			fmt.Printf("    %v\n", strings.Join(command, " "))
		}
	} else {
		fmt.Println("⚠️   No commands specified to run on events. Events will be printed to stdout.")
	}

	fmt.Println("\nlistening for events...")
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
					if len(options.commands) == 0 {
						fmt.Printf("%s  %s\n", time.Now().Format("2006-01-02 15:04:05"), event)
						continue
					}
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
