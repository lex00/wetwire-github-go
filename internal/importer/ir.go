// Package importer provides YAML parsing into intermediate representation.
package importer

import "gopkg.in/yaml.v3"

// IRWorkflow represents a parsed GitHub Actions workflow.
type IRWorkflow struct {
	Name        string              `yaml:"name,omitempty"`
	On          IRTriggers          `yaml:"on,omitempty"`
	Env         map[string]any      `yaml:"env,omitempty"`
	Defaults    *IRDefaults         `yaml:"defaults,omitempty"`
	Concurrency *IRConcurrency      `yaml:"concurrency,omitempty"`
	Jobs        map[string]*IRJob   `yaml:"jobs,omitempty"`
	Permissions map[string]string   `yaml:"permissions,omitempty"`
}

// IRTriggers represents workflow triggers.
// Can be a string, slice, or map depending on YAML structure.
type IRTriggers struct {
	Push              *IRPushTrigger         `yaml:"push,omitempty"`
	PullRequest       *IRPullRequestTrigger  `yaml:"pull_request,omitempty"`
	PullRequestTarget *IRPullRequestTrigger  `yaml:"pull_request_target,omitempty"`
	WorkflowDispatch  *IRWorkflowDispatch    `yaml:"workflow_dispatch,omitempty"`
	WorkflowCall      *IRWorkflowCall        `yaml:"workflow_call,omitempty"`
	Schedule          []IRSchedule           `yaml:"schedule,omitempty"`
	RepositoryDispatch *IRRepositoryDispatch `yaml:"repository_dispatch,omitempty"`
	Release           *IRReleaseTrigger      `yaml:"release,omitempty"`
	Issues            *IRIssuesTrigger       `yaml:"issues,omitempty"`
	Raw               any                    `yaml:"-"` // Original raw value
}

// UnmarshalYAML implements custom unmarshalling for flexible trigger syntax.
func (t *IRTriggers) UnmarshalYAML(value *yaml.Node) error {
	// Handle string: "on: push"
	if value.Kind == yaml.ScalarNode {
		t.setSimpleTrigger(value.Value)
		t.Raw = value.Value
		return nil
	}

	// Handle sequence: "on: [push, pull_request]"
	if value.Kind == yaml.SequenceNode {
		var triggers []string
		if err := value.Decode(&triggers); err != nil {
			return err
		}
		for _, trigger := range triggers {
			t.setSimpleTrigger(trigger)
		}
		t.Raw = triggers
		return nil
	}

	// Handle mapping: "on: {push: {...}}"
	if value.Kind == yaml.MappingNode {
		type rawTriggers IRTriggers
		var rt rawTriggers
		if err := value.Decode(&rt); err != nil {
			return err
		}
		*t = IRTriggers(rt)
		return nil
	}

	return nil
}

// setSimpleTrigger sets a trigger from a simple string name.
func (t *IRTriggers) setSimpleTrigger(name string) {
	switch name {
	case "push":
		t.Push = &IRPushTrigger{}
	case "pull_request":
		t.PullRequest = &IRPullRequestTrigger{}
	case "pull_request_target":
		t.PullRequestTarget = &IRPullRequestTrigger{}
	case "workflow_dispatch":
		t.WorkflowDispatch = &IRWorkflowDispatch{}
	case "workflow_call":
		t.WorkflowCall = &IRWorkflowCall{}
	case "repository_dispatch":
		t.RepositoryDispatch = &IRRepositoryDispatch{}
	case "release":
		t.Release = &IRReleaseTrigger{}
	case "issues":
		t.Issues = &IRIssuesTrigger{}
	}
}

// IRPushTrigger represents a push trigger.
type IRPushTrigger struct {
	Branches       []string `yaml:"branches,omitempty"`
	BranchesIgnore []string `yaml:"branches-ignore,omitempty"`
	Tags           []string `yaml:"tags,omitempty"`
	TagsIgnore     []string `yaml:"tags-ignore,omitempty"`
	Paths          []string `yaml:"paths,omitempty"`
	PathsIgnore    []string `yaml:"paths-ignore,omitempty"`
}

// IRPullRequestTrigger represents a pull_request trigger.
type IRPullRequestTrigger struct {
	Branches       []string `yaml:"branches,omitempty"`
	BranchesIgnore []string `yaml:"branches-ignore,omitempty"`
	Paths          []string `yaml:"paths,omitempty"`
	PathsIgnore    []string `yaml:"paths-ignore,omitempty"`
	Types          []string `yaml:"types,omitempty"`
}

// IRWorkflowDispatch represents manual workflow dispatch inputs.
type IRWorkflowDispatch struct {
	Inputs map[string]IRInput `yaml:"inputs,omitempty"`
}

// IRWorkflowCall represents reusable workflow call inputs and outputs.
type IRWorkflowCall struct {
	Inputs  map[string]IRInput  `yaml:"inputs,omitempty"`
	Outputs map[string]IROutput `yaml:"outputs,omitempty"`
	Secrets map[string]IRSecret `yaml:"secrets,omitempty"`
}

// IRInput represents a workflow input.
type IRInput struct {
	Description string `yaml:"description,omitempty"`
	Required    bool   `yaml:"required,omitempty"`
	Default     any    `yaml:"default,omitempty"`
	Type        string `yaml:"type,omitempty"`
	Options     []any  `yaml:"options,omitempty"`
}

// IROutput represents a workflow output.
type IROutput struct {
	Description string `yaml:"description,omitempty"`
	Value       string `yaml:"value,omitempty"`
}

// IRSecret represents a workflow secret.
type IRSecret struct {
	Description string `yaml:"description,omitempty"`
	Required    bool   `yaml:"required,omitempty"`
}

// IRSchedule represents a cron schedule.
type IRSchedule struct {
	Cron string `yaml:"cron,omitempty"`
}

// IRRepositoryDispatch represents repository dispatch events.
type IRRepositoryDispatch struct {
	Types []string `yaml:"types,omitempty"`
}

// IRReleaseTrigger represents release events.
type IRReleaseTrigger struct {
	Types []string `yaml:"types,omitempty"`
}

// IRIssuesTrigger represents issues events.
type IRIssuesTrigger struct {
	Types []string `yaml:"types,omitempty"`
}

// IRDefaults represents workflow defaults.
type IRDefaults struct {
	Run *IRRunDefaults `yaml:"run,omitempty"`
}

// IRRunDefaults represents default run settings.
type IRRunDefaults struct {
	Shell            string `yaml:"shell,omitempty"`
	WorkingDirectory string `yaml:"working-directory,omitempty"`
}

// IRConcurrency represents concurrency settings.
type IRConcurrency struct {
	Group            string `yaml:"group,omitempty"`
	CancelInProgress bool   `yaml:"cancel-in-progress,omitempty"`
}

// IRJob represents a job in a workflow.
type IRJob struct {
	Name            string              `yaml:"name,omitempty"`
	RunsOn          any                 `yaml:"runs-on,omitempty"`
	Needs           any                 `yaml:"needs,omitempty"` // string or []string
	If              string              `yaml:"if,omitempty"`
	Env             map[string]any      `yaml:"env,omitempty"`
	Defaults        *IRDefaults         `yaml:"defaults,omitempty"`
	Steps           []IRStep            `yaml:"steps,omitempty"`
	TimeoutMinutes  int                 `yaml:"timeout-minutes,omitempty"`
	Strategy        *IRStrategy         `yaml:"strategy,omitempty"`
	ContinueOnError any                 `yaml:"continue-on-error,omitempty"`
	Container       *IRContainer        `yaml:"container,omitempty"`
	Services        map[string]*IRService `yaml:"services,omitempty"`
	Outputs         map[string]string   `yaml:"outputs,omitempty"`
	Permissions     map[string]string   `yaml:"permissions,omitempty"`
	Concurrency     *IRConcurrency      `yaml:"concurrency,omitempty"`
	Uses            string              `yaml:"uses,omitempty"` // For reusable workflows
	With            map[string]any      `yaml:"with,omitempty"`
	Secrets         any                 `yaml:"secrets,omitempty"` // string "inherit" or map
}

// IRStrategy represents a job strategy.
type IRStrategy struct {
	Matrix      *IRMatrix `yaml:"matrix,omitempty"`
	FailFast    *bool     `yaml:"fail-fast,omitempty"`
	MaxParallel int       `yaml:"max-parallel,omitempty"`
}

// IRMatrix represents a matrix configuration.
type IRMatrix struct {
	Include []map[string]any `yaml:"include,omitempty"`
	Exclude []map[string]any `yaml:"exclude,omitempty"`
	Values  map[string][]any `yaml:"-"` // Extracted from other fields
	Raw     map[string]any   `yaml:"-"` // Original raw data
}

// IRContainer represents a container configuration.
type IRContainer struct {
	Image       string            `yaml:"image,omitempty"`
	Credentials *IRCredentials    `yaml:"credentials,omitempty"`
	Env         map[string]string `yaml:"env,omitempty"`
	Ports       []any             `yaml:"ports,omitempty"`
	Volumes     []string          `yaml:"volumes,omitempty"`
	Options     string            `yaml:"options,omitempty"`
}

// IRCredentials represents container credentials.
type IRCredentials struct {
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

// IRService represents a service container.
type IRService struct {
	Image       string            `yaml:"image,omitempty"`
	Credentials *IRCredentials    `yaml:"credentials,omitempty"`
	Env         map[string]string `yaml:"env,omitempty"`
	Ports       []any             `yaml:"ports,omitempty"`
	Volumes     []string          `yaml:"volumes,omitempty"`
	Options     string            `yaml:"options,omitempty"`
}

// IRStep represents a step in a job.
type IRStep struct {
	ID              string         `yaml:"id,omitempty"`
	Name            string         `yaml:"name,omitempty"`
	Uses            string         `yaml:"uses,omitempty"`
	Run             string         `yaml:"run,omitempty"`
	Shell           string         `yaml:"shell,omitempty"`
	With            map[string]any `yaml:"with,omitempty"`
	Env             map[string]any `yaml:"env,omitempty"`
	If              string         `yaml:"if,omitempty"`
	ContinueOnError any            `yaml:"continue-on-error,omitempty"`
	TimeoutMinutes  int            `yaml:"timeout-minutes,omitempty"`
	WorkingDirectory string        `yaml:"working-directory,omitempty"`
}

// GetNeeds returns the job dependencies as a slice of strings.
func (j *IRJob) GetNeeds() []string {
	if j.Needs == nil {
		return nil
	}
	switch v := j.Needs.(type) {
	case string:
		return []string{v}
	case []any:
		needs := make([]string, 0, len(v))
		for _, n := range v {
			if s, ok := n.(string); ok {
				needs = append(needs, s)
			}
		}
		return needs
	case []string:
		return v
	}
	return nil
}

// GetRunsOn returns the runs-on value as a string or the first element if it's a slice.
func (j *IRJob) GetRunsOn() string {
	if j.RunsOn == nil {
		return ""
	}
	switch v := j.RunsOn.(type) {
	case string:
		return v
	case []any:
		if len(v) > 0 {
			if s, ok := v[0].(string); ok {
				return s
			}
		}
	case []string:
		if len(v) > 0 {
			return v[0]
		}
	}
	return ""
}
