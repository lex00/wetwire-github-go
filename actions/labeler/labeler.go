// Package labeler provides a typed wrapper for actions/labeler.
package labeler

// Labeler wraps the actions/labeler@v5 action.
// Automatically label pull requests based on file patterns.
type Labeler struct {
	// Token for API access
	RepoToken string `yaml:"repo-token,omitempty"`

	// Path to configuration file
	ConfigurationPath string `yaml:"configuration-path,omitempty"`

	// Remove labels not matching rules
	SyncLabels bool `yaml:"sync-labels,omitempty"`

	// Enable globbing for hidden files
	Dot bool `yaml:"dot,omitempty"`

	// PR number to label (optional)
	PRNumber int `yaml:"pr-number,omitempty"`
}

// Action returns the action reference.
func (a Labeler) Action() string {
	return "actions/labeler@v5"
}

// Inputs returns the action inputs as a map.
func (a Labeler) Inputs() map[string]any {
	with := make(map[string]any)

	if a.RepoToken != "" {
		with["repo-token"] = a.RepoToken
	}
	if a.ConfigurationPath != "" {
		with["configuration-path"] = a.ConfigurationPath
	}
	if a.SyncLabels {
		with["sync-labels"] = a.SyncLabels
	}
	if a.Dot {
		with["dot"] = a.Dot
	}
	if a.PRNumber != 0 {
		with["pr-number"] = a.PRNumber
	}

	return with
}
