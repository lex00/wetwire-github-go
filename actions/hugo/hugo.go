// Package hugo provides a typed wrapper for peaceiris/actions-hugo.
package hugo

// Hugo wraps the peaceiris/actions-hugo@v3 action.
// Setup Hugo static site generator.
type Hugo struct {
	// The Hugo version to download (if necessary) and use
	HugoVersion string `yaml:"hugo-version,omitempty"`

	// Set to true to use the extended edition of Hugo
	Extended bool `yaml:"extended,omitempty"`

	// GitHub token for downloading Hugo releases
	GitHubToken string `yaml:"github-token,omitempty"`
}

// Action returns the action reference.
func (a Hugo) Action() string {
	return "peaceiris/actions-hugo@v3"
}

// Inputs returns the action inputs as a map.
func (a Hugo) Inputs() map[string]any {
	with := make(map[string]any)

	if a.HugoVersion != "" {
		with["hugo-version"] = a.HugoVersion
	}
	if a.Extended {
		with["extended"] = a.Extended
	}
	if a.GitHubToken != "" {
		with["github-token"] = a.GitHubToken
	}

	return with
}
