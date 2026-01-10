package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// DeployDispatch triggers on repository_dispatch events with specific event types.
// Event types allow filtering which dispatch events trigger this workflow.
var DeployDispatch = workflow.RepositoryDispatchTrigger{
	Types: []string{"deploy", "deploy-staging", "deploy-production"},
}

// DispatchTriggers configures the workflow to respond to repository dispatch events.
var DispatchTriggers = workflow.Triggers{
	RepositoryDispatch: &DeployDispatch,
}
