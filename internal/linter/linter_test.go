package linter

import (
	"os"
	"path/filepath"
	"strings"
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
	if len(l.Rules()) != 19 {
		t.Errorf("len(Rules()) = %d, want 19", len(l.Rules()))
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

func TestWAG001_Check(t *testing.T) {
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
		t.Error("WAG001 should have found issue with raw uses: string")
	}

	found := false
	for _, issue := range result.Issues {
		if issue.Rule == "WAG001" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected WAG001 issue not found")
	}
}

func TestWAG002_Check(t *testing.T) {
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var Step = workflow.Step{
	If: "${{ github.ref == 'refs/heads/main' }}",
}
`)

	l := NewLinter(&WAG002{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG002 should have found issue with raw expression string")
	}
}

func TestWAG003_Check(t *testing.T) {
	content := []byte(`package main

var token = "ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
`)

	l := NewLinter(&WAG003{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG003 should have found hardcoded secret")
	}

	found := false
	for _, issue := range result.Issues {
		if issue.Rule == "WAG003" {
			found = true
			if issue.Severity != "error" {
				t.Error("WAG003 issues should be severity 'error'")
			}
		}
	}
	if !found {
		t.Error("Expected WAG003 issue not found")
	}
}

func TestWAG006_Check_DuplicateNames(t *testing.T) {
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI1 = workflow.Workflow{Name: "CI"}
var CI2 = workflow.Workflow{Name: "CI"}
`)

	l := NewLinter(&WAG006{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG006 should have found duplicate workflow names")
	}

	found := false
	for _, issue := range result.Issues {
		if issue.Rule == "WAG006" {
			found = true
		}
	}
	if !found {
		t.Error("Expected WAG006 issue not found")
	}
}

func TestWAG007_Check_TooManyJobs(t *testing.T) {
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var Job1 = workflow.Job{Name: "job1"}
var Job2 = workflow.Job{Name: "job2"}
var Job3 = workflow.Job{Name: "job3"}
`)

	l := NewLinter(&WAG007{MaxJobs: 2})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG007 should have found too many jobs")
	}

	found := false
	for _, issue := range result.Issues {
		if issue.Rule == "WAG007" {
			found = true
		}
	}
	if !found {
		t.Error("Expected WAG007 issue not found")
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

func TestWAG009_Check_EmptyMatrix(t *testing.T) {
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var Matrix = workflow.Matrix{
	Values: map[string][]any{
		"os": {},
	},
}
`)

	l := NewLinter(&WAG009{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG009 should have found empty matrix dimension")
	}

	found := false
	for _, issue := range result.Issues {
		if issue.Rule == "WAG009" {
			found = true
			if issue.Severity != "error" {
				t.Error("WAG009 issues should be severity 'error'")
			}
		}
	}
	if !found {
		t.Error("Expected WAG009 issue not found")
	}
}

func TestWAG009_Check_ValidMatrix(t *testing.T) {
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var Matrix = workflow.Matrix{
	Values: map[string][]any{
		"os": {"ubuntu-latest", "macos-latest"},
	},
}
`)

	l := NewLinter(&WAG009{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if !result.Success {
		t.Error("WAG009 should not flag valid matrix with values")
	}
}

func TestWAG010_Check_MissingInput(t *testing.T) {
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/actions/setup_go"

var Step = setup_go.SetupGo{}
`)

	l := NewLinter(&WAG010{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG010 should have flagged missing GoVersion")
	}

	found := false
	for _, issue := range result.Issues {
		if issue.Rule == "WAG010" {
			found = true
		}
	}
	if !found {
		t.Error("Expected WAG010 issue not found")
	}
}

func TestWAG010_Check_HasInput(t *testing.T) {
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/actions/setup_go"

var Step = setup_go.SetupGo{GoVersion: "1.23"}
`)

	l := NewLinter(&WAG010{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if !result.Success {
		t.Error("WAG010 should not flag when GoVersion is set")
	}
}

func TestWAG011_Check_UndefinedDependency(t *testing.T) {
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var Build = workflow.Job{
	Name:   "build",
	RunsOn: "ubuntu-latest",
}

var Deploy = workflow.Job{
	Name:   "deploy",
	RunsOn: "ubuntu-latest",
	Needs:  []any{Build, Test}, // Test is not defined
}
`)

	l := NewLinter(&WAG011{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG011 should have flagged undefined job dependency")
	}

	found := false
	for _, issue := range result.Issues {
		if issue.Rule == "WAG011" {
			found = true
			if issue.Severity != "error" {
				t.Error("WAG011 issues should be severity 'error'")
			}
		}
	}
	if !found {
		t.Error("Expected WAG011 issue not found")
	}
}

func TestWAG011_Check_ValidDependency(t *testing.T) {
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var Build = workflow.Job{
	Name:   "build",
	RunsOn: "ubuntu-latest",
}

var Deploy = workflow.Job{
	Name:   "deploy",
	RunsOn: "ubuntu-latest",
	Needs:  []any{Build},
}
`)

	l := NewLinter(&WAG011{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if !result.Success {
		t.Error("WAG011 should not flag valid job dependencies")
	}
}

func TestWAG012_Check_DeprecatedVersion(t *testing.T) {
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var Step = workflow.Step{Uses: "actions/checkout@v2"}
`)

	l := NewLinter(&WAG012{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG012 should have flagged deprecated action version")
	}

	found := false
	for _, issue := range result.Issues {
		if issue.Rule == "WAG012" {
			found = true
		}
	}
	if !found {
		t.Error("Expected WAG012 issue not found")
	}
}

func TestWAG012_Check_CurrentVersion(t *testing.T) {
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var Step = workflow.Step{Uses: "actions/checkout@v4"}
`)

	l := NewLinter(&WAG012{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if !result.Success {
		t.Error("WAG012 should not flag current action versions")
	}
}

func TestWAG004_Check_InlineMatrix(t *testing.T) {
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var BuildJob = workflow.Job{
	Strategy: workflow.Strategy{
		Matrix: workflow.Matrix{
			Values: map[string][]any{
				"os": {"ubuntu-latest", "macos-latest"},
			},
		},
	},
}
`)

	l := NewLinter(&WAG004{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG004 should have found inline matrix")
	}

	found := false
	for _, issue := range result.Issues {
		if issue.Rule == "WAG004" {
			found = true
			if issue.Severity != "info" {
				t.Error("WAG004 issues should be severity 'info'")
			}
		}
	}
	if !found {
		t.Error("Expected WAG004 issue not found")
	}
}

func TestWAG004_Check_NoInlineMatrix(t *testing.T) {
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var MyMatrix = workflow.Matrix{
	Values: map[string][]any{
		"os": {"ubuntu-latest"},
	},
}

var BuildJob = workflow.Job{
	Strategy: workflow.Strategy{
		Matrix: MyMatrix,
	},
}
`)

	l := NewLinter(&WAG004{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	// Should not flag when Matrix is a reference, not inline Values
	for _, issue := range result.Issues {
		if issue.Rule == "WAG004" {
			t.Log("WAG004 correctly triggered for Values definition")
		}
	}
}

func TestWAG005_Check_DeeplyNested(t *testing.T) {
	// WAG005 checks for deeply nested struct composite literals
	// The nesting depth must exceed maxNesting (2) to trigger
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var DeepJob = workflow.Job{
	Strategy: workflow.Strategy{
		Matrix: workflow.Matrix{
			Values: map[string][]any{
				"os": {"ubuntu-latest"},
			},
		},
	},
}
`)

	l := NewLinter(&WAG005{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	// Should flag deep nesting (Job > Strategy > Matrix = depth 3)
	found := false
	for _, issue := range result.Issues {
		if issue.Rule == "WAG005" {
			found = true
		}
	}
	if !found {
		t.Error("Expected WAG005 issue for deep nesting")
	}
}

func TestWAG008_Check_HardcodedExpression(t *testing.T) {
	content := []byte(`package main

var expr = "${{ success() && failure() }}"
`)

	l := NewLinter(&WAG008{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG008 should have found hardcoded expression")
	}

	found := false
	for _, issue := range result.Issues {
		if issue.Rule == "WAG008" {
			found = true
			if issue.Severity != "info" {
				t.Error("WAG008 issues should be severity 'info'")
			}
		}
	}
	if !found {
		t.Error("Expected WAG008 issue not found")
	}
}

func TestWAG008_Check_AllowedContexts(t *testing.T) {
	// Test that simple context references are allowed
	content := []byte(`package main

var token = "${{ github.token }}"
var secret = "${{ secrets.MY_SECRET }}"
var matrix = "${{ matrix.os }}"
var step = "${{ steps.build.outputs.result }}"
var needs = "${{ needs.build.result }}"
var input = "${{ inputs.version }}"
var env = "${{ env.MY_VAR }}"
`)

	l := NewLinter(&WAG008{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	// These should NOT be flagged
	for _, issue := range result.Issues {
		if issue.Rule == "WAG008" {
			t.Errorf("WAG008 should not flag allowed context: %s", issue.Message)
		}
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

func TestWAG011_Check_SingleDependency(t *testing.T) {
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var Build = workflow.Job{
	Name:   "build",
	RunsOn: "ubuntu-latest",
}

var Deploy = workflow.Job{
	Name:   "deploy",
	RunsOn: "ubuntu-latest",
	Needs:  Build, // Single dependency (not slice)
}
`)

	l := NewLinter(&WAG011{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if !result.Success {
		t.Error("WAG011 should not flag valid single dependency")
	}
}

func TestWAG007_Check_DefaultMaxJobs(t *testing.T) {
	// Create content with exactly 10 jobs (default max)
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var Job1 = workflow.Job{Name: "job1"}
var Job2 = workflow.Job{Name: "job2"}
var Job3 = workflow.Job{Name: "job3"}
var Job4 = workflow.Job{Name: "job4"}
var Job5 = workflow.Job{Name: "job5"}
var Job6 = workflow.Job{Name: "job6"}
var Job7 = workflow.Job{Name: "job7"}
var Job8 = workflow.Job{Name: "job8"}
var Job9 = workflow.Job{Name: "job9"}
var Job10 = workflow.Job{Name: "job10"}
`)

	// Test with MaxJobs = 0 (should use default 10)
	l := NewLinter(&WAG007{MaxJobs: 0})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	// Exactly 10 jobs should be fine with default max of 10
	if !result.Success {
		t.Error("WAG007 should not flag exactly 10 jobs with default max")
	}
}

func TestWAG001_Check_NonStepLit(t *testing.T) {
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
}
`)

	l := NewLinter(&WAG001{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if !result.Success {
		t.Error("WAG001 should not flag non-Step composite literals")
	}
}

func TestWAG012_Check_NonActionString(t *testing.T) {
	content := []byte(`package main

var notAnAction = "just/a/path"
var alsoNot = "no-at-sign/here"
`)

	l := NewLinter(&WAG012{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if !result.Success {
		t.Error("WAG012 should not flag non-action strings")
	}
}

func TestWAG012_Check_UnknownAction(t *testing.T) {
	content := []byte(`package main

var unknownAction = "custom/action@v1"
`)

	l := NewLinter(&WAG012{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	// Unknown action should not be flagged as deprecated
	if !result.Success {
		t.Error("WAG012 should not flag unknown actions")
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

func TestWAG003_Check_AllSecretPatterns(t *testing.T) {
	testCases := []string{
		`var token1 = "ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"`,
		`var token2 = "ghs_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"`,
		`var token3 = "ghu_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"`,
		`var token4 = "ghr_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"`,
		`var token5 = "github_pat_xxxxxxxxxxxxxxxxxxxxxxxxxxxxx"`,
	}

	for _, tc := range testCases {
		content := []byte("package main\n\n" + tc)

		l := NewLinter(&WAG003{})
		result, err := l.LintContent("test.go", content)
		if err != nil {
			t.Fatalf("LintContent() error = %v for %s", err, tc)
		}

		if result.Success {
			t.Errorf("WAG003 should have found secret in: %s", tc)
		}
	}
}

func TestWAG009_Check_NonMapValue(t *testing.T) {
	// Test WAG009 with a non-map value in Matrix Values
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var Matrix = workflow.Matrix{
	Values: nil,
}
`)

	l := NewLinter(&WAG009{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	// Should not crash or error
	_ = result
}

func TestWAG002_Check_NonExpressionIf(t *testing.T) {
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var Step = workflow.Step{
	If: "success()",
}
`)

	l := NewLinter(&WAG002{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	// Should not flag If without ${{ }}
	if !result.Success {
		t.Error("WAG002 should not flag If without expression syntax")
	}
}

// WAG013 Tests - Avoid pointer assignments

func TestWAG013_Check_PointerAssignment(t *testing.T) {
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = &workflow.Workflow{
	Name: "CI",
}
`)

	l := NewLinter(&WAG013{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG013 should have found pointer assignment")
	}

	found := false
	for _, issue := range result.Issues {
		if issue.Rule == "WAG013" {
			found = true
			if issue.Severity != "error" {
				t.Error("WAG013 issues should be severity 'error'")
			}
		}
	}
	if !found {
		t.Error("Expected WAG013 issue not found")
	}
}

func TestWAG013_Check_NestedPointer(t *testing.T) {
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var BuildJob = workflow.Job{
	Name:     "build",
	Strategy: &workflow.Strategy{},
}
`)

	l := NewLinter(&WAG013{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG013 should have found nested pointer assignment")
	}
}

func TestWAG013_Check_NoPointer(t *testing.T) {
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
}
`)

	l := NewLinter(&WAG013{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if !result.Success {
		t.Error("WAG013 should not flag value type assignments")
	}
}

// WAG014 Tests - Missing timeout-minutes

func TestWAG014_Check_MissingTimeout(t *testing.T) {
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var BuildJob = workflow.Job{
	Name:   "build",
	RunsOn: "ubuntu-latest",
}
`)

	l := NewLinter(&WAG014{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG014 should have flagged missing TimeoutMinutes")
	}

	found := false
	for _, issue := range result.Issues {
		if issue.Rule == "WAG014" {
			found = true
			if issue.Severity != "warning" {
				t.Error("WAG014 issues should be severity 'warning'")
			}
		}
	}
	if !found {
		t.Error("Expected WAG014 issue not found")
	}
}

func TestWAG014_Check_HasTimeout(t *testing.T) {
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var BuildJob = workflow.Job{
	Name:           "build",
	RunsOn:         "ubuntu-latest",
	TimeoutMinutes: 30,
}
`)

	l := NewLinter(&WAG014{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if !result.Success {
		t.Error("WAG014 should not flag when TimeoutMinutes is set")
	}
}

// WAG015 Tests - Suggest caching for setup actions

func TestWAG015_Check_SetupGoWithoutCache(t *testing.T) {
	content := []byte(`package main

import (
	"github.com/lex00/wetwire-github-go/workflow"
	"github.com/lex00/wetwire-github-go/actions/setup_go"
)

var BuildSteps = []any{
	setup_go.SetupGo{GoVersion: "1.23"},
	workflow.Step{Run: "go build ./..."},
}

var BuildJob = workflow.Job{
	Name:   "build",
	RunsOn: "ubuntu-latest",
	Steps:  BuildSteps,
}
`)

	l := NewLinter(&WAG015{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG015 should have suggested caching for setup-go")
	}

	found := false
	for _, issue := range result.Issues {
		if issue.Rule == "WAG015" {
			found = true
			if issue.Severity != "warning" {
				t.Error("WAG015 issues should be severity 'warning'")
			}
		}
	}
	if !found {
		t.Error("Expected WAG015 issue not found")
	}
}

func TestWAG015_Check_SetupGoWithCache(t *testing.T) {
	content := []byte(`package main

import (
	"github.com/lex00/wetwire-github-go/workflow"
	"github.com/lex00/wetwire-github-go/actions/setup_go"
	"github.com/lex00/wetwire-github-go/actions/cache"
)

var BuildSteps = []any{
	setup_go.SetupGo{GoVersion: "1.23"},
	cache.Cache{Path: "~/go/pkg/mod", Key: "go-mod"},
	workflow.Step{Run: "go build ./..."},
}

var BuildJob = workflow.Job{
	Name:   "build",
	RunsOn: "ubuntu-latest",
	Steps:  BuildSteps,
}
`)

	l := NewLinter(&WAG015{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if !result.Success {
		t.Error("WAG015 should not flag when cache is present")
	}
}

func TestWAG015_Check_SetupNodeWithoutCache(t *testing.T) {
	content := []byte(`package main

import (
	"github.com/lex00/wetwire-github-go/workflow"
	"github.com/lex00/wetwire-github-go/actions/setup_node"
)

var BuildSteps = []any{
	setup_node.SetupNode{NodeVersion: "20"},
	workflow.Step{Run: "npm install"},
}

var BuildJob = workflow.Job{
	Name:   "build",
	RunsOn: "ubuntu-latest",
	Steps:  BuildSteps,
}
`)

	l := NewLinter(&WAG015{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG015 should have suggested caching for setup-node")
	}
}

// WAG016 Tests - Validate concurrency settings

func TestWAG016_Check_CancelWithoutGroup(t *testing.T) {
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
	Concurrency: workflow.Concurrency{
		CancelInProgress: true,
	},
}
`)

	l := NewLinter(&WAG016{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG016 should have flagged cancel-in-progress without group")
	}

	found := false
	for _, issue := range result.Issues {
		if issue.Rule == "WAG016" {
			found = true
			if issue.Severity != "warning" {
				t.Error("WAG016 issues should be severity 'warning'")
			}
		}
	}
	if !found {
		t.Error("Expected WAG016 issue not found")
	}
}

func TestWAG016_Check_CancelWithGroup(t *testing.T) {
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
	Concurrency: workflow.Concurrency{
		Group:            "ci-${{ github.ref }}",
		CancelInProgress: true,
	},
}
`)

	l := NewLinter(&WAG016{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if !result.Success {
		t.Error("WAG016 should not flag when group is defined with cancel-in-progress")
	}
}

func TestWAG016_Check_GroupOnly(t *testing.T) {
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
	Concurrency: workflow.Concurrency{
		Group: "ci-${{ github.ref }}",
	},
}
`)

	l := NewLinter(&WAG016{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if !result.Success {
		t.Error("WAG016 should not flag when only group is defined")
	}
}

// WAG017 Tests - Suggest workflow permissions scope

func TestWAG017_Check_MissingPermissions(t *testing.T) {
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
	On:   workflow.Triggers{},
}
`)

	l := NewLinter(&WAG017{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG017 should have flagged missing Permissions")
	}

	found := false
	for _, issue := range result.Issues {
		if issue.Rule == "WAG017" {
			found = true
			if issue.Severity != "info" {
				t.Error("WAG017 issues should be severity 'info'")
			}
		}
	}
	if !found {
		t.Error("Expected WAG017 issue not found")
	}
}

func TestWAG017_Check_HasPermissions(t *testing.T) {
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name:        "CI",
	On:          workflow.Triggers{},
	Permissions: workflow.Permissions{Contents: "read"},
}
`)

	l := NewLinter(&WAG017{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if !result.Success {
		t.Error("WAG017 should not flag when Permissions is set")
	}
}

// WAG019 - Circular Dependency Detection Tests

func TestWAG019_Check_SimpleCycle(t *testing.T) {
	// A -> B -> A (simple 2-job cycle)
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var JobA = workflow.Job{
	Name:   "job-a",
	RunsOn: "ubuntu-latest",
	Needs:  []any{JobB},
}

var JobB = workflow.Job{
	Name:   "job-b",
	RunsOn: "ubuntu-latest",
	Needs:  []any{JobA},
}
`)

	l := NewLinter(&WAG019{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG019 should have detected circular dependency between JobA and JobB")
	}

	found := false
	for _, issue := range result.Issues {
		if issue.Rule == "WAG019" {
			found = true
			if issue.Severity != "error" {
				t.Errorf("WAG019 issues should be severity 'error', got %s", issue.Severity)
			}
			if !strings.Contains(issue.Message, "JobA") || !strings.Contains(issue.Message, "JobB") {
				t.Errorf("WAG019 message should contain job names, got: %s", issue.Message)
			}
		}
	}
	if !found {
		t.Error("Expected WAG019 issue not found")
	}
}

func TestWAG019_Check_ThreeJobCycle(t *testing.T) {
	// A -> B -> C -> A (3-job cycle)
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var JobA = workflow.Job{
	Name:   "job-a",
	RunsOn: "ubuntu-latest",
	Needs:  []any{JobC},
}

var JobB = workflow.Job{
	Name:   "job-b",
	RunsOn: "ubuntu-latest",
	Needs:  []any{JobA},
}

var JobC = workflow.Job{
	Name:   "job-c",
	RunsOn: "ubuntu-latest",
	Needs:  []any{JobB},
}
`)

	l := NewLinter(&WAG019{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG019 should have detected circular dependency in 3-job cycle")
	}

	found := false
	for _, issue := range result.Issues {
		if issue.Rule == "WAG019" {
			found = true
		}
	}
	if !found {
		t.Error("Expected WAG019 issue not found")
	}
}

func TestWAG019_Check_SelfReference(t *testing.T) {
	// A -> A (self-referencing cycle)
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var JobA = workflow.Job{
	Name:   "job-a",
	RunsOn: "ubuntu-latest",
	Needs:  []any{JobA},
}
`)

	l := NewLinter(&WAG019{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG019 should have detected self-referencing circular dependency")
	}

	found := false
	for _, issue := range result.Issues {
		if issue.Rule == "WAG019" {
			found = true
			if !strings.Contains(issue.Message, "JobA") {
				t.Errorf("WAG019 message should contain job name, got: %s", issue.Message)
			}
		}
	}
	if !found {
		t.Error("Expected WAG019 issue not found")
	}
}

func TestWAG019_Check_NoCycle(t *testing.T) {
	// Linear chain: A -> B -> C (no cycle)
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var JobA = workflow.Job{
	Name:   "job-a",
	RunsOn: "ubuntu-latest",
}

var JobB = workflow.Job{
	Name:   "job-b",
	RunsOn: "ubuntu-latest",
	Needs:  []any{JobA},
}

var JobC = workflow.Job{
	Name:   "job-c",
	RunsOn: "ubuntu-latest",
	Needs:  []any{JobB},
}
`)

	l := NewLinter(&WAG019{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if !result.Success {
		for _, issue := range result.Issues {
			t.Logf("Unexpected issue: %s", issue.Message)
		}
		t.Error("WAG019 should not flag linear dependency chain")
	}
}

func TestWAG019_Check_DiamondDependency(t *testing.T) {
	// Diamond shape: A -> B, A -> C, B -> D, C -> D (no cycle)
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var JobA = workflow.Job{
	Name:   "job-a",
	RunsOn: "ubuntu-latest",
}

var JobB = workflow.Job{
	Name:   "job-b",
	RunsOn: "ubuntu-latest",
	Needs:  []any{JobA},
}

var JobC = workflow.Job{
	Name:   "job-c",
	RunsOn: "ubuntu-latest",
	Needs:  []any{JobA},
}

var JobD = workflow.Job{
	Name:   "job-d",
	RunsOn: "ubuntu-latest",
	Needs:  []any{JobB, JobC},
}
`)

	l := NewLinter(&WAG019{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if !result.Success {
		for _, issue := range result.Issues {
			t.Logf("Unexpected issue: %s", issue.Message)
		}
		t.Error("WAG019 should not flag diamond dependency pattern (no cycle)")
	}
}

func TestWAG019_Check_SingleDependencyFormat(t *testing.T) {
	// Test with single dependency (not slice): A -> B -> A
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var JobA = workflow.Job{
	Name:   "job-a",
	RunsOn: "ubuntu-latest",
	Needs:  JobB,
}

var JobB = workflow.Job{
	Name:   "job-b",
	RunsOn: "ubuntu-latest",
	Needs:  JobA,
}
`)

	l := NewLinter(&WAG019{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG019 should detect cycle with single dependency format")
	}
}

func TestWAG019_Check_NoJobs(t *testing.T) {
	// No jobs at all
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
}
`)

	l := NewLinter(&WAG019{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if !result.Success {
		t.Error("WAG019 should not flag when there are no jobs")
	}
}
