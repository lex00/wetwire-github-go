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

func TestGCPAuth_Inputs_ProjectID(t *testing.T) {
	a := GCPAuth{
		ProjectID: "my-gcp-project",
	}

	inputs := a.Inputs()

	if inputs["project_id"] != "my-gcp-project" {
		t.Errorf("inputs[project_id] = %v, want my-gcp-project", inputs["project_id"])
	}
}

func TestGCPAuth_Inputs_Audience(t *testing.T) {
	a := GCPAuth{
		Audience: "//iam.googleapis.com/projects/123/locations/global/workloadIdentityPools/pool/providers/prov",
	}

	inputs := a.Inputs()

	if inputs["audience"] != "//iam.googleapis.com/projects/123/locations/global/workloadIdentityPools/pool/providers/prov" {
		t.Errorf("inputs[audience] = %v, want audience URL", inputs["audience"])
	}
}

func TestGCPAuth_Inputs_CreateCredentialsFile(t *testing.T) {
	a := GCPAuth{
		WorkloadIdentityProvider: "projects/123/locations/global/workloadIdentityPools/pool/providers/prov",
		ServiceAccount:           "sa@project.iam.gserviceaccount.com",
		CreateCredentialsFile:    true,
	}

	inputs := a.Inputs()

	if inputs["create_credentials_file"] != true {
		t.Errorf("inputs[create_credentials_file] = %v, want true", inputs["create_credentials_file"])
	}
}

func TestGCPAuth_Inputs_CreateCredentialsFile_False(t *testing.T) {
	a := GCPAuth{
		WorkloadIdentityProvider: "projects/123/locations/global/workloadIdentityPools/pool/providers/prov",
		ServiceAccount:           "sa@project.iam.gserviceaccount.com",
		CreateCredentialsFile:    false,
	}

	inputs := a.Inputs()

	if _, exists := inputs["create_credentials_file"]; exists {
		t.Errorf("inputs[create_credentials_file] should not be set when false")
	}
}

func TestGCPAuth_Inputs_ExportEnvironmentVariables(t *testing.T) {
	a := GCPAuth{
		WorkloadIdentityProvider:   "projects/123/locations/global/workloadIdentityPools/pool/providers/prov",
		ServiceAccount:             "sa@project.iam.gserviceaccount.com",
		ExportEnvironmentVariables: true,
	}

	inputs := a.Inputs()

	if inputs["export_environment_variables"] != true {
		t.Errorf("inputs[export_environment_variables] = %v, want true", inputs["export_environment_variables"])
	}
}

func TestGCPAuth_Inputs_ExportEnvironmentVariables_False(t *testing.T) {
	a := GCPAuth{
		WorkloadIdentityProvider:   "projects/123/locations/global/workloadIdentityPools/pool/providers/prov",
		ServiceAccount:             "sa@project.iam.gserviceaccount.com",
		ExportEnvironmentVariables: false,
	}

	inputs := a.Inputs()

	if _, exists := inputs["export_environment_variables"]; exists {
		t.Errorf("inputs[export_environment_variables] should not be set when false")
	}
}

func TestGCPAuth_Inputs_Delegates(t *testing.T) {
	a := GCPAuth{
		WorkloadIdentityProvider: "projects/123/locations/global/workloadIdentityPools/pool/providers/prov",
		ServiceAccount:           "sa@project.iam.gserviceaccount.com",
		Delegates:                "delegate1@project.iam.gserviceaccount.com,delegate2@project.iam.gserviceaccount.com",
	}

	inputs := a.Inputs()

	if inputs["delegates"] != "delegate1@project.iam.gserviceaccount.com,delegate2@project.iam.gserviceaccount.com" {
		t.Errorf("inputs[delegates] = %v, want delegate list", inputs["delegates"])
	}
}

func TestGCPAuth_Inputs_CleanupCredentials(t *testing.T) {
	a := GCPAuth{
		WorkloadIdentityProvider: "projects/123/locations/global/workloadIdentityPools/pool/providers/prov",
		ServiceAccount:           "sa@project.iam.gserviceaccount.com",
		CleanupCredentials:       true,
	}

	inputs := a.Inputs()

	if inputs["cleanup_credentials"] != true {
		t.Errorf("inputs[cleanup_credentials] = %v, want true", inputs["cleanup_credentials"])
	}
}

func TestGCPAuth_Inputs_CleanupCredentials_False(t *testing.T) {
	a := GCPAuth{
		WorkloadIdentityProvider: "projects/123/locations/global/workloadIdentityPools/pool/providers/prov",
		ServiceAccount:           "sa@project.iam.gserviceaccount.com",
		CleanupCredentials:       false,
	}

	inputs := a.Inputs()

	if _, exists := inputs["cleanup_credentials"]; exists {
		t.Errorf("inputs[cleanup_credentials] should not be set when false")
	}
}

func TestGCPAuth_Inputs_AccessTokenScopes(t *testing.T) {
	a := GCPAuth{
		WorkloadIdentityProvider: "projects/123/locations/global/workloadIdentityPools/pool/providers/prov",
		ServiceAccount:           "sa@project.iam.gserviceaccount.com",
		AccessTokenScopes:        "https://www.googleapis.com/auth/cloud-platform",
	}

	inputs := a.Inputs()

	if inputs["access_token_scopes"] != "https://www.googleapis.com/auth/cloud-platform" {
		t.Errorf("inputs[access_token_scopes] = %v, want cloud-platform scope", inputs["access_token_scopes"])
	}
}

func TestGCPAuth_Inputs_AccessTokenSubject(t *testing.T) {
	a := GCPAuth{
		WorkloadIdentityProvider: "projects/123/locations/global/workloadIdentityPools/pool/providers/prov",
		ServiceAccount:           "sa@project.iam.gserviceaccount.com",
		AccessTokenSubject:       "user@example.com",
	}

	inputs := a.Inputs()

	if inputs["access_token_subject"] != "user@example.com" {
		t.Errorf("inputs[access_token_subject] = %v, want user@example.com", inputs["access_token_subject"])
	}
}

func TestGCPAuth_Inputs_IDTokenAudience(t *testing.T) {
	a := GCPAuth{
		WorkloadIdentityProvider: "projects/123/locations/global/workloadIdentityPools/pool/providers/prov",
		ServiceAccount:           "sa@project.iam.gserviceaccount.com",
		IDTokenAudience:          "https://example.com",
	}

	inputs := a.Inputs()

	if inputs["id_token_audience"] != "https://example.com" {
		t.Errorf("inputs[id_token_audience] = %v, want https://example.com", inputs["id_token_audience"])
	}
}

func TestGCPAuth_Inputs_IDTokenIncludeEmail(t *testing.T) {
	a := GCPAuth{
		WorkloadIdentityProvider: "projects/123/locations/global/workloadIdentityPools/pool/providers/prov",
		ServiceAccount:           "sa@project.iam.gserviceaccount.com",
		IDTokenIncludeEmail:      true,
	}

	inputs := a.Inputs()

	if inputs["id_token_include_email"] != true {
		t.Errorf("inputs[id_token_include_email] = %v, want true", inputs["id_token_include_email"])
	}
}

func TestGCPAuth_Inputs_IDTokenIncludeEmail_False(t *testing.T) {
	a := GCPAuth{
		WorkloadIdentityProvider: "projects/123/locations/global/workloadIdentityPools/pool/providers/prov",
		ServiceAccount:           "sa@project.iam.gserviceaccount.com",
		IDTokenIncludeEmail:      false,
	}

	inputs := a.Inputs()

	if _, exists := inputs["id_token_include_email"]; exists {
		t.Errorf("inputs[id_token_include_email] should not be set when false")
	}
}

func TestGCPAuth_Inputs_AllFields(t *testing.T) {
	a := GCPAuth{
		ProjectID:                  "my-project",
		WorkloadIdentityProvider:   "projects/123/locations/global/workloadIdentityPools/pool/providers/prov",
		ServiceAccount:             "sa@project.iam.gserviceaccount.com",
		Audience:                   "//iam.googleapis.com/projects/123/locations/global/workloadIdentityPools/pool/providers/prov",
		CredentialsJSON:            "${{ secrets.GCP_CREDENTIALS }}",
		CreateCredentialsFile:      true,
		ExportEnvironmentVariables: true,
		TokenFormat:                "id_token",
		Delegates:                  "delegate@project.iam.gserviceaccount.com",
		CleanupCredentials:         true,
		AccessTokenLifetime:        "3600s",
		AccessTokenScopes:          "https://www.googleapis.com/auth/cloud-platform",
		AccessTokenSubject:         "user@example.com",
		IDTokenAudience:            "https://example.com",
		IDTokenIncludeEmail:        true,
	}

	inputs := a.Inputs()

	expectedFields := []string{
		"project_id",
		"workload_identity_provider",
		"service_account",
		"audience",
		"credentials_json",
		"create_credentials_file",
		"export_environment_variables",
		"token_format",
		"delegates",
		"cleanup_credentials",
		"access_token_lifetime",
		"access_token_scopes",
		"access_token_subject",
		"id_token_audience",
		"id_token_include_email",
	}

	if len(inputs) != len(expectedFields) {
		t.Errorf("inputs has %d fields, want %d", len(inputs), len(expectedFields))
	}

	for _, field := range expectedFields {
		if _, exists := inputs[field]; !exists {
			t.Errorf("inputs[%s] should be set", field)
		}
	}
}
