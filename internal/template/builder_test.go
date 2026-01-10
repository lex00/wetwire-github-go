package template

import (
	"reflect"
	"strings"
	"testing"

	"github.com/lex00/wetwire-github-go/internal/discover"
	"github.com/lex00/wetwire-github-go/internal/runner"
	"github.com/lex00/wetwire-github-go/workflow"
)

func TestNewBuilder(t *testing.T) {
	b := NewBuilder()
	if b == nil {
		t.Fatal("NewBuilder() returned nil")
	}
}

func TestBuilder_Build_Empty(t *testing.T) {
	b := NewBuilder()

	discovered := &discover.DiscoveryResult{
		Workflows: []discover.DiscoveredWorkflow{},
		Jobs:      []discover.DiscoveredJob{},
	}
	extracted := &runner.ExtractionResult{
		Workflows: []runner.ExtractedWorkflow{},
		Jobs:      []runner.ExtractedJob{},
	}

	result, err := b.Build(discovered, extracted)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	if len(result.Workflows) != 0 {
		t.Errorf("Expected 0 workflows, got %d", len(result.Workflows))
	}
}

func TestBuilder_Build_SimpleWorkflow(t *testing.T) {
	b := NewBuilder()

	discovered := &discover.DiscoveryResult{
		Workflows: []discover.DiscoveredWorkflow{
			{Name: "CI", File: "ci.go", Line: 10, Jobs: []string{"Build"}},
		},
		Jobs: []discover.DiscoveredJob{
			{Name: "Build", File: "ci.go", Line: 20, Dependencies: []string{}},
		},
	}

	extracted := &runner.ExtractionResult{
		Workflows: []runner.ExtractedWorkflow{
			{
				Name: "CI",
				Data: map[string]any{
					"Name": "CI",
					"On":   workflow.Triggers{Push: &workflow.PushTrigger{}},
				},
			},
		},
		Jobs: []runner.ExtractedJob{
			{
				Name: "Build",
				Data: map[string]any{
					"Name":   "build",
					"RunsOn": "ubuntu-latest",
					"Steps":  []workflow.Step{{Run: "echo hello"}},
				},
			},
		},
	}

	result, err := b.Build(discovered, extracted)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	if len(result.Workflows) != 1 {
		t.Fatalf("Expected 1 workflow, got %d", len(result.Workflows))
	}

	wf := result.Workflows[0]
	if wf.Name != "CI" {
		t.Errorf("Workflow name = %q, want %q", wf.Name, "CI")
	}

	if len(wf.YAML) == 0 {
		t.Error("Workflow YAML is empty")
	}
}

func TestBuilder_Build_WithDependencies(t *testing.T) {
	b := NewBuilder()

	discovered := &discover.DiscoveryResult{
		Workflows: []discover.DiscoveredWorkflow{
			{Name: "CI", File: "ci.go", Line: 10, Jobs: []string{"Build", "Test", "Deploy"}},
		},
		Jobs: []discover.DiscoveredJob{
			{Name: "Build", File: "ci.go", Line: 20, Dependencies: []string{}},
			{Name: "Test", File: "ci.go", Line: 30, Dependencies: []string{"Build"}},
			{Name: "Deploy", File: "ci.go", Line: 40, Dependencies: []string{"Build", "Test"}},
		},
	}

	extracted := &runner.ExtractionResult{
		Workflows: []runner.ExtractedWorkflow{
			{
				Name: "CI",
				Data: map[string]any{
					"Name": "CI",
				},
			},
		},
		Jobs: []runner.ExtractedJob{
			{Name: "Build", Data: map[string]any{"Name": "build", "RunsOn": "ubuntu-latest"}},
			{Name: "Test", Data: map[string]any{"Name": "test", "RunsOn": "ubuntu-latest"}},
			{Name: "Deploy", Data: map[string]any{"Name": "deploy", "RunsOn": "ubuntu-latest"}},
		},
	}

	result, err := b.Build(discovered, extracted)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	if len(result.Workflows) != 1 {
		t.Fatalf("Expected 1 workflow, got %d", len(result.Workflows))
	}

	// Check job ordering
	wf := result.Workflows[0]
	jobs := wf.Jobs

	// Build should come before Test and Deploy
	buildIdx := indexOf(jobs, "Build")
	testIdx := indexOf(jobs, "Test")
	deployIdx := indexOf(jobs, "Deploy")

	if buildIdx > testIdx {
		t.Error("Build should come before Test")
	}
	if buildIdx > deployIdx {
		t.Error("Build should come before Deploy")
	}
	if testIdx > deployIdx {
		t.Error("Test should come before Deploy")
	}
}

func indexOf(jobs []string, name string) int {
	for i, j := range jobs {
		if j == name {
			return i
		}
	}
	return -1
}

func TestBuilder_Build_CycleDetection(t *testing.T) {
	b := NewBuilder()

	discovered := &discover.DiscoveryResult{
		Workflows: []discover.DiscoveredWorkflow{
			{Name: "CI", File: "ci.go", Line: 10, Jobs: []string{"A", "B", "C"}},
		},
		Jobs: []discover.DiscoveredJob{
			{Name: "A", File: "ci.go", Line: 20, Dependencies: []string{"C"}},
			{Name: "B", File: "ci.go", Line: 30, Dependencies: []string{"A"}},
			{Name: "C", File: "ci.go", Line: 40, Dependencies: []string{"B"}},
		},
	}

	extracted := &runner.ExtractionResult{
		Workflows: []runner.ExtractedWorkflow{
			{Name: "CI", Data: map[string]any{"Name": "CI"}},
		},
		Jobs: []runner.ExtractedJob{
			{Name: "A", Data: map[string]any{"Name": "a"}},
			{Name: "B", Data: map[string]any{"Name": "b"}},
			{Name: "C", Data: map[string]any{"Name": "c"}},
		},
	}

	_, err := b.Build(discovered, extracted)
	if err == nil {
		t.Error("Build() expected error for cycle")
	}
	if !strings.Contains(err.Error(), "cycle") {
		t.Errorf("Error should mention cycle, got: %v", err)
	}
}

func TestOrderJobs(t *testing.T) {
	jobs := []discover.DiscoveredJob{
		{Name: "Deploy", Dependencies: []string{"Test"}},
		{Name: "Test", Dependencies: []string{"Build"}},
		{Name: "Build", Dependencies: []string{}},
	}

	result, err := OrderJobs(jobs)
	if err != nil {
		t.Fatalf("OrderJobs() error = %v", err)
	}

	expected := []string{"Build", "Test", "Deploy"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("OrderJobs() = %v, want %v", result, expected)
	}
}

func TestOrderJobs_Cycle(t *testing.T) {
	jobs := []discover.DiscoveredJob{
		{Name: "A", Dependencies: []string{"B"}},
		{Name: "B", Dependencies: []string{"A"}},
	}

	_, err := OrderJobs(jobs)
	if err == nil {
		t.Error("OrderJobs() expected error for cycle")
	}
}

func TestValidateJobDependencies_Valid(t *testing.T) {
	jobs := []discover.DiscoveredJob{
		{Name: "Build", Dependencies: []string{}},
		{Name: "Test", Dependencies: []string{"Build"}},
		{Name: "Deploy", Dependencies: []string{"Build", "Test"}},
	}

	errors := ValidateJobDependencies(jobs)
	if len(errors) != 0 {
		t.Errorf("ValidateJobDependencies() returned errors: %v", errors)
	}
}

func TestValidateJobDependencies_UnknownDependency(t *testing.T) {
	jobs := []discover.DiscoveredJob{
		{Name: "Build", Dependencies: []string{}},
		{Name: "Test", Dependencies: []string{"Unknown"}},
	}

	errors := ValidateJobDependencies(jobs)
	if len(errors) == 0 {
		t.Error("ValidateJobDependencies() should return error for unknown dependency")
	}

	found := false
	for _, err := range errors {
		if strings.Contains(err, "Unknown") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Error should mention unknown dependency, got: %v", errors)
	}
}

func TestValidateJobDependencies_Cycle(t *testing.T) {
	jobs := []discover.DiscoveredJob{
		{Name: "A", Dependencies: []string{"B"}},
		{Name: "B", Dependencies: []string{"A"}},
	}

	errors := ValidateJobDependencies(jobs)
	if len(errors) == 0 {
		t.Error("ValidateJobDependencies() should return error for cycle")
	}

	found := false
	for _, err := range errors {
		if strings.Contains(err, "cycle") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Error should mention cycle, got: %v", errors)
	}
}

func TestBuilder_filterAndOrderJobs(t *testing.T) {
	b := NewBuilder()

	sortedJobs := []string{"A", "B", "C", "D", "E"}
	jobNames := []string{"B", "D", "E"}

	result := b.filterAndOrderJobs(jobNames, sortedJobs)

	expected := []string{"B", "D", "E"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("filterAndOrderJobs() = %v, want %v", result, expected)
	}
}

func TestBuilder_Build_MissingWorkflowData(t *testing.T) {
	b := NewBuilder()

	discovered := &discover.DiscoveryResult{
		Workflows: []discover.DiscoveredWorkflow{
			{Name: "CI", File: "ci.go", Line: 10, Jobs: []string{}},
		},
		Jobs: []discover.DiscoveredJob{},
	}

	// No extraction data for the workflow
	extracted := &runner.ExtractionResult{
		Workflows: []runner.ExtractedWorkflow{},
		Jobs:      []runner.ExtractedJob{},
	}

	result, err := b.Build(discovered, extracted)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	// Should have an error about missing extraction data
	if len(result.Errors) == 0 {
		t.Error("Expected error about missing extraction data")
	}
}

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
			"ID":                "test",
			"Name":              "Test",
			"Run":               "go test ./...",
			"WorkingDirectory":  "/app",
			"TimeoutMinutes":    30.0,
			"If":                "success()",
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
		name        string
		discovered  discover.DiscoveredWorkflow
		data        map[string]any
		jobNames    []string
		jobMap      map[string]*runner.ExtractedJob
		sortedJobs  []string
		wantErr     bool
		validate    func(*workflow.Workflow) bool
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

func TestBuilder_Build_SerializationError(t *testing.T) {
	b := NewBuilder()

	discovered := &discover.DiscoveryResult{
		Workflows: []discover.DiscoveredWorkflow{
			{Name: "CI", File: "ci.go", Line: 10, Jobs: []string{"Build"}},
		},
		Jobs: []discover.DiscoveredJob{
			{Name: "Build", File: "ci.go", Line: 20, Dependencies: []string{}},
		},
	}

	// Create a workflow that might cause serialization issues
	extracted := &runner.ExtractionResult{
		Workflows: []runner.ExtractedWorkflow{
			{
				Name: "CI",
				Data: map[string]any{
					"Name": "CI",
				},
			},
		},
		Jobs: []runner.ExtractedJob{
			{
				Name: "Build",
				Data: map[string]any{
					"Name":   "build",
					"RunsOn": "ubuntu-latest",
				},
			},
		},
	}

	result, err := b.Build(discovered, extracted)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	// Even with potential serialization issues, build should not fail
	// It should capture errors in result.Errors
	if len(result.Workflows) == 0 && len(result.Errors) == 0 {
		t.Error("Expected either workflows or errors, got neither")
	}
}

func TestBuilder_Build_JobOrdering(t *testing.T) {
	b := NewBuilder()

	discovered := &discover.DiscoveryResult{
		Workflows: []discover.DiscoveredWorkflow{
			{Name: "CI", File: "ci.go", Line: 10, Jobs: []string{"Deploy", "Test", "Build"}},
		},
		Jobs: []discover.DiscoveredJob{
			{Name: "Build", File: "ci.go", Line: 20, Dependencies: []string{}},
			{Name: "Test", File: "ci.go", Line: 30, Dependencies: []string{"Build"}},
			{Name: "Deploy", File: "ci.go", Line: 40, Dependencies: []string{"Test"}},
		},
	}

	extracted := &runner.ExtractionResult{
		Workflows: []runner.ExtractedWorkflow{
			{
				Name: "CI",
				Data: map[string]any{
					"Name": "CI",
				},
			},
		},
		Jobs: []runner.ExtractedJob{
			{Name: "Build", Data: map[string]any{"Name": "build", "RunsOn": "ubuntu-latest"}},
			{Name: "Test", Data: map[string]any{"Name": "test", "RunsOn": "ubuntu-latest"}},
			{Name: "Deploy", Data: map[string]any{"Name": "deploy", "RunsOn": "ubuntu-latest"}},
		},
	}

	result, err := b.Build(discovered, extracted)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	if len(result.Workflows) != 1 {
		t.Fatalf("Expected 1 workflow, got %d", len(result.Workflows))
	}

	wf := result.Workflows[0]
	jobs := wf.Jobs

	// Verify jobs are ordered: Build, Test, Deploy
	if len(jobs) != 3 {
		t.Fatalf("Expected 3 jobs, got %d", len(jobs))
	}

	if jobs[0] != "Build" {
		t.Errorf("First job should be Build, got %s", jobs[0])
	}
	if jobs[1] != "Test" {
		t.Errorf("Second job should be Test, got %s", jobs[1])
	}
	if jobs[2] != "Deploy" {
		t.Errorf("Third job should be Deploy, got %s", jobs[2])
	}
}
