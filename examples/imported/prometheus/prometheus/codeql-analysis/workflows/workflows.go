package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var CodeQL = workflow.Workflow{
	Name: "CodeQL",
	On:   CodeQLTriggers,
	Jobs: map[string]workflow.Job{
		"analyze": Analyze,
	},
}
