// Package first_interaction provides a typed wrapper for actions/first-interaction.
package first_interaction

// FirstInteraction wraps the actions/first-interaction@v1 action.
// Greet first-time contributors when they open an issue or pull request.
type FirstInteraction struct {
	// Token with permissions to post issue and PR comments
	RepoToken string `yaml:"repo-token,omitempty"`

	// Comment to post on an individual's first issue
	IssueMessage string `yaml:"issue-message,omitempty"`

	// Comment to post on an individual's first pull request
	PRMessage string `yaml:"pr-message,omitempty"`
}

// Action returns the action reference.
func (a FirstInteraction) Action() string {
	return "actions/first-interaction@v1"
}

// Inputs returns the action inputs as a map.
func (a FirstInteraction) Inputs() map[string]any {
	with := make(map[string]any)

	if a.RepoToken != "" {
		with["repo-token"] = a.RepoToken
	}
	if a.IssueMessage != "" {
		with["issue-message"] = a.IssueMessage
	}
	if a.PRMessage != "" {
		with["pr-message"] = a.PRMessage
	}

	return with
}
