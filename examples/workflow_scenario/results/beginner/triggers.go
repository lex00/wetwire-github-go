package main

import "github.com/wetwire/wetwire-github-go/workflow"

// CI Triggers - Run on push to main and pull requests
var CITriggers = workflow.Triggers{
	Push:        CIPush,
	PullRequest: CIPullRequest,
}

var CIPush = workflow.PushTrigger{
	Branches: List("main"),
}

var CIPullRequest = workflow.PullRequestTrigger{
	Branches: List("main"),
}
