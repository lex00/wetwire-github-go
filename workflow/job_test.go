package workflow_test

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestJob(t *testing.T) {
	t.Run("basic job", func(t *testing.T) {
		job := workflow.Job{
			Name:   "test",
			RunsOn: "ubuntu-latest",
			Steps: []any{
				workflow.Step{Run: "echo 'hello'"},
			},
		}

		if job.Name != "test" {
			t.Errorf("expected Name='test', got %q", job.Name)
		}

		if job.RunsOn != "ubuntu-latest" {
			t.Errorf("expected RunsOn='ubuntu-latest', got %q", job.RunsOn)
		}

		if len(job.Steps) != 1 {
			t.Errorf("expected 1 step, got %d", len(job.Steps))
		}
	})

	t.Run("job with matrix runner", func(t *testing.T) {
		job := workflow.Job{
			Name:   "test",
			RunsOn: workflow.MatrixContext.Get("os"),
			Steps: []any{
				workflow.Step{Run: "go test"},
			},
		}

		expr, ok := job.RunsOn.(workflow.Expression)
		if !ok {
			t.Fatal("expected RunsOn to be an Expression")
		}

		if expr.Raw() != "matrix.os" {
			t.Errorf("expected 'matrix.os', got %q", expr.Raw())
		}
	})

	t.Run("job with needs", func(t *testing.T) {
		buildJob := workflow.Job{Name: "build", RunsOn: "ubuntu-latest"}
		testJob := workflow.Job{Name: "test", RunsOn: "ubuntu-latest"}

		deployJob := workflow.Job{
			Name:   "deploy",
			RunsOn: "ubuntu-latest",
			Needs:  []any{buildJob, testJob},
		}

		if len(deployJob.Needs) != 2 {
			t.Errorf("expected 2 needs, got %d", len(deployJob.Needs))
		}
	})

	t.Run("job with condition", func(t *testing.T) {
		job := workflow.Job{
			Name:   "deploy",
			RunsOn: "ubuntu-latest",
			If:     workflow.Branch("main").And(workflow.Success()),
		}

		expr, ok := job.If.(workflow.Expression)
		if !ok {
			t.Fatal("expected If to be an Expression")
		}

		expected := "(github.ref == 'refs/heads/main') && (success())"
		if expr.Raw() != expected {
			t.Errorf("expected %q, got %q", expected, expr.Raw())
		}
	})

	t.Run("job with outputs", func(t *testing.T) {
		job := workflow.Job{
			Name:   "build",
			RunsOn: "ubuntu-latest",
			Outputs: map[string]any{
				"version": workflow.Steps.Get("version", "value"),
				"tag":     "v1.0.0",
			},
		}

		if len(job.Outputs) != 2 {
			t.Errorf("expected 2 outputs, got %d", len(job.Outputs))
		}

		versionExpr, ok := job.Outputs["version"].(workflow.Expression)
		if !ok {
			t.Fatal("expected version to be an Expression")
		}

		if versionExpr.Raw() != "steps.version.outputs.value" {
			t.Errorf("expected 'steps.version.outputs.value', got %q", versionExpr.Raw())
		}
	})

	t.Run("job with timeout", func(t *testing.T) {
		job := workflow.Job{
			Name:           "long-job",
			RunsOn:         "ubuntu-latest",
			TimeoutMinutes: 120,
		}

		if job.TimeoutMinutes != 120 {
			t.Errorf("expected TimeoutMinutes=120, got %d", job.TimeoutMinutes)
		}
	})

	t.Run("job with continue on error", func(t *testing.T) {
		job := workflow.Job{
			Name:            "experimental",
			RunsOn:          "ubuntu-latest",
			ContinueOnError: true,
		}

		if !job.ContinueOnError {
			t.Error("expected ContinueOnError=true")
		}
	})
}

func TestPermissions(t *testing.T) {
	t.Run("all permissions", func(t *testing.T) {
		perms := workflow.Permissions{
			Actions:            workflow.PermissionRead,
			Checks:             workflow.PermissionWrite,
			Contents:           workflow.PermissionWrite,
			Deployments:        workflow.PermissionRead,
			Discussions:        workflow.PermissionRead,
			IDToken:            workflow.PermissionWrite,
			Issues:             workflow.PermissionWrite,
			Packages:           workflow.PermissionWrite,
			Pages:              workflow.PermissionWrite,
			PullRequests:       workflow.PermissionWrite,
			RepositoryProjects: workflow.PermissionRead,
			SecurityEvents:     workflow.PermissionWrite,
			Statuses:           workflow.PermissionWrite,
		}

		if perms.Actions != "read" {
			t.Errorf("expected Actions='read', got %q", perms.Actions)
		}

		if perms.Contents != "write" {
			t.Errorf("expected Contents='write', got %q", perms.Contents)
		}

		if perms.IDToken != "write" {
			t.Errorf("expected IDToken='write', got %q", perms.IDToken)
		}
	})

	t.Run("permission constants", func(t *testing.T) {
		if workflow.PermissionRead != "read" {
			t.Errorf("expected PermissionRead='read', got %q", workflow.PermissionRead)
		}

		if workflow.PermissionWrite != "write" {
			t.Errorf("expected PermissionWrite='write', got %q", workflow.PermissionWrite)
		}

		if workflow.PermissionNone != "none" {
			t.Errorf("expected PermissionNone='none', got %q", workflow.PermissionNone)
		}
	})
}

func TestEnvironment(t *testing.T) {
	env := workflow.Environment{
		Name: "production",
		URL:  "https://example.com",
	}

	if env.Name != "production" {
		t.Errorf("expected Name='production', got %q", env.Name)
	}

	if env.URL != "https://example.com" {
		t.Errorf("expected URL='https://example.com', got %q", env.URL)
	}
}

func TestConcurrency(t *testing.T) {
	t.Run("with cancel in progress", func(t *testing.T) {
		conc := workflow.Concurrency{
			Group:            "ci-${{ github.ref }}",
			CancelInProgress: true,
		}

		if conc.Group != "ci-${{ github.ref }}" {
			t.Errorf("expected Group='ci-${{ github.ref }}', got %q", conc.Group)
		}

		if !conc.CancelInProgress {
			t.Error("expected CancelInProgress=true")
		}
	})

	t.Run("without cancel in progress", func(t *testing.T) {
		conc := workflow.Concurrency{
			Group:            "deploy",
			CancelInProgress: false,
		}

		if conc.CancelInProgress {
			t.Error("expected CancelInProgress=false")
		}
	})
}

func TestJobDefaults(t *testing.T) {
	job := workflow.Job{
		Name:   "build",
		RunsOn: "ubuntu-latest",
		Defaults: &workflow.JobDefaults{
			Run: &workflow.RunDefaults{
				Shell:            "bash",
				WorkingDirectory: "./src",
			},
		},
	}

	if job.Defaults == nil {
		t.Fatal("expected Defaults to be set")
	}

	if job.Defaults.Run == nil {
		t.Fatal("expected Run defaults to be set")
	}

	if job.Defaults.Run.Shell != "bash" {
		t.Errorf("expected Shell='bash', got %q", job.Defaults.Run.Shell)
	}

	if job.Defaults.Run.WorkingDirectory != "./src" {
		t.Errorf("expected WorkingDirectory='./src', got %q", job.Defaults.Run.WorkingDirectory)
	}
}

func TestContainer(t *testing.T) {
	t.Run("basic container", func(t *testing.T) {
		container := workflow.Container{
			Image: "golang:1.23",
		}

		if container.Image != "golang:1.23" {
			t.Errorf("expected Image='golang:1.23', got %q", container.Image)
		}
	})

	t.Run("container with env", func(t *testing.T) {
		container := workflow.Container{
			Image: "golang:1.23",
			Env: workflow.Env{
				"CGO_ENABLED": "0",
				"GOOS":        "linux",
			},
		}

		if container.Env["CGO_ENABLED"] != "0" {
			t.Errorf("expected CGO_ENABLED='0', got %v", container.Env["CGO_ENABLED"])
		}
	})

	t.Run("container with credentials", func(t *testing.T) {
		container := workflow.Container{
			Image: "private-registry/image:latest",
			Credentials: &workflow.Credentials{
				Username: "user",
				Password: "${{ secrets.REGISTRY_PASSWORD }}",
			},
		}

		if container.Credentials == nil {
			t.Fatal("expected Credentials to be set")
		}

		if container.Credentials.Username != "user" {
			t.Errorf("expected Username='user', got %q", container.Credentials.Username)
		}
	})

	t.Run("container with ports and volumes", func(t *testing.T) {
		container := workflow.Container{
			Image:   "postgres:15",
			Ports:   []any{5432, "8080:80"},
			Volumes: []string{"/data:/var/lib/postgresql/data"},
			Options: "--health-cmd pg_isready --health-interval 10s",
		}

		if len(container.Ports) != 2 {
			t.Errorf("expected 2 ports, got %d", len(container.Ports))
		}

		if len(container.Volumes) != 1 {
			t.Errorf("expected 1 volume, got %d", len(container.Volumes))
		}

		if container.Options == "" {
			t.Error("expected Options to be set")
		}
	})
}

func TestService(t *testing.T) {
	service := workflow.Service{
		Image: "postgres:15",
		Env: workflow.Env{
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
		Ports:   []any{"5432:5432"},
		Volumes: []string{"/tmp/postgres:/var/lib/postgresql/data"},
		Options: "--health-cmd pg_isready",
	}

	if service.Image != "postgres:15" {
		t.Errorf("expected Image='postgres:15', got %q", service.Image)
	}

	if service.Env["POSTGRES_PASSWORD"] != "test" {
		t.Errorf("expected POSTGRES_PASSWORD='test', got %v", service.Env["POSTGRES_PASSWORD"])
	}

	if len(service.Ports) != 1 {
		t.Errorf("expected 1 port, got %d", len(service.Ports))
	}

	if len(service.Volumes) != 1 {
		t.Errorf("expected 1 volume, got %d", len(service.Volumes))
	}
}

func TestStrategy(t *testing.T) {
	t.Run("basic strategy", func(t *testing.T) {
		strategy := workflow.Strategy{
			Matrix: &workflow.Matrix{
				Values: map[string][]any{
					"go": {"1.22", "1.23"},
				},
			},
		}

		if strategy.Matrix == nil {
			t.Fatal("expected Matrix to be set")
		}

		if len(strategy.Matrix.Values["go"]) != 2 {
			t.Errorf("expected 2 go versions, got %d", len(strategy.Matrix.Values["go"]))
		}
	})

	t.Run("strategy with fail-fast and max-parallel", func(t *testing.T) {
		failFast := false
		strategy := workflow.Strategy{
			Matrix: &workflow.Matrix{
				Values: map[string][]any{
					"os": {"ubuntu-latest", "macos-latest"},
				},
			},
			FailFast:    &failFast,
			MaxParallel: 2,
		}

		if strategy.FailFast == nil {
			t.Fatal("expected FailFast to be set")
		}

		if *strategy.FailFast {
			t.Error("expected FailFast=false")
		}

		if strategy.MaxParallel != 2 {
			t.Errorf("expected MaxParallel=2, got %d", strategy.MaxParallel)
		}
	})
}
