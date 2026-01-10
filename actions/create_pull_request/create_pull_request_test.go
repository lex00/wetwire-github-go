package create_pull_request

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestCreatePullRequest_Action(t *testing.T) {
	a := CreatePullRequest{}
	if got := a.Action(); got != "peter-evans/create-pull-request@v6" {
		t.Errorf("Action() = %q, want %q", got, "peter-evans/create-pull-request@v6")
	}
}

func TestCreatePullRequest_Inputs_Empty(t *testing.T) {
	a := CreatePullRequest{}
	inputs := a.Inputs()

	if a.Action() != "peter-evans/create-pull-request@v6" {
		t.Errorf("Action() = %q, want %q", a.Action(), "peter-evans/create-pull-request@v6")
	}

	if len(inputs) != 0 {
		t.Errorf("Empty CreatePullRequest should have no inputs, got %v", inputs)
	}
}

func TestCreatePullRequest_Inputs_Basic(t *testing.T) {
	a := CreatePullRequest{
		CommitMessage: "Update dependencies",
		Title:         "chore: update dependencies",
		Body:          "This PR updates all dependencies to their latest versions",
		Branch:        "update-deps",
	}

	inputs := a.Inputs()

	if inputs["commit-message"] != "Update dependencies" {
		t.Errorf("commit-message = %v, want %q", inputs["commit-message"], "Update dependencies")
	}
	if inputs["title"] != "chore: update dependencies" {
		t.Errorf("title = %v, want %q", inputs["title"], "chore: update dependencies")
	}
	if inputs["body"] != "This PR updates all dependencies to their latest versions" {
		t.Errorf("body = %v, want expected value", inputs["body"])
	}
	if inputs["branch"] != "update-deps" {
		t.Errorf("branch = %v, want %q", inputs["branch"], "update-deps")
	}
}

func TestCreatePullRequest_Inputs_WithToken(t *testing.T) {
	a := CreatePullRequest{
		Token:         "${{ secrets.PAT }}",
		CommitMessage: "Update workflow",
		Title:         "Update workflow",
	}

	inputs := a.Inputs()

	if inputs["token"] != "${{ secrets.PAT }}" {
		t.Errorf("token = %v, want %q", inputs["token"], "${{ secrets.PAT }}")
	}
}

func TestCreatePullRequest_Inputs_WithPath(t *testing.T) {
	a := CreatePullRequest{
		Path:          "submodule/",
		CommitMessage: "Update submodule",
		Title:         "Update submodule",
	}

	inputs := a.Inputs()

	if inputs["path"] != "submodule/" {
		t.Errorf("path = %v, want %q", inputs["path"], "submodule/")
	}
}

func TestCreatePullRequest_Inputs_WithAddPaths(t *testing.T) {
	a := CreatePullRequest{
		AddPaths:      "src/\ntests/",
		CommitMessage: "Update source and tests",
		Title:         "Update source and tests",
	}

	inputs := a.Inputs()

	if inputs["add-paths"] != "src/\ntests/" {
		t.Errorf("add-paths = %v, want expected value", inputs["add-paths"])
	}
}

func TestCreatePullRequest_Inputs_WithCommitter(t *testing.T) {
	a := CreatePullRequest{
		Committer:     "GitHub Actions Bot <bot@github.com>",
		CommitMessage: "Auto-update",
		Title:         "Auto-update",
	}

	inputs := a.Inputs()

	if inputs["committer"] != "GitHub Actions Bot <bot@github.com>" {
		t.Errorf("committer = %v, want expected value", inputs["committer"])
	}
}

func TestCreatePullRequest_Inputs_WithAuthor(t *testing.T) {
	a := CreatePullRequest{
		Author:        "Custom Author <author@example.com>",
		CommitMessage: "Custom commit",
		Title:         "Custom PR",
	}

	inputs := a.Inputs()

	if inputs["author"] != "Custom Author <author@example.com>" {
		t.Errorf("author = %v, want expected value", inputs["author"])
	}
}

func TestCreatePullRequest_Inputs_WithSignoff(t *testing.T) {
	a := CreatePullRequest{
		Signoff:       true,
		CommitMessage: "Signed commit",
		Title:         "Signed PR",
	}

	inputs := a.Inputs()

	if inputs["signoff"] != true {
		t.Errorf("signoff = %v, want %v", inputs["signoff"], true)
	}
}

func TestCreatePullRequest_Inputs_WithBranchSuffix(t *testing.T) {
	a := CreatePullRequest{
		Branch:        "auto-update",
		BranchSuffix:  "timestamp",
		CommitMessage: "Auto-update",
		Title:         "Auto-update PR",
	}

	inputs := a.Inputs()

	if inputs["branch"] != "auto-update" {
		t.Errorf("branch = %v, want %q", inputs["branch"], "auto-update")
	}
	if inputs["branch-suffix"] != "timestamp" {
		t.Errorf("branch-suffix = %v, want %q", inputs["branch-suffix"], "timestamp")
	}
}

func TestCreatePullRequest_Inputs_WithDeleteBranch(t *testing.T) {
	a := CreatePullRequest{
		DeleteBranch:  true,
		CommitMessage: "Update",
		Title:         "Update PR",
	}

	inputs := a.Inputs()

	if inputs["delete-branch"] != true {
		t.Errorf("delete-branch = %v, want %v", inputs["delete-branch"], true)
	}
}

func TestCreatePullRequest_Inputs_WithBodyPath(t *testing.T) {
	a := CreatePullRequest{
		BodyPath:      ".github/pull_request_template.md",
		CommitMessage: "Update",
		Title:         "Update PR",
	}

	inputs := a.Inputs()

	if inputs["body-path"] != ".github/pull_request_template.md" {
		t.Errorf("body-path = %v, want %q", inputs["body-path"], ".github/pull_request_template.md")
	}
}

func TestCreatePullRequest_Inputs_WithLabels(t *testing.T) {
	a := CreatePullRequest{
		Labels:        "bug,dependencies",
		CommitMessage: "Fix dependency bug",
		Title:         "Fix dependency bug",
	}

	inputs := a.Inputs()

	if inputs["labels"] != "bug,dependencies" {
		t.Errorf("labels = %v, want %q", inputs["labels"], "bug,dependencies")
	}
}

func TestCreatePullRequest_Inputs_WithAssignees(t *testing.T) {
	a := CreatePullRequest{
		Assignees:     "user1,user2",
		CommitMessage: "Update",
		Title:         "Update PR",
	}

	inputs := a.Inputs()

	if inputs["assignees"] != "user1,user2" {
		t.Errorf("assignees = %v, want %q", inputs["assignees"], "user1,user2")
	}
}

func TestCreatePullRequest_Inputs_WithReviewers(t *testing.T) {
	a := CreatePullRequest{
		Reviewers:     "reviewer1\nreviewer2",
		CommitMessage: "Update",
		Title:         "Update PR",
	}

	inputs := a.Inputs()

	if inputs["reviewers"] != "reviewer1\nreviewer2" {
		t.Errorf("reviewers = %v, want expected value", inputs["reviewers"])
	}
}

func TestCreatePullRequest_Inputs_WithTeamReviewers(t *testing.T) {
	a := CreatePullRequest{
		TeamReviewers: "team1,team2",
		CommitMessage: "Update",
		Title:         "Update PR",
	}

	inputs := a.Inputs()

	if inputs["team-reviewers"] != "team1,team2" {
		t.Errorf("team-reviewers = %v, want %q", inputs["team-reviewers"], "team1,team2")
	}
}

func TestCreatePullRequest_Inputs_WithMilestone(t *testing.T) {
	a := CreatePullRequest{
		Milestone:     5,
		CommitMessage: "Update",
		Title:         "Update PR",
	}

	inputs := a.Inputs()

	if inputs["milestone"] != 5 {
		t.Errorf("milestone = %v, want %d", inputs["milestone"], 5)
	}
}

func TestCreatePullRequest_Inputs_WithMilestoneZero(t *testing.T) {
	a := CreatePullRequest{
		Milestone:     0,
		CommitMessage: "Update",
		Title:         "Update PR",
	}

	inputs := a.Inputs()

	if _, ok := inputs["milestone"]; ok {
		t.Error("milestone should not be in inputs when set to 0")
	}
}

func TestCreatePullRequest_Inputs_WithDraft(t *testing.T) {
	a := CreatePullRequest{
		Draft:         true,
		CommitMessage: "WIP: Update",
		Title:         "WIP: Update PR",
	}

	inputs := a.Inputs()

	if inputs["draft"] != true {
		t.Errorf("draft = %v, want %v", inputs["draft"], true)
	}
}

func TestCreatePullRequest_Inputs_AllFields(t *testing.T) {
	a := CreatePullRequest{
		Token:         "${{ secrets.PAT }}",
		Path:          "subdir/",
		AddPaths:      "src/\ntests/",
		CommitMessage: "feat: add new feature",
		Committer:     "Bot <bot@example.com>",
		Author:        "Author <author@example.com>",
		Signoff:       true,
		Branch:        "feature-branch",
		BranchSuffix:  "timestamp",
		DeleteBranch:  true,
		Title:         "Add new feature",
		Body:          "## Summary\n\nThis PR adds a new feature",
		BodyPath:      ".github/pr_template.md",
		Labels:        "feature,enhancement",
		Assignees:     "user1,user2",
		Reviewers:     "reviewer1,reviewer2",
		TeamReviewers: "team1,team2",
		Milestone:     3,
		Draft:         true,
	}

	inputs := a.Inputs()

	expected := map[string]any{
		"token":          "${{ secrets.PAT }}",
		"path":           "subdir/",
		"add-paths":      "src/\ntests/",
		"commit-message": "feat: add new feature",
		"committer":      "Bot <bot@example.com>",
		"author":         "Author <author@example.com>",
		"signoff":        true,
		"branch":         "feature-branch",
		"branch-suffix":  "timestamp",
		"delete-branch":  true,
		"title":          "Add new feature",
		"body":           "## Summary\n\nThis PR adds a new feature",
		"body-path":      ".github/pr_template.md",
		"labels":         "feature,enhancement",
		"assignees":      "user1,user2",
		"reviewers":      "reviewer1,reviewer2",
		"team-reviewers": "team1,team2",
		"milestone":      3,
		"draft":          true,
	}

	for key, want := range expected {
		if got := inputs[key]; got != want {
			t.Errorf("inputs[%q] = %v, want %v", key, got, want)
		}
	}

	if len(inputs) != len(expected) {
		t.Errorf("inputs length = %d, want %d", len(inputs), len(expected))
	}
}

func TestCreatePullRequest_Inputs_MultilineLabels(t *testing.T) {
	a := CreatePullRequest{
		Labels:        "bug\nurgent\nsecurity",
		CommitMessage: "security: fix vulnerability",
		Title:         "Security fix",
	}

	inputs := a.Inputs()

	if inputs["labels"] != "bug\nurgent\nsecurity" {
		t.Errorf("labels = %v, want expected multiline value", inputs["labels"])
	}
}

func TestCreatePullRequest_Inputs_BranchSuffixOptions(t *testing.T) {
	testCases := []struct {
		name   string
		suffix string
	}{
		{"random", "random"},
		{"timestamp", "timestamp"},
		{"short-commit-hash", "short-commit-hash"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := CreatePullRequest{
				BranchSuffix:  tc.suffix,
				CommitMessage: "Update",
				Title:         "Update PR",
			}

			inputs := a.Inputs()

			if inputs["branch-suffix"] != tc.suffix {
				t.Errorf("branch-suffix = %v, want %q", inputs["branch-suffix"], tc.suffix)
			}
		})
	}
}

func TestCreatePullRequest_ImplementsStepAction(t *testing.T) {
	var _ workflow.StepAction = CreatePullRequest{}
}
