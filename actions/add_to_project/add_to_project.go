// Package add_to_project provides a typed wrapper for actions/add-to-project.
package add_to_project

// AddToProject wraps the actions/add-to-project@v1 action.
// Automate adding issues and pull requests to GitHub projects.
type AddToProject struct {
	// URL of the project to add issues to (required)
	ProjectURL string `yaml:"project-url,omitempty"`

	// GitHub personal access token with write access to the project (required)
	GithubToken string `yaml:"github-token,omitempty"`

	// Comma-separated list of labels to use as a filter for issues to be added
	Labeled string `yaml:"labeled,omitempty"`

	// Behavior of the labels filter: AND, OR, or NOT (default: OR)
	LabelOperator string `yaml:"label-operator,omitempty"`
}

// Action returns the action reference.
func (a AddToProject) Action() string {
	return "actions/add-to-project@v1"
}

// Inputs returns the action inputs as a map.
func (a AddToProject) Inputs() map[string]any {
	with := make(map[string]any)

	if a.ProjectURL != "" {
		with["project-url"] = a.ProjectURL
	}
	if a.GithubToken != "" {
		with["github-token"] = a.GithubToken
	}
	if a.Labeled != "" {
		with["labeled"] = a.Labeled
	}
	if a.LabelOperator != "" {
		with["label-operator"] = a.LabelOperator
	}

	return with
}
