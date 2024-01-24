# Watcher

Watcher is a tool written in Go that monitors file system events and executes specified commands in response to those events.

## Features

- Watch a directory or file for file system events.
- Run custom commands on different types of events (create, write, chmod, remove, rename).
- Optionally watch subdirectories recursively.

## Getting Started

### Prerequisites

- Go (version 1.13 or later)

### Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/knbr13/watcher.git
    ```

2. Change to the project directory:

    ```bash
    cd your-file-watcher
    ```

3. Build the project:

    ```bash
    go build
    ```

4. Run the executable:

    ```bash
    ./your-file-watcher
    ```

### Usage

```bash
./your-file-watcher -cmd "your-command" -path "/path/to/watch" -events "specify-events" -r
```
#### Command Line Options:
    -cmd: Specify the command to run on events.
    -path: Set the path to the directory to watch for events.
    -events: Specify the events to watch (write, create, chmod, remove, rename, all).
    -r: Watch subdirectories recursively.