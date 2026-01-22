package cicd

import (
	"github.com/wetwire/wetwire-github-go/pkg/actions/cache"
	"github.com/wetwire/wetwire-github-go/pkg/actions/checkout"
	"github.com/wetwire/wetwire-github-go/pkg/actions/setup_go"
	"github.com/wetwire/wetwire-github-go/pkg/workflow"
)

// Build job - matrix testing on multiple Go versions and OS
var BuildMatrix = workflow.Matrix{
	Values: map[string][]any{
		"go": {"1.23", "1.24"},
		"os": {"ubuntu-latest", "macos-latest"},
	},
}

var BuildStrategy = workflow.Strategy{
	Matrix: BuildMatrix,
}

var BuildSteps = []any{
	checkout.Checkout{},
	setup_go.SetupGo{
		GoVersion: workflow.Matrix.Get("go"),
	},
	cache.Cache{
		Path: "~/go/pkg/mod",
		Key:  workflow.Expr("${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}"),
		RestoreKeys: List(
			workflow.Expr("${{ runner.os }}-go-"),
		),
	},
	workflow.Step{
		Name: "Build",
		Run:  "go build ./...",
	},
}

var Build = workflow.Job{
	Name:     "build",
	RunsOn:   workflow.Matrix.Get("os"),
	Strategy: BuildStrategy,
	Steps:    BuildSteps,
}

// Test job - run tests with coverage
var TestSteps = []any{
	checkout.Checkout{},
	setup_go.SetupGo{
		GoVersion: "1.23",
	},
	cache.Cache{
		Path: "~/go/pkg/mod",
		Key:  workflow.Expr("${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}"),
		RestoreKeys: List(
			workflow.Expr("${{ runner.os }}-go-"),
		),
	},
	workflow.Step{
		Name: "Run tests with coverage",
		Run:  "go test -v -coverprofile=coverage.out ./...",
	},
	workflow.Step{
		Name: "Display coverage",
		Run:  "go tool cover -func=coverage.out",
	},
}

var Test = workflow.Job{
	Name:   "test",
	RunsOn: "ubuntu-latest",
	Steps:  TestSteps,
}

// Deploy Staging job - deploys to staging on main branch
var DeployStagingSteps = []any{
	checkout.Checkout{},
	workflow.Step{
		Name: "Deploy to staging",
		Run:  "echo 'Deploying to staging environment...'",
		Env: workflow.Env{
			"ENVIRONMENT": "staging",
		},
	},
}

var DeployStaging = workflow.Job{
	Name:   "deploy-staging",
	RunsOn: "ubuntu-latest",
	Needs:  []any{Build, Test},
	If:     workflow.Expr("github.ref == 'refs/heads/main'"),
	Steps:  DeployStagingSteps,
}

// Deploy Production job - requires manual approval via environment gate
var DeployProductionEnvironment = workflow.Environment{
	Name: "production",
	URL:  "https://example.com",
}

var DeployProductionSteps = []any{
	checkout.Checkout{},
	workflow.Step{
		Name: "Deploy to production",
		Run:  "echo 'Deploying to production environment...'",
		Env: workflow.Env{
			"ENVIRONMENT": "production",
		},
	},
}

var DeployProduction = workflow.Job{
	Name:        "deploy-production",
	RunsOn:      "ubuntu-latest",
	Needs:       []any{Build, Test},
	If:          workflow.Expr("github.ref == 'refs/heads/main'"),
	Environment: DeployProductionEnvironment,
	Steps:       DeployProductionSteps,
}
