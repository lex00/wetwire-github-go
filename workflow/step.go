package workflow

// Step represents a single step in a job.
type Step struct {
	// ID is a unique identifier for this step.
	ID string `yaml:"id,omitempty"`

	// Name is the display name for this step.
	Name string `yaml:"name,omitempty"`

	// If is a conditional expression to determine if this step runs.
	If any `yaml:"if,omitempty"`

	// Uses specifies an action to run.
	Uses string `yaml:"uses,omitempty"`

	// With provides input parameters to the action.
	With map[string]any `yaml:"with,omitempty"`

	// Run executes a shell command.
	Run string `yaml:"run,omitempty"`

	// Shell specifies the shell to use for Run commands.
	Shell string `yaml:"shell,omitempty"`

	// Env sets environment variables for this step.
	Env map[string]any `yaml:"env,omitempty"`

	// WorkingDirectory sets the working directory for Run commands.
	WorkingDirectory string `yaml:"working-directory,omitempty"`

	// ContinueOnError allows the job to continue if this step fails.
	ContinueOnError bool `yaml:"continue-on-error,omitempty"`

	// TimeoutMinutes sets the maximum time for this step.
	TimeoutMinutes int `yaml:"timeout-minutes,omitempty"`
}

// Output returns an OutputRef for referencing this step's outputs.
// The step must have an ID set to reference outputs.
//
// Example:
//
//	checkoutStep := checkout.Checkout{}.ToStepWithID("checkout")
//	// Later reference: checkoutStep.Output("ref")
func (s Step) Output(name string) OutputRef {
	return OutputRef{StepID: s.ID, Output: name}
}

// OutputRef represents a reference to a step output.
// When serialized to YAML, becomes: ${{ steps.step_id.outputs.name }}
type OutputRef struct {
	StepID string
	Output string
}

// String returns the GitHub Actions expression for this output reference.
func (o OutputRef) String() string {
	return Steps.Get(o.StepID, o.Output).String()
}

// Expression returns the OutputRef as an Expression for use in conditionals.
func (o OutputRef) Expression() Expression {
	return Steps.Get(o.StepID, o.Output)
}

// StepAction is implemented by action wrappers to convert to Step.
type StepAction interface {
	// ToStep converts the action to a workflow Step.
	ToStep() Step
}
