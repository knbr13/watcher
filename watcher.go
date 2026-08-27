package main

import (
	"context"
	"errors"
	"fmt"
	"os"
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

func watchEvents(ctx context.Context, watcher *fsnotify.Watcher, cf CommandsFile, rootPath string, opts *runOptions, wg *sync.WaitGroup) {
	if watcher == nil {
		panic("watcher is nil!")
	}
	eventTimes := make(map[string]time.Time)
	lastOps := make(map[string]fsnotify.Op)

	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			// A newly created directory needs its own watch added, or
			// files created inside it later would go unnoticed.
			if opts.Recursive && event.Op.Has(fsnotify.Create) {
				if info, statErr := os.Stat(event.Name); statErr == nil && info.IsDir() {
					if err := addPathRecursively(event.Name, watcher, opts.ExcludeDirs); err != nil {
						fmt.Fprintf(logger, "watcher: error: failed to watch new directory %q: %s\n", event.Name, err.Error())
					}
				}
			}

			if !(time.Since(eventTimes[event.Name]) > cf.Debounce || lastOps[event.Name] != event.Op) {
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
				handleEvent(ctx, wg, cf.Write, event, opts)
			case fsnotify.Create:
				handleEvent(ctx, wg, cf.Create, event, opts)
			case fsnotify.Remove:
				handleEvent(ctx, wg, cf.Remove, event, opts)
			case fsnotify.Rename:
				handleEvent(ctx, wg, cf.Rename, event, opts)
			case fsnotify.Chmod:
				handleEvent(ctx, wg, cf.Chmod, event, opts)
			}
			handleEvent(ctx, wg, cf.Common, event, opts)

			eventTimes[event.Name] = time.Now()
			lastOps[event.Name] = event.Op
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Fprintf(logger, "watcher: error: %s\n", err.Error())
		}
	}
}

func handleEvent(ctx context.Context, wg *sync.WaitGroup, rules []Rule, event fsnotify.Event, opts *runOptions) {
	for _, rule := range rules {
		wg.Add(1)
		go func(rule Rule) {
			defer wg.Done()
			if !matchesPattern(event.Name, rule.Pattern) {
				return
			}
			data := getEventData(event)
			var cmdWg sync.WaitGroup
			var errOccurred atomic.Bool
			runOne := func(cmdStr string) {
				if opts.DryRun {
					fmt.Printf("[dry-run] would run: %s\n", cmdStr)
					return
				}
				if cmd := wrapCmd(parseCommand(cmdStr), data); cmd != nil {
					exitCode, err := runCommand(ctx, cmd, rule.Timeout)
					if err != nil {
						fmt.Fprintf(logger, "watcher: error running command %q: %s\n", cmdStr, err.Error())
					}
					if exitCode != 0 {
						errOccurred.Store(true)
					}
				}
			}
			for _, cmdStr := range rule.Commands {
				cmdStr = expandTemplate(cmdStr, data)
				if rule.Sequential {
					runOne(cmdStr)
					continue
				}
				cmdWg.Add(1)
				go func(cmdStr string) {
					defer cmdWg.Done()
					runOne(cmdStr)
				}(cmdStr)
			}
			cmdWg.Wait()
			if errOccurred.Load() {
				runPostCommands(ctx, rule.OnFailure, data, opts)
			} else {
				runPostCommands(ctx, rule.OnSuccess, data, opts)
			}
		}(rule)
	}
}

// runCommand starts cmd and waits for it to finish, killing the whole
// process tree if timeout elapses (when timeout > 0) or ctx is canceled
// first, whichever comes first.
func runCommand(ctx context.Context, cmd *exec.Cmd, timeout time.Duration) (int, error) {
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	if err := cmd.Start(); err != nil {
		return -1, err
	}

	done := make(chan error, 1)
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

		// non-exit errors (e.g., command not found)
		return -1, err

	case <-ctx.Done():
		killProcessTree(cmd)
		<-done
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return -1, fmt.Errorf("command timed out after %v", timeout)
		}
		return -1, fmt.Errorf("command canceled: watcher is shutting down")
	}
}

func runPostCommands(ctx context.Context, cmds []string, data EventData, opts *runOptions) {
	for _, cmdStr := range cmds {
		cmdStr = expandTemplate(cmdStr, data)
		if opts.DryRun {
			fmt.Printf("[dry-run] would run: %s\n", cmdStr)
			continue
		}
		if cmd := wrapCmd(parseCommand(cmdStr), data); cmd != nil {
			_, err := runCommand(ctx, cmd, 0)
			if err != nil {
				fmt.Fprintf(logger, "watcher: error running post command %q: %s\n", cmdStr, err.Error())
			}
		}
	}
}
