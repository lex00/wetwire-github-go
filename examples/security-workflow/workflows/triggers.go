package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// SecurityPush triggers on push to main branch.
var SecurityPush = workflow.PushTrigger{
	Branches: []string{"main"},
}

// SecurityPullRequest triggers on pull request to main branch.
var SecurityPullRequest = workflow.PullRequestTrigger{
	Branches: []string{"main"},
}

// SecurityTriggers combines push and pull request triggers.
var SecurityTriggers = workflow.Triggers{
	Push:        &SecurityPush,
	PullRequest: &SecurityPullRequest,
}
