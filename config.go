package main

import "gopkg.in/yaml.v3"

type Rule struct {
	Commands   []string `yaml:"commands"`
	OnSuccess  []string `yaml:"on_success"`
	OnFailure  []string `yaml:"on_failure"`
	Pattern    string   `yaml:"pattern"`
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
