package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// Build is the main build job in the reusable workflow.
// It accepts inputs and produces outputs that can be used by callers.
var Build = workflow.Job{
	Name:   "Build",
	RunsOn: "ubuntu-latest",
	Outputs: map[string]any{
		"artifact": "${{ steps.build.outputs.artifact }}",
		"version":  "${{ steps.build.outputs.version }}",
	},
	Steps: BuildSteps,
}

// UseOutput demonstrates a job that consumes outputs from a reusable workflow call.
// Note: This job references outputs via needs context from the caller workflow.
var UseOutput = workflow.Job{
	Name:   "Use Output",
	RunsOn: "ubuntu-latest",
	Needs:  []any{"call-build"},
	Steps:  UseOutputSteps,
}
