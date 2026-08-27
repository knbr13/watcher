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

func TestParseRejectsUnknownFields(t *testing.T) {
	yamlData := `
writ:
  - pattern: "*.go"
    commands: ["go build"]
`
	var cf CommandsFile
	if err := cf.Parse([]byte(yamlData)); err == nil {
		t.Fatal("expected Parse to reject an unknown top-level field, got nil error")
	}
}

func TestParseAcceptsKnownFields(t *testing.T) {
	yamlData := `
write:
  - pattern: "*.go"
    commands: ["go build"]
`
	var cf CommandsFile
	if err := cf.Parse([]byte(yamlData)); err != nil {
		t.Fatalf("Parse: %v", err)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		cf      CommandsFile
		wantErr bool
	}{
		{
			name: "valid config",
			cf: CommandsFile{
				Write: []Rule{{Pattern: "*.go", Commands: []string{"go build"}}},
			},
			wantErr: false,
		},
		{
			name:    "negative debounce",
			cf:      CommandsFile{Debounce: -time.Second},
			wantErr: true,
		},
		{
			name: "empty pattern",
			cf: CommandsFile{
				Write: []Rule{{Pattern: "", Commands: []string{"go build"}}},
			},
			wantErr: true,
		},
		{
			name: "negative timeout",
			cf: CommandsFile{
				Write: []Rule{{Pattern: "*.go", Commands: []string{"go build"}, Timeout: -time.Second}},
			},
			wantErr: true,
		},
		{
			name: "empty commands",
			cf: CommandsFile{
				Write: []Rule{{Pattern: "*.go", Commands: nil}},
			},
			wantErr: true,
		},
		{
			name: "blank command entry",
			cf: CommandsFile{
				Write: []Rule{{Pattern: "*.go", Commands: []string{"  "}}},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cf.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
