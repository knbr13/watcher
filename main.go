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
		fatalf("err: %s\n", err.Error())
	}

	if args.Path != "" && !validPath(args.Path) {
		fatalf("err: invalid path %q\n", args.Path)
	}

	if args.Path == "" {
		args.Path = wd
	}

	c := &CommandsFile{}

	data, err := os.ReadFile(args.File)
	if err != nil {
		fatalf("err: %s\n", err.Error())
	}

	err = c.Parse(data)
	if err != nil {
		fatalf("err: %s\n", err.Error())
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fatalf("err: %s\n", err.Error())
	}
	defer watcher.Close()

	if args.Recursive {
		addPathRecursively(args.Path, watcher)
	} else {
		watcher.Add(args.Path)
	}

	watchEvents(watcher, *c)
}
