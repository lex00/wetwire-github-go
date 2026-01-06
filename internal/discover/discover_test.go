package discover

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverer_Discover(t *testing.T) {
	// Create a temporary directory with test files
	tmpDir := t.TempDir()

	// Create a test Go file with workflow and job declarations
	testFile := filepath.Join(tmpDir, "workflows.go")
	content := `package main

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var CI = workflow.Workflow{
	Name: "CI",
	Jobs: []any{Build, Test},
}

var Build = workflow.Job{
	Name:   "build",
	RunsOn: "ubuntu-latest",
}

var Test = workflow.Job{
	Name:   "test",
	RunsOn: "ubuntu-latest",
	Needs:  []any{Build},
}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	d := NewDiscoverer()
	result, err := d.Discover(tmpDir)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	// Check workflows
	if len(result.Workflows) != 1 {
		t.Errorf("len(result.Workflows) = %d, want 1", len(result.Workflows))
	} else {
		w := result.Workflows[0]
		if w.Name != "CI" {
			t.Errorf("workflow.Name = %q, want %q", w.Name, "CI")
		}
		if len(w.Jobs) != 2 {
			t.Errorf("len(workflow.Jobs) = %d, want 2", len(w.Jobs))
		}
	}

	// Check jobs
	if len(result.Jobs) != 2 {
		t.Errorf("len(result.Jobs) = %d, want 2", len(result.Jobs))
	}
}

func TestDiscoverer_SkipsVendor(t *testing.T) {
	tmpDir := t.TempDir()

	// Create vendor directory with a Go file
	vendorDir := filepath.Join(tmpDir, "vendor")
	if err := os.MkdirAll(vendorDir, 0755); err != nil {
		t.Fatal(err)
	}

	vendorFile := filepath.Join(vendorDir, "test.go")
	content := `package vendor

import "github.com/lex00/wetwire-github-go/workflow"

var VendorWorkflow = workflow.Workflow{Name: "Vendor"}
`
	if err := os.WriteFile(vendorFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	d := NewDiscoverer()
	result, err := d.Discover(tmpDir)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	if len(result.Workflows) != 0 {
		t.Errorf("Expected no workflows from vendor, got %d", len(result.Workflows))
	}
}

func TestDiscoverer_SkipsHiddenDirs(t *testing.T) {
	tmpDir := t.TempDir()

	// Create hidden directory
	hiddenDir := filepath.Join(tmpDir, ".hidden")
	if err := os.MkdirAll(hiddenDir, 0755); err != nil {
		t.Fatal(err)
	}

	hiddenFile := filepath.Join(hiddenDir, "test.go")
	content := `package hidden

import "github.com/lex00/wetwire-github-go/workflow"

var HiddenWorkflow = workflow.Workflow{Name: "Hidden"}
`
	if err := os.WriteFile(hiddenFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	d := NewDiscoverer()
	result, err := d.Discover(tmpDir)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	if len(result.Workflows) != 0 {
		t.Errorf("Expected no workflows from hidden dir, got %d", len(result.Workflows))
	}
}

func TestDiscoverer_SkipsTestFiles(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "workflows_test.go")
	content := `package main

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var TestWorkflow = workflow.Workflow{Name: "Test"}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	d := NewDiscoverer()
	result, err := d.Discover(tmpDir)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	if len(result.Workflows) != 0 {
		t.Errorf("Expected no workflows from test files, got %d", len(result.Workflows))
	}
}

func TestDiscoverer_NoWorkflowImport(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "other.go")
	content := `package main

import "fmt"

type Workflow struct {
	Name string
}

var CI = Workflow{Name: "CI"}

func main() {
	fmt.Println("hello")
}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	d := NewDiscoverer()
	result, err := d.Discover(tmpDir)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	// Should not find any workflows since workflow package is not imported
	if len(result.Workflows) != 0 {
		t.Errorf("Expected no workflows without workflow import, got %d", len(result.Workflows))
	}
}

func TestNewDiscoverer(t *testing.T) {
	d := NewDiscoverer()
	if d == nil {
		t.Error("NewDiscoverer() returned nil")
	}
	if d.fset == nil {
		t.Error("NewDiscoverer().fset is nil")
	}
}

func TestIsBuiltinIdent(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"true", true},
		{"false", true},
		{"nil", true},
		{"string", true},
		{"int", true},
		{"bool", true},
		{"workflow", true},
		{"Job", true},
		{"Workflow", true},
		{"Build", false},
		{"Test", false},
		{"CI", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isBuiltinIdent(tt.name)
			if got != tt.want {
				t.Errorf("isBuiltinIdent(%q) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
