// Package gcp_auth provides a typed wrapper for google-github-actions/auth.
package gcp_auth

// GCPAuth wraps the google-github-actions/auth@v2 action.
// Authenticate to Google Cloud using Workload Identity Federation or service account keys.
type GCPAuth struct {
	// Google Cloud project ID.
	ProjectID string `yaml:"project_id,omitempty"`

	// Full identifier of the Workload Identity Provider.
	WorkloadIdentityProvider string `yaml:"workload_identity_provider,omitempty"`

	// Email or unique ID of the service account.
	ServiceAccount string `yaml:"service_account,omitempty"`

	// Audience for GitHub OIDC token.
	Audience string `yaml:"audience,omitempty"`

	// Google Cloud JSON service account key.
	CredentialsJSON string `yaml:"credentials_json,omitempty"`

	// Generate credentials file for gcloud and SDKs.
	CreateCredentialsFile bool `yaml:"create_credentials_file,omitempty"`

	// Export environment variables like GOOGLE_CLOUD_PROJECT.
	ExportEnvironmentVariables bool `yaml:"export_environment_variables,omitempty"`

	// Output format: access_token or id_token.
	TokenFormat string `yaml:"token_format,omitempty"`

	// Additional service accounts for impersonation chain.
	Delegates string `yaml:"delegates,omitempty"`

	// Remove credentials after completion.
	CleanupCredentials bool `yaml:"cleanup_credentials,omitempty"`

	// Access token lifetime (e.g., "3600s").
	AccessTokenLifetime string `yaml:"access_token_lifetime,omitempty"`

	// OAuth 2.0 scopes for access token.
	AccessTokenScopes string `yaml:"access_token_scopes,omitempty"`

	// Email for Domain-Wide Delegation.
	AccessTokenSubject string `yaml:"access_token_subject,omitempty"`

	// Audience for ID token.
	IDTokenAudience string `yaml:"id_token_audience,omitempty"`

	// Include service account email in ID token.
	IDTokenIncludeEmail bool `yaml:"id_token_include_email,omitempty"`
}

// Action returns the action reference.
func (a GCPAuth) Action() string {
	return "google-github-actions/auth@v2"
}

// Inputs returns the action inputs as a map.
func (a GCPAuth) Inputs() map[string]any {
	with := make(map[string]any)

	if a.ProjectID != "" {
		with["project_id"] = a.ProjectID
	}
	if a.WorkloadIdentityProvider != "" {
		with["workload_identity_provider"] = a.WorkloadIdentityProvider
	}
	if a.ServiceAccount != "" {
		with["service_account"] = a.ServiceAccount
	}
	if a.Audience != "" {
		with["audience"] = a.Audience
	}
	if a.CredentialsJSON != "" {
		with["credentials_json"] = a.CredentialsJSON
	}
	if a.CreateCredentialsFile {
		with["create_credentials_file"] = a.CreateCredentialsFile
	}
	if a.ExportEnvironmentVariables {
		with["export_environment_variables"] = a.ExportEnvironmentVariables
	}
	if a.TokenFormat != "" {
		with["token_format"] = a.TokenFormat
	}
	if a.Delegates != "" {
		with["delegates"] = a.Delegates
	}
	if a.CleanupCredentials {
		with["cleanup_credentials"] = a.CleanupCredentials
	}
	if a.AccessTokenLifetime != "" {
		with["access_token_lifetime"] = a.AccessTokenLifetime
	}
	if a.AccessTokenScopes != "" {
		with["access_token_scopes"] = a.AccessTokenScopes
	}
	if a.AccessTokenSubject != "" {
		with["access_token_subject"] = a.AccessTokenSubject
	}
	if a.IDTokenAudience != "" {
		with["id_token_audience"] = a.IDTokenAudience
	}
	if a.IDTokenIncludeEmail {
		with["id_token_include_email"] = a.IDTokenIncludeEmail
	}

	return with
}
