// Package download_artifact provides a typed wrapper for actions/download-artifact.
package download_artifact

// DownloadArtifact wraps the actions/download-artifact@v4 action.
// Download a build artifact previously uploaded in the workflow.
type DownloadArtifact struct {
	// Name of the artifact to download. If unspecified, all artifacts are downloaded
	Name string `yaml:"name,omitempty"`

	// Destination path. Defaults to $GITHUB_WORKSPACE
	Path string `yaml:"path,omitempty"`

	// A glob pattern to filter artifacts by name
	Pattern string `yaml:"pattern,omitempty"`

	// When multiple artifacts are matched, this changes the behavior of the destination directories
	MergeMultiple bool `yaml:"merge-multiple,omitempty"`

	// The GitHub token used to authenticate with the GitHub API
	GithubToken string `yaml:"github-token,omitempty"`

	// The repository to download artifacts from
	Repository string `yaml:"repository,omitempty"`

	// The id of the workflow run to download artifacts from
	RunID string `yaml:"run-id,omitempty"`
}

// Action returns the action reference.
func (a DownloadArtifact) Action() string {
	return "actions/download-artifact@v4"
}

// Inputs returns the action inputs as a map.
func (a DownloadArtifact) Inputs() map[string]any {
	with := make(map[string]any)

	if a.Name != "" {
		with["name"] = a.Name
	}
	if a.Path != "" {
		with["path"] = a.Path
	}
	if a.Pattern != "" {
		with["pattern"] = a.Pattern
	}
	if a.MergeMultiple {
		with["merge-multiple"] = a.MergeMultiple
	}
	if a.GithubToken != "" {
		with["github-token"] = a.GithubToken
	}
	if a.Repository != "" {
		with["repository"] = a.Repository
	}
	if a.RunID != "" {
		with["run-id"] = a.RunID
	}

	return with
}
