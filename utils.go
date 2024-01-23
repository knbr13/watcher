package main

import (
	"io/fs"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func addSubdirectories(root string, watcher *fsnotify.Watcher) error {
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err

		}
		if !d.IsDir() {
			return nil
		}
		if slices.Contains(excludedFolders, strings.ToLower(d.Name())) {
			return nil
		}
		err = watcher.Add(path)
		return err
	})
	return err
}

func parseCommands(cmd string) []*exec.Cmd {
	cmd = strings.TrimSpace(cmd)
	cmds := strings.Split(cmd, ";") // echo hello;echo world // []string{"echo hello     ", "echo world"}
	execCommands := make([]*exec.Cmd, 0, len(cmds))
	for _, cmd := range cmds {
		cmd = strings.TrimSpace(cmd)
		if cmd == "" {
			continue
		}
		execCommands = append(execCommands, exec.Command("sh", "-c", cmd))
	}
	return execCommands
}

var excludedFolders = []string{
	"node_modules",
	"vendor",
	".git",
	".svn",
	".hg",
	".bzr",
	"_vendor",
	"godeps",
	"thirdparty",
	"bin",
	"obj",
	"testdata",
	"examples",
	"tmp",
	"build",
}
