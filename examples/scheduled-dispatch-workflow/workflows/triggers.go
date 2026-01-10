package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// DailySchedule runs workflow daily at midnight UTC.
var DailySchedule = workflow.ScheduleTrigger{
	Cron: "0 0 * * *",
}

// DispatchInputs defines the inputs for manual workflow triggers.
var DispatchInputs = map[string]workflow.WorkflowInput{
	"environment": {
		Description: "Target deployment environment",
		Required:    true,
		Type:        "choice",
		Options:     []string{"dev", "staging", "prod"},
		Default:     "dev",
	},
	"dry_run": {
		Description: "Run without making actual changes",
		Required:    false,
		Type:        "boolean",
		Default:     false,
	},
	"version": {
		Description: "Version to deploy (e.g., v1.2.3)",
		Required:    false,
		Type:        "string",
		Default:     "latest",
	},
}

// ManualDispatch allows manual triggering with typed inputs.
var ManualDispatch = workflow.WorkflowDispatchTrigger{
	Inputs: DispatchInputs,
}

// WorkflowTriggers combines schedule and dispatch triggers.
var WorkflowTriggers = workflow.Triggers{
	Schedule:         []workflow.ScheduleTrigger{DailySchedule},
	WorkflowDispatch: &ManualDispatch,
}
