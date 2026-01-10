package workflow_test

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestExpressionString(t *testing.T) {
	expr := workflow.Expression("github.ref")
	expected := "${{ github.ref }}"
	if expr.String() != expected {
		t.Errorf("expected %q, got %q", expected, expr.String())
	}
}

func TestExpressionRaw(t *testing.T) {
	expr := workflow.Expression("github.ref")
	if expr.Raw() != "github.ref" {
		t.Errorf("expected 'github.ref', got %q", expr.Raw())
	}
}

func TestExpressionAnd(t *testing.T) {
	expr := workflow.Branch("main").And(workflow.Push())
	expected := "(github.ref == 'refs/heads/main') && (github.event_name == 'push')"
	if expr.Raw() != expected {
		t.Errorf("expected %q, got %q", expected, expr.Raw())
	}
}

func TestExpressionOr(t *testing.T) {
	expr := workflow.Branch("main").Or(workflow.Branch("develop"))
	expected := "(github.ref == 'refs/heads/main') || (github.ref == 'refs/heads/develop')"
	if expr.Raw() != expected {
		t.Errorf("expected %q, got %q", expected, expr.Raw())
	}
}

func TestExpressionNot(t *testing.T) {
	expr := workflow.Push().Not()
	expected := "!(github.event_name == 'push')"
	if expr.Raw() != expected {
		t.Errorf("expected %q, got %q", expected, expr.Raw())
	}
}

func TestGitHubContext(t *testing.T) {
	tests := []struct {
		name     string
		expr     workflow.Expression
		expected string
	}{
		{"Ref", workflow.GitHub.Ref(), "github.ref"},
		{"RefName", workflow.GitHub.RefName(), "github.ref_name"},
		{"SHA", workflow.GitHub.SHA(), "github.sha"},
		{"Actor", workflow.GitHub.Actor(), "github.actor"},
		{"Repository", workflow.GitHub.Repository(), "github.repository"},
		{"EventName", workflow.GitHub.EventName(), "github.event_name"},
		{"Token", workflow.GitHub.Token(), "github.token"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expr.Raw() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, tt.expr.Raw())
			}
		})
	}
}

func TestGitHubEvent(t *testing.T) {
	expr := workflow.GitHub.Event("pull_request.number")
	if expr.Raw() != "github.event.pull_request.number" {
		t.Errorf("unexpected expression: %s", expr.Raw())
	}
}

func TestRunnerContext(t *testing.T) {
	tests := []struct {
		name     string
		expr     workflow.Expression
		expected string
	}{
		{"OS", workflow.Runner.OS(), "runner.os"},
		{"Arch", workflow.Runner.Arch(), "runner.arch"},
		{"Name", workflow.Runner.Name(), "runner.name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expr.Raw() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, tt.expr.Raw())
			}
		})
	}
}

func TestSecretsContext(t *testing.T) {
	expr := workflow.Secrets.Get("DEPLOY_TOKEN")
	if expr.Raw() != "secrets.DEPLOY_TOKEN" {
		t.Errorf("unexpected expression: %s", expr.Raw())
	}

	ghToken := workflow.Secrets.GITHUB_TOKEN()
	if ghToken.Raw() != "secrets.GITHUB_TOKEN" {
		t.Errorf("unexpected expression: %s", ghToken.Raw())
	}
}

func TestMatrixContext(t *testing.T) {
	expr := workflow.MatrixContext.Get("os")
	if expr.Raw() != "matrix.os" {
		t.Errorf("unexpected expression: %s", expr.Raw())
	}
}

func TestStepsContext(t *testing.T) {
	expr := workflow.Steps.Get("checkout", "ref")
	if expr.Raw() != "steps.checkout.outputs.ref" {
		t.Errorf("unexpected expression: %s", expr.Raw())
	}

	outcome := workflow.Steps.Outcome("build")
	if outcome.Raw() != "steps.build.outcome" {
		t.Errorf("unexpected expression: %s", outcome.Raw())
	}
}

func TestNeedsContext(t *testing.T) {
	expr := workflow.Needs.Get("build", "version")
	if expr.Raw() != "needs.build.outputs.version" {
		t.Errorf("unexpected expression: %s", expr.Raw())
	}

	result := workflow.Needs.Result("build")
	if result.Raw() != "needs.build.result" {
		t.Errorf("unexpected expression: %s", result.Raw())
	}
}

func TestInputsContext(t *testing.T) {
	expr := workflow.Inputs.Get("environment")
	if expr.Raw() != "inputs.environment" {
		t.Errorf("unexpected expression: %s", expr.Raw())
	}
}

func TestVarsContext(t *testing.T) {
	expr := workflow.Vars.Get("MY_VAR")
	if expr.Raw() != "vars.MY_VAR" {
		t.Errorf("unexpected expression: %s", expr.Raw())
	}
}

func TestEnvContext(t *testing.T) {
	expr := workflow.EnvContext.Get("CI")
	if expr.Raw() != "env.CI" {
		t.Errorf("unexpected expression: %s", expr.Raw())
	}
}

func TestConditionBuilders(t *testing.T) {
	tests := []struct {
		name     string
		expr     workflow.Expression
		expected string
	}{
		{"Always", workflow.Always(), "always()"},
		{"Failure", workflow.Failure(), "failure()"},
		{"Success", workflow.Success(), "success()"},
		{"Cancelled", workflow.Cancelled(), "cancelled()"},
		{"Branch", workflow.Branch("main"), "github.ref == 'refs/heads/main'"},
		{"Tag", workflow.Tag("v1.0.0"), "github.ref == 'refs/tags/v1.0.0'"},
		{"TagPrefix", workflow.TagPrefix("v"), "startsWith(github.ref, 'refs/tags/v')"},
		{"Push", workflow.Push(), "github.event_name == 'push'"},
		{"PullRequest", workflow.PullRequest(), "github.event_name == 'pull_request'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expr.Raw() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, tt.expr.Raw())
			}
		})
	}
}

func TestOutputRef(t *testing.T) {
	step := workflow.Step{
		ID:   "checkout",
		Uses: "actions/checkout@v4",
	}

	ref := step.Output("ref")
	expected := "${{ steps.checkout.outputs.ref }}"
	if ref.String() != expected {
		t.Errorf("expected %q, got %q", expected, ref.String())
	}
}

func TestConditionInStep(t *testing.T) {
	step := workflow.Step{
		If:  workflow.Branch("main").And(workflow.Push()),
		Run: "deploy.sh",
	}

	// If field accepts any type, including Expression
	if step.If == nil {
		t.Error("expected If to be set")
	}
}

func TestExpressionInEnv(t *testing.T) {
	step := workflow.Step{
		Run: "echo $TOKEN",
		Env: workflow.Env{
			"TOKEN": workflow.Secrets.Get("DEPLOY_TOKEN"),
		},
	}

	expr, ok := step.Env["TOKEN"].(workflow.Expression)
	if !ok {
		t.Error("expected TOKEN to be an Expression")
	}
	if expr.Raw() != "secrets.DEPLOY_TOKEN" {
		t.Errorf("unexpected expression: %s", expr.Raw())
	}
}

func TestAdditionalGitHubContext(t *testing.T) {
	tests := []struct {
		name     string
		expr     workflow.Expression
		expected string
	}{
		{"RefType", workflow.GitHub.RefType(), "github.ref_type"},
		{"RepositoryOwner", workflow.GitHub.RepositoryOwner(), "github.repository_owner"},
		{"Workspace", workflow.GitHub.Workspace(), "github.workspace"},
		{"RunID", workflow.GitHub.RunID(), "github.run_id"},
		{"RunNumber", workflow.GitHub.RunNumber(), "github.run_number"},
		{"RunAttempt", workflow.GitHub.RunAttempt(), "github.run_attempt"},
		{"Job", workflow.GitHub.Job(), "github.job"},
		{"ServerURL", workflow.GitHub.ServerURL(), "github.server_url"},
		{"APIURL", workflow.GitHub.APIURL(), "github.api_url"},
		{"GraphQLURL", workflow.GitHub.GraphQLURL(), "github.graphql_url"},
		{"HeadRef", workflow.GitHub.HeadRef(), "github.head_ref"},
		{"BaseRef", workflow.GitHub.BaseRef(), "github.base_ref"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expr.Raw() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, tt.expr.Raw())
			}
		})
	}
}

func TestAdditionalRunnerContext(t *testing.T) {
	tests := []struct {
		name     string
		expr     workflow.Expression
		expected string
	}{
		{"Temp", workflow.Runner.Temp(), "runner.temp"},
		{"ToolCache", workflow.Runner.ToolCache(), "runner.tool_cache"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expr.Raw() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, tt.expr.Raw())
			}
		})
	}
}

func TestStepsContextAdditional(t *testing.T) {
	t.Run("Conclusion", func(t *testing.T) {
		expr := workflow.Steps.Conclusion("deploy")
		expected := "steps.deploy.conclusion"
		if expr.Raw() != expected {
			t.Errorf("expected %q, got %q", expected, expr.Raw())
		}
	})
}

func TestStringFunctions(t *testing.T) {
	t.Run("Contains", func(t *testing.T) {
		expr := workflow.Contains(
			workflow.GitHub.Ref(),
			workflow.Expression("'refs/heads/main'"),
		)
		// The function uses Expression.String() which wraps in ${{ }}
		expected := "contains(${{ github.ref }}, ${{ 'refs/heads/main' }})"
		if expr.Raw() != expected {
			t.Errorf("expected %q, got %q", expected, expr.Raw())
		}
	})

	t.Run("StartsWith", func(t *testing.T) {
		expr := workflow.StartsWith(
			workflow.GitHub.Ref(),
			workflow.Expression("'refs/tags/'"),
		)
		expected := "startsWith(${{ github.ref }}, ${{ 'refs/tags/' }})"
		if expr.Raw() != expected {
			t.Errorf("expected %q, got %q", expected, expr.Raw())
		}
	})

	t.Run("EndsWith", func(t *testing.T) {
		expr := workflow.EndsWith(
			workflow.GitHub.RefName(),
			workflow.Expression("'-rc'"),
		)
		expected := "endsWith(${{ github.ref_name }}, ${{ '-rc' }})"
		if expr.Raw() != expected {
			t.Errorf("expected %q, got %q", expected, expr.Raw())
		}
	})

	t.Run("Format", func(t *testing.T) {
		expr := workflow.Format(
			"Version: {0}",
			workflow.GitHub.SHA(),
		)
		expected := "format('Version: {0}', github.sha)"
		if expr.Raw() != expected {
			t.Errorf("expected %q, got %q", expected, expr.Raw())
		}
	})

	t.Run("Join", func(t *testing.T) {
		expr := workflow.Join(
			workflow.Expression("github.event.commits.*.message"),
			", ",
		)
		expected := "join(${{ github.event.commits.*.message }}, ', ')"
		if expr.Raw() != expected {
			t.Errorf("expected %q, got %q", expected, expr.Raw())
		}
	})

	t.Run("ToJSON", func(t *testing.T) {
		expr := workflow.ToJSON(workflow.GitHub.Event("pull_request"))
		expected := "toJSON(${{ github.event.pull_request }})"
		if expr.Raw() != expected {
			t.Errorf("expected %q, got %q", expected, expr.Raw())
		}
	})

	t.Run("FromJSON", func(t *testing.T) {
		expr := workflow.FromJSON(workflow.Steps.Get("config", "json"))
		expected := "fromJSON(${{ steps.config.outputs.json }})"
		if expr.Raw() != expected {
			t.Errorf("expected %q, got %q", expected, expr.Raw())
		}
	})
}

func TestComplexExpressionChaining(t *testing.T) {
	expr := workflow.Branch("main").
		And(workflow.Push()).
		Or(workflow.Branch("develop").And(workflow.PullRequest()))

	expected := "((github.ref == 'refs/heads/main') && (github.event_name == 'push')) || ((github.ref == 'refs/heads/develop') && (github.event_name == 'pull_request'))"
	if expr.Raw() != expected {
		t.Errorf("expected %q, got %q", expected, expr.Raw())
	}
}

func TestOutputRefExpression(t *testing.T) {
	step := workflow.Step{
		ID:   "build",
		Uses: "actions/setup-go@v4",
	}

	ref := step.Output("version")

	// Test String() method
	if ref.String() != "${{ steps.build.outputs.version }}" {
		t.Errorf("expected '${{ steps.build.outputs.version }}', got %q", ref.String())
	}

	// Test Expression() method
	expr := ref.Expression()
	if expr.Raw() != "steps.build.outputs.version" {
		t.Errorf("expected 'steps.build.outputs.version', got %q", expr.Raw())
	}
}
