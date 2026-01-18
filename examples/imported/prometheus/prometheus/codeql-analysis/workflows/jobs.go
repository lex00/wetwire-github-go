package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var Analyze = workflow.Job{
	Name:   "Analyze",
	RunsOn: "ubuntu-latest",
	Steps:  AnalyzeSteps,
}
