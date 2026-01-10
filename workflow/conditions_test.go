package workflow_test

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestConditionBuilder(t *testing.T) {
	t.Run("And chain", func(t *testing.T) {
		cond := workflow.NewCondition(workflow.Branch("main")).
			And(workflow.Push()).
			Build()

		expected := "(github.ref == 'refs/heads/main') && (github.event_name == 'push')"
		if cond.Raw() != expected {
			t.Errorf("expected %q, got %q", expected, cond.Raw())
		}
	})

	t.Run("Or chain", func(t *testing.T) {
		cond := workflow.NewCondition(workflow.Branch("main")).
			Or(workflow.Branch("develop")).
			Build()

		expected := "(github.ref == 'refs/heads/main') || (github.ref == 'refs/heads/develop')"
		if cond.Raw() != expected {
			t.Errorf("expected %q, got %q", expected, cond.Raw())
		}
	})

	t.Run("Not", func(t *testing.T) {
		cond := workflow.NewCondition(workflow.Push()).
			Not().
			Build()

		expected := "!(github.event_name == 'push')"
		if cond.Raw() != expected {
			t.Errorf("expected %q, got %q", expected, cond.Raw())
		}
	})

	t.Run("Complex chain", func(t *testing.T) {
		cond := workflow.NewCondition(workflow.Branch("main")).
			And(workflow.Push()).
			Or(workflow.Branch("develop")).
			Build()

		expected := "((github.ref == 'refs/heads/main') && (github.event_name == 'push')) || (github.ref == 'refs/heads/develop')"
		if cond.Raw() != expected {
			t.Errorf("expected %q, got %q", expected, cond.Raw())
		}
	})
}

func TestCommonConditions(t *testing.T) {
	tests := []struct {
		name     string
		expr     workflow.Expression
		expected string
	}{
		{
			name:     "OnMainBranch",
			expr:     workflow.OnMainBranch(),
			expected: "github.ref == 'refs/heads/main'",
		},
		{
			name:     "OnDefaultBranch",
			expr:     workflow.OnDefaultBranch(),
			expected: "github.ref == format('refs/heads/{0}', github.event.repository.default_branch)",
		},
		{
			name:     "IsPullRequest",
			expr:     workflow.IsPullRequest(),
			expected: "github.event_name == 'pull_request'",
		},
		{
			name:     "IsPush",
			expr:     workflow.IsPush(),
			expected: "github.event_name == 'push'",
		},
		{
			name:     "IsRelease",
			expr:     workflow.IsRelease(),
			expected: "github.event_name == 'release'",
		},
		{
			name:     "IsTag",
			expr:     workflow.IsTag(),
			expected: "startsWith(github.ref, 'refs/tags/')",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expr.Raw() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, tt.expr.Raw())
			}
		})
	}
}

func TestPreviousJobConditions(t *testing.T) {
	t.Run("PreviousJobSucceeded", func(t *testing.T) {
		expr := workflow.PreviousJobSucceeded("build")
		expected := "(needs.build.result) && ('success')"
		if expr.Raw() != expected {
			t.Errorf("expected %q, got %q", expected, expr.Raw())
		}
	})

	t.Run("PreviousJobFailed", func(t *testing.T) {
		expr := workflow.PreviousJobFailed("build")
		expected := "needs.build.result == 'failure'"
		if expr.Raw() != expected {
			t.Errorf("expected %q, got %q", expected, expr.Raw())
		}
	})
}

func TestStringCondition(t *testing.T) {
	sc := workflow.StringCondition("custom-condition")
	if sc.String() != "custom-condition" {
		t.Errorf("expected 'custom-condition', got %q", sc.String())
	}
}
