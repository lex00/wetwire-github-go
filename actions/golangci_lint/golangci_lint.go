// Package golangci_lint provides a typed wrapper for golangci/golangci-lint-action.
package golangci_lint

// GolangciLint wraps the golangci/golangci-lint-action@v6 action.
// Run golangci-lint for Go code linting.
type GolangciLint struct {
	// Version of golangci-lint to use (e.g., "v1.61", "latest").
	Version string `yaml:"version,omitempty"`

	// Working directory relative to repository root.
	WorkingDirectory string `yaml:"working-directory,omitempty"`

	// Golangci-lint command line arguments.
	Args string `yaml:"args,omitempty"`

	// Only show new issues (for PRs, compared to base branch).
	OnlyNewIssues bool `yaml:"only-new-issues,omitempty"`

	// Skip Go build cache.
	SkipBuildCache bool `yaml:"skip-build-cache,omitempty"`

	// Skip Go package cache.
	SkipPkgCache bool `yaml:"skip-pkg-cache,omitempty"`

	// Enable GitHub Actions problem matchers.
	ProblemMatchers bool `yaml:"problem-matchers,omitempty"`

	// GitHub token for API requests.
	GithubToken string `yaml:"github-token,omitempty"`

	// Install golangci-lint only (don't run).
	InstallMode string `yaml:"install-mode,omitempty"`

	// Force the usage of Go modules.
	GoModules bool `yaml:"go-modules,omitempty"`

	// Skip cache entirely.
	SkipCache bool `yaml:"skip-cache,omitempty"`
}

// Action returns the action reference.
func (a GolangciLint) Action() string {
	return "golangci/golangci-lint-action@v6"
}

// Inputs returns the action inputs as a map.
func (a GolangciLint) Inputs() map[string]any {
	with := make(map[string]any)

	if a.Version != "" {
		with["version"] = a.Version
	}
	if a.WorkingDirectory != "" {
		with["working-directory"] = a.WorkingDirectory
	}
	if a.Args != "" {
		with["args"] = a.Args
	}
	if a.OnlyNewIssues {
		with["only-new-issues"] = a.OnlyNewIssues
	}
	if a.SkipBuildCache {
		with["skip-build-cache"] = a.SkipBuildCache
	}
	if a.SkipPkgCache {
		with["skip-pkg-cache"] = a.SkipPkgCache
	}
	if a.ProblemMatchers {
		with["problem-matchers"] = a.ProblemMatchers
	}
	if a.GithubToken != "" {
		with["github-token"] = a.GithubToken
	}
	if a.InstallMode != "" {
		with["install-mode"] = a.InstallMode
	}
	if a.GoModules {
		with["go-modules"] = a.GoModules
	}
	if a.SkipCache {
		with["skip-cache"] = a.SkipCache
	}

	return with
}
