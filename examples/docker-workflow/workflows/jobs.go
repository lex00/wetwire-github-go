package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// BuildPermissions allows pushing to GHCR.
var BuildPermissions = workflow.Permissions{
	Contents: "read",
	Packages: "write",
}

// Build compiles and optionally pushes the Docker image.
var Build = workflow.Job{
	Name:        "Build and Push",
	RunsOn:      "ubuntu-latest",
	Permissions: &BuildPermissions,
	Steps:       BuildSteps,
}
