package workflow_test

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

// mockAction is a test helper implementing StepAction
type mockAction struct{}

func (mockAction) Action() string {
	return "actions/mock@v1"
}

func (mockAction) Inputs() map[string]any {
	return map[string]any{
		"param": "value",
	}
}

func TestStep(t *testing.T) {
	t.Run("basic step", func(t *testing.T) {
		step := workflow.Step{
			Name: "Build",
			Run:  "go build ./...",
		}

		if step.Name != "Build" {
			t.Errorf("expected Name='Build', got %q", step.Name)
		}

		if step.Run != "go build ./..." {
			t.Errorf("expected Run='go build ./...', got %q", step.Run)
		}
	})

	t.Run("step with ID", func(t *testing.T) {
		step := workflow.Step{
			ID:   "checkout",
			Name: "Checkout code",
			Uses: "actions/checkout@v4",
		}

		if step.ID != "checkout" {
			t.Errorf("expected ID='checkout', got %q", step.ID)
		}
	})

	t.Run("step with condition", func(t *testing.T) {
		step := workflow.Step{
			If:  workflow.Success(),
			Run: "echo 'success'",
		}

		expr, ok := step.If.(workflow.Expression)
		if !ok {
			t.Fatal("expected If to be an Expression")
		}

		if expr.Raw() != "success()" {
			t.Errorf("expected 'success()', got %q", expr.Raw())
		}
	})

	t.Run("step with env", func(t *testing.T) {
		step := workflow.Step{
			Run: "echo $VAR",
			Env: workflow.Env{
				"VAR":    "value",
				"SECRET": workflow.Secrets.Get("MY_SECRET"),
			},
		}

		if step.Env["VAR"] != "value" {
			t.Errorf("expected VAR='value', got %v", step.Env["VAR"])
		}

		secret, ok := step.Env["SECRET"].(workflow.Expression)
		if !ok {
			t.Fatal("expected SECRET to be an Expression")
		}

		if secret.Raw() != "secrets.MY_SECRET" {
			t.Errorf("expected 'secrets.MY_SECRET', got %q", secret.Raw())
		}
	})

	t.Run("step with working directory", func(t *testing.T) {
		step := workflow.Step{
			Run:              "make build",
			WorkingDirectory: "./subdir",
		}

		if step.WorkingDirectory != "./subdir" {
			t.Errorf("expected WorkingDirectory='./subdir', got %q", step.WorkingDirectory)
		}
	})

	t.Run("step with shell", func(t *testing.T) {
		step := workflow.Step{
			Run:   "echo 'test'",
			Shell: "bash",
		}

		if step.Shell != "bash" {
			t.Errorf("expected Shell='bash', got %q", step.Shell)
		}
	})

	t.Run("step with timeout", func(t *testing.T) {
		step := workflow.Step{
			Run:            "long-running-command",
			TimeoutMinutes: 30,
		}

		if step.TimeoutMinutes != 30 {
			t.Errorf("expected TimeoutMinutes=30, got %d", step.TimeoutMinutes)
		}
	})

	t.Run("step with continue on error", func(t *testing.T) {
		step := workflow.Step{
			Run:             "flaky-test",
			ContinueOnError: true,
		}

		if !step.ContinueOnError {
			t.Error("expected ContinueOnError=true")
		}
	})

	t.Run("action step with inputs", func(t *testing.T) {
		step := workflow.Step{
			Uses: "actions/checkout@v4",
			With: workflow.With{
				"fetch-depth": 0,
				"submodules":  "recursive",
				"token":       workflow.Secrets.GITHUB_TOKEN(),
			},
		}

		if step.Uses != "actions/checkout@v4" {
			t.Errorf("expected Uses='actions/checkout@v4', got %q", step.Uses)
		}

		if step.With["fetch-depth"] != 0 {
			t.Errorf("expected fetch-depth=0, got %v", step.With["fetch-depth"])
		}

		if step.With["submodules"] != "recursive" {
			t.Errorf("expected submodules='recursive', got %v", step.With["submodules"])
		}

		token, ok := step.With["token"].(workflow.Expression)
		if !ok {
			t.Fatal("expected token to be an Expression")
		}

		if token.Raw() != "secrets.GITHUB_TOKEN" {
			t.Errorf("expected 'secrets.GITHUB_TOKEN', got %q", token.Raw())
		}
	})
}

func TestStepOutput(t *testing.T) {
	t.Run("step output method", func(t *testing.T) {
		step := workflow.Step{
			ID:   "deploy",
			Run:  "echo '::set-output name=url::https://example.com'",
		}

		ref := step.Output("url")

		if ref.StepID != "deploy" {
			t.Errorf("expected StepID='deploy', got %q", ref.StepID)
		}

		if ref.Output != "url" {
			t.Errorf("expected Output='url', got %q", ref.Output)
		}

		expected := "${{ steps.deploy.outputs.url }}"
		if ref.String() != expected {
			t.Errorf("expected %q, got %q", expected, ref.String())
		}
	})
}

func TestToStep(t *testing.T) {
	// Uses mockAction defined at package level
	action := mockAction{}
	step := workflow.ToStep(action)

	if step.Uses != "actions/mock@v1" {
		t.Errorf("expected Uses='actions/mock@v1', got %q", step.Uses)
	}

	if step.With["param"] != "value" {
		t.Errorf("expected param='value', got %v", step.With["param"])
	}
}
