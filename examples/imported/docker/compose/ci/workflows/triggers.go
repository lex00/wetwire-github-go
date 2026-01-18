package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var CiPush = workflow.PushTrigger{
	Branches: []string{"main"},
	Tags:     []string{"v*"},
}

var CiTriggers = workflow.Triggers{
	Push:             &CiPush,
	WorkflowDispatch: &workflow.WorkflowDispatchTrigger{},
}
