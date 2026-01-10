package serialize_test

import (
	"testing"

	"github.com/lex00/wetwire-github-go/internal/serialize"
	"github.com/lex00/wetwire-github-go/workflow"
)

// createBenchmarkWorkflow creates a realistic workflow for benchmarking.
func createBenchmarkWorkflow() *workflow.Workflow {
	return &workflow.Workflow{
		Name: "CI",
		On: workflow.Triggers{
			Push: &workflow.PushTrigger{
				Branches: []string{"main", "develop"},
				Paths:    []string{"src/**", "*.go"},
			},
			PullRequest: &workflow.PullRequestTrigger{
				Types:    []string{"opened", "synchronize", "reopened"},
				Branches: []string{"main"},
			},
		},
		Env: workflow.Env{
			"GO_VERSION": "1.23",
			"CI":         "true",
		},
		Permissions: &workflow.Permissions{
			Contents:     "read",
			PullRequests: "write",
		},
		Concurrency: &workflow.Concurrency{
			Group:            "ci-${{ github.ref }}",
			CancelInProgress: true,
		},
		Jobs: map[string]workflow.Job{
			"build": {
				Name:   "Build",
				RunsOn: "ubuntu-latest",
				Steps: []any{
					workflow.Step{Uses: "actions/checkout@v4"},
					workflow.Step{
						Uses: "actions/setup-go@v5",
						With: workflow.With{"go-version": "1.23"},
					},
					workflow.Step{Run: "go build ./..."},
				},
			},
			"test": {
				Name:   "Test",
				RunsOn: "ubuntu-latest",
				Needs:  []any{"build"},
				Steps: []any{
					workflow.Step{Uses: "actions/checkout@v4"},
					workflow.Step{
						Uses: "actions/setup-go@v5",
						With: workflow.With{"go-version": "1.23"},
					},
					workflow.Step{Run: "go test -v ./..."},
				},
			},
			"lint": {
				Name:   "Lint",
				RunsOn: "ubuntu-latest",
				Steps: []any{
					workflow.Step{Uses: "actions/checkout@v4"},
					workflow.Step{
						Uses: "golangci/golangci-lint-action@v4",
						With: workflow.With{"version": "latest"},
					},
				},
			},
		},
	}
}

// createBenchmarkJob creates a realistic job for benchmarking.
func createBenchmarkJob() *workflow.Job {
	return &workflow.Job{
		Name:   "Test",
		RunsOn: "ubuntu-latest",
		Env: workflow.Env{
			"GOPATH":      "/home/runner/go",
			"GOCACHE":     "/home/runner/.cache/go-build",
			"GOLANGCI_LINT_CACHE": "/home/runner/.cache/golangci-lint",
		},
		Strategy: &workflow.Strategy{
			Matrix: &workflow.Matrix{
				Values: map[string][]any{
					"go":       {"1.22", "1.23"},
					"os":       {"ubuntu-latest", "macos-latest", "windows-latest"},
					"include":  {map[string]any{"go": "1.23", "os": "ubuntu-latest", "primary": true}},
				},
			},
			FailFast: workflow.Ptr(false),
		},
		TimeoutMinutes: 30,
		Steps: []any{
			workflow.Step{
				Name: "Checkout",
				Uses: "actions/checkout@v4",
				With: workflow.With{
					"fetch-depth": 0,
					"submodules":  "recursive",
				},
			},
			workflow.Step{
				Name: "Setup Go",
				Uses: "actions/setup-go@v5",
				With: workflow.With{
					"go-version": "${{ matrix.go }}",
					"cache":      true,
				},
			},
			workflow.Step{
				Name: "Run Tests",
				Run:  "go test -v -race -coverprofile=coverage.txt ./...",
				Env: workflow.Env{
					"GOPROXY": "https://proxy.golang.org",
				},
			},
			workflow.Step{
				Name: "Upload Coverage",
				Uses: "codecov/codecov-action@v3",
				With: workflow.With{
					"file": "coverage.txt",
				},
			},
		},
	}
}

// createBenchmarkStep creates a realistic step for benchmarking.
func createBenchmarkStep() workflow.Step {
	return workflow.Step{
		ID:   "build",
		Name: "Build Binary",
		If:   workflow.Success(),
		Uses: "actions/setup-go@v5",
		With: workflow.With{
			"go-version":   "1.23",
			"cache":        true,
			"cache-dependency-path": "go.sum",
		},
		Env: workflow.Env{
			"CGO_ENABLED": "0",
			"GOOS":        "linux",
			"GOARCH":      "amd64",
		},
		TimeoutMinutes:   10,
		WorkingDirectory: "./cmd/app",
	}
}

// BenchmarkSerializeWorkflow benchmarks workflow to YAML serialization.
func BenchmarkSerializeWorkflow(b *testing.B) {
	w := createBenchmarkWorkflow()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		data, err := serialize.ToYAML(w)
		if err != nil {
			b.Fatal(err)
		}
		if len(data) == 0 {
			b.Fatal("expected output")
		}
	}
}

// BenchmarkSerializeJob benchmarks job serialization.
func BenchmarkSerializeJob(b *testing.B) {
	// Create a workflow with just one job to benchmark job serialization
	w := &workflow.Workflow{
		Name: "Test",
		On: workflow.Triggers{
			Push: &workflow.PushTrigger{},
		},
		Jobs: map[string]workflow.Job{
			"test": *createBenchmarkJob(),
		},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		data, err := serialize.ToYAML(w)
		if err != nil {
			b.Fatal(err)
		}
		if len(data) == 0 {
			b.Fatal("expected output")
		}
	}
}

// BenchmarkSerializeStep benchmarks step serialization.
func BenchmarkSerializeStep(b *testing.B) {
	step := createBenchmarkStep()
	w := &workflow.Workflow{
		Name: "Test",
		On: workflow.Triggers{
			Push: &workflow.PushTrigger{},
		},
		Jobs: map[string]workflow.Job{
			"build": {
				RunsOn: "ubuntu-latest",
				Steps:  []any{step},
			},
		},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		data, err := serialize.ToYAML(w)
		if err != nil {
			b.Fatal(err)
		}
		if len(data) == 0 {
			b.Fatal("expected output")
		}
	}
}

// BenchmarkSerializeComplexWorkflow benchmarks a complex workflow with many jobs.
func BenchmarkSerializeComplexWorkflow(b *testing.B) {
	jobs := make(map[string]workflow.Job)
	for i := 0; i < 10; i++ {
		jobName := "job" + string(rune('a'+i))
		jobs[jobName] = workflow.Job{
			Name:   "Job " + string(rune('A'+i)),
			RunsOn: "ubuntu-latest",
			Steps: []any{
				workflow.Step{Uses: "actions/checkout@v4"},
				workflow.Step{Run: "echo 'Running job " + jobName + "'"},
				workflow.Step{
					Name: "Build",
					Run:  "go build ./...",
					Env: workflow.Env{
						"GOOS":   "linux",
						"GOARCH": "amd64",
					},
				},
			},
			Env: workflow.Env{
				"JOB_NAME": jobName,
			},
			TimeoutMinutes: 30,
		}
	}

	w := &workflow.Workflow{
		Name: "Complex CI",
		On: workflow.Triggers{
			Push: &workflow.PushTrigger{
				Branches: []string{"main", "develop", "release/*"},
			},
			PullRequest: &workflow.PullRequestTrigger{
				Branches: []string{"main"},
			},
			Schedule: []workflow.ScheduleTrigger{
				{Cron: "0 0 * * *"},
			},
		},
		Jobs: jobs,
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		data, err := serialize.ToYAML(w)
		if err != nil {
			b.Fatal(err)
		}
		if len(data) == 0 {
			b.Fatal("expected output")
		}
	}
}

// BenchmarkSerializeMatrixStrategy benchmarks serialization of matrix strategies.
func BenchmarkSerializeMatrixStrategy(b *testing.B) {
	w := &workflow.Workflow{
		Name: "Matrix Test",
		On: workflow.Triggers{
			Push: &workflow.PushTrigger{},
		},
		Jobs: map[string]workflow.Job{
			"test": {
				RunsOn: workflow.MatrixContext.Get("os"),
				Strategy: &workflow.Strategy{
					Matrix: &workflow.Matrix{
						Values: map[string][]any{
							"go":   {"1.21", "1.22", "1.23"},
							"os":   {"ubuntu-latest", "macos-latest", "windows-latest"},
							"arch": {"amd64", "arm64"},
						},
						Include: []map[string]any{
							{"go": "1.23", "os": "ubuntu-latest", "coverage": true},
						},
						Exclude: []map[string]any{
							{"os": "windows-latest", "arch": "arm64"},
						},
					},
					FailFast:    workflow.Ptr(false),
					MaxParallel: 4,
				},
				Steps: []any{
					workflow.Step{Uses: "actions/checkout@v4"},
					workflow.Step{Run: "go test ./..."},
				},
			},
		},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		data, err := serialize.ToYAML(w)
		if err != nil {
			b.Fatal(err)
		}
		if len(data) == 0 {
			b.Fatal("expected output")
		}
	}
}

// BenchmarkSerializeMinimalWorkflow benchmarks minimal workflow serialization.
func BenchmarkSerializeMinimalWorkflow(b *testing.B) {
	w := &workflow.Workflow{
		Name: "CI",
		On: workflow.Triggers{
			Push: &workflow.PushTrigger{},
		},
		Jobs: map[string]workflow.Job{
			"build": {
				RunsOn: "ubuntu-latest",
				Steps: []any{
					workflow.Step{Run: "echo hello"},
				},
			},
		},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		data, err := serialize.ToYAML(w)
		if err != nil {
			b.Fatal(err)
		}
		if len(data) == 0 {
			b.Fatal("expected output")
		}
	}
}
