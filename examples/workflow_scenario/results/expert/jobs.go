package main

import (
	"github.com/wetwire/wetwire-github/actions/cache"
	"github.com/wetwire/wetwire-github/actions/checkout"
	"github.com/wetwire/wetwire-github/actions/setup_go"
	"github.com/wetwire/wetwire-github/workflow"
)

// Build job with matrix strategy
var BuildJob = workflow.Job{
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
	checkout.Checkout{
		With: checkout.CheckoutInputs{
			FetchDepth: 0,
		},
	},
	setup_go.SetupGo{
		With: setup_go.SetupGoInputs{
			GoVersion: workflow.Matrix.Get("go"),
		},
	},
	cache.Cache{
		With: cache.CacheInputs{
			Path: "~/go/pkg/mod",
			Key:  "go-mod-${{ hashFiles('**/go.sum') }}",
		},
	},
	workflow.Step{
		Name: "Build",
		Run:  "go build ./...",
	},
}

// Test job
var TestJob = workflow.Job{
	Name:   "test",
	RunsOn: "ubuntu-latest",
	Steps:  TestSteps,
}

var TestSteps = []any{
	checkout.Checkout{},
	setup_go.SetupGo{
		With: setup_go.SetupGoInputs{
			GoVersion: "1.23",
		},
	},
	workflow.Step{
		Name: "Run tests with coverage",
		Run:  "go test -race -coverprofile=coverage.out ./...",
	},
}

// Deploy staging job
var DeployStagingJob = workflow.Job{
	Name:        "deploy-staging",
	RunsOn:      "ubuntu-latest",
	Needs:       []any{BuildJob, TestJob},
	If:          "github.ref == 'refs/heads/main'",
	Environment: StagingEnvironment,
	Steps:       DeployStagingSteps,
}

var StagingEnvironment = workflow.Environment{
	Name: "staging",
}

var DeployStagingSteps = []any{
	checkout.Checkout{},
	workflow.Step{
		Name: "Deploy to staging",
		Run:  "echo \"Deploying to staging\"",
	},
}

// Deploy production job
var DeployProductionJob = workflow.Job{
	Name:        "deploy-production",
	RunsOn:      "ubuntu-latest",
	Needs:       []any{BuildJob, TestJob},
	If:          "github.ref == 'refs/heads/main'",
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
		Name: "Deploy to production",
		Run:  "echo \"Deploying to production\"",
	},
}
