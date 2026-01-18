package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var Merge = workflow.Workflow{
	Name: "merge",
	On:   MergeTriggers,
	Jobs: map[string]workflow.Job{
		"bin-image":         BinImage,
		"desktop-edge-test": DesktopEdgeTest,
		"e2e":               E2e,
	},
}
