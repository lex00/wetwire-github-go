package lint

import "testing"

// WAG018 Tests - Detect dangerous pull_request_target patterns

func TestWAG018_Check_PullRequestTargetWithCheckout(t *testing.T) {
	content := []byte(`package main

import (
	"github.com/lex00/wetwire-github-go/workflow"
	"github.com/lex00/wetwire-github-go/actions/checkout"
)

var PRTargetTriggers = workflow.Triggers{
	PullRequestTarget: &workflow.PullRequestTargetTrigger{},
}

var DangerousWorkflow = workflow.Workflow{
	Name: "PR Target",
	On:   PRTargetTriggers,
	Jobs: map[string]workflow.Job{
		"test": {
			RunsOn: "ubuntu-latest",
			Steps: []any{
				checkout.Checkout{},
			},
		},
	},
}
`)

	l := NewLinter(&WAG018{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG018 should have flagged pull_request_target with checkout action")
	}

	found := false
	for _, issue := range result.Issues {
		if issue.Rule == "WAG018" {
			found = true
			if issue.Severity != SeverityWarning {
				t.Error("WAG018 issues should be severity 'warning'")
			}
		}
	}
	if !found {
		t.Error("Expected WAG018 issue not found")
	}
}

func TestWAG018_Check_PullRequestTargetWithoutCheckout(t *testing.T) {
	content := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var PRTargetTriggers = workflow.Triggers{
	PullRequestTarget: &workflow.PullRequestTargetTrigger{},
}

var SafeWorkflow = workflow.Workflow{
	Name: "PR Target",
	On:   PRTargetTriggers,
	Jobs: map[string]workflow.Job{
		"test": {
			RunsOn: "ubuntu-latest",
			Steps: []any{
				workflow.Step{Run: "echo 'No checkout'"},
			},
		},
	},
}
`)

	l := NewLinter(&WAG018{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	// Should not flag when there's no checkout action
	if !result.Success {
		t.Error("WAG018 should not flag pull_request_target without checkout")
	}
}

func TestWAG018_Check_PullRequestWithCheckout(t *testing.T) {
	content := []byte(`package main

import (
	"github.com/lex00/wetwire-github-go/workflow"
	"github.com/lex00/wetwire-github-go/actions/checkout"
)

var PRTriggers = workflow.Triggers{
	PullRequest: &workflow.PullRequestTrigger{},
}

var SafeWorkflow = workflow.Workflow{
	Name: "PR",
	On:   PRTriggers,
	Jobs: map[string]workflow.Job{
		"test": {
			RunsOn: "ubuntu-latest",
			Steps: []any{
				checkout.Checkout{},
			},
		},
	},
}
`)

	l := NewLinter(&WAG018{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	// Should not flag pull_request (only pull_request_target is dangerous)
	if !result.Success {
		t.Error("WAG018 should not flag pull_request with checkout")
	}
}

func TestWAG018_Check_InlinePullRequestTarget(t *testing.T) {
	content := []byte(`package main

import (
	"github.com/lex00/wetwire-github-go/workflow"
	"github.com/lex00/wetwire-github-go/actions/checkout"
)

var DangerousWorkflow = workflow.Workflow{
	Name: "PR Target",
	On: workflow.Triggers{
		PullRequestTarget: &workflow.PullRequestTargetTrigger{},
	},
	Jobs: map[string]workflow.Job{
		"test": {
			RunsOn: "ubuntu-latest",
			Steps: []any{
				checkout.Checkout{},
			},
		},
	},
}
`)

	l := NewLinter(&WAG018{})
	result, err := l.LintContent("test.go", content)
	if err != nil {
		t.Fatalf("LintContent() error = %v", err)
	}

	if result.Success {
		t.Error("WAG018 should have flagged inline pull_request_target with checkout")
	}
}
