package discover

import (
	"os"
	"path/filepath"
	"testing"
)

// createBenchmarkDir creates a temporary directory with realistic workflow files for benchmarking.
func createBenchmarkDir(b *testing.B) string {
	b.Helper()
	tmpDir := b.TempDir()

	// Create a realistic Go workflow file with multiple workflows and jobs
	workflowContent := `package main

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

// CI is the main continuous integration workflow
var CI = workflow.Workflow{
	Name: "CI",
	Jobs: []any{Build, Test, Lint},
}

// Release is the release workflow
var Release = workflow.Workflow{
	Name: "Release",
	Jobs: []any{Build, Test, Publish, Deploy},
}

// Nightly is the nightly build workflow
var Nightly = workflow.Workflow{
	Name: "Nightly",
	Jobs: []any{Build, Test, Integration},
}

var Build = workflow.Job{
	Name:   "build",
	RunsOn: "ubuntu-latest",
}

var Test = workflow.Job{
	Name:   "test",
	RunsOn: "ubuntu-latest",
	Needs:  []any{Build},
}

var Lint = workflow.Job{
	Name:   "lint",
	RunsOn: "ubuntu-latest",
}

var Publish = workflow.Job{
	Name:   "publish",
	RunsOn: "ubuntu-latest",
	Needs:  []any{Build, Test},
}

var Deploy = workflow.Job{
	Name:   "deploy",
	RunsOn: "ubuntu-latest",
	Needs:  []any{Publish},
}

var Integration = workflow.Job{
	Name:   "integration",
	RunsOn: "ubuntu-latest",
	Needs:  []any{Build, Test},
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "workflows.go"), []byte(workflowContent), 0644); err != nil {
		b.Fatal(err)
	}

	// Create a second file with more jobs
	jobsContent := `package main

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var BuildLinux = workflow.Job{
	Name:   "build-linux",
	RunsOn: "ubuntu-latest",
}

var BuildMacOS = workflow.Job{
	Name:   "build-macos",
	RunsOn: "macos-latest",
}

var BuildWindows = workflow.Job{
	Name:   "build-windows",
	RunsOn: "windows-latest",
}

var TestLinux = workflow.Job{
	Name:   "test-linux",
	RunsOn: "ubuntu-latest",
	Needs:  []any{BuildLinux},
}

var TestMacOS = workflow.Job{
	Name:   "test-macos",
	RunsOn: "macos-latest",
	Needs:  []any{BuildMacOS},
}

var TestWindows = workflow.Job{
	Name:   "test-windows",
	RunsOn: "windows-latest",
	Needs:  []any{BuildWindows},
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "jobs.go"), []byte(jobsContent), 0644); err != nil {
		b.Fatal(err)
	}

	// Create a subdirectory with more workflows
	subDir := filepath.Join(tmpDir, "ci")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		b.Fatal(err)
	}

	ciContent := `package ci

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var DockerBuild = workflow.Workflow{
	Name: "Docker Build",
	Jobs: []any{DockerJob},
}

var DockerJob = workflow.Job{
	Name:   "docker",
	RunsOn: "ubuntu-latest",
}
`
	if err := os.WriteFile(filepath.Join(subDir, "docker.go"), []byte(ciContent), 0644); err != nil {
		b.Fatal(err)
	}

	return tmpDir
}

// BenchmarkDiscoverWorkflows benchmarks workflow discovery.
func BenchmarkDiscoverWorkflows(b *testing.B) {
	tmpDir := createBenchmarkDir(b)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d := NewDiscoverer()
		result, err := d.Discover(tmpDir)
		if err != nil {
			b.Fatal(err)
		}
		// Prevent compiler optimization
		if len(result.Workflows) == 0 {
			b.Fatal("expected workflows")
		}
	}
}

// BenchmarkDiscoverJobs benchmarks job discovery.
func BenchmarkDiscoverJobs(b *testing.B) {
	tmpDir := createBenchmarkDir(b)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d := NewDiscoverer()
		result, err := d.Discover(tmpDir)
		if err != nil {
			b.Fatal(err)
		}
		// Prevent compiler optimization
		if len(result.Jobs) == 0 {
			b.Fatal("expected jobs")
		}
	}
}

// BenchmarkDiscoverAll benchmarks complete discovery of workflows and jobs.
func BenchmarkDiscoverAll(b *testing.B) {
	tmpDir := createBenchmarkDir(b)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d := NewDiscoverer()
		result, err := d.Discover(tmpDir)
		if err != nil {
			b.Fatal(err)
		}
		// Prevent compiler optimization
		if len(result.Workflows) == 0 || len(result.Jobs) == 0 {
			b.Fatal("expected resources")
		}
	}
}

// BenchmarkDiscoverLargeProject benchmarks discovery on a larger simulated project.
func BenchmarkDiscoverLargeProject(b *testing.B) {
	tmpDir := b.TempDir()

	// Create 10 files with workflows and jobs
	for i := 0; i < 10; i++ {
		content := `package main

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var Workflow` + string(rune('A'+i)) + ` = workflow.Workflow{
	Name: "Workflow ` + string(rune('A'+i)) + `",
	Jobs: []any{Job` + string(rune('A'+i)) + `1, Job` + string(rune('A'+i)) + `2},
}

var Job` + string(rune('A'+i)) + `1 = workflow.Job{
	Name:   "job-1",
	RunsOn: "ubuntu-latest",
}

var Job` + string(rune('A'+i)) + `2 = workflow.Job{
	Name:   "job-2",
	RunsOn: "ubuntu-latest",
	Needs:  []any{Job` + string(rune('A'+i)) + `1},
}
`
		filename := filepath.Join(tmpDir, "workflow"+string(rune('a'+i))+".go")
		if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
			b.Fatal(err)
		}
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d := NewDiscoverer()
		result, err := d.Discover(tmpDir)
		if err != nil {
			b.Fatal(err)
		}
		if len(result.Workflows) == 0 {
			b.Fatal("expected workflows")
		}
	}
}

// BenchmarkDiscoverSingleFile benchmarks discovery on a single file.
func BenchmarkDiscoverSingleFile(b *testing.B) {
	tmpDir := b.TempDir()

	content := `package main

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var CI = workflow.Workflow{
	Name: "CI",
	Jobs: []any{Build},
}

var Build = workflow.Job{
	Name:   "build",
	RunsOn: "ubuntu-latest",
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "ci.go"), []byte(content), 0644); err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d := NewDiscoverer()
		result, err := d.Discover(tmpDir)
		if err != nil {
			b.Fatal(err)
		}
		if len(result.Workflows) == 0 {
			b.Fatal("expected workflows")
		}
	}
}
