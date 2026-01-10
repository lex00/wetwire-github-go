package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// BuildMatrix defines Go versions and OS combinations to test.
var BuildMatrix = workflow.Matrix{
	Values: map[string][]any{
		"go": {"1.23", "1.24"},
		"os": {"ubuntu-latest", "macos-latest"},
	},
}

// BuildStrategy uses matrix to test multiple configurations.
var BuildStrategy = workflow.Strategy{
	Matrix: &BuildMatrix,
}

// Build compiles and tests the project across multiple Go versions.
var Build = workflow.Job{
	Name:     "Build and Test",
	RunsOn:   "${{ matrix.os }}",
	Strategy: &BuildStrategy,
	Steps:    BuildSteps,
}

// Lint runs static analysis on the codebase.
var Lint = workflow.Job{
	Name:   "Lint",
	RunsOn: "ubuntu-latest",
	Steps:  LintSteps,
}
