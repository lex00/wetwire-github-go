// Package workflows defines GitHub Actions workflow declarations.
package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// ArtifactPipeline is a multi-stage build, test, and release pipeline.
// It demonstrates artifact passing between jobs using upload/download actions.
var ArtifactPipeline = workflow.Workflow{
	Name: "Artifact Pipeline",
	On:   PipelineTriggers,
	Jobs: map[string]workflow.Job{
		"build":   Build,
		"test":    Test,
		"release": Release,
	},
}
