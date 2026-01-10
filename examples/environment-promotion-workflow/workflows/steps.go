package workflows

import (
	"github.com/lex00/wetwire-github-go/actions/checkout"
	"github.com/lex00/wetwire-github-go/actions/setup_go"
	"github.com/lex00/wetwire-github-go/workflow"
)

// CheckoutCode checks out the repository.
var CheckoutCode = checkout.Checkout{}

// SetupGoVersion sets up Go for building.
var SetupGoVersion = setup_go.SetupGo{
	GoVersion: "1.24",
}

// BuildBinary builds the application binary.
var BuildBinary = workflow.Step{
	Name: "Build binary",
	ID:   "build",
	Run: `VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
go build -ldflags="-X main.version=$VERSION" -o app ./...
echo "version=$VERSION" >> $GITHUB_OUTPUT`,
}

// RunTests runs the test suite.
var RunTests = workflow.Step{
	Name: "Run tests",
	Run:  "go test -v ./...",
}

// DeployDevStep deploys to the development environment.
var DeployDevStep = workflow.Step{
	Name: "Deploy to dev",
	Env: map[string]any{
		"DEPLOY_TOKEN": workflow.Secrets.Get("DEV_DEPLOY_TOKEN"),
		"ENVIRONMENT":  "dev",
		"API_URL":      "https://api-dev.example.com",
	},
	Run: `echo "Deploying to dev environment..."
echo "Using deploy token: ${DEPLOY_TOKEN:0:4}..."
echo "Deployment target: $ENVIRONMENT"
echo "API URL: $API_URL"
# Add your dev deployment commands here
# e.g., kubectl apply -f k8s/dev/`,
}

// VerifyDevStep verifies the dev deployment.
var VerifyDevStep = workflow.Step{
	Name: "Verify dev deployment",
	Run: `echo "Running smoke tests on dev..."
# Add your dev verification commands here
# e.g., curl -f https://dev.example.com/health`,
}

// DeployDevSteps are the steps for the dev deployment job.
var DeployDevSteps = []any{
	CheckoutCode,
	SetupGoVersion,
	BuildBinary,
	RunTests,
	DeployDevStep,
	VerifyDevStep,
}

// PromoteStagingStep promotes to the staging environment.
var PromoteStagingStep = workflow.Step{
	Name: "Promote to staging",
	Env: map[string]any{
		"DEPLOY_TOKEN": workflow.Secrets.Get("STAGING_DEPLOY_TOKEN"),
		"ENVIRONMENT":  "staging",
		"API_URL":      "https://api-staging.example.com",
	},
	Run: `echo "Promoting to staging environment..."
echo "Using deploy token: ${DEPLOY_TOKEN:0:4}..."
echo "Deployment target: $ENVIRONMENT"
echo "API URL: $API_URL"
# Add your staging deployment commands here
# e.g., kubectl apply -f k8s/staging/`,
}

// VerifyStagingStep verifies the staging deployment.
var VerifyStagingStep = workflow.Step{
	Name: "Verify staging deployment",
	Run: `echo "Running smoke tests on staging..."
# Add your staging verification commands here
# e.g., curl -f https://staging.example.com/health`,
}

// PromoteStagingSteps are the steps for the staging promotion job.
var PromoteStagingSteps = []any{
	CheckoutCode,
	PromoteStagingStep,
	VerifyStagingStep,
}

// PromoteProductionStep promotes to the production environment.
var PromoteProductionStep = workflow.Step{
	Name: "Promote to production",
	Env: map[string]any{
		"DEPLOY_TOKEN": workflow.Secrets.Get("PRODUCTION_DEPLOY_TOKEN"),
		"ENVIRONMENT":  "production",
		"API_URL":      "https://api.example.com",
	},
	Run: `echo "Promoting to production environment..."
echo "Using deploy token: ${DEPLOY_TOKEN:0:4}..."
echo "Deployment target: $ENVIRONMENT"
echo "API URL: $API_URL"
# Add your production deployment commands here
# e.g., kubectl apply -f k8s/production/`,
}

// VerifyProductionStep verifies the production deployment.
var VerifyProductionStep = workflow.Step{
	Name: "Verify production deployment",
	Run: `echo "Running smoke tests on production..."
# Add your production verification commands here
# e.g., curl -f https://example.com/health`,
}

// PromoteProductionSteps are the steps for the production promotion job.
var PromoteProductionSteps = []any{
	CheckoutCode,
	PromoteProductionStep,
	VerifyProductionStep,
}
