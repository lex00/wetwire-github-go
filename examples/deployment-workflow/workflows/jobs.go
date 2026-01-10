package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// StagingEnvironment configures the staging deployment environment.
// Environment protection rules can be configured in GitHub repository settings.
var StagingEnvironment = workflow.Environment{
	Name: "staging",
	URL:  "https://staging.example.com",
}

// ProductionEnvironment configures the production deployment environment.
// Configure required reviewers and wait timers in GitHub repository settings
// to enforce manual approval gates before production deployments.
var ProductionEnvironment = workflow.Environment{
	Name: "production",
	URL:  "https://example.com",
}

// Build compiles the application and runs tests.
var Build = workflow.Job{
	Name:   "Build",
	RunsOn: "ubuntu-latest",
	Outputs: map[string]any{
		"version": "${{ steps.build.outputs.version }}",
	},
	Steps: BuildSteps,
}

// DeployStaging deploys the application to the staging environment.
// This job depends on the build job completing successfully.
var DeployStaging = workflow.Job{
	Name:        "Deploy to Staging",
	RunsOn:      "ubuntu-latest",
	Needs:       []any{Build},
	Environment: &StagingEnvironment,
	Steps:       DeployStagingSteps,
}

// DeployProduction deploys the application to the production environment.
// This job depends on the staging deployment completing successfully.
// Manual approval is enforced through environment protection rules.
var DeployProduction = workflow.Job{
	Name:        "Deploy to Production",
	RunsOn:      "ubuntu-latest",
	Needs:       []any{DeployStaging},
	Environment: &ProductionEnvironment,
	Steps:       DeployProductionSteps,
}
