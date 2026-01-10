package setup_terraform

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestSetupTerraform_Action(t *testing.T) {
	s := SetupTerraform{}
	if got := s.Action(); got != "hashicorp/setup-terraform@v3" {
		t.Errorf("Action() = %q, want %q", got, "hashicorp/setup-terraform@v3")
	}
}

func TestSetupTerraform_Inputs(t *testing.T) {
	s := SetupTerraform{
		TerraformVersion: "1.7.0",
		TerraformWrapper: true,
	}

	inputs := s.Inputs()

	if inputs["terraform_version"] != "1.7.0" {
		t.Errorf("inputs[terraform_version] = %v, want %q", inputs["terraform_version"], "1.7.0")
	}

	if inputs["terraform_wrapper"] != true {
		t.Errorf("inputs[terraform_wrapper] = %v, want true", inputs["terraform_wrapper"])
	}
}

func TestSetupTerraform_Inputs_Empty(t *testing.T) {
	s := SetupTerraform{}
	inputs := s.Inputs()

	// Empty SetupTerraform should have no inputs
	if len(inputs) != 0 {
		t.Errorf("empty SetupTerraform.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestSetupTerraform_Inputs_CliConfigCredentialsHostname(t *testing.T) {
	s := SetupTerraform{
		CliConfigCredentialsHostname: "terraform.example.com",
	}

	inputs := s.Inputs()

	if inputs["cli_config_credentials_hostname"] != "terraform.example.com" {
		t.Errorf("inputs[cli_config_credentials_hostname] = %v, want %q", inputs["cli_config_credentials_hostname"], "terraform.example.com")
	}
}

func TestSetupTerraform_Inputs_CliConfigCredentialsToken(t *testing.T) {
	s := SetupTerraform{
		CliConfigCredentialsToken: "${{ secrets.TF_TOKEN }}",
	}

	inputs := s.Inputs()

	if inputs["cli_config_credentials_token"] != "${{ secrets.TF_TOKEN }}" {
		t.Errorf("inputs[cli_config_credentials_token] = %v, want %q", inputs["cli_config_credentials_token"], "${{ secrets.TF_TOKEN }}")
	}
}

func TestSetupTerraform_Inputs_TerraformVersion(t *testing.T) {
	s := SetupTerraform{
		TerraformVersion: "latest",
	}

	inputs := s.Inputs()

	if inputs["terraform_version"] != "latest" {
		t.Errorf("inputs[terraform_version] = %v, want %q", inputs["terraform_version"], "latest")
	}
}

func TestSetupTerraform_Inputs_TerraformVersionConstraint(t *testing.T) {
	s := SetupTerraform{
		TerraformVersion: "<1.13.0",
	}

	inputs := s.Inputs()

	if inputs["terraform_version"] != "<1.13.0" {
		t.Errorf("inputs[terraform_version] = %v, want %q", inputs["terraform_version"], "<1.13.0")
	}
}

func TestSetupTerraform_Inputs_TerraformWrapperTrue(t *testing.T) {
	s := SetupTerraform{
		TerraformWrapper: true,
	}

	inputs := s.Inputs()

	if inputs["terraform_wrapper"] != true {
		t.Errorf("inputs[terraform_wrapper] = %v, want true", inputs["terraform_wrapper"])
	}
}

func TestSetupTerraform_Inputs_TerraformWrapperFalse(t *testing.T) {
	// Test that false boolean values are not included
	s := SetupTerraform{
		TerraformWrapper: false,
	}

	inputs := s.Inputs()

	if len(inputs) != 0 {
		t.Errorf("inputs for false bool has %d entries, want 0. Got: %v", len(inputs), inputs)
	}

	if _, exists := inputs["terraform_wrapper"]; exists {
		t.Errorf("inputs[terraform_wrapper] should not exist for false value")
	}
}

func TestSetupTerraform_Inputs_AllFields(t *testing.T) {
	s := SetupTerraform{
		CliConfigCredentialsHostname: "app.terraform.io",
		CliConfigCredentialsToken:    "token123",
		TerraformVersion:             "1.7.0",
		TerraformWrapper:             true,
	}

	inputs := s.Inputs()

	expected := map[string]any{
		"cli_config_credentials_hostname": "app.terraform.io",
		"cli_config_credentials_token":    "token123",
		"terraform_version":               "1.7.0",
		"terraform_wrapper":               true,
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

func TestSetupTerraform_ImplementsStepAction(t *testing.T) {
	s := SetupTerraform{}
	// Verify SetupTerraform implements StepAction interface
	var _ workflow.StepAction = s
}

func TestSetupTerraform_Inputs_HostnameOnly(t *testing.T) {
	s := SetupTerraform{
		CliConfigCredentialsHostname: "custom.terraform.io",
	}

	inputs := s.Inputs()

	if inputs["cli_config_credentials_hostname"] != "custom.terraform.io" {
		t.Errorf("inputs[cli_config_credentials_hostname] = %v, want %q", inputs["cli_config_credentials_hostname"], "custom.terraform.io")
	}

	if len(inputs) != 1 {
		t.Errorf("inputs has %d entries, want 1", len(inputs))
	}
}

func TestSetupTerraform_Inputs_TokenOnly(t *testing.T) {
	s := SetupTerraform{
		CliConfigCredentialsToken: "secret-token",
	}

	inputs := s.Inputs()

	if inputs["cli_config_credentials_token"] != "secret-token" {
		t.Errorf("inputs[cli_config_credentials_token] = %v, want %q", inputs["cli_config_credentials_token"], "secret-token")
	}

	if len(inputs) != 1 {
		t.Errorf("inputs has %d entries, want 1", len(inputs))
	}
}

func TestSetupTerraform_Inputs_VersionOnly(t *testing.T) {
	s := SetupTerraform{
		TerraformVersion: "1.5.7",
	}

	inputs := s.Inputs()

	if inputs["terraform_version"] != "1.5.7" {
		t.Errorf("inputs[terraform_version] = %v, want %q", inputs["terraform_version"], "1.5.7")
	}

	if len(inputs) != 1 {
		t.Errorf("inputs has %d entries, want 1", len(inputs))
	}
}

func TestSetupTerraform_Inputs_WrapperOnly(t *testing.T) {
	s := SetupTerraform{
		TerraformWrapper: true,
	}

	inputs := s.Inputs()

	if inputs["terraform_wrapper"] != true {
		t.Errorf("inputs[terraform_wrapper] = %v, want true", inputs["terraform_wrapper"])
	}

	if len(inputs) != 1 {
		t.Errorf("inputs has %d entries, want 1", len(inputs))
	}
}

func TestSetupTerraform_Inputs_CredentialsOnly(t *testing.T) {
	s := SetupTerraform{
		CliConfigCredentialsHostname: "app.terraform.io",
		CliConfigCredentialsToken:    "${{ secrets.TF_API_TOKEN }}",
	}

	inputs := s.Inputs()

	expected := map[string]any{
		"cli_config_credentials_hostname": "app.terraform.io",
		"cli_config_credentials_token":    "${{ secrets.TF_API_TOKEN }}",
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
