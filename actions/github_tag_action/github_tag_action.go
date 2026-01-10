// Package github_tag_action provides a typed wrapper for anothrNick/github-tag-action.
package github_tag_action

// GitHubTagAction wraps the anothrNick/github-tag-action@v1 action.
// Automatically bump and tag with SemVer based on merged PR labels.
type GitHubTagAction struct {
	// GitHub token for authentication (required)
	GitHubToken string `yaml:"github_token,omitempty"`

	// Default version bump type: major, minor, or patch
	DefaultBump string `yaml:"default_bump,omitempty"`

	// Prefix to prepend to the tag (e.g., "v")
	TagPrefix string `yaml:"tag_prefix,omitempty"`

	// If true, perform a dry run without creating the tag
	DryRun bool `yaml:"dry_run,omitempty"`

	// Custom tag to use instead of auto-generated
	CustomTag string `yaml:"custom_tag,omitempty"`

	// Initial version if no tags exist
	InitialVersion string `yaml:"initial_version,omitempty"`

	// Comma-separated list of branches for releases
	ReleasesBranches string `yaml:"release_branches,omitempty"`

	// Comma-separated list of branches for prereleases
	PrereleaseBranches string `yaml:"prerelease_branches,omitempty"`
}

// Action returns the action reference.
func (a GitHubTagAction) Action() string {
	return "anothrNick/github-tag-action@v1"
}

// Inputs returns the action inputs as a map.
func (a GitHubTagAction) Inputs() map[string]any {
	with := make(map[string]any)

	if a.GitHubToken != "" {
		with["github_token"] = a.GitHubToken
	}
	if a.DefaultBump != "" {
		with["default_bump"] = a.DefaultBump
	}
	if a.TagPrefix != "" {
		with["tag_prefix"] = a.TagPrefix
	}
	if a.DryRun {
		with["dry_run"] = a.DryRun
	}
	if a.CustomTag != "" {
		with["custom_tag"] = a.CustomTag
	}
	if a.InitialVersion != "" {
		with["initial_version"] = a.InitialVersion
	}
	if a.ReleasesBranches != "" {
		with["release_branches"] = a.ReleasesBranches
	}
	if a.PrereleaseBranches != "" {
		with["prerelease_branches"] = a.PrereleaseBranches
	}

	return with
}
