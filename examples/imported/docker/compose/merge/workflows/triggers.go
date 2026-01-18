package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var MergePush = workflow.PushTrigger{
	Branches: []string{"main"},
	Tags:     []string{"v*"},
}

var MergeTriggers = workflow.Triggers{
	Push: &MergePush,
}
