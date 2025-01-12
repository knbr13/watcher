package main

import (
	"fmt"
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

func main() {
	var args args
	arg.MustParse(&args)

	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s\n", err.Error())
		os.Exit(1)
	}

	if args.Path != "" && !validPath(args.Path) {
		fmt.Fprintf(os.Stderr, "err: %s\n", "invalid path: "+args.Path)
		os.Exit(1)
	}

	if args.Path == "" {
		args.Path = wd
	}

	c := &CommandsFile{}

	data, err := os.ReadFile(args.File)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s\n", err.Error())
		os.Exit(1)
	}

	err = c.Parse(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s\n", err.Error())
		os.Exit(1)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s\n", err.Error())
		os.Exit(1)
	}
	defer watcher.Close()

	if args.Recursive {
		addPathRecursively(args.Path, watcher)
	} else {
		watcher.Add(args.Path)
	}

	watchEvents(watcher, *c)
}
