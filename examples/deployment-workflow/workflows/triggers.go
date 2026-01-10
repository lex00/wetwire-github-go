package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// DeployPush triggers on push to main branch.
var DeployPush = workflow.PushTrigger{
	Branches: []string{"main"},
}

// DeployDispatchInputs defines inputs for manual deployment triggering.
var DeployDispatchInputs = map[string]workflow.WorkflowInput{
	"environment": {
		Description: "Target environment for deployment",
		Required:    true,
		Type:        "choice",
		Options:     []string{"staging", "production"},
		Default:     "staging",
	},
	"skip_staging": {
		Description: "Skip staging and deploy directly to production (requires manual approval)",
		Required:    false,
		Type:        "boolean",
		Default:     false,
	},
}

// DeployDispatch allows manual workflow triggering with environment selection.
var DeployDispatch = workflow.WorkflowDispatchTrigger{
	Inputs: DeployDispatchInputs,
}

// DeployTriggers combines push and manual dispatch triggers.
var DeployTriggers = workflow.Triggers{
	Push:             &DeployPush,
	WorkflowDispatch: &DeployDispatch,
}
