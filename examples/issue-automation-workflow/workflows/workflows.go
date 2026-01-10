// Package workflows defines GitHub Actions workflow declarations.
package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// IssueAutomation is a workflow that automates issue and PR management.
// It responds to issue creation, comments, and PR reviews with automated actions.
var IssueAutomation = workflow.Workflow{
	Name: "Issue Automation",
	On:   IssueAutomationTriggers,
	Jobs: map[string]workflow.Job{
		"auto-label":      AutoLabel,
		"respond-comment": RespondToComment,
		"enforce-review":  EnforceReviewPolicy,
	},
}
