package deploy_pages

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestDeployPages_Action(t *testing.T) {
	d := DeployPages{}
	if got := d.Action(); got != "actions/deploy-pages@v4" {
		t.Errorf("Action() = %q, want %q", got, "actions/deploy-pages@v4")
	}
}

func TestDeployPages_Inputs_Empty(t *testing.T) {
	d := DeployPages{}
	inputs := d.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty DeployPages.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestDeployPages_Inputs_AllFields(t *testing.T) {
	d := DeployPages{
		Token:             "${{ secrets.GITHUB_TOKEN }}",
		Timeout:           600000,
		ErrorCount:        3,
		ReportingInterval: 5000,
		ArtifactName:      "github-pages",
	}

	inputs := d.Inputs()

	expected := map[string]any{
		"token":              "${{ secrets.GITHUB_TOKEN }}",
		"timeout":            600000,
		"error_count":        3,
		"reporting_interval": 5000,
		"artifact_name":      "github-pages",
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

func TestDeployPages_Inputs_Token(t *testing.T) {
	d := DeployPages{
		Token: "ghp_token123",
	}

	inputs := d.Inputs()

	if inputs["token"] != "ghp_token123" {
		t.Errorf("inputs[token] = %v, want %q", inputs["token"], "ghp_token123")
	}
	if len(inputs) != 1 {
		t.Errorf("inputs has %d entries, want 1", len(inputs))
	}
}

func TestDeployPages_Inputs_Timeout(t *testing.T) {
	d := DeployPages{
		Timeout: 300000,
	}

	inputs := d.Inputs()

	if inputs["timeout"] != 300000 {
		t.Errorf("inputs[timeout] = %v, want %d", inputs["timeout"], 300000)
	}
	if len(inputs) != 1 {
		t.Errorf("inputs has %d entries, want 1", len(inputs))
	}
}

func TestDeployPages_Inputs_ErrorCount(t *testing.T) {
	d := DeployPages{
		ErrorCount: 5,
	}

	inputs := d.Inputs()

	if inputs["error_count"] != 5 {
		t.Errorf("inputs[error_count] = %v, want %d", inputs["error_count"], 5)
	}
	if len(inputs) != 1 {
		t.Errorf("inputs has %d entries, want 1", len(inputs))
	}
}

func TestDeployPages_Inputs_ReportingInterval(t *testing.T) {
	d := DeployPages{
		ReportingInterval: 10000,
	}

	inputs := d.Inputs()

	if inputs["reporting_interval"] != 10000 {
		t.Errorf("inputs[reporting_interval] = %v, want %d", inputs["reporting_interval"], 10000)
	}
	if len(inputs) != 1 {
		t.Errorf("inputs has %d entries, want 1", len(inputs))
	}
}

func TestDeployPages_Inputs_ArtifactName(t *testing.T) {
	d := DeployPages{
		ArtifactName: "my-pages-artifact",
	}

	inputs := d.Inputs()

	if inputs["artifact_name"] != "my-pages-artifact" {
		t.Errorf("inputs[artifact_name] = %v, want %q", inputs["artifact_name"], "my-pages-artifact")
	}
	if len(inputs) != 1 {
		t.Errorf("inputs has %d entries, want 1", len(inputs))
	}
}

func TestDeployPages_ImplementsStepAction(t *testing.T) {
	d := DeployPages{}
	// Verify DeployPages implements StepAction interface
	var _ workflow.StepAction = d
}

func TestDeployPages_InSteps(t *testing.T) {
	// Test that DeployPages can be used in a steps slice
	steps := []any{
		DeployPages{
			ArtifactName: "github-pages",
			Timeout:      600000,
		},
	}

	if len(steps) != 1 {
		t.Errorf("steps has %d entries, want 1", len(steps))
	}

	dp, ok := steps[0].(DeployPages)
	if !ok {
		t.Fatal("steps[0] is not DeployPages")
	}

	if dp.ArtifactName != "github-pages" {
		t.Errorf("ArtifactName = %q, want %q", dp.ArtifactName, "github-pages")
	}
	if dp.Timeout != 600000 {
		t.Errorf("Timeout = %d, want %d", dp.Timeout, 600000)
	}
}

func TestDeployPages_Inputs_ZeroTimeout(t *testing.T) {
	// Test that Timeout = 0 is not included
	d := DeployPages{
		Timeout: 0,
	}

	inputs := d.Inputs()

	if _, exists := inputs["timeout"]; exists {
		t.Errorf("inputs[timeout] should not exist for Timeout=0")
	}
}

func TestDeployPages_Inputs_ZeroErrorCount(t *testing.T) {
	// Test that ErrorCount = 0 is not included
	d := DeployPages{
		ErrorCount: 0,
	}

	inputs := d.Inputs()

	if _, exists := inputs["error_count"]; exists {
		t.Errorf("inputs[error_count] should not exist for ErrorCount=0")
	}
}

func TestDeployPages_Inputs_ZeroReportingInterval(t *testing.T) {
	// Test that ReportingInterval = 0 is not included
	d := DeployPages{
		ReportingInterval: 0,
	}

	inputs := d.Inputs()

	if _, exists := inputs["reporting_interval"]; exists {
		t.Errorf("inputs[reporting_interval] should not exist for ReportingInterval=0")
	}
}
