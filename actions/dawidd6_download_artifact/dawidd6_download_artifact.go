// Package dawidd6_download_artifact provides a typed wrapper for dawidd6/action-download-artifact.
package dawidd6_download_artifact

// DownloadArtifact wraps the dawidd6/action-download-artifact@v6 action.
// Download artifacts from a different workflow run or repository.
type DownloadArtifact struct {
	// GitHub token for authentication
	GitHubToken string `yaml:"github_token,omitempty"`

	// Workflow file name or ID to download artifacts from
	Workflow string `yaml:"workflow,omitempty"`

	// Artifact name to download
	Name string `yaml:"name,omitempty"`

	// Download path for the artifact
	Path string `yaml:"path,omitempty"`

	// Branch to download artifacts from
	Branch string `yaml:"branch,omitempty"`

	// Repository to download artifacts from (owner/repo format)
	Repo string `yaml:"repo,omitempty"`

	// Run ID of the workflow to download artifacts from
	RunID string `yaml:"run_id,omitempty"`

	// Run number of the workflow to download artifacts from
	RunNumber string `yaml:"run_number,omitempty"`

	// Behavior if no artifact is found: error, warn, or ignore
	IfNoArtifactFound string `yaml:"if_no_artifact_found,omitempty"`

	// Allow artifacts from fork pull requests
	AllowForks bool `yaml:"allow_forks,omitempty"`

	// Check artifacts from all workflow runs
	CheckArtifacts bool `yaml:"check_artifacts,omitempty"`

	// Search for artifacts across all workflow runs
	SearchArtifacts bool `yaml:"search_artifacts,omitempty"`
}

// Action returns the action reference.
func (a DownloadArtifact) Action() string {
	return "dawidd6/action-download-artifact@v6"
}

// Inputs returns the action inputs as a map.
func (a DownloadArtifact) Inputs() map[string]any {
	with := make(map[string]any)

	if a.GitHubToken != "" {
		with["github_token"] = a.GitHubToken
	}
	if a.Workflow != "" {
		with["workflow"] = a.Workflow
	}
	if a.Name != "" {
		with["name"] = a.Name
	}
	if a.Path != "" {
		with["path"] = a.Path
	}
	if a.Branch != "" {
		with["branch"] = a.Branch
	}
	if a.Repo != "" {
		with["repo"] = a.Repo
	}
	if a.RunID != "" {
		with["run_id"] = a.RunID
	}
	if a.RunNumber != "" {
		with["run_number"] = a.RunNumber
	}
	if a.IfNoArtifactFound != "" {
		with["if_no_artifact_found"] = a.IfNoArtifactFound
	}
	if a.AllowForks {
		with["allow_forks"] = a.AllowForks
	}
	if a.CheckArtifacts {
		with["check_artifacts"] = a.CheckArtifacts
	}
	if a.SearchArtifacts {
		with["search_artifacts"] = a.SearchArtifacts
	}

	return with
}
