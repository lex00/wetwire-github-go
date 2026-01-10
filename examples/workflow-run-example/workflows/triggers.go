package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// CIWorkflowRun triggers when the "CI" workflow completes.
// This is useful for deploying after tests pass or sending notifications.
var CIWorkflowRun = workflow.WorkflowRunTrigger{
	Workflows: []string{"CI"},
	Types:     []string{"completed"},
	Branches:  []string{"main"},
}

// DeployTriggers configures the workflow to run after CI completes on main.
var DeployTriggers = workflow.Triggers{
	WorkflowRun: &CIWorkflowRun,
}
