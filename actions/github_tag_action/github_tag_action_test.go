package github_tag_action

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestGitHubTagAction_Action(t *testing.T) {
	a := GitHubTagAction{}
	if got := a.Action(); got != "anothrNick/github-tag-action@v1" {
		t.Errorf("Action() = %q, want %q", got, "anothrNick/github-tag-action@v1")
	}
}

func TestGitHubTagAction_Inputs_Empty(t *testing.T) {
	a := GitHubTagAction{}
	inputs := a.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty GitHubTagAction.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestGitHubTagAction_ImplementsStepAction(t *testing.T) {
	var _ workflow.StepAction = GitHubTagAction{}
}

func TestGitHubTagAction_Inputs_GitHubToken(t *testing.T) {
	a := GitHubTagAction{
		GitHubToken: "${{ secrets.GITHUB_TOKEN }}",
	}

	inputs := a.Inputs()

	if inputs["github_token"] != "${{ secrets.GITHUB_TOKEN }}" {
		t.Errorf("inputs[github_token] = %v, want %q", inputs["github_token"], "${{ secrets.GITHUB_TOKEN }}")
	}
}

func TestGitHubTagAction_Inputs_DefaultBump(t *testing.T) {
	a := GitHubTagAction{
		DefaultBump: "patch",
	}

	inputs := a.Inputs()

	if inputs["default_bump"] != "patch" {
		t.Errorf("inputs[default_bump] = %v, want %q", inputs["default_bump"], "patch")
	}
}

func TestGitHubTagAction_Inputs_TagPrefix(t *testing.T) {
	a := GitHubTagAction{
		TagPrefix: "v",
	}

	inputs := a.Inputs()

	if inputs["tag_prefix"] != "v" {
		t.Errorf("inputs[tag_prefix] = %v, want %q", inputs["tag_prefix"], "v")
	}
}

func TestGitHubTagAction_Inputs_DryRun(t *testing.T) {
	a := GitHubTagAction{
		DryRun: true,
	}

	inputs := a.Inputs()

	if inputs["dry_run"] != true {
		t.Errorf("inputs[dry_run] = %v, want true", inputs["dry_run"])
	}
}

func TestGitHubTagAction_Inputs_CustomTag(t *testing.T) {
	a := GitHubTagAction{
		CustomTag: "v1.2.3",
	}

	inputs := a.Inputs()

	if inputs["custom_tag"] != "v1.2.3" {
		t.Errorf("inputs[custom_tag] = %v, want %q", inputs["custom_tag"], "v1.2.3")
	}
}

func TestGitHubTagAction_Inputs_InitialVersion(t *testing.T) {
	a := GitHubTagAction{
		InitialVersion: "0.1.0",
	}

	inputs := a.Inputs()

	if inputs["initial_version"] != "0.1.0" {
		t.Errorf("inputs[initial_version] = %v, want %q", inputs["initial_version"], "0.1.0")
	}
}

func TestGitHubTagAction_Inputs_ReleasesBranches(t *testing.T) {
	a := GitHubTagAction{
		ReleasesBranches: "main,master",
	}

	inputs := a.Inputs()

	if inputs["release_branches"] != "main,master" {
		t.Errorf("inputs[release_branches] = %v, want %q", inputs["release_branches"], "main,master")
	}
}

func TestGitHubTagAction_Inputs_PrereleaseBranches(t *testing.T) {
	a := GitHubTagAction{
		PrereleaseBranches: "develop,staging",
	}

	inputs := a.Inputs()

	if inputs["prerelease_branches"] != "develop,staging" {
		t.Errorf("inputs[prerelease_branches] = %v, want %q", inputs["prerelease_branches"], "develop,staging")
	}
}

func TestGitHubTagAction_Inputs_AllFields(t *testing.T) {
	a := GitHubTagAction{
		GitHubToken:        "${{ secrets.GITHUB_TOKEN }}",
		DefaultBump:        "minor",
		TagPrefix:          "release-",
		DryRun:             true,
		CustomTag:          "custom-v2.0.0",
		InitialVersion:     "1.0.0",
		ReleasesBranches:   "main,release/*",
		PrereleaseBranches: "develop,feature/*",
	}

	inputs := a.Inputs()

	expected := map[string]any{
		"github_token":        "${{ secrets.GITHUB_TOKEN }}",
		"default_bump":        "minor",
		"tag_prefix":          "release-",
		"dry_run":             true,
		"custom_tag":          "custom-v2.0.0",
		"initial_version":     "1.0.0",
		"release_branches":    "main,release/*",
		"prerelease_branches": "develop,feature/*",
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

func TestGitHubTagAction_Inputs_FalseDryRun(t *testing.T) {
	a := GitHubTagAction{
		DryRun: false,
	}

	inputs := a.Inputs()

	if len(inputs) != 0 {
		t.Errorf("inputs for false DryRun has %d entries, want 0. Got: %v", len(inputs), inputs)
	}
}
