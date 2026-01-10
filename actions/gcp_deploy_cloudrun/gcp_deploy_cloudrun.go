// Package gcp_deploy_cloudrun provides a typed wrapper for google-github-actions/deploy-cloudrun.
package gcp_deploy_cloudrun

// GCPDeployCloudRun wraps the google-github-actions/deploy-cloudrun@v2 action.
// Deploy a container to Google Cloud Run.
type GCPDeployCloudRun struct {
	// ID or fully-qualified identifier of the Cloud Run service.
	Service string `yaml:"service,omitempty"`

	// ID or fully-qualified identifier of the Cloud Run job.
	Job string `yaml:"job,omitempty"`

	// Fully-qualified container image name.
	Image string `yaml:"image,omitempty"`

	// Path to source code for deployment.
	Source string `yaml:"source,omitempty"`

	// YAML service description.
	Metadata string `yaml:"metadata,omitempty"`

	// Environment variables (KEY=VALUE pairs).
	EnvVars string `yaml:"env_vars,omitempty"`

	// Environment variable update strategy: merge or overwrite.
	EnvVarsUpdateStrategy string `yaml:"env_vars_update_strategy,omitempty"`

	// Secrets (KEY=VALUE pairs).
	Secrets string `yaml:"secrets,omitempty"`

	// Secrets update strategy: merge or overwrite.
	SecretsUpdateStrategy string `yaml:"secrets_update_strategy,omitempty"`

	// Labels for the service.
	Labels string `yaml:"labels,omitempty"`

	// Traffic tag for new revision.
	Tag string `yaml:"tag,omitempty"`

	// Maximum request execution time.
	Timeout string `yaml:"timeout,omitempty"`

	// Additional gcloud run flags.
	Flags string `yaml:"flags,omitempty"`

	// Don't route traffic to new revision.
	NoTraffic bool `yaml:"no_traffic,omitempty"`

	// Google Cloud project ID.
	ProjectID string `yaml:"project_id,omitempty"`

	// Region for Cloud Run deployment.
	Region string `yaml:"region,omitempty"`

	// Suffix for revision name.
	Suffix string `yaml:"suffix,omitempty"`

	// Skip default labels from GitHub Actions.
	SkipDefaultLabels bool `yaml:"skip_default_labels,omitempty"`
}

// Action returns the action reference.
func (a GCPDeployCloudRun) Action() string {
	return "google-github-actions/deploy-cloudrun@v2"
}

// Inputs returns the action inputs as a map.
func (a GCPDeployCloudRun) Inputs() map[string]any {
	with := make(map[string]any)

	if a.Service != "" {
		with["service"] = a.Service
	}
	if a.Job != "" {
		with["job"] = a.Job
	}
	if a.Image != "" {
		with["image"] = a.Image
	}
	if a.Source != "" {
		with["source"] = a.Source
	}
	if a.Metadata != "" {
		with["metadata"] = a.Metadata
	}
	if a.EnvVars != "" {
		with["env_vars"] = a.EnvVars
	}
	if a.EnvVarsUpdateStrategy != "" {
		with["env_vars_update_strategy"] = a.EnvVarsUpdateStrategy
	}
	if a.Secrets != "" {
		with["secrets"] = a.Secrets
	}
	if a.SecretsUpdateStrategy != "" {
		with["secrets_update_strategy"] = a.SecretsUpdateStrategy
	}
	if a.Labels != "" {
		with["labels"] = a.Labels
	}
	if a.Tag != "" {
		with["tag"] = a.Tag
	}
	if a.Timeout != "" {
		with["timeout"] = a.Timeout
	}
	if a.Flags != "" {
		with["flags"] = a.Flags
	}
	if a.NoTraffic {
		with["no_traffic"] = a.NoTraffic
	}
	if a.ProjectID != "" {
		with["project_id"] = a.ProjectID
	}
	if a.Region != "" {
		with["region"] = a.Region
	}
	if a.Suffix != "" {
		with["suffix"] = a.Suffix
	}
	if a.SkipDefaultLabels {
		with["skip_default_labels"] = a.SkipDefaultLabels
	}

	return with
}
