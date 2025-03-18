package main

import (
	"fmt"
	"io"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/fsnotify/fsnotify"
)

func init() {
	fmt.Println(`
                 _      __ ______ / /_ _____ / /_   ___   _____
                | | /| / // __  // __// ___// __ \ / _ \ / ___/
                | |/ |/ // /_/ // /_ / /__ / / / //  __// /    
                |__/|__/ \____/ \__/ \___//_/ /_/ \___//_/      
    `)
}

var logger = io.Discard

func main() {
	var args args
	arg.MustParse(&args)

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

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fatalf("watcher: error: %s\n", err.Error())
	}
	defer watcher.Close()

	if args.Recursive {
		err = addPathRecursively(watcher, args.Path)
	} else {
		err = watcher.Add(args.Path)
	}
	if err != nil {
		fatalf("watcher: error: %s\n", err.Error())
	}

	watchEvents(watcher, *c)
}
