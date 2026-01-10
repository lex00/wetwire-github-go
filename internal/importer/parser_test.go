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

func TestParser_Parse_WorkflowDispatch(t *testing.T) {
	yaml := []byte(`name: Manual
on: workflow_dispatch
jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - run: echo "manual"
`)

	workflow, err := ParseWorkflow(yaml)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if workflow.On.WorkflowDispatch == nil {
		t.Error("workflow.On.WorkflowDispatch is nil")
	}
}

func TestParser_Parse_WorkflowCall(t *testing.T) {
	yaml := []byte(`name: Reusable
on: workflow_call
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - run: echo "reusable"
`)

	workflow, err := ParseWorkflow(yaml)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if workflow.On.WorkflowCall == nil {
		t.Error("workflow.On.WorkflowCall is nil")
	}
}

func TestParser_Parse_Schedule(t *testing.T) {
	yaml := []byte(`name: Scheduled
on:
  schedule:
    - cron: "0 0 * * *"
    - cron: "0 12 * * *"
jobs:
  backup:
    runs-on: ubuntu-latest
    steps:
      - run: echo "scheduled"
`)

	workflow, err := ParseWorkflow(yaml)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(workflow.On.Schedule) != 2 {
		t.Errorf("len(workflow.On.Schedule) = %d, want 2", len(workflow.On.Schedule))
	}
	if workflow.On.Schedule[0].Cron != "0 0 * * *" {
		t.Errorf("Schedule[0].Cron = %q, want %q", workflow.On.Schedule[0].Cron, "0 0 * * *")
	}
}

func TestParser_Parse_RepositoryDispatch(t *testing.T) {
	yaml := []byte(`name: Dispatch
on: repository_dispatch
jobs:
  handle:
    runs-on: ubuntu-latest
    steps:
      - run: echo "dispatched"
`)

	workflow, err := ParseWorkflow(yaml)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if workflow.On.RepositoryDispatch == nil {
		t.Error("workflow.On.RepositoryDispatch is nil")
	}
}

func TestParser_Parse_Release(t *testing.T) {
	yaml := []byte(`name: Release
on: release
jobs:
  publish:
    runs-on: ubuntu-latest
    steps:
      - run: echo "release"
`)

	workflow, err := ParseWorkflow(yaml)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if workflow.On.Release == nil {
		t.Error("workflow.On.Release is nil")
	}
}

func TestParser_Parse_Issues(t *testing.T) {
	yaml := []byte(`name: Issues
on: issues
jobs:
  triage:
    runs-on: ubuntu-latest
    steps:
      - run: echo "issue"
`)

	workflow, err := ParseWorkflow(yaml)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if workflow.On.Issues == nil {
		t.Error("workflow.On.Issues is nil")
	}
}

func TestParser_Parse_PullRequestTarget(t *testing.T) {
	yaml := []byte(`name: PR Target
on: pull_request_target
jobs:
  check:
    runs-on: ubuntu-latest
    steps:
      - run: echo "pr target"
`)

	workflow, err := ParseWorkflow(yaml)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if workflow.On.PullRequestTarget == nil {
		t.Error("workflow.On.PullRequestTarget is nil")
	}
}

func TestParser_Parse_FileNotFound(t *testing.T) {
	_, err := ParseWorkflowFile("/nonexistent/path/to/file.yml")
	if err == nil {
		t.Error("ParseFile() expected error for nonexistent file")
	}
}

func TestParser_ParseFile_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "invalid.yml")

	invalidYAML := []byte("invalid: yaml: [[[")
	if err := os.WriteFile(path, invalidYAML, 0644); err != nil {
		t.Fatal(err)
	}

	_, err := ParseWorkflowFile(path)
	if err == nil {
		t.Error("ParseFile() expected error for invalid YAML")
	}
}

func TestParser_Parse_ComplexMatrix(t *testing.T) {
	yaml := []byte(`name: Matrix
on: push
jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        go: ["1.22", "1.23"]
        include:
          - os: ubuntu-latest
            go: "1.21"
        exclude:
          - os: windows-latest
            go: "1.22"
    steps:
      - run: echo "test"
`)

	workflow, err := ParseWorkflow(yaml)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	matrix := workflow.Jobs["test"].Strategy.Matrix
	if matrix == nil {
		t.Fatal("Matrix is nil")
	}

	if len(matrix.Values["os"]) != 3 {
		t.Errorf("Matrix.Values[os] length = %d, want 3", len(matrix.Values["os"]))
	}

	if len(matrix.Include) != 1 {
		t.Errorf("Matrix.Include length = %d, want 1", len(matrix.Include))
	}

	if len(matrix.Exclude) != 1 {
		t.Errorf("Matrix.Exclude length = %d, want 1", len(matrix.Exclude))
	}
}

func TestParser_Parse_JobWithContainer(t *testing.T) {
	yaml := []byte(`name: Container
on: push
jobs:
  test:
    runs-on: ubuntu-latest
    container:
      image: node:18
      credentials:
        username: user
        password: pass
      env:
        NODE_ENV: test
      ports:
        - 80
      volumes:
        - /data
      options: --cpus 1
    steps:
      - run: echo "containerized"
`)

	workflow, err := ParseWorkflow(yaml)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	container := workflow.Jobs["test"].Container
	if container == nil {
		t.Fatal("Container is nil")
	}

	if container.Image != "node:18" {
		t.Errorf("Container.Image = %q, want %q", container.Image, "node:18")
	}

	if container.Credentials == nil || container.Credentials.Username != "user" {
		t.Error("Container credentials not parsed correctly")
	}

	if len(container.Env) == 0 {
		t.Error("Container env is empty")
	}

	if container.Options != "--cpus 1" {
		t.Errorf("Container.Options = %q, want %q", container.Options, "--cpus 1")
	}
}

func TestParser_Parse_JobWithServices(t *testing.T) {
	yaml := []byte(`name: Services
on: push
jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:14
        env:
          POSTGRES_PASSWORD: postgres
        ports:
          - 5432:5432
        options: >-
          --health-cmd pg_isready
    steps:
      - run: echo "services"
`)

	workflow, err := ParseWorkflow(yaml)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	services := workflow.Jobs["test"].Services
	if services == nil || services["postgres"] == nil {
		t.Fatal("Services not parsed correctly")
	}

	postgres := services["postgres"]
	if postgres.Image != "postgres:14" {
		t.Errorf("Service.Image = %q, want %q", postgres.Image, "postgres:14")
	}
}

func TestParser_Parse_JobWithOutputs(t *testing.T) {
	yaml := []byte(`name: Outputs
on: push
jobs:
  build:
    runs-on: ubuntu-latest
    outputs:
      version: ${{ steps.get-version.outputs.version }}
      artifact: ${{ steps.build.outputs.path }}
    steps:
      - id: get-version
        run: echo "version=1.0.0" >> $GITHUB_OUTPUT
`)

	workflow, err := ParseWorkflow(yaml)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	outputs := workflow.Jobs["build"].Outputs
	if len(outputs) != 2 {
		t.Errorf("len(outputs) = %d, want 2", len(outputs))
	}

	if _, ok := outputs["version"]; !ok {
		t.Error("Missing version output")
	}
}

func TestParser_Parse_ReusableWorkflowJob(t *testing.T) {
	yaml := []byte(`name: Caller
on: push
jobs:
  call-workflow:
    uses: ./.github/workflows/reusable.yml
    with:
      config: production
    secrets: inherit
`)

	workflow, err := ParseWorkflow(yaml)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	job := workflow.Jobs["call-workflow"]
	if job.Uses != "./.github/workflows/reusable.yml" {
		t.Errorf("Job.Uses = %q, want %q", job.Uses, "./.github/workflows/reusable.yml")
	}

	if job.With == nil || job.With["config"] != "production" {
		t.Error("Job with parameters not parsed correctly")
	}

	if job.Secrets != "inherit" {
		t.Errorf("Job.Secrets = %v, want %q", job.Secrets, "inherit")
	}
}

func TestParser_Parse_WorkflowPermissions(t *testing.T) {
	yaml := []byte(`name: Permissions
on: push
permissions:
  contents: read
  issues: write
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

	if workflow.Permissions["contents"] != "read" {
		t.Error("Workflow permissions not parsed correctly")
	}
}

func TestParser_Parse_WorkflowDefaults(t *testing.T) {
	yaml := []byte(`name: Defaults
on: push
defaults:
  run:
    shell: bash
    working-directory: ./src
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

	if workflow.Defaults == nil || workflow.Defaults.Run == nil {
		t.Fatal("Defaults not parsed")
	}

	if workflow.Defaults.Run.Shell != "bash" {
		t.Errorf("Defaults.Run.Shell = %q, want %q", workflow.Defaults.Run.Shell, "bash")
	}

	if workflow.Defaults.Run.WorkingDirectory != "./src" {
		t.Errorf("Defaults.Run.WorkingDirectory = %q, want %q", workflow.Defaults.Run.WorkingDirectory, "./src")
	}
}

func TestParser_Parse_WorkflowConcurrency(t *testing.T) {
	yaml := []byte(`name: Concurrency
on: push
concurrency:
  group: ${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true
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

	if workflow.Concurrency == nil {
		t.Fatal("Concurrency not parsed")
	}

	if !workflow.Concurrency.CancelInProgress {
		t.Error("Concurrency.CancelInProgress should be true")
	}
}

func TestIRJob_GetRunsOn_StringSlice(t *testing.T) {
	job := &IRJob{
		RunsOn: []string{"ubuntu-latest", "macos-latest"},
	}
	if job.GetRunsOn() != "ubuntu-latest" {
		t.Errorf("GetRunsOn() = %q, want %q", job.GetRunsOn(), "ubuntu-latest")
	}
}

func TestIRJob_GetNeeds_StringSlice(t *testing.T) {
	job := &IRJob{
		Needs: []string{"build", "test"},
	}
	needs := job.GetNeeds()
	if len(needs) != 2 || needs[0] != "build" || needs[1] != "test" {
		t.Errorf("GetNeeds() = %v, want [build test]", needs)
	}
}

func TestBuildReferenceGraph_ReusableWorkflow(t *testing.T) {
	workflow := &IRWorkflow{
		Jobs: map[string]*IRJob{
			"call-workflow": {
				Uses: "./.github/workflows/reusable.yml",
			},
		},
	}

	graph := BuildReferenceGraph(workflow)

	// Reusable workflow should be in used actions
	found := false
	for _, action := range graph.UsedActions {
		if action == "./.github/workflows/reusable.yml" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Reusable workflow not found in UsedActions")
	}
}
