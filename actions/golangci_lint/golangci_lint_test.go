package golangci_lint

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestGolangciLint_Action(t *testing.T) {
	a := GolangciLint{}
	if got := a.Action(); got != "golangci/golangci-lint-action@v6" {
		t.Errorf("Action() = %q, want %q", got, "golangci/golangci-lint-action@v6")
	}
}

func TestGolangciLint_Inputs(t *testing.T) {
	a := GolangciLint{
		Version:          "v1.61",
		WorkingDirectory: "./src",
		Args:             "--timeout=5m",
	}

	inputs := a.Inputs()

	if inputs["version"] != "v1.61" {
		t.Errorf("inputs[version] = %v, want v1.61", inputs["version"])
	}
	if inputs["working-directory"] != "./src" {
		t.Errorf("inputs[working-directory] = %v, want ./src", inputs["working-directory"])
	}
	if inputs["args"] != "--timeout=5m" {
		t.Errorf("inputs[args] = %v, want --timeout=5m", inputs["args"])
	}
}

func TestGolangciLint_Inputs_Empty(t *testing.T) {
	a := GolangciLint{}
	inputs := a.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty GolangciLint.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestGolangciLint_Inputs_OnlyNewIssues(t *testing.T) {
	a := GolangciLint{
		OnlyNewIssues: true,
	}

	inputs := a.Inputs()

	if inputs["only-new-issues"] != true {
		t.Errorf("inputs[only-new-issues] = %v, want true", inputs["only-new-issues"])
	}
}

func TestGolangciLint_Inputs_CacheOptions(t *testing.T) {
	a := GolangciLint{
		SkipBuildCache: true,
		SkipPkgCache:   true,
		SkipCache:      true,
	}

	inputs := a.Inputs()

	if inputs["skip-build-cache"] != true {
		t.Errorf("inputs[skip-build-cache] = %v, want true", inputs["skip-build-cache"])
	}
	if inputs["skip-pkg-cache"] != true {
		t.Errorf("inputs[skip-pkg-cache] = %v, want true", inputs["skip-pkg-cache"])
	}
	if inputs["skip-cache"] != true {
		t.Errorf("inputs[skip-cache] = %v, want true", inputs["skip-cache"])
	}
}

func TestGolangciLint_ImplementsStepAction(t *testing.T) {
	a := GolangciLint{}
	var _ workflow.StepAction = a
}
