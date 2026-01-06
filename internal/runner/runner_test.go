package runner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lex00/wetwire-github-go/internal/discover"
)

func TestNewRunner(t *testing.T) {
	r := NewRunner()
	if r == nil {
		t.Error("NewRunner() returned nil")
	}
	if r.TempDir == "" {
		t.Error("NewRunner().TempDir is empty")
	}
}

func TestRunner_parseGoMod(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23

require (
	github.com/some/dep v1.0.0
)
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()
	modulePath, err := r.parseGoMod(tmpDir)
	if err != nil {
		t.Fatalf("parseGoMod() error = %v", err)
	}

	if modulePath != "github.com/example/test" {
		t.Errorf("parseGoMod() = %q, want %q", modulePath, "github.com/example/test")
	}
}

func TestRunner_parseGoMod_NotFound(t *testing.T) {
	tmpDir := t.TempDir()

	r := NewRunner()
	_, err := r.parseGoMod(tmpDir)
	if err == nil {
		t.Error("parseGoMod() expected error for missing go.mod")
	}
}

func TestRunner_parseGoMod_NoModule(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()
	_, err := r.parseGoMod(tmpDir)
	if err == nil {
		t.Error("parseGoMod() expected error for missing module directive")
	}
}

func TestRunner_generateGoMod(t *testing.T) {
	r := NewRunner()
	result := r.generateGoMod("github.com/example/test", "/path/to/project")

	if !strings.Contains(result, "module wetwire-extract") {
		t.Error("generateGoMod() missing module directive")
	}

	if !strings.Contains(result, "require github.com/example/test") {
		t.Error("generateGoMod() missing require directive")
	}

	if !strings.Contains(result, "replace github.com/example/test =>") {
		t.Error("generateGoMod() missing replace directive")
	}
}

func TestRunner_getPackagePath(t *testing.T) {
	r := NewRunner()
	baseDir := "/project"

	tests := []struct {
		modulePath string
		baseDir    string
		file       string
		want       string
	}{
		{"github.com/example/test", baseDir, "/project/workflows.go", "github.com/example/test"},
		{"github.com/example/test", baseDir, "/project/pkg/workflows.go", "github.com/example/test/pkg"},
		{"github.com/example/test", baseDir, "/project/internal/ci/workflows.go", "github.com/example/test/internal/ci"},
	}

	for _, tt := range tests {
		got := r.getPackagePath(tt.modulePath, tt.baseDir, tt.file)
		if got != tt.want {
			t.Errorf("getPackagePath(%q, %q, %q) = %q, want %q", tt.modulePath, tt.baseDir, tt.file, got, tt.want)
		}
	}
}

func TestRunner_pkgAlias(t *testing.T) {
	r := NewRunner()

	tests := []struct {
		input string
		want  string
	}{
		{"github.com/example/test", "test"},
		{"github.com/example/my-pkg", "my_pkg"},
		{"github.com/org/repo/internal/ci", "ci"},
	}

	for _, tt := range tests {
		got := r.pkgAlias(tt.input)
		if got != tt.want {
			t.Errorf("pkgAlias(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestRunner_generateProgram(t *testing.T) {
	r := NewRunner()
	baseDir := "/project"

	discovered := &discover.DiscoveryResult{
		Workflows: []discover.DiscoveredWorkflow{
			{Name: "CI", File: "/project/workflows.go", Line: 10},
		},
		Jobs: []discover.DiscoveredJob{
			{Name: "Build", File: "/project/jobs.go", Line: 5},
		},
	}

	program, err := r.generateProgram("github.com/example/test", baseDir, discovered)
	if err != nil {
		t.Fatalf("generateProgram() error = %v", err)
	}

	expectedStrings := []string{
		"package main",
		"encoding/json",
		"ExtractionResult",
		"ExtractedWorkflow",
		"ExtractedJob",
		"toMap",
		"json.Marshal",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(program, expected) {
			t.Errorf("generateProgram() missing %q\n\nGenerated:\n%s", expected, program)
		}
	}
}

func TestRunner_ExtractValues_Empty(t *testing.T) {
	r := NewRunner()

	discovered := &discover.DiscoveryResult{
		Workflows: []discover.DiscoveredWorkflow{},
		Jobs:      []discover.DiscoveredJob{},
	}

	result, err := r.ExtractValues(".", discovered)
	if err != nil {
		t.Fatalf("ExtractValues() error = %v", err)
	}

	if len(result.Workflows) != 0 {
		t.Errorf("len(result.Workflows) = %d, want 0", len(result.Workflows))
	}

	if len(result.Jobs) != 0 {
		t.Errorf("len(result.Jobs) = %d, want 0", len(result.Jobs))
	}
}

func TestFindGoBinary(t *testing.T) {
	path, err := FindGoBinary()
	if err != nil {
		t.Skipf("Go binary not found, skipping: %v", err)
	}

	if path == "" {
		t.Error("FindGoBinary() returned empty path")
	}
}
