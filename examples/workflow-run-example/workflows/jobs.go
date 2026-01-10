package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// Deploy runs deployment when the triggering CI workflow succeeded.
// Uses github.event.workflow_run.conclusion to check the result.
var Deploy = workflow.Job{
	Name:   "Deploy",
	RunsOn: "ubuntu-latest",
	If:     "${{ github.event.workflow_run.conclusion == 'success' }}",
	Steps:  DeploySteps,
}

// Notify sends notifications for any workflow completion (success or failure).
// This job always runs regardless of the triggering workflow's conclusion.
var Notify = workflow.Job{
	Name:   "Notify",
	RunsOn: "ubuntu-latest",
	Steps:  NotifySteps,
}
