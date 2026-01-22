package cicd

import (
	"github.com/wetwire/wetwire-github-go/pkg/workflow"
)

var CI = workflow.Workflow{
	Name: "CI/CD",
	On:   CITriggers,
	Jobs: map[string]workflow.Job{
		"build":              Build,
		"test":               Test,
		"deploy-staging":     DeployStaging,
		"deploy-production":  DeployProduction,
	},
}
