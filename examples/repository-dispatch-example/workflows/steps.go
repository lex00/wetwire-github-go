package workflows

import (
	"github.com/lex00/wetwire-github-go/actions/checkout"
	"github.com/lex00/wetwire-github-go/workflow"
)

// ValidateSteps validate the dispatch event and its payload.
var ValidateSteps = []any{
	// Display event information
	workflow.Step{
		Name: "Show Event Info",
		Run: `echo "Event type: ${{ github.event.action }}"
echo "Client payload: ${{ toJSON(github.event.client_payload) }}"`,
	},
	// Validate required payload fields
	workflow.Step{
		Name: "Validate Payload",
		ID:   "validate",
		Run: `# Extract payload fields
environment="${{ github.event.client_payload.environment }}"
version="${{ github.event.client_payload.version }}"
ref="${{ github.event.client_payload.ref }}"

# Validate required fields
if [ -z "$environment" ]; then
  echo "::error::Missing required field: environment"
  exit 1
fi

if [ -z "$version" ]; then
  echo "::warning::No version specified, using 'latest'"
  version="latest"
fi

if [ -z "$ref" ]; then
  ref="main"
fi

# Set outputs
echo "environment=$environment" >> $GITHUB_OUTPUT
echo "version=$version" >> $GITHUB_OUTPUT
echo "ref=$ref" >> $GITHUB_OUTPUT
echo "Validation successful!"`,
	},
}

// DeploySteps perform the actual deployment using client_payload data.
var DeploySteps = []any{
	// Checkout at the specified ref
	checkout.Checkout{
		Ref: "${{ needs.validate.outputs.ref }}",
	},
	// Show deployment configuration
	workflow.Step{
		Name: "Show Deployment Config",
		Run: `echo "Deploying version: ${{ needs.validate.outputs.version }}"
echo "Target environment: ${{ needs.validate.outputs.environment }}"
echo "Git ref: ${{ needs.validate.outputs.ref }}"
echo "Triggered by: ${{ github.event.sender.login }}"`,
	},
	// Simulate deployment
	workflow.Step{
		Name: "Deploy Application",
		Env: workflow.Env{
			"ENVIRONMENT": "${{ needs.validate.outputs.environment }}",
			"VERSION":     "${{ needs.validate.outputs.version }}",
			"DEPLOY_KEY":  "${{ secrets.DEPLOY_KEY }}",
		},
		Run: `echo "Starting deployment to $ENVIRONMENT..."
echo "Version: $VERSION"
# In a real workflow, this would call deployment scripts
echo "Deployment complete!"`,
	},
}

// StagingDeploySteps handle staging-specific deployment.
var StagingDeploySteps = []any{
	checkout.Checkout{
		Ref: "${{ github.event.client_payload.ref || 'main' }}",
	},
	workflow.Step{
		Name: "Deploy to Staging",
		Run: `echo "Deploying to STAGING environment"
echo "Version: ${{ github.event.client_payload.version || 'latest' }}"
echo "Features: ${{ github.event.client_payload.features || 'all' }}"
echo "Staging deployment complete!"`,
	},
}

// ProductionDeploySteps handle production-specific deployment with extra safety.
var ProductionDeploySteps = []any{
	checkout.Checkout{
		Ref: "${{ github.event.client_payload.ref || 'main' }}",
	},
	workflow.Step{
		Name: "Verify Production Deployment",
		Run: `echo "Verifying production deployment prerequisites..."
echo "Version: ${{ github.event.client_payload.version }}"
echo "Approver: ${{ github.event.client_payload.approved_by }}"

# Require explicit approval for production
if [ -z "${{ github.event.client_payload.approved_by }}" ]; then
  echo "::error::Production deployments require approval"
  exit 1
fi

echo "Production deployment approved!"`,
	},
	workflow.Step{
		Name: "Deploy to Production",
		Run: `echo "Deploying to PRODUCTION environment"
echo "Version: ${{ github.event.client_payload.version }}"
echo "Approved by: ${{ github.event.client_payload.approved_by }}"
echo "Production deployment complete!"`,
	},
}
