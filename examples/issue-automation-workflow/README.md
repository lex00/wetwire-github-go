# Issue Automation Workflow Example

A complete example demonstrating how to define GitHub Actions workflows with Issues, IssueComment, and PullRequestReview triggers using wetwire-github-go.

## Features Demonstrated

- **Issues trigger** with opened and labeled events for auto-labeling
- **IssueComment trigger** with created event for responding to commands
- **PullRequestReview trigger** with submitted event for enforcing review policies
- **Typed action wrappers** for github-script

## Project Structure

```
issue-automation-workflow/
├── go.mod                    # Module with replace directive
├── README.md                 # This file
├── CLAUDE.md                 # AI assistant context
└── workflows/
    ├── workflows.go          # Workflow declarations
    ├── jobs.go               # Job definitions
    ├── triggers.go           # Trigger configurations
    └── steps.go              # Step sequences with github-script
```

## Usage

### Generate YAML

```bash
cd examples/issue-automation-workflow
go mod tidy
wetwire-github build ./workflows
```

This generates `.github/workflows/issue-automation.yml`.

### View Generated YAML

```bash
cat .github/workflows/issue-automation.yml
```

### Validate with actionlint

```bash
wetwire-github validate .github/workflows/issue-automation.yml
```

### Local Development

When developing wetwire-github-go locally, add a replace directive to go.mod:

```go
replace github.com/lex00/wetwire-github-go => ../..
```

Then run `go mod tidy` before building.

## Key Patterns

### Issues Trigger

Trigger workflows when issues are opened or labeled:

```go
var IssuesOpened = workflow.IssuesTrigger{
    Types: []string{"opened", "labeled"},
}

var WorkflowTriggers = workflow.Triggers{
    Issues: &IssuesOpened,
}
```

### IssueComment Trigger

Trigger workflows when comments are created on issues:

```go
var CommentCreated = workflow.IssueCommentTrigger{
    Types: []string{"created"},
}

var WorkflowTriggers = workflow.Triggers{
    IssueComment: &CommentCreated,
}
```

### PullRequestReview Trigger

Trigger workflows when PR reviews are submitted:

```go
var ReviewSubmitted = workflow.PullRequestReviewTrigger{
    Types: []string{"submitted"},
}

var WorkflowTriggers = workflow.Triggers{
    PullRequestReview: &ReviewSubmitted,
}
```

### Conditional Jobs

Run jobs based on event type and action:

```go
var AutoLabel = workflow.Job{
    If: "${{ github.event_name == 'issues' && github.event.action == 'opened' }}",
    // ...
}

var RespondToComment = workflow.Job{
    If: "${{ github.event_name == 'issue_comment' && github.event.action == 'created' }}",
    // ...
}
```

### GitHub Script for Automation

Use github-script for complex automation logic:

```go
github_script.GithubScript{
    Script: `
const issue = context.payload.issue;
await github.rest.issues.addLabels({
    owner: context.repo.owner,
    repo: context.repo.repo,
    issue_number: issue.number,
    labels: ['bug']
});
`,
}
```

## Automation Patterns

### Auto-Labeling Issues

Labels are automatically applied based on keywords in the issue title and body:

| Keyword | Label |
|---------|-------|
| bug, error, fix, broken, crash | `bug` |
| feature, enhance, improve, add | `enhancement` |
| docs, documentation, readme, typo | `documentation` |
| question, how to, help, confused | `question` |

### Comment Commands

Users can interact with issues using slash commands:

| Command | Action |
|---------|--------|
| `/help` | Show available commands |
| `/assign` | Assign issue to commenter |
| `/close` | Close the issue |

### Review Policy Enforcement

The workflow enforces these review policies:

1. **Self-approval prevention**: Automatically dismisses reviews where the author approves their own PR
2. **Approval tracking**: Labels PRs as `ready-to-merge` when they have 2+ approvals
3. **Review status**: Labels PRs as `needs-review` until they have sufficient approvals
