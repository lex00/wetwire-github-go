package workflow

// Job represents a workflow job.
type Job struct {
	// Name is the display name for this job.
	Name string `yaml:"name,omitempty"`

	// RunsOn specifies the runner environment.
	// Can be a string ("ubuntu-latest") or expression (Matrix.Get("os")).
	RunsOn any `yaml:"runs-on"`

	// Needs lists jobs that must complete before this job runs.
	// Use []any{OtherJob1, OtherJob2} to reference other job variables.
	Needs []any `yaml:"needs,omitempty"`

	// If is a conditional expression to determine if this job runs.
	If any `yaml:"if,omitempty"`

	// Permissions sets GITHUB_TOKEN permissions for this job.
	Permissions *Permissions `yaml:"permissions,omitempty"`

	// Environment specifies the deployment environment.
	Environment *Environment `yaml:"environment,omitempty"`

	// Concurrency controls concurrent job execution.
	Concurrency *Concurrency `yaml:"concurrency,omitempty"`

	// Outputs defines job outputs available to dependent jobs.
	Outputs map[string]any `yaml:"outputs,omitempty"`

	// Env sets environment variables for all steps.
	Env map[string]any `yaml:"env,omitempty"`

	// Defaults sets default settings for all steps.
	Defaults *JobDefaults `yaml:"defaults,omitempty"`

	// Strategy configures matrix builds and failure handling.
	Strategy *Strategy `yaml:"strategy,omitempty"`

	// Container runs the job in a container.
	Container *Container `yaml:"container,omitempty"`

	// Services defines service containers for this job.
	Services map[string]Service `yaml:"services,omitempty"`

	// Steps are the tasks that run in this job.
	Steps []Step `yaml:"steps"`

	// TimeoutMinutes sets the maximum time for this job.
	TimeoutMinutes int `yaml:"timeout-minutes,omitempty"`

	// ContinueOnError allows the workflow to continue if this job fails.
	ContinueOnError bool `yaml:"continue-on-error,omitempty"`
}

// Permissions configures GITHUB_TOKEN permissions.
type Permissions struct {
	Actions            string `yaml:"actions,omitempty"`
	Checks             string `yaml:"checks,omitempty"`
	Contents           string `yaml:"contents,omitempty"`
	Deployments        string `yaml:"deployments,omitempty"`
	Discussions        string `yaml:"discussions,omitempty"`
	IDToken            string `yaml:"id-token,omitempty"`
	Issues             string `yaml:"issues,omitempty"`
	Packages           string `yaml:"packages,omitempty"`
	Pages              string `yaml:"pages,omitempty"`
	PullRequests       string `yaml:"pull-requests,omitempty"`
	RepositoryProjects string `yaml:"repository-projects,omitempty"`
	SecurityEvents     string `yaml:"security-events,omitempty"`
	Statuses           string `yaml:"statuses,omitempty"`
}

// PermissionLevel constants for Permissions fields.
const (
	PermissionRead  = "read"
	PermissionWrite = "write"
	PermissionNone  = "none"
)

// Environment configures a deployment environment.
type Environment struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url,omitempty"`
}

// Concurrency controls concurrent workflow/job execution.
type Concurrency struct {
	Group            string `yaml:"group"`
	CancelInProgress bool   `yaml:"cancel-in-progress,omitempty"`
}

// JobDefaults sets default settings for job steps.
type JobDefaults struct {
	Run *RunDefaults `yaml:"run,omitempty"`
}

// RunDefaults sets default settings for run steps.
type RunDefaults struct {
	Shell            string `yaml:"shell,omitempty"`
	WorkingDirectory string `yaml:"working-directory,omitempty"`
}

// Container configures a container for running job steps.
type Container struct {
	Image       string            `yaml:"image"`
	Credentials *Credentials      `yaml:"credentials,omitempty"`
	Env         map[string]any    `yaml:"env,omitempty"`
	Ports       []any             `yaml:"ports,omitempty"`
	Volumes     []string          `yaml:"volumes,omitempty"`
	Options     string            `yaml:"options,omitempty"`
}

// Credentials for container registry authentication.
type Credentials struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// Service configures a service container.
type Service struct {
	Image       string            `yaml:"image"`
	Credentials *Credentials      `yaml:"credentials,omitempty"`
	Env         map[string]any    `yaml:"env,omitempty"`
	Ports       []any             `yaml:"ports,omitempty"`
	Volumes     []string          `yaml:"volumes,omitempty"`
	Options     string            `yaml:"options,omitempty"`
}
