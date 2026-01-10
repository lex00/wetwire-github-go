package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// PipelinePush triggers on push to main branch or version tags.
var PipelinePush = workflow.PushTrigger{
	Branches: []string{"main"},
	Tags:     []string{"v*"},
}

// PipelineTriggers activates on push to main or version tags.
var PipelineTriggers = workflow.Triggers{
	Push: &PipelinePush,
}
