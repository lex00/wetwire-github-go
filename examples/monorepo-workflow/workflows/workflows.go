// Package workflows defines GitHub Actions workflow declarations for a monorepo.
package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// MonorepoCI is the main continuous integration workflow for a monorepo.
// It detects changes per service and only builds affected components.
var MonorepoCI = workflow.Workflow{
	Name: "Monorepo CI",
	On:   MonorepoTriggers,
	Jobs: map[string]workflow.Job{
		"detect-changes": DetectChanges,
		"build-api":      BuildAPI,
		"build-web":      BuildWeb,
		"build-shared":   BuildShared,
	},
}
