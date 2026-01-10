package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// IssuesOpened triggers when issues are opened or labeled.
// Used for auto-labeling based on issue content.
var IssuesOpened = workflow.IssuesTrigger{
	Types: []string{"opened", "labeled"},
}

// CommentCreated triggers when issue comments are created.
// Used for responding to issue comments.
var CommentCreated = workflow.IssueCommentTrigger{
	Types: []string{"created"},
}

// ReviewSubmitted triggers when PR reviews are submitted.
// Used for enforcing review policies.
var ReviewSubmitted = workflow.PullRequestReviewTrigger{
	Types: []string{"submitted"},
}

// IssueAutomationTriggers combines all issue automation triggers.
var IssueAutomationTriggers = workflow.Triggers{
	Issues:            &IssuesOpened,
	IssueComment:     &CommentCreated,
	PullRequestReview: &ReviewSubmitted,
}
