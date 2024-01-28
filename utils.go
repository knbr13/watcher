package main

import (
	"io/fs"
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
		if !d.IsDir() {
			return nil
		}
		if slices.Contains(excludedFolders, strings.ToLower(d.Name())) {
			return nil
		}
		return watcher.Add(path)
	})
	return err
}

func parseCommands(cmd string) [][]string {
	cmd = strings.TrimSpace(cmd)
	cmds := strings.Split(cmd, ";")
	var res [][]string
	for _, cmd := range cmds {
		cmd = strings.TrimSpace(cmd)
		cmds = strings.Split(cmd, " ")
		if len(cmds) < 1 || cmds[0] == "" {
			continue
		}
		res = append(res, cmds)
	}
	return res
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
