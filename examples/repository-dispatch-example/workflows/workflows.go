// Package workflows defines GitHub Actions workflow declarations.
package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// APITriggeredDeploy is a workflow triggered by repository dispatch events.
// It demonstrates API-triggered deployments with custom parameters:
// - Event type filtering for different deployment targets
// - Client payload access for deployment configuration
// - Conditional jobs based on event type
var APITriggeredDeploy = workflow.Workflow{
	Name: "API Triggered Deploy",
	On:   DispatchTriggers,
	Jobs: map[string]workflow.Job{
		"validate":          Validate,
		"deploy":            Deploy,
		"deploy-staging":    DeployStaging,
		"deploy-production": DeployProduction,
	},
}
