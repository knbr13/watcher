package main

import (
	"gopkg.in/yaml.v3"
)

type args struct {
	Path      string `arg:"-p,--path"`
	File      string `arg:"-f,--file, required"`
	Recursive bool   `arg:"-r,--recursive"`
}

type CommandsFile struct {
	Write  []string `yaml:"write"`
	Chmod  []string `yaml:"chmod"`
	Rename []string `yaml:"rename"`
	Remove []string `yaml:"remove"`
	Create []string `yaml:"create"`
	Common []string `yaml:"common"`
}

func (c *CommandsFile) Parse(data []byte) error {
	return yaml.Unmarshal(data, c)
}
