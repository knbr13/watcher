package main

import (
	"testing"
)

func TestShouldProcess(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		include  []string
		exclude  []string
		expected bool
	}{
		{
			name:     "no filters",
			path:     "main.go",
			include:  nil,
			exclude:  nil,
			expected: true,
		},
		{
			name:     "include match extension",
			path:     "main.go",
			include:  []string{"*.go"},
			exclude:  nil,
			expected: true,
		},
		{
			name:     "include mismatch extension",
			path:     "main.go",
			include:  []string{"*.txt"},
			exclude:  nil,
			expected: false,
		},
		{
			name:     "exclude match",
			path:     "main.go",
			include:  nil,
			exclude:  []string{"main.go"},
			expected: false,
		},
		{
			name:     "include and exclude match (exclude wins)",
			path:     "main.go",
			include:  []string{"*.go"},
			exclude:  []string{"main.go"},
			expected: false,
		},
		{
			name:     "include match in subdirectory",
			path:     "src/main.go",
			include:  []string{"src/*.go"},
			exclude:  nil,
			expected: true,
		},
		{
			name:     "include match by base name in subdirectory",
			path:     "src/main.go",
			include:  []string{"*.go"},
			exclude:  nil,
			expected: true,
		},
		{
			name:     "exclude subdirectory",
			path:     "vendor/lib.go",
			include:  nil,
			exclude:  []string{"vendor/*"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldProcess(tt.path, tt.include, tt.exclude); got != tt.expected {
				t.Errorf("shouldProcess() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestMatchesPattern(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		pattern  string
		expected bool
	}{
		{"empty pattern", "main.go", "", true},
		{"exact match", "main.go", "main.go", true},
		{"extension match", "main.go", "*.go", true},
		{"extension mismatch", "main.go", "*.txt", false},
		{"base name match in subdir", "src/main.go", "main.go", true},
		{"subdir match", "src/main.go", "src/*.go", true},
		{"subdir mismatch", "pkg/main.go", "src/*.go", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchesPattern(tt.path, tt.pattern); got != tt.expected {
				t.Errorf("matchesPattern() = %v, want %v", got, tt.expected)
			}
		})
	}
}
