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
	if len(l.Rules()) != 8 {
		t.Errorf("len(Rules()) = %d, want 8", len(l.Rules()))
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
	}

	expectedIDs := []string{
		"WAG001", "WAG002", "WAG003", "WAG004",
		"WAG005", "WAG006", "WAG007", "WAG008",
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
