package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// DevEnvironment configures the development deployment environment.
// No approval required - deploys automatically on push.
var DevEnvironment = workflow.Environment{
	Name: "dev",
	URL:  "https://dev.example.com",
}

// StagingEnvironment configures the staging deployment environment.
// Configure optional reviewers in GitHub repository settings.
var StagingEnvironment = workflow.Environment{
	Name: "staging",
	URL:  "https://staging.example.com",
}

// ProductionEnvironment configures the production deployment environment.
// Configure required reviewers in GitHub repository settings for approval gates.
var ProductionEnvironment = workflow.Environment{
	Name: "production",
	URL:  "https://example.com",
}

// DeployDev deploys the application to the development environment.
// This job runs automatically on push to main branch.
var DeployDev = workflow.Job{
	Name:        "Deploy to Dev",
	RunsOn:      "ubuntu-latest",
	Environment: &DevEnvironment,
	Outputs: map[string]any{
		"version": "${{ steps.build.outputs.version }}",
	},
	Steps: DeployDevSteps,
}

// PromoteStaging promotes the application to the staging environment.
// This job depends on the dev deployment completing successfully.
// Approval is controlled through GitHub environment protection rules.
var PromoteStaging = workflow.Job{
	Name:        "Promote to Staging",
	RunsOn:      "ubuntu-latest",
	Needs:       []any{DeployDev},
	Environment: &StagingEnvironment,
	Steps:       PromoteStagingSteps,
}

// PromoteProduction promotes the application to the production environment.
// This job depends on the staging promotion completing successfully.
// Approval is controlled through GitHub environment protection rules.
var PromoteProduction = workflow.Job{
	Name:        "Promote to Production",
	RunsOn:      "ubuntu-latest",
	Needs:       []any{PromoteStaging},
	Environment: &ProductionEnvironment,
	Steps:       PromoteProductionSteps,
}
