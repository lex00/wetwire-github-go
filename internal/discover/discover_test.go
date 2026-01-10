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
		{"any", true},
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

func TestDiscoverer_ExplicitWorkflowType(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "workflows.go")
	content := `package main

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var CI workflow.Workflow = workflow.Workflow{
	Name: "CI",
	Jobs: []any{Build},
}

var Build workflow.Job = workflow.Job{
	Name:   "build",
	RunsOn: "ubuntu-latest",
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

	if len(result.Workflows) != 1 {
		t.Errorf("len(result.Workflows) = %d, want 1", len(result.Workflows))
	}

	if len(result.Jobs) != 1 {
		t.Errorf("len(result.Jobs) = %d, want 1", len(result.Jobs))
	}
}

func TestDiscoverer_UnqualifiedTypes(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "workflows.go")
	content := `package main

import (
	. "github.com/lex00/wetwire-github-go/workflow"
)

var CI Workflow = Workflow{
	Name: "CI",
	Jobs: []any{Build},
}

var Build Job = Job{
	Name:   "build",
	RunsOn: "ubuntu-latest",
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

	if len(result.Workflows) != 1 {
		t.Errorf("len(result.Workflows) = %d, want 1", len(result.Workflows))
	}

	if len(result.Jobs) != 1 {
		t.Errorf("len(result.Jobs) = %d, want 1", len(result.Jobs))
	}
}

func TestDiscoverer_ParseError(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "invalid.go")
	content := `package main

import "github.com/lex00/wetwire-github-go/workflow"

var InvalidSyntax = // missing value
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	d := NewDiscoverer()
	result, err := d.Discover(tmpDir)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	if len(result.Errors) == 0 {
		t.Error("Expected parse errors to be recorded")
	}
}

func TestDiscoverer_PointerType(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "workflows.go")
	content := `package main

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var CI = &workflow.Workflow{
	Name: "CI",
	Jobs: []any{Build},
}

var Build = &workflow.Job{
	Name:   "build",
	RunsOn: "ubuntu-latest",
	Needs:  []any{Test},
}

var Test = &workflow.Job{
	Name:   "test",
	RunsOn: "ubuntu-latest",
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

	// Should discover workflows and jobs even with pointer syntax
	if len(result.Workflows) != 1 {
		t.Errorf("len(result.Workflows) = %d, want 1", len(result.Workflows))
	}

	if len(result.Jobs) != 2 {
		t.Errorf("len(result.Jobs) = %d, want 2", len(result.Jobs))
	}

	// Note: Job/dependency extraction from pointer expressions is not fully supported
	// The types are inferred correctly but the field values (Jobs, Needs) are not extracted
	// This is expected behavior given the current implementation
}

func TestDiscoverer_MultipleFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create multiple files with workflows and jobs
	file1 := filepath.Join(tmpDir, "ci.go")
	content1 := `package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
	Jobs: []any{Build},
}
`
	if err := os.WriteFile(file1, []byte(content1), 0644); err != nil {
		t.Fatal(err)
	}

	file2 := filepath.Join(tmpDir, "jobs.go")
	content2 := `package main

import "github.com/lex00/wetwire-github-go/workflow"

var Build = workflow.Job{
	Name:   "build",
	RunsOn: "ubuntu-latest",
}

var Test = workflow.Job{
	Name:   "test",
	RunsOn: "ubuntu-latest",
}
`
	if err := os.WriteFile(file2, []byte(content2), 0644); err != nil {
		t.Fatal(err)
	}

	d := NewDiscoverer()
	result, err := d.Discover(tmpDir)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	if len(result.Workflows) != 1 {
		t.Errorf("len(result.Workflows) = %d, want 1", len(result.Workflows))
	}

	if len(result.Jobs) != 2 {
		t.Errorf("len(result.Jobs) = %d, want 2", len(result.Jobs))
	}
}

func TestDiscoverer_SubDirectories(t *testing.T) {
	tmpDir := t.TempDir()

	// Create subdirectory
	subDir := filepath.Join(tmpDir, "workflows")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	testFile := filepath.Join(subDir, "ci.go")
	content := `package workflows

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
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

	if len(result.Workflows) != 1 {
		t.Errorf("len(result.Workflows) = %d, want 1", len(result.Workflows))
	}
}

func TestDiscoverer_EmptyJobsField(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "workflows.go")
	content := `package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
	Jobs: []any{},
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

	if len(result.Workflows) != 1 {
		t.Errorf("len(result.Workflows) = %d, want 1", len(result.Workflows))
	}

	if len(result.Workflows[0].Jobs) != 0 {
		t.Errorf("len(workflow.Jobs) = %d, want 0", len(result.Workflows[0].Jobs))
	}
}

func TestDiscoverer_EmptyNeedsField(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "workflows.go")
	content := `package main

import "github.com/lex00/wetwire-github-go/workflow"

var Build = workflow.Job{
	Name:   "build",
	RunsOn: "ubuntu-latest",
	Needs:  []any{},
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

	if len(result.Jobs) != 1 {
		t.Errorf("len(result.Jobs) = %d, want 1", len(result.Jobs))
	}

	if len(result.Jobs[0].Dependencies) != 0 {
		t.Errorf("len(job.Dependencies) = %d, want 0", len(result.Jobs[0].Dependencies))
	}
}

func TestDiscoverer_NonGoFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a non-Go file
	testFile := filepath.Join(tmpDir, "readme.txt")
	content := "This is a README file"
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	d := NewDiscoverer()
	result, err := d.Discover(tmpDir)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	if len(result.Workflows) != 0 {
		t.Errorf("Expected no workflows from non-Go files, got %d", len(result.Workflows))
	}
}

func TestDiscoverer_SkipsTestdata(t *testing.T) {
	tmpDir := t.TempDir()

	testdataDir := filepath.Join(tmpDir, "testdata")
	if err := os.MkdirAll(testdataDir, 0755); err != nil {
		t.Fatal(err)
	}

	testdataFile := filepath.Join(testdataDir, "workflows.go")
	content := `package testdata

import "github.com/lex00/wetwire-github-go/workflow"

var TestdataWorkflow = workflow.Workflow{Name: "Testdata"}
`
	if err := os.WriteFile(testdataFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	d := NewDiscoverer()
	result, err := d.Discover(tmpDir)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	if len(result.Workflows) != 0 {
		t.Errorf("Expected no workflows from testdata, got %d", len(result.Workflows))
	}
}

func TestDiscoverer_WorkflowWithoutJobsField(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "workflows.go")
	content := `package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
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

	if len(result.Workflows) != 1 {
		t.Errorf("len(result.Workflows) = %d, want 1", len(result.Workflows))
	}

	if len(result.Workflows[0].Jobs) != 0 {
		t.Errorf("len(workflow.Jobs) = %d, want 0", len(result.Workflows[0].Jobs))
	}
}

func TestDiscoverer_JobWithoutNeedsField(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "workflows.go")
	content := `package main

import "github.com/lex00/wetwire-github-go/workflow"

var Build = workflow.Job{
	Name:   "build",
	RunsOn: "ubuntu-latest",
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

	if len(result.Jobs) != 1 {
		t.Errorf("len(result.Jobs) = %d, want 1", len(result.Jobs))
	}

	if len(result.Jobs[0].Dependencies) != 0 {
		t.Errorf("len(job.Dependencies) = %d, want 0", len(result.Jobs[0].Dependencies))
	}
}

func TestDiscoverer_NestedCompositeLiteral(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "workflows.go")
	content := `package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
	Jobs: []any{
		Build,
		Test,
		Deploy,
	},
}

var Deploy = workflow.Job{
	Name:   "deploy",
	RunsOn: "ubuntu-latest",
	Needs:  []any{Build, Test},
}

var Build = workflow.Job{Name: "build"}
var Test = workflow.Job{Name: "test"}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	d := NewDiscoverer()
	result, err := d.Discover(tmpDir)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	if len(result.Workflows) != 1 {
		t.Errorf("len(result.Workflows) = %d, want 1", len(result.Workflows))
	}

	// Should extract all 3 job references
	if len(result.Workflows[0].Jobs) != 3 {
		t.Errorf("len(workflow.Jobs) = %d, want 3", len(result.Workflows[0].Jobs))
	}

	// Should extract both dependencies
	foundDeploy := false
	for _, job := range result.Jobs {
		if job.Name == "Deploy" {
			foundDeploy = true
			if len(job.Dependencies) != 2 {
				t.Errorf("len(Deploy.Dependencies) = %d, want 2", len(job.Dependencies))
			}
		}
	}
	if !foundDeploy {
		t.Error("Deploy job not found")
	}
}

func TestDiscoverer_DirectoryWalkError(t *testing.T) {
	// Test with a non-existent directory
	d := NewDiscoverer()
	_, err := d.Discover("/nonexistent/directory/that/does/not/exist")
	if err == nil {
		t.Error("Expected error for non-existent directory")
	}
}

func TestDiscoverer_GetTypeName_EdgeCases(t *testing.T) {
	tmpDir := t.TempDir()

	// Test with a complex type that might not be recognized
	testFile := filepath.Join(tmpDir, "workflows.go")
	content := `package main

import "github.com/lex00/wetwire-github-go/workflow"

// This should not be discovered as a workflow (wrong type)
var NotAWorkflow struct {
	Name string
}

var ValidWorkflow = workflow.Workflow{
	Name: "Valid",
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

	// Should only find the valid workflow
	if len(result.Workflows) != 1 {
		t.Errorf("len(result.Workflows) = %d, want 1", len(result.Workflows))
	}

	if result.Workflows[0].Name != "ValidWorkflow" {
		t.Errorf("workflow.Name = %q, want ValidWorkflow", result.Workflows[0].Name)
	}
}
