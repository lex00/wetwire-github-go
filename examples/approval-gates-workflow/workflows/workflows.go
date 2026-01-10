// Package workflows defines GitHub Actions workflow declarations.
package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// ApprovalGates is the deployment workflow with explicit approval gates.
// It demonstrates a four-stage deployment pipeline:
// build -> deploy staging -> approve production -> deploy production.
//
// The ApproveProduction job acts as an explicit approval gate, using
// GitHub environment protection rules to require manual approval before
// proceeding to production deployment.
var ApprovalGates = workflow.Workflow{
	Name: "Approval Gates",
	On:   ApprovalTriggers,
	Jobs: map[string]workflow.Job{
		"build":              BuildJob,
		"deploy-staging":     DeployStaging,
		"approve-production": ApproveProduction,
		"deploy-production":  DeployProduction,
	},
}
