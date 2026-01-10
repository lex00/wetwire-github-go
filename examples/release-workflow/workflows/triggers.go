package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// ReleasePush triggers on version tags (v*).
var ReleasePush = workflow.PushTrigger{
	Tags: []string{"v*"},
}

// ReleaseTriggers activates on version tag pushes.
var ReleaseTriggers = workflow.Triggers{
	Push: &ReleasePush,
}
