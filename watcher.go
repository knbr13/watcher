package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
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
				handleEvent(cf.Write, event)
			case fsnotify.Create.String():
				handleEvent(cf.Create, event)
			case fsnotify.Remove.String():
				handleEvent(cf.Remove, event)
			case fsnotify.Rename.String():
				handleEvent(cf.Rename, event)
			case fsnotify.Chmod.String():
				handleEvent(cf.Chmod, event)
			}
			handleEvent(cf.Common, event)

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

func handleEvent(rules []Rule, event fsnotify.Event) {
	fName := filepath.Base(event.Name)

	for _, rule := range rules {
		go func(rule Rule) {
			if !matchesPattern(fName, rule.Pattern) {
				return
			}
			var wg sync.WaitGroup
			var errOccurred atomic.Bool
			for _, cmd := range rule.Commands {
				timeout, _ := time.ParseDuration(rule.Timeout.String())
				if rule.Sequential {
					if cmd := wrapCmd(parseCommand(cmd, event)); cmd != nil {
						exitCode, _ := runCommand(cmd, timeout)
						if exitCode != 0 {
							errOccurred.Store(true)
						}
					}
					continue
				}
				wg.Add(1)
				go func(cmd string) {
					defer wg.Done()
					if cmd := wrapCmd(parseCommand(cmd, event)); cmd != nil {
						exitCode, _ := runCommand(cmd, timeout)
						if exitCode != 0 {
							errOccurred.Store(true)
						}
					}
				}(cmd)
			}
			wg.Wait()
			if errOccurred.Load() {
				runPostCommands(rule.OnFailure, event)
			} else {
				runPostCommands(rule.OnSuccess, event)
			}
		}(rule)
	}
}

func runCommand(cmd *exec.Cmd, timeout time.Duration) (int, error) {
	if timeout <= 0 {
		err := cmd.Run()
		if err == nil {
			return 0, nil
		}

		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode(), err
		}

		// non-exit errors (e.g., command not found)
		return -1, err
	}

	if err := cmd.Start(); err != nil {
		return -1, err
	}

	done := make(chan error)
	go func() { done <- cmd.Wait() }()

	select {
	case err := <-done:
		if err == nil {
			return 0, nil
		}

		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode(), err
		}
		return -1, err

	case <-time.After(timeout):
		cmd.Process.Kill()
		<-done
		return -1, fmt.Errorf("command timed out after %v", timeout)
	}
}

func runPostCommands(cmds []string, event fsnotify.Event) {
	for _, cmd := range cmds {
		if cmd := wrapCmd(parseCommand(cmd, event)); cmd != nil {
			_, _ = runCommand(cmd, 0)
		}
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
