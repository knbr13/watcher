# Watcher

Watcher is a simple tool written in Go that monitors file system events on a specific path and executes specified commands in response to those events.

## Features

- Watch a directory or file for file system events.
- Run custom commands on different types of events (create, write, chmod, remove, rename).
- Optionally watch subdirectories recursively.

### Usage

```bash
./watcher --path "/path/to/watch" --file "/path/to/commands-file-in-yaml-format" -r
```

#### Command Line Options:
    "-r, --recursive": watch subdirectories recursively.
    "-p, --path": set the path to the directory to watch for events.
    "-f, --file": set the path to the file that contains the commands to run on each event, check out the `commands.yaml` file to see how this file should look like.
