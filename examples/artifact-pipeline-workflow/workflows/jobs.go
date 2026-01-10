package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// ReleasePermissions allows creating releases and writing contents.
var ReleasePermissions = workflow.Permissions{
	Contents: "write",
}

// Build compiles the application and uploads binaries as artifacts.
var Build = workflow.Job{
	Name:   "Build",
	RunsOn: "ubuntu-latest",
	Steps:  BuildSteps,
}

// Test runs the test suite after downloading build artifacts.
// Depends on Build job completing successfully.
var Test = workflow.Job{
	Name:   "Test",
	RunsOn: "ubuntu-latest",
	Needs:  []any{Build},
	Steps:  TestSteps,
}

// ReleaseCondition checks if the ref starts with a version tag.
var ReleaseCondition = workflow.StartsWith(
	workflow.GitHub.Ref(),
	workflow.Expression("'refs/tags/v'"),
)

// Release creates a GitHub release with binaries.
// Only runs on version tags (v*) after build and test succeed.
var Release = workflow.Job{
	Name:        "Release",
	RunsOn:      "ubuntu-latest",
	Needs:       []any{Build, Test},
	If:          ReleaseCondition,
	Permissions: &ReleasePermissions,
	Steps:       ReleaseSteps,
}
