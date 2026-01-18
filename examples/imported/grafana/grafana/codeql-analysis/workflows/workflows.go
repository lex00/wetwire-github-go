package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var CodeQLChecks = workflow.Workflow{
	Name: "CodeQL checks",
	On:   CodeQLChecksTriggers,
	Jobs: map[string]workflow.Job{
		"analyze":        Analyze,
		"detect-changes": DetectChanges,
	},
}
