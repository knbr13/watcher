package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/fsnotify/fsnotify"
)

// shutdownGracePeriod bounds how long watcher waits for in-flight commands
// to finish after a shutdown signal before giving up and exiting anyway.
const shutdownGracePeriod = 10 * time.Second

// version is set at build time via -ldflags "-X main.version=...".
var version = "dev"

var logger = io.Discard

func main() {
	var args args
	arg.MustParse(&args)

	if args.Version {
		fmt.Printf("watcher %s\n", version)
		return
	}

	fmt.Println(`
                 _      __ ______ / /_ _____ / /_   ___   _____
                | | /| / // __  // __// ___// __ \ / _ \ / ___/
                | |/ |/ // /_/ // /_ / /__ / / / //  __// /
                |__/|__/ \____/ \__/ \___//_/ /_/ \___//_/
    `)

	if args.File == "" {
		fatalf("watcher: error: --file is required\n")
	}

	if args.Debug {
		logger = os.Stderr
	}

	wd, err := os.Getwd()
	if err != nil {
		fatalf("watcher: error: %s\n", err.Error())
	}

	if args.Path != "" && !validPath(args.Path) {
		fatalf("watcher: error: invalid path %q\n", args.Path)
	}

	if args.Path == "" {
		args.Path = wd
	}

	c := &CommandsFile{}

	data, err := os.ReadFile(args.File)
	if err != nil {
		fatalf("watcher: error: %s\n", err.Error())
	}

	err = c.Parse(data)
	if err != nil {
		fatalf("watcher: error: %s\n", err.Error())
	}

	if args.Debounce != 0 {
		c.Debounce = args.Debounce
	}
	if c.Debounce == 0 {
		c.Debounce = 400 * time.Millisecond
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fatalf("watcher: error: %s\n", err.Error())
	}
	defer watcher.Close()

	if args.Recursive {
		err = addPathRecursively(args.Path, watcher)
	} else {
		err = watcher.Add(args.Path)
	}
	if err != nil {
		fatalf("watcher: error: %s\n", err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\n👋  Shutting down watcher...")
		cancel()
	}()

	watchEvents(ctx, watcher, *c, args.Path, &wg)

	drained := make(chan struct{})
	go func() {
		wg.Wait()
		close(drained)
	}()

	select {
	case <-drained:
	case <-time.After(shutdownGracePeriod):
		fmt.Fprintf(os.Stderr, "watcher: warning: commands still running after %v grace period, exiting anyway\n", shutdownGracePeriod)
	}
}
