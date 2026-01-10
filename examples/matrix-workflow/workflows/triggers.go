package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// MatrixPush triggers on push to main branch.
var MatrixPush = workflow.PushTrigger{
	Branches: []string{"main"},
}

// MatrixPullRequest triggers on pull request to main branch.
var MatrixPullRequest = workflow.PullRequestTrigger{
	Branches: []string{"main"},
}

// MatrixTriggers combines push and pull request triggers.
var MatrixTriggers = workflow.Triggers{
	Push:        &MatrixPush,
	PullRequest: &MatrixPullRequest,
}
