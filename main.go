package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type watcherOptions struct {
	path            string
	commands        [][]string
	registredEvents []fsnotify.Op
	recursive       bool
}

func (opt *watcherOptions) print() {
	//TODO: colorize the option with emoji
	fmt.Println("watcher version 0.1.0")
	fmt.Println("path: ", opt.path)
	fmt.Println("commands: ", opt.commands)
	fmt.Println("events: ", opt.registredEvents)
	fmt.Println("recursive: ", opt.recursive)
}

func validateAndParseFlags(
	commands string,
	path string,
	events string,
	recursive bool,
) (opt watcherOptions) {
	// a function to validate the flag and build the watcherOptions struct
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %v\n", err)
		os.Exit(1)
	} else {
		opt.path = path
	}

	if recursive && !fileInfo.IsDir() {
		fmt.Fprintf(os.Stderr, "err: %v\n", "recursive flag can only be used with directories")
		os.Exit(1)
	} else {
		opt.recursive = recursive
	}

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

func watchEvents(watcher *fsnotify.Watcher, options watcherOptions) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			for _, op := range options.registredEvents {
				if event.Has(op) {
					log.Printf("%v on %v => executing command %v\n", event.Op, event.Name, options.commands)
					for _, s := range options.commands {
						cmd := exec.Command(s[0], s[1:]...)
						cmd.Stdout = os.Stdout
						cmd.Stderr = os.Stderr
						err := cmd.Start()
						if err != nil {
							fmt.Fprintf(os.Stderr, "error: can't start command: %s\n", err.Error())
							continue
						}
						// TODO: a way to organize the output of the commands in a better way
					}
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error: ", err)
		}
	}
}

func main() {
	var path, command string
	var recursive bool
	var events string

	wd, _ := os.Getwd()

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
