// Package setup_helm provides a typed wrapper for azure/setup-helm.
package setup_helm

// SetupHelm wraps the azure/setup-helm@v4 action.
// Install Helm CLI on a GitHub Actions runner.
type SetupHelm struct {
	// The version of Helm to install
	Version string `yaml:"version,omitempty"`

	// GitHub token for downloading Helm releases
	Token string `yaml:"token,omitempty"`

	// Set the download base URL
	DownloadBaseURL string `yaml:"downloadBaseURL,omitempty"`
}

// Action returns the action reference.
func (a SetupHelm) Action() string {
	return "azure/setup-helm@v4"
}

// Inputs returns the action inputs as a map.
func (a SetupHelm) Inputs() map[string]any {
	with := make(map[string]any)

	if a.Version != "" {
		with["version"] = a.Version
	}
	if a.Token != "" {
		with["token"] = a.Token
	}
	if a.DownloadBaseURL != "" {
		with["downloadBaseURL"] = a.DownloadBaseURL
	}

	return with
}
