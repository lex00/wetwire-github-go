package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// DockerPush triggers on push to main branch.
var DockerPush = workflow.PushTrigger{
	Branches: []string{"main"},
}

// DockerPullRequest triggers on pull request to main branch.
var DockerPullRequest = workflow.PullRequestTrigger{
	Branches: []string{"main"},
}

// DockerTriggers combines push and pull request triggers.
var DockerTriggers = workflow.Triggers{
	Push:        &DockerPush,
	PullRequest: &DockerPullRequest,
}
