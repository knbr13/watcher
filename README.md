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
    cd watcher
    ```

3. Build the project:

    ```bash
    go build
    ```

4. Run the executable:

    ```bash
    ./watcher
    ```

### Usage

```bash
./watcher -cmd "your-commands" -path "/path/to/watch" -events "specify-events" -r
``l

#### Command Line Options:
    -cmd: Specify the commands to run on events, separated by ';'.
    -path: Set the path to the directory to watch for events.
    -events: Specify the events to watch (write, create, chmod, remove, rename, all), separated by ','.
    -r: Watch subdirectories recursively.
