package importer

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestIRTriggers_UnmarshalYAML_String(t *testing.T) {
	yamlContent := []byte(`push`)

	var triggers IRTriggers
	err := yaml.Unmarshal(yamlContent, &triggers)
	if err != nil {
		t.Fatalf("UnmarshalYAML() error = %v", err)
	}

	if triggers.Push == nil {
		t.Error("Push trigger should be set")
	}
	if triggers.Raw != "push" {
		t.Errorf("Raw = %v, want %q", triggers.Raw, "push")
	}
}

func TestIRTriggers_UnmarshalYAML_Sequence(t *testing.T) {
	yamlContent := []byte(`[push, pull_request, workflow_dispatch]`)

	var triggers IRTriggers
	err := yaml.Unmarshal(yamlContent, &triggers)
	if err != nil {
		t.Fatalf("UnmarshalYAML() error = %v", err)
	}

	if triggers.Push == nil {
		t.Error("Push trigger should be set")
	}
	if triggers.PullRequest == nil {
		t.Error("PullRequest trigger should be set")
	}
	if triggers.WorkflowDispatch == nil {
		t.Error("WorkflowDispatch trigger should be set")
	}
}

func TestIRTriggers_UnmarshalYAML_Mapping(t *testing.T) {
	yamlContent := []byte(`
push:
  branches: [main]
pull_request:
  branches: [main]
`)

	var triggers IRTriggers
	err := yaml.Unmarshal(yamlContent, &triggers)
	if err != nil {
		t.Fatalf("UnmarshalYAML() error = %v", err)
	}

	if triggers.Push == nil {
		t.Error("Push trigger should be set")
	}
	if triggers.PullRequest == nil {
		t.Error("PullRequest trigger should be set")
	}
}

func TestIRTriggers_SetSimpleTrigger_AllTypes(t *testing.T) {
	tests := []struct {
		name     string
		trigger  string
		checkFn  func(*IRTriggers) bool
	}{
		{
			name:    "push",
			trigger: "push",
			checkFn: func(t *IRTriggers) bool { return t.Push != nil },
		},
		{
			name:    "pull_request",
			trigger: "pull_request",
			checkFn: func(t *IRTriggers) bool { return t.PullRequest != nil },
		},
		{
			name:    "pull_request_target",
			trigger: "pull_request_target",
			checkFn: func(t *IRTriggers) bool { return t.PullRequestTarget != nil },
		},
		{
			name:    "workflow_dispatch",
			trigger: "workflow_dispatch",
			checkFn: func(t *IRTriggers) bool { return t.WorkflowDispatch != nil },
		},
		{
			name:    "workflow_call",
			trigger: "workflow_call",
			checkFn: func(t *IRTriggers) bool { return t.WorkflowCall != nil },
		},
		{
			name:    "repository_dispatch",
			trigger: "repository_dispatch",
			checkFn: func(t *IRTriggers) bool { return t.RepositoryDispatch != nil },
		},
		{
			name:    "release",
			trigger: "release",
			checkFn: func(t *IRTriggers) bool { return t.Release != nil },
		},
		{
			name:    "issues",
			trigger: "issues",
			checkFn: func(t *IRTriggers) bool { return t.Issues != nil },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var triggers IRTriggers
			triggers.setSimpleTrigger(tt.trigger)

			if !tt.checkFn(&triggers) {
				t.Errorf("setSimpleTrigger(%q) did not set expected trigger", tt.trigger)
			}
		})
	}
}

func TestIRTriggers_SetSimpleTrigger_Unknown(t *testing.T) {
	var triggers IRTriggers
	triggers.setSimpleTrigger("unknown_trigger")

	// Should not panic, just do nothing
	if triggers.Push != nil {
		t.Error("Unknown trigger should not set any fields")
	}
}

func TestIRJob_GetNeeds_AnySliceWithNonString(t *testing.T) {
	job := &IRJob{
		Needs: []any{"build", 123, "test"},
	}
	needs := job.GetNeeds()

	// Should only include string values
	if len(needs) != 2 {
		t.Errorf("GetNeeds() = %v, want 2 string elements", needs)
	}
	if needs[0] != "build" || needs[1] != "test" {
		t.Errorf("GetNeeds() = %v, want [build test]", needs)
	}
}

func TestIRJob_GetNeeds_InvalidType(t *testing.T) {
	job := &IRJob{
		Needs: 123,
	}
	needs := job.GetNeeds()

	if needs != nil {
		t.Errorf("GetNeeds() = %v, want nil for invalid type", needs)
	}
}

func TestIRJob_GetRunsOn_EmptySlice(t *testing.T) {
	job := &IRJob{
		RunsOn: []any{},
	}

	result := job.GetRunsOn()
	if result != "" {
		t.Errorf("GetRunsOn() = %q, want empty for empty slice", result)
	}
}

func TestIRJob_GetRunsOn_SliceWithNonString(t *testing.T) {
	job := &IRJob{
		RunsOn: []any{123, "ubuntu-latest"},
	}

	result := job.GetRunsOn()
	if result != "" {
		t.Errorf("GetRunsOn() = %q, want empty when first element is not string", result)
	}
}

func TestIRJob_GetRunsOn_InvalidType(t *testing.T) {
	job := &IRJob{
		RunsOn: 123,
	}

	result := job.GetRunsOn()
	if result != "" {
		t.Errorf("GetRunsOn() = %q, want empty for invalid type", result)
	}
}

func TestIRJob_GetRunsOn_EmptyStringSlice(t *testing.T) {
	job := &IRJob{
		RunsOn: []string{},
	}

	result := job.GetRunsOn()
	if result != "" {
		t.Errorf("GetRunsOn() = %q, want empty for empty string slice", result)
	}
}
