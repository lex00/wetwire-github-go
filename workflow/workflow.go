package workflow

// Workflow represents a GitHub Actions workflow.
type Workflow struct {
	// Name is the display name of the workflow.
	Name string `yaml:"name,omitempty"`

	// On defines the events that trigger this workflow.
	On Triggers `yaml:"on"`

	// Env sets environment variables for all jobs.
	Env map[string]any `yaml:"env,omitempty"`

	// Defaults sets default settings for all jobs.
	Defaults *WorkflowDefaults `yaml:"defaults,omitempty"`

	// Concurrency controls concurrent workflow execution.
	Concurrency *Concurrency `yaml:"concurrency,omitempty"`

	// Permissions sets GITHUB_TOKEN permissions for all jobs.
	Permissions *Permissions `yaml:"permissions,omitempty"`

	// Jobs contains the workflow jobs.
	// This is typically populated by the build process from discovered Job variables.
	Jobs map[string]Job `yaml:"jobs,omitempty"`
}

// WorkflowDefaults sets default settings for all jobs in a workflow.
type WorkflowDefaults struct {
	Run *RunDefaults `yaml:"run,omitempty"`
}

// ResourceType returns "workflow" for interface compliance.
func (w Workflow) ResourceType() string {
	return "workflow"
}
