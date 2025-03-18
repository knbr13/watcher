package main

type args struct {
	Path      string `arg:"-p,--path"`
	File      string `arg:"-f,--file, required"`
	Recursive bool   `arg:"-r,--recursive"`
	Debug     bool   `arg:"-d,--debug"`
}
