package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// DeployStaging deploys the application to the staging environment.
// Only runs on the main branch after build and test jobs pass.
var DeployStaging = workflow.Job{
	Name:   "Deploy Staging",
	RunsOn: "ubuntu-latest",
	Needs:  []any{Build, Test},
	If:     "${{ github.ref == 'refs/heads/main' }}",
	Environment: &workflow.Environment{
		Name: "staging",
		URL:  "https://staging.example.com",
	},
	Steps: DeployStagingSteps,
}

// DeployStagingSteps defines the deployment steps for staging.
var DeployStagingSteps = []any{
	workflow.Step{Uses: "actions/checkout@v4"},
	workflow.Step{
		Run: "./scripts/deploy.sh",
		Env: map[string]any{
			"ENVIRONMENT": "staging",
			"DEPLOY_KEY":  "${{ secrets.STAGING_DEPLOY_KEY }}",
		},
	},
}

// DeployProduction deploys the application to the production environment.
// Requires manual approval and only runs on the main branch after build and test pass.
var DeployProduction = workflow.Job{
	Name:   "Deploy Production",
	RunsOn: "ubuntu-latest",
	Needs:  []any{Build, Test},
	If:     "${{ github.ref == 'refs/heads/main' }}",
	Environment: &workflow.Environment{
		Name: "production",
		URL:  "https://example.com",
	},
	Steps: DeployProductionSteps,
}

// DeployProductionSteps defines the deployment steps for production.
var DeployProductionSteps = []any{
	workflow.Step{Uses: "actions/checkout@v4"},
	workflow.Step{
		Run: "./scripts/deploy.sh",
		Env: map[string]any{
			"ENVIRONMENT": "production",
			"DEPLOY_KEY":  "${{ secrets.PRODUCTION_DEPLOY_KEY }}",
		},
	},
}
