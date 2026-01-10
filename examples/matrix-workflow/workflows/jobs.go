package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// TestMatrix defines Go versions and OS combinations to test.
var TestMatrix = workflow.Matrix{
	Values: map[string][]any{
		"go": {"1.22", "1.23"},
		"os": {"ubuntu-latest", "macos-latest"},
	},
}

// TestStrategy uses matrix to test multiple configurations.
var TestStrategy = workflow.Strategy{
	Matrix: &TestMatrix,
}

// Test runs tests across multiple Go versions and operating systems.
var Test = workflow.Job{
	Name:     "Test",
	RunsOn:   "${{ matrix.os }}",
	Strategy: &TestStrategy,
	Steps:    TestSteps,
}
