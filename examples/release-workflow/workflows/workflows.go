// Package workflows defines GitHub Actions workflow declarations.
package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// Release is the automated release workflow.
// It creates GitHub releases when version tags (v*) are pushed.
var Release = workflow.Workflow{
	Name: "Release",
	On:   ReleaseTriggers,
	Jobs: map[string]workflow.Job{
		"release": CreateRelease,
	},
}
