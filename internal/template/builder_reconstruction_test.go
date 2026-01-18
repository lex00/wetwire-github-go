package template

import (
	"testing"

	"github.com/lex00/wetwire-github-go/internal/discover"
	"github.com/lex00/wetwire-github-go/internal/runner"
	"github.com/lex00/wetwire-github-go/workflow"
)

func TestBuilder_reconstructTriggers(t *testing.T) {
	b := NewBuilder()

	tests := []struct {
		name string
		data map[string]any
		want func(workflow.Triggers) bool
	}{
		{
			name: "push trigger",
			data: map[string]any{
				"Push": map[string]any{
					"Branches": []any{"main", "develop"},
					"Tags":     []any{"v*"},
					"Paths":    []any{"src/**"},
				},
			},
			want: func(t workflow.Triggers) bool {
				return t.Push != nil &&
					len(t.Push.Branches) == 2 &&
					t.Push.Branches[0] == "main" &&
					len(t.Push.Tags) == 1 &&
					len(t.Push.Paths) == 1
			},
		},
		{
			name: "pull_request trigger",
			data: map[string]any{
				"PullRequest": map[string]any{
					"Branches": []any{"main"},
					"Types":    []any{"opened", "synchronize"},
				},
			},
			want: func(t workflow.Triggers) bool {
				return t.PullRequest != nil &&
					len(t.PullRequest.Branches) == 1 &&
					len(t.PullRequest.Types) == 2
			},
		},
		{
			name: "workflow_dispatch trigger",
			data: map[string]any{
				"WorkflowDispatch": map[string]any{},
			},
			want: func(t workflow.Triggers) bool {
				return t.WorkflowDispatch != nil
			},
		},
		{
			name: "workflow_call trigger",
			data: map[string]any{
				"WorkflowCall": map[string]any{},
			},
			want: func(t workflow.Triggers) bool {
				return t.WorkflowCall != nil
			},
		},
		{
			name: "schedule trigger",
			data: map[string]any{
				"Schedule": []any{
					map[string]any{"Cron": "0 0 * * *"},
					map[string]any{"Cron": "0 12 * * *"},
				},
			},
			want: func(t workflow.Triggers) bool {
				return len(t.Schedule) == 2 &&
					t.Schedule[0].Cron == "0 0 * * *" &&
					t.Schedule[1].Cron == "0 12 * * *"
			},
		},
		{
			name: "multiple triggers",
			data: map[string]any{
				"Push": map[string]any{
					"Branches": []any{"main"},
				},
				"PullRequest": map[string]any{
					"Branches": []any{"main"},
				},
				"WorkflowDispatch": map[string]any{},
			},
			want: func(t workflow.Triggers) bool {
				return t.Push != nil && t.PullRequest != nil && t.WorkflowDispatch != nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := b.reconstructTriggers(tt.data)
			if !tt.want(result) {
				t.Errorf("reconstructTriggers() validation failed")
			}
		})
	}
}

func TestBuilder_reconstructStepsAsAny(t *testing.T) {
	b := NewBuilder()

	data := []any{
		map[string]any{
			"ID":   "checkout",
			"Name": "Checkout code",
			"Uses": "actions/checkout@v4",
		},
		map[string]any{
			"ID":    "build",
			"Name":  "Build",
			"Run":   "go build ./...",
			"Shell": "bash",
		},
		map[string]any{
			"ID":               "test",
			"Name":             "Test",
			"Run":              "go test ./...",
			"WorkingDirectory": "/app",
			"TimeoutMinutes":   30.0,
			"If":               "success()",
		},
		map[string]any{
			"Name": "With env",
			"Run":  "echo $TEST",
			"Env": map[string]any{
				"TEST": "value",
			},
		},
		map[string]any{
			"Name": "With with",
			"Uses": "actions/setup-go@v5",
			"With": map[string]any{
				"go-version": "1.23",
			},
		},
	}

	result := b.reconstructStepsAsAny(data)

	if len(result) != 5 {
		t.Fatalf("Expected 5 steps, got %d", len(result))
	}

	// Check first step
	step0, ok := result[0].(workflow.Step)
	if !ok {
		t.Fatalf("First element is not a workflow.Step, got %T", result[0])
	}
	if step0.ID != "checkout" {
		t.Errorf("First step ID = %q, want %q", step0.ID, "checkout")
	}
	if step0.Uses != "actions/checkout@v4" {
		t.Errorf("First step Uses = %q, want %q", step0.Uses, "actions/checkout@v4")
	}

	// Check second step
	step1, ok := result[1].(workflow.Step)
	if !ok {
		t.Fatalf("Second element is not a workflow.Step")
	}
	if step1.Run != "go build ./..." {
		t.Errorf("Second step Run = %q, want %q", step1.Run, "go build ./...")
	}
	if step1.Shell != "bash" {
		t.Errorf("Second step Shell = %q, want %q", step1.Shell, "bash")
	}

	// Check third step with advanced options
	step2, ok := result[2].(workflow.Step)
	if !ok {
		t.Fatalf("Third element is not a workflow.Step")
	}
	if step2.WorkingDirectory != "/app" {
		t.Errorf("Third step WorkingDirectory = %q, want %q", step2.WorkingDirectory, "/app")
	}
	if step2.TimeoutMinutes != 30 {
		t.Errorf("Third step TimeoutMinutes = %d, want %d", step2.TimeoutMinutes, 30)
	}
	if step2.If != "success()" {
		t.Errorf("Third step If = %q, want %q", step2.If, "success()")
	}

	// Check fourth step with env
	step3, ok := result[3].(workflow.Step)
	if !ok {
		t.Fatalf("Fourth element is not a workflow.Step")
	}
	if step3.Env == nil {
		t.Fatal("Fourth step Env is nil")
	}
	if val, ok := step3.Env["TEST"]; !ok || val != "value" {
		t.Errorf("Fourth step Env[TEST] = %v, want %q", val, "value")
	}

	// Check fifth step with with
	step4, ok := result[4].(workflow.Step)
	if !ok {
		t.Fatalf("Fifth element is not a workflow.Step")
	}
	if step4.With == nil {
		t.Fatal("Fifth step With is nil")
	}
	if val, ok := step4.With["go-version"]; !ok || val != "1.23" {
		t.Errorf("Fifth step With[go-version] = %v, want %q", val, "1.23")
	}
}

func TestBuilder_anySliceToStrings(t *testing.T) {
	tests := []struct {
		name  string
		input []any
		want  []string
	}{
		{
			name:  "all strings",
			input: []any{"a", "b", "c"},
			want:  []string{"a", "b", "c"},
		},
		{
			name:  "mixed types - only strings extracted",
			input: []any{"a", 123, "b", true, "c"},
			want:  []string{"a", "b", "c"},
		},
		{
			name:  "empty slice",
			input: []any{},
			want:  []string{},
		},
		{
			name:  "no strings",
			input: []any{123, true, 45.6},
			want:  []string{},
		},
		{
			name:  "nil slice",
			input: nil,
			want:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := anySliceToStrings(tt.input)
			if len(result) != len(tt.want) {
				t.Errorf("anySliceToStrings() length = %d, want %d", len(result), len(tt.want))
				return
			}
			for i := range result {
				if result[i] != tt.want[i] {
					t.Errorf("anySliceToStrings()[%d] = %q, want %q", i, result[i], tt.want[i])
				}
			}
		})
	}
}

func TestBuilder_buildJob(t *testing.T) {
	b := NewBuilder()

	tests := []struct {
		name string
		job  *runner.ExtractedJob
		want func(*workflow.Job) bool
	}{
		{
			name: "basic job",
			job: &runner.ExtractedJob{
				Name: "Build",
				Data: map[string]any{
					"Name":   "build",
					"RunsOn": "ubuntu-latest",
				},
			},
			want: func(j *workflow.Job) bool {
				return j.Name == "build" && j.RunsOn == "ubuntu-latest"
			},
		},
		{
			name: "job with env",
			job: &runner.ExtractedJob{
				Name: "Test",
				Data: map[string]any{
					"Name":   "test",
					"RunsOn": "ubuntu-latest",
					"Env": map[string]any{
						"GO_VERSION": "1.23",
					},
				},
			},
			want: func(j *workflow.Job) bool {
				return j.Env != nil && j.Env["GO_VERSION"] == "1.23"
			},
		},
		{
			name: "job with timeout",
			job: &runner.ExtractedJob{
				Name: "LongJob",
				Data: map[string]any{
					"Name":           "long-job",
					"RunsOn":         "ubuntu-latest",
					"TimeoutMinutes": 60,
				},
			},
			want: func(j *workflow.Job) bool {
				return j.TimeoutMinutes == 60
			},
		},
		{
			name: "job with continue-on-error",
			job: &runner.ExtractedJob{
				Name: "Lint",
				Data: map[string]any{
					"Name":            "lint",
					"RunsOn":          "ubuntu-latest",
					"ContinueOnError": true,
				},
			},
			want: func(j *workflow.Job) bool {
				return j.ContinueOnError == true
			},
		},
		{
			name: "job with needs",
			job: &runner.ExtractedJob{
				Name: "Deploy",
				Data: map[string]any{
					"Name":   "deploy",
					"RunsOn": "ubuntu-latest",
					"Needs":  []any{"build", "test"},
				},
			},
			want: func(j *workflow.Job) bool {
				return len(j.Needs) == 2
			},
		},
		{
			name: "job with outputs",
			job: &runner.ExtractedJob{
				Name: "Version",
				Data: map[string]any{
					"Name":   "version",
					"RunsOn": "ubuntu-latest",
					"Outputs": map[string]any{
						"version": "${{ steps.get_version.outputs.version }}",
					},
				},
			},
			want: func(j *workflow.Job) bool {
				return j.Outputs != nil && len(j.Outputs) == 1
			},
		},
		{
			name: "job with steps as []workflow.Step",
			job: &runner.ExtractedJob{
				Name: "WithSteps",
				Data: map[string]any{
					"Name":   "with-steps",
					"RunsOn": "ubuntu-latest",
					"Steps": []workflow.Step{
						{Run: "echo hello"},
						{Run: "echo world"},
					},
				},
			},
			want: func(j *workflow.Job) bool {
				return len(j.Steps) == 2
			},
		},
		{
			name: "job with steps as []any",
			job: &runner.ExtractedJob{
				Name: "WithAnySteps",
				Data: map[string]any{
					"Name":   "with-any-steps",
					"RunsOn": "ubuntu-latest",
					"Steps": []any{
						map[string]any{"Run": "echo test"},
					},
				},
			},
			want: func(j *workflow.Job) bool {
				return len(j.Steps) == 1
			},
		},
		{
			name: "job with if condition",
			job: &runner.ExtractedJob{
				Name: "ConditionalJob",
				Data: map[string]any{
					"Name":   "conditional",
					"RunsOn": "ubuntu-latest",
					"If":     "github.ref == 'refs/heads/main'",
				},
			},
			want: func(j *workflow.Job) bool {
				return j.If == "github.ref == 'refs/heads/main'"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := b.buildJob(tt.job)
			if err != nil {
				t.Fatalf("buildJob() error = %v", err)
			}
			if result == nil {
				t.Fatal("buildJob() returned nil")
			}
			if !tt.want(result) {
				t.Errorf("buildJob() validation failed")
			}
		})
	}
}

func TestBuilder_buildWorkflow(t *testing.T) {
	b := NewBuilder()

	tests := []struct {
		name       string
		discovered discover.DiscoveredWorkflow
		data       map[string]any
		jobNames   []string
		jobMap     map[string]*runner.ExtractedJob
		sortedJobs []string
		wantErr    bool
		validate   func(*workflow.Workflow) bool
	}{
		{
			name: "basic workflow",
			discovered: discover.DiscoveredWorkflow{
				Name: "CI",
				Jobs: []string{"Build"},
			},
			data: map[string]any{
				"Name": "Continuous Integration",
			},
			jobNames: []string{"Build"},
			jobMap: map[string]*runner.ExtractedJob{
				"Build": {
					Name: "Build",
					Data: map[string]any{
						"Name":   "build",
						"RunsOn": "ubuntu-latest",
					},
				},
			},
			sortedJobs: []string{"Build"},
			wantErr:    false,
			validate: func(wf *workflow.Workflow) bool {
				return wf.Name == "Continuous Integration" && len(wf.Jobs) == 1
			},
		},
		{
			name: "workflow with env",
			discovered: discover.DiscoveredWorkflow{
				Name: "CI",
				Jobs: []string{},
			},
			data: map[string]any{
				"Name": "CI",
				"Env": map[string]any{
					"NODE_VERSION": "18",
				},
			},
			jobNames:   []string{},
			jobMap:     map[string]*runner.ExtractedJob{},
			sortedJobs: []string{},
			wantErr:    false,
			validate: func(wf *workflow.Workflow) bool {
				return wf.Env != nil && wf.Env["NODE_VERSION"] == "18"
			},
		},
		{
			name: "workflow with triggers as workflow.Triggers",
			discovered: discover.DiscoveredWorkflow{
				Name: "CI",
				Jobs: []string{},
			},
			data: map[string]any{
				"Name": "CI",
				"On":   workflow.Triggers{Push: &workflow.PushTrigger{Branches: []string{"main"}}},
			},
			jobNames:   []string{},
			jobMap:     map[string]*runner.ExtractedJob{},
			sortedJobs: []string{},
			wantErr:    false,
			validate: func(wf *workflow.Workflow) bool {
				return wf.On.Push != nil && len(wf.On.Push.Branches) == 1
			},
		},
		{
			name: "workflow with triggers as map",
			discovered: discover.DiscoveredWorkflow{
				Name: "CI",
				Jobs: []string{},
			},
			data: map[string]any{
				"Name": "CI",
				"On": map[string]any{
					"Push": map[string]any{
						"Branches": []any{"main", "develop"},
					},
				},
			},
			jobNames:   []string{},
			jobMap:     map[string]*runner.ExtractedJob{},
			sortedJobs: []string{},
			wantErr:    false,
			validate: func(wf *workflow.Workflow) bool {
				return wf.On.Push != nil && len(wf.On.Push.Branches) == 2
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := b.buildWorkflow(tt.discovered, tt.data, tt.jobNames, tt.jobMap, tt.sortedJobs)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildWorkflow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && result == nil {
				t.Fatal("buildWorkflow() returned nil without error")
			}
			if err == nil && !tt.validate(result) {
				t.Errorf("buildWorkflow() validation failed")
			}
		})
	}
}
