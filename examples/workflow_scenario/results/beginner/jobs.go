package main

import (
	"github.com/wetwire/wetwire-github-go/actions/checkout"
	"github.com/wetwire/wetwire-github-go/actions/setup_go"
	"github.com/wetwire/wetwire-github-go/workflow"
)

// Build Job - Runs on matrix of Go versions and OS
var Build = workflow.Job{
	Name:     "build",
	RunsOn:   workflow.Matrix.Get("os"),
	Strategy: BuildStrategy,
	Steps:    BuildSteps,
}

var BuildStrategy = workflow.Strategy{
	Matrix: BuildMatrix,
}

var BuildMatrix = workflow.Matrix{
	Values: map[string][]any{
		"go": {"1.23", "1.24"},
		"os": {"ubuntu-latest", "macos-latest"},
	},
}

var BuildSteps = []any{
	checkout.Checkout{},
	setup_go.SetupGo{
		GoVersion: workflow.Matrix.Get("go"),
	},
	workflow.Step{
		Name: "Build",
		Run:  "go build ./...",
	},
}

// Test Job - Runs tests with coverage
var Test = workflow.Job{
	Name:     "test",
	RunsOn:   workflow.Matrix.Get("os"),
	Strategy: BuildStrategy,
	Steps:    TestSteps,
}

var TestSteps = []any{
	checkout.Checkout{},
	setup_go.SetupGo{
		GoVersion: workflow.Matrix.Get("go"),
	},
	workflow.Step{
		Name: "Run tests with coverage",
		Run:  "go test -v -race -coverprofile=coverage.out -covermode=atomic ./...",
	},
	workflow.Step{
		Name: "Upload coverage",
		Run:  "go tool cover -func=coverage.out",
	},
}

// Lint Job - Runs linter checks
var Lint = workflow.Job{
	Name:   "lint",
	RunsOn: "ubuntu-latest",
	Steps:  LintSteps,
}

var LintSteps = []any{
	checkout.Checkout{},
	setup_go.SetupGo{
		GoVersion: "1.23",
	},
	workflow.Step{
		Name: "golangci-lint",
		Run:  "go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest && golangci-lint run ./...",
	},
}

// Deploy Staging Job - Deploys to staging after build and test pass
var DeployStaging = workflow.Job{
	Name:   "deploy-staging",
	RunsOn: "ubuntu-latest",
	Needs:  []any{Build, Test, Lint},
	If:     workflow.Branch("main"),
	Steps:  DeployStagingSteps,
}

var DeployStagingSteps = []any{
	checkout.Checkout{},
	workflow.Step{
		Name: "Deploy to Staging",
		Run:  "echo 'Deploying to staging environment...'",
		Env: workflow.Env{
			"DEPLOY_TOKEN": workflow.Secrets.Get("STAGING_DEPLOY_TOKEN"),
			"ENVIRONMENT":  "staging",
		},
	},
	workflow.Step{
		Name: "Verify Deployment",
		Run:  "echo 'Verifying staging deployment...'",
	},
}

// Deploy Production Job - Deploys to production with manual approval
var DeployProduction = workflow.Job{
	Name:        "deploy-production",
	RunsOn:      "ubuntu-latest",
	Needs:       []any{DeployStaging},
	If:          workflow.Branch("main"),
	Environment: ProductionEnvironment,
	Steps:       DeployProductionSteps,
}

var ProductionEnvironment = workflow.Environment{
	Name: "production",
	URL:  "https://example.com",
}

var DeployProductionSteps = []any{
	checkout.Checkout{},
	workflow.Step{
		Name: "Deploy to Production",
		Run:  "echo 'Deploying to production environment...'",
		Env: workflow.Env{
			"DEPLOY_TOKEN": workflow.Secrets.Get("PRODUCTION_DEPLOY_TOKEN"),
			"ENVIRONMENT":  "production",
		},
	},
	workflow.Step{
		Name: "Verify Production Deployment",
		Run:  "echo 'Verifying production deployment...'",
	},
	workflow.Step{
		Name: "Notify Team",
		Run:  "echo 'Production deployment completed successfully!'",
	},
}
