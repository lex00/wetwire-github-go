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
