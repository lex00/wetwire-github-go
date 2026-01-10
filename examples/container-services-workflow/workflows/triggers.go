package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// CIPush triggers on push to main branch.
var CIPush = workflow.PushTrigger{
	Branches: []string{"main"},
}

// CIPullRequest triggers on pull request to main branch.
var CIPullRequest = workflow.PullRequestTrigger{
	Branches: []string{"main"},
}

// CITriggers combines push and pull request triggers.
var CITriggers = workflow.Triggers{
	Push:        &CIPush,
	PullRequest: &CIPullRequest,
}
