// Package azure_login provides a typed wrapper for azure/login.
package azure_login

// AzureLogin wraps the azure/login@v2 action.
// Login to Azure using service principal or OIDC.
type AzureLogin struct {
	// Azure credentials JSON from `az ad sp create-for-rbac`.
	Creds string `yaml:"creds,omitempty"`

	// Client ID for Azure service principal (OIDC).
	ClientID string `yaml:"client-id,omitempty"`

	// Tenant ID for Azure service principal.
	TenantID string `yaml:"tenant-id,omitempty"`

	// Azure subscription ID.
	SubscriptionID string `yaml:"subscription-id,omitempty"`

	// Enable Azure PowerShell session alongside CLI.
	EnableAzPSSession bool `yaml:"enable-AzPSSession,omitempty"`

	// Azure environment (azurecloud, azureusgovernment, azurechinacloud, etc.).
	Environment string `yaml:"environment,omitempty"`

	// Allow login without subscriptions (tenant-level access).
	AllowNoSubscriptions bool `yaml:"allow-no-subscriptions,omitempty"`

	// Token audience for OIDC.
	Audience string `yaml:"audience,omitempty"`

	// Authentication type: SERVICE_PRINCIPAL or IDENTITY.
	AuthType string `yaml:"auth-type,omitempty"`
}

// Action returns the action reference.
func (a AzureLogin) Action() string {
	return "azure/login@v2"
}

// Inputs returns the action inputs as a map.
func (a AzureLogin) Inputs() map[string]any {
	with := make(map[string]any)

	if a.Creds != "" {
		with["creds"] = a.Creds
	}
	if a.ClientID != "" {
		with["client-id"] = a.ClientID
	}
	if a.TenantID != "" {
		with["tenant-id"] = a.TenantID
	}
	if a.SubscriptionID != "" {
		with["subscription-id"] = a.SubscriptionID
	}
	if a.EnableAzPSSession {
		with["enable-AzPSSession"] = a.EnableAzPSSession
	}
	if a.Environment != "" {
		with["environment"] = a.Environment
	}
	if a.AllowNoSubscriptions {
		with["allow-no-subscriptions"] = a.AllowNoSubscriptions
	}
	if a.Audience != "" {
		with["audience"] = a.Audience
	}
	if a.AuthType != "" {
		with["auth-type"] = a.AuthType
	}

	return with
}
