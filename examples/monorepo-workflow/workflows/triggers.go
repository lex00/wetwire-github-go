package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// MonorepoPush triggers on push to main branch with path filters.
// Only triggers when files in services/** or shared/** are modified.
var MonorepoPush = workflow.PushTrigger{
	Branches: []string{"main"},
	Paths: []string{
		"services/api/**",
		"services/web/**",
		"shared/**",
	},
}

// MonorepoPullRequest triggers on pull request to main branch with path filters.
// Only triggers when files in services/** or shared/** are modified.
var MonorepoPullRequest = workflow.PullRequestTrigger{
	Branches: []string{"main"},
	Paths: []string{
		"services/api/**",
		"services/web/**",
		"shared/**",
	},
}

// MonorepoTriggers combines push and pull request triggers with path filtering.
var MonorepoTriggers = workflow.Triggers{
	Push:        &MonorepoPush,
	PullRequest: &MonorepoPullRequest,
}
