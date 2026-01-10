// Package codeql_init provides a typed wrapper for github/codeql-action/init.
package codeql_init

// CodeQLInit wraps the github/codeql-action/init@v3 action.
// Initialize CodeQL for scanning and set up the analysis environment.
type CodeQLInit struct {
	// Languages to analyze (comma-separated: go, javascript, python, etc.).
	Languages string `yaml:"languages,omitempty"`

	// Queries to run (security-extended, security-and-quality, or path to queries).
	Queries string `yaml:"queries,omitempty"`

	// Configuration file for CodeQL.
	ConfigFile string `yaml:"config-file,omitempty"`

	// Path to external CodeQL configuration.
	ExternalRepositoryToken string `yaml:"external-repository-token,omitempty"`

	// Tools URL for CodeQL bundle.
	Tools string `yaml:"tools,omitempty"`

	// Enable debug mode.
	Debug bool `yaml:"debug,omitempty"`

	// RAM limit for CodeQL in MB.
	RAM string `yaml:"ram,omitempty"`

	// Number of threads for CodeQL.
	Threads string `yaml:"threads,omitempty"`

	// Tracing settings for compiled languages.
	BuildMode string `yaml:"build-mode,omitempty"`
}

// Action returns the action reference.
func (a CodeQLInit) Action() string {
	return "github/codeql-action/init@v3"
}

// Inputs returns the action inputs as a map.
func (a CodeQLInit) Inputs() map[string]any {
	with := make(map[string]any)

	if a.Languages != "" {
		with["languages"] = a.Languages
	}
	if a.Queries != "" {
		with["queries"] = a.Queries
	}
	if a.ConfigFile != "" {
		with["config-file"] = a.ConfigFile
	}
	if a.ExternalRepositoryToken != "" {
		with["external-repository-token"] = a.ExternalRepositoryToken
	}
	if a.Tools != "" {
		with["tools"] = a.Tools
	}
	if a.Debug {
		with["debug"] = a.Debug
	}
	if a.RAM != "" {
		with["ram"] = a.RAM
	}
	if a.Threads != "" {
		with["threads"] = a.Threads
	}
	if a.BuildMode != "" {
		with["build-mode"] = a.BuildMode
	}

	return with
}
