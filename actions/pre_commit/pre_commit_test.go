package pre_commit

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestPreCommit_Action(t *testing.T) {
	a := PreCommit{}
	if got := a.Action(); got != "pre-commit/action@v3.0.1" {
		t.Errorf("Action() = %q, want %q", got, "pre-commit/action@v3.0.1")
	}
}

func TestPreCommit_Inputs_Empty(t *testing.T) {
	a := PreCommit{}
	inputs := a.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty PreCommit.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestPreCommit_ImplementsStepAction(t *testing.T) {
	a := PreCommit{}
	var _ workflow.StepAction = a
}

func TestPreCommit_Inputs_ExtraArgs(t *testing.T) {
	a := PreCommit{
		ExtraArgs: "--all-files",
	}

	inputs := a.Inputs()

	if inputs["extra_args"] != "--all-files" {
		t.Errorf("inputs[extra_args] = %v, want %q", inputs["extra_args"], "--all-files")
	}
}

func TestPreCommit_Inputs_ExtraArgs_WithConfig(t *testing.T) {
	a := PreCommit{
		ExtraArgs: "--config .pre-commit-config.yaml --all-files",
	}

	inputs := a.Inputs()

	if inputs["extra_args"] != "--config .pre-commit-config.yaml --all-files" {
		t.Errorf("inputs[extra_args] = %v, want %q", inputs["extra_args"], "--config .pre-commit-config.yaml --all-files")
	}
}

func TestPreCommit_Inputs_Token(t *testing.T) {
	a := PreCommit{
		Token: "${{ secrets.GITHUB_TOKEN }}",
	}

	inputs := a.Inputs()

	if inputs["token"] != "${{ secrets.GITHUB_TOKEN }}" {
		t.Errorf("inputs[token] = %v, want %q", inputs["token"], "${{ secrets.GITHUB_TOKEN }}")
	}
}

func TestPreCommit_Inputs_AllFields(t *testing.T) {
	a := PreCommit{
		ExtraArgs: "--all-files --verbose",
		Token:     "${{ secrets.GITHUB_TOKEN }}",
	}

	inputs := a.Inputs()

	expected := map[string]any{
		"extra_args": "--all-files --verbose",
		"token":      "${{ secrets.GITHUB_TOKEN }}",
	}

	if len(inputs) != len(expected) {
		t.Errorf("inputs has %d entries, want %d", len(inputs), len(expected))
	}

	for key, want := range expected {
		if got := inputs[key]; got != want {
			t.Errorf("inputs[%q] = %v, want %v", key, got, want)
		}
	}
}

func TestPreCommit_Inputs_CommonUsage(t *testing.T) {
	// Test common usage pattern: run all hooks on all files
	a := PreCommit{
		ExtraArgs: "--all-files",
	}

	inputs := a.Inputs()

	if len(inputs) != 1 {
		t.Errorf("common usage has %d entries, want 1", len(inputs))
	}

	if inputs["extra_args"] != "--all-files" {
		t.Errorf("inputs[extra_args] = %v, want --all-files", inputs["extra_args"])
	}
}

func TestPreCommit_Inputs_HookSelection(t *testing.T) {
	// Test selecting specific hooks
	a := PreCommit{
		ExtraArgs: "--hook-stage commit trailing-whitespace",
	}

	inputs := a.Inputs()

	if inputs["extra_args"] != "--hook-stage commit trailing-whitespace" {
		t.Errorf("inputs[extra_args] = %v, want %q", inputs["extra_args"], "--hook-stage commit trailing-whitespace")
	}
}

func TestPreCommit_Inputs_EmptyExtraArgs(t *testing.T) {
	a := PreCommit{
		ExtraArgs: "",
	}

	inputs := a.Inputs()

	if _, exists := inputs["extra_args"]; exists {
		t.Error("inputs[extra_args] should not exist for empty ExtraArgs")
	}
}
