// Package workflows defines GitHub Actions workflow declarations for container services.
package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// ContainerServices is a workflow demonstrating container and service configurations.
// It runs integration tests against PostgreSQL and Redis service containers.
var ContainerServices = workflow.Workflow{
	Name: "Container Services",
	On:   CITriggers,
	Jobs: map[string]workflow.Job{
		"unit-test":        UnitTest,
		"integration-test": IntegrationTest,
	},
}
