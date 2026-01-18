package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// Test runs the test suite with coverage reporting.
// Executes tests with race detection and generates a coverage profile.
var Test = workflow.Job{
	Name:   "Test",
	RunsOn: "ubuntu-latest",
	Steps:  TestSteps,
}

// TestSteps defines the steps for running tests with coverage.
var TestSteps = []any{
	workflow.Step{Uses: "actions/checkout@v4"},
	workflow.Step{
		Uses: "actions/setup-go@v5",
		With: map[string]any{
			"go-version": "1.24",
		},
	},
	workflow.Step{
		Uses: "actions/cache@v4",
		With: map[string]any{
			"path": "~/go/pkg/mod",
			"key":  "go-${{ hashFiles('**/go.sum') }}",
		},
	},
	workflow.Step{Run: "go test -v -race -coverprofile=coverage.out ./..."},
	workflow.Step{
		Uses: "actions/upload-artifact@v4",
		With: map[string]any{
			"name": "coverage-report",
			"path": "coverage.out",
		},
	},
}
