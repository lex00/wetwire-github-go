// Package deploy_pages provides a typed wrapper for actions/deploy-pages.
package deploy_pages

// DeployPages wraps the actions/deploy-pages@v4 action.
// Deploys an artifact to GitHub Pages.
type DeployPages struct {
	// GitHub token for authentication
	Token string `yaml:"token,omitempty"`

	// Deployment timeout in milliseconds
	Timeout int `yaml:"timeout,omitempty"`

	// Acceptable error count during deployment
	ErrorCount int `yaml:"error_count,omitempty"`

	// Progress reporting interval in milliseconds
	ReportingInterval int `yaml:"reporting_interval,omitempty"`

	// Name of the artifact to deploy
	ArtifactName string `yaml:"artifact_name,omitempty"`
}

// Action returns the action reference.
func (a DeployPages) Action() string {
	return "actions/deploy-pages@v4"
}

// Inputs returns the action inputs as a map.
func (a DeployPages) Inputs() map[string]any {
	m := make(map[string]any)

	if a.Token != "" {
		m["token"] = a.Token
	}
	if a.Timeout != 0 {
		m["timeout"] = a.Timeout
	}
	if a.ErrorCount != 0 {
		m["error_count"] = a.ErrorCount
	}
	if a.ReportingInterval != 0 {
		m["reporting_interval"] = a.ReportingInterval
	}
	if a.ArtifactName != "" {
		m["artifact_name"] = a.ArtifactName
	}

	return m
}
