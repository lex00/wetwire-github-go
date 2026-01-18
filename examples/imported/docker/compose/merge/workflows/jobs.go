package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var BinImage = workflow.Job{
	RunsOn: "ubuntu-22.04",
	Steps:  BinImageSteps,
}

var DesktopEdgeTest = workflow.Job{
	RunsOn: "ubuntu-latest",
	Needs:  []any{"bin-image"},
	Steps:  DesktopEdgeTestSteps,
}

var E2e = workflow.Job{
	Name:           "Build and test",
	RunsOn:         "${{ matrix.os }}",
	TimeoutMinutes: 15,
	Steps:          E2eSteps,
}
