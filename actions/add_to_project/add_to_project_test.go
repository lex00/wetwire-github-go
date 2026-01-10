package add_to_project

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestAddToProject_Action(t *testing.T) {
	a := AddToProject{}
	if got := a.Action(); got != "actions/add-to-project@v1" {
		t.Errorf("Action() = %q, want %q", got, "actions/add-to-project@v1")
	}
}

func TestAddToProject_Inputs(t *testing.T) {
	a := AddToProject{
		ProjectURL:  "https://github.com/orgs/myorg/projects/1",
		GithubToken: "${{ secrets.ADD_TO_PROJECT_PAT }}",
		Labeled:     "bug,needs-triage",
		LabelOperator: "OR",
	}

	inputs := a.Inputs()

	if inputs["project-url"] != "https://github.com/orgs/myorg/projects/1" {
		t.Errorf("inputs[project-url] = %v, want %q", inputs["project-url"], "https://github.com/orgs/myorg/projects/1")
	}

	if inputs["github-token"] != "${{ secrets.ADD_TO_PROJECT_PAT }}" {
		t.Errorf("inputs[github-token] = %v, want %q", inputs["github-token"], "${{ secrets.ADD_TO_PROJECT_PAT }}")
	}

	if inputs["labeled"] != "bug,needs-triage" {
		t.Errorf("inputs[labeled] = %v, want %q", inputs["labeled"], "bug,needs-triage")
	}

	if inputs["label-operator"] != "OR" {
		t.Errorf("inputs[label-operator] = %v, want %q", inputs["label-operator"], "OR")
	}
}

func TestAddToProject_Inputs_Empty(t *testing.T) {
	a := AddToProject{}
	inputs := a.Inputs()

	// Empty struct should have no inputs
	if len(inputs) != 0 {
		t.Errorf("empty AddToProject.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestAddToProject_Inputs_RequiredOnly(t *testing.T) {
	a := AddToProject{
		ProjectURL:  "https://github.com/orgs/test/projects/5",
		GithubToken: "${{ secrets.GITHUB_TOKEN }}",
	}

	inputs := a.Inputs()

	if len(inputs) != 2 {
		t.Errorf("inputs has %d entries, want 2", len(inputs))
	}

	if inputs["project-url"] != "https://github.com/orgs/test/projects/5" {
		t.Errorf("inputs[project-url] = %v, want %q", inputs["project-url"], "https://github.com/orgs/test/projects/5")
	}

	if inputs["github-token"] != "${{ secrets.GITHUB_TOKEN }}" {
		t.Errorf("inputs[github-token] = %v, want %q", inputs["github-token"], "${{ secrets.GITHUB_TOKEN }}")
	}
}

func TestAddToProject_Inputs_WithLabeled(t *testing.T) {
	a := AddToProject{
		Labeled: "enhancement,feature",
	}

	inputs := a.Inputs()

	if inputs["labeled"] != "enhancement,feature" {
		t.Errorf("inputs[labeled] = %v, want %q", inputs["labeled"], "enhancement,feature")
	}
}

func TestAddToProject_Inputs_LabelOperatorAND(t *testing.T) {
	a := AddToProject{
		LabelOperator: "AND",
	}

	inputs := a.Inputs()

	if inputs["label-operator"] != "AND" {
		t.Errorf("inputs[label-operator] = %v, want %q", inputs["label-operator"], "AND")
	}
}

func TestAddToProject_Inputs_LabelOperatorNOT(t *testing.T) {
	a := AddToProject{
		LabelOperator: "NOT",
	}

	inputs := a.Inputs()

	if inputs["label-operator"] != "NOT" {
		t.Errorf("inputs[label-operator] = %v, want %q", inputs["label-operator"], "NOT")
	}
}

func TestAddToProject_ImplementsStepAction(t *testing.T) {
	a := AddToProject{}
	// Verify AddToProject implements StepAction interface
	var _ workflow.StepAction = a
}

func TestAddToProject_Inputs_AllFields(t *testing.T) {
	a := AddToProject{
		ProjectURL:    "https://github.com/users/username/projects/10",
		GithubToken:   "${{ secrets.PROJECT_TOKEN }}",
		Labeled:       "priority-high,bug",
		LabelOperator: "AND",
	}

	inputs := a.Inputs()

	expected := map[string]any{
		"project-url":    "https://github.com/users/username/projects/10",
		"github-token":   "${{ secrets.PROJECT_TOKEN }}",
		"labeled":        "priority-high,bug",
		"label-operator": "AND",
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

func TestAddToProject_Inputs_OrgProject(t *testing.T) {
	a := AddToProject{
		ProjectURL: "https://github.com/orgs/my-org/projects/42",
	}

	inputs := a.Inputs()

	if inputs["project-url"] != "https://github.com/orgs/my-org/projects/42" {
		t.Errorf("inputs[project-url] = %v, want %q", inputs["project-url"], "https://github.com/orgs/my-org/projects/42")
	}
}

func TestAddToProject_Inputs_UserProject(t *testing.T) {
	a := AddToProject{
		ProjectURL: "https://github.com/users/johndoe/projects/1",
	}

	inputs := a.Inputs()

	if inputs["project-url"] != "https://github.com/users/johndoe/projects/1" {
		t.Errorf("inputs[project-url] = %v, want %q", inputs["project-url"], "https://github.com/users/johndoe/projects/1")
	}
}

func TestAddToProject_Inputs_SingleLabel(t *testing.T) {
	a := AddToProject{
		Labeled: "bug",
	}

	inputs := a.Inputs()

	if inputs["labeled"] != "bug" {
		t.Errorf("inputs[labeled] = %v, want %q", inputs["labeled"], "bug")
	}
}

func TestAddToProject_Inputs_MultipleLabels(t *testing.T) {
	a := AddToProject{
		Labeled: "bug,enhancement,documentation,good-first-issue",
	}

	inputs := a.Inputs()

	if inputs["labeled"] != "bug,enhancement,documentation,good-first-issue" {
		t.Errorf("inputs[labeled] = %v, want %q", inputs["labeled"], "bug,enhancement,documentation,good-first-issue")
	}
}
