package github_script

import (
	"testing"
)

func TestGithubScript_Action(t *testing.T) {
	g := GithubScript{}
	if got := g.Action(); got != "actions/github-script@v7" {
		t.Errorf("Action() = %q, want %q", got, "actions/github-script@v7")
	}
}

func TestGithubScript_ToStep(t *testing.T) {
	g := GithubScript{
		Script: `github.rest.issues.createComment({
			owner: context.repo.owner,
			repo: context.repo.repo,
			issue_number: context.issue.number,
			body: 'Hello from github-script!'
		})`,
		GithubToken: "${{ secrets.GITHUB_TOKEN }}",
		Debug:       true,
	}

	step := g.ToStep()

	if step.Uses != "actions/github-script@v7" {
		t.Errorf("step.Uses = %q, want %q", step.Uses, "actions/github-script@v7")
	}

	if step.With["script"] == nil {
		t.Error("step.With[script] should not be nil")
	}

	if step.With["github-token"] != "${{ secrets.GITHUB_TOKEN }}" {
		t.Errorf("step.With[github-token] = %v, want %q", step.With["github-token"], "${{ secrets.GITHUB_TOKEN }}")
	}

	if step.With["debug"] != true {
		t.Errorf("step.With[debug] = %v, want true", step.With["debug"])
	}
}

func TestGithubScript_ToStep_EmptyWithMaps(t *testing.T) {
	g := GithubScript{}
	step := g.ToStep()

	// Empty github-script should have no with values
	if len(step.With) != 0 {
		t.Errorf("empty GithubScript.ToStep() has %d with entries, want 0", len(step.With))
	}
}

func TestGithubScript_ToStep_AllFields(t *testing.T) {
	g := GithubScript{
		Script:         "console.log('test')",
		GithubToken:    "${{ secrets.MY_TOKEN }}",
		Debug:          true,
		UserAgent:      "my-custom-agent",
		Previews:       "mercy-preview",
		ResultEncoding: "json",
		Retries:        3,
		RetryExemptStatusCodes: "400,401",
	}

	step := g.ToStep()

	if step.With["script"] != "console.log('test')" {
		t.Errorf("step.With[script] = %v, want %q", step.With["script"], "console.log('test')")
	}

	if step.With["github-token"] != "${{ secrets.MY_TOKEN }}" {
		t.Errorf("step.With[github-token] = %v, want %q", step.With["github-token"], "${{ secrets.MY_TOKEN }}")
	}

	if step.With["debug"] != true {
		t.Errorf("step.With[debug] = %v, want true", step.With["debug"])
	}

	if step.With["user-agent"] != "my-custom-agent" {
		t.Errorf("step.With[user-agent] = %v, want %q", step.With["user-agent"], "my-custom-agent")
	}

	if step.With["previews"] != "mercy-preview" {
		t.Errorf("step.With[previews] = %v, want %q", step.With["previews"], "mercy-preview")
	}

	if step.With["result-encoding"] != "json" {
		t.Errorf("step.With[result-encoding] = %v, want %q", step.With["result-encoding"], "json")
	}

	if step.With["retries"] != 3 {
		t.Errorf("step.With[retries] = %v, want 3", step.With["retries"])
	}

	if step.With["retry-exempt-status-codes"] != "400,401" {
		t.Errorf("step.With[retry-exempt-status-codes] = %v, want %q", step.With["retry-exempt-status-codes"], "400,401")
	}
}

func TestGithubScript_ToStep_ScriptOnly(t *testing.T) {
	g := GithubScript{
		Script: "return 'hello'",
	}

	step := g.ToStep()

	if step.With["script"] != "return 'hello'" {
		t.Errorf("step.With[script] = %v, want %q", step.With["script"], "return 'hello'")
	}

	// Should only have 1 entry
	if len(step.With) != 1 {
		t.Errorf("GithubScript with only Script should have 1 with entry, got %d", len(step.With))
	}
}
