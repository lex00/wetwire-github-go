package main

import "github.com/wetwire/wetwire-github-go/workflow"

// Main CI/CD Workflow
var CI = workflow.Workflow{
	Name: "CI/CD Pipeline",
	On:   CITriggers,
	Jobs: map[string]workflow.Job{
		"build":              Build,
		"test":               Test,
		"lint":               Lint,
		"deploy-staging":     DeployStaging,
		"deploy-production":  DeployProduction,
	},
}
