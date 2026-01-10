// Package workflows defines GitHub Actions workflow declarations.
package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// Deploy is the multi-environment deployment workflow.
// It demonstrates a sequential deployment pipeline:
// build -> deploy to staging -> deploy to production.
//
// Manual approval gates for production are configured through
// GitHub environment protection rules, not in the workflow definition.
var Deploy = workflow.Workflow{
	Name: "Deploy",
	On:   DeployTriggers,
	Jobs: map[string]workflow.Job{
		"build":             Build,
		"deploy-staging":    DeployStaging,
		"deploy-production": DeployProduction,
	},
}
