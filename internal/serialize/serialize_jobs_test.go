package serialize_test

import (
	"strings"
	"testing"

	"github.com/lex00/wetwire-github-go/internal/serialize"
	"github.com/lex00/wetwire-github-go/workflow"
)

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
				Steps: []any{
					workflow.Step{Uses: "actions/checkout@v4"},
					workflow.Step{Run: "go build ./..."},
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
				Steps: []any{
					workflow.Step{
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
				Steps: []any{
					workflow.Step{Run: "go test ./..."},
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

func TestConditionInStep(t *testing.T) {
	w := &workflow.Workflow{
		Name: "CI",
		On: workflow.Triggers{
			Push: &workflow.PushTrigger{},
		},
		Jobs: map[string]workflow.Job{
			"deploy": {
				RunsOn: "ubuntu-latest",
				Steps: []any{
					workflow.Step{
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
				Steps: []any{
					workflow.Step{Run: "go test ./..."},
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

// TestStepActionInterface tests serialization of action wrappers implementing StepAction.
func TestStepActionInterface(t *testing.T) {
	// Create a mock StepAction implementation
	type mockAction struct {
		action string
		inputs map[string]any
	}

	mockAction1 := mockAction{
		action: "actions/checkout@v4",
		inputs: map[string]any{
			"fetch-depth": 0,
			"submodules":  "recursive",
		},
	}

	w := &workflow.Workflow{
		Name: "Test Action",
		On: workflow.Triggers{
			Push: &workflow.PushTrigger{},
		},
		Jobs: map[string]workflow.Job{
			"test": {
				RunsOn: "ubuntu-latest",
				Steps: []any{
					// Test that StepAction interface types are properly serialized
					workflow.Step{
						Uses: mockAction1.action,
						With: mockAction1.inputs,
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

	if !strings.Contains(yamlStr, "uses: actions/checkout@v4") {
		t.Errorf("expected action reference, got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "fetch-depth:") {
		t.Errorf("expected input field, got:\n%s", yamlStr)
	}
}

// TestNilAndEmptyValues tests handling of nil pointers and empty values.
func TestNilAndEmptyValues(t *testing.T) {
	w := &workflow.Workflow{
		Name: "Test",
		On: workflow.Triggers{
			Push: &workflow.PushTrigger{},
		},
		Jobs: map[string]workflow.Job{
			"test": {
				RunsOn: "ubuntu-latest",
				// Test nil pointers
				Permissions: nil,
				Environment: nil,
				Concurrency: nil,
				Strategy:    nil,
				Container:   nil,
				// Empty collections
				Needs:    []any{},
				Env:      workflow.Env{},
				Outputs:  workflow.Env{},
				Services: map[string]workflow.Service{},
				Steps: []any{
					workflow.Step{Run: "echo test"},
				},
			},
		},
	}

	yaml, err := serialize.ToYAML(w)
	if err != nil {
		t.Fatalf("ToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	// Verify empty/nil values are omitted
	if strings.Contains(yamlStr, "permissions:") {
		t.Errorf("nil permissions should be omitted, got:\n%s", yamlStr)
	}
	if strings.Contains(yamlStr, "environment:") {
		t.Errorf("nil environment should be omitted, got:\n%s", yamlStr)
	}
	if strings.Contains(yamlStr, "needs:") {
		t.Errorf("empty needs should be omitted, got:\n%s", yamlStr)
	}
}

// TestMatrixIncludeExclude tests matrix include and exclude functionality.
func TestMatrixIncludeExclude(t *testing.T) {
	w := &workflow.Workflow{
		Name: "Test",
		On: workflow.Triggers{
			Push: &workflow.PushTrigger{},
		},
		Jobs: map[string]workflow.Job{
			"test": {
				RunsOn: "ubuntu-latest",
				Strategy: &workflow.Strategy{
					Matrix: &workflow.Matrix{
						Values: map[string][]any{
							"os": {"ubuntu-latest", "macos-latest"},
						},
						Include: []map[string]any{
							{"os": "windows-latest", "experimental": true},
						},
						Exclude: []map[string]any{
							{"os": "macos-latest"},
						},
					},
				},
				Steps: []any{
					workflow.Step{Run: "echo test"},
				},
			},
		},
	}

	yaml, err := serialize.ToYAML(w)
	if err != nil {
		t.Fatalf("ToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	if !strings.Contains(yamlStr, "include:") {
		t.Errorf("expected 'include:', got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "exclude:") {
		t.Errorf("expected 'exclude:', got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "windows-latest") {
		t.Errorf("expected windows-latest in include, got:\n%s", yamlStr)
	}
}

// TestJobNeeds tests job dependency serialization.
func TestJobNeeds(t *testing.T) {
	buildJob := workflow.Job{Name: "build"}
	testJob := workflow.Job{Name: "test"}

	w := &workflow.Workflow{
		Name: "Test",
		On: workflow.Triggers{
			Push: &workflow.PushTrigger{},
		},
		Jobs: map[string]workflow.Job{
			"build": buildJob,
			"test":  testJob,
			"deploy": {
				Name:   "deploy",
				RunsOn: "ubuntu-latest",
				Needs:  []any{buildJob, testJob},
				Steps: []any{
					workflow.Step{Run: "echo deploy"},
				},
			},
		},
	}

	yaml, err := serialize.ToYAML(w)
	if err != nil {
		t.Fatalf("ToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	if !strings.Contains(yamlStr, "needs:") {
		t.Errorf("expected 'needs:', got:\n%s", yamlStr)
	}
	// Check that job names are extracted
	if !strings.Contains(yamlStr, "- build") || !strings.Contains(yamlStr, "- test") {
		t.Errorf("expected job names in needs, got:\n%s", yamlStr)
	}
}

// TestStepPointer tests that both Step and *Step work in steps array.
func TestStepPointer(t *testing.T) {
	step1 := workflow.Step{Run: "echo step1"}
	step2 := &workflow.Step{Run: "echo step2"}

	w := &workflow.Workflow{
		Name: "Test",
		On: workflow.Triggers{
			Push: &workflow.PushTrigger{},
		},
		Jobs: map[string]workflow.Job{
			"test": {
				RunsOn: "ubuntu-latest",
				Steps: []any{
					step1, // value
					step2, // pointer
				},
			},
		},
	}

	yaml, err := serialize.ToYAML(w)
	if err != nil {
		t.Fatalf("ToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	if !strings.Contains(yamlStr, "echo step1") {
		t.Errorf("expected step1, got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "echo step2") {
		t.Errorf("expected step2, got:\n%s", yamlStr)
	}
}

// TestExpressionInRunsOn tests Expression used in runs-on field.
func TestExpressionInRunsOn(t *testing.T) {
	w := &workflow.Workflow{
		Name: "Test",
		On: workflow.Triggers{
			Push: &workflow.PushTrigger{},
		},
		Jobs: map[string]workflow.Job{
			"test": {
				RunsOn: workflow.MatrixContext.Get("os"),
				Strategy: &workflow.Strategy{
					Matrix: &workflow.Matrix{
						Values: map[string][]any{
							"os": {"ubuntu-latest", "macos-latest"},
						},
					},
				},
				Steps: []any{
					workflow.Step{Run: "echo test"},
				},
			},
		},
	}

	yaml, err := serialize.ToYAML(w)
	if err != nil {
		t.Fatalf("ToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	if !strings.Contains(yamlStr, "${{ matrix.os }}") {
		t.Errorf("expected matrix expression in runs-on, got:\n%s", yamlStr)
	}
}

// TestStepWithAllFields tests step with all possible fields.
func TestStepWithAllFields(t *testing.T) {
	w := &workflow.Workflow{
		Name: "Test",
		On: workflow.Triggers{
			Push: &workflow.PushTrigger{},
		},
		Jobs: map[string]workflow.Job{
			"test": {
				RunsOn: "ubuntu-latest",
				Steps: []any{
					workflow.Step{
						ID:               "step1",
						Name:             "Test Step",
						If:               workflow.Success(),
						Uses:             "actions/checkout@v4",
						With:             workflow.With{"fetch-depth": 0},
						Run:              "echo test",
						Shell:            "bash",
						Env:              workflow.Env{"KEY": "value"},
						WorkingDirectory: "/tmp",
						ContinueOnError:  true,
						TimeoutMinutes:   30,
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

	expectedFields := []string{
		"id: step1",
		"name: Test Step",
		"if:",
		"uses: actions/checkout@v4",
		"with:",
		"run: echo test",
		"shell: bash",
		"env:",
		"working-directory: /tmp",
		"continue-on-error: true",
		"timeout-minutes: 30",
	}

	for _, field := range expectedFields {
		if !strings.Contains(yamlStr, field) {
			t.Errorf("expected field %q, got:\n%s", field, yamlStr)
		}
	}
}

// TestJobWithAllFields tests job with all possible fields.
func TestJobWithAllFields(t *testing.T) {
	w := &workflow.Workflow{
		Name: "Test",
		On: workflow.Triggers{
			Push: &workflow.PushTrigger{},
		},
		Jobs: map[string]workflow.Job{
			"test": {
				Name:   "Test Job",
				RunsOn: "ubuntu-latest",
				Permissions: &workflow.Permissions{
					Contents: "read",
					Issues:   "write",
				},
				Environment: &workflow.Environment{
					Name: "production",
					URL:  "https://example.com",
				},
				Concurrency: &workflow.Concurrency{
					Group:            "prod-deploy",
					CancelInProgress: false,
				},
				Outputs: workflow.Env{
					"result": workflow.Steps.Get("test", "result"),
				},
				Env: workflow.Env{
					"NODE_ENV": "production",
				},
				TimeoutMinutes:  60,
				ContinueOnError: true,
				Steps: []any{
					workflow.Step{Run: "echo test"},
				},
			},
		},
	}

	yaml, err := serialize.ToYAML(w)
	if err != nil {
		t.Fatalf("ToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	expectedFields := []string{
		"name: Test Job",
		"runs-on: ubuntu-latest",
		"permissions:",
		"environment:",
		"concurrency:",
		"outputs:",
		"env:",
		"timeout-minutes: 60",
		"continue-on-error: true",
	}

	for _, field := range expectedFields {
		if !strings.Contains(yamlStr, field) {
			t.Errorf("expected field %q, got:\n%s", field, yamlStr)
		}
	}
}

// TestRealStepAction tests using real action wrappers implementing StepAction.
func TestRealStepAction(t *testing.T) {
	// Create a mock that implements StepAction interface
	type mockStepAction struct{}

	mockAction := mockStepAction{}

	w := &workflow.Workflow{
		Name: "Test StepAction",
		On: workflow.Triggers{
			Push: &workflow.PushTrigger{},
		},
		Jobs: map[string]workflow.Job{
			"test": {
				RunsOn: "ubuntu-latest",
				Steps: []any{
					mockAction, // This should be converted via anyStepToMap
				},
			},
		},
	}

	_, err := serialize.ToYAML(w)
	// This should fail because mockStepAction doesn't implement StepAction
	if err == nil {
		t.Fatal("Expected error for non-StepAction type, got nil")
	}
}

// TestSerializeConditionTypes tests different condition value types.
func TestSerializeConditionTypes(t *testing.T) {
	tests := []struct {
		name      string
		condition any
		expected  string
	}{
		{
			name:      "expression",
			condition: workflow.Success(),
			expected:  "success()",
		},
		{
			name:      "string literal",
			condition: "always()",
			expected:  "always()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &workflow.Workflow{
				Name: "Test",
				On: workflow.Triggers{
					Push: &workflow.PushTrigger{},
				},
				Jobs: map[string]workflow.Job{
					"test": {
						If:     tt.condition,
						RunsOn: "ubuntu-latest",
						Steps: []any{
							workflow.Step{Run: "echo test"},
						},
					},
				},
			}

			yaml, err := serialize.ToYAML(w)
			if err != nil {
				t.Fatalf("ToYAML failed: %v", err)
			}

			yamlStr := string(yaml)
			if !strings.Contains(yamlStr, tt.expected) {
				t.Errorf("expected %q in output, got:\n%s", tt.expected, yamlStr)
			}
		})
	}
}

// TestNeedsWithStringAndReflection tests serializeNeeds with different types.
func TestNeedsWithStringAndReflection(t *testing.T) {
	// Test with direct string needs
	w := &workflow.Workflow{
		Name: "Test",
		On: workflow.Triggers{
			Push: &workflow.PushTrigger{},
		},
		Jobs: map[string]workflow.Job{
			"build": {
				Name:   "build",
				RunsOn: "ubuntu-latest",
				Steps: []any{
					workflow.Step{Run: "echo build"},
				},
			},
			"test": {
				Name:   "test",
				RunsOn: "ubuntu-latest",
				Needs:  []any{"build"}, // String reference
				Steps: []any{
					workflow.Step{Run: "echo test"},
				},
			},
		},
	}

	yaml, err := serialize.ToYAML(w)
	if err != nil {
		t.Fatalf("ToYAML failed: %v", err)
	}

	yamlStr := string(yaml)
	if !strings.Contains(yamlStr, "needs:") {
		t.Errorf("expected needs field, got:\n%s", yamlStr)
	}
}

// TestContainerAndServices tests container configuration.
func TestContainerAndServices(t *testing.T) {
	w := &workflow.Workflow{
		Name: "Test",
		On: workflow.Triggers{
			Push: &workflow.PushTrigger{},
		},
		Jobs: map[string]workflow.Job{
			"test": {
				RunsOn: "ubuntu-latest",
				Container: &workflow.Container{
					Image: "node:16",
					Credentials: &workflow.Credentials{
						Username: "${{ github.actor }}",
						Password: "${{ secrets.GITHUB_TOKEN }}",
					},
					Env: workflow.Env{
						"NODE_ENV": "test",
					},
					Ports:   []any{3000, 8080},
					Volumes: []string{"/tmp:/tmp"},
					Options: "--cpus 2",
				},
				Steps: []any{
					workflow.Step{Run: "npm test"},
				},
			},
		},
	}

	yaml, err := serialize.ToYAML(w)
	if err != nil {
		t.Fatalf("ToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	if !strings.Contains(yamlStr, "container:") {
		t.Errorf("expected container field, got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "image: node:16") {
		t.Errorf("expected container image, got:\n%s", yamlStr)
	}
}
