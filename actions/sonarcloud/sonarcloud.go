// Package sonarcloud provides a typed wrapper for SonarSource/sonarcloud-github-action.
package sonarcloud

// SonarCloud wraps the SonarSource/sonarcloud-github-action@v3 action.
// Detect bugs and vulnerabilities in your code with SonarCloud analysis.
type SonarCloud struct {
	// ProjectBaseDir is the project base directory for the analysis.
	// Default is the repository root.
	ProjectBaseDir string `yaml:"projectBaseDir,omitempty"`

	// Args are additional arguments passed to the sonar-scanner.
	// Example: "-Dsonar.verbose=true"
	Args string `yaml:"args,omitempty"`
}

// Action returns the action reference.
func (a SonarCloud) Action() string {
	return "SonarSource/sonarcloud-github-action@v3"
}

// Inputs returns the action inputs as a map.
func (a SonarCloud) Inputs() map[string]any {
	with := make(map[string]any)

	if a.ProjectBaseDir != "" {
		with["projectBaseDir"] = a.ProjectBaseDir
	}
	if a.Args != "" {
		with["args"] = a.Args
	}

	return with
}
