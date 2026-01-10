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

func TestGolangciLint_Inputs_ProblemMatchers(t *testing.T) {
	a := GolangciLint{
		Version:         "v1.61",
		ProblemMatchers: true,
	}

	inputs := a.Inputs()

	if inputs["problem-matchers"] != true {
		t.Errorf("inputs[problem-matchers] = %v, want true", inputs["problem-matchers"])
	}
}

func TestGolangciLint_Inputs_GithubToken(t *testing.T) {
	a := GolangciLint{
		Version:     "v1.61",
		GithubToken: "${{ secrets.GITHUB_TOKEN }}",
	}

	inputs := a.Inputs()

	if inputs["github-token"] != "${{ secrets.GITHUB_TOKEN }}" {
		t.Errorf("inputs[github-token] = %v, want secret reference", inputs["github-token"])
	}
}

func TestGolangciLint_Inputs_InstallMode(t *testing.T) {
	a := GolangciLint{
		Version:     "v1.61",
		InstallMode: "binary",
	}

	inputs := a.Inputs()

	if inputs["install-mode"] != "binary" {
		t.Errorf("inputs[install-mode] = %v, want binary", inputs["install-mode"])
	}
}

func TestGolangciLint_Inputs_GoModules(t *testing.T) {
	a := GolangciLint{
		Version:   "v1.61",
		GoModules: true,
	}

	inputs := a.Inputs()

	if inputs["go-modules"] != true {
		t.Errorf("inputs[go-modules] = %v, want true", inputs["go-modules"])
	}
}

func TestGolangciLint_Inputs_BooleanFalse(t *testing.T) {
	a := GolangciLint{
		Version:         "v1.61",
		OnlyNewIssues:   false,
		SkipBuildCache:  false,
		SkipPkgCache:    false,
		ProblemMatchers: false,
		GoModules:       false,
		SkipCache:       false,
	}

	inputs := a.Inputs()

	if _, ok := inputs["only-new-issues"]; ok {
		t.Error("only-new-issues=false should not be in inputs")
	}
	if _, ok := inputs["skip-build-cache"]; ok {
		t.Error("skip-build-cache=false should not be in inputs")
	}
	if _, ok := inputs["skip-pkg-cache"]; ok {
		t.Error("skip-pkg-cache=false should not be in inputs")
	}
	if _, ok := inputs["problem-matchers"]; ok {
		t.Error("problem-matchers=false should not be in inputs")
	}
	if _, ok := inputs["go-modules"]; ok {
		t.Error("go-modules=false should not be in inputs")
	}
	if _, ok := inputs["skip-cache"]; ok {
		t.Error("skip-cache=false should not be in inputs")
	}
}

func TestGolangciLint_ImplementsStepAction(t *testing.T) {
	a := GolangciLint{}
	var _ workflow.StepAction = a
}
