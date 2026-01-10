package serialize_test

import (
	"strings"
	"testing"

	"github.com/lex00/wetwire-github-go/codeowners"
	"github.com/lex00/wetwire-github-go/dependabot"
	"github.com/lex00/wetwire-github-go/internal/serialize"
	"github.com/lex00/wetwire-github-go/templates"
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

func TestExpressionInEnv(t *testing.T) {
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

// TestExpressionSerialization tests various expression types and their serialization.
func TestExpressionSerialization(t *testing.T) {
	tests := []struct {
		name     string
		expr     workflow.Expression
		expected string
	}{
		{
			name:     "secrets expression",
			expr:     workflow.Secrets.Get("MY_SECRET"),
			expected: "secrets.MY_SECRET",
		},
		{
			name:     "github context",
			expr:     workflow.GitHub.SHA(),
			expected: "github.sha",
		},
		{
			name:     "matrix context",
			expr:     workflow.MatrixContext.Get("os"),
			expected: "matrix.os",
		},
		{
			name:     "combined expression",
			expr:     workflow.Branch("main").And(workflow.Success()),
			expected: "(github.ref == 'refs/heads/main') && (success())",
		},
		{
			name:     "negated expression",
			expr:     workflow.Failure().Not(),
			expected: "!(failure())",
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
						RunsOn: "ubuntu-latest",
						Steps: []any{
							workflow.Step{
								If:  tt.expr,
								Run: "echo test",
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
			if !strings.Contains(yamlStr, tt.expected) {
				t.Errorf("expected %q in output, got:\n%s", tt.expected, yamlStr)
			}
		})
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
				Permissions:  nil,
				Environment:  nil,
				Concurrency:  nil,
				Strategy:     nil,
				Container:    nil,
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

// TestAllTriggers tests all trigger types.
func TestAllTriggers(t *testing.T) {
	w := &workflow.Workflow{
		Name: "Test All Triggers",
		On: workflow.Triggers{
			Push:              &workflow.PushTrigger{Branches: []string{"main"}},
			PullRequest:       &workflow.PullRequestTrigger{Types: []string{"opened"}},
			PullRequestTarget: &workflow.PullRequestTargetTrigger{Branches: []string{"main"}},
			Schedule:          []workflow.ScheduleTrigger{{Cron: "0 0 * * *"}},
			WorkflowDispatch:  &workflow.WorkflowDispatchTrigger{},
			WorkflowRun: &workflow.WorkflowRunTrigger{
				Workflows: []string{"CI"},
			},
			RepositoryDispatch: &workflow.RepositoryDispatchTrigger{
				Types: []string{"deploy"},
			},
			IssueComment: &workflow.IssueCommentTrigger{Types: []string{"created"}},
			Issues:       &workflow.IssuesTrigger{Types: []string{"opened"}},
			Release:      &workflow.ReleaseTrigger{Types: []string{"published"}},
			Create:       &workflow.CreateTrigger{},
			Delete:       &workflow.DeleteTrigger{},
		},
	}

	yaml, err := serialize.ToYAML(w)
	if err != nil {
		t.Fatalf("ToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	expectedTriggers := []string{
		"push:", "pull_request:", "pull_request_target:",
		"schedule:", "workflow_dispatch:", "workflow_run:",
		"repository_dispatch:", "issue_comment:", "issues:",
		"release:", "create:", "delete:",
	}

	for _, trigger := range expectedTriggers {
		if !strings.Contains(yamlStr, trigger) {
			t.Errorf("expected trigger %q, got:\n%s", trigger, yamlStr)
		}
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
					step1,  // value
					step2,  // pointer
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

// TestComplexWorkflowCall tests workflow_call with inputs, outputs, and secrets.
func TestComplexWorkflowCall(t *testing.T) {
	w := &workflow.Workflow{
		Name: "Reusable",
		On: workflow.Triggers{
			WorkflowCall: &workflow.WorkflowCallTrigger{
				Inputs: map[string]workflow.WorkflowInput{
					"environment": {
						Type:        "string",
						Required:    true,
						Description: "Deployment environment",
					},
				},
				Outputs: map[string]workflow.WorkflowOutput{
					"deployment_id": {
						Description: "The deployment ID",
						Value:       workflow.Expression("jobs.deploy.outputs.id"),
					},
				},
				Secrets: map[string]workflow.WorkflowSecret{
					"deploy_key": {
						Description: "Deployment key",
						Required:    true,
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

	if !strings.Contains(yamlStr, "workflow_call:") {
		t.Errorf("expected 'workflow_call:', got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "inputs:") {
		t.Errorf("expected 'inputs:', got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "outputs:") {
		t.Errorf("expected 'outputs:', got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "secrets:") {
		t.Errorf("expected 'secrets:', got:\n%s", yamlStr)
	}
}

// TestPushTriggerPathFilters tests push trigger with path filters.
func TestPushTriggerPathFilters(t *testing.T) {
	w := &workflow.Workflow{
		Name: "Test",
		On: workflow.Triggers{
			Push: &workflow.PushTrigger{
				Branches:       []string{"main"},
				BranchesIgnore: []string{"dev"},
				Tags:           []string{"v*"},
				TagsIgnore:     []string{"v*-beta"},
				Paths:          []string{"src/**"},
				PathsIgnore:    []string{"docs/**"},
			},
		},
	}

	yaml, err := serialize.ToYAML(w)
	if err != nil {
		t.Fatalf("ToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	expectedFields := []string{
		"branches:", "branches-ignore:", "tags:",
		"tags-ignore:", "paths:", "paths-ignore:",
	}

	for _, field := range expectedFields {
		if !strings.Contains(yamlStr, field) {
			t.Errorf("expected field %q, got:\n%s", field, yamlStr)
		}
	}
}

// ===== Dependabot Tests =====

// TestBasicDependabot tests basic dependabot configuration serialization.
func TestBasicDependabot(t *testing.T) {
	d := &dependabot.Dependabot{
		Version: 2,
		Updates: []dependabot.Update{
			{
				PackageEcosystem: "go",
				Directory:        "/",
				Schedule: dependabot.Schedule{
					Interval: "daily",
				},
			},
		},
	}

	yaml, err := serialize.DependabotToYAML(d)
	if err != nil {
		t.Fatalf("DependabotToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	if !strings.Contains(yamlStr, "version: 2") {
		t.Errorf("expected 'version: 2', got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "package-ecosystem: go") {
		t.Errorf("expected 'package-ecosystem: go', got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "interval: daily") {
		t.Errorf("expected 'interval: daily', got:\n%s", yamlStr)
	}
}

// TestDependabotWithAllFields tests dependabot with all optional fields.
func TestDependabotWithAllFields(t *testing.T) {
	d := &dependabot.Dependabot{
		Version:              2,
		EnableBetaEcosystems: true,
		Updates: []dependabot.Update{
			{
				PackageEcosystem: "npm",
				Directory:        "/frontend",
				Schedule: dependabot.Schedule{
					Interval: "weekly",
					Day:      "monday",
					Time:     "10:00",
					Timezone: "America/New_York",
				},
				Allow: []dependabot.Allow{
					{DependencyName: "lodash"},
					{DependencyType: "production"},
				},
				Ignore: []dependabot.Ignore{
					{
						DependencyName: "webpack",
						Versions:       []string{">= 5.0.0, < 6.0.0"},
					},
				},
				Labels:                 []string{"dependencies", "npm"},
				Assignees:              []string{"@reviewer1"},
				Reviewers:              []string{"@team/reviewers"},
				Milestone:              5,
				OpenPullRequestsLimit:  10,
				RebaseStrategy:         "auto",
				VersioningStrategy:     "increase",
				Vendor:                 true,
				TargetBranch:           "develop",
				CommitMessage: &dependabot.CommitMessage{
					Prefix:             "chore",
					PrefixDevelopment:  "dev",
					Include:            "scope",
				},
				PullRequestBranchName: &dependabot.PullRequestBranchName{
					Separator: "/",
				},
				InsecureExternalCodeExecution: "allow",
				Groups: map[string]dependabot.Group{
					"production": {
						Patterns:        []string{"*"},
						DependencyType:  "production",
						UpdateTypes:     []string{"minor", "patch"},
						ExcludePatterns: []string{"test-*"},
						AppliesTo:       "version-updates",
					},
				},
			},
		},
		Registries: map[string]dependabot.Registry{
			"npm-private": {
				Type:         "npm-registry",
				URL:          "https://npm.example.com",
				Token:        "${{ secrets.NPM_TOKEN }}",
				ReplacesBase: true,
			},
		},
	}

	yaml, err := serialize.DependabotToYAML(d)
	if err != nil {
		t.Fatalf("DependabotToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	expectedFields := []string{
		"enable-beta-ecosystems: true",
		"package-ecosystem: npm",
		"schedule:",
		"day: monday",
		"time: \"10:00\"",
		"timezone: America/New_York",
		"allow:",
		"ignore:",
		"labels:",
		"assignees:",
		"reviewers:",
		"milestone: 5",
		"open-pull-requests-limit: 10",
		"rebase-strategy: auto",
		"versioning-strategy: increase",
		"vendor: true",
		"target-branch: develop",
		"commit-message:",
		"prefix: chore",
		"pull-request-branch-name:",
		"separator: /",
		"groups:",
		"registries:",
	}

	for _, field := range expectedFields {
		if !strings.Contains(yamlStr, field) {
			t.Errorf("expected field %q, got:\n%s", field, yamlStr)
		}
	}
}

// TestDependabotMultipleUpdates tests multiple update configurations.
func TestDependabotMultipleUpdates(t *testing.T) {
	d := &dependabot.Dependabot{
		Version: 2,
		Updates: []dependabot.Update{
			{
				PackageEcosystem: "go",
				Directory:        "/",
				Schedule: dependabot.Schedule{
					Interval: "daily",
				},
			},
			{
				PackageEcosystem: "docker",
				Directory:        "/",
				Schedule: dependabot.Schedule{
					Interval: "weekly",
				},
			},
			{
				PackageEcosystem: "github-actions",
				Directory:        "/",
				Schedule: dependabot.Schedule{
					Interval: "weekly",
				},
			},
		},
	}

	yaml, err := serialize.DependabotToYAML(d)
	if err != nil {
		t.Fatalf("DependabotToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	if !strings.Contains(yamlStr, "package-ecosystem: go") {
		t.Errorf("expected go ecosystem, got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "package-ecosystem: docker") {
		t.Errorf("expected docker ecosystem, got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "package-ecosystem: github-actions") {
		t.Errorf("expected github-actions ecosystem, got:\n%s", yamlStr)
	}
}

// ===== Issue Template Tests =====

// TestBasicIssueTemplate tests basic issue template serialization.
func TestBasicIssueTemplate(t *testing.T) {
	tmpl := &templates.IssueTemplate{
		Name:        "Bug Report",
		Description: "File a bug report",
		Body: []templates.FormElement{
			templates.Markdown{
				Value: "## Bug Report\nPlease describe the bug.",
			},
			templates.Input{
				Label:    "Summary",
				Required: true,
			},
		},
	}

	yaml, err := serialize.IssueTemplateToYAML(tmpl)
	if err != nil {
		t.Fatalf("IssueTemplateToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	if !strings.Contains(yamlStr, "name: Bug Report") {
		t.Errorf("expected name, got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "description: File a bug report") {
		t.Errorf("expected description, got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "type: markdown") {
		t.Errorf("expected markdown type, got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "type: input") {
		t.Errorf("expected input type, got:\n%s", yamlStr)
	}
}

// TestIssueTemplateAllElements tests all form element types.
func TestIssueTemplateAllElements(t *testing.T) {
	tmpl := &templates.IssueTemplate{
		Name:        "Complete Form",
		Description: "Test all form elements",
		Title:       "New Issue",
		Labels:      []string{"bug", "needs-triage"},
		Assignees:   []string{"@maintainer"},
		Projects:    []string{"org/project/1"},
		Body: []templates.FormElement{
			templates.Markdown{
				ID:    "intro",
				Value: "## Welcome",
			},
			templates.Input{
				ID:          "summary",
				Label:       "Summary",
				Description: "Brief description",
				Placeholder: "Enter summary...",
				Value:       "Default value",
				Required:    true,
			},
			templates.Textarea{
				ID:          "details",
				Label:       "Details",
				Description: "Detailed description",
				Placeholder: "Enter details...",
				Value:       "Default details",
				Render:      "markdown",
				Required:    true,
			},
			templates.Dropdown{
				ID:          "severity",
				Label:       "Severity",
				Description: "Select severity",
				Options:     []string{"Low", "Medium", "High", "Critical"},
				Multiple:    false,
				Default:     2,
				Required:    true,
			},
			templates.Checkboxes{
				ID:          "checklist",
				Label:       "Pre-submission checklist",
				Description: "Please check all boxes",
				Options: []templates.CheckboxOption{
					{Label: "I have searched existing issues", Required: true},
					{Label: "I have read the documentation", Required: true},
					{Label: "I agree to the code of conduct", Required: false},
				},
			},
		},
	}

	yaml, err := serialize.IssueTemplateToYAML(tmpl)
	if err != nil {
		t.Fatalf("IssueTemplateToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	expectedFields := []string{
		"name: Complete Form",
		"description: Test all form elements",
		"title: New Issue",
		"labels:",
		"assignees:",
		"projects:",
		"type: markdown",
		"type: input",
		"type: textarea",
		"type: dropdown",
		"type: checkboxes",
		"render: markdown",
		"placeholder:",
		"options:",
		"default: 2",
		"validations:",
		"required: true",
	}

	for _, field := range expectedFields {
		if !strings.Contains(yamlStr, field) {
			t.Errorf("expected field %q, got:\n%s", field, yamlStr)
		}
	}
}

// TestIssueTemplateMinimal tests minimal issue template.
func TestIssueTemplateMinimal(t *testing.T) {
	tmpl := &templates.IssueTemplate{
		Name:        "Simple",
		Description: "Simple template",
		Body: []templates.FormElement{
			templates.Input{
				Label: "Title",
			},
		},
	}

	yaml, err := serialize.IssueTemplateToYAML(tmpl)
	if err != nil {
		t.Fatalf("IssueTemplateToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	// Should not have optional fields
	if strings.Contains(yamlStr, "title:") {
		t.Errorf("should not have title field, got:\n%s", yamlStr)
	}
	if strings.Contains(yamlStr, "labels:") {
		t.Errorf("should not have labels field, got:\n%s", yamlStr)
	}
}

// ===== Discussion Template Tests =====

// TestBasicDiscussionTemplate tests basic discussion template serialization.
func TestBasicDiscussionTemplate(t *testing.T) {
	tmpl := &templates.DiscussionTemplate{
		Title:       "General Discussion",
		Description: "Start a general discussion",
		Body: []templates.FormElement{
			templates.Markdown{
				Value: "## Discussion Topic",
			},
			templates.Textarea{
				Label:    "Description",
				Required: true,
			},
		},
	}

	yaml, err := serialize.DiscussionTemplateToYAML(tmpl)
	if err != nil {
		t.Fatalf("DiscussionTemplateToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	if !strings.Contains(yamlStr, "title: General Discussion") {
		t.Errorf("expected title, got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "description: Start a general discussion") {
		t.Errorf("expected description, got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "type: markdown") {
		t.Errorf("expected markdown type, got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "type: textarea") {
		t.Errorf("expected textarea type, got:\n%s", yamlStr)
	}
}

// TestDiscussionTemplateWithLabels tests discussion template with labels.
func TestDiscussionTemplateWithLabels(t *testing.T) {
	tmpl := &templates.DiscussionTemplate{
		Title:       "Feature Request",
		Description: "Request a new feature",
		Labels:      []string{"enhancement", "discussion"},
		Body: []templates.FormElement{
			templates.Input{
				Label:    "Feature Name",
				Required: true,
			},
		},
	}

	yaml, err := serialize.DiscussionTemplateToYAML(tmpl)
	if err != nil {
		t.Fatalf("DiscussionTemplateToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	if !strings.Contains(yamlStr, "labels:") {
		t.Errorf("expected labels field, got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "enhancement") {
		t.Errorf("expected enhancement label, got:\n%s", yamlStr)
	}
}

// ===== CODEOWNERS Tests =====

// TestBasicCodeowners tests basic CODEOWNERS serialization.
func TestBasicCodeowners(t *testing.T) {
	owners := &codeowners.Owners{
		Rules: []codeowners.Rule{
			{
				Pattern: "*",
				Owners:  []string{"@org/default-team"},
			},
			{
				Pattern: "*.go",
				Owners:  []string{"@go-team"},
			},
		},
	}

	text, err := serialize.CodeownersToText(owners)
	if err != nil {
		t.Fatalf("CodeownersToText failed: %v", err)
	}

	textStr := string(text)

	if !strings.Contains(textStr, "* @org/default-team") {
		t.Errorf("expected default rule, got:\n%s", textStr)
	}
	if !strings.Contains(textStr, "*.go @go-team") {
		t.Errorf("expected Go rule, got:\n%s", textStr)
	}
	if !strings.Contains(textStr, "# CODEOWNERS file generated by wetwire-github") {
		t.Errorf("expected header comment, got:\n%s", textStr)
	}
}

// TestCodeownersWithComments tests CODEOWNERS with comments.
func TestCodeownersWithComments(t *testing.T) {
	owners := &codeowners.Owners{
		Rules: []codeowners.Rule{
			{
				Pattern: "*.go",
				Owners:  []string{"@go-team"},
				Comment: "Go source files",
			},
			{
				Pattern: "*.js",
				Owners:  []string{"@js-team"},
				Comment: "JavaScript files",
			},
		},
	}

	text, err := serialize.CodeownersToText(owners)
	if err != nil {
		t.Fatalf("CodeownersToText failed: %v", err)
	}

	textStr := string(text)

	if !strings.Contains(textStr, "# Go source files") {
		t.Errorf("expected Go comment, got:\n%s", textStr)
	}
	if !strings.Contains(textStr, "# JavaScript files") {
		t.Errorf("expected JS comment, got:\n%s", textStr)
	}
}

// TestCodeownersMultipleOwners tests multiple owners per pattern.
func TestCodeownersMultipleOwners(t *testing.T) {
	owners := &codeowners.Owners{
		Rules: []codeowners.Rule{
			{
				Pattern: "/docs/",
				Owners:  []string{"@docs-team", "@user1", "@user2"},
			},
		},
	}

	text, err := serialize.CodeownersToText(owners)
	if err != nil {
		t.Fatalf("CodeownersToText failed: %v", err)
	}

	textStr := string(text)

	if !strings.Contains(textStr, "/docs/ @docs-team @user1 @user2") {
		t.Errorf("expected multiple owners, got:\n%s", textStr)
	}
}

// TestCodeownersEmptyOwners tests pattern without owners.
func TestCodeownersEmptyOwners(t *testing.T) {
	owners := &codeowners.Owners{
		Rules: []codeowners.Rule{
			{
				Pattern: "*.tmp",
				Owners:  []string{},
				Comment: "Temporary files",
			},
		},
	}

	text, err := serialize.CodeownersToText(owners)
	if err != nil {
		t.Fatalf("CodeownersToText failed: %v", err)
	}

	textStr := string(text)

	if !strings.Contains(textStr, "# Temporary files") {
		t.Errorf("expected comment, got:\n%s", textStr)
	}
	if !strings.Contains(textStr, "*.tmp") {
		t.Errorf("expected pattern, got:\n%s", textStr)
	}
}

// TestCodeownersRulesToText tests ExtractedCodeownersRule serialization.
func TestCodeownersRulesToText(t *testing.T) {
	rules := []serialize.ExtractedCodeownersRule{
		{
			Pattern: "*.go",
			Owners:  []string{"@go-team"},
			Comment: "Go files",
		},
		{
			Pattern: "/src/**",
			Owners:  []string{"@src-team"},
		},
	}

	text, err := serialize.CodeownersRulesToText(rules)
	if err != nil {
		t.Fatalf("CodeownersRulesToText failed: %v", err)
	}

	textStr := string(text)

	if !strings.Contains(textStr, "*.go @go-team") {
		t.Errorf("expected Go rule, got:\n%s", textStr)
	}
	if !strings.Contains(textStr, "/src/** @src-team") {
		t.Errorf("expected src rule, got:\n%s", textStr)
	}
	if !strings.Contains(textStr, "# Go files") {
		t.Errorf("expected comment, got:\n%s", textStr)
	}
}

// ===== Additional Edge Case Tests =====

// TestEmptyWorkflow tests empty workflow handling.
func TestEmptyWorkflow(t *testing.T) {
	w := &workflow.Workflow{
		Name: "Empty",
		On:   workflow.Triggers{},
		Jobs: map[string]workflow.Job{},
	}

	yaml, err := serialize.ToYAML(w)
	if err != nil {
		t.Fatalf("ToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	if !strings.Contains(yamlStr, "name: Empty") {
		t.Errorf("expected name, got:\n%s", yamlStr)
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
					mockAction,  // This should be converted via anyStepToMap
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
				Needs:  []any{"build"},  // String reference
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

// TestPullRequestTriggerAllFields tests all pull_request trigger fields.
func TestPullRequestTriggerAllFields(t *testing.T) {
	w := &workflow.Workflow{
		Name: "Test",
		On: workflow.Triggers{
			PullRequest: &workflow.PullRequestTrigger{
				Types:          []string{"opened", "synchronize", "reopened"},
				Branches:       []string{"main", "develop"},
				BranchesIgnore: []string{"experimental"},
				Paths:          []string{"src/**", "*.go"},
				PathsIgnore:    []string{"docs/**"},
			},
		},
	}

	yaml, err := serialize.ToYAML(w)
	if err != nil {
		t.Fatalf("ToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	expectedFields := []string{
		"pull_request:",
		"types:",
		"opened",
		"branches:",
		"branches-ignore:",
		"paths:",
		"paths-ignore:",
	}

	for _, field := range expectedFields {
		if !strings.Contains(yamlStr, field) {
			t.Errorf("expected field %q, got:\n%s", field, yamlStr)
		}
	}
}

// TestAllPermissionsFields tests all permission scopes.
func TestAllPermissionsFields(t *testing.T) {
	w := &workflow.Workflow{
		Name: "Test",
		On: workflow.Triggers{
			Push: &workflow.PushTrigger{},
		},
		Permissions: &workflow.Permissions{
			Actions:            "read",
			Checks:             "write",
			Contents:           "write",
			Deployments:        "write",
			Discussions:        "read",
			IDToken:            "write",
			Issues:             "write",
			Packages:           "write",
			Pages:              "write",
			PullRequests:       "write",
			RepositoryProjects: "write",
			SecurityEvents:     "write",
			Statuses:           "write",
		},
	}

	yaml, err := serialize.ToYAML(w)
	if err != nil {
		t.Fatalf("ToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	expectedPermissions := []string{
		"actions:", "checks:", "contents:", "deployments:",
		"discussions:", "id-token:", "issues:", "packages:",
		"pages:", "pull-requests:", "repository-projects:",
		"security-events:", "statuses:",
	}

	for _, perm := range expectedPermissions {
		if !strings.Contains(yamlStr, perm) {
			t.Errorf("expected permission %q, got:\n%s", perm, yamlStr)
		}
	}
}

// TestDependabotRegistryAllFields tests all registry fields.
func TestDependabotRegistryAllFields(t *testing.T) {
	d := &dependabot.Dependabot{
		Version: 2,
		Updates: []dependabot.Update{
			{
				PackageEcosystem: "npm",
				Directory:        "/",
				Schedule: dependabot.Schedule{
					Interval: "daily",
				},
			},
		},
		Registries: map[string]dependabot.Registry{
			"npm-registry": {
				Type:         "npm-registry",
				URL:          "https://npm.example.com",
				Username:     "user",
				Password:     "${{ secrets.NPM_PASSWORD }}",
				Token:        "${{ secrets.NPM_TOKEN }}",
				Key:          "${{ secrets.NPM_KEY }}",
				Organization: "my-org",
				ReplacesBase: true,
			},
		},
	}

	yaml, err := serialize.DependabotToYAML(d)
	if err != nil {
		t.Fatalf("DependabotToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	expectedFields := []string{
		"type: npm-registry",
		"url: https://npm.example.com",
		"username: user",
		"password:",
		"token:",
		"key:",
		"organization: my-org",
		"replaces-base: true",
	}

	for _, field := range expectedFields {
		if !strings.Contains(yamlStr, field) {
			t.Errorf("expected field %q, got:\n%s", field, yamlStr)
		}
	}
}

// TestDependabotIgnoreAllFields tests all ignore fields.
func TestDependabotIgnoreAllFields(t *testing.T) {
	d := &dependabot.Dependabot{
		Version: 2,
		Updates: []dependabot.Update{
			{
				PackageEcosystem: "npm",
				Directory:        "/",
				Schedule: dependabot.Schedule{
					Interval: "daily",
				},
				Ignore: []dependabot.Ignore{
					{
						DependencyName: "webpack",
						Versions:       []string{"5.x", "6.x"},
						UpdateTypes:    []string{"version-update:semver-major"},
					},
				},
			},
		},
	}

	yaml, err := serialize.DependabotToYAML(d)
	if err != nil {
		t.Fatalf("DependabotToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	if !strings.Contains(yamlStr, "dependency-name: webpack") {
		t.Errorf("expected dependency-name, got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "versions:") {
		t.Errorf("expected versions, got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "update-types:") {
		t.Errorf("expected update-types, got:\n%s", yamlStr)
	}
}

// TestWorkflowRunTriggerAllFields tests all workflow_run trigger fields.
func TestWorkflowRunTriggerAllFields(t *testing.T) {
	w := &workflow.Workflow{
		Name: "Test",
		On: workflow.Triggers{
			WorkflowRun: &workflow.WorkflowRunTrigger{
				Workflows: []string{"CI", "Build"},
				Types:     []string{"completed"},
				Branches:  []string{"main", "develop"},
			},
		},
	}

	yaml, err := serialize.ToYAML(w)
	if err != nil {
		t.Fatalf("ToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	expectedFields := []string{
		"workflow_run:",
		"workflows:",
		"types:",
		"branches:",
	}

	for _, field := range expectedFields {
		if !strings.Contains(yamlStr, field) {
			t.Errorf("expected field %q, got:\n%s", field, yamlStr)
		}
	}
}

// TestMoreTriggerTypes tests additional trigger types with types field.
func TestMoreTriggerTypes(t *testing.T) {
	w := &workflow.Workflow{
		Name: "Test",
		On: workflow.Triggers{
			Label:                      &workflow.LabelTrigger{Types: []string{"created"}},
			Milestone:                  &workflow.MilestoneTrigger{Types: []string{"opened"}},
			Project:                    &workflow.ProjectTrigger{Types: []string{"created"}},
			ProjectCard:                &workflow.ProjectCardTrigger{Types: []string{"moved"}},
			ProjectColumn:              &workflow.ProjectColumnTrigger{Types: []string{"created"}},
			PullRequestReview:          &workflow.PullRequestReviewTrigger{Types: []string{"submitted"}},
			PullRequestReviewComment:   &workflow.PullRequestReviewCommentTrigger{Types: []string{"created"}},
			Watch:                      &workflow.WatchTrigger{Types: []string{"started"}},
			CheckRun:                   &workflow.CheckRunTrigger{Types: []string{"completed"}},
			CheckSuite:                 &workflow.CheckSuiteTrigger{Types: []string{"completed"}},
			Discussion:                 &workflow.DiscussionTrigger{Types: []string{"created"}},
			DiscussionComment:          &workflow.DiscussionCommentTrigger{Types: []string{"created"}},
			MergeGroup:                 &workflow.MergeGroupTrigger{Types: []string{"checks_requested"}},
		},
	}

	yaml, err := serialize.ToYAML(w)
	if err != nil {
		t.Fatalf("ToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	expectedTriggers := []string{
		"label:", "milestone:", "project:", "project_card:", "project_column:",
		"pull_request_review:", "pull_request_review_comment:", "watch:",
		"check_run:", "check_suite:", "discussion:", "discussion_comment:",
		"merge_group:",
	}

	for _, trigger := range expectedTriggers {
		if !strings.Contains(yamlStr, trigger) {
			t.Errorf("expected trigger %q, got:\n%s", trigger, yamlStr)
		}
	}
}

// TestSimpleTriggers tests triggers without configuration.
func TestSimpleTriggers(t *testing.T) {
	w := &workflow.Workflow{
		Name: "Test",
		On: workflow.Triggers{
			Fork:      &workflow.ForkTrigger{},
			Gollum:    &workflow.GollumTrigger{},
			Public:    &workflow.PublicTrigger{},
			PageBuild: &workflow.PageBuildTrigger{},
			Status:    &workflow.StatusTrigger{},
		},
	}

	yaml, err := serialize.ToYAML(w)
	if err != nil {
		t.Fatalf("ToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	expectedTriggers := []string{
		"fork:", "gollum:", "public:", "page_build:", "status:",
	}

	for _, trigger := range expectedTriggers {
		if !strings.Contains(yamlStr, trigger) {
			t.Errorf("expected trigger %q, got:\n%s", trigger, yamlStr)
		}
	}
}

// TestWorkflowWithDefaults tests workflow and job defaults.
func TestWorkflowWithDefaults(t *testing.T) {
	w := &workflow.Workflow{
		Name: "Test",
		On: workflow.Triggers{
			Push: &workflow.PushTrigger{},
		},
		Defaults: &workflow.WorkflowDefaults{
			Run: &workflow.RunDefaults{
				Shell:            "bash",
				WorkingDirectory: "/tmp",
			},
		},
		Jobs: map[string]workflow.Job{
			"test": {
				RunsOn: "ubuntu-latest",
				Defaults: &workflow.JobDefaults{
					Run: &workflow.RunDefaults{
						Shell: "sh",
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

	if !strings.Contains(yamlStr, "defaults:") {
		t.Errorf("expected defaults field, got:\n%s", yamlStr)
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

// TestPullRequestTargetTrigger tests pull_request_target trigger.
func TestPullRequestTargetTrigger(t *testing.T) {
	w := &workflow.Workflow{
		Name: "Test",
		On: workflow.Triggers{
			PullRequestTarget: &workflow.PullRequestTargetTrigger{
				Types:          []string{"opened"},
				Branches:       []string{"main"},
				BranchesIgnore: []string{"dev"},
				Paths:          []string{"src/**"},
				PathsIgnore:    []string{"docs/**"},
			},
		},
	}

	yaml, err := serialize.ToYAML(w)
	if err != nil {
		t.Fatalf("ToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	if !strings.Contains(yamlStr, "pull_request_target:") {
		t.Errorf("expected pull_request_target trigger, got:\n%s", yamlStr)
	}
}
