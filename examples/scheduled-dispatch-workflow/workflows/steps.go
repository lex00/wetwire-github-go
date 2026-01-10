package workflows

import (
	"github.com/lex00/wetwire-github-go/actions/checkout"
	"github.com/lex00/wetwire-github-go/workflow"
)

// MaintenanceSteps are the steps for scheduled maintenance tasks.
var MaintenanceSteps = []any{
	checkout.Checkout{},
	workflow.Step{
		Name: "Run cleanup tasks",
		Run:  "echo 'Running scheduled cleanup...'",
	},
	workflow.Step{
		Name: "Generate reports",
		Run:  "echo 'Generating daily reports...'",
	},
	workflow.Step{
		Name: "Check system health",
		Run:  "echo 'Checking system health...'",
	},
}

// DeploySteps are the steps for manual deployment.
var DeploySteps = []any{
	checkout.Checkout{},
	workflow.Step{
		Name: "Show deployment parameters",
		Run:  "echo \"Environment: ${{ inputs.environment }}\"\necho \"Dry run: ${{ inputs.dry_run }}\"\necho \"Version: ${{ inputs.version }}\"",
	},
	workflow.Step{
		Name: "Validate version",
		Run:  "echo \"Validating version: ${{ inputs.version }}\"",
	},
	workflow.Step{
		Name: "Deploy to environment",
		If:   "${{ inputs.dry_run == false }}",
		Run:  "echo \"Deploying ${{ inputs.version }} to ${{ inputs.environment }}...\"",
	},
	workflow.Step{
		Name: "Dry run deployment",
		If:   "${{ inputs.dry_run == true }}",
		Run:  "echo \"DRY RUN: Would deploy ${{ inputs.version }} to ${{ inputs.environment }}\"",
	},
}
