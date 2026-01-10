package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// DockerPermissions allows pushing to GHCR.
var DockerPermissions = workflow.Permissions{
	Contents: "read",
	Packages: "write",
}

// DockerPublish builds and pushes Docker images to GHCR.
var DockerPublish = workflow.Job{
	Name:        "Publish Docker Image",
	RunsOn:      "ubuntu-latest",
	Permissions: &DockerPermissions,
	Steps:       DockerPublishSteps,
}

// ReleasePermissions allows creating releases.
var ReleasePermissions = workflow.Permissions{
	Contents: "write",
}

// CreateRelease creates a GitHub release with auto-generated notes.
var CreateRelease = workflow.Job{
	Name:        "Create Release",
	RunsOn:      "ubuntu-latest",
	Permissions: &ReleasePermissions,
	Steps:       CreateReleaseSteps,
}

// BuildMatrix defines the platforms for multi-platform builds.
var BuildMatrix = workflow.Matrix{
	Values: map[string][]any{
		"goos":   {"linux", "darwin", "windows"},
		"goarch": {"amd64", "arm64"},
	},
	Exclude: []map[string]any{
		{"goos": "windows", "goarch": "arm64"},
	},
}

// BuildStrategy configures the matrix strategy for builds.
var BuildStrategy = workflow.Strategy{
	Matrix: &BuildMatrix,
}

// ArtifactPermissions allows uploading artifacts and release assets.
var ArtifactPermissions = workflow.Permissions{
	Contents: "write",
}

// BuildArtifacts builds multi-platform binaries and uploads to release.
var BuildArtifacts = workflow.Job{
	Name:        "Build Artifacts",
	RunsOn:      "ubuntu-latest",
	Strategy:    &BuildStrategy,
	Permissions: &ArtifactPermissions,
	Steps:       BuildArtifactSteps,
}

// Notify sends announcements after release artifacts are ready.
var Notify = workflow.Job{
	Name:   "Notify",
	RunsOn: "ubuntu-latest",
	Needs:  []any{BuildArtifacts},
	Steps:  NotifySteps,
}
