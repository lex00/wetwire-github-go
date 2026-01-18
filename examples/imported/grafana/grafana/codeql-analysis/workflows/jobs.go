package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var Analyze = workflow.Job{
	Name:   "Analyze",
	RunsOn: "ubuntu-x64-large-io",
	Needs:  []any{"detect-changes"},
	If:     "github.repository == 'grafana/grafana'",
	Steps:  AnalyzeSteps,
}

var DetectChanges = workflow.Job{
	Name:   "Detect whether code changed",
	RunsOn: "ubuntu-latest",
	Steps:  DetectChangesSteps,
}
