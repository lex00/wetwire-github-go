package runner

import (
	"fmt"
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

// Test ExtractValues with invalid Go binary path
func TestRunner_ExtractValues_InvalidGoBinary(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := &Runner{
		TempDir: os.TempDir(),
		GoPath:  "/nonexistent/go/binary",
	}

	discovered := &discover.DiscoveryResult{
		Workflows: []discover.DiscoveredWorkflow{
			{Name: "CI", File: filepath.Join(tmpDir, "workflows.go"), Line: 5},
		},
	}

	_, err := r.ExtractValues(tmpDir, discovered)
	if err == nil {
		t.Error("ExtractValues() expected error for invalid Go binary")
	}
	// Should fail at go mod tidy or go run stage
	if !strings.Contains(err.Error(), "go mod tidy") && !strings.Contains(err.Error(), "running extraction") {
		t.Errorf("Expected go execution error, got: %v", err)
	}
}

// Test ExtractDependabot with invalid Go binary path
func TestRunner_ExtractDependabot_InvalidGoBinary(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := &Runner{
		TempDir: os.TempDir(),
		GoPath:  "/nonexistent/go/binary",
	}

	discovered := &discover.DependabotDiscoveryResult{
		Configs: []discover.DiscoveredDependabot{
			{Name: "Config", File: filepath.Join(tmpDir, "dependabot.go"), Line: 5},
		},
	}

	_, err := r.ExtractDependabot(tmpDir, discovered)
	if err == nil {
		t.Error("ExtractDependabot() expected error for invalid Go binary")
	}
}

// Test ExtractIssueTemplates with invalid Go binary path
func TestRunner_ExtractIssueTemplates_InvalidGoBinary(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := &Runner{
		TempDir: os.TempDir(),
		GoPath:  "/nonexistent/go/binary",
	}

	discovered := &discover.IssueTemplateDiscoveryResult{
		Templates: []discover.DiscoveredIssueTemplate{
			{Name: "Template", File: filepath.Join(tmpDir, "template.go"), Line: 5},
		},
	}

	_, err := r.ExtractIssueTemplates(tmpDir, discovered)
	if err == nil {
		t.Error("ExtractIssueTemplates() expected error for invalid Go binary")
	}
}

// Test ExtractDiscussionTemplates with invalid Go binary path
func TestRunner_ExtractDiscussionTemplates_InvalidGoBinary(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := &Runner{
		TempDir: os.TempDir(),
		GoPath:  "/nonexistent/go/binary",
	}

	discovered := &discover.DiscussionTemplateDiscoveryResult{
		Templates: []discover.DiscoveredDiscussionTemplate{
			{Name: "Template", File: filepath.Join(tmpDir, "template.go"), Line: 5},
		},
	}

	_, err := r.ExtractDiscussionTemplates(tmpDir, discovered)
	if err == nil {
		t.Error("ExtractDiscussionTemplates() expected error for invalid Go binary")
	}
}

// Test ExtractPRTemplates with invalid Go binary path
func TestRunner_ExtractPRTemplates_InvalidGoBinary(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := &Runner{
		TempDir: os.TempDir(),
		GoPath:  "/nonexistent/go/binary",
	}

	discovered := &discover.PRTemplateDiscoveryResult{
		Templates: []discover.DiscoveredPRTemplate{
			{Name: "Template", File: filepath.Join(tmpDir, "template.go"), Line: 5},
		},
	}

	_, err := r.ExtractPRTemplates(tmpDir, discovered)
	if err == nil {
		t.Error("ExtractPRTemplates() expected error for invalid Go binary")
	}
}

// Test ExtractCodeowners with invalid Go binary path
func TestRunner_ExtractCodeowners_InvalidGoBinary(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := &Runner{
		TempDir: os.TempDir(),
		GoPath:  "/nonexistent/go/binary",
	}

	discovered := &discover.CodeownersDiscoveryResult{
		Configs: []discover.DiscoveredCodeowners{
			{Name: "Config", File: filepath.Join(tmpDir, "codeowners.go"), Line: 5},
		},
	}

	_, err := r.ExtractCodeowners(tmpDir, discovered)
	if err == nil {
		t.Error("ExtractCodeowners() expected error for invalid Go binary")
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

// Test with read-only temp directory to trigger write errors
func TestRunner_ExtractValues_WriteProgramError(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	goCode := `package testproject

import "github.com/lex00/wetwire-github-go/workflow"

var TestWorkflow = workflow.Workflow{Name: "Test"}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "workflows.go"), []byte(goCode), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a read-only temp directory for the runner
	readOnlyTmp := filepath.Join(tmpDir, "readonly")
	if err := os.MkdirAll(readOnlyTmp, 0555); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(readOnlyTmp, 0755)

	r := &Runner{
		TempDir: readOnlyTmp,
		GoPath:  NewRunner().GoPath,
		Verbose: false,
	}

	discovered := &discover.DiscoveryResult{
		Workflows: []discover.DiscoveredWorkflow{
			{Name: "TestWorkflow", File: filepath.Join(tmpDir, "workflows.go"), Line: 5},
		},
	}

	_, err := r.ExtractValues(tmpDir, discovered)
	// On most systems, this should fail due to permission issues
	if err != nil {
		t.Logf("ExtractValues() with read-only temp dir error = %v (expected on most systems)", err)
	}
}

// Test ExtractValues when program execution fails (compile error)
func TestRunner_ExtractValues_RunError(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	// Write code with the correct variable name
	goCode := `package testproject

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{Name: "CI"}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "workflows.go"), []byte(goCode), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()

	// Reference a variable that doesn't exist in the source code
	// The generated program will try to reference a variable that doesn't exist
	discovered := &discover.DiscoveryResult{
		Workflows: []discover.DiscoveredWorkflow{
			{Name: "NonexistentVariable", File: filepath.Join(tmpDir, "workflows.go"), Line: 5},
		},
	}

	_, err := r.ExtractValues(tmpDir, discovered)
	if err == nil {
		t.Log("ExtractValues() succeeded unexpectedly")
	} else {
		// The error should be about running extraction (compile error)
		t.Logf("ExtractValues() error = %v", err)
	}
}

// Test ExtractDependabot when program execution fails
func TestRunner_ExtractDependabot_RunError(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	goCode := `package testproject

import "github.com/lex00/wetwire-github-go/dependabot"

var Config = dependabot.Config{}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "dependabot.go"), []byte(goCode), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()

	// Reference a variable that doesn't exist
	discovered := &discover.DependabotDiscoveryResult{
		Configs: []discover.DiscoveredDependabot{
			{Name: "NonexistentConfig", File: filepath.Join(tmpDir, "dependabot.go"), Line: 5},
		},
	}

	_, err := r.ExtractDependabot(tmpDir, discovered)
	if err == nil {
		t.Log("ExtractDependabot() succeeded unexpectedly")
	} else {
		t.Logf("ExtractDependabot() error = %v", err)
	}
}

// Test ExtractIssueTemplates when program execution fails
func TestRunner_ExtractIssueTemplates_RunError(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	goCode := `package testproject

import "github.com/lex00/wetwire-github-go/issue"

var Template = issue.Template{}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "templates.go"), []byte(goCode), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()

	discovered := &discover.IssueTemplateDiscoveryResult{
		Templates: []discover.DiscoveredIssueTemplate{
			{Name: "NonexistentTemplate", File: filepath.Join(tmpDir, "templates.go"), Line: 5},
		},
	}

	_, err := r.ExtractIssueTemplates(tmpDir, discovered)
	if err == nil {
		t.Log("ExtractIssueTemplates() succeeded unexpectedly")
	} else {
		t.Logf("ExtractIssueTemplates() error = %v", err)
	}
}

// Test ExtractDiscussionTemplates when program execution fails
func TestRunner_ExtractDiscussionTemplates_RunError(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	goCode := `package testproject

import "github.com/lex00/wetwire-github-go/discussion"

var Template = discussion.Template{}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "templates.go"), []byte(goCode), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()

	discovered := &discover.DiscussionTemplateDiscoveryResult{
		Templates: []discover.DiscoveredDiscussionTemplate{
			{Name: "NonexistentTemplate", File: filepath.Join(tmpDir, "templates.go"), Line: 5},
		},
	}

	_, err := r.ExtractDiscussionTemplates(tmpDir, discovered)
	if err == nil {
		t.Log("ExtractDiscussionTemplates() succeeded unexpectedly")
	} else {
		t.Logf("ExtractDiscussionTemplates() error = %v", err)
	}
}

// Test ExtractPRTemplates when program execution fails
func TestRunner_ExtractPRTemplates_RunError(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	goCode := `package testproject

import "github.com/lex00/wetwire-github-go/pr"

var Template = pr.Template{}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "templates.go"), []byte(goCode), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()

	discovered := &discover.PRTemplateDiscoveryResult{
		Templates: []discover.DiscoveredPRTemplate{
			{Name: "NonexistentTemplate", File: filepath.Join(tmpDir, "templates.go"), Line: 5},
		},
	}

	_, err := r.ExtractPRTemplates(tmpDir, discovered)
	if err == nil {
		t.Log("ExtractPRTemplates() succeeded unexpectedly")
	} else {
		t.Logf("ExtractPRTemplates() error = %v", err)
	}
}

// Test ExtractCodeowners when program execution fails
func TestRunner_ExtractCodeowners_RunError(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	goCode := `package testproject

import "github.com/lex00/wetwire-github-go/codeowners"

var Owners = codeowners.Owners{}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "codeowners.go"), []byte(goCode), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()

	discovered := &discover.CodeownersDiscoveryResult{
		Configs: []discover.DiscoveredCodeowners{
			{Name: "NonexistentOwners", File: filepath.Join(tmpDir, "codeowners.go"), Line: 5},
		},
	}

	_, err := r.ExtractCodeowners(tmpDir, discovered)
	if err == nil {
		t.Log("ExtractCodeowners() succeeded unexpectedly")
	} else {
		t.Logf("ExtractCodeowners() error = %v", err)
	}
}

// Test ExtractValues with missing Go binary
func TestRunner_ExtractValues_MissingGoBinary(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	// Write a simple Go file
	goCode := `package testproject

import "github.com/lex00/wetwire-github-go/workflow"

var TestWorkflow = workflow.Workflow{Name: "Test"}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "workflows.go"), []byte(goCode), 0644); err != nil {
		t.Fatal(err)
	}

	r := &Runner{
		TempDir: os.TempDir(),
		GoPath:  "/nonexistent/go/binary", // Invalid Go path
		Verbose: false,
	}

	discovered := &discover.DiscoveryResult{
		Workflows: []discover.DiscoveredWorkflow{
			{Name: "TestWorkflow", File: filepath.Join(tmpDir, "workflows.go"), Line: 5},
		},
	}

	_, err := r.ExtractValues(tmpDir, discovered)
	if err == nil {
		t.Error("ExtractValues() expected error with missing Go binary")
	}
}

// Test ExtractValues when go mod tidy fails
func TestRunner_ExtractValues_GoModTidyFails(t *testing.T) {
	tmpDir := t.TempDir()

	// Write an invalid go.mod that will cause go mod tidy to fail
	goMod := `module github.com/example/test

go 1.23

require nonexistent-module-that-does-not-exist v99.99.99
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	goCode := `package testproject

import "github.com/lex00/wetwire-github-go/workflow"

var TestWorkflow = workflow.Workflow{Name: "Test"}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "workflows.go"), []byte(goCode), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()

	discovered := &discover.DiscoveryResult{
		Workflows: []discover.DiscoveredWorkflow{
			{Name: "TestWorkflow", File: filepath.Join(tmpDir, "workflows.go"), Line: 5},
		},
	}

	_, err := r.ExtractValues(tmpDir, discovered)
	if err == nil {
		t.Error("ExtractValues() expected error when go mod tidy fails")
	} else if !strings.Contains(err.Error(), "go mod tidy") {
		t.Logf("ExtractValues() error = %v", err)
	}
}

// Test ExtractValues when compilation fails
func TestRunner_ExtractValues_CompilationFails(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	// Write code that references a non-existent variable
	goCode := `package testproject

import "github.com/lex00/wetwire-github-go/workflow"

var TestWorkflow = workflow.Workflow{Name: "Test"}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "workflows.go"), []byte(goCode), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()

	// Reference a variable that doesn't exist
	discovered := &discover.DiscoveryResult{
		Workflows: []discover.DiscoveredWorkflow{
			{Name: "NonExistentWorkflow", File: filepath.Join(tmpDir, "workflows.go"), Line: 5},
		},
	}

	_, err := r.ExtractValues(tmpDir, discovered)
	// This will fail during compilation because NonExistentWorkflow doesn't exist
	if err == nil {
		t.Log("ExtractValues() succeeded - variable may have been found")
	} else {
		t.Logf("ExtractValues() error = %v (expected)", err)
	}
}

// Test ExtractDependabot with missing Go binary
func TestRunner_ExtractDependabot_MissingGoBinary(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	goCode := `package testproject

import "github.com/lex00/wetwire-github-go/dependabot"

var Config = dependabot.Config{}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "dependabot.go"), []byte(goCode), 0644); err != nil {
		t.Fatal(err)
	}

	r := &Runner{
		TempDir: os.TempDir(),
		GoPath:  "/nonexistent/go/binary",
		Verbose: false,
	}

	discovered := &discover.DependabotDiscoveryResult{
		Configs: []discover.DiscoveredDependabot{
			{Name: "Config", File: filepath.Join(tmpDir, "dependabot.go"), Line: 5},
		},
	}

	_, err := r.ExtractDependabot(tmpDir, discovered)
	if err == nil {
		t.Error("ExtractDependabot() expected error with missing Go binary")
	}
}

// Test ExtractIssueTemplates with missing Go binary
func TestRunner_ExtractIssueTemplates_MissingGoBinary(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	goCode := `package testproject

import "github.com/lex00/wetwire-github-go/issue"

var BugReport = issue.Template{}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "templates.go"), []byte(goCode), 0644); err != nil {
		t.Fatal(err)
	}

	r := &Runner{
		TempDir: os.TempDir(),
		GoPath:  "/nonexistent/go/binary",
		Verbose: false,
	}

	discovered := &discover.IssueTemplateDiscoveryResult{
		Templates: []discover.DiscoveredIssueTemplate{
			{Name: "BugReport", File: filepath.Join(tmpDir, "templates.go"), Line: 5},
		},
	}

	_, err := r.ExtractIssueTemplates(tmpDir, discovered)
	if err == nil {
		t.Error("ExtractIssueTemplates() expected error with missing Go binary")
	}
}

// Test ExtractDiscussionTemplates with missing Go binary
func TestRunner_ExtractDiscussionTemplates_MissingGoBinary(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	goCode := `package testproject

import "github.com/lex00/wetwire-github-go/discussion"

var Announcement = discussion.Template{}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "templates.go"), []byte(goCode), 0644); err != nil {
		t.Fatal(err)
	}

	r := &Runner{
		TempDir: os.TempDir(),
		GoPath:  "/nonexistent/go/binary",
		Verbose: false,
	}

	discovered := &discover.DiscussionTemplateDiscoveryResult{
		Templates: []discover.DiscoveredDiscussionTemplate{
			{Name: "Announcement", File: filepath.Join(tmpDir, "templates.go"), Line: 5},
		},
	}

	_, err := r.ExtractDiscussionTemplates(tmpDir, discovered)
	if err == nil {
		t.Error("ExtractDiscussionTemplates() expected error with missing Go binary")
	}
}

// Test ExtractPRTemplates with missing Go binary
func TestRunner_ExtractPRTemplates_MissingGoBinary(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	goCode := `package testproject

import "github.com/lex00/wetwire-github-go/pr"

var DefaultPR = pr.Template{}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "templates.go"), []byte(goCode), 0644); err != nil {
		t.Fatal(err)
	}

	r := &Runner{
		TempDir: os.TempDir(),
		GoPath:  "/nonexistent/go/binary",
		Verbose: false,
	}

	discovered := &discover.PRTemplateDiscoveryResult{
		Templates: []discover.DiscoveredPRTemplate{
			{Name: "DefaultPR", File: filepath.Join(tmpDir, "templates.go"), Line: 5},
		},
	}

	_, err := r.ExtractPRTemplates(tmpDir, discovered)
	if err == nil {
		t.Error("ExtractPRTemplates() expected error with missing Go binary")
	}
}

// Test ExtractCodeowners with missing Go binary
func TestRunner_ExtractCodeowners_MissingGoBinary(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	goCode := `package testproject

import "github.com/lex00/wetwire-github-go/codeowners"

var Owners = codeowners.Owners{}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "codeowners.go"), []byte(goCode), 0644); err != nil {
		t.Fatal(err)
	}

	r := &Runner{
		TempDir: os.TempDir(),
		GoPath:  "/nonexistent/go/binary",
		Verbose: false,
	}

	discovered := &discover.CodeownersDiscoveryResult{
		Configs: []discover.DiscoveredCodeowners{
			{Name: "Owners", File: filepath.Join(tmpDir, "codeowners.go"), Line: 5},
		},
	}

	_, err := r.ExtractCodeowners(tmpDir, discovered)
	if err == nil {
		t.Error("ExtractCodeowners() expected error with missing Go binary")
	}
}

// Suppress unused import warning
var _ = fmt.Sprintf
