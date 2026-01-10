// Package codecov provides a typed wrapper for codecov/codecov-action.
package codecov

// Codecov wraps the codecov/codecov-action@v5 action.
// Upload code coverage reports to Codecov.
type Codecov struct {
	// Repository upload token. Not required for public repos using GitHub Actions.
	Token string `yaml:"token,omitempty"`

	// Comma-separated list of coverage report files to upload.
	Files string `yaml:"files,omitempty"`

	// Directory to search for coverage reports.
	Directory string `yaml:"directory,omitempty"`

	// Comma-separated list of flags to associate with the upload.
	Flags string `yaml:"flags,omitempty"`

	// Custom name for the upload.
	Name string `yaml:"name,omitempty"`

	// Whether to fail the CI if an error is encountered during upload.
	FailCIIfError bool `yaml:"fail_ci_if_error,omitempty"`

	// Enable verbose logging.
	Verbose bool `yaml:"verbose,omitempty"`

	// Working directory for the action.
	WorkingDirectory string `yaml:"working-directory,omitempty"`

	// Environment variables to include in the upload.
	EnvVars string `yaml:"env_vars,omitempty"`

	// Override the detected OS.
	OS string `yaml:"os,omitempty"`

	// Override the repository slug (owner/repo).
	Slug string `yaml:"slug,omitempty"`

	// Version of the Codecov CLI to use.
	Version string `yaml:"version,omitempty"`

	// Don't upload files to Codecov.
	DryRun bool `yaml:"dry_run,omitempty"`

	// Use OIDC instead of token for authentication.
	UseOIDC bool `yaml:"use_oidc,omitempty"`

	// Path to codecov.yml configuration file.
	CodecovYMLPath string `yaml:"codecov_yml_path,omitempty"`

	// Plugins to run. Use "noop" to turn off automatic fixes.
	Plugin string `yaml:"plugin,omitempty"`
}

// Action returns the action reference.
func (a Codecov) Action() string {
	return "codecov/codecov-action@v5"
}

// Inputs returns the action inputs as a map.
func (a Codecov) Inputs() map[string]any {
	with := make(map[string]any)

	if a.Token != "" {
		with["token"] = a.Token
	}
	if a.Files != "" {
		with["files"] = a.Files
	}
	if a.Directory != "" {
		with["directory"] = a.Directory
	}
	if a.Flags != "" {
		with["flags"] = a.Flags
	}
	if a.Name != "" {
		with["name"] = a.Name
	}
	if a.FailCIIfError {
		with["fail_ci_if_error"] = a.FailCIIfError
	}
	if a.Verbose {
		with["verbose"] = a.Verbose
	}
	if a.WorkingDirectory != "" {
		with["working-directory"] = a.WorkingDirectory
	}
	if a.EnvVars != "" {
		with["env_vars"] = a.EnvVars
	}
	if a.OS != "" {
		with["os"] = a.OS
	}
	if a.Slug != "" {
		with["slug"] = a.Slug
	}
	if a.Version != "" {
		with["version"] = a.Version
	}
	if a.DryRun {
		with["dry_run"] = a.DryRun
	}
	if a.UseOIDC {
		with["use_oidc"] = a.UseOIDC
	}
	if a.CodecovYMLPath != "" {
		with["codecov_yml_path"] = a.CodecovYMLPath
	}
	if a.Plugin != "" {
		with["plugin"] = a.Plugin
	}

	return with
}
