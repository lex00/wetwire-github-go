// Package workflows defines GitHub Actions workflow declarations.
package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// Matrix is the multi-OS/version testing workflow.
// It tests across Go versions and operating systems.
var Matrix = workflow.Workflow{
	Name: "Matrix Test",
	On:   MatrixTriggers,
	Jobs: map[string]workflow.Job{
		"test": Test,
	},
}
