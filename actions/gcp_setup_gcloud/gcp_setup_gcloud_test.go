package gcp_setup_gcloud

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestGCPSetupGcloud_Action(t *testing.T) {
	a := GCPSetupGcloud{}
	if got := a.Action(); got != "google-github-actions/setup-gcloud@v2" {
		t.Errorf("Action() = %q, want %q", got, "google-github-actions/setup-gcloud@v2")
	}
}

func TestGCPSetupGcloud_Inputs(t *testing.T) {
	a := GCPSetupGcloud{
		Version:   "390.0.0",
		ProjectID: "my-project",
	}

	inputs := a.Inputs()

	if inputs["version"] != "390.0.0" {
		t.Errorf("inputs[version] = %v, want 390.0.0", inputs["version"])
	}
	if inputs["project_id"] != "my-project" {
		t.Errorf("inputs[project_id] = %v, want my-project", inputs["project_id"])
	}
}

func TestGCPSetupGcloud_Inputs_Empty(t *testing.T) {
	a := GCPSetupGcloud{}
	inputs := a.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty GCPSetupGcloud.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestGCPSetupGcloud_Inputs_Components(t *testing.T) {
	a := GCPSetupGcloud{
		InstallComponents: "gke-gcloud-auth-plugin,kubectl",
	}

	inputs := a.Inputs()

	if inputs["install_components"] != "gke-gcloud-auth-plugin,kubectl" {
		t.Errorf("inputs[install_components] = %v, want components list", inputs["install_components"])
	}
}

func TestGCPSetupGcloud_Inputs_Options(t *testing.T) {
	a := GCPSetupGcloud{
		SkipInstall: true,
		Cache:       true,
	}

	inputs := a.Inputs()

	if inputs["skip_install"] != true {
		t.Errorf("inputs[skip_install] = %v, want true", inputs["skip_install"])
	}
	if inputs["cache"] != true {
		t.Errorf("inputs[cache] = %v, want true", inputs["cache"])
	}
}

func TestGCPSetupGcloud_ImplementsStepAction(t *testing.T) {
	a := GCPSetupGcloud{}
	var _ workflow.StepAction = a
}
