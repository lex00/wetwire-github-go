package workflow_test

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestWorkflowDeclaration(t *testing.T) {
	// Test the "no parens" pattern - struct literal declarations
	w := workflow.Workflow{
		Name: "CI",
		On: workflow.Triggers{
			Push: &workflow.PushTrigger{
				Branches: workflow.List("main"),
			},
		},
	}

	if w.Name != "CI" {
		t.Errorf("expected Name=CI, got %s", w.Name)
	}
	if w.On.Push == nil {
		t.Error("expected Push trigger to be set")
	}
	if len(w.On.Push.Branches) != 1 || w.On.Push.Branches[0] != "main" {
		t.Errorf("expected Branches=[main], got %v", w.On.Push.Branches)
	}
}

func TestJobDeclaration(t *testing.T) {
	job := workflow.Job{
		Name:   "build",
		RunsOn: "ubuntu-latest",
		Steps: workflow.List(
			workflow.Step{Run: "go build ./..."},
			workflow.Step{Run: "go test ./..."},
		),
	}

	if job.Name != "build" {
		t.Errorf("expected Name=build, got %s", job.Name)
	}
	if job.RunsOn != "ubuntu-latest" {
		t.Errorf("expected RunsOn=ubuntu-latest, got %v", job.RunsOn)
	}
	if len(job.Steps) != 2 {
		t.Errorf("expected 2 steps, got %d", len(job.Steps))
	}
}

func TestStepWithAction(t *testing.T) {
	step := workflow.Step{
		Uses: "actions/checkout@v4",
		With: workflow.With{
			"fetch-depth": 0,
		},
	}

	if step.Uses != "actions/checkout@v4" {
		t.Errorf("expected Uses=actions/checkout@v4, got %s", step.Uses)
	}
	if step.With["fetch-depth"] != 0 {
		t.Errorf("expected fetch-depth=0, got %v", step.With["fetch-depth"])
	}
}

func TestMatrixDeclaration(t *testing.T) {
	matrix := workflow.Matrix{
		Values: map[string][]any{
			"go": {"1.22", "1.23"},
			"os": {"ubuntu-latest", "macos-latest"},
		},
	}

	if len(matrix.Values["go"]) != 2 {
		t.Errorf("expected 2 go versions, got %d", len(matrix.Values["go"]))
	}
	if len(matrix.Values["os"]) != 2 {
		t.Errorf("expected 2 os options, got %d", len(matrix.Values["os"]))
	}
}

func TestStrategyDeclaration(t *testing.T) {
	strategy := workflow.Strategy{
		Matrix: &workflow.Matrix{
			Values: map[string][]any{
				"go": {"1.22", "1.23"},
			},
		},
		FailFast:    workflow.Ptr(false),
		MaxParallel: 2,
	}

	if strategy.Matrix == nil {
		t.Error("expected Matrix to be set")
	}
	if *strategy.FailFast != false {
		t.Error("expected FailFast=false")
	}
	if strategy.MaxParallel != 2 {
		t.Errorf("expected MaxParallel=2, got %d", strategy.MaxParallel)
	}
}

func TestTriggerDeclaration(t *testing.T) {
	triggers := workflow.Triggers{
		Push: &workflow.PushTrigger{
			Branches: workflow.List("main", "develop"),
		},
		PullRequest: &workflow.PullRequestTrigger{
			Branches: workflow.List("main"),
			Types:    workflow.List("opened", "synchronize"),
		},
		Schedule: []workflow.ScheduleTrigger{
			{Cron: "0 0 * * *"},
		},
	}

	if len(triggers.Push.Branches) != 2 {
		t.Errorf("expected 2 push branches, got %d", len(triggers.Push.Branches))
	}
	if len(triggers.PullRequest.Types) != 2 {
		t.Errorf("expected 2 PR types, got %d", len(triggers.PullRequest.Types))
	}
	if len(triggers.Schedule) != 1 {
		t.Errorf("expected 1 schedule, got %d", len(triggers.Schedule))
	}
}

func TestWorkflowCallTrigger(t *testing.T) {
	trigger := workflow.WorkflowCallTrigger{
		Inputs: map[string]workflow.WorkflowInput{
			"environment": {
				Type:        "string",
				Required:    true,
				Description: "Deployment environment",
			},
		},
		Secrets: map[string]workflow.WorkflowSecret{
			"deploy-token": {
				Required:    true,
				Description: "Token for deployment",
			},
		},
		Outputs: map[string]workflow.WorkflowOutput{
			"artifact-url": {
				Value:       workflow.Steps.Get("upload", "url"),
				Description: "URL of the uploaded artifact",
			},
		},
	}

	if trigger.Inputs["environment"].Type != "string" {
		t.Error("expected input type=string")
	}
	if !trigger.Secrets["deploy-token"].Required {
		t.Error("expected secret to be required")
	}
}

func TestJobWithNeeds(t *testing.T) {
	// Simulate the pattern where jobs reference each other
	buildJob := workflow.Job{
		Name:   "build",
		RunsOn: "ubuntu-latest",
		Steps:  workflow.List(workflow.Step{Run: "go build"}),
	}

	testJob := workflow.Job{
		Name:   "test",
		RunsOn: "ubuntu-latest",
		Steps:  workflow.List(workflow.Step{Run: "go test"}),
	}

	// In actual usage, Needs would contain the job variables
	// but for testing we verify the structure accepts []any
	deployJob := workflow.Job{
		Name:   "deploy",
		RunsOn: "ubuntu-latest",
		Needs:  []any{buildJob, testJob},
		Steps:  workflow.List(workflow.Step{Run: "deploy.sh"}),
	}

	if len(deployJob.Needs) != 2 {
		t.Errorf("expected 2 needs, got %d", len(deployJob.Needs))
	}
}

func TestPermissions(t *testing.T) {
	job := workflow.Job{
		Name:   "release",
		RunsOn: "ubuntu-latest",
		Permissions: &workflow.Permissions{
			Contents:     workflow.PermissionWrite,
			PullRequests: workflow.PermissionRead,
		},
	}

	if job.Permissions.Contents != "write" {
		t.Errorf("expected contents=write, got %s", job.Permissions.Contents)
	}
}

func TestContainer(t *testing.T) {
	job := workflow.Job{
		Name:   "build",
		RunsOn: "ubuntu-latest",
		Container: &workflow.Container{
			Image: "golang:1.23",
			Env: workflow.Env{
				"CGO_ENABLED": "0",
			},
		},
	}

	if job.Container.Image != "golang:1.23" {
		t.Errorf("expected image=golang:1.23, got %s", job.Container.Image)
	}
}

func TestServices(t *testing.T) {
	job := workflow.Job{
		Name:   "test",
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
	}

	if job.Services["postgres"].Image != "postgres:15" {
		t.Error("expected postgres service image")
	}
}

func TestConcurrency(t *testing.T) {
	w := workflow.Workflow{
		Name: "CI",
		Concurrency: &workflow.Concurrency{
			Group:            "ci-${{ github.ref }}",
			CancelInProgress: true,
		},
	}

	if !w.Concurrency.CancelInProgress {
		t.Error("expected CancelInProgress=true")
	}
}

func TestEnvironment(t *testing.T) {
	job := workflow.Job{
		Name:   "deploy",
		RunsOn: "ubuntu-latest",
		Environment: &workflow.Environment{
			Name: "production",
			URL:  "https://example.com",
		},
	}

	if job.Environment.Name != "production" {
		t.Errorf("expected environment=production, got %s", job.Environment.Name)
	}
}
