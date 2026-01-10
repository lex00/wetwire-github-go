package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// Maintenance runs scheduled cleanup and reporting tasks.
// This job runs on schedule events.
var Maintenance = workflow.Job{
	Name:   "Scheduled Maintenance",
	RunsOn: "ubuntu-latest",
	If:     "${{ github.event_name == 'schedule' }}",
	Steps:  MaintenanceSteps,
}

// Deploy runs manual deployment with input parameters.
// This job runs on workflow_dispatch events.
var Deploy = workflow.Job{
	Name:   "Deploy",
	RunsOn: "ubuntu-latest",
	If:     "${{ github.event_name == 'workflow_dispatch' }}",
	Steps:  DeploySteps,
}
