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

// BuildSteps are the steps for the build job.
var BuildSteps = []any{
	CheckoutCode,
	SetupGoVersion,
	BuildBinary,
	RunTests,
}

// DeployStagingStep deploys to the staging environment.
var DeployStagingStep = workflow.Step{
	Name: "Deploy to staging",
	Env: map[string]any{
		"DEPLOY_TOKEN": workflow.Secrets.Get("STAGING_DEPLOY_TOKEN"),
		"ENVIRONMENT":  "staging",
	},
	Run: `echo "Deploying to staging environment..."
echo "Using deploy token: ${DEPLOY_TOKEN:0:4}..."
echo "Deployment target: $ENVIRONMENT"
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

// DeployStagingSteps are the steps for the staging deployment job.
var DeployStagingSteps = []any{
	CheckoutCode,
	DeployStagingStep,
	VerifyStagingStep,
}

// DeployProductionStep deploys to the production environment.
var DeployProductionStep = workflow.Step{
	Name: "Deploy to production",
	Env: map[string]any{
		"DEPLOY_TOKEN": workflow.Secrets.Get("PRODUCTION_DEPLOY_TOKEN"),
		"ENVIRONMENT":  "production",
	},
	Run: `echo "Deploying to production environment..."
echo "Using deploy token: ${DEPLOY_TOKEN:0:4}..."
echo "Deployment target: $ENVIRONMENT"
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

// DeployProductionSteps are the steps for the production deployment job.
var DeployProductionSteps = []any{
	CheckoutCode,
	DeployProductionStep,
	VerifyProductionStep,
}
