package main

import (
	"runtime"
	"strings"
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
		{
			name:     "include with negation carves out an exception",
			path:     "src/testdata/main.go",
			include:  []string{"**/*.go", "!**/testdata/**"},
			exclude:  nil,
			expected: false,
		},
		{
			name:     "include with negation still matches everything else",
			path:     "src/main.go",
			include:  []string{"**/*.go", "!**/testdata/**"},
			exclude:  nil,
			expected: true,
		},
		{
			name:     "exclude with negation lets a specific file through",
			path:     "vendor/keep.go",
			include:  nil,
			exclude:  []string{"vendor/**", "!vendor/keep.go"},
			expected: true,
		},
		{
			name:     "brace expansion in include",
			path:     "uploads/photo.png",
			include:  []string{"uploads/*.{jpg,png}"},
			exclude:  nil,
			expected: true,
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
		{"doublestar match nested dirs", "src/pkg/deep/main.go", "**/*.go", true},
		{"doublestar match at root", "main.go", "**/*.go", true},
		{"brace alternatives match jpg", "photo.jpg", "*.{jpg,png}", true},
		{"brace alternatives match png", "photo.png", "*.{jpg,png}", true},
		{"brace alternatives mismatch", "photo.gif", "*.{jpg,png}", false},
		{"negated pattern excludes match", "main.go", "!*.go", false},
		{"negated pattern lets non-match through", "main.txt", "!*.go", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchesPattern(tt.path, tt.pattern); got != tt.expected {
				t.Errorf("matchesPattern() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestExpandTemplate(t *testing.T) {
	data := EventData{
		Path:      "/home/user/project/main.go",
		Base:      "main.go",
		Dir:       "/home/user/project",
		Ext:       ".go",
		Op:        "WRITE",
		Time:      "2023-01-01T00:00:00Z",
		Timestamp: 1672531200,
		PWD:       "/home/user/project",
	}

	tests := []struct {
		name     string
		template string
		expected string
	}{
		{
			name:     "no placeholders",
			template: "echo hello",
			expected: "echo hello",
		},
		{
			name:     "path placeholder",
			template: "echo {{.Path}}",
			expected: "echo /home/user/project/main.go",
		},
		{
			name:     "multiple placeholders",
			template: "echo {{.Base}} in {{.Dir}}",
			expected: "echo main.go in /home/user/project",
		},
		{
			name:     "operation placeholder",
			template: "Event: {{.Op}}",
			expected: "Event: WRITE",
		},
		{
			name:     "timestamp placeholder",
			template: "at {{.Timestamp}}",
			expected: "at 1672531200",
		},
		{
			name:     "invalid template",
			template: "echo {{.Invalid}}",
			expected: "echo {{.Invalid}}", // Should return original on error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := expandTemplate(tt.template, data); got != tt.expected {
				t.Errorf("expandTemplate() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestShellQuote(t *testing.T) {
	tests := []string{
		"plain",
		"with space",
		"it's got a quote",
		`"double quoted"`,
		"; rm -rf / #",
		"$(whoami)",
		"back`tick`",
	}
	for _, s := range tests {
		t.Run(s, func(t *testing.T) {
			got := shellQuote(s)
			var want string
			if runtime.GOOS == "windows" {
				want = `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
			} else {
				want = `'` + strings.ReplaceAll(s, `'`, `'\''`) + `'`
			}
			if got != want {
				t.Errorf("shellQuote(%q) = %q, want %q", s, got, want)
			}
		})
	}
}

func TestBuildExcludeDirs(t *testing.T) {
	merged := buildExcludeDirs([]string{"Custom", "ANOTHER"})

	for _, want := range []string{"node_modules", "vendor", "custom", "another"} {
		found := false
		for _, got := range merged {
			if got == want {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("buildExcludeDirs() missing %q, got %v", want, merged)
		}
	}
}

func TestExpandTemplateQuoteFunc(t *testing.T) {
	data := EventData{Base: "; rm -rf ~ #"}
	got := expandTemplate("echo {{.Base | quote}}", data)
	want := "echo " + shellQuote(data.Base)
	if got != want {
		t.Errorf("expandTemplate with quote = %q, want %q", got, want)
	}
}
