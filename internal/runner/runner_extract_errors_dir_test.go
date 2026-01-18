package runner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lex00/wetwire-github-go/internal/discover"
)

// Test ExtractValues with invalid directory
func TestRunner_ExtractValues_InvalidDir(t *testing.T) {
	r := NewRunner()

	discovered := &discover.DiscoveryResult{
		Workflows: []discover.DiscoveredWorkflow{
			{Name: "CI", File: "/nonexistent/workflows.go", Line: 10},
		},
	}

	_, err := r.ExtractValues("/nonexistent/directory", discovered)
	if err == nil {
		t.Error("ExtractValues() expected error for nonexistent directory")
	}
}

// Test ExtractDependabot error path - invalid directory
func TestRunner_ExtractDependabot_InvalidDir(t *testing.T) {
	r := NewRunner()

	discovered := &discover.DependabotDiscoveryResult{
		Configs: []discover.DiscoveredDependabot{
			{Name: "Config", File: "/nonexistent/dependabot.go", Line: 10},
		},
	}

	_, err := r.ExtractDependabot("/nonexistent/directory", discovered)
	if err == nil {
		t.Error("ExtractDependabot() expected error for nonexistent directory")
	}
}

// Test ExtractIssueTemplates error path - invalid directory
func TestRunner_ExtractIssueTemplates_InvalidDir(t *testing.T) {
	r := NewRunner()

	discovered := &discover.IssueTemplateDiscoveryResult{
		Templates: []discover.DiscoveredIssueTemplate{
			{Name: "Template", File: "/nonexistent/template.go", Line: 10},
		},
	}

	_, err := r.ExtractIssueTemplates("/nonexistent/directory", discovered)
	if err == nil {
		t.Error("ExtractIssueTemplates() expected error for nonexistent directory")
	}
}

// Test ExtractDiscussionTemplates error path - invalid directory
func TestRunner_ExtractDiscussionTemplates_InvalidDir(t *testing.T) {
	r := NewRunner()

	discovered := &discover.DiscussionTemplateDiscoveryResult{
		Templates: []discover.DiscoveredDiscussionTemplate{
			{Name: "Template", File: "/nonexistent/template.go", Line: 10},
		},
	}

	_, err := r.ExtractDiscussionTemplates("/nonexistent/directory", discovered)
	if err == nil {
		t.Error("ExtractDiscussionTemplates() expected error for nonexistent directory")
	}
}

// Test ExtractPRTemplates error path - invalid directory
func TestRunner_ExtractPRTemplates_InvalidDir(t *testing.T) {
	r := NewRunner()

	discovered := &discover.PRTemplateDiscoveryResult{
		Templates: []discover.DiscoveredPRTemplate{
			{Name: "Template", File: "/nonexistent/template.go", Line: 10},
		},
	}

	_, err := r.ExtractPRTemplates("/nonexistent/directory", discovered)
	if err == nil {
		t.Error("ExtractPRTemplates() expected error for nonexistent directory")
	}
}

// Test ExtractCodeowners error path - invalid directory
func TestRunner_ExtractCodeowners_InvalidDir(t *testing.T) {
	r := NewRunner()

	discovered := &discover.CodeownersDiscoveryResult{
		Configs: []discover.DiscoveredCodeowners{
			{Name: "Config", File: "/nonexistent/codeowners.go", Line: 10},
		},
	}

	_, err := r.ExtractCodeowners("/nonexistent/directory", discovered)
	if err == nil {
		t.Error("ExtractCodeowners() expected error for nonexistent directory")
	}
}

// Test ExtractValues with invalid TempDir
func TestRunner_ExtractValues_InvalidTempDir(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a Go file for the workflow
	workflowFile := `package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{Name: "CI"}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "workflows.go"), []byte(workflowFile), 0644); err != nil {
		t.Fatal(err)
	}

	r := &Runner{
		TempDir: "/nonexistent/temp/dir",
		GoPath:  "go",
	}

	discovered := &discover.DiscoveryResult{
		Workflows: []discover.DiscoveredWorkflow{
			{Name: "CI", File: filepath.Join(tmpDir, "workflows.go"), Line: 5},
		},
	}

	_, err := r.ExtractValues(tmpDir, discovered)
	if err == nil {
		t.Error("ExtractValues() expected error for invalid TempDir")
	}
	if !strings.Contains(err.Error(), "creating temp dir") {
		t.Errorf("Expected 'creating temp dir' error, got: %v", err)
	}
}

// Test ExtractDependabot with invalid TempDir
func TestRunner_ExtractDependabot_InvalidTempDir(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := &Runner{
		TempDir: "/nonexistent/temp/dir",
		GoPath:  "go",
	}

	discovered := &discover.DependabotDiscoveryResult{
		Configs: []discover.DiscoveredDependabot{
			{Name: "Config", File: filepath.Join(tmpDir, "dependabot.go"), Line: 5},
		},
	}

	_, err := r.ExtractDependabot(tmpDir, discovered)
	if err == nil {
		t.Error("ExtractDependabot() expected error for invalid TempDir")
	}
	if !strings.Contains(err.Error(), "creating temp dir") {
		t.Errorf("Expected 'creating temp dir' error, got: %v", err)
	}
}

// Test ExtractIssueTemplates with invalid TempDir
func TestRunner_ExtractIssueTemplates_InvalidTempDir(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := &Runner{
		TempDir: "/nonexistent/temp/dir",
		GoPath:  "go",
	}

	discovered := &discover.IssueTemplateDiscoveryResult{
		Templates: []discover.DiscoveredIssueTemplate{
			{Name: "Template", File: filepath.Join(tmpDir, "template.go"), Line: 5},
		},
	}

	_, err := r.ExtractIssueTemplates(tmpDir, discovered)
	if err == nil {
		t.Error("ExtractIssueTemplates() expected error for invalid TempDir")
	}
	if !strings.Contains(err.Error(), "creating temp dir") {
		t.Errorf("Expected 'creating temp dir' error, got: %v", err)
	}
}

// Test ExtractDiscussionTemplates with invalid TempDir
func TestRunner_ExtractDiscussionTemplates_InvalidTempDir(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := &Runner{
		TempDir: "/nonexistent/temp/dir",
		GoPath:  "go",
	}

	discovered := &discover.DiscussionTemplateDiscoveryResult{
		Templates: []discover.DiscoveredDiscussionTemplate{
			{Name: "Template", File: filepath.Join(tmpDir, "template.go"), Line: 5},
		},
	}

	_, err := r.ExtractDiscussionTemplates(tmpDir, discovered)
	if err == nil {
		t.Error("ExtractDiscussionTemplates() expected error for invalid TempDir")
	}
	if !strings.Contains(err.Error(), "creating temp dir") {
		t.Errorf("Expected 'creating temp dir' error, got: %v", err)
	}
}

// Test ExtractPRTemplates with invalid TempDir
func TestRunner_ExtractPRTemplates_InvalidTempDir(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := &Runner{
		TempDir: "/nonexistent/temp/dir",
		GoPath:  "go",
	}

	discovered := &discover.PRTemplateDiscoveryResult{
		Templates: []discover.DiscoveredPRTemplate{
			{Name: "Template", File: filepath.Join(tmpDir, "template.go"), Line: 5},
		},
	}

	_, err := r.ExtractPRTemplates(tmpDir, discovered)
	if err == nil {
		t.Error("ExtractPRTemplates() expected error for invalid TempDir")
	}
	if !strings.Contains(err.Error(), "creating temp dir") {
		t.Errorf("Expected 'creating temp dir' error, got: %v", err)
	}
}

// Test ExtractCodeowners with invalid TempDir
func TestRunner_ExtractCodeowners_InvalidTempDir(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := &Runner{
		TempDir: "/nonexistent/temp/dir",
		GoPath:  "go",
	}

	discovered := &discover.CodeownersDiscoveryResult{
		Configs: []discover.DiscoveredCodeowners{
			{Name: "Config", File: filepath.Join(tmpDir, "codeowners.go"), Line: 5},
		},
	}

	_, err := r.ExtractCodeowners(tmpDir, discovered)
	if err == nil {
		t.Error("ExtractCodeowners() expected error for invalid TempDir")
	}
	if !strings.Contains(err.Error(), "creating temp dir") {
		t.Errorf("Expected 'creating temp dir' error, got: %v", err)
	}
}

// Test ExtractValues with missing go.mod in directory with workflows
func TestRunner_ExtractValues_MissingGoMod(t *testing.T) {
	tmpDir := t.TempDir()

	// No go.mod file

	r := NewRunner()

	discovered := &discover.DiscoveryResult{
		Workflows: []discover.DiscoveredWorkflow{
			{Name: "CI", File: filepath.Join(tmpDir, "workflows.go"), Line: 5},
		},
	}

	_, err := r.ExtractValues(tmpDir, discovered)
	if err == nil {
		t.Error("ExtractValues() expected error for missing go.mod")
	}
	if !strings.Contains(err.Error(), "parsing go.mod") {
		t.Errorf("Expected 'parsing go.mod' error, got: %v", err)
	}
}

// Test ExtractDependabot with missing go.mod
func TestRunner_ExtractDependabot_MissingGoMod(t *testing.T) {
	tmpDir := t.TempDir()

	r := NewRunner()

	discovered := &discover.DependabotDiscoveryResult{
		Configs: []discover.DiscoveredDependabot{
			{Name: "Config", File: filepath.Join(tmpDir, "dependabot.go"), Line: 5},
		},
	}

	_, err := r.ExtractDependabot(tmpDir, discovered)
	if err == nil {
		t.Error("ExtractDependabot() expected error for missing go.mod")
	}
	if !strings.Contains(err.Error(), "parsing go.mod") {
		t.Errorf("Expected 'parsing go.mod' error, got: %v", err)
	}
}

// Test ExtractIssueTemplates with missing go.mod
func TestRunner_ExtractIssueTemplates_MissingGoMod(t *testing.T) {
	tmpDir := t.TempDir()

	r := NewRunner()

	discovered := &discover.IssueTemplateDiscoveryResult{
		Templates: []discover.DiscoveredIssueTemplate{
			{Name: "Template", File: filepath.Join(tmpDir, "template.go"), Line: 5},
		},
	}

	_, err := r.ExtractIssueTemplates(tmpDir, discovered)
	if err == nil {
		t.Error("ExtractIssueTemplates() expected error for missing go.mod")
	}
	if !strings.Contains(err.Error(), "parsing go.mod") {
		t.Errorf("Expected 'parsing go.mod' error, got: %v", err)
	}
}

// Test ExtractDiscussionTemplates with missing go.mod
func TestRunner_ExtractDiscussionTemplates_MissingGoMod(t *testing.T) {
	tmpDir := t.TempDir()

	r := NewRunner()

	discovered := &discover.DiscussionTemplateDiscoveryResult{
		Templates: []discover.DiscoveredDiscussionTemplate{
			{Name: "Template", File: filepath.Join(tmpDir, "template.go"), Line: 5},
		},
	}

	_, err := r.ExtractDiscussionTemplates(tmpDir, discovered)
	if err == nil {
		t.Error("ExtractDiscussionTemplates() expected error for missing go.mod")
	}
	if !strings.Contains(err.Error(), "parsing go.mod") {
		t.Errorf("Expected 'parsing go.mod' error, got: %v", err)
	}
}

// Test ExtractPRTemplates with missing go.mod
func TestRunner_ExtractPRTemplates_MissingGoMod(t *testing.T) {
	tmpDir := t.TempDir()

	r := NewRunner()

	discovered := &discover.PRTemplateDiscoveryResult{
		Templates: []discover.DiscoveredPRTemplate{
			{Name: "Template", File: filepath.Join(tmpDir, "template.go"), Line: 5},
		},
	}

	_, err := r.ExtractPRTemplates(tmpDir, discovered)
	if err == nil {
		t.Error("ExtractPRTemplates() expected error for missing go.mod")
	}
	if !strings.Contains(err.Error(), "parsing go.mod") {
		t.Errorf("Expected 'parsing go.mod' error, got: %v", err)
	}
}

// Test ExtractCodeowners with missing go.mod
func TestRunner_ExtractCodeowners_MissingGoMod(t *testing.T) {
	tmpDir := t.TempDir()

	r := NewRunner()

	discovered := &discover.CodeownersDiscoveryResult{
		Configs: []discover.DiscoveredCodeowners{
			{Name: "Config", File: filepath.Join(tmpDir, "codeowners.go"), Line: 5},
		},
	}

	_, err := r.ExtractCodeowners(tmpDir, discovered)
	if err == nil {
		t.Error("ExtractCodeowners() expected error for missing go.mod")
	}
	if !strings.Contains(err.Error(), "parsing go.mod") {
		t.Errorf("Expected 'parsing go.mod' error, got: %v", err)
	}
}

// Test ExtractValues with read-only temp directory to trigger write errors
func TestRunner_ExtractValues_ReadOnlyTempDir(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("Skipping test when running as root")
	}

	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a temp directory that will be used as TempDir, then make it read-only
	readOnlyDir := filepath.Join(tmpDir, "readonly")
	if err := os.MkdirAll(readOnlyDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Make the directory read-only to prevent writing
	if err := os.Chmod(readOnlyDir, 0555); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(readOnlyDir, 0755) // Restore permissions for cleanup

	r := &Runner{
		TempDir: readOnlyDir,
		GoPath:  "go",
	}

	discovered := &discover.DiscoveryResult{
		Workflows: []discover.DiscoveredWorkflow{
			{Name: "CI", File: filepath.Join(tmpDir, "workflows.go"), Line: 5},
		},
	}

	_, err := r.ExtractValues(tmpDir, discovered)
	if err == nil {
		t.Error("ExtractValues() expected error for read-only temp directory")
	}
	// Should fail at temp dir creation
	if !strings.Contains(err.Error(), "creating temp dir") && !strings.Contains(err.Error(), "permission denied") {
		t.Logf("Got error: %v (acceptable if permission-related)", err)
	}
}

// Test ExtractDependabot with read-only temp directory
func TestRunner_ExtractDependabot_ReadOnlyTempDir(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("Skipping test when running as root")
	}

	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	readOnlyDir := filepath.Join(tmpDir, "readonly")
	if err := os.MkdirAll(readOnlyDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(readOnlyDir, 0555); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(readOnlyDir, 0755)

	r := &Runner{
		TempDir: readOnlyDir,
		GoPath:  "go",
	}

	discovered := &discover.DependabotDiscoveryResult{
		Configs: []discover.DiscoveredDependabot{
			{Name: "Config", File: filepath.Join(tmpDir, "dependabot.go"), Line: 5},
		},
	}

	_, err := r.ExtractDependabot(tmpDir, discovered)
	if err == nil {
		t.Error("ExtractDependabot() expected error for read-only temp directory")
	}
}

// Test ExtractIssueTemplates with read-only temp directory
func TestRunner_ExtractIssueTemplates_ReadOnlyTempDir(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("Skipping test when running as root")
	}

	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	readOnlyDir := filepath.Join(tmpDir, "readonly")
	if err := os.MkdirAll(readOnlyDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(readOnlyDir, 0555); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(readOnlyDir, 0755)

	r := &Runner{
		TempDir: readOnlyDir,
		GoPath:  "go",
	}

	discovered := &discover.IssueTemplateDiscoveryResult{
		Templates: []discover.DiscoveredIssueTemplate{
			{Name: "Template", File: filepath.Join(tmpDir, "template.go"), Line: 5},
		},
	}

	_, err := r.ExtractIssueTemplates(tmpDir, discovered)
	if err == nil {
		t.Error("ExtractIssueTemplates() expected error for read-only temp directory")
	}
}

// Test ExtractDiscussionTemplates with read-only temp directory
func TestRunner_ExtractDiscussionTemplates_ReadOnlyTempDir(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("Skipping test when running as root")
	}

	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	readOnlyDir := filepath.Join(tmpDir, "readonly")
	if err := os.MkdirAll(readOnlyDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(readOnlyDir, 0555); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(readOnlyDir, 0755)

	r := &Runner{
		TempDir: readOnlyDir,
		GoPath:  "go",
	}

	discovered := &discover.DiscussionTemplateDiscoveryResult{
		Templates: []discover.DiscoveredDiscussionTemplate{
			{Name: "Template", File: filepath.Join(tmpDir, "template.go"), Line: 5},
		},
	}

	_, err := r.ExtractDiscussionTemplates(tmpDir, discovered)
	if err == nil {
		t.Error("ExtractDiscussionTemplates() expected error for read-only temp directory")
	}
}

// Test ExtractPRTemplates with read-only temp directory
func TestRunner_ExtractPRTemplates_ReadOnlyTempDir(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("Skipping test when running as root")
	}

	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	readOnlyDir := filepath.Join(tmpDir, "readonly")
	if err := os.MkdirAll(readOnlyDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(readOnlyDir, 0555); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(readOnlyDir, 0755)

	r := &Runner{
		TempDir: readOnlyDir,
		GoPath:  "go",
	}

	discovered := &discover.PRTemplateDiscoveryResult{
		Templates: []discover.DiscoveredPRTemplate{
			{Name: "Template", File: filepath.Join(tmpDir, "template.go"), Line: 5},
		},
	}

	_, err := r.ExtractPRTemplates(tmpDir, discovered)
	if err == nil {
		t.Error("ExtractPRTemplates() expected error for read-only temp directory")
	}
}

// Test ExtractCodeowners with read-only temp directory
func TestRunner_ExtractCodeowners_ReadOnlyTempDir(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("Skipping test when running as root")
	}

	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	readOnlyDir := filepath.Join(tmpDir, "readonly")
	if err := os.MkdirAll(readOnlyDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(readOnlyDir, 0555); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(readOnlyDir, 0755)

	r := &Runner{
		TempDir: readOnlyDir,
		GoPath:  "go",
	}

	discovered := &discover.CodeownersDiscoveryResult{
		Configs: []discover.DiscoveredCodeowners{
			{Name: "Config", File: filepath.Join(tmpDir, "codeowners.go"), Line: 5},
		},
	}

	_, err := r.ExtractCodeowners(tmpDir, discovered)
	if err == nil {
		t.Error("ExtractCodeowners() expected error for read-only temp directory")
	}
}
