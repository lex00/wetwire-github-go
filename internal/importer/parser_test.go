package importer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParser_Parse_Simple(t *testing.T) {
	yaml := []byte(`name: CI
on: push
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: echo "Hello"
`)

	workflow, err := ParseWorkflow(yaml)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if workflow.Name != "CI" {
		t.Errorf("workflow.Name = %q, want %q", workflow.Name, "CI")
	}

	if workflow.On.Push == nil {
		t.Error("workflow.On.Push is nil")
	}

	if len(workflow.Jobs) != 1 {
		t.Errorf("len(workflow.Jobs) = %d, want 1", len(workflow.Jobs))
	}

	job := workflow.Jobs["build"]
	if job == nil {
		t.Fatal("workflow.Jobs[build] is nil")
	}

	if job.GetRunsOn() != "ubuntu-latest" {
		t.Errorf("job.GetRunsOn() = %q, want %q", job.GetRunsOn(), "ubuntu-latest")
	}

	if len(job.Steps) != 2 {
		t.Errorf("len(job.Steps) = %d, want 2", len(job.Steps))
	}
}

func TestParser_Parse_MultipleTriggers(t *testing.T) {
	yaml := []byte(`name: CI
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - run: echo "test"
`)

	workflow, err := ParseWorkflow(yaml)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if workflow.On.Push == nil {
		t.Error("workflow.On.Push is nil")
	}

	if workflow.On.PullRequest == nil {
		t.Error("workflow.On.PullRequest is nil")
	}
}

func TestParser_Parse_DetailedTriggers(t *testing.T) {
	yaml := []byte(`name: CI
on:
  push:
    branches:
      - main
      - develop
  pull_request:
    branches:
      - main
    types:
      - opened
      - synchronize
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - run: echo "test"
`)

	workflow, err := ParseWorkflow(yaml)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if workflow.On.Push == nil {
		t.Fatal("workflow.On.Push is nil")
	}

	if len(workflow.On.Push.Branches) != 2 {
		t.Errorf("len(workflow.On.Push.Branches) = %d, want 2", len(workflow.On.Push.Branches))
	}

	if workflow.On.PullRequest == nil {
		t.Fatal("workflow.On.PullRequest is nil")
	}

	if len(workflow.On.PullRequest.Types) != 2 {
		t.Errorf("len(workflow.On.PullRequest.Types) = %d, want 2", len(workflow.On.PullRequest.Types))
	}
}

func TestParser_Parse_JobNeeds(t *testing.T) {
	yaml := []byte(`name: CI
on: push
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - run: echo "build"
  test:
    needs: build
    runs-on: ubuntu-latest
    steps:
      - run: echo "test"
  deploy:
    needs: [build, test]
    runs-on: ubuntu-latest
    steps:
      - run: echo "deploy"
`)

	workflow, err := ParseWorkflow(yaml)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	build := workflow.Jobs["build"]
	if len(build.GetNeeds()) != 0 {
		t.Errorf("build.GetNeeds() = %v, want empty", build.GetNeeds())
	}

	test := workflow.Jobs["test"]
	needs := test.GetNeeds()
	if len(needs) != 1 || needs[0] != "build" {
		t.Errorf("test.GetNeeds() = %v, want [build]", needs)
	}

	deploy := workflow.Jobs["deploy"]
	deployNeeds := deploy.GetNeeds()
	if len(deployNeeds) != 2 {
		t.Errorf("deploy.GetNeeds() = %v, want [build, test]", deployNeeds)
	}
}

func TestParser_Parse_Matrix(t *testing.T) {
	yaml := []byte(`name: CI
on: push
jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go: ["1.22", "1.23"]
        os: [ubuntu-latest, macos-latest]
    steps:
      - run: echo "test"
`)

	workflow, err := ParseWorkflow(yaml)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	job := workflow.Jobs["test"]
	if job == nil || job.Strategy == nil || job.Strategy.Matrix == nil {
		t.Fatal("job.Strategy.Matrix is nil")
	}

	if len(job.Strategy.Matrix.Values) != 2 {
		t.Errorf("len(Matrix.Values) = %d, want 2", len(job.Strategy.Matrix.Values))
	}
}

func TestParser_Parse_StepWithDetails(t *testing.T) {
	yaml := []byte(`name: CI
on: push
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - id: checkout
        name: Checkout code
        uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - name: Run tests
        run: go test ./...
        env:
          GO_VERSION: "1.23"
        if: success()
        working-directory: ./src
`)

	workflow, err := ParseWorkflow(yaml)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	steps := workflow.Jobs["test"].Steps
	if len(steps) != 2 {
		t.Fatalf("len(steps) = %d, want 2", len(steps))
	}

	// Check first step
	step1 := steps[0]
	if step1.ID != "checkout" {
		t.Errorf("step1.ID = %q, want %q", step1.ID, "checkout")
	}
	if step1.Uses != "actions/checkout@v4" {
		t.Errorf("step1.Uses = %q, want %q", step1.Uses, "actions/checkout@v4")
	}
	if step1.With["fetch-depth"] == nil {
		t.Error("step1.With[fetch-depth] is nil")
	}

	// Check second step
	step2 := steps[1]
	if step2.Run != "go test ./..." {
		t.Errorf("step2.Run = %q, want %q", step2.Run, "go test ./...")
	}
	if step2.If != "success()" {
		t.Errorf("step2.If = %q, want %q", step2.If, "success()")
	}
	if step2.WorkingDirectory != "./src" {
		t.Errorf("step2.WorkingDirectory = %q, want %q", step2.WorkingDirectory, "./src")
	}
}

func TestParser_ParseFile(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "ci.yml")

	yaml := []byte(`name: CI
on: push
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - run: echo "test"
`)

	if err := os.WriteFile(path, yaml, 0644); err != nil {
		t.Fatal(err)
	}

	workflow, err := ParseWorkflowFile(path)
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}

	if workflow.Name != "CI" {
		t.Errorf("workflow.Name = %q, want %q", workflow.Name, "CI")
	}
}

func TestBuildReferenceGraph(t *testing.T) {
	yaml := []byte(`name: CI
on: push
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - id: checkout
        uses: actions/checkout@v4
      - run: echo "build"
  test:
    needs: build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/setup-go@v5
      - run: go test ./...
`)

	workflow, err := ParseWorkflow(yaml)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	graph := BuildReferenceGraph(workflow)

	// Check job dependencies
	if len(graph.JobDependencies["build"]) != 0 {
		t.Error("build should have no dependencies")
	}
	testNeeds := graph.JobDependencies["test"]
	if len(testNeeds) != 1 || testNeeds[0] != "build" {
		t.Errorf("test dependencies = %v, want [build]", testNeeds)
	}

	// Check used actions
	if len(graph.UsedActions) != 2 {
		t.Errorf("len(UsedActions) = %d, want 2", len(graph.UsedActions))
	}

	// Check step outputs
	if _, ok := graph.StepOutputs["checkout"]; !ok {
		t.Error("StepOutputs missing 'checkout'")
	}
}

func TestNewParser(t *testing.T) {
	p := NewParser()
	if p == nil {
		t.Error("NewParser() returned nil")
	}
}

func TestParser_Parse_Invalid(t *testing.T) {
	_, err := ParseWorkflow([]byte("invalid: yaml: ["))
	if err == nil {
		t.Error("Parse() expected error for invalid YAML")
	}
}

func TestIRJob_GetRunsOn_Slice(t *testing.T) {
	job := &IRJob{
		RunsOn: []any{"ubuntu-latest", "macos-latest"},
	}
	if job.GetRunsOn() != "ubuntu-latest" {
		t.Errorf("GetRunsOn() = %q, want %q", job.GetRunsOn(), "ubuntu-latest")
	}
}

func TestIRJob_GetRunsOn_Nil(t *testing.T) {
	job := &IRJob{}
	if job.GetRunsOn() != "" {
		t.Errorf("GetRunsOn() = %q, want empty", job.GetRunsOn())
	}
}

func TestIRJob_GetNeeds_Nil(t *testing.T) {
	job := &IRJob{}
	if len(job.GetNeeds()) != 0 {
		t.Errorf("GetNeeds() = %v, want empty", job.GetNeeds())
	}
}
