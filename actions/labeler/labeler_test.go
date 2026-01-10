package labeler

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestLabeler_Action(t *testing.T) {
	l := Labeler{}
	if got := l.Action(); got != "actions/labeler@v5" {
		t.Errorf("Action() = %q, want %q", got, "actions/labeler@v5")
	}
}

func TestLabeler_Inputs(t *testing.T) {
	l := Labeler{
		RepoToken:         "${{ github.token }}",
		ConfigurationPath: ".github/labeler.yml",
		SyncLabels:        true,
		Dot:               true,
		PRNumber:          123,
	}

	inputs := l.Inputs()

	if inputs["repo-token"] != "${{ github.token }}" {
		t.Errorf("inputs[repo-token] = %v, want %q", inputs["repo-token"], "${{ github.token }}")
	}

	if inputs["configuration-path"] != ".github/labeler.yml" {
		t.Errorf("inputs[configuration-path] = %v, want %q", inputs["configuration-path"], ".github/labeler.yml")
	}

	if inputs["sync-labels"] != true {
		t.Errorf("inputs[sync-labels] = %v, want true", inputs["sync-labels"])
	}

	if inputs["dot"] != true {
		t.Errorf("inputs[dot] = %v, want true", inputs["dot"])
	}

	if inputs["pr-number"] != 123 {
		t.Errorf("inputs[pr-number] = %v, want 123", inputs["pr-number"])
	}
}

func TestLabeler_Inputs_Empty(t *testing.T) {
	l := Labeler{}
	inputs := l.Inputs()

	// Empty labeler should have no inputs
	if len(inputs) != 0 {
		t.Errorf("empty Labeler.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestLabeler_Inputs_RepoTokenOnly(t *testing.T) {
	l := Labeler{
		RepoToken: "${{ secrets.GITHUB_TOKEN }}",
	}

	inputs := l.Inputs()

	if len(inputs) != 1 {
		t.Errorf("Labeler.Inputs() has %d entries, want 1", len(inputs))
	}

	if inputs["repo-token"] != "${{ secrets.GITHUB_TOKEN }}" {
		t.Errorf("inputs[repo-token] = %v, want %q", inputs["repo-token"], "${{ secrets.GITHUB_TOKEN }}")
	}
}

func TestLabeler_Inputs_ConfigurationPathOnly(t *testing.T) {
	l := Labeler{
		ConfigurationPath: ".github/custom-labeler.yml",
	}

	inputs := l.Inputs()

	if len(inputs) != 1 {
		t.Errorf("Labeler.Inputs() has %d entries, want 1", len(inputs))
	}

	if inputs["configuration-path"] != ".github/custom-labeler.yml" {
		t.Errorf("inputs[configuration-path] = %v, want %q", inputs["configuration-path"], ".github/custom-labeler.yml")
	}
}

func TestLabeler_Inputs_BoolFields(t *testing.T) {
	l := Labeler{
		SyncLabels: true,
		Dot:        true,
	}

	inputs := l.Inputs()

	if inputs["sync-labels"] != true {
		t.Errorf("inputs[sync-labels] = %v, want true", inputs["sync-labels"])
	}

	if inputs["dot"] != true {
		t.Errorf("inputs[dot] = %v, want true", inputs["dot"])
	}
}

func TestLabeler_Inputs_BoolFields_False(t *testing.T) {
	l := Labeler{
		SyncLabels: false,
		Dot:        false,
	}

	inputs := l.Inputs()

	// False boolean values should not be included in inputs
	if _, exists := inputs["sync-labels"]; exists {
		t.Errorf("inputs[sync-labels] should not exist for false value")
	}

	if _, exists := inputs["dot"]; exists {
		t.Errorf("inputs[dot] should not exist for false value")
	}
}

func TestLabeler_Inputs_PRNumber(t *testing.T) {
	l := Labeler{
		PRNumber: 456,
	}

	inputs := l.Inputs()

	if inputs["pr-number"] != 456 {
		t.Errorf("inputs[pr-number] = %v, want 456", inputs["pr-number"])
	}
}

func TestLabeler_Inputs_PRNumber_Zero(t *testing.T) {
	l := Labeler{
		PRNumber: 0,
	}

	inputs := l.Inputs()

	// Zero PR number should not be included in inputs
	if _, exists := inputs["pr-number"]; exists {
		t.Errorf("inputs[pr-number] should not exist for zero value")
	}
}

func TestLabeler_Inputs_PartialFields(t *testing.T) {
	l := Labeler{
		ConfigurationPath: ".github/labeler.yml",
		Dot:               true,
	}

	inputs := l.Inputs()

	if len(inputs) != 2 {
		t.Errorf("Labeler.Inputs() has %d entries, want 2", len(inputs))
	}

	if inputs["configuration-path"] != ".github/labeler.yml" {
		t.Errorf("inputs[configuration-path] = %v, want %q", inputs["configuration-path"], ".github/labeler.yml")
	}

	if inputs["dot"] != true {
		t.Errorf("inputs[dot] = %v, want true", inputs["dot"])
	}

	// Verify other fields are not present
	if _, exists := inputs["repo-token"]; exists {
		t.Errorf("inputs[repo-token] should not exist")
	}

	if _, exists := inputs["sync-labels"]; exists {
		t.Errorf("inputs[sync-labels] should not exist")
	}

	if _, exists := inputs["pr-number"]; exists {
		t.Errorf("inputs[pr-number] should not exist")
	}
}

func TestLabeler_ImplementsStepAction(t *testing.T) {
	l := Labeler{}
	// Verify Labeler implements StepAction interface
	var _ workflow.StepAction = l
}
