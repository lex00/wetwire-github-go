package aws_configure_credentials

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestAWSConfigureCredentials_Action(t *testing.T) {
	a := AWSConfigureCredentials{}
	if got := a.Action(); got != "aws-actions/configure-aws-credentials@v4" {
		t.Errorf("Action() = %q, want %q", got, "aws-actions/configure-aws-credentials@v4")
	}
}

func TestAWSConfigureCredentials_Inputs(t *testing.T) {
	a := AWSConfigureCredentials{
		AWSRegion:    "us-east-1",
		RoleToAssume: "arn:aws:iam::123456789012:role/GitHubActionsRole",
	}

	inputs := a.Inputs()

	if inputs["aws-region"] != "us-east-1" {
		t.Errorf("inputs[aws-region] = %v, want %q", inputs["aws-region"], "us-east-1")
	}
	if inputs["role-to-assume"] != "arn:aws:iam::123456789012:role/GitHubActionsRole" {
		t.Errorf("inputs[role-to-assume] = %v, want role ARN", inputs["role-to-assume"])
	}
}

func TestAWSConfigureCredentials_Inputs_Empty(t *testing.T) {
	a := AWSConfigureCredentials{}
	inputs := a.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty AWSConfigureCredentials.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestAWSConfigureCredentials_Inputs_WithAccessKeys(t *testing.T) {
	a := AWSConfigureCredentials{
		AWSRegion:          "us-west-2",
		AWSAccessKeyID:     "${{ secrets.AWS_ACCESS_KEY_ID }}",
		AWSSecretAccessKey: "${{ secrets.AWS_SECRET_ACCESS_KEY }}",
		RoleToAssume:       "arn:aws:iam::123456789012:role/DeployRole",
	}

	inputs := a.Inputs()

	if inputs["aws-access-key-id"] != "${{ secrets.AWS_ACCESS_KEY_ID }}" {
		t.Errorf("inputs[aws-access-key-id] = %v, want secret reference", inputs["aws-access-key-id"])
	}
	if inputs["aws-secret-access-key"] != "${{ secrets.AWS_SECRET_ACCESS_KEY }}" {
		t.Errorf("inputs[aws-secret-access-key] = %v, want secret reference", inputs["aws-secret-access-key"])
	}
}

func TestAWSConfigureCredentials_Inputs_WithRoleOptions(t *testing.T) {
	a := AWSConfigureCredentials{
		AWSRegion:           "eu-west-1",
		RoleToAssume:        "arn:aws:iam::123456789012:role/MyRole",
		RoleDurationSeconds: 3600,
		RoleSessionName:     "GitHubDeploy",
		RoleExternalID:      "my-external-id",
	}

	inputs := a.Inputs()

	if inputs["role-duration-seconds"] != 3600 {
		t.Errorf("inputs[role-duration-seconds] = %v, want 3600", inputs["role-duration-seconds"])
	}
	if inputs["role-session-name"] != "GitHubDeploy" {
		t.Errorf("inputs[role-session-name] = %v, want GitHubDeploy", inputs["role-session-name"])
	}
	if inputs["role-external-id"] != "my-external-id" {
		t.Errorf("inputs[role-external-id] = %v, want my-external-id", inputs["role-external-id"])
	}
}

func TestAWSConfigureCredentials_Inputs_BoolFields(t *testing.T) {
	a := AWSConfigureCredentials{
		AWSRegion:        "us-east-1",
		RoleChaining:     true,
		OutputCredentials: true,
		MaskAWSAccountID: true,
	}

	inputs := a.Inputs()

	if inputs["role-chaining"] != true {
		t.Errorf("inputs[role-chaining] = %v, want true", inputs["role-chaining"])
	}
	if inputs["output-credentials"] != true {
		t.Errorf("inputs[output-credentials] = %v, want true", inputs["output-credentials"])
	}
	if inputs["mask-aws-account-id"] != true {
		t.Errorf("inputs[mask-aws-account-id] = %v, want true", inputs["mask-aws-account-id"])
	}
}

func TestAWSConfigureCredentials_ImplementsStepAction(t *testing.T) {
	a := AWSConfigureCredentials{}
	var _ workflow.StepAction = a
}
