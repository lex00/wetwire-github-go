# issue-automation-workflow

Example workflow demonstrating issue and PR automation triggers.

## What This Is

A reference implementation showing how to declare GitHub Actions workflows with Issues, IssueComment, and PullRequestReview triggers using typed Go structs. This example creates a workflow that automates issue labeling, responds to commands in comments, and enforces PR review policies.

## Key Files

- `workflows/workflows.go` - Main workflow declaration
- `workflows/jobs.go` - AutoLabel, RespondToComment, and EnforceReviewPolicy job definitions
- `workflows/triggers.go` - Issues, IssueComment, and PullRequestReview trigger configurations
- `workflows/steps.go` - Step sequences using github-script for automation logic

## Patterns Used

1. **Flat variables** - All structs are package-level variables, not nested
2. **Typed wrappers** - Uses `github_script.GithubScript{}` for automation
3. **Issues trigger** - Responds to issue opened and labeled events
4. **IssueComment trigger** - Responds to comment created events
5. **PullRequestReview trigger** - Responds to review submitted events
6. **Conditional jobs** - Jobs that run based on event type and action

## Build Command

```bash
wetwire-github build ./workflows
```

Output: `.github/workflows/issue-automation.yml`
