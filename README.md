# Watcher

Watcher is a tool written in Go that monitors file system events and executes specified commands in response to those events.

## Features

- Watch a directory or file for file system events.
- Run custom commands on different types of events (create, write, chmod, remove, rename).
- Optionally watch subdirectories recursively.

### Prerequisites

- Go (version 1.13 or later)

### Usage

```bash
./watcher --path "/path/to/watch" --file "/path/to/commands-file-in-yaml-format" -r
```

#### Command Line Options:
    "-p, --path": Set the path to the directory to watch for events.
    "-f, --file": Set the path to the file that contains the commands to run on each event.
    "-r, --recursive": Watch subdirectories recursively.
