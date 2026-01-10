package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// StagingEnvironment configures the staging deployment environment.
var StagingEnvironment = workflow.Environment{
	Name: "staging",
	URL:  "https://staging.example.com",
}

// ApprovalGateEnvironment configures the production environment for approval gate.
// Configure required reviewers in GitHub repository settings to enforce manual approval.
var ApprovalGateEnvironment = workflow.Environment{
	Name: "production",
}

// ProductionEnvironment configures the production deployment environment.
var ProductionEnvironment = workflow.Environment{
	Name: "production",
	URL:  "https://example.com",
}

// BuildJob compiles the application and runs tests.
var BuildJob = workflow.Job{
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
	Needs:       []any{BuildJob},
	Environment: &StagingEnvironment,
	Steps:       DeployStagingSteps,
}

// ApproveProduction is an explicit approval gate job.
// This job uses the production environment to trigger GitHub's approval flow.
// The job depends on staging deployment and blocks production until approved.
var ApproveProduction = workflow.Job{
	Name:        "Approve Production Deployment",
	RunsOn:      "ubuntu-latest",
	Needs:       []any{DeployStaging},
	Environment: &ApprovalGateEnvironment,
	Steps:       ApprovalSteps,
}

// DeployProduction deploys the application to the production environment.
// This job depends on the approval gate job completing successfully.
var DeployProduction = workflow.Job{
	Name:        "Deploy to Production",
	RunsOn:      "ubuntu-latest",
	Needs:       []any{ApproveProduction},
	Environment: &ProductionEnvironment,
	Steps:       DeployProductionSteps,
}
