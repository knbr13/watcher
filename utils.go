package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func addPathRecursively(root string, watcher *fsnotify.Watcher) error {
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() || slices.Contains(excludedFolders, strings.ToLower(d.Name())) {
			return nil
		}
		return watcher.Add(path)
	})
	return err
}

func parseCommand(cmd string) *exec.Cmd {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return nil
	}
	return exec.Command("sh", "-c", cmd)
}

func wrapCmd(cmd *exec.Cmd, event fsnotify.Event) *exec.Cmd {
	if cmd != nil {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = append(os.Environ(),
			fmt.Sprintf("WATCHER_EVENT=%s", event.Op),
			fmt.Sprintf("WATCHER_PATH=%s", event.Name),
		)
	}
	return cmd
}

func shouldProcess(path string, include, exclude []string) bool {
	// Check exclude first
	if matches(path, exclude) {
		return false
	}

	// If include is empty, everything is included
	if len(include) == 0 {
		return true
	}

	return matches(path, include)
}

func matches(path string, patterns []string) bool {
	for _, p := range patterns {
		if matchesPattern(path, p) {
			return true
		}
	}
	return false
}

func matchesPattern(path, pattern string) bool {
	if pattern == "" {
		return true
	}
	// Try matching the full path
	if m, _ := filepath.Match(pattern, path); m {
		return true
	}
	// Try matching the base name
	if m, _ := filepath.Match(pattern, filepath.Base(path)); m {
		return true
	}
	return false
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
