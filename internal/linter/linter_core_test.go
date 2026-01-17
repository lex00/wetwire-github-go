package linter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewLinter(t *testing.T) {
	l := NewLinter(&WAG001{})
	if l == nil {
		t.Error("NewLinter() returned nil")
	}
	if len(l.Rules()) != 1 {
		t.Errorf("len(Rules()) = %d, want 1", len(l.Rules()))
	}
}

func TestDefaultLinter(t *testing.T) {
	l := DefaultLinter()
	if l == nil {
		t.Error("DefaultLinter() returned nil")
	}
	if len(l.Rules()) != 20 {
		t.Errorf("len(Rules()) = %d, want 20", len(l.Rules()))
	}
}

func TestLinter_LintContent_Valid(t *testing.T) {
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
}
`)

	l := DefaultLinter()
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if !result.Success {
		for _, issue := range result.Issues {
			t.Logf("Issue: %s:%d: [%s] %s", issue.File, issue.Line, issue.Rule, issue.Message)
		}
	}
}

func TestLinter_LintDir(t *testing.T) {
	tmpDir := t.TempDir()

	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
}
`)

	if err := os.WriteFile(filepath.Join(tmpDir, "workflows.go"), content, 0644); err != nil {
		t.Fatal(err)
	}

	l := DefaultLinter()
	result, err := l.LintDir(tmpDir)
	if err != nil {
		t.Fatalf("LintDir() error = %v", err)
	}

	if result == nil {
		t.Error("LintDir() returned nil result")
	}
}

func TestRuleIDs(t *testing.T) {
	rules := []Rule{
		&WAG001{},
		&WAG002{},
		&WAG003{},
		&WAG004{},
		&WAG005{},
		&WAG006{},
		&WAG007{},
		&WAG008{},
		&WAG009{},
		&WAG010{},
		&WAG011{},
		&WAG012{},
		&WAG013{},
		&WAG014{},
		&WAG015{},
		&WAG016{},
		&WAG017{},
	}

	expectedIDs := []string{
		"WAG001", "WAG002", "WAG003", "WAG004",
		"WAG005", "WAG006", "WAG007", "WAG008",
		"WAG009", "WAG010", "WAG011", "WAG012",
		"WAG013", "WAG014", "WAG015", "WAG016",
		"WAG017",
	}

	for i, rule := range rules {
		if rule.ID() != expectedIDs[i] {
			t.Errorf("Rule %d ID() = %q, want %q", i, rule.ID(), expectedIDs[i])
		}
		if rule.Description() == "" {
			t.Errorf("Rule %s has empty description", rule.ID())
		}
	}
}

func TestLinter_AddRule(t *testing.T) {
	l := NewLinter()
	if len(l.Rules()) != 0 {
		t.Error("New linter should have no rules")
	}

	l.AddRule(&WAG001{})
	if len(l.Rules()) != 1 {
		t.Errorf("After AddRule, len(Rules()) = %d, want 1", len(l.Rules()))
	}
}

func TestLinter_LintContent_ParseError(t *testing.T) {
	content := []byte(`package main
func invalid syntax {
`)

	l := DefaultLinter()
	_, err := l.LintContent("test.go", content)
	if err == nil {
		t.Error("LintContent() expected error for invalid syntax")
	}
}

func TestFixer_Interface(t *testing.T) {
	// Test that Fixer interface works
	var rule Rule = &WAG001{}
	_, isFixer := rule.(Fixer)
	// WAG001 is fixable - it should implement Fixer
	if !isFixer {
		t.Error("WAG001 should implement Fixer interface")
	}
}

func TestLinter_LintDir_SkipsDirectories(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a hidden directory that should be skipped
	hiddenDir := filepath.Join(tmpDir, ".hidden")
	if err := os.MkdirAll(hiddenDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create vendor directory that should be skipped
	vendorDir := filepath.Join(tmpDir, "vendor")
	if err := os.MkdirAll(vendorDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create testdata directory that should be skipped
	testdataDir := filepath.Join(tmpDir, "testdata")
	if err := os.MkdirAll(testdataDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create files in skipped directories
	badContent := []byte(`package main
var token = "ghp_secrettoken"
`)
	if err := os.WriteFile(filepath.Join(hiddenDir, "bad.go"), badContent, 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(vendorDir, "bad.go"), badContent, 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(testdataDir, "bad.go"), badContent, 0644); err != nil {
		t.Fatal(err)
	}

	// Create a test file that should also be skipped
	if err := os.WriteFile(filepath.Join(tmpDir, "workflows_test.go"), badContent, 0644); err != nil {
		t.Fatal(err)
	}

	// Create a valid file in main directory
	validContent := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{Name: "CI"}
`)
	if err := os.WriteFile(filepath.Join(tmpDir, "workflows.go"), validContent, 0644); err != nil {
		t.Fatal(err)
	}

	l := NewLinter(&WAG003{})
	result, err := l.LintDir(tmpDir)
	if err != nil {
		t.Fatalf("LintDir() error = %v", err)
	}

	// Should not have found any secrets (all bad files should be skipped)
	for _, issue := range result.Issues {
		if issue.Rule == "WAG003" {
			t.Errorf("WAG003 found secret in file that should have been skipped: %s", issue.File)
		}
	}
}

func TestLinter_LintDir_ParseError(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a file with invalid Go syntax
	invalidContent := []byte(`package main
func invalid { syntax
`)
	if err := os.WriteFile(filepath.Join(tmpDir, "invalid.go"), invalidContent, 0644); err != nil {
		t.Fatal(err)
	}

	l := DefaultLinter()
	result, err := l.LintDir(tmpDir)
	if err != nil {
		t.Fatalf("LintDir() error = %v", err)
	}

	// Should have captured a parse error
	found := false
	for _, issue := range result.Issues {
		if issue.Rule == "parse-error" {
			found = true
		}
	}
	if !found {
		t.Error("Expected parse-error issue not found")
	}
}

func TestLinter_LintFile_NotExist(t *testing.T) {
	l := DefaultLinter()
	_, err := l.LintFile("/nonexistent/path/file.go")
	if err == nil {
		t.Error("LintFile() expected error for non-existent file")
	}
}

func TestGetTypeName_SelectorExpr(t *testing.T) {
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var Step = workflow.Step{
	Name: "Test",
	Run:  "echo hello",
}
`)

	l := NewLinter(&WAG001{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	// This tests the SelectorExpr path in getTypeName
	_ = result
}
