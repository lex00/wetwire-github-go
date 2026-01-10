package first_interaction

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestFirstInteraction_Action(t *testing.T) {
	f := FirstInteraction{}
	if got := f.Action(); got != "actions/first-interaction@v1" {
		t.Errorf("Action() = %q, want %q", got, "actions/first-interaction@v1")
	}
}

func TestFirstInteraction_Inputs(t *testing.T) {
	f := FirstInteraction{
		RepoToken:    "${{ secrets.GITHUB_TOKEN }}",
		IssueMessage: "Welcome! Thanks for opening your first issue.",
		PRMessage:    "Thanks for your first pull request!",
	}

	inputs := f.Inputs()

	if inputs["repo-token"] != "${{ secrets.GITHUB_TOKEN }}" {
		t.Errorf("inputs[repo-token] = %v, want %q", inputs["repo-token"], "${{ secrets.GITHUB_TOKEN }}")
	}

	if inputs["issue-message"] != "Welcome! Thanks for opening your first issue." {
		t.Errorf("inputs[issue-message] = %v, want %q", inputs["issue-message"], "Welcome! Thanks for opening your first issue.")
	}

	if inputs["pr-message"] != "Thanks for your first pull request!" {
		t.Errorf("inputs[pr-message] = %v, want %q", inputs["pr-message"], "Thanks for your first pull request!")
	}
}

func TestFirstInteraction_Inputs_Empty(t *testing.T) {
	f := FirstInteraction{}
	inputs := f.Inputs()

	// Empty struct should have no inputs
	if len(inputs) != 0 {
		t.Errorf("empty FirstInteraction.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestFirstInteraction_Inputs_OnlyRepoToken(t *testing.T) {
	f := FirstInteraction{
		RepoToken: "${{ github.token }}",
	}

	inputs := f.Inputs()

	if len(inputs) != 1 {
		t.Errorf("inputs has %d entries, want 1", len(inputs))
	}

	if inputs["repo-token"] != "${{ github.token }}" {
		t.Errorf("inputs[repo-token] = %v, want %q", inputs["repo-token"], "${{ github.token }}")
	}
}

func TestFirstInteraction_Inputs_OnlyIssueMessage(t *testing.T) {
	f := FirstInteraction{
		IssueMessage: "Thanks for the issue!",
	}

	inputs := f.Inputs()

	if len(inputs) != 1 {
		t.Errorf("inputs has %d entries, want 1", len(inputs))
	}

	if inputs["issue-message"] != "Thanks for the issue!" {
		t.Errorf("inputs[issue-message] = %v, want %q", inputs["issue-message"], "Thanks for the issue!")
	}
}

func TestFirstInteraction_Inputs_OnlyPRMessage(t *testing.T) {
	f := FirstInteraction{
		PRMessage: "Thanks for the PR!",
	}

	inputs := f.Inputs()

	if len(inputs) != 1 {
		t.Errorf("inputs has %d entries, want 1", len(inputs))
	}

	if inputs["pr-message"] != "Thanks for the PR!" {
		t.Errorf("inputs[pr-message] = %v, want %q", inputs["pr-message"], "Thanks for the PR!")
	}
}

func TestFirstInteraction_ImplementsStepAction(t *testing.T) {
	f := FirstInteraction{}
	// Verify FirstInteraction implements StepAction interface
	var _ workflow.StepAction = f
}

func TestFirstInteraction_Inputs_AllFields(t *testing.T) {
	f := FirstInteraction{
		RepoToken:    "${{ secrets.GITHUB_TOKEN }}",
		IssueMessage: "Welcome to the project! Thank you for opening your first issue.",
		PRMessage:    "Thank you for your first pull request! We appreciate your contribution.",
	}

	inputs := f.Inputs()

	expected := map[string]any{
		"repo-token":    "${{ secrets.GITHUB_TOKEN }}",
		"issue-message": "Welcome to the project! Thank you for opening your first issue.",
		"pr-message":    "Thank you for your first pull request! We appreciate your contribution.",
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

func TestFirstInteraction_Inputs_MultilineMessages(t *testing.T) {
	f := FirstInteraction{
		IssueMessage: "# Welcome!\n\nThanks for opening your first issue.\n\nPlease make sure to follow our guidelines.",
		PRMessage:    "# Thanks!\n\nWe appreciate your first pull request.\n\nPlease ensure all tests pass.",
	}

	inputs := f.Inputs()

	expectedIssueMessage := "# Welcome!\n\nThanks for opening your first issue.\n\nPlease make sure to follow our guidelines."
	expectedPRMessage := "# Thanks!\n\nWe appreciate your first pull request.\n\nPlease ensure all tests pass."

	if inputs["issue-message"] != expectedIssueMessage {
		t.Errorf("inputs[issue-message] = %v, want %q", inputs["issue-message"], expectedIssueMessage)
	}

	if inputs["pr-message"] != expectedPRMessage {
		t.Errorf("inputs[pr-message] = %v, want %q", inputs["pr-message"], expectedPRMessage)
	}
}
