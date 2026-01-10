// Package workflows defines GitHub Actions workflow declarations.
package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// ScheduledDispatch is a workflow that supports both scheduled and manual triggers.
// Scheduled runs perform maintenance tasks, while manual dispatches allow deployment.
var ScheduledDispatch = workflow.Workflow{
	Name: "Scheduled and Dispatch",
	On:   WorkflowTriggers,
	Jobs: map[string]workflow.Job{
		"maintenance": Maintenance,
		"deploy":      Deploy,
	},
}
