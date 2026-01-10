// Package add_and_commit provides a typed wrapper for EndBug/add-and-commit.
package add_and_commit

// AddAndCommit wraps the EndBug/add-and-commit@v9 action.
// Add and commit files to a Git repository.
type AddAndCommit struct {
	// The files to add, separated by spaces or newlines. Default: "."
	Add string `yaml:"add,omitempty"`

	// The name of the user who will author the commit.
	AuthorName string `yaml:"author_name,omitempty"`

	// The email of the user who will author the commit.
	AuthorEmail string `yaml:"author_email,omitempty"`

	// The commit message.
	Message string `yaml:"message,omitempty"`

	// Whether to push the commit to the remote. Default: "true"
	// Can be "true", "false", or a branch name to push to.
	Push string `yaml:"push,omitempty"`

	// How to fill missing author name/email. Options: "github_actor", "user_info", "github_actions"
	DefaultAuthor string `yaml:"default_author,omitempty"`
}

// Action returns the action reference.
func (a AddAndCommit) Action() string {
	return "EndBug/add-and-commit@v9"
}

// Inputs returns the action inputs as a map.
func (a AddAndCommit) Inputs() map[string]any {
	with := make(map[string]any)

	if a.Add != "" {
		with["add"] = a.Add
	}
	if a.AuthorName != "" {
		with["author_name"] = a.AuthorName
	}
	if a.AuthorEmail != "" {
		with["author_email"] = a.AuthorEmail
	}
	if a.Message != "" {
		with["message"] = a.Message
	}
	if a.Push != "" {
		with["push"] = a.Push
	}
	if a.DefaultAuthor != "" {
		with["default_author"] = a.DefaultAuthor
	}

	return with
}
