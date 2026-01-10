// Package setup_terraform provides a typed wrapper for hashicorp/setup-terraform.
package setup_terraform

// SetupTerraform wraps the hashicorp/setup-terraform@v3 action.
// Setup Terraform CLI in your GitHub Actions workflow.
type SetupTerraform struct {
	// The hostname of a HCP Terraform/Terraform Enterprise instance
	CliConfigCredentialsHostname string `yaml:"cli_config_credentials_hostname,omitempty"`

	// The API token for a HCP Terraform/Terraform Enterprise instance
	CliConfigCredentialsToken string `yaml:"cli_config_credentials_token,omitempty"`

	// The version of Terraform CLI to install (e.g., "1.7.0", "latest", "<1.13.0")
	TerraformVersion string `yaml:"terraform_version,omitempty"`

	// Whether to install a wrapper to expose Terraform outputs
	TerraformWrapper bool `yaml:"terraform_wrapper,omitempty"`
}

// Action returns the action reference.
func (a SetupTerraform) Action() string {
	return "hashicorp/setup-terraform@v3"
}

// Inputs returns the action inputs as a map.
func (a SetupTerraform) Inputs() map[string]any {
	with := make(map[string]any)

	if a.CliConfigCredentialsHostname != "" {
		with["cli_config_credentials_hostname"] = a.CliConfigCredentialsHostname
	}
	if a.CliConfigCredentialsToken != "" {
		with["cli_config_credentials_token"] = a.CliConfigCredentialsToken
	}
	if a.TerraformVersion != "" {
		with["terraform_version"] = a.TerraformVersion
	}
	if a.TerraformWrapper {
		with["terraform_wrapper"] = a.TerraformWrapper
	}

	return with
}
