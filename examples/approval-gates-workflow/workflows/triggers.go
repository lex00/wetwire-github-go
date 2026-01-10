package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// ApprovalPush triggers on push to main branch.
var ApprovalPush = workflow.PushTrigger{
	Branches: []string{"main"},
}

// ApprovalDispatch allows manual workflow triggering.
var ApprovalDispatch = workflow.WorkflowDispatchTrigger{}

// ApprovalTriggers combines push and manual dispatch triggers.
var ApprovalTriggers = workflow.Triggers{
	Push:             &ApprovalPush,
	WorkflowDispatch: &ApprovalDispatch,
}
