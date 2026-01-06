package serialize_test

import (
	"strings"
	"testing"

	"github.com/lex00/wetwire-github-go/internal/serialize"
	"github.com/lex00/wetwire-github-go/workflow"
)

func TestBasicWorkflow(t *testing.T) {
	w := &workflow.Workflow{
		Name: "CI",
		On: workflow.Triggers{
			Push: &workflow.PushTrigger{
				Branches: []string{"main"},
			},
		},
	}

	yaml, err := serialize.ToYAML(w)
	if err != nil {
		t.Fatalf("ToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	// Check required fields
	if !strings.Contains(yamlStr, "name: CI") {
		t.Errorf("expected 'name: CI', got:\n%s", yamlStr)
	}
	// "on" may be quoted since it's a YAML reserved word
	if !strings.Contains(yamlStr, "on:") && !strings.Contains(yamlStr, `"on":`) {
		t.Errorf("expected 'on:' or '\"on\":', got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "push:") {
		t.Errorf("expected 'push:', got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "branches:") {
		t.Errorf("expected 'branches:', got:\n%s", yamlStr)
	}
}

func TestWorkflowWithJobs(t *testing.T) {
	w := &workflow.Workflow{
		Name: "CI",
		On: workflow.Triggers{
			Push: &workflow.PushTrigger{
				Branches: []string{"main"},
			},
		},
		Jobs: map[string]workflow.Job{
			"build": {
				Name:   "Build",
				RunsOn: "ubuntu-latest",
				Steps: []workflow.Step{
					{Uses: "actions/checkout@v4"},
					{Run: "go build ./..."},
				},
			},
		},
	}

	yaml, err := serialize.ToYAML(w)
	if err != nil {
		t.Fatalf("ToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	if !strings.Contains(yamlStr, "jobs:") {
		t.Errorf("expected 'jobs:', got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "build:") {
		t.Errorf("expected 'build:', got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "runs-on: ubuntu-latest") {
		t.Errorf("expected 'runs-on: ubuntu-latest', got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "uses: actions/checkout@v4") {
		t.Errorf("expected checkout step, got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "run: go build ./...") {
		t.Errorf("expected build step, got:\n%s", yamlStr)
	}
}

func TestStepWithAction(t *testing.T) {
	w := &workflow.Workflow{
		Name: "CI",
		On: workflow.Triggers{
			Push: &workflow.PushTrigger{},
		},
		Jobs: map[string]workflow.Job{
			"build": {
				RunsOn: "ubuntu-latest",
				Steps: []workflow.Step{
					{
						Uses: "actions/checkout@v4",
						With: workflow.With{
							"fetch-depth": 0,
							"submodules":  "recursive",
						},
					},
				},
			},
		},
	}

	yaml, err := serialize.ToYAML(w)
	if err != nil {
		t.Fatalf("ToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	if !strings.Contains(yamlStr, "with:") {
		t.Errorf("expected 'with:', got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "fetch-depth:") {
		t.Errorf("expected 'fetch-depth:', got:\n%s", yamlStr)
	}
}

func TestMatrixStrategy(t *testing.T) {
	w := &workflow.Workflow{
		Name: "CI",
		On: workflow.Triggers{
			Push: &workflow.PushTrigger{},
		},
		Jobs: map[string]workflow.Job{
			"test": {
				RunsOn: "ubuntu-latest",
				Strategy: &workflow.Strategy{
					Matrix: &workflow.Matrix{
						Values: map[string][]any{
							"go": {"1.22", "1.23"},
						},
					},
					FailFast: workflow.Ptr(false),
				},
				Steps: []workflow.Step{
					{Run: "go test ./..."},
				},
			},
		},
	}

	yaml, err := serialize.ToYAML(w)
	if err != nil {
		t.Fatalf("ToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	if !strings.Contains(yamlStr, "strategy:") {
		t.Errorf("expected 'strategy:', got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "matrix:") {
		t.Errorf("expected 'matrix:', got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "fail-fast: false") {
		t.Errorf("expected 'fail-fast: false', got:\n%s", yamlStr)
	}
}

func TestExpressionInEnv(t *testing.T) {
	w := &workflow.Workflow{
		Name: "CI",
		On: workflow.Triggers{
			Push: &workflow.PushTrigger{},
		},
		Jobs: map[string]workflow.Job{
			"deploy": {
				RunsOn: "ubuntu-latest",
				Steps: []workflow.Step{
					{
						Run: "deploy.sh",
						Env: workflow.Env{
							"TOKEN": workflow.Secrets.Get("DEPLOY_TOKEN"),
						},
					},
				},
			},
		},
	}

	yaml, err := serialize.ToYAML(w)
	if err != nil {
		t.Fatalf("ToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	if !strings.Contains(yamlStr, "${{ secrets.DEPLOY_TOKEN }}") {
		t.Errorf("expected secret expression, got:\n%s", yamlStr)
	}
}

func TestConditionInStep(t *testing.T) {
	w := &workflow.Workflow{
		Name: "CI",
		On: workflow.Triggers{
			Push: &workflow.PushTrigger{},
		},
		Jobs: map[string]workflow.Job{
			"deploy": {
				RunsOn: "ubuntu-latest",
				Steps: []workflow.Step{
					{
						If:  workflow.Branch("main"),
						Run: "deploy.sh",
					},
				},
			},
		},
	}

	yaml, err := serialize.ToYAML(w)
	if err != nil {
		t.Fatalf("ToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	if !strings.Contains(yamlStr, "if:") {
		t.Errorf("expected 'if:', got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "github.ref == 'refs/heads/main'") {
		t.Errorf("expected branch condition, got:\n%s", yamlStr)
	}
}

func TestScheduleTrigger(t *testing.T) {
	w := &workflow.Workflow{
		Name: "Nightly",
		On: workflow.Triggers{
			Schedule: []workflow.ScheduleTrigger{
				{Cron: "0 0 * * *"},
			},
		},
	}

	yaml, err := serialize.ToYAML(w)
	if err != nil {
		t.Fatalf("ToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	if !strings.Contains(yamlStr, "schedule:") {
		t.Errorf("expected 'schedule:', got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "cron:") {
		t.Errorf("expected 'cron:', got:\n%s", yamlStr)
	}
}

func TestWorkflowDispatch(t *testing.T) {
	w := &workflow.Workflow{
		Name: "Manual",
		On: workflow.Triggers{
			WorkflowDispatch: &workflow.WorkflowDispatchTrigger{
				Inputs: map[string]workflow.WorkflowInput{
					"environment": {
						Type:        "string",
						Required:    true,
						Description: "Deployment environment",
					},
				},
			},
		},
	}

	yaml, err := serialize.ToYAML(w)
	if err != nil {
		t.Fatalf("ToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	if !strings.Contains(yamlStr, "workflow_dispatch:") {
		t.Errorf("expected 'workflow_dispatch:', got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "inputs:") {
		t.Errorf("expected 'inputs:', got:\n%s", yamlStr)
	}
}

func TestPermissions(t *testing.T) {
	w := &workflow.Workflow{
		Name: "Release",
		On: workflow.Triggers{
			Push: &workflow.PushTrigger{},
		},
		Permissions: &workflow.Permissions{
			Contents:     "write",
			PullRequests: "read",
		},
	}

	yaml, err := serialize.ToYAML(w)
	if err != nil {
		t.Fatalf("ToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	if !strings.Contains(yamlStr, "permissions:") {
		t.Errorf("expected 'permissions:', got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "contents: write") {
		t.Errorf("expected 'contents: write', got:\n%s", yamlStr)
	}
}

func TestConcurrency(t *testing.T) {
	w := &workflow.Workflow{
		Name: "CI",
		On: workflow.Triggers{
			Push: &workflow.PushTrigger{},
		},
		Concurrency: &workflow.Concurrency{
			Group:            "ci-${{ github.ref }}",
			CancelInProgress: true,
		},
	}

	yaml, err := serialize.ToYAML(w)
	if err != nil {
		t.Fatalf("ToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	if !strings.Contains(yamlStr, "concurrency:") {
		t.Errorf("expected 'concurrency:', got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "cancel-in-progress: true") {
		t.Errorf("expected 'cancel-in-progress: true', got:\n%s", yamlStr)
	}
}

func TestServiceContainer(t *testing.T) {
	w := &workflow.Workflow{
		Name: "Test",
		On: workflow.Triggers{
			Push: &workflow.PushTrigger{},
		},
		Jobs: map[string]workflow.Job{
			"test": {
				RunsOn: "ubuntu-latest",
				Services: map[string]workflow.Service{
					"postgres": {
						Image: "postgres:15",
						Env: workflow.Env{
							"POSTGRES_PASSWORD": "test",
						},
						Ports: []any{"5432:5432"},
					},
				},
				Steps: []workflow.Step{
					{Run: "go test ./..."},
				},
			},
		},
	}

	yaml, err := serialize.ToYAML(w)
	if err != nil {
		t.Fatalf("ToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	if !strings.Contains(yamlStr, "services:") {
		t.Errorf("expected 'services:', got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "postgres:") {
		t.Errorf("expected 'postgres:', got:\n%s", yamlStr)
	}
}
