package gcp_auth

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestGCPAuth_Action(t *testing.T) {
	a := GCPAuth{}
	if got := a.Action(); got != "google-github-actions/auth@v2" {
		t.Errorf("Action() = %q, want %q", got, "google-github-actions/auth@v2")
	}
}

func TestGCPAuth_Inputs_WorkloadIdentity(t *testing.T) {
	a := GCPAuth{
		WorkloadIdentityProvider: "projects/123456789/locations/global/workloadIdentityPools/my-pool/providers/my-provider",
		ServiceAccount:           "my-service-account@my-project.iam.gserviceaccount.com",
	}

	inputs := a.Inputs()

	if inputs["workload_identity_provider"] == nil {
		t.Error("inputs[workload_identity_provider] should be set")
	}
	if inputs["service_account"] == nil {
		t.Error("inputs[service_account] should be set")
	}
}

func TestGCPAuth_Inputs_Empty(t *testing.T) {
	a := GCPAuth{}
	inputs := a.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty GCPAuth.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestGCPAuth_Inputs_CredentialsJSON(t *testing.T) {
	a := GCPAuth{
		CredentialsJSON: "${{ secrets.GCP_CREDENTIALS }}",
	}

	inputs := a.Inputs()

	if inputs["credentials_json"] != "${{ secrets.GCP_CREDENTIALS }}" {
		t.Errorf("inputs[credentials_json] = %v, want secret reference", inputs["credentials_json"])
	}
}

func TestGCPAuth_Inputs_TokenOptions(t *testing.T) {
	a := GCPAuth{
		WorkloadIdentityProvider: "projects/123/locations/global/workloadIdentityPools/pool/providers/prov",
		ServiceAccount:           "sa@project.iam.gserviceaccount.com",
		TokenFormat:              "access_token",
		AccessTokenLifetime:      "7200s",
	}

	inputs := a.Inputs()

	if inputs["token_format"] != "access_token" {
		t.Errorf("inputs[token_format] = %v, want access_token", inputs["token_format"])
	}
	if inputs["access_token_lifetime"] != "7200s" {
		t.Errorf("inputs[access_token_lifetime] = %v, want 7200s", inputs["access_token_lifetime"])
	}
}

func TestGCPAuth_ImplementsStepAction(t *testing.T) {
	a := GCPAuth{}
	var _ workflow.StepAction = a
}
