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
	if len(l.Rules()) != 12 {
		t.Errorf("len(Rules()) = %d, want 12", len(l.Rules()))
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
	}

	expectedIDs := []string{
		"WAG001", "WAG002", "WAG003", "WAG004",
		"WAG005", "WAG006", "WAG007", "WAG008",
		"WAG009", "WAG010", "WAG011", "WAG012",
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

var Step = setup_go.SetupGo{}.ToStep()
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

var Step = setup_go.SetupGo{GoVersion: "1.23"}.ToStep()
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
