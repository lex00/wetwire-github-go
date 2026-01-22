package main

import "github.com/wetwire/wetwire-github/workflow"

var CICDTriggers = workflow.Triggers{
	Push:        CIPush,
	PullRequest: CIPullRequest,
}

var CIPush = workflow.PushTrigger{
	Branches: List("main"),
}

var CIPullRequest = workflow.PullRequestTrigger{
	Branches: List("main"),
}
