package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var ConsistencyChecks = workflow.Job{
	Name:   "Code Consistency Checks",
	RunsOn: "ubuntu-latest",
	Steps:  ConsistencyChecksSteps,
}

var E2eTests = workflow.Job{
	Name:   "End-to-end Tests",
	RunsOn: "ubuntu-latest",
	Steps:  E2eTestsSteps,
}

var RaceTests = workflow.Job{
	Name:   "Race Tests",
	RunsOn: "ubuntu-latest",
	Steps:  RaceTestsSteps,
}

var UnitTests = workflow.Job{
	Name:   "Unit Tests",
	RunsOn: "ubuntu-latest",
	Steps:  UnitTestsSteps,
}
