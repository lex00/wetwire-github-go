package github_script

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestGithubScript_Action(t *testing.T) {
	g := GithubScript{}
	if got := g.Action(); got != "actions/github-script@v7" {
		t.Errorf("Action() = %q, want %q", got, "actions/github-script@v7")
	}
}

func TestGithubScript_Inputs(t *testing.T) {
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

	inputs := g.Inputs()

	if g.Action() != "actions/github-script@v7" {
		t.Errorf("Action() = %q, want %q", g.Action(), "actions/github-script@v7")
	}

	if inputs["script"] == nil {
		t.Error("inputs[script] should not be nil")
	}

	if inputs["github-token"] != "${{ secrets.GITHUB_TOKEN }}" {
		t.Errorf("inputs[github-token] = %v, want %q", inputs["github-token"], "${{ secrets.GITHUB_TOKEN }}")
	}

	if inputs["debug"] != true {
		t.Errorf("inputs[debug] = %v, want true", inputs["debug"])
	}
}

func TestGithubScript_Inputs_EmptyWithMaps(t *testing.T) {
	g := GithubScript{}
	inputs := g.Inputs()

	// Empty github-script should have no inputs values
	if len(inputs) != 0 {
		t.Errorf("empty GithubScript.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestGithubScript_Inputs_AllFields(t *testing.T) {
	g := GithubScript{
		Script:                 "console.log('test')",
		GithubToken:           "${{ secrets.MY_TOKEN }}",
		Debug:                  true,
		UserAgent:              "my-custom-agent",
		Previews:               "mercy-preview",
		ResultEncoding:         "json",
		Retries:                3,
		RetryExemptStatusCodes: "400,401",
	}

	inputs := g.Inputs()

	if inputs["script"] != "console.log('test')" {
		t.Errorf("inputs[script] = %v, want %q", inputs["script"], "console.log('test')")
	}

	if inputs["github-token"] != "${{ secrets.MY_TOKEN }}" {
		t.Errorf("inputs[github-token] = %v, want %q", inputs["github-token"], "${{ secrets.MY_TOKEN }}")
	}

	if inputs["debug"] != true {
		t.Errorf("inputs[debug] = %v, want true", inputs["debug"])
	}

	if inputs["user-agent"] != "my-custom-agent" {
		t.Errorf("inputs[user-agent] = %v, want %q", inputs["user-agent"], "my-custom-agent")
	}

	if inputs["previews"] != "mercy-preview" {
		t.Errorf("inputs[previews] = %v, want %q", inputs["previews"], "mercy-preview")
	}

	if inputs["result-encoding"] != "json" {
		t.Errorf("inputs[result-encoding] = %v, want %q", inputs["result-encoding"], "json")
	}

	if inputs["retries"] != 3 {
		t.Errorf("inputs[retries] = %v, want 3", inputs["retries"])
	}

	if inputs["retry-exempt-status-codes"] != "400,401" {
		t.Errorf("inputs[retry-exempt-status-codes] = %v, want %q", inputs["retry-exempt-status-codes"], "400,401")
	}
}

func TestGithubScript_Inputs_ScriptOnly(t *testing.T) {
	g := GithubScript{
		Script: "return 'hello'",
	}

	inputs := g.Inputs()

	if inputs["script"] != "return 'hello'" {
		t.Errorf("inputs[script] = %v, want %q", inputs["script"], "return 'hello'")
	}

	// Should only have 1 entry
	if len(inputs) != 1 {
		t.Errorf("GithubScript with only Script should have 1 input entry, got %d", len(inputs))
	}
}

func TestGithubScript_ImplementsStepAction(t *testing.T) {
	var _ workflow.StepAction = GithubScript{}
}
