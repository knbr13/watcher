package main

import (
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func TestConfigDebounce(t *testing.T) {
	yamlData := `
debounce: 500ms
include:
  - "*.go"
`
	var cf CommandsFile
	err := yaml.Unmarshal([]byte(yamlData), &cf)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if cf.Debounce != 500*time.Millisecond {
		t.Errorf("Expected debounce 500ms, got %v", cf.Debounce)
	}
}

func TestConfigFull(t *testing.T) {
	yamlData := `
debounce: 1s
include: ["*.go"]
exclude: ["vendor/*"]
write:
  - pattern: "*.go"
    commands: ["go build"]
    timeout: 5s
`
	var cf CommandsFile
	err := yaml.Unmarshal([]byte(yamlData), &cf)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if cf.Debounce != time.Second {
		t.Errorf("Expected 1s, got %v", cf.Debounce)
	}
	if len(cf.Include) != 1 || cf.Include[0] != "*.go" {
		t.Errorf("Include mismatch")
	}
	if len(cf.Write) != 1 || cf.Write[0].Timeout != 5*time.Second {
		t.Errorf("Write rule mismatch")
	}
}
