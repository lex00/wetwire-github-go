// Package super_linter provides a typed wrapper for super-linter/super-linter.
package super_linter

// SuperLinter wraps the super-linter/super-linter@v7 action.
// Run Super-Linter to lint code across multiple languages and formats.
type SuperLinter struct {
	// Validate all codebase or only changed files.
	ValidateAllCodebase bool `yaml:"validate_all_codebase,omitempty"`

	// Default branch to compare against.
	DefaultBranch string `yaml:"default_branch,omitempty"`

	// GitHub token for API access.
	GithubToken string `yaml:"github_token,omitempty"`

	// Regex to exclude files from linting.
	FilterRegexExclude string `yaml:"filter_regex_exclude,omitempty"`

	// Regex to include files for linting.
	FilterRegexInclude string `yaml:"filter_regex_include,omitempty"`

	// Log level for output (DEBUG, INFO, NOTICE, WARNING, ERROR).
	LogLevel string `yaml:"log_level,omitempty"`

	// Output format (tap).
	OutputFormat string `yaml:"output_format,omitempty"`

	// Output details (simpler or detailed).
	OutputDetails string `yaml:"output_details,omitempty"`

	// Enable linting of Go files.
	ValidateGo bool `yaml:"validate_go,omitempty"`

	// Enable linting of JavaScript files.
	ValidateJavascript bool `yaml:"validate_javascript,omitempty"`

	// Enable linting of TypeScript files.
	ValidateTypescript bool `yaml:"validate_typescript,omitempty"`

	// Enable linting of Python files.
	ValidatePython bool `yaml:"validate_python,omitempty"`

	// Enable linting of YAML files.
	ValidateYaml bool `yaml:"validate_yaml,omitempty"`

	// Enable linting of JSON files.
	ValidateJson bool `yaml:"validate_json,omitempty"`

	// Enable linting of Markdown files.
	ValidateMarkdown bool `yaml:"validate_markdown,omitempty"`

	// Enable linting of Dockerfile files.
	ValidateDockerfile bool `yaml:"validate_dockerfile,omitempty"`

	// Enable linting of shell scripts.
	ValidateBash bool `yaml:"validate_bash,omitempty"`

	// Path to workspace.
	DefaultWorkspace string `yaml:"default_workspace,omitempty"`

	// Linter rules path.
	LinterRulesPath string `yaml:"linter_rules_path,omitempty"`
}

// Action returns the action reference.
func (a SuperLinter) Action() string {
	return "super-linter/super-linter@v7"
}

// Inputs returns the action inputs as a map.
func (a SuperLinter) Inputs() map[string]any {
	with := make(map[string]any)

	if a.ValidateAllCodebase {
		with["validate_all_codebase"] = a.ValidateAllCodebase
	}
	if a.DefaultBranch != "" {
		with["default_branch"] = a.DefaultBranch
	}
	if a.GithubToken != "" {
		with["github_token"] = a.GithubToken
	}
	if a.FilterRegexExclude != "" {
		with["filter_regex_exclude"] = a.FilterRegexExclude
	}
	if a.FilterRegexInclude != "" {
		with["filter_regex_include"] = a.FilterRegexInclude
	}
	if a.LogLevel != "" {
		with["log_level"] = a.LogLevel
	}
	if a.OutputFormat != "" {
		with["output_format"] = a.OutputFormat
	}
	if a.OutputDetails != "" {
		with["output_details"] = a.OutputDetails
	}
	if a.ValidateGo {
		with["validate_go"] = a.ValidateGo
	}
	if a.ValidateJavascript {
		with["validate_javascript"] = a.ValidateJavascript
	}
	if a.ValidateTypescript {
		with["validate_typescript"] = a.ValidateTypescript
	}
	if a.ValidatePython {
		with["validate_python"] = a.ValidatePython
	}
	if a.ValidateYaml {
		with["validate_yaml"] = a.ValidateYaml
	}
	if a.ValidateJson {
		with["validate_json"] = a.ValidateJson
	}
	if a.ValidateMarkdown {
		with["validate_markdown"] = a.ValidateMarkdown
	}
	if a.ValidateDockerfile {
		with["validate_dockerfile"] = a.ValidateDockerfile
	}
	if a.ValidateBash {
		with["validate_bash"] = a.ValidateBash
	}
	if a.DefaultWorkspace != "" {
		with["default_workspace"] = a.DefaultWorkspace
	}
	if a.LinterRulesPath != "" {
		with["linter_rules_path"] = a.LinterRulesPath
	}

	return with
}
