package main

import "time"

type args struct {
	Path      string        `arg:"-p,--path"`
	File      string        `arg:"-f,--file" help:"path to configuration file (required unless --version is passed)"`
	Recursive bool          `arg:"-r,--recursive"`
	Debug     bool          `arg:"-d,--debug"`
	Debounce  time.Duration `arg:"-b,--debounce" help:"debounce interval (e.g. 400ms)"`
	Version   bool          `arg:"-v,--version" help:"print version information and exit"`
	DryRun    bool          `arg:"--dry-run" help:"log which commands would run without executing them"`
}
