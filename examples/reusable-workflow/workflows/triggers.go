package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// ReusableInputs defines the typed inputs for the reusable workflow.
var ReusableInputs = map[string]workflow.WorkflowInput{
	"go_version": {
		Description: "Go version to use for building",
		Required:    true,
		Type:        "string",
		Default:     "1.24",
	},
	"run_tests": {
		Description: "Whether to run tests",
		Required:    false,
		Type:        "boolean",
		Default:     true,
	},
	"build_target": {
		Description: "Build target (linux, darwin, windows)",
		Required:    false,
		Type:        "choice",
		Options:     []string{"linux", "darwin", "windows"},
		Default:     "linux",
	},
}

// ReusableOutputs defines the typed outputs from the reusable workflow.
var ReusableOutputs = map[string]workflow.WorkflowOutput{
	"artifact_name": {
		Description: "Name of the built artifact",
		Value:       "${{ jobs.build.outputs.artifact }}",
	},
	"build_version": {
		Description: "Version of the build",
		Value:       "${{ jobs.build.outputs.version }}",
	},
}

// ReusableSecrets defines the secrets that can be passed to the reusable workflow.
var ReusableSecrets = map[string]workflow.WorkflowSecret{
	"DEPLOY_TOKEN": {
		Description: "Token for deployment operations",
		Required:    false,
	},
	"REGISTRY_PASSWORD": {
		Description: "Password for container registry",
		Required:    false,
	},
}

// ReusableWorkflowCall defines the workflow_call trigger with inputs, outputs, and secrets.
var ReusableWorkflowCall = workflow.WorkflowCallTrigger{
	Inputs:  ReusableInputs,
	Outputs: ReusableOutputs,
	Secrets: ReusableSecrets,
}

// ReusableTriggers combines the workflow_call trigger for the reusable workflow.
var ReusableTriggers = workflow.Triggers{
	WorkflowCall: &ReusableWorkflowCall,
}

// CallerPush triggers on push to main branch.
var CallerPush = workflow.PushTrigger{
	Branches: []string{"main"},
}

// CallerPullRequest triggers on pull request to main branch.
var CallerPullRequest = workflow.PullRequestTrigger{
	Branches: []string{"main"},
}

// CallerTriggers combines push and pull request triggers for the caller workflow.
var CallerTriggers = workflow.Triggers{
	Push:        &CallerPush,
	PullRequest: &CallerPullRequest,
}
