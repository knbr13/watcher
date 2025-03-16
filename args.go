package main

import (
	"gopkg.in/yaml.v3"
)

type args struct {
	Path      string `arg:"-p,--path"`
	File      string `arg:"-f,--file, required"`
	Recursive bool   `arg:"-r,--recursive"`
}

type Rule struct {
	Pattern    string   `yaml:"pattern"`
	Commands   []string `yaml:"commands"`
	Sequential bool     `yaml:"sequential"`
}

type CommandsFile struct {
	Write  []Rule `yaml:"write"`
	Chmod  []Rule `yaml:"chmod"`
	Rename []Rule `yaml:"rename"`
	Remove []Rule `yaml:"remove"`
	Create []Rule `yaml:"create"`
	Common []Rule `yaml:"common"`
}

func (c *CommandsFile) Parse(data []byte) error {
	return yaml.Unmarshal(data, c)
}
