package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// ReleasePermissions allows creating releases and writing contents.
var ReleasePermissions = workflow.Permissions{
	Contents: "write",
}

// CreateRelease creates a GitHub release with auto-generated notes.
var CreateRelease = workflow.Job{
	Name:        "Create Release",
	RunsOn:      "ubuntu-latest",
	Permissions: &ReleasePermissions,
	Steps:       ReleaseSteps,
}
