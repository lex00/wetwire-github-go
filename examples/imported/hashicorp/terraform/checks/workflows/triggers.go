package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var QuickChecksPush = workflow.PushTrigger{
	Branches: []string{"*"},
	Tags:     []string{"v[0-9]+.[0-9]+.[0-9]+*"},
}

var QuickChecksPullRequest = workflow.PullRequestTrigger{
	Types: []string{"opened", "ready_for_review", "reopened", "synchronize"},
}

var QuickChecksTriggers = workflow.Triggers{
	Push:        &QuickChecksPush,
	PullRequest: &QuickChecksPullRequest,
}
