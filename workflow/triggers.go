package workflow

// Triggers defines all events that can trigger a workflow.
type Triggers struct {
	// Push runs workflow on push events.
	Push *PushTrigger `yaml:"push,omitempty"`

	// PullRequest runs workflow on pull_request events.
	PullRequest *PullRequestTrigger `yaml:"pull_request,omitempty"`

	// PullRequestTarget runs workflow on pull_request_target events.
	PullRequestTarget *PullRequestTargetTrigger `yaml:"pull_request_target,omitempty"`

	// Schedule runs workflow on a schedule.
	Schedule []ScheduleTrigger `yaml:"schedule,omitempty"`

	// WorkflowDispatch allows manual triggering.
	WorkflowDispatch *WorkflowDispatchTrigger `yaml:"workflow_dispatch,omitempty"`

	// WorkflowCall allows this workflow to be called by other workflows.
	WorkflowCall *WorkflowCallTrigger `yaml:"workflow_call,omitempty"`

	// WorkflowRun runs after another workflow completes.
	WorkflowRun *WorkflowRunTrigger `yaml:"workflow_run,omitempty"`

	// RepositoryDispatch runs on repository_dispatch events.
	RepositoryDispatch *RepositoryDispatchTrigger `yaml:"repository_dispatch,omitempty"`

	// Create runs when a branch or tag is created.
	Create *CreateTrigger `yaml:"create,omitempty"`

	// Delete runs when a branch or tag is deleted.
	Delete *DeleteTrigger `yaml:"delete,omitempty"`

	// Fork runs when the repository is forked.
	Fork *ForkTrigger `yaml:"fork,omitempty"`

	// Gollum runs when wiki pages are created or updated.
	Gollum *GollumTrigger `yaml:"gollum,omitempty"`

	// IssueComment runs on issue comment events.
	IssueComment *IssueCommentTrigger `yaml:"issue_comment,omitempty"`

	// Issues runs on issue events.
	Issues *IssuesTrigger `yaml:"issues,omitempty"`

	// Label runs on label events.
	Label *LabelTrigger `yaml:"label,omitempty"`

	// Milestone runs on milestone events.
	Milestone *MilestoneTrigger `yaml:"milestone,omitempty"`

	// PageBuild runs on page build events.
	PageBuild *PageBuildTrigger `yaml:"page_build,omitempty"`

	// Project runs on project events.
	Project *ProjectTrigger `yaml:"project,omitempty"`

	// ProjectCard runs on project card events.
	ProjectCard *ProjectCardTrigger `yaml:"project_card,omitempty"`

	// ProjectColumn runs on project column events.
	ProjectColumn *ProjectColumnTrigger `yaml:"project_column,omitempty"`

	// Public runs when repository is made public.
	Public *PublicTrigger `yaml:"public,omitempty"`

	// PullRequestReview runs on pull request review events.
	PullRequestReview *PullRequestReviewTrigger `yaml:"pull_request_review,omitempty"`

	// PullRequestReviewComment runs on pull request review comment events.
	PullRequestReviewComment *PullRequestReviewCommentTrigger `yaml:"pull_request_review_comment,omitempty"`

	// Release runs on release events.
	Release *ReleaseTrigger `yaml:"release,omitempty"`

	// Status runs on commit status events.
	Status *StatusTrigger `yaml:"status,omitempty"`

	// Watch runs when repository is starred.
	Watch *WatchTrigger `yaml:"watch,omitempty"`

	// CheckRun runs on check run events.
	CheckRun *CheckRunTrigger `yaml:"check_run,omitempty"`

	// CheckSuite runs on check suite events.
	CheckSuite *CheckSuiteTrigger `yaml:"check_suite,omitempty"`

	// Discussion runs on discussion events.
	Discussion *DiscussionTrigger `yaml:"discussion,omitempty"`

	// DiscussionComment runs on discussion comment events.
	DiscussionComment *DiscussionCommentTrigger `yaml:"discussion_comment,omitempty"`

	// MergeGroup runs on merge group events.
	MergeGroup *MergeGroupTrigger `yaml:"merge_group,omitempty"`
}

// PushTrigger configures push event triggers.
type PushTrigger struct {
	Branches       []string `yaml:"branches,omitempty"`
	BranchesIgnore []string `yaml:"branches-ignore,omitempty"`
	Tags           []string `yaml:"tags,omitempty"`
	TagsIgnore     []string `yaml:"tags-ignore,omitempty"`
	Paths          []string `yaml:"paths,omitempty"`
	PathsIgnore    []string `yaml:"paths-ignore,omitempty"`
}

// PullRequestTrigger configures pull_request event triggers.
type PullRequestTrigger struct {
	Types          []string `yaml:"types,omitempty"`
	Branches       []string `yaml:"branches,omitempty"`
	BranchesIgnore []string `yaml:"branches-ignore,omitempty"`
	Paths          []string `yaml:"paths,omitempty"`
	PathsIgnore    []string `yaml:"paths-ignore,omitempty"`
}

// PullRequestTargetTrigger configures pull_request_target event triggers.
type PullRequestTargetTrigger struct {
	Types          []string `yaml:"types,omitempty"`
	Branches       []string `yaml:"branches,omitempty"`
	BranchesIgnore []string `yaml:"branches-ignore,omitempty"`
	Paths          []string `yaml:"paths,omitempty"`
	PathsIgnore    []string `yaml:"paths-ignore,omitempty"`
}

// ScheduleTrigger configures scheduled triggers using cron syntax.
type ScheduleTrigger struct {
	Cron string `yaml:"cron"`
}

// WorkflowDispatchTrigger configures manual workflow triggers with inputs.
type WorkflowDispatchTrigger struct {
	Inputs map[string]WorkflowInput `yaml:"inputs,omitempty"`
}

// WorkflowInput defines an input for workflow_dispatch or workflow_call.
type WorkflowInput struct {
	Description string   `yaml:"description,omitempty"`
	Required    bool     `yaml:"required,omitempty"`
	Default     any      `yaml:"default,omitempty"`
	Type        string   `yaml:"type,omitempty"` // string, boolean, choice, environment
	Options     []string `yaml:"options,omitempty"`
}

// WorkflowCallTrigger configures reusable workflow triggers.
type WorkflowCallTrigger struct {
	Inputs  map[string]WorkflowInput  `yaml:"inputs,omitempty"`
	Outputs map[string]WorkflowOutput `yaml:"outputs,omitempty"`
	Secrets map[string]WorkflowSecret `yaml:"secrets,omitempty"`
}

// WorkflowOutput defines an output from a reusable workflow.
type WorkflowOutput struct {
	Description string     `yaml:"description,omitempty"`
	Value       Expression `yaml:"value"`
}

// WorkflowSecret defines a secret input for workflow_call.
type WorkflowSecret struct {
	Description string `yaml:"description,omitempty"`
	Required    bool   `yaml:"required,omitempty"`
}

// WorkflowRunTrigger configures workflow_run event triggers.
type WorkflowRunTrigger struct {
	Workflows []string `yaml:"workflows,omitempty"`
	Types     []string `yaml:"types,omitempty"` // completed, requested, in_progress
	Branches  []string `yaml:"branches,omitempty"`
}

// RepositoryDispatchTrigger configures repository_dispatch event triggers.
type RepositoryDispatchTrigger struct {
	Types []string `yaml:"types,omitempty"`
}

// CreateTrigger configures create event triggers (branch/tag creation).
type CreateTrigger struct{}

// DeleteTrigger configures delete event triggers (branch/tag deletion).
type DeleteTrigger struct{}

// ForkTrigger configures fork event triggers.
type ForkTrigger struct{}

// GollumTrigger configures gollum (wiki) event triggers.
type GollumTrigger struct{}

// IssueCommentTrigger configures issue_comment event triggers.
type IssueCommentTrigger struct {
	Types []string `yaml:"types,omitempty"` // created, edited, deleted
}

// IssuesTrigger configures issues event triggers.
type IssuesTrigger struct {
	Types []string `yaml:"types,omitempty"`
}

// LabelTrigger configures label event triggers.
type LabelTrigger struct {
	Types []string `yaml:"types,omitempty"` // created, edited, deleted
}

// MilestoneTrigger configures milestone event triggers.
type MilestoneTrigger struct {
	Types []string `yaml:"types,omitempty"`
}

// PageBuildTrigger configures page_build event triggers.
type PageBuildTrigger struct{}

// ProjectTrigger configures project event triggers.
type ProjectTrigger struct {
	Types []string `yaml:"types,omitempty"`
}

// ProjectCardTrigger configures project_card event triggers.
type ProjectCardTrigger struct {
	Types []string `yaml:"types,omitempty"`
}

// ProjectColumnTrigger configures project_column event triggers.
type ProjectColumnTrigger struct {
	Types []string `yaml:"types,omitempty"`
}

// PublicTrigger configures public event triggers.
type PublicTrigger struct{}

// PullRequestReviewTrigger configures pull_request_review event triggers.
type PullRequestReviewTrigger struct {
	Types []string `yaml:"types,omitempty"` // submitted, edited, dismissed
}

// PullRequestReviewCommentTrigger configures pull_request_review_comment event triggers.
type PullRequestReviewCommentTrigger struct {
	Types []string `yaml:"types,omitempty"` // created, edited, deleted
}

// ReleaseTrigger configures release event triggers.
type ReleaseTrigger struct {
	Types []string `yaml:"types,omitempty"`
}

// StatusTrigger configures status event triggers.
type StatusTrigger struct{}

// WatchTrigger configures watch event triggers.
type WatchTrigger struct {
	Types []string `yaml:"types,omitempty"` // started
}

// CheckRunTrigger configures check_run event triggers.
type CheckRunTrigger struct {
	Types []string `yaml:"types,omitempty"` // created, rerequested, completed, requested_action
}

// CheckSuiteTrigger configures check_suite event triggers.
type CheckSuiteTrigger struct {
	Types []string `yaml:"types,omitempty"` // completed, requested, rerequested
}

// DiscussionTrigger configures discussion event triggers.
type DiscussionTrigger struct {
	Types []string `yaml:"types,omitempty"`
}

// DiscussionCommentTrigger configures discussion_comment event triggers.
type DiscussionCommentTrigger struct {
	Types []string `yaml:"types,omitempty"` // created, edited, deleted
}

// MergeGroupTrigger configures merge_group event triggers.
type MergeGroupTrigger struct {
	Types []string `yaml:"types,omitempty"` // checks_requested
}
