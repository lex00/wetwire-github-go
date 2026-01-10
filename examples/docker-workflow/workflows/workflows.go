// Package workflows defines GitHub Actions workflow declarations.
package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// Docker is the main Docker build and push workflow.
// It builds images on PRs and pushes to GHCR on main branch.
var Docker = workflow.Workflow{
	Name: "Docker",
	On:   DockerTriggers,
	Jobs: map[string]workflow.Job{
		"build": Build,
	},
}
