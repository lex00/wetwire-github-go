// Package gh_pages_peaceiris provides a typed wrapper for peaceiris/actions-gh-pages.
package gh_pages_peaceiris

// GHPagesPeaceiris wraps the peaceiris/actions-gh-pages@v4 action.
// Deploy static files to GitHub Pages.
type GHPagesPeaceiris struct {
	// SSH private key from repository secret value for pushing
	DeployKey string `yaml:"deploy_key,omitempty"`

	// Generated GITHUB_TOKEN for pushing to remote branch
	GithubToken string `yaml:"github_token,omitempty"`

	// Personal access token for pushing to remote branch
	PersonalToken string `yaml:"personal_token,omitempty"`

	// Target branch for deployment (default: gh-pages)
	PublishBranch string `yaml:"publish_branch,omitempty"`

	// Input directory for deployment (default: public)
	PublishDir string `yaml:"publish_dir,omitempty"`

	// Destination subdirectory for deployment
	DestinationDir string `yaml:"destination_dir,omitempty"`

	// External repository in owner/repo format
	ExternalRepository string `yaml:"external_repository,omitempty"`

	// Whether empty commits should be made to publication branch
	AllowEmptyCommit bool `yaml:"allow_empty_commit,omitempty"`

	// Whether existing files should be retained before deploying
	KeepFiles bool `yaml:"keep_files,omitempty"`

	// Keep only the latest commit on GitHub Pages branch
	ForceOrphan bool `yaml:"force_orphan,omitempty"`

	// Git user.name configuration
	UserName string `yaml:"user_name,omitempty"`

	// Git user.email configuration
	UserEmail string `yaml:"user_email,omitempty"`

	// Custom commit message with triggered commit hash
	CommitMessage string `yaml:"commit_message,omitempty"`

	// Custom full commit message without commit hash
	FullCommitMessage string `yaml:"full_commit_message,omitempty"`

	// Tag name for release
	TagName string `yaml:"tag_name,omitempty"`

	// Tag message for release
	TagMessage string `yaml:"tag_message,omitempty"`

	// Enable GitHub Pages built-in Jekyll
	EnableJekyll bool `yaml:"enable_jekyll,omitempty"`

	// Alias for enable_jekyll to disable .nojekyll file
	DisableNojekyll bool `yaml:"disable_nojekyll,omitempty"`

	// Custom domain configuration
	CNAME string `yaml:"cname,omitempty"`

	// Files or directories to exclude from publish directory
	ExcludeAssets string `yaml:"exclude_assets,omitempty"`
}

// Action returns the action reference.
func (a GHPagesPeaceiris) Action() string {
	return "peaceiris/actions-gh-pages@v4"
}

// Inputs returns the action inputs as a map.
func (a GHPagesPeaceiris) Inputs() map[string]any {
	with := make(map[string]any)

	if a.DeployKey != "" {
		with["deploy_key"] = a.DeployKey
	}
	if a.GithubToken != "" {
		with["github_token"] = a.GithubToken
	}
	if a.PersonalToken != "" {
		with["personal_token"] = a.PersonalToken
	}
	if a.PublishBranch != "" {
		with["publish_branch"] = a.PublishBranch
	}
	if a.PublishDir != "" {
		with["publish_dir"] = a.PublishDir
	}
	if a.DestinationDir != "" {
		with["destination_dir"] = a.DestinationDir
	}
	if a.ExternalRepository != "" {
		with["external_repository"] = a.ExternalRepository
	}
	if a.AllowEmptyCommit {
		with["allow_empty_commit"] = a.AllowEmptyCommit
	}
	if a.KeepFiles {
		with["keep_files"] = a.KeepFiles
	}
	if a.ForceOrphan {
		with["force_orphan"] = a.ForceOrphan
	}
	if a.UserName != "" {
		with["user_name"] = a.UserName
	}
	if a.UserEmail != "" {
		with["user_email"] = a.UserEmail
	}
	if a.CommitMessage != "" {
		with["commit_message"] = a.CommitMessage
	}
	if a.FullCommitMessage != "" {
		with["full_commit_message"] = a.FullCommitMessage
	}
	if a.TagName != "" {
		with["tag_name"] = a.TagName
	}
	if a.TagMessage != "" {
		with["tag_message"] = a.TagMessage
	}
	if a.EnableJekyll {
		with["enable_jekyll"] = a.EnableJekyll
	}
	if a.DisableNojekyll {
		with["disable_nojekyll"] = a.DisableNojekyll
	}
	if a.CNAME != "" {
		with["cname"] = a.CNAME
	}
	if a.ExcludeAssets != "" {
		with["exclude_assets"] = a.ExcludeAssets
	}

	return with
}
