// Package fossa provides a typed wrapper for fossas/fossa-action.
package fossa

// Fossa wraps the fossas/fossa-action@v1 action.
// Run FOSSA license compliance and security analysis.
type Fossa struct {
	// FOSSA API key for authentication
	APIKey string `yaml:"api-key,omitempty"`

	// Override the detected branch name
	Branch string `yaml:"branch,omitempty"`

	// Override the detected revision/commit hash
	Revision string `yaml:"revision,omitempty"`

	// Custom FOSSA CLI container image to use
	Container string `yaml:"container,omitempty"`
}

// Action returns the action reference.
func (a Fossa) Action() string {
	return "fossas/fossa-action@v1"
}

// Inputs returns the action inputs as a map.
func (a Fossa) Inputs() map[string]any {
	with := make(map[string]any)

	if a.APIKey != "" {
		with["api-key"] = a.APIKey
	}
	if a.Branch != "" {
		with["branch"] = a.Branch
	}
	if a.Revision != "" {
		with["revision"] = a.Revision
	}
	if a.Container != "" {
		with["container"] = a.Container
	}

	return with
}
