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

func TestAWSConfigureCredentials_Inputs_SessionToken(t *testing.T) {
	a := AWSConfigureCredentials{
		AWSRegion:       "us-east-1",
		AWSSessionToken: "session-token-123",
	}

	inputs := a.Inputs()

	if inputs["aws-session-token"] != "session-token-123" {
		t.Errorf("inputs[aws-session-token] = %v, want session-token-123", inputs["aws-session-token"])
	}
}

func TestAWSConfigureCredentials_Inputs_WebIdentityToken(t *testing.T) {
	a := AWSConfigureCredentials{
		AWSRegion:            "us-east-1",
		WebIdentityTokenFile: "/path/to/token",
		Audience:             "sts.amazonaws.com",
	}

	inputs := a.Inputs()

	if inputs["web-identity-token-file"] != "/path/to/token" {
		t.Errorf("inputs[web-identity-token-file] = %v, want /path/to/token", inputs["web-identity-token-file"])
	}
	if inputs["audience"] != "sts.amazonaws.com" {
		t.Errorf("inputs[audience] = %v, want sts.amazonaws.com", inputs["audience"])
	}
}

func TestAWSConfigureCredentials_Inputs_HTTPProxy(t *testing.T) {
	a := AWSConfigureCredentials{
		AWSRegion: "us-east-1",
		HTTPProxy: "http://proxy.example.com:8080",
	}

	inputs := a.Inputs()

	if inputs["http-proxy"] != "http://proxy.example.com:8080" {
		t.Errorf("inputs[http-proxy] = %v, want http://proxy.example.com:8080", inputs["http-proxy"])
	}
}

func TestAWSConfigureCredentials_Inputs_SessionPolicies(t *testing.T) {
	a := AWSConfigureCredentials{
		AWSRegion:              "us-east-1",
		RoleToAssume:           "arn:aws:iam::123456789012:role/MyRole",
		InlineSessionPolicy:    "{\"Version\":\"2012-10-17\"}",
		ManagedSessionPolicies: "arn:aws:iam::aws:policy/ReadOnlyAccess",
	}

	inputs := a.Inputs()

	if inputs["inline-session-policy"] != "{\"Version\":\"2012-10-17\"}" {
		t.Errorf("inputs[inline-session-policy] = %v, want policy JSON", inputs["inline-session-policy"])
	}
	if inputs["managed-session-policies"] != "arn:aws:iam::aws:policy/ReadOnlyAccess" {
		t.Errorf("inputs[managed-session-policies] = %v, want policy ARN", inputs["managed-session-policies"])
	}
}

func TestAWSConfigureCredentials_Inputs_AdvancedBoolFields(t *testing.T) {
	a := AWSConfigureCredentials{
		AWSRegion:               "us-east-1",
		RoleSkipSessionTagging:  true,
		UnsetCurrentCredentials: true,
		DisableRetry:            true,
		SpecialCharactersWorkaround: true,
	}

	inputs := a.Inputs()

	if inputs["role-skip-session-tagging"] != true {
		t.Errorf("inputs[role-skip-session-tagging] = %v, want true", inputs["role-skip-session-tagging"])
	}
	if inputs["unset-current-credentials"] != true {
		t.Errorf("inputs[unset-current-credentials] = %v, want true", inputs["unset-current-credentials"])
	}
	if inputs["disable-retry"] != true {
		t.Errorf("inputs[disable-retry] = %v, want true", inputs["disable-retry"])
	}
	if inputs["special-characters-workaround"] != true {
		t.Errorf("inputs[special-characters-workaround] = %v, want true", inputs["special-characters-workaround"])
	}
}

func TestAWSConfigureCredentials_Inputs_RetryMaxAttempts(t *testing.T) {
	a := AWSConfigureCredentials{
		AWSRegion:        "us-east-1",
		RetryMaxAttempts: 5,
	}

	inputs := a.Inputs()

	if inputs["retry-max-attempts"] != 5 {
		t.Errorf("inputs[retry-max-attempts] = %v, want 5", inputs["retry-max-attempts"])
	}
}

func TestAWSConfigureCredentials_ImplementsStepAction(t *testing.T) {
	a := AWSConfigureCredentials{}
	var _ workflow.StepAction = a
}
