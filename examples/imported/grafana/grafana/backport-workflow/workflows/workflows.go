package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var BackportWorkflow = workflow.Workflow{
	Name: "Backport (workflow)",
	On:   BackportWorkflowTriggers,
	Jobs: map[string]workflow.Job{
		"backport": Backport,
	},
}
