package main

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Rule struct {
	Commands   []string      `yaml:"commands"`
	OnSuccess  []string      `yaml:"on_success"`
	OnFailure  []string      `yaml:"on_failure"`
	Pattern    string        `yaml:"pattern"`
	Sequential bool          `yaml:"sequential"`
	Timeout    time.Duration `yaml:"timeout"`
}

type CommandsFile struct {
	Debounce    time.Duration `yaml:"debounce"`
	Include     []string      `yaml:"include"`
	Watch       []string      `yaml:"watch"`
	Exclude     []string      `yaml:"exclude"`
	ExcludeDirs []string      `yaml:"exclude_dirs"`
	Write       []Rule        `yaml:"write"`
	Chmod       []Rule        `yaml:"chmod"`
	Rename      []Rule        `yaml:"rename"`
	Remove      []Rule        `yaml:"remove"`
	Create      []Rule        `yaml:"create"`
	Common      []Rule        `yaml:"common"`
}

// Parse decodes a config file, rejecting unknown top-level or rule fields so
// typos (e.g. "patern" instead of "pattern") fail loudly instead of being
// silently ignored.
func (c *CommandsFile) Parse(data []byte) error {
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	if err := dec.Decode(c); err != nil {
		return fmt.Errorf("parsing config: %w", err)
	}
	return nil
}

// Validate checks the decoded config for values that would otherwise fail
// confusingly (or silently do nothing) at runtime.
func (c *CommandsFile) Validate() error {
	if c.Debounce < 0 {
		return fmt.Errorf("debounce must not be negative, got %v", c.Debounce)
	}

	ruleSets := []struct {
		name  string
		rules []Rule
	}{
		{"write", c.Write},
		{"chmod", c.Chmod},
		{"rename", c.Rename},
		{"remove", c.Remove},
		{"create", c.Create},
		{"common", c.Common},
	}
	for _, set := range ruleSets {
		for i, rule := range set.rules {
			if err := rule.validate(); err != nil {
				return fmt.Errorf("%s[%d]: %w", set.name, i, err)
			}
		}
	}
	return nil
}

func (r Rule) validate() error {
	if strings.TrimSpace(r.Pattern) == "" {
		return fmt.Errorf("pattern must not be empty")
	}
	if r.Timeout < 0 {
		return fmt.Errorf("timeout must not be negative, got %v", r.Timeout)
	}
	if len(r.Commands) == 0 {
		return fmt.Errorf("commands must not be empty")
	}
	for _, cmd := range r.Commands {
		if strings.TrimSpace(cmd) == "" {
			return fmt.Errorf("commands must not contain an empty entry")
		}
	}
	return nil
}
