package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"text/template"
	"time"

	"github.com/bmatcuk/doublestar/v4"
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

func parseCommand(cmdStr string) *exec.Cmd {
	cmdStr = strings.TrimSpace(cmdStr)
	if cmdStr == "" {
		return nil
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		shell := os.Getenv("ComSpec")
		if shell == "" {
			shell = "cmd.exe"
		}
		cmd = exec.Command(shell, "/C", cmdStr)
	} else {
		cmd = exec.Command("sh", "-c", cmdStr)
	}
	setProcAttrs(cmd)
	return cmd
}

func wrapCmd(cmd *exec.Cmd, data EventData) *exec.Cmd {
	if cmd != nil {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = append(os.Environ(),
			fmt.Sprintf("WATCHER_EVENT=%s", data.Op),
			fmt.Sprintf("WATCHER_PATH=%s", data.Path),
			fmt.Sprintf("FILE=%s", data.Path),
			fmt.Sprintf("FILE_BASE=%s", data.Base),
			fmt.Sprintf("FILE_DIR=%s", data.Dir),
			fmt.Sprintf("FILE_EXT=%s", data.Ext),
			fmt.Sprintf("EVENT_TYPE=%s", data.Op),
			fmt.Sprintf("EVENT_TIME=%s", data.Time),
			fmt.Sprintf("TIMESTAMP=%d", data.Timestamp),
			fmt.Sprintf("PWD=%s", data.PWD),
		)
	}
	return cmd
}

type EventData struct {
	Path      string
	Base      string
	Dir       string
	Ext       string
	Op        string
	Time      string
	Timestamp int64
	PWD       string
}

func getEventData(event fsnotify.Event) EventData {
	absPath, err := filepath.Abs(event.Name)
	if err != nil {
		absPath = event.Name
	}
	now := time.Now()
	pwd, _ := os.Getwd()
	return EventData{
		Path:      absPath,
		Base:      filepath.Base(event.Name),
		Dir:       filepath.Dir(event.Name),
		Ext:       filepath.Ext(event.Name),
		Op:        strings.ToUpper(event.Op.String()),
		Time:      now.Format(time.RFC3339),
		Timestamp: now.Unix(),
		PWD:       pwd,
	}
}

// shellQuote quotes s so it is safe to splice into the shell command
// selected by parseCommand as a single literal argument, regardless of what
// characters s contains (e.g. a filename crafted to break out of the
// command). Intended for use as the "quote" template function, e.g.
// {{.Base | quote}}.
func shellQuote(s string) string {
	if runtime.GOOS == "windows" {
		// cmd.exe has no fully safe quoting story (e.g. "%...%" is still
		// expanded inside double quotes), but this neutralizes the common
		// case of quotes/spaces/metacharacters in a file name.
		return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
	}
	return `'` + strings.ReplaceAll(s, `'`, `'\''`) + `'`
}

func expandTemplate(text string, data EventData) string {
	tmpl, err := template.New("cmd").Funcs(template.FuncMap{"quote": shellQuote}).Parse(text)
	if err != nil {
		fmt.Fprintf(logger, "watcher: error parsing template %q: %s\n", text, err.Error())
		return text
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		fmt.Fprintf(logger, "watcher: error executing template %q: %s\n", text, err.Error())
		return text
	}
	return buf.String()
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

// matches reports whether path is matched by patterns, treating a leading
// "!" on any pattern as an exception that always wins for that path (e.g.
// ["**/*.go", "!**/testdata/**"] matches every .go file except those under
// testdata). A patterns list containing only negations is treated as
// "everything except these".
func matches(path string, patterns []string) bool {
	if len(patterns) == 0 {
		return false
	}

	hasPositive := false
	matched := false
	excluded := false
	for _, p := range patterns {
		if negated, ok := strings.CutPrefix(p, "!"); ok {
			if globMatch(path, negated) {
				excluded = true
			}
			continue
		}
		hasPositive = true
		if globMatch(path, p) {
			matched = true
		}
	}
	if !hasPositive {
		matched = true
	}
	return matched && !excluded
}

// matchesPattern matches a single rule pattern against path, honoring an
// optional leading "!" to negate the match.
func matchesPattern(path, pattern string) bool {
	if pattern == "" {
		return true
	}
	if negated, ok := strings.CutPrefix(pattern, "!"); ok {
		return !globMatch(path, negated)
	}
	return globMatch(path, pattern)
}

// globMatch matches path (or its base name) against a doublestar pattern,
// which supports "**", "?", character classes, and "{a,b}" alternatives.
func globMatch(path, pattern string) bool {
	if pattern == "" {
		return true
	}
	path = filepath.ToSlash(path)
	pattern = filepath.ToSlash(pattern)
	if m, _ := doublestar.Match(pattern, path); m {
		return true
	}
	if m, _ := doublestar.Match(pattern, filepath.Base(path)); m {
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
