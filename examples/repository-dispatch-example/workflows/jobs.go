package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// ValidateOutputs defines the outputs from the validation job.
var ValidateOutputs = map[string]any{
	"environment": "${{ steps.validate.outputs.environment }}",
	"version":     "${{ steps.validate.outputs.version }}",
	"ref":         "${{ steps.validate.outputs.ref }}",
}

// Validate validates the dispatch event payload before deployment.
// Runs for all dispatch event types.
var Validate = workflow.Job{
	Name:    "Validate Dispatch",
	RunsOn:  "ubuntu-latest",
	Outputs: ValidateOutputs,
	Steps:   ValidateSteps,
}

// Deploy runs the main deployment job after validation.
// Uses the generic "deploy" event type.
var Deploy = workflow.Job{
	Name:   "Deploy",
	RunsOn: "ubuntu-latest",
	Needs:  []any{Validate},
	If:     "${{ github.event.action == 'deploy' }}",
	Steps:  DeploySteps,
}

// DeployStaging handles staging-specific deployments.
// Triggered by the "deploy-staging" event type.
var DeployStaging = workflow.Job{
	Name:   "Deploy Staging",
	RunsOn: "ubuntu-latest",
	If:     "${{ github.event.action == 'deploy-staging' }}",
	Steps:  StagingDeploySteps,
}

// DeployProduction handles production deployments with extra verification.
// Triggered by the "deploy-production" event type.
var DeployProduction = workflow.Job{
	Name:   "Deploy Production",
	RunsOn: "ubuntu-latest",
	If:     "${{ github.event.action == 'deploy-production' }}",
	Steps:  ProductionDeploySteps,
}
