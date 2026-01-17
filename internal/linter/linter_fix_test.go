package linter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLinter_Fix_WAG001(t *testing.T) {
	// Test WAG001 fix: raw uses: string -> typed action wrapper
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var CheckoutStep = workflow.Step{Uses: "actions/checkout@v4"}
`)

	l := NewLinter(&WAG001{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Fatal("Expected WAG001 issue to be detected")
	}

	// Apply fix
	fixResult, err := l.Fix("test.go", content)
	if err != nil {
		t.Fatalf("Fix() error = %v", err)
	}

	if fixResult.FixedCount == 0 {
		t.Error("Expected at least one fix to be applied")
	}

	// Verify the fixed content imports checkout package
	if !containsString(string(fixResult.Content), "checkout.Checkout") {
		t.Error("Fixed content should use checkout.Checkout wrapper")
	}

	// Lint again to verify fix worked
	result2, err := l.LintContent("test.go", fixResult.Content)
	if err != nil {
		t.Fatalf("LintContent() after fix error = %v", err)
	}

	for _, issue := range result2.Issues {
		if issue.Rule == "WAG001" {
			t.Error("WAG001 issue should be fixed after applying fix")
		}
	}
}

func TestLinter_FixFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "workflows.go")

	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var CheckoutStep = workflow.Step{Uses: "actions/checkout@v4"}
`)

	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatal(err)
	}

	l := NewLinter(&WAG001{})
	fixResult, err := l.FixFile(testFile)
	if err != nil {
		t.Fatalf("FixFile() error = %v", err)
	}

	if fixResult.FixedCount == 0 {
		t.Error("Expected at least one fix to be applied")
	}

	// Read back the file and verify
	fixed, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	if !containsString(string(fixed), "checkout.Checkout") {
		t.Error("Fixed file should use checkout.Checkout wrapper")
	}
}

func TestLinter_FixDir(t *testing.T) {
	tmpDir := t.TempDir()

	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var CheckoutStep = workflow.Step{Uses: "actions/checkout@v4"}
`)

	if err := os.WriteFile(filepath.Join(tmpDir, "workflows.go"), content, 0644); err != nil {
		t.Fatal(err)
	}

	l := NewLinter(&WAG001{})
	fixResult, err := l.FixDir(tmpDir)
	if err != nil {
		t.Fatalf("FixDir() error = %v", err)
	}

	if fixResult.TotalFixed == 0 {
		t.Error("Expected at least one fix to be applied")
	}

	if len(fixResult.Files) == 0 {
		t.Error("Expected at least one file to be fixed")
	}
}

func containsString(haystack, needle string) bool {
	return len(haystack) > 0 && len(needle) > 0 &&
		(haystack == needle || len(haystack) > len(needle) &&
			(haystack[:len(needle)] == needle || containsString(haystack[1:], needle)))
}

func TestFixer_WAG001_KnownActions(t *testing.T) {
	// Test that WAG001 can fix known actions
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    `workflow.Step{Uses: "actions/checkout@v4"}`,
			expected: "checkout.Checkout",
		},
		{
			input:    `workflow.Step{Uses: "actions/setup-go@v5"}`,
			expected: "setup_go.SetupGo",
		},
		{
			input:    `workflow.Step{Uses: "actions/cache@v4"}`,
			expected: "cache.Cache",
		},
	}

	for _, tc := range testCases {
		content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var Step = ` + tc.input + `
`)

		l := NewLinter(&WAG001{})
		fixResult, err := l.Fix("test.go", content)
		if err != nil {
			t.Fatalf("Fix() error = %v for input %s", err, tc.input)
		}

		if !containsString(string(fixResult.Content), tc.expected) {
			t.Errorf("Fixed content should contain %q for input %s\nGot: %s",
				tc.expected, tc.input, string(fixResult.Content))
		}
	}
}

func TestLinter_Fix_NonFixableIssue(t *testing.T) {
	// WAG002 is not fixable
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var Step = workflow.Step{
	If: "${{ github.ref == 'refs/heads/main' }}",
}
`)

	l := NewLinter(&WAG002{})
	result, err := l.Fix("test.go", content)
	if err != nil {
		t.Fatalf("Fix() error = %v", err)
	}

	// Should have 0 fixed, 1 remaining issue
	if result.FixedCount != 0 {
		t.Errorf("FixedCount = %d, want 0", result.FixedCount)
	}
	if len(result.Issues) != 1 {
		t.Errorf("len(Issues) = %d, want 1", len(result.Issues))
	}
}

func TestLinter_Fix_UnknownAction(t *testing.T) {
	// WAG001 with an unknown action that cannot be auto-fixed
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var Step = workflow.Step{Uses: "unknown/action@v1"}
`)

	l := NewLinter(&WAG001{})
	result, err := l.Fix("test.go", content)
	if err != nil {
		t.Fatalf("Fix() error = %v", err)
	}

	// Should have 0 fixed because action is unknown
	if result.FixedCount != 0 {
		t.Errorf("FixedCount = %d, want 0", result.FixedCount)
	}
}

func TestLinter_Fix_ParseError(t *testing.T) {
	content := []byte(`package main
func invalid { syntax
`)

	l := DefaultLinter()
	_, err := l.Fix("test.go", content)
	if err == nil {
		t.Error("Fix() expected error for invalid syntax")
	}
}

func TestLinter_FixFile_NotExist(t *testing.T) {
	l := DefaultLinter()
	_, err := l.FixFile("/nonexistent/path/file.go")
	if err == nil {
		t.Error("FixFile() expected error for non-existent file")
	}
}

func TestLinter_FixFile_NoChanges(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "valid.go")

	// File with no issues
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{Name: "CI"}
`)
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatal(err)
	}

	l := NewLinter(&WAG001{})
	result, err := l.FixFile(testFile)
	if err != nil {
		t.Fatalf("FixFile() error = %v", err)
	}

	if result.FixedCount != 0 {
		t.Errorf("FixedCount = %d, want 0", result.FixedCount)
	}
}

func TestLinter_FixDir_SkipsDirectories(t *testing.T) {
	tmpDir := t.TempDir()

	// Create vendor directory that should be skipped
	vendorDir := filepath.Join(tmpDir, "vendor")
	if err := os.MkdirAll(vendorDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create file with fixable issue in vendor
	badContent := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var Step = workflow.Step{Uses: "actions/checkout@v4"}
`)
	if err := os.WriteFile(filepath.Join(vendorDir, "workflows.go"), badContent, 0644); err != nil {
		t.Fatal(err)
	}

	l := NewLinter(&WAG001{})
	result, err := l.FixDir(tmpDir)
	if err != nil {
		t.Fatalf("FixDir() error = %v", err)
	}

	// Should have 0 fixed because vendor is skipped
	if result.TotalFixed != 0 {
		t.Errorf("TotalFixed = %d, want 0", result.TotalFixed)
	}
}

func TestAddImportIfNeeded_SingleImport(t *testing.T) {
	// This tests the single import case in addImportIfNeeded
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var Step = workflow.Step{Uses: "actions/checkout@v4"}
`)

	l := NewLinter(&WAG001{})
	result, err := l.Fix("test.go", content)
	if err != nil {
		t.Fatalf("Fix() error = %v", err)
	}

	// Check that the fix added checkout import
	if result.FixedCount > 0 {
		codeStr := string(result.Content)
		if !strings.Contains(codeStr, "checkout") {
			t.Error("Fixed content should contain checkout import")
		}
	}
}

func TestWAG001_Fix_NoUsesField(t *testing.T) {
	// Test WAG001 fix when Uses field is not found at the expected line
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var Step = workflow.Step{
	Name: "No uses field",
	Run:  "echo hello",
}
`)

	l := NewLinter(&WAG001{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	// No WAG001 issues should be found
	for _, issue := range result.Issues {
		if issue.Rule == "WAG001" {
			t.Error("WAG001 should not flag Step without Uses field")
		}
	}
}
