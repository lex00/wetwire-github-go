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
		Steps: []any{
			workflow.Step{Run: "go build ./..."},
			workflow.Step{Run: "go test ./..."},
		},
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
		Steps:  []any{workflow.Step{Run: "go build"}},
	}

	testJob := workflow.Job{
		Name:   "test",
		RunsOn: "ubuntu-latest",
		Steps:  []any{workflow.Step{Run: "go test"}},
	}

	// In actual usage, Needs would contain the job variables
	// but for testing we verify the structure accepts []any
	deployJob := workflow.Job{
		Name:   "deploy",
		RunsOn: "ubuntu-latest",
		Needs:  []any{buildJob, testJob},
		Steps:  []any{workflow.Step{Run: "deploy.sh"}},
	}

	if len(deployJob.Needs) != 2 {
		t.Errorf("expected 2 needs, got %d", len(deployJob.Needs))
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


func TestWorkflowWithEnv(t *testing.T) {
	w := workflow.Workflow{
		Name: "CI",
		Env: workflow.Env{
			"GO_VERSION": "1.23",
			"NODE_ENV":   "production",
		},
	}

	if len(w.Env) != 2 {
		t.Errorf("expected 2 env vars, got %d", len(w.Env))
	}

	if w.Env["GO_VERSION"] != "1.23" {
		t.Errorf("expected GO_VERSION='1.23', got %v", w.Env["GO_VERSION"])
	}
}

func TestWorkflowWithDefaults(t *testing.T) {
	w := workflow.Workflow{
		Name: "CI",
		Defaults: &workflow.WorkflowDefaults{
			Run: &workflow.RunDefaults{
				Shell:            "bash",
				WorkingDirectory: "./app",
			},
		},
	}

	if w.Defaults == nil {
		t.Fatal("expected Defaults to be set")
	}

	if w.Defaults.Run.Shell != "bash" {
		t.Errorf("expected Shell='bash', got %q", w.Defaults.Run.Shell)
	}
}

func TestWorkflowWithPermissions(t *testing.T) {
	w := workflow.Workflow{
		Name: "Release",
		Permissions: &workflow.Permissions{
			Contents:     workflow.PermissionWrite,
			PullRequests: workflow.PermissionRead,
		},
	}

	if w.Permissions == nil {
		t.Fatal("expected Permissions to be set")
	}

	if w.Permissions.Contents != "write" {
		t.Errorf("expected Contents='write', got %q", w.Permissions.Contents)
	}

	if w.Permissions.PullRequests != "read" {
		t.Errorf("expected PullRequests='read', got %q", w.Permissions.PullRequests)
	}
}

func TestWorkflowWithJobs(t *testing.T) {
	w := workflow.Workflow{
		Name: "CI",
		Jobs: map[string]workflow.Job{
			"build": {
				Name:   "Build",
				RunsOn: "ubuntu-latest",
			},
			"test": {
				Name:   "Test",
				RunsOn: "ubuntu-latest",
			},
		},
	}

	if len(w.Jobs) != 2 {
		t.Errorf("expected 2 jobs, got %d", len(w.Jobs))
	}

	buildJob, ok := w.Jobs["build"]
	if !ok {
		t.Fatal("expected build job to exist")
	}

	if buildJob.Name != "Build" {
		t.Errorf("expected Name='Build', got %q", buildJob.Name)
	}
}

func TestWorkflowResourceType(t *testing.T) {
	w := workflow.Workflow{Name: "Test"}

	if w.ResourceType() != "workflow" {
		t.Errorf("expected ResourceType='workflow', got %q", w.ResourceType())
	}
}

func TestComplexWorkflow(t *testing.T) {
	w := workflow.Workflow{
		Name: "Complete CI/CD",
		On: workflow.Triggers{
			Push: &workflow.PushTrigger{
				Branches: workflow.List("main"),
			},
			PullRequest: &workflow.PullRequestTrigger{
				Types: workflow.List("opened", "synchronize"),
			},
			Schedule: []workflow.ScheduleTrigger{
				{Cron: "0 0 * * 0"},
			},
		},
		Env: workflow.Env{
			"CI": "true",
		},
		Concurrency: &workflow.Concurrency{
			Group:            "ci-${{ github.ref }}",
			CancelInProgress: true,
		},
		Permissions: &workflow.Permissions{
			Contents:     workflow.PermissionRead,
			PullRequests: workflow.PermissionWrite,
		},
		Defaults: &workflow.WorkflowDefaults{
			Run: &workflow.RunDefaults{
				Shell: "bash",
			},
		},
		Jobs: map[string]workflow.Job{
			"build": {
				Name:   "Build",
				RunsOn: "ubuntu-latest",
				Steps: []any{
					workflow.Step{
						Uses: "actions/checkout@v4",
					},
					workflow.Step{
						Run: "make build",
					},
				},
			},
		},
	}

	// Verify all fields are set correctly
	if w.Name != "Complete CI/CD" {
		t.Errorf("expected Name='Complete CI/CD', got %q", w.Name)
	}

	if w.On.Push == nil {
		t.Fatal("expected Push trigger to be set")
	}

	if w.Env["CI"] != "true" {
		t.Errorf("expected CI='true', got %v", w.Env["CI"])
	}

	if w.Concurrency == nil || !w.Concurrency.CancelInProgress {
		t.Error("expected CancelInProgress=true")
	}

	if w.Permissions == nil || w.Permissions.Contents != "read" {
		t.Error("expected Permissions.Contents='read'")
	}

	if w.Defaults == nil || w.Defaults.Run.Shell != "bash" {
		t.Error("expected Defaults.Run.Shell='bash'")
	}

	if len(w.Jobs) != 1 {
		t.Errorf("expected 1 job, got %d", len(w.Jobs))
	}
}
