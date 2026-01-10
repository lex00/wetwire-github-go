// Package github_script provides a typed wrapper for actions/github-script.
package github_script

// GithubScript wraps the actions/github-script@v7 action.
// Run JavaScript in your workflows using the GitHub API and workflow contexts.
type GithubScript struct {
	// The script to run. Required.
	Script string `yaml:"script,omitempty"`

	// The GitHub token to use for authentication. Defaults to github.token.
	GithubToken string `yaml:"github-token,omitempty"`

	// Whether to enable debug logging.
	Debug bool `yaml:"debug,omitempty"`

	// An optional user-agent string.
	UserAgent string `yaml:"user-agent,omitempty"`

	// A comma-separated list of API previews to accept.
	Previews string `yaml:"previews,omitempty"`

	// How the result will be encoded. Can be "string" or "json".
	ResultEncoding string `yaml:"result-encoding,omitempty"`

	// The number of times to retry a request.
	Retries int `yaml:"retries,omitempty"`

	// A comma-separated list of status codes that will NOT be retried.
	RetryExemptStatusCodes string `yaml:"retry-exempt-status-codes,omitempty"`
}

// Action returns the action reference.
func (a GithubScript) Action() string {
	return "actions/github-script@v7"
}

// Inputs returns the action inputs as a map.
func (a GithubScript) Inputs() map[string]any {
	with := make(map[string]any)

	if a.Script != "" {
		with["script"] = a.Script
	}
	if a.GithubToken != "" {
		with["github-token"] = a.GithubToken
	}
	if a.Debug {
		with["debug"] = a.Debug
	}
	if a.UserAgent != "" {
		with["user-agent"] = a.UserAgent
	}
	if a.Previews != "" {
		with["previews"] = a.Previews
	}
	if a.ResultEncoding != "" {
		with["result-encoding"] = a.ResultEncoding
	}
	if a.Retries != 0 {
		with["retries"] = a.Retries
	}
	if a.RetryExemptStatusCodes != "" {
		with["retry-exempt-status-codes"] = a.RetryExemptStatusCodes
	}

	return with
}
