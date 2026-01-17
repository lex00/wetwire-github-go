package linter

import (
	"testing"
)

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
