package stale

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestStale_Action(t *testing.T) {
	s := Stale{}
	if got := s.Action(); got != "actions/stale@v9" {
		t.Errorf("Action() = %q, want %q", got, "actions/stale@v9")
	}
}

func TestStale_Inputs(t *testing.T) {
	s := Stale{
		RepoToken:         "${{ secrets.GITHUB_TOKEN }}",
		StaleIssueMessage: "This issue is stale",
		StalePRMessage:    "This PR is stale",
		DaysBeforeStale:   30,
		DaysBeforeClose:   7,
		StaleIssueLabel:   "stale",
		StalePRLabel:      "stale",
		OperationsPerRun:  50,
	}

	inputs := s.Inputs()

	if inputs["repo-token"] != "${{ secrets.GITHUB_TOKEN }}" {
		t.Errorf("inputs[repo-token] = %v, want %q", inputs["repo-token"], "${{ secrets.GITHUB_TOKEN }}")
	}

	if inputs["stale-issue-message"] != "This issue is stale" {
		t.Errorf("inputs[stale-issue-message] = %v, want %q", inputs["stale-issue-message"], "This issue is stale")
	}

	if inputs["stale-pr-message"] != "This PR is stale" {
		t.Errorf("inputs[stale-pr-message] = %v, want %q", inputs["stale-pr-message"], "This PR is stale")
	}

	if inputs["days-before-stale"] != 30 {
		t.Errorf("inputs[days-before-stale] = %v, want 30", inputs["days-before-stale"])
	}

	if inputs["days-before-close"] != 7 {
		t.Errorf("inputs[days-before-close] = %v, want 7", inputs["days-before-close"])
	}

	if inputs["stale-issue-label"] != "stale" {
		t.Errorf("inputs[stale-issue-label] = %v, want %q", inputs["stale-issue-label"], "stale")
	}

	if inputs["stale-pr-label"] != "stale" {
		t.Errorf("inputs[stale-pr-label] = %v, want %q", inputs["stale-pr-label"], "stale")
	}

	if inputs["operations-per-run"] != 50 {
		t.Errorf("inputs[operations-per-run] = %v, want 50", inputs["operations-per-run"])
	}
}

func TestStale_Inputs_Empty(t *testing.T) {
	s := Stale{}
	inputs := s.Inputs()

	// Empty stale should have no inputs
	if len(inputs) != 0 {
		t.Errorf("empty Stale.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestStale_Inputs_BoolFields(t *testing.T) {
	s := Stale{
		RemoveStaleWhenUpdated:      true,
		DebugOnly:                   true,
		Ascending:                   true,
		DeleteBranch:                true,
		ExemptAssignees:             true,
		ExemptMilestones:            true,
		EnableStatistics:            true,
		IgnoreIssues:                true,
		IgnorePRs:                   true,
		RemoveIssueStaleWhenUpdated: true,
		RemovePRStaleWhenUpdated:    true,
	}

	inputs := s.Inputs()

	if inputs["remove-stale-when-updated"] != true {
		t.Errorf("inputs[remove-stale-when-updated] = %v, want true", inputs["remove-stale-when-updated"])
	}

	if inputs["debug-only"] != true {
		t.Errorf("inputs[debug-only] = %v, want true", inputs["debug-only"])
	}

	if inputs["ascending"] != true {
		t.Errorf("inputs[ascending] = %v, want true", inputs["ascending"])
	}

	if inputs["delete-branch"] != true {
		t.Errorf("inputs[delete-branch] = %v, want true", inputs["delete-branch"])
	}

	if inputs["exempt-assignees"] != true {
		t.Errorf("inputs[exempt-assignees] = %v, want true", inputs["exempt-assignees"])
	}

	if inputs["exempt-milestones"] != true {
		t.Errorf("inputs[exempt-milestones] = %v, want true", inputs["exempt-milestones"])
	}

	if inputs["enable-statistics"] != true {
		t.Errorf("inputs[enable-statistics] = %v, want true", inputs["enable-statistics"])
	}

	if inputs["ignore-issues"] != true {
		t.Errorf("inputs[ignore-issues] = %v, want true", inputs["ignore-issues"])
	}

	if inputs["ignore-prs"] != true {
		t.Errorf("inputs[ignore-prs] = %v, want true", inputs["ignore-prs"])
	}

	if inputs["remove-issue-stale-when-updated"] != true {
		t.Errorf("inputs[remove-issue-stale-when-updated] = %v, want true", inputs["remove-issue-stale-when-updated"])
	}

	if inputs["remove-pr-stale-when-updated"] != true {
		t.Errorf("inputs[remove-pr-stale-when-updated] = %v, want true", inputs["remove-pr-stale-when-updated"])
	}
}

func TestStale_Inputs_Labels(t *testing.T) {
	s := Stale{
		ExemptIssueLabels:       "wontfix,help wanted",
		ExemptPRLabels:          "wip,needs-review",
		OnlyLabels:              "bug",
		OnlyIssueLabels:         "enhancement",
		OnlyPRLabels:            "feature",
		LabelsToRemoveWhenStale: "needs-triage",
		LabelsToAddWhenUnstale:  "active",
		CloseIssueLabel:         "closed-stale",
		ClosePRLabel:            "closed-stale",
		AnyOfLabels:             "label1,label2",
		AnyOfIssueLabels:        "issue-label1,issue-label2",
		AnyOfPRLabels:           "pr-label1,pr-label2",
	}

	inputs := s.Inputs()

	if inputs["exempt-issue-labels"] != "wontfix,help wanted" {
		t.Errorf("inputs[exempt-issue-labels] = %v, want %q", inputs["exempt-issue-labels"], "wontfix,help wanted")
	}

	if inputs["exempt-pr-labels"] != "wip,needs-review" {
		t.Errorf("inputs[exempt-pr-labels] = %v, want %q", inputs["exempt-pr-labels"], "wip,needs-review")
	}

	if inputs["only-labels"] != "bug" {
		t.Errorf("inputs[only-labels] = %v, want %q", inputs["only-labels"], "bug")
	}

	if inputs["only-issue-labels"] != "enhancement" {
		t.Errorf("inputs[only-issue-labels] = %v, want %q", inputs["only-issue-labels"], "enhancement")
	}

	if inputs["only-pr-labels"] != "feature" {
		t.Errorf("inputs[only-pr-labels] = %v, want %q", inputs["only-pr-labels"], "feature")
	}

	if inputs["labels-to-remove-when-stale"] != "needs-triage" {
		t.Errorf("inputs[labels-to-remove-when-stale] = %v, want %q", inputs["labels-to-remove-when-stale"], "needs-triage")
	}

	if inputs["labels-to-add-when-unstale"] != "active" {
		t.Errorf("inputs[labels-to-add-when-unstale] = %v, want %q", inputs["labels-to-add-when-unstale"], "active")
	}

	if inputs["close-issue-label"] != "closed-stale" {
		t.Errorf("inputs[close-issue-label] = %v, want %q", inputs["close-issue-label"], "closed-stale")
	}

	if inputs["close-pr-label"] != "closed-stale" {
		t.Errorf("inputs[close-pr-label] = %v, want %q", inputs["close-pr-label"], "closed-stale")
	}

	if inputs["any-of-labels"] != "label1,label2" {
		t.Errorf("inputs[any-of-labels] = %v, want %q", inputs["any-of-labels"], "label1,label2")
	}

	if inputs["any-of-issue-labels"] != "issue-label1,issue-label2" {
		t.Errorf("inputs[any-of-issue-labels] = %v, want %q", inputs["any-of-issue-labels"], "issue-label1,issue-label2")
	}

	if inputs["any-of-pr-labels"] != "pr-label1,pr-label2" {
		t.Errorf("inputs[any-of-pr-labels] = %v, want %q", inputs["any-of-pr-labels"], "pr-label1,pr-label2")
	}
}

func TestStale_Inputs_Messages(t *testing.T) {
	s := Stale{
		CloseIssueMessage: "Closing this issue due to inactivity",
		ClosePRMessage:    "Closing this PR due to inactivity",
	}

	inputs := s.Inputs()

	if inputs["close-issue-message"] != "Closing this issue due to inactivity" {
		t.Errorf("inputs[close-issue-message] = %v, want %q", inputs["close-issue-message"], "Closing this issue due to inactivity")
	}

	if inputs["close-pr-message"] != "Closing this PR due to inactivity" {
		t.Errorf("inputs[close-pr-message] = %v, want %q", inputs["close-pr-message"], "Closing this PR due to inactivity")
	}
}

func TestStale_Inputs_DaysOverrides(t *testing.T) {
	s := Stale{
		DaysBeforeIssueStale: 45,
		DaysBeforePRStale:    30,
		DaysBeforeIssueClose: 10,
		DaysBeforePRClose:    5,
	}

	inputs := s.Inputs()

	if inputs["days-before-issue-stale"] != 45 {
		t.Errorf("inputs[days-before-issue-stale] = %v, want 45", inputs["days-before-issue-stale"])
	}

	if inputs["days-before-pr-stale"] != 30 {
		t.Errorf("inputs[days-before-pr-stale] = %v, want 30", inputs["days-before-pr-stale"])
	}

	if inputs["days-before-issue-close"] != 10 {
		t.Errorf("inputs[days-before-issue-close] = %v, want 10", inputs["days-before-issue-close"])
	}

	if inputs["days-before-pr-close"] != 5 {
		t.Errorf("inputs[days-before-pr-close] = %v, want 5", inputs["days-before-pr-close"])
	}
}

func TestStale_Inputs_ExemptFlags(t *testing.T) {
	s := Stale{
		ExemptIssueAssignees:     true,
		ExemptPRAssignees:        true,
		ExemptIssueMilestones:    true,
		ExemptPRMilestones:       true,
		ExemptAllMilestones:      true,
		ExemptAllIssueMilestones: true,
		ExemptAllPRMilestones:    true,
		IgnoreUpdates:            true,
		IncludeOnlyAssigned:      true,
	}

	inputs := s.Inputs()

	if inputs["exempt-issue-assignees"] != true {
		t.Errorf("inputs[exempt-issue-assignees] = %v, want true", inputs["exempt-issue-assignees"])
	}

	if inputs["exempt-pr-assignees"] != true {
		t.Errorf("inputs[exempt-pr-assignees] = %v, want true", inputs["exempt-pr-assignees"])
	}

	if inputs["exempt-issue-milestones"] != true {
		t.Errorf("inputs[exempt-issue-milestones] = %v, want true", inputs["exempt-issue-milestones"])
	}

	if inputs["exempt-pr-milestones"] != true {
		t.Errorf("inputs[exempt-pr-milestones] = %v, want true", inputs["exempt-pr-milestones"])
	}

	if inputs["exempt-all-milestones"] != true {
		t.Errorf("inputs[exempt-all-milestones] = %v, want true", inputs["exempt-all-milestones"])
	}

	if inputs["exempt-all-issue-milestones"] != true {
		t.Errorf("inputs[exempt-all-issue-milestones] = %v, want true", inputs["exempt-all-issue-milestones"])
	}

	if inputs["exempt-all-pr-milestones"] != true {
		t.Errorf("inputs[exempt-all-pr-milestones] = %v, want true", inputs["exempt-all-pr-milestones"])
	}

	if inputs["ignore-updates"] != true {
		t.Errorf("inputs[ignore-updates] = %v, want true", inputs["ignore-updates"])
	}

	if inputs["include-only-assigned"] != true {
		t.Errorf("inputs[include-only-assigned] = %v, want true", inputs["include-only-assigned"])
	}
}

func TestStale_Inputs_StartDate(t *testing.T) {
	s := Stale{
		StartDate: "2024-01-01T00:00:00Z",
	}

	inputs := s.Inputs()

	if inputs["start-date"] != "2024-01-01T00:00:00Z" {
		t.Errorf("inputs[start-date] = %v, want %q", inputs["start-date"], "2024-01-01T00:00:00Z")
	}
}

func TestStale_ImplementsStepAction(t *testing.T) {
	s := Stale{}
	// Verify Stale implements StepAction interface
	var _ workflow.StepAction = s
}

func TestStale_Inputs_ComprehensiveCoverage(t *testing.T) {
	// Test all fields to ensure comprehensive coverage
	s := Stale{
		RepoToken:                   "token",
		StaleIssueMessage:           "stale issue",
		StalePRMessage:              "stale pr",
		CloseIssueMessage:           "close issue",
		ClosePRMessage:              "close pr",
		DaysBeforeStale:             60,
		DaysBeforeClose:             7,
		DaysBeforeIssueStale:        45,
		DaysBeforePRStale:           30,
		DaysBeforeIssueClose:        10,
		DaysBeforePRClose:           5,
		StaleIssueLabel:             "stale",
		StalePRLabel:                "stale-pr",
		ExemptIssueLabels:           "exempt",
		ExemptPRLabels:              "exempt-pr",
		OnlyLabels:                  "only",
		OnlyIssueLabels:             "only-issue",
		OnlyPRLabels:                "only-pr",
		OperationsPerRun:            30,
		RemoveStaleWhenUpdated:      true,
		RemoveIssueStaleWhenUpdated: true,
		RemovePRStaleWhenUpdated:    true,
		DebugOnly:                   true,
		Ascending:                   true,
		DeleteBranch:                true,
		StartDate:                   "2024-01-01",
		ExemptAssignees:             true,
		ExemptIssueAssignees:        true,
		ExemptPRAssignees:           true,
		ExemptMilestones:            true,
		ExemptIssueMilestones:       true,
		ExemptPRMilestones:          true,
		ExemptAllMilestones:         true,
		ExemptAllIssueMilestones:    true,
		ExemptAllPRMilestones:       true,
		EnableStatistics:            true,
		LabelsToRemoveWhenStale:     "remove",
		LabelsToAddWhenUnstale:      "add",
		IgnoreIssues:                true,
		IgnorePRs:                   true,
		IgnoreUpdates:               true,
		CloseIssueLabel:             "closed",
		ClosePRLabel:                "closed-pr",
		AnyOfLabels:                 "any",
		AnyOfIssueLabels:            "any-issue",
		AnyOfPRLabels:               "any-pr",
		IncludeOnlyAssigned:         true,
	}

	inputs := s.Inputs()

	// Verify we have all the inputs we expect
	expectedCount := 47 // Total number of non-zero fields
	if len(inputs) != expectedCount {
		t.Errorf("Inputs() returned %d entries, want %d", len(inputs), expectedCount)
	}

	// Spot check a few critical fields
	if inputs["repo-token"] != "token" {
		t.Errorf("inputs[repo-token] = %v, want %q", inputs["repo-token"], "token")
	}

	if inputs["operations-per-run"] != 30 {
		t.Errorf("inputs[operations-per-run] = %v, want 30", inputs["operations-per-run"])
	}

	if inputs["start-date"] != "2024-01-01" {
		t.Errorf("inputs[start-date] = %v, want %q", inputs["start-date"], "2024-01-01")
	}
}
