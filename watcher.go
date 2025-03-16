package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
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

			fName := filepath.Base(event.Name)

			switch event.Op.String() {
			case fsnotify.Write.String():
				handleEvent(cf.Write, fName)
			case fsnotify.Create.String():
				handleEvent(cf.Create, fName)
			case fsnotify.Remove.String():
				handleEvent(cf.Remove, fName)
			case fsnotify.Rename.String():
				handleEvent(cf.Rename, fName)
			case fsnotify.Chmod.String():
				handleEvent(cf.Chmod, fName)
			}
			handleEvent(cf.Common, fName)

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

func handleEvent(rules []Rule, fName string) {
	for _, rule := range rules {
		go func(rule Rule) {
			if !matchesPattern(fName, rule.Pattern) {
				return
			}
			for _, cmd := range rule.Commands {
				if rule.Sequential {
					if cmd := wrapCmd(parseCommand(cmd)); cmd != nil {
						err := cmd.Run()
						if err != nil {
							fmt.Fprintf(os.Stderr, "error running command %q: %s\n", cmd, err)
							return
						}
					}
					continue
				}
				go func(cmd string) {
					if cmd := wrapCmd(parseCommand(cmd)); cmd != nil {
						err := cmd.Run()
						if err != nil {
							fmt.Fprintf(os.Stderr, "error running command %q: %s\n", cmd, err)
							return
						}
					}
				}(cmd)
			}
		}(rule)
	}
}

func addPathRecursively(watcher *fsnotify.Watcher, root string) error {
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "watch error: %s\n", err.Error())
			return nil
		}
		if !d.IsDir() || slices.Contains(excludedFolders, strings.ToLower(d.Name())) {
			return nil
		}
		return watcher.Add(path)
	})
	return err
}
