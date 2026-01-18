package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var Backport = workflow.Job{
	RunsOn: "ubuntu-latest",
	If:     "github.event.workflow_run.head_repository.fork == false && github.repository == 'grafana/grafana' && github.event.workflow_run.conclusion == 'success'",
	Steps:  BackportSteps,
}
