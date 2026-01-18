package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// BuildMatrix defines the matrix strategy for testing multiple Go versions and operating systems.
var BuildMatrix = workflow.Matrix{
	Values: map[string][]any{
		"go": {"1.23", "1.24"},
		"os": {"ubuntu-latest", "macos-latest"},
	},
}

// BuildStrategy configures the matrix build strategy.
var BuildStrategy = workflow.Strategy{
	Matrix: &BuildMatrix,
}

// Build compiles and tests the application across multiple Go versions and operating systems.
// Uses matrix strategy to test on Go 1.23 and 1.24, on both Ubuntu and macOS.
var Build = workflow.Job{
	Name:     "Build",
	RunsOn:   "${{ matrix.os }}",
	Strategy: &BuildStrategy,
	Steps:    BuildSteps,
}

// BuildSteps defines the steps for the build job.
var BuildSteps = []any{
	workflow.Step{Uses: "actions/checkout@v4"},
	workflow.Step{
		Uses: "actions/setup-go@v5",
		With: map[string]any{
			"go-version": "${{ matrix.go }}",
		},
	},
	workflow.Step{
		Uses: "actions/cache@v4",
		With: map[string]any{
			"path": "~/go/pkg/mod",
			"key":  "go-${{ hashFiles('**/go.sum') }}",
		},
	},
	workflow.Step{Run: "go build ./..."},
	workflow.Step{Run: "go vet ./..."},
}
