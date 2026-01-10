// Package create_pull_request provides a typed wrapper for peter-evans/create-pull-request.
package create_pull_request

// CreatePullRequest wraps the peter-evans/create-pull-request@v6 action.
// Create a pull request for changes to your repository in the actions workspace.
type CreatePullRequest struct {
	// GITHUB_TOKEN or a PAT with repo scope.
	// Default: ${{ github.token }}
	Token string `yaml:"token,omitempty"`

	// Relative path under $GITHUB_WORKSPACE to the repository.
	// Default: .
	Path string `yaml:"path,omitempty"`

	// A comma or newline-separated list of file paths to commit.
	// Paths should follow git's pathspec syntax.
	AddPaths string `yaml:"add-paths,omitempty"`

	// The message to use when committing changes.
	CommitMessage string `yaml:"commit-message,omitempty"`

	// The committer name and email address in the format Display Name <email@address.com>.
	// Default: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>
	Committer string `yaml:"committer,omitempty"`

	// The author name and email address in the format Display Name <email@address.com>.
	// Default: the committer value
	Author string `yaml:"author,omitempty"`

	// Add Signed-off-by line by the committer at the end of the commit log message.
	Signoff bool `yaml:"signoff,omitempty"`

	// The pull request branch name.
	Branch string `yaml:"branch,omitempty"`

	// The branch suffix type.
	// Valid values: random, timestamp, short-commit-hash
	BranchSuffix string `yaml:"branch-suffix,omitempty"`

	// Delete the branch when closing pull requests, and when undeleted after merging.
	DeleteBranch bool `yaml:"delete-branch,omitempty"`

	// The title of the pull request.
	Title string `yaml:"title,omitempty"`

	// The body of the pull request.
	Body string `yaml:"body,omitempty"`

	// The path to a file containing the pull request body.
	BodyPath string `yaml:"body-path,omitempty"`

	// A comma or newline-separated list of labels.
	Labels string `yaml:"labels,omitempty"`

	// A comma or newline-separated list of assignees (GitHub usernames).
	Assignees string `yaml:"assignees,omitempty"`

	// A comma or newline-separated list of reviewers (GitHub usernames).
	Reviewers string `yaml:"reviewers,omitempty"`

	// A comma or newline-separated list of team reviewers (GitHub teams).
	TeamReviewers string `yaml:"team-reviewers,omitempty"`

	// The number of the milestone to associate the pull request with.
	Milestone int `yaml:"milestone,omitempty"`

	// Create a draft pull request.
	Draft bool `yaml:"draft,omitempty"`
}

// Action returns the action reference.
func (a CreatePullRequest) Action() string {
	return "peter-evans/create-pull-request@v6"
}

// Inputs returns the action inputs as a map.
func (a CreatePullRequest) Inputs() map[string]any {
	with := make(map[string]any)

	if a.Token != "" {
		with["token"] = a.Token
	}
	if a.Path != "" {
		with["path"] = a.Path
	}
	if a.AddPaths != "" {
		with["add-paths"] = a.AddPaths
	}
	if a.CommitMessage != "" {
		with["commit-message"] = a.CommitMessage
	}
	if a.Committer != "" {
		with["committer"] = a.Committer
	}
	if a.Author != "" {
		with["author"] = a.Author
	}
	if a.Signoff {
		with["signoff"] = a.Signoff
	}
	if a.Branch != "" {
		with["branch"] = a.Branch
	}
	if a.BranchSuffix != "" {
		with["branch-suffix"] = a.BranchSuffix
	}
	if a.DeleteBranch {
		with["delete-branch"] = a.DeleteBranch
	}
	if a.Title != "" {
		with["title"] = a.Title
	}
	if a.Body != "" {
		with["body"] = a.Body
	}
	if a.BodyPath != "" {
		with["body-path"] = a.BodyPath
	}
	if a.Labels != "" {
		with["labels"] = a.Labels
	}
	if a.Assignees != "" {
		with["assignees"] = a.Assignees
	}
	if a.Reviewers != "" {
		with["reviewers"] = a.Reviewers
	}
	if a.TeamReviewers != "" {
		with["team-reviewers"] = a.TeamReviewers
	}
	if a.Milestone > 0 {
		with["milestone"] = a.Milestone
	}
	if a.Draft {
		with["draft"] = a.Draft
	}

	return with
}
