// Package pre_commit provides a typed wrapper for pre-commit/action.
package pre_commit

// PreCommit wraps the pre-commit/action@v3.0.1 action.
// Run pre-commit hooks as a GitHub Action.
type PreCommit struct {
	// ExtraArgs are additional arguments to pass to pre-commit run.
	// Example: "--all-files" or "--config .pre-commit-config.yaml"
	ExtraArgs string `yaml:"extra_args,omitempty"`

	// Token is the GitHub token for committing fixes (if auto-fix is enabled).
	Token string `yaml:"token,omitempty"`
}

// Action returns the action reference.
func (a PreCommit) Action() string {
	return "pre-commit/action@v3.0.1"
}

// Inputs returns the action inputs as a map.
func (a PreCommit) Inputs() map[string]any {
	with := make(map[string]any)

	if a.ExtraArgs != "" {
		with["extra_args"] = a.ExtraArgs
	}
	if a.Token != "" {
		with["token"] = a.Token
	}

	return with
}
