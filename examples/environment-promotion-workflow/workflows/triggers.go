package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// PromotionPush triggers on push to main branch.
var PromotionPush = workflow.PushTrigger{
	Branches: []string{"main"},
}

// PromotionDispatchInputs defines inputs for manual promotion triggering.
var PromotionDispatchInputs = map[string]workflow.WorkflowInput{
	"environment": {
		Description: "Target environment for promotion",
		Required:    true,
		Type:        "choice",
		Options:     []string{"dev", "staging", "production"},
		Default:     "dev",
	},
}

// PromotionDispatch allows manual workflow triggering with environment selection.
var PromotionDispatch = workflow.WorkflowDispatchTrigger{
	Inputs: PromotionDispatchInputs,
}

// PromotionTriggers combines push and manual dispatch triggers.
var PromotionTriggers = workflow.Triggers{
	Push:             &PromotionPush,
	WorkflowDispatch: &PromotionDispatch,
}
