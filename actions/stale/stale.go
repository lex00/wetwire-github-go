// Package stale provides a typed wrapper for actions/stale.
package stale

// Stale wraps the actions/stale@v9 action.
// Marks issues and pull requests as stale and closes them after a period of inactivity.
type Stale struct {
	// Token for API access (default: ${{ github.token }})
	RepoToken string `yaml:"repo-token,omitempty"`

	// Message to post on issues when they are marked as stale
	StaleIssueMessage string `yaml:"stale-issue-message,omitempty"`

	// Message to post on pull requests when they are marked as stale
	StalePRMessage string `yaml:"stale-pr-message,omitempty"`

	// Message to post on issues when they are closed
	CloseIssueMessage string `yaml:"close-issue-message,omitempty"`

	// Message to post on pull requests when they are closed
	ClosePRMessage string `yaml:"close-pr-message,omitempty"`

	// Number of days of inactivity before an issue is marked as stale (default: 60)
	DaysBeforeStale int `yaml:"days-before-stale,omitempty"`

	// Number of days of inactivity before a stale issue is closed (default: 7)
	DaysBeforeClose int `yaml:"days-before-close,omitempty"`

	// Number of days of inactivity before an issue is marked as stale (overrides days-before-stale for issues)
	DaysBeforeIssueStale int `yaml:"days-before-issue-stale,omitempty"`

	// Number of days of inactivity before a pull request is marked as stale (overrides days-before-stale for PRs)
	DaysBeforePRStale int `yaml:"days-before-pr-stale,omitempty"`

	// Number of days of inactivity before a stale issue is closed (overrides days-before-close for issues)
	DaysBeforeIssueClose int `yaml:"days-before-issue-close,omitempty"`

	// Number of days of inactivity before a stale pull request is closed (overrides days-before-close for PRs)
	DaysBeforePRClose int `yaml:"days-before-pr-close,omitempty"`

	// Label to apply to stale issues (default: Stale)
	StaleIssueLabel string `yaml:"stale-issue-label,omitempty"`

	// Label to apply to stale pull requests (default: Stale)
	StalePRLabel string `yaml:"stale-pr-label,omitempty"`

	// Labels that exempt issues from being marked as stale (comma-separated)
	ExemptIssueLabels string `yaml:"exempt-issue-labels,omitempty"`

	// Labels that exempt pull requests from being marked as stale (comma-separated)
	ExemptPRLabels string `yaml:"exempt-pr-labels,omitempty"`

	// Only process issues with these labels (comma-separated)
	OnlyLabels string `yaml:"only-labels,omitempty"`

	// Only process issues with these labels (comma-separated)
	OnlyIssueLabels string `yaml:"only-issue-labels,omitempty"`

	// Only process pull requests with these labels (comma-separated)
	OnlyPRLabels string `yaml:"only-pr-labels,omitempty"`

	// Maximum number of operations per run (to avoid API rate limits, default: 30)
	OperationsPerRun int `yaml:"operations-per-run,omitempty"`

	// Remove stale labels when issues/PRs are updated or commented on
	RemoveStaleWhenUpdated bool `yaml:"remove-stale-when-updated,omitempty"`

	// Remove issue stale labels when issues are updated or commented on
	RemoveIssueStaleWhenUpdated bool `yaml:"remove-issue-stale-when-updated,omitempty"`

	// Remove PR stale labels when PRs are updated or commented on
	RemovePRStaleWhenUpdated bool `yaml:"remove-pr-stale-when-updated,omitempty"`

	// Enable debug logging
	DebugOnly bool `yaml:"debug-only,omitempty"`

	// Order to process issues: asc (oldest first) or desc (newest first, default: false = asc)
	Ascending bool `yaml:"ascending,omitempty"`

	// Delete the branch when closing a stale pull request
	DeleteBranch bool `yaml:"delete-branch,omitempty"`

	// Date used to determine staleness in ISO 8601 format
	StartDate string `yaml:"start-date,omitempty"`

	// Exempt all issues/PRs with assignees
	ExemptAssignees bool `yaml:"exempt-assignees,omitempty"`

	// Exempt all issues with assignees
	ExemptIssueAssignees bool `yaml:"exempt-issue-assignees,omitempty"`

	// Exempt all PRs with assignees
	ExemptPRAssignees bool `yaml:"exempt-pr-assignees,omitempty"`

	// Exempt all issues/PRs with milestones
	ExemptMilestones bool `yaml:"exempt-milestones,omitempty"`

	// Exempt all issues with milestones
	ExemptIssueMilestones bool `yaml:"exempt-issue-milestones,omitempty"`

	// Exempt all PRs with milestones
	ExemptPRMilestones bool `yaml:"exempt-pr-milestones,omitempty"`

	// Exempt all issues/PRs with these milestones (comma-separated)
	ExemptAllMilestones bool `yaml:"exempt-all-milestones,omitempty"`

	// Exempt all issues with these milestones (comma-separated)
	ExemptAllIssueMilestones bool `yaml:"exempt-all-issue-milestones,omitempty"`

	// Exempt all PRs with these milestones (comma-separated)
	ExemptAllPRMilestones bool `yaml:"exempt-all-pr-milestones,omitempty"`

	// Enable statistics logging
	EnableStatistics bool `yaml:"enable-statistics,omitempty"`

	// Labels that exempt issues from being closed (comma-separated)
	LabelsToRemoveWhenStale string `yaml:"labels-to-remove-when-stale,omitempty"`

	// Labels to add when an issue is marked as stale (comma-separated)
	LabelsToAddWhenUnstale string `yaml:"labels-to-add-when-unstale,omitempty"`

	// Ignore all issues
	IgnoreIssues bool `yaml:"ignore-issues,omitempty"`

	// Ignore all pull requests
	IgnorePRs bool `yaml:"ignore-prs,omitempty"`

	// Ignore all updates to issues/PRs
	IgnoreUpdates bool `yaml:"ignore-updates,omitempty"`

	// Label to apply when closing an issue
	CloseIssueLabel string `yaml:"close-issue-label,omitempty"`

	// Label to apply when closing a pull request
	ClosePRLabel string `yaml:"close-pr-label,omitempty"`

	// Comma-separated list of labels that can be assigned by anyone
	AnyOfLabels string `yaml:"any-of-labels,omitempty"`

	// Comma-separated list of labels that can be assigned to issues by anyone
	AnyOfIssueLabels string `yaml:"any-of-issue-labels,omitempty"`

	// Comma-separated list of labels that can be assigned to PRs by anyone
	AnyOfPRLabels string `yaml:"any-of-pr-labels,omitempty"`

	// Include only issues/PRs with activity before this date in ISO 8601 format
	IncludeOnlyAssigned bool `yaml:"include-only-assigned,omitempty"`
}

// Action returns the action reference.
func (a Stale) Action() string {
	return "actions/stale@v9"
}

// Inputs returns the action inputs as a map.
func (a Stale) Inputs() map[string]any {
	with := make(map[string]any)

	if a.RepoToken != "" {
		with["repo-token"] = a.RepoToken
	}
	if a.StaleIssueMessage != "" {
		with["stale-issue-message"] = a.StaleIssueMessage
	}
	if a.StalePRMessage != "" {
		with["stale-pr-message"] = a.StalePRMessage
	}
	if a.CloseIssueMessage != "" {
		with["close-issue-message"] = a.CloseIssueMessage
	}
	if a.ClosePRMessage != "" {
		with["close-pr-message"] = a.ClosePRMessage
	}
	if a.DaysBeforeStale != 0 {
		with["days-before-stale"] = a.DaysBeforeStale
	}
	if a.DaysBeforeClose != 0 {
		with["days-before-close"] = a.DaysBeforeClose
	}
	if a.DaysBeforeIssueStale != 0 {
		with["days-before-issue-stale"] = a.DaysBeforeIssueStale
	}
	if a.DaysBeforePRStale != 0 {
		with["days-before-pr-stale"] = a.DaysBeforePRStale
	}
	if a.DaysBeforeIssueClose != 0 {
		with["days-before-issue-close"] = a.DaysBeforeIssueClose
	}
	if a.DaysBeforePRClose != 0 {
		with["days-before-pr-close"] = a.DaysBeforePRClose
	}
	if a.StaleIssueLabel != "" {
		with["stale-issue-label"] = a.StaleIssueLabel
	}
	if a.StalePRLabel != "" {
		with["stale-pr-label"] = a.StalePRLabel
	}
	if a.ExemptIssueLabels != "" {
		with["exempt-issue-labels"] = a.ExemptIssueLabels
	}
	if a.ExemptPRLabels != "" {
		with["exempt-pr-labels"] = a.ExemptPRLabels
	}
	if a.OnlyLabels != "" {
		with["only-labels"] = a.OnlyLabels
	}
	if a.OnlyIssueLabels != "" {
		with["only-issue-labels"] = a.OnlyIssueLabels
	}
	if a.OnlyPRLabels != "" {
		with["only-pr-labels"] = a.OnlyPRLabels
	}
	if a.OperationsPerRun != 0 {
		with["operations-per-run"] = a.OperationsPerRun
	}
	if a.RemoveStaleWhenUpdated {
		with["remove-stale-when-updated"] = a.RemoveStaleWhenUpdated
	}
	if a.RemoveIssueStaleWhenUpdated {
		with["remove-issue-stale-when-updated"] = a.RemoveIssueStaleWhenUpdated
	}
	if a.RemovePRStaleWhenUpdated {
		with["remove-pr-stale-when-updated"] = a.RemovePRStaleWhenUpdated
	}
	if a.DebugOnly {
		with["debug-only"] = a.DebugOnly
	}
	if a.Ascending {
		with["ascending"] = a.Ascending
	}
	if a.DeleteBranch {
		with["delete-branch"] = a.DeleteBranch
	}
	if a.StartDate != "" {
		with["start-date"] = a.StartDate
	}
	if a.ExemptAssignees {
		with["exempt-assignees"] = a.ExemptAssignees
	}
	if a.ExemptIssueAssignees {
		with["exempt-issue-assignees"] = a.ExemptIssueAssignees
	}
	if a.ExemptPRAssignees {
		with["exempt-pr-assignees"] = a.ExemptPRAssignees
	}
	if a.ExemptMilestones {
		with["exempt-milestones"] = a.ExemptMilestones
	}
	if a.ExemptIssueMilestones {
		with["exempt-issue-milestones"] = a.ExemptIssueMilestones
	}
	if a.ExemptPRMilestones {
		with["exempt-pr-milestones"] = a.ExemptPRMilestones
	}
	if a.ExemptAllMilestones {
		with["exempt-all-milestones"] = a.ExemptAllMilestones
	}
	if a.ExemptAllIssueMilestones {
		with["exempt-all-issue-milestones"] = a.ExemptAllIssueMilestones
	}
	if a.ExemptAllPRMilestones {
		with["exempt-all-pr-milestones"] = a.ExemptAllPRMilestones
	}
	if a.EnableStatistics {
		with["enable-statistics"] = a.EnableStatistics
	}
	if a.LabelsToRemoveWhenStale != "" {
		with["labels-to-remove-when-stale"] = a.LabelsToRemoveWhenStale
	}
	if a.LabelsToAddWhenUnstale != "" {
		with["labels-to-add-when-unstale"] = a.LabelsToAddWhenUnstale
	}
	if a.IgnoreIssues {
		with["ignore-issues"] = a.IgnoreIssues
	}
	if a.IgnorePRs {
		with["ignore-prs"] = a.IgnorePRs
	}
	if a.IgnoreUpdates {
		with["ignore-updates"] = a.IgnoreUpdates
	}
	if a.CloseIssueLabel != "" {
		with["close-issue-label"] = a.CloseIssueLabel
	}
	if a.ClosePRLabel != "" {
		with["close-pr-label"] = a.ClosePRLabel
	}
	if a.AnyOfLabels != "" {
		with["any-of-labels"] = a.AnyOfLabels
	}
	if a.AnyOfIssueLabels != "" {
		with["any-of-issue-labels"] = a.AnyOfIssueLabels
	}
	if a.AnyOfPRLabels != "" {
		with["any-of-pr-labels"] = a.AnyOfPRLabels
	}
	if a.IncludeOnlyAssigned {
		with["include-only-assigned"] = a.IncludeOnlyAssigned
	}

	return with
}
