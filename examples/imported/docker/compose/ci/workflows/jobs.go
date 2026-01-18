package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var Binary = workflow.Job{
	RunsOn: "ubuntu-latest",
	Needs:  []any{"prepare"},
	Steps:  BinarySteps,
}

var Coverage = workflow.Job{
	RunsOn: "ubuntu-latest",
	Needs:  []any{"test", "e2e"},
	Steps:  CoverageSteps,
}

var E2e = workflow.Job{
	Name:   "e2e (${{ matrix.mode }}, ${{ matrix.channel }})",
	RunsOn: "ubuntu-latest",
	Steps:  E2eSteps,
}

var Prepare = workflow.Job{
	RunsOn: "ubuntu-latest",
	Steps:  PrepareSteps,
}

var Release = workflow.Job{
	RunsOn: "ubuntu-latest",
	Needs:  []any{"binary"},
	Steps:  ReleaseSteps,
}

var Test = workflow.Job{
	RunsOn: "ubuntu-latest",
	Steps:  TestSteps,
}

var Validate = workflow.Job{
	RunsOn: "ubuntu-latest",
	Steps:  ValidateSteps,
}
