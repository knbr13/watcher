package main

import (
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
	parts := strings.Split(cmd, " ")
	if len(parts) < 2 {
		return nil
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
