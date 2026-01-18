// Package workflows provides example GitHub Actions workflow declarations for a Go web application.
package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// CI is the main CI/CD workflow for the web application.
// It builds, tests, and deploys to staging and production environments.
var CI = workflow.Workflow{
	Name: "CI/CD",
	On:   CITriggers,
	Jobs: map[string]workflow.Job{
		"build":             Build,
		"test":              Test,
		"deploy-staging":    DeployStaging,
		"deploy-production": DeployProduction,
	},
}

// CITriggers defines when the CI/CD workflow runs.
var CITriggers = workflow.Triggers{
	Push: &workflow.PushTrigger{
		Branches: []string{"main"},
	},
	PullRequest: &workflow.PullRequestTrigger{
		Branches: []string{"main"},
	},
}
