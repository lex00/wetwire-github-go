package main

import "github.com/wetwire/wetwire-github/workflow"

var CICD = workflow.Workflow{
	Name: "CI/CD",
	On:   CICDTriggers,
	Jobs: map[string]any{
		"build":              BuildJob,
		"test":               TestJob,
		"deploy-staging":     DeployStagingJob,
		"deploy-production":  DeployProductionJob,
	},
}
