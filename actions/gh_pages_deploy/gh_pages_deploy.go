// Package gh_pages_deploy provides a typed wrapper for JamesIves/github-pages-deploy-action.
package gh_pages_deploy

// GitHubPagesDeploy wraps the JamesIves/github-pages-deploy-action@v4 action.
// Deploy to GitHub Pages from your GitHub Actions workflow.
type GitHubPagesDeploy struct {
	// Private SSH key to be used with repository deployment key
	SSHKey string `yaml:"ssh-key,omitempty"`

	// Personal access token for deployment. Defaults to repository-scoped token
	Token string `yaml:"token,omitempty"`

	// Target branch for deployment (e.g., gh-pages or docs)
	Branch string `yaml:"branch,omitempty"`

	// Source folder to deploy (required)
	Folder string `yaml:"folder,omitempty"`

	// Optional destination directory on the deployment branch
	TargetFolder string `yaml:"target-folder,omitempty"`

	// Customize the commit message for the deployment
	CommitMessage string `yaml:"commit-message,omitempty"`

	// Delete hashed files from the target folder on the deployment branch with each deploy
	Clean bool `yaml:"clean,omitempty"`

	// Preserve specific files/folders during cleanup
	CleanExclude string `yaml:"clean-exclude,omitempty"`

	// Use --dry-run flag on git push without actually pushing
	DryRun bool `yaml:"dry-run,omitempty"`

	// Force-push to overwrite existing deployments
	Force bool `yaml:"force,omitempty"`

	// Custom name for GitHub config used during commit pushes
	GitConfigName string `yaml:"git-config-name,omitempty"`

	// Custom email for GitHub config used during commit pushes
	GitConfigEmail string `yaml:"git-config-email,omitempty"`

	// Deploy to a different repository (format: Owner/repo-name)
	RepositoryName string `yaml:"repository-name,omitempty"`

	// Add a version tag to the commit
	Tag string `yaml:"tag,omitempty"`

	// Maintain a single commit on deployment branch instead of full history
	SingleCommit bool `yaml:"single-commit,omitempty"`

	// Suppress action output and git messages
	Silent bool `yaml:"silent,omitempty"`

	// Number of rebase attempts before suspending the job
	AttemptLimit int `yaml:"attempt-limit,omitempty"`
}

// Action returns the action reference.
func (a GitHubPagesDeploy) Action() string {
	return "JamesIves/github-pages-deploy-action@v4"
}

// Inputs returns the action inputs as a map.
func (a GitHubPagesDeploy) Inputs() map[string]any {
	with := make(map[string]any)

	if a.SSHKey != "" {
		with["ssh-key"] = a.SSHKey
	}
	if a.Token != "" {
		with["token"] = a.Token
	}
	if a.Branch != "" {
		with["branch"] = a.Branch
	}
	if a.Folder != "" {
		with["folder"] = a.Folder
	}
	if a.TargetFolder != "" {
		with["target-folder"] = a.TargetFolder
	}
	if a.CommitMessage != "" {
		with["commit-message"] = a.CommitMessage
	}
	if a.Clean {
		with["clean"] = a.Clean
	}
	if a.CleanExclude != "" {
		with["clean-exclude"] = a.CleanExclude
	}
	if a.DryRun {
		with["dry-run"] = a.DryRun
	}
	if a.Force {
		with["force"] = a.Force
	}
	if a.GitConfigName != "" {
		with["git-config-name"] = a.GitConfigName
	}
	if a.GitConfigEmail != "" {
		with["git-config-email"] = a.GitConfigEmail
	}
	if a.RepositoryName != "" {
		with["repository-name"] = a.RepositoryName
	}
	if a.Tag != "" {
		with["tag"] = a.Tag
	}
	if a.SingleCommit {
		with["single-commit"] = a.SingleCommit
	}
	if a.Silent {
		with["silent"] = a.Silent
	}
	if a.AttemptLimit != 0 {
		with["attempt-limit"] = a.AttemptLimit
	}

	return with
}
