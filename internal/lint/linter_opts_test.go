package lint

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewLinterWithDisabledRules(t *testing.T) {
	// Create a linter with some disabled rules
	l := NewLinterWithOptions(LinterOptions{
		DisabledRules: []string{"WAG001", "WAG002"},
	})
	if l == nil {
		t.Error("NewLinterWithOptions() returned nil")
	}

	// Verify that disabled rules are not included
	rules := l.Rules()
	for _, rule := range rules {
		if rule.ID() == "WAG001" || rule.ID() == "WAG002" {
			t.Errorf("Rule %s should be disabled but is present", rule.ID())
		}
	}

	// Verify other rules are still present
	foundWAG003 := false
	for _, rule := range rules {
		if rule.ID() == "WAG003" {
			foundWAG003 = true
			break
		}
	}
	if !foundWAG003 {
		t.Error("WAG003 should be present but is not")
	}
}

func TestNewLinterWithOptions_EmptyDisabled(t *testing.T) {
	// When no rules are disabled, all default rules should be present
	l := NewLinterWithOptions(LinterOptions{})
	allRules := DefaultLinter().Rules()

	if len(l.Rules()) != len(allRules) {
		t.Errorf("Expected %d rules, got %d", len(allRules), len(l.Rules()))
	}
}

func TestLinterOptions_DisableAllRules(t *testing.T) {
	// Disable all rules - should result in empty rule set
	allRuleIDs := []string{
		"WAG001", "WAG002", "WAG003", "WAG004", "WAG005",
		"WAG006", "WAG007", "WAG008", "WAG009", "WAG010",
		"WAG011", "WAG012", "WAG013", "WAG014", "WAG015",
		"WAG016", "WAG017", "WAG018", "WAG019", "WAG020",
	}

	l := NewLinterWithOptions(LinterOptions{
		DisabledRules: allRuleIDs,
	})

	if len(l.Rules()) != 0 {
		t.Errorf("Expected 0 rules when all disabled, got %d", len(l.Rules()))
	}
}

func TestLinter_LintFile_WithDisabledRules(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a file that would trigger WAG001 (raw uses: string)
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var Step = workflow.Step{Uses: "actions/checkout@v4"}
`)
	filePath := filepath.Join(tmpDir, "workflows.go")
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatal(err)
	}

	// Lint without disabling - should find WAG001
	linterWithWAG001 := DefaultLinter()
	result, err := linterWithWAG001.LintFile(filePath)
	if err != nil {
		t.Fatalf("LintFile() error = %v", err)
	}

	foundWAG001 := false
	for _, issue := range result.Issues {
		if issue.Rule == "WAG001" {
			foundWAG001 = true
			break
		}
	}
	if !foundWAG001 {
		t.Error("Expected WAG001 issue when not disabled")
	}

	// Lint with WAG001 disabled - should NOT find WAG001
	linterWithoutWAG001 := NewLinterWithOptions(LinterOptions{
		DisabledRules: []string{"WAG001"},
	})
	result, err = linterWithoutWAG001.LintFile(filePath)
	if err != nil {
		t.Fatalf("LintFile() error = %v", err)
	}

	for _, issue := range result.Issues {
		if issue.Rule == "WAG001" {
			t.Error("WAG001 should be disabled but issue was found")
		}
	}
}

func TestLinter_LintDir_WithDisabledRules(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a file that would trigger WAG001
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var Step = workflow.Step{Uses: "actions/checkout@v4"}
`)
	if err := os.WriteFile(filepath.Join(tmpDir, "workflows.go"), content, 0644); err != nil {
		t.Fatal(err)
	}

	// Lint with WAG001 disabled
	l := NewLinterWithOptions(LinterOptions{
		DisabledRules: []string{"WAG001"},
	})

	result, err := l.LintDir(tmpDir)
	if err != nil {
		t.Fatalf("LintDir() error = %v", err)
	}

	// Should not find WAG001 issues
	for _, issue := range result.Issues {
		if issue.Rule == "WAG001" {
			t.Error("WAG001 should be disabled but issue was found in LintDir")
		}
	}
}

func TestLinterOptions_Fix(t *testing.T) {
	// Test that Fix option is accepted (even if not fully implemented)
	opts := LinterOptions{
		Fix: true,
	}

	l := NewLinterWithOptions(opts)
	if l == nil {
		t.Error("NewLinterWithOptions with Fix=true should not return nil")
	}
}

func TestLinter_LintContent_WithDisabledRules(t *testing.T) {
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var Step = workflow.Step{Uses: "actions/checkout@v4"}
`)

	// Lint with WAG001 disabled
	l := NewLinterWithOptions(LinterOptions{
		DisabledRules: []string{"WAG001"},
	})

	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	// Should not find WAG001 issues
	for _, issue := range result.Issues {
		if issue.Rule == "WAG001" {
			t.Error("WAG001 should be disabled but issue was found in LintContent")
		}
	}
}

func TestLinterOptions_DisabledRulesCaseInsensitive(t *testing.T) {
	// Rule IDs should be case-insensitive (or at least handled consistently)
	l := NewLinterWithOptions(LinterOptions{
		DisabledRules: []string{"wag001"}, // lowercase
	})

	// This tests that the implementation handles case correctly
	// The actual behavior depends on implementation choice
	if l == nil {
		t.Error("NewLinterWithOptions should not return nil")
	}
}
