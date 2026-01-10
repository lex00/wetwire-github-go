// Package workflows defines GitHub Actions workflow declarations.
package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// Publish is the main publishing workflow.
// It builds and pushes Docker images, creates releases, and tags Go modules.
// Triggered on version tag pushes (v*).
var Publish = workflow.Workflow{
	Name: "Publish",
	On:   PublishTriggers,
	Jobs: map[string]workflow.Job{
		"docker":  DockerPublish,
		"release": CreateRelease,
	},
}

// Release is the workflow that runs after a release is published.
// It builds multi-platform artifacts and uploads them to the release.
var Release = workflow.Workflow{
	Name: "Release",
	On:   ReleaseTriggers,
	Jobs: map[string]workflow.Job{
		"build-artifacts": BuildArtifacts,
		"notify":          Notify,
	},
}
