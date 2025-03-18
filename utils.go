package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/fsnotify/fsnotify"
	"github.com/google/shlex"
)

// Add to imports
import (
	"path/filepath"
	"time"
)

func expandVars(cmd string, event fsnotify.Event) string {
	base := filepath.Base(event.Name)
	dir := filepath.Dir(event.Name)
	abs, _ := filepath.Abs(event.Name)

	return os.Expand(cmd, func(key string) string {
		switch key {
		case "FILE":
			return event.Name
		case "FILE_BASE":
			return base
		case "FILE_DIR":
			return dir
		case "FILE_ABS":
			return abs
		case "FILE_EXT":
			return filepath.Ext(event.Name)

		case "EVENT_TYPE":
			return event.Op.String()
		case "EVENT_TIME":
			return time.Now().Format(time.RFC3339)

		case "PWD":
			wd, _ := os.Getwd()
			return wd
		case "TIMESTAMP":
			return fmt.Sprintf("%d", time.Now().Unix())

		default:
			return os.Getenv(key)
		}
	})
}

func parseCommand(cmdTemplate string, event fsnotify.Event) *exec.Cmd {
	expanded := expandVars(cmdTemplate, event)
	expanded = strings.TrimSpace(expanded)
	if expanded == "" {
		return nil
	}

	parts, err := shlex.Split(expanded)
	if err != nil || len(parts) == 0 {
		return nil
	}

	if len(parts) == 1 {
		return exec.Command(parts[0])
	}
	return exec.Command(parts[0], parts[1:]...)
}

func wrapCmd(cmd *exec.Cmd) *exec.Cmd {
	if cmd != nil {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd
}

func matchesPattern(path, pattern string) bool {
	matched, err := doublestar.Match(pattern, path)
	if err != nil {
		fmt.Printf("error matching pattern %q: %s\n", pattern, err.Error())
		return false
	}
	return matched
}

var excludedFolders = []string{
	"node_modules",
	"vendor",
	".git",
	".svn",
	".hg",
	".bzr",
	".vscode",
	"_vendor",
	"godeps",
	"dist",
	"thirdparty",
	"bin",
	"__pycache__",
	".cache",
	"obj",
	"testdata",
	"examples",
	"tmp",
	"build",
}

func validPath(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// fatalf prints a formatted error message to stderr and exits with status code 1
func fatalf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg, args...)
	os.Exit(1)
}
