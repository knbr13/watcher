package main

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fsnotify/fsnotify"
)

// fmt.Println("👀  Watcher v0.1.0")
// fmt.Printf("📂  Path: %s\n", opt.path)
// fmt.Printf("🔍  Events: %v\n", opt.registredEvents)
// fmt.Printf("🔄  Recursive: %v\n", opt.recursive)

func watchEvents(watcher *fsnotify.Watcher, cf CommandsFile, rootPath string) {
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

			// Global Path filtering
			relPath, err := filepath.Rel(rootPath, event.Name)
			if err != nil {
				relPath = event.Name
			}

			include := cf.Include
			if len(include) == 0 {
				include = cf.Watch
			}

			if !shouldProcess(relPath, include, cf.Exclude) {
				continue
			}

			switch event.Op {
			case fsnotify.Write:
				handleEvent(cf.Write, event)
			case fsnotify.Create:
				handleEvent(cf.Create, event)
			case fsnotify.Remove:
				handleEvent(cf.Remove, event)
			case fsnotify.Rename:
				handleEvent(cf.Rename, event)
			case fsnotify.Chmod:
				handleEvent(cf.Chmod, event)
			}
			handleEvent(cf.Common, event)

			eventTime = time.Now()
			lastEvent = event.Op
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Fprintf(logger, "watcher: error: %s\n", err.Error())
		}
	}
}

func handleEvent(rules []Rule, event fsnotify.Event) {
	for _, rule := range rules {
		go func(rule Rule) {
			if !matchesPattern(event.Name, rule.Pattern) {
				return
			}
			var wg sync.WaitGroup
			var errOccurred atomic.Bool
			for _, cmdStr := range rule.Commands {
				timeout := rule.Timeout
				if rule.Sequential {
					if cmd := wrapCmd(parseCommand(cmdStr), event); cmd != nil {
						exitCode, err := runCommand(cmd, timeout)
						if err != nil {
							fmt.Fprintf(logger, "watcher: error running command %q: %s\n", cmdStr, err.Error())
						}
						if exitCode != 0 {
							errOccurred.Store(true)
						}
					}
					continue
				}
				wg.Add(1)
				go func(cmdStr string) {
					defer wg.Done()
					if cmd := wrapCmd(parseCommand(cmdStr), event); cmd != nil {
						exitCode, err := runCommand(cmd, timeout)
						if err != nil {
							fmt.Fprintf(logger, "watcher: error running command %q: %s\n", cmdStr, err.Error())
						}
						if exitCode != 0 {
							errOccurred.Store(true)
						}
					}
				}(cmdStr)
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
	for _, cmdStr := range cmds {
		if cmd := wrapCmd(parseCommand(cmdStr), event); cmd != nil {
			_, err := runCommand(cmd, 0)
			if err != nil {
				fmt.Fprintf(logger, "watcher: error running post command %q: %s\n", cmdStr, err.Error())
			}
		}
	}
}
