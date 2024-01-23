package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

func main() {
	var path, command string
	var write, create, chmod, remove, rename bool
	
	wd, _ := os.Getwd()

	flag.StringVar(&command, "cmd", "", "command to run when new event occur")
	flag.StringVar(&path, "path", wd, "path to the directory to watch for events on")
	flag.BoolVar(&write, "write", true, "write event")
	flag.BoolVar(&create, "create", false, "create event")
	flag.BoolVar(&chmod, "chmod", false, "chmod event")
	flag.BoolVar(&remove, "remove", false, "remove event")
	flag.BoolVar(&rename, "rename", false, "rename event")
	flag.Parse()

	m := map[fsnotify.Op]bool{
		fsnotify.Remove: remove,
		fsnotify.Chmod:  chmod,
		fsnotify.Create: create,
		fsnotify.Rename: rename,
		fsnotify.Write:  write,
	}

	var atLeastOneTrue bool
	for _, v := range m {
		if v {
			atLeastOneTrue = true
		}
	}

	if !atLeastOneTrue {
		fmt.Fprint(os.Stderr, "err: at least one event should be specified\n")
		os.Exit(1)
	}

	_, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %v\n", err)
		os.Exit(1)
	}

	watcher, _ := fsnotify.NewWatcher()
	err = addSubdirectories(path, watcher)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %v\n", err)
		os.Exit(1)
	}

	eventTime := time.Now()

	for {
		select {
		case event := <-watcher.Events:
			for k, ok := range m {
				if event.Has(k) && ok && time.Since(eventTime) > time.Millisecond*500 {
					cmds := parseCommands(command)
					for _, c := range cmds {
						c.Stderr = os.Stderr
						c.Stdout = os.Stdout

						err = c.Start()
						if err != nil {
							fmt.Fprintf(os.Stderr, "err: %v\n", err)
							os.Exit(1)
						}
					}
					eventTime = time.Now()
				}
			}
		case err := <-watcher.Errors:
			fmt.Fprintf(os.Stderr, "err: %v\n", err)
			os.Exit(1)
		}
	}
}
