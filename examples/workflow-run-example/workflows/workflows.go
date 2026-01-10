// Package workflows defines GitHub Actions workflow declarations.
package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// DeployAfterCI is a workflow that runs after the CI workflow completes.
// It demonstrates the workflow_run trigger pattern for:
// - Deploying after CI passes
// - Downloading artifacts from the triggering workflow
// - Accessing workflow run context (conclusion, head_sha, etc.)
var DeployAfterCI = workflow.Workflow{
	Name: "Deploy After CI",
	On:   DeployTriggers,
	Jobs: map[string]workflow.Job{
		"deploy": Deploy,
		"notify": Notify,
	},
}
