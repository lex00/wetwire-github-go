// Package pulumi provides a typed wrapper for pulumi/actions.
package pulumi

// Pulumi wraps the pulumi/actions@v6 action.
// Deploy infrastructure using Pulumi in GitHub Actions.
type Pulumi struct {
	// The Pulumi command to run (up, preview, destroy, refresh)
	Command string `yaml:"command,omitempty"`

	// The name of the Pulumi stack to operate on
	StackName string `yaml:"stack-name,omitempty"`

	// The working directory to run Pulumi commands in
	WorkDir string `yaml:"work-dir,omitempty"`

	// The URL of the Pulumi Cloud backend
	CloudURL string `yaml:"cloud-url,omitempty"`

	// Configuration values as a JSON map
	ConfigMap string `yaml:"config-map,omitempty"`

	// The secrets provider to use for encrypting secrets
	SecretsProvider string `yaml:"secrets-provider,omitempty"`

	// Colorize output (auto, always, never, raw)
	Color string `yaml:"color,omitempty"`

	// Show the diff for the update
	Diff bool `yaml:"diff,omitempty"`

	// Comment on the PR with the results of the Pulumi operation
	CommentOnPR bool `yaml:"comment-on-pr,omitempty"`

	// Edit existing PR comment instead of creating a new one
	EditPRComment bool `yaml:"edit-pr-comment,omitempty"`
}

// Action returns the action reference.
func (a Pulumi) Action() string {
	return "pulumi/actions@v6"
}

// Inputs returns the action inputs as a map.
func (a Pulumi) Inputs() map[string]any {
	with := make(map[string]any)

	if a.Command != "" {
		with["command"] = a.Command
	}
	if a.StackName != "" {
		with["stack-name"] = a.StackName
	}
	if a.WorkDir != "" {
		with["work-dir"] = a.WorkDir
	}
	if a.CloudURL != "" {
		with["cloud-url"] = a.CloudURL
	}
	if a.ConfigMap != "" {
		with["config-map"] = a.ConfigMap
	}
	if a.SecretsProvider != "" {
		with["secrets-provider"] = a.SecretsProvider
	}
	if a.Color != "" {
		with["color"] = a.Color
	}
	if a.Diff {
		with["diff"] = a.Diff
	}
	if a.CommentOnPR {
		with["comment-on-pr"] = a.CommentOnPR
	}
	if a.EditPRComment {
		with["edit-pr-comment"] = a.EditPRComment
	}

	return with
}
