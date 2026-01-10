package aws_ecr_login

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestAWSECRLogin_Action(t *testing.T) {
	a := AWSECRLogin{}
	if got := a.Action(); got != "aws-actions/amazon-ecr-login@v2" {
		t.Errorf("Action() = %q, want %q", got, "aws-actions/amazon-ecr-login@v2")
	}
}

func TestAWSECRLogin_Inputs(t *testing.T) {
	a := AWSECRLogin{
		Registries: "123456789012,987654321098",
	}

	inputs := a.Inputs()

	if inputs["registries"] != "123456789012,987654321098" {
		t.Errorf("inputs[registries] = %v, want comma-separated account IDs", inputs["registries"])
	}
}

func TestAWSECRLogin_Inputs_Empty(t *testing.T) {
	a := AWSECRLogin{}
	inputs := a.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty AWSECRLogin.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestAWSECRLogin_Inputs_Public(t *testing.T) {
	a := AWSECRLogin{
		RegistryType: "public",
	}

	inputs := a.Inputs()

	if inputs["registry-type"] != "public" {
		t.Errorf("inputs[registry-type] = %v, want public", inputs["registry-type"])
	}
}

func TestAWSECRLogin_Inputs_BoolFields(t *testing.T) {
	a := AWSECRLogin{
		MaskPassword: true,
		SkipLogout:   true,
	}

	inputs := a.Inputs()

	if inputs["mask-password"] != true {
		t.Errorf("inputs[mask-password] = %v, want true", inputs["mask-password"])
	}
	if inputs["skip-logout"] != true {
		t.Errorf("inputs[skip-logout] = %v, want true", inputs["skip-logout"])
	}
}

func TestAWSECRLogin_ImplementsStepAction(t *testing.T) {
	a := AWSECRLogin{}
	var _ workflow.StepAction = a
}
