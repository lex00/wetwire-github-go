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
			Label:                    &workflow.LabelTrigger{Types: []string{"created"}},
			Milestone:                &workflow.MilestoneTrigger{Types: []string{"opened"}},
			Project:                  &workflow.ProjectTrigger{Types: []string{"created"}},
			ProjectCard:              &workflow.ProjectCardTrigger{Types: []string{"moved"}},
			ProjectColumn:            &workflow.ProjectColumnTrigger{Types: []string{"created"}},
			PullRequestReview:        &workflow.PullRequestReviewTrigger{Types: []string{"submitted"}},
			PullRequestReviewComment: &workflow.PullRequestReviewCommentTrigger{Types: []string{"created"}},
			Watch:                    &workflow.WatchTrigger{Types: []string{"started"}},
			CheckRun:                 &workflow.CheckRunTrigger{Types: []string{"completed"}},
			CheckSuite:               &workflow.CheckSuiteTrigger{Types: []string{"completed"}},
			Discussion:               &workflow.DiscussionTrigger{Types: []string{"created"}},
			DiscussionComment:        &workflow.DiscussionCommentTrigger{Types: []string{"created"}},
			MergeGroup:               &workflow.MergeGroupTrigger{Types: []string{"checks_requested"}},
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
