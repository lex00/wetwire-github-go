package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// AutoLabel automatically labels issues based on their content.
// This job runs when issues are opened.
var AutoLabel = workflow.Job{
	Name:   "Auto Label Issues",
	RunsOn: "ubuntu-latest",
	If:     "${{ github.event_name == 'issues' && github.event.action == 'opened' }}",
	Steps:  AutoLabelSteps,
}

// RespondToComment responds to issue comments with specific commands.
// This job runs when issue comments are created.
var RespondToComment = workflow.Job{
	Name:   "Respond to Comment",
	RunsOn: "ubuntu-latest",
	If:     "${{ github.event_name == 'issue_comment' && github.event.action == 'created' }}",
	Steps:  RespondToCommentSteps,
}

// EnforceReviewPolicy enforces review policies on pull requests.
// This job runs when PR reviews are submitted.
var EnforceReviewPolicy = workflow.Job{
	Name:   "Enforce Review Policy",
	RunsOn: "ubuntu-latest",
	If:     "${{ github.event_name == 'pull_request_review' && github.event.action == 'submitted' }}",
	Steps:  EnforceReviewPolicySteps,
}
