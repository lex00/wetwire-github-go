package dependabot

// Update defines update configuration for a package ecosystem.
type Update struct {
	// PackageEcosystem is the package manager type.
	// Supported values: bundler, cargo, composer, docker, github-actions,
	// gitsubmodule, gomod, gradle, maven, mix, npm, nuget, pip, pub, swift, terraform.
	PackageEcosystem string `yaml:"package-ecosystem"`

	// Directory is the location of the manifest files.
	// Defaults to "/" (repository root).
	Directory string `yaml:"directory,omitempty"`

	// Directories specifies multiple manifest locations.
	Directories []string `yaml:"directories,omitempty"`

	// Schedule defines when Dependabot checks for updates.
	Schedule Schedule `yaml:"schedule"`

	// Allow defines which dependencies to allow updates for.
	Allow []Allow `yaml:"allow,omitempty"`

	// Ignore defines which dependencies or versions to ignore.
	Ignore []Ignore `yaml:"ignore,omitempty"`

	// Labels are the labels to add to pull requests.
	// Defaults to ["dependencies"].
	Labels []string `yaml:"labels,omitempty"`

	// Assignees are the GitHub users to assign to pull requests.
	Assignees []string `yaml:"assignees,omitempty"`

	// Reviewers are the GitHub users or teams to request review from.
	Reviewers []string `yaml:"reviewers,omitempty"`

	// Milestone is the numeric milestone identifier.
	Milestone int `yaml:"milestone,omitempty"`

	// OpenPullRequestsLimit is the maximum number of open PRs.
	// Defaults to 5.
	OpenPullRequestsLimit int `yaml:"open-pull-requests-limit,omitempty"`

	// RebaseStrategy controls rebasing behavior.
	// Values: "auto" (default), "disabled".
	RebaseStrategy string `yaml:"rebase-strategy,omitempty"`

	// VersioningStrategy controls how manifest versions are updated.
	// Values: "auto", "increase", "increase-if-necessary", "lockfile-only", "widen".
	VersioningStrategy string `yaml:"versioning-strategy,omitempty"`

	// Vendor enables dependency vendoring.
	Vendor bool `yaml:"vendor,omitempty"`

	// TargetBranch is the alternative branch for changes.
	TargetBranch string `yaml:"target-branch,omitempty"`

	// Registries selects which registries to use.
	// Use "*" for all registries, or a list of registry names.
	Registries any `yaml:"registries,omitempty"`

	// Groups defines dependency grouping rules.
	Groups map[string]Group `yaml:"groups,omitempty"`

	// CommitMessage customizes commit message format.
	CommitMessage *CommitMessage `yaml:"commit-message,omitempty"`

	// PullRequestBranchName customizes branch naming.
	PullRequestBranchName *PullRequestBranchName `yaml:"pull-request-branch-name,omitempty"`

	// InsecureExternalCodeExecution controls external code execution.
	// Values: "allow", "deny".
	InsecureExternalCodeExecution string `yaml:"insecure-external-code-execution,omitempty"`
}

// Allow defines dependencies to allow updates for.
type Allow struct {
	// DependencyName matches a specific dependency.
	DependencyName string `yaml:"dependency-name,omitempty"`

	// DependencyType matches dependencies by type.
	// Values: "direct", "indirect", "all", "production", "development".
	DependencyType string `yaml:"dependency-type,omitempty"`
}

// Ignore defines dependencies or versions to ignore.
type Ignore struct {
	// DependencyName matches a specific dependency.
	DependencyName string `yaml:"dependency-name,omitempty"`

	// Versions specifies version patterns to ignore.
	Versions []string `yaml:"versions,omitempty"`

	// UpdateTypes specifies update types to ignore.
	// Values: "version-update:semver-major", "version-update:semver-minor", "version-update:semver-patch".
	UpdateTypes []string `yaml:"update-types,omitempty"`
}

// CommitMessage customizes commit message format.
type CommitMessage struct {
	// Prefix is added before the commit message (max 50 chars).
	Prefix string `yaml:"prefix,omitempty"`

	// PrefixDevelopment is used for development dependencies (max 50 chars).
	PrefixDevelopment string `yaml:"prefix-development,omitempty"`

	// Include adds additional information.
	// Value: "scope".
	Include string `yaml:"include,omitempty"`
}

// PullRequestBranchName customizes PR branch naming.
type PullRequestBranchName struct {
	// Separator between parts of the branch name.
	// Values: "-", "_", "/" (default "/").
	Separator string `yaml:"separator"`
}
