package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var QuickChecks = workflow.Workflow{
	Name: "Quick Checks",
	On:   QuickChecksTriggers,
	Jobs: map[string]workflow.Job{
		"consistency-checks": ConsistencyChecks,
		"e2e-tests":          E2eTests,
		"race-tests":         RaceTests,
		"unit-tests":         UnitTests,
	},
}
