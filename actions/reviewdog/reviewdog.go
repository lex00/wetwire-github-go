// Package reviewdog provides a typed wrapper for reviewdog/action-setup.
package reviewdog

// Reviewdog wraps the reviewdog/action-setup@v1 action.
// Install reviewdog for automated code review tool integration.
type Reviewdog struct {
	// Version of reviewdog to install (e.g., "latest", "v0.17.0").
	ReviewdogVersion string `yaml:"reviewdog_version,omitempty"`
}

// Action returns the action reference.
func (a Reviewdog) Action() string {
	return "reviewdog/action-setup@v1"
}

// Inputs returns the action inputs as a map.
func (a Reviewdog) Inputs() map[string]any {
	with := make(map[string]any)

	if a.ReviewdogVersion != "" {
		with["reviewdog_version"] = a.ReviewdogVersion
	}

	return with
}

// ReviewdogReporter wraps the reviewdog/reviewdog-action action for running reviewdog.
// Run reviewdog to post review comments from linter outputs.
type ReviewdogReporter struct {
	// GitHub token for API access.
	GithubToken string `yaml:"github_token,omitempty"`

	// Workdir relative to the root directory.
	Workdir string `yaml:"workdir,omitempty"`

	// Reporter type: github-pr-check, github-pr-review, github-check.
	Reporter string `yaml:"reporter,omitempty"`

	// Filter mode for reviewdog (added, diff_context, file, nofilter).
	Filter string `yaml:"filter,omitempty"`

	// Exit code for reviewdog when errors are found.
	FailOnError bool `yaml:"fail_on_error,omitempty"`

	// Level for reviewdog (info, warning, error).
	Level string `yaml:"level,omitempty"`

	// Reviewdog flags (e.g., "-diff='git diff FETCH_HEAD'").
	ReviewdogFlags string `yaml:"reviewdog_flags,omitempty"`

	// Tool name for reviewdog.
	Name string `yaml:"name,omitempty"`
}

// Action returns the action reference.
func (a ReviewdogReporter) Action() string {
	return "reviewdog/action-reviewdog@v1"
}

// Inputs returns the action inputs as a map.
func (a ReviewdogReporter) Inputs() map[string]any {
	with := make(map[string]any)

	if a.GithubToken != "" {
		with["github_token"] = a.GithubToken
	}
	if a.Workdir != "" {
		with["workdir"] = a.Workdir
	}
	if a.Reporter != "" {
		with["reporter"] = a.Reporter
	}
	if a.Filter != "" {
		with["filter"] = a.Filter
	}
	if a.FailOnError {
		with["fail_on_error"] = a.FailOnError
	}
	if a.Level != "" {
		with["level"] = a.Level
	}
	if a.ReviewdogFlags != "" {
		with["reviewdog_flags"] = a.ReviewdogFlags
	}
	if a.Name != "" {
		with["name"] = a.Name
	}

	return with
}
