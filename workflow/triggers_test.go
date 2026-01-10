package workflow_test

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestPushTrigger(t *testing.T) {
	trigger := workflow.PushTrigger{
		Branches:       workflow.List("main", "develop"),
		BranchesIgnore: workflow.List("feature/**"),
		Tags:           workflow.List("v*"),
		TagsIgnore:     workflow.List("v*-beta"),
		Paths:          workflow.List("src/**"),
		PathsIgnore:    workflow.List("docs/**"),
	}

	if len(trigger.Branches) != 2 {
		t.Errorf("expected 2 branches, got %d", len(trigger.Branches))
	}

	if trigger.Branches[0] != "main" {
		t.Errorf("expected main, got %s", trigger.Branches[0])
	}

	if len(trigger.Tags) != 1 {
		t.Errorf("expected 1 tag, got %d", len(trigger.Tags))
	}

	if len(trigger.PathsIgnore) != 1 {
		t.Errorf("expected 1 paths-ignore, got %d", len(trigger.PathsIgnore))
	}
}

func TestPullRequestTrigger(t *testing.T) {
	trigger := workflow.PullRequestTrigger{
		Types:          workflow.List("opened", "synchronize", "reopened"),
		Branches:       workflow.List("main"),
		BranchesIgnore: workflow.List("wip/**"),
		Paths:          workflow.List("src/**", "tests/**"),
		PathsIgnore:    workflow.List("*.md"),
	}

	if len(trigger.Types) != 3 {
		t.Errorf("expected 3 types, got %d", len(trigger.Types))
	}

	if trigger.Types[0] != "opened" {
		t.Errorf("expected opened, got %s", trigger.Types[0])
	}

	if len(trigger.Branches) != 1 {
		t.Errorf("expected 1 branch, got %d", len(trigger.Branches))
	}

	if len(trigger.Paths) != 2 {
		t.Errorf("expected 2 paths, got %d", len(trigger.Paths))
	}
}

func TestPullRequestTargetTrigger(t *testing.T) {
	trigger := workflow.PullRequestTargetTrigger{
		Types:          workflow.List("opened"),
		Branches:       workflow.List("main"),
		BranchesIgnore: workflow.List("staging"),
		Paths:          workflow.List("app/**"),
		PathsIgnore:    workflow.List("tests/**"),
	}

	if len(trigger.Types) != 1 {
		t.Errorf("expected 1 type, got %d", len(trigger.Types))
	}

	if trigger.Branches[0] != "main" {
		t.Errorf("expected main, got %s", trigger.Branches[0])
	}
}

func TestScheduleTrigger(t *testing.T) {
	trigger := workflow.ScheduleTrigger{
		Cron: "0 0 * * *",
	}

	if trigger.Cron != "0 0 * * *" {
		t.Errorf("expected '0 0 * * *', got %s", trigger.Cron)
	}
}

func TestWorkflowDispatchTrigger(t *testing.T) {
	trigger := workflow.WorkflowDispatchTrigger{
		Inputs: map[string]workflow.WorkflowInput{
			"environment": {
				Description: "Deployment environment",
				Required:    true,
				Default:     "staging",
				Type:        "string",
			},
			"debug": {
				Description: "Enable debug mode",
				Required:    false,
				Default:     false,
				Type:        "boolean",
			},
			"region": {
				Description: "AWS region",
				Type:        "choice",
				Options:     workflow.List("us-east-1", "us-west-2", "eu-west-1"),
			},
		},
	}

	if len(trigger.Inputs) != 3 {
		t.Errorf("expected 3 inputs, got %d", len(trigger.Inputs))
	}

	env := trigger.Inputs["environment"]
	if env.Type != "string" {
		t.Errorf("expected string type, got %s", env.Type)
	}

	if !env.Required {
		t.Error("expected environment to be required")
	}

	debug := trigger.Inputs["debug"]
	if debug.Type != "boolean" {
		t.Errorf("expected boolean type, got %s", debug.Type)
	}

	region := trigger.Inputs["region"]
	if region.Type != "choice" {
		t.Errorf("expected choice type, got %s", region.Type)
	}

	if len(region.Options) != 3 {
		t.Errorf("expected 3 options, got %d", len(region.Options))
	}
}

func TestWorkflowRunTrigger(t *testing.T) {
	trigger := workflow.WorkflowRunTrigger{
		Workflows: workflow.List("Build", "Test"),
		Types:     workflow.List("completed"),
		Branches:  workflow.List("main", "develop"),
	}

	if len(trigger.Workflows) != 2 {
		t.Errorf("expected 2 workflows, got %d", len(trigger.Workflows))
	}

	if trigger.Types[0] != "completed" {
		t.Errorf("expected completed, got %s", trigger.Types[0])
	}

	if len(trigger.Branches) != 2 {
		t.Errorf("expected 2 branches, got %d", len(trigger.Branches))
	}
}

func TestRepositoryDispatchTrigger(t *testing.T) {
	trigger := workflow.RepositoryDispatchTrigger{
		Types: workflow.List("custom-event", "another-event"),
	}

	if len(trigger.Types) != 2 {
		t.Errorf("expected 2 types, got %d", len(trigger.Types))
	}

	if trigger.Types[0] != "custom-event" {
		t.Errorf("expected custom-event, got %s", trigger.Types[0])
	}
}

func TestEmptyTriggers(t *testing.T) {
	tests := []struct {
		name    string
		trigger any
	}{
		{"CreateTrigger", workflow.CreateTrigger{}},
		{"DeleteTrigger", workflow.DeleteTrigger{}},
		{"ForkTrigger", workflow.ForkTrigger{}},
		{"GollumTrigger", workflow.GollumTrigger{}},
		{"PageBuildTrigger", workflow.PageBuildTrigger{}},
		{"PublicTrigger", workflow.PublicTrigger{}},
		{"StatusTrigger", workflow.StatusTrigger{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// These triggers have no fields, just verify they can be created
			if tt.trigger == nil {
				t.Error("expected trigger to be created")
			}
		})
	}
}

func TestIssueCommentTrigger(t *testing.T) {
	trigger := workflow.IssueCommentTrigger{
		Types: workflow.List("created", "edited", "deleted"),
	}

	if len(trigger.Types) != 3 {
		t.Errorf("expected 3 types, got %d", len(trigger.Types))
	}
}

func TestIssuesTrigger(t *testing.T) {
	trigger := workflow.IssuesTrigger{
		Types: workflow.List("opened", "closed", "reopened"),
	}

	if len(trigger.Types) != 3 {
		t.Errorf("expected 3 types, got %d", len(trigger.Types))
	}
}

func TestLabelTrigger(t *testing.T) {
	trigger := workflow.LabelTrigger{
		Types: workflow.List("created", "edited", "deleted"),
	}

	if len(trigger.Types) != 3 {
		t.Errorf("expected 3 types, got %d", len(trigger.Types))
	}
}

func TestMilestoneTrigger(t *testing.T) {
	trigger := workflow.MilestoneTrigger{
		Types: workflow.List("created", "closed", "opened"),
	}

	if len(trigger.Types) != 3 {
		t.Errorf("expected 3 types, got %d", len(trigger.Types))
	}
}

func TestProjectTrigger(t *testing.T) {
	trigger := workflow.ProjectTrigger{
		Types: workflow.List("created", "updated", "closed"),
	}

	if len(trigger.Types) != 3 {
		t.Errorf("expected 3 types, got %d", len(trigger.Types))
	}
}

func TestProjectCardTrigger(t *testing.T) {
	trigger := workflow.ProjectCardTrigger{
		Types: workflow.List("created", "moved", "deleted"),
	}

	if len(trigger.Types) != 3 {
		t.Errorf("expected 3 types, got %d", len(trigger.Types))
	}
}

func TestProjectColumnTrigger(t *testing.T) {
	trigger := workflow.ProjectColumnTrigger{
		Types: workflow.List("created", "updated", "moved"),
	}

	if len(trigger.Types) != 3 {
		t.Errorf("expected 3 types, got %d", len(trigger.Types))
	}
}

func TestPullRequestReviewTrigger(t *testing.T) {
	trigger := workflow.PullRequestReviewTrigger{
		Types: workflow.List("submitted", "edited", "dismissed"),
	}

	if len(trigger.Types) != 3 {
		t.Errorf("expected 3 types, got %d", len(trigger.Types))
	}
}

func TestPullRequestReviewCommentTrigger(t *testing.T) {
	trigger := workflow.PullRequestReviewCommentTrigger{
		Types: workflow.List("created", "edited", "deleted"),
	}

	if len(trigger.Types) != 3 {
		t.Errorf("expected 3 types, got %d", len(trigger.Types))
	}
}

func TestReleaseTrigger(t *testing.T) {
	trigger := workflow.ReleaseTrigger{
		Types: workflow.List("published", "created", "released"),
	}

	if len(trigger.Types) != 3 {
		t.Errorf("expected 3 types, got %d", len(trigger.Types))
	}
}

func TestWatchTrigger(t *testing.T) {
	trigger := workflow.WatchTrigger{
		Types: workflow.List("started"),
	}

	if len(trigger.Types) != 1 {
		t.Errorf("expected 1 type, got %d", len(trigger.Types))
	}

	if trigger.Types[0] != "started" {
		t.Errorf("expected started, got %s", trigger.Types[0])
	}
}

func TestCheckRunTrigger(t *testing.T) {
	trigger := workflow.CheckRunTrigger{
		Types: workflow.List("created", "rerequested", "completed", "requested_action"),
	}

	if len(trigger.Types) != 4 {
		t.Errorf("expected 4 types, got %d", len(trigger.Types))
	}
}

func TestCheckSuiteTrigger(t *testing.T) {
	trigger := workflow.CheckSuiteTrigger{
		Types: workflow.List("completed", "requested", "rerequested"),
	}

	if len(trigger.Types) != 3 {
		t.Errorf("expected 3 types, got %d", len(trigger.Types))
	}
}

func TestDiscussionTrigger(t *testing.T) {
	trigger := workflow.DiscussionTrigger{
		Types: workflow.List("created", "edited", "deleted"),
	}

	if len(trigger.Types) != 3 {
		t.Errorf("expected 3 types, got %d", len(trigger.Types))
	}
}

func TestDiscussionCommentTrigger(t *testing.T) {
	trigger := workflow.DiscussionCommentTrigger{
		Types: workflow.List("created", "edited", "deleted"),
	}

	if len(trigger.Types) != 3 {
		t.Errorf("expected 3 types, got %d", len(trigger.Types))
	}
}

func TestMergeGroupTrigger(t *testing.T) {
	trigger := workflow.MergeGroupTrigger{
		Types: workflow.List("checks_requested"),
	}

	if len(trigger.Types) != 1 {
		t.Errorf("expected 1 type, got %d", len(trigger.Types))
	}

	if trigger.Types[0] != "checks_requested" {
		t.Errorf("expected checks_requested, got %s", trigger.Types[0])
	}
}

func TestTriggersComposite(t *testing.T) {
	triggers := workflow.Triggers{
		Push: &workflow.PushTrigger{
			Branches: workflow.List("main"),
		},
		PullRequest: &workflow.PullRequestTrigger{
			Types: workflow.List("opened", "synchronize"),
		},
		Schedule: []workflow.ScheduleTrigger{
			{Cron: "0 0 * * *"},
			{Cron: "0 12 * * *"},
		},
		WorkflowDispatch: &workflow.WorkflowDispatchTrigger{},
	}

	if triggers.Push == nil {
		t.Error("expected Push trigger to be set")
	}

	if triggers.PullRequest == nil {
		t.Error("expected PullRequest trigger to be set")
	}

	if len(triggers.Schedule) != 2 {
		t.Errorf("expected 2 schedules, got %d", len(triggers.Schedule))
	}

	if triggers.WorkflowDispatch == nil {
		t.Error("expected WorkflowDispatch trigger to be set")
	}
}
