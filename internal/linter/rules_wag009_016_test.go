package linter

import (
	"testing"
)

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
