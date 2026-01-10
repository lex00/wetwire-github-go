// Package gcp_setup_gcloud provides a typed wrapper for google-github-actions/setup-gcloud.
package gcp_setup_gcloud

// GCPSetupGcloud wraps the google-github-actions/setup-gcloud@v2 action.
// Set up and configure the Google Cloud SDK (gcloud).
type GCPSetupGcloud struct {
	// Version of Cloud SDK to install (e.g., "290.0.1" or "latest").
	Version string `yaml:"version,omitempty"`

	// Google Cloud project ID to configure as default.
	ProjectID string `yaml:"project_id,omitempty"`

	// Additional gcloud components to install (comma-separated).
	InstallComponents string `yaml:"install_components,omitempty"`

	// Skip installation and use system gcloud.
	SkipInstall bool `yaml:"skip_install,omitempty"`

	// Cache downloaded artifacts for future runs.
	Cache bool `yaml:"cache,omitempty"`
}

// Action returns the action reference.
func (a GCPSetupGcloud) Action() string {
	return "google-github-actions/setup-gcloud@v2"
}

// Inputs returns the action inputs as a map.
func (a GCPSetupGcloud) Inputs() map[string]any {
	with := make(map[string]any)

	if a.Version != "" {
		with["version"] = a.Version
	}
	if a.ProjectID != "" {
		with["project_id"] = a.ProjectID
	}
	if a.InstallComponents != "" {
		with["install_components"] = a.InstallComponents
	}
	if a.SkipInstall {
		with["skip_install"] = a.SkipInstall
	}
	if a.Cache {
		with["cache"] = a.Cache
	}

	return with
}
