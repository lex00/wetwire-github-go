package reviewdog

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestReviewdog_Action(t *testing.T) {
	a := Reviewdog{}
	if got := a.Action(); got != "reviewdog/action-setup@v1" {
		t.Errorf("Action() = %q, want %q", got, "reviewdog/action-setup@v1")
	}
}

func TestReviewdog_Inputs(t *testing.T) {
	a := Reviewdog{
		ReviewdogVersion: "v0.17.0",
	}

	inputs := a.Inputs()

	if inputs["reviewdog_version"] != "v0.17.0" {
		t.Errorf("inputs[reviewdog_version] = %v, want v0.17.0", inputs["reviewdog_version"])
	}
}

func TestReviewdog_Inputs_Empty(t *testing.T) {
	a := Reviewdog{}
	inputs := a.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty Reviewdog.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestReviewdog_Inputs_Latest(t *testing.T) {
	a := Reviewdog{
		ReviewdogVersion: "latest",
	}

	inputs := a.Inputs()

	if inputs["reviewdog_version"] != "latest" {
		t.Errorf("inputs[reviewdog_version] = %v, want latest", inputs["reviewdog_version"])
	}
}

func TestReviewdog_ImplementsStepAction(t *testing.T) {
	a := Reviewdog{}
	var _ workflow.StepAction = a
}

func TestReviewdogReporter_Action(t *testing.T) {
	a := ReviewdogReporter{}
	if got := a.Action(); got != "reviewdog/action-reviewdog@v1" {
		t.Errorf("Action() = %q, want %q", got, "reviewdog/action-reviewdog@v1")
	}
}

func TestReviewdogReporter_Inputs(t *testing.T) {
	a := ReviewdogReporter{
		GithubToken: "${{ secrets.GITHUB_TOKEN }}",
		Reporter:    "github-pr-review",
		Filter:      "diff_context",
		Level:       "warning",
	}

	inputs := a.Inputs()

	if inputs["github_token"] != "${{ secrets.GITHUB_TOKEN }}" {
		t.Errorf("inputs[github_token] = %v, want secret reference", inputs["github_token"])
	}
	if inputs["reporter"] != "github-pr-review" {
		t.Errorf("inputs[reporter] = %v, want github-pr-review", inputs["reporter"])
	}
	if inputs["filter"] != "diff_context" {
		t.Errorf("inputs[filter] = %v, want diff_context", inputs["filter"])
	}
	if inputs["level"] != "warning" {
		t.Errorf("inputs[level] = %v, want warning", inputs["level"])
	}
}

func TestReviewdogReporter_Inputs_Empty(t *testing.T) {
	a := ReviewdogReporter{}
	inputs := a.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty ReviewdogReporter.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestReviewdogReporter_Inputs_FailOnError(t *testing.T) {
	a := ReviewdogReporter{
		FailOnError: true,
	}

	inputs := a.Inputs()

	if inputs["fail_on_error"] != true {
		t.Errorf("inputs[fail_on_error] = %v, want true", inputs["fail_on_error"])
	}
}

func TestReviewdogReporter_Inputs_FailOnErrorFalse(t *testing.T) {
	a := ReviewdogReporter{
		FailOnError: false,
	}

	inputs := a.Inputs()

	if _, ok := inputs["fail_on_error"]; ok {
		t.Error("fail_on_error=false should not be in inputs")
	}
}

func TestReviewdogReporter_Inputs_AllFields(t *testing.T) {
	a := ReviewdogReporter{
		GithubToken:    "${{ secrets.GITHUB_TOKEN }}",
		Workdir:        "./src",
		Reporter:       "github-pr-check",
		Filter:         "added",
		FailOnError:    true,
		Level:          "error",
		ReviewdogFlags: "-diff='git diff FETCH_HEAD'",
		Name:           "my-linter",
	}

	inputs := a.Inputs()

	expected := map[string]any{
		"github_token":    "${{ secrets.GITHUB_TOKEN }}",
		"workdir":         "./src",
		"reporter":        "github-pr-check",
		"filter":          "added",
		"fail_on_error":   true,
		"level":           "error",
		"reviewdog_flags": "-diff='git diff FETCH_HEAD'",
		"name":            "my-linter",
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

func TestReviewdogReporter_Inputs_Workdir(t *testing.T) {
	a := ReviewdogReporter{
		Workdir: "./packages/lib",
	}

	inputs := a.Inputs()

	if inputs["workdir"] != "./packages/lib" {
		t.Errorf("inputs[workdir] = %v, want ./packages/lib", inputs["workdir"])
	}
}

func TestReviewdogReporter_Inputs_ReviewdogFlags(t *testing.T) {
	a := ReviewdogReporter{
		ReviewdogFlags: "-reporter=github-check",
	}

	inputs := a.Inputs()

	if inputs["reviewdog_flags"] != "-reporter=github-check" {
		t.Errorf("inputs[reviewdog_flags] = %v, want -reporter=github-check", inputs["reviewdog_flags"])
	}
}

func TestReviewdogReporter_Inputs_Name(t *testing.T) {
	a := ReviewdogReporter{
		Name: "eslint",
	}

	inputs := a.Inputs()

	if inputs["name"] != "eslint" {
		t.Errorf("inputs[name] = %v, want eslint", inputs["name"])
	}
}

func TestReviewdogReporter_ImplementsStepAction(t *testing.T) {
	a := ReviewdogReporter{}
	var _ workflow.StepAction = a
}
