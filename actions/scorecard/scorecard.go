// Package scorecard provides a typed wrapper for ossf/scorecard-action.
package scorecard

// Scorecard wraps the ossf/scorecard-action@v2.4.0 action.
// Run OpenSSF Scorecard to assess the security posture of open source projects.
type Scorecard struct {
	// File path to store the results.
	ResultsFile string `yaml:"results_file,omitempty"`

	// Format of the results (sarif or json).
	ResultsFormat string `yaml:"results_format,omitempty"`

	// Whether to publish the results to the repository.
	PublishResults bool `yaml:"publish_results,omitempty"`

	// GitHub token for API access.
	RepoToken string `yaml:"repo_token,omitempty"`
}

// Action returns the action reference.
func (a Scorecard) Action() string {
	return "ossf/scorecard-action@v2.4.0"
}

// Inputs returns the action inputs as a map.
func (a Scorecard) Inputs() map[string]any {
	with := make(map[string]any)

	if a.ResultsFile != "" {
		with["results_file"] = a.ResultsFile
	}
	if a.ResultsFormat != "" {
		with["results_format"] = a.ResultsFormat
	}
	if a.PublishResults {
		with["publish_results"] = a.PublishResults
	}
	if a.RepoToken != "" {
		with["repo_token"] = a.RepoToken
	}

	return with
}
