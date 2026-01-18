package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var Ci = workflow.Workflow{
	Name: "ci",
	On:   CiTriggers,
	Jobs: map[string]workflow.Job{
		"binary":   Binary,
		"coverage": Coverage,
		"e2e":      E2e,
		"prepare":  Prepare,
		"release":  Release,
		"test":     Test,
		"validate": Validate,
	},
}
