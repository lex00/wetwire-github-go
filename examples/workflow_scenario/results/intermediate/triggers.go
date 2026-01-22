package cicd

import (
	"github.com/wetwire/wetwire-github-go/pkg/workflow"
)

var CIPush = workflow.PushTrigger{
	Branches: List("main"),
}

var CIPullRequest = workflow.PullRequestTrigger{
	Branches: List("main"),
}

var CITriggers = workflow.Triggers{
	Push:        CIPush,
	PullRequest: CIPullRequest,
}
