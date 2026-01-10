package codeql_analyze

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestCodeQLAnalyze_Action(t *testing.T) {
	a := CodeQLAnalyze{}
	if got := a.Action(); got != "github/codeql-action/analyze@v3" {
		t.Errorf("Action() = %q, want %q", got, "github/codeql-action/analyze@v3")
	}
}

func TestCodeQLAnalyze_Inputs(t *testing.T) {
	a := CodeQLAnalyze{
		Category: "my-analysis",
		Output:   "./results",
	}

	inputs := a.Inputs()

	if inputs["category"] != "my-analysis" {
		t.Errorf("inputs[category] = %v, want my-analysis", inputs["category"])
	}
	if inputs["output"] != "./results" {
		t.Errorf("inputs[output] = %v, want ./results", inputs["output"])
	}
}

func TestCodeQLAnalyze_Inputs_Empty(t *testing.T) {
	a := CodeQLAnalyze{}
	inputs := a.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty CodeQLAnalyze.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestCodeQLAnalyze_Inputs_BoolFields(t *testing.T) {
	a := CodeQLAnalyze{
		Upload:         true,
		UploadDatabase: true,
	}

	inputs := a.Inputs()

	if inputs["upload"] != true {
		t.Errorf("inputs[upload] = %v, want true", inputs["upload"])
	}
	if inputs["upload-database"] != true {
		t.Errorf("inputs[upload-database] = %v, want true", inputs["upload-database"])
	}
}

func TestCodeQLAnalyze_Inputs_AllFields(t *testing.T) {
	a := CodeQLAnalyze{
		Category:       "security",
		Output:         "./sarif-results",
		Upload:         true,
		UploadDatabase: true,
		CheckoutPath:   "/workspace/repo",
		RAM:            "4096",
		Threads:        "4",
	}

	inputs := a.Inputs()

	if inputs["category"] != "security" {
		t.Errorf("inputs[category] = %v, want security", inputs["category"])
	}
	if inputs["output"] != "./sarif-results" {
		t.Errorf("inputs[output] = %v, want ./sarif-results", inputs["output"])
	}
	if inputs["upload"] != true {
		t.Errorf("inputs[upload] = %v, want true", inputs["upload"])
	}
	if inputs["upload-database"] != true {
		t.Errorf("inputs[upload-database] = %v, want true", inputs["upload-database"])
	}
	if inputs["checkout-path"] != "/workspace/repo" {
		t.Errorf("inputs[checkout-path] = %v, want /workspace/repo", inputs["checkout-path"])
	}
	if inputs["ram"] != "4096" {
		t.Errorf("inputs[ram] = %v, want 4096", inputs["ram"])
	}
	if inputs["threads"] != "4" {
		t.Errorf("inputs[threads] = %v, want 4", inputs["threads"])
	}
}

func TestCodeQLAnalyze_ImplementsStepAction(t *testing.T) {
	a := CodeQLAnalyze{}
	var _ workflow.StepAction = a
}
