// Package checkout provides a typed wrapper for actions/checkout.
package checkout

// Checkout wraps the actions/checkout@v4 action.
// Checkout a Git repository at a particular version.
type Checkout struct {
	// Repository name with owner (e.g., actions/checkout)
	Repository string `yaml:"repository,omitempty"`

	// The branch, tag or SHA to checkout
	Ref string `yaml:"ref,omitempty"`

	// Personal access token (PAT) used to fetch the repository
	Token string `yaml:"token,omitempty"`

	// SSH key used to fetch the repository
	SSHKey string `yaml:"ssh-key,omitempty"`

	// Known hosts in addition to the user and global host key database
	SSHKnownHosts string `yaml:"ssh-known-hosts,omitempty"`

	// Whether to perform strict host key checking
	SSHStrict bool `yaml:"ssh-strict,omitempty"`

	// Whether to configure the token or SSH key with the local git config
	PersistCredentials bool `yaml:"persist-credentials,omitempty"`

	// Relative path under $GITHUB_WORKSPACE to place the repository
	Path string `yaml:"path,omitempty"`

	// Whether to execute git clean -ffdx && git reset --hard HEAD before fetching
	Clean bool `yaml:"clean,omitempty"`

	// Partially clone against a given filter
	Filter string `yaml:"filter,omitempty"`

	// Do a sparse checkout on given patterns
	SparseCheckout string `yaml:"sparse-checkout,omitempty"`

	// Specifies whether to use cone-mode when doing a sparse checkout
	SparseCheckoutConeMode bool `yaml:"sparse-checkout-cone-mode,omitempty"`

	// Number of commits to fetch. 0 indicates all history for all branches and tags
	FetchDepth int `yaml:"fetch-depth,omitempty"`

	// Whether to fetch tags, even if fetch-depth > 0
	FetchTags bool `yaml:"fetch-tags,omitempty"`

	// Whether to show progress status output when fetching
	ShowProgress bool `yaml:"show-progress,omitempty"`

	// Whether to download Git-LFS files
	LFS bool `yaml:"lfs,omitempty"`

	// Whether to checkout submodules: true to checkout submodules, recursive to recursively checkout
	Submodules string `yaml:"submodules,omitempty"`

	// Add repository path as safe.directory for Git global config
	SetSafeDirectory bool `yaml:"set-safe-directory,omitempty"`

	// The base URL for the GitHub instance to clone from
	GithubServerURL string `yaml:"github-server-url,omitempty"`
}

// Action returns the action reference.
func (a Checkout) Action() string {
	return "actions/checkout@v4"
}

// Inputs returns the action inputs as a map.
func (a Checkout) Inputs() map[string]any {
	with := make(map[string]any)

	if a.Repository != "" {
		with["repository"] = a.Repository
	}
	if a.Ref != "" {
		with["ref"] = a.Ref
	}
	if a.Token != "" {
		with["token"] = a.Token
	}
	if a.SSHKey != "" {
		with["ssh-key"] = a.SSHKey
	}
	if a.SSHKnownHosts != "" {
		with["ssh-known-hosts"] = a.SSHKnownHosts
	}
	if a.SSHStrict {
		with["ssh-strict"] = a.SSHStrict
	}
	if a.PersistCredentials {
		with["persist-credentials"] = a.PersistCredentials
	}
	if a.Path != "" {
		with["path"] = a.Path
	}
	if a.Clean {
		with["clean"] = a.Clean
	}
	if a.Filter != "" {
		with["filter"] = a.Filter
	}
	if a.SparseCheckout != "" {
		with["sparse-checkout"] = a.SparseCheckout
	}
	if a.SparseCheckoutConeMode {
		with["sparse-checkout-cone-mode"] = a.SparseCheckoutConeMode
	}
	if a.FetchDepth != 0 {
		with["fetch-depth"] = a.FetchDepth
	}
	if a.FetchTags {
		with["fetch-tags"] = a.FetchTags
	}
	if a.ShowProgress {
		with["show-progress"] = a.ShowProgress
	}
	if a.LFS {
		with["lfs"] = a.LFS
	}
	if a.Submodules != "" {
		with["submodules"] = a.Submodules
	}
	if a.SetSafeDirectory {
		with["set-safe-directory"] = a.SetSafeDirectory
	}
	if a.GithubServerURL != "" {
		with["github-server-url"] = a.GithubServerURL
	}

	return with
}
