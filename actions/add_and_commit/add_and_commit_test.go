package add_and_commit

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestAddAndCommit_Action(t *testing.T) {
	a := AddAndCommit{}
	if got := a.Action(); got != "EndBug/add-and-commit@v9" {
		t.Errorf("Action() = %q, want %q", got, "EndBug/add-and-commit@v9")
	}
}

func TestAddAndCommit_Inputs_Empty(t *testing.T) {
	a := AddAndCommit{}
	inputs := a.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty AddAndCommit.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestAddAndCommit_Inputs_Add(t *testing.T) {
	a := AddAndCommit{
		Add: ".",
	}

	inputs := a.Inputs()

	if inputs["add"] != "." {
		t.Errorf("inputs[add] = %v, want %q", inputs["add"], ".")
	}
}

func TestAddAndCommit_Inputs_AuthorName(t *testing.T) {
	a := AddAndCommit{
		AuthorName: "GitHub Actions Bot",
	}

	inputs := a.Inputs()

	if inputs["author_name"] != "GitHub Actions Bot" {
		t.Errorf("inputs[author_name] = %v, want %q", inputs["author_name"], "GitHub Actions Bot")
	}
}

func TestAddAndCommit_Inputs_AuthorEmail(t *testing.T) {
	a := AddAndCommit{
		AuthorEmail: "actions@github.com",
	}

	inputs := a.Inputs()

	if inputs["author_email"] != "actions@github.com" {
		t.Errorf("inputs[author_email] = %v, want %q", inputs["author_email"], "actions@github.com")
	}
}

func TestAddAndCommit_Inputs_Message(t *testing.T) {
	a := AddAndCommit{
		Message: "chore: auto-generated changes",
	}

	inputs := a.Inputs()

	if inputs["message"] != "chore: auto-generated changes" {
		t.Errorf("inputs[message] = %v, want %q", inputs["message"], "chore: auto-generated changes")
	}
}

func TestAddAndCommit_Inputs_Push(t *testing.T) {
	a := AddAndCommit{
		Push: "true",
	}

	inputs := a.Inputs()

	if inputs["push"] != "true" {
		t.Errorf("inputs[push] = %v, want %q", inputs["push"], "true")
	}
}

func TestAddAndCommit_Inputs_Push_Branch(t *testing.T) {
	a := AddAndCommit{
		Push: "origin/main",
	}

	inputs := a.Inputs()

	if inputs["push"] != "origin/main" {
		t.Errorf("inputs[push] = %v, want %q", inputs["push"], "origin/main")
	}
}

func TestAddAndCommit_Inputs_DefaultAuthor(t *testing.T) {
	a := AddAndCommit{
		DefaultAuthor: "github_actor",
	}

	inputs := a.Inputs()

	if inputs["default_author"] != "github_actor" {
		t.Errorf("inputs[default_author] = %v, want %q", inputs["default_author"], "github_actor")
	}
}

func TestAddAndCommit_Inputs_AllFields(t *testing.T) {
	a := AddAndCommit{
		Add:           "src/ docs/",
		AuthorName:    "Bot User",
		AuthorEmail:   "bot@example.com",
		Message:       "feat: automated update",
		Push:          "false",
		DefaultAuthor: "user_info",
	}

	inputs := a.Inputs()

	if inputs["add"] != "src/ docs/" {
		t.Errorf("inputs[add] = %v, want %q", inputs["add"], "src/ docs/")
	}
	if inputs["author_name"] != "Bot User" {
		t.Errorf("inputs[author_name] = %v, want %q", inputs["author_name"], "Bot User")
	}
	if inputs["author_email"] != "bot@example.com" {
		t.Errorf("inputs[author_email] = %v, want %q", inputs["author_email"], "bot@example.com")
	}
	if inputs["message"] != "feat: automated update" {
		t.Errorf("inputs[message] = %v, want %q", inputs["message"], "feat: automated update")
	}
	if inputs["push"] != "false" {
		t.Errorf("inputs[push] = %v, want %q", inputs["push"], "false")
	}
	if inputs["default_author"] != "user_info" {
		t.Errorf("inputs[default_author] = %v, want %q", inputs["default_author"], "user_info")
	}
	if len(inputs) != 6 {
		t.Errorf("inputs has %d entries, want 6", len(inputs))
	}
}

func TestAddAndCommit_ImplementsStepAction(t *testing.T) {
	var _ workflow.StepAction = AddAndCommit{}
}
