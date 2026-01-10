// Package workflows defines GitHub Actions workflow declarations.
package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// BuildReusable is a reusable workflow that can be called from other workflows.
// It accepts inputs (go_version, run_tests, build_target) and produces outputs
// (artifact_name, build_version) that callers can use.
var BuildReusable = workflow.Workflow{
	Name: "Reusable Build Workflow",
	On:   ReusableTriggers,
	Jobs: map[string]workflow.Job{
		"build": Build,
	},
}

// CICaller is a workflow that demonstrates calling a reusable workflow.
// Note: The actual workflow call using `uses` at the job level requires
// manual YAML editing, as the current type system focuses on workflow_call
// trigger definitions. This workflow shows the pattern for consuming outputs.
var CICaller = workflow.Workflow{
	Name: "CI Caller",
	On:   CallerTriggers,
	Jobs: map[string]workflow.Job{
		"use-output": UseOutput,
	},
}
