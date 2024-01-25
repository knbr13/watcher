package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

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
	var path, command string
	var recursive bool
	var events string

	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %v%s\n", "invalid path: ", path)
	}

	flag.StringVar(&command, "cmd", "", "command to run when new event occur")
	flag.StringVar(&path, "path", wd, "path to the directory to watch for events on")
	flag.StringVar(&events, "events", "all", "events to watch for (write, create, chmod, remove, rename, all)")
	flag.BoolVar(&recursive, "r", false, "watch subdirectories recursively")
	flag.Parse()

	options := validateAndParseFlags(
		command,
		path,
		events,
		recursive,
	)
	options.print()
	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()

	if options.recursive {
		addSubdirectories(options.path, watcher)
	} else {
		watcher.Add(options.path)
	}

	go watchEvents(watcher, options)
	<-make(chan struct{})
}

// validateAndParseFlags function to validate the flag and build the watcherOptions struct
func validateAndParseFlags(
	commands string,
	path string,
	events string,
	recursive bool,
) (opt watcherOptions) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %v%s\n", "invalid path: ", path)
		os.Exit(1)
	}
	opt.path = path

	if recursive && !fileInfo.IsDir() {
		fmt.Fprintf(os.Stderr, "err: you can't use recursive flag with a file\n")
		os.Exit(1)
	}
	opt.recursive = recursive

	// parse events
	if events == "all" {
		opt.registredEvents = []fsnotify.Op{
			fsnotify.Write,
			fsnotify.Create,
			fsnotify.Chmod,
			fsnotify.Remove,
			fsnotify.Rename,
		}
	} else {
		events = strings.ToLower(events)
		eventsList := strings.Split(events, ",")
		for _, event := range eventsList {
			event = strings.TrimSpace(event)
			switch event {
			case "write":
				opt.registredEvents = append(opt.registredEvents, fsnotify.Write)
			case "create":
				opt.registredEvents = append(opt.registredEvents, fsnotify.Create)
			case "chmod":
				opt.registredEvents = append(opt.registredEvents, fsnotify.Chmod)
			case "remove":
				opt.registredEvents = append(opt.registredEvents, fsnotify.Remove)
			case "rename":
				opt.registredEvents = append(opt.registredEvents, fsnotify.Rename)
			default:
				fmt.Fprintf(os.Stderr, "err: %v%s\n", "invalid event: ", event)
				os.Exit(1)
			}
		}
	}
	opt.commands = parseCommands(commands)
	return
}
