# Workflow Run Trigger Example

A complete example demonstrating how to define GitHub Actions workflows that respond to other workflow completions using wetwire-github-go.

## Features Demonstrated

- **WorkflowRunTrigger** to respond when "CI" workflow completes
- **Conditional deployment** only when CI succeeds
- **Accessing triggering workflow context**:
  - `github.event.workflow_run.conclusion` - success/failure/cancelled
  - `github.event.workflow_run.head_sha` - commit that triggered CI
  - `github.event.workflow_run.head_branch` - branch that triggered CI
  - `github.event.workflow_run.id` - run ID for artifact downloads
- **Downloading artifacts** from the triggering workflow
- **Notification job** that runs regardless of CI conclusion

## Project Structure

```
workflow-run-example/
├── go.mod                    # Module with replace directive
├── README.md                 # This file
├── CLAUDE.md                 # AI assistant context
└── workflows/
    ├── workflows.go          # Workflow declarations
    ├── jobs.go               # Deploy and Notify job definitions
    ├── triggers.go           # WorkflowRun trigger configuration
    └── steps.go              # Step sequences with artifact handling
```

## Usage

### Generate YAML

```bash
cd examples/workflow-run-example
go mod tidy
wetwire-github build ./workflows
```

This generates `.github/workflows/deploy-after-ci.yml`.

### View Generated YAML

```bash
cat .github/workflows/deploy-after-ci.yml
```

### Validate with actionlint

```bash
wetwire-github validate .github/workflows/deploy-after-ci.yml
```

## Key Patterns

### WorkflowRun Trigger

Configure a workflow to run after another workflow completes:

```go
var CIWorkflowRun = workflow.WorkflowRunTrigger{
    Workflows: []string{"CI"},           // Name of triggering workflow
    Types:     []string{"completed"},    // Run when workflow completes
    Branches:  []string{"main"},         // Only for main branch
}

var DeployTriggers = workflow.Triggers{
    WorkflowRun: &CIWorkflowRun,
}
```

### Checking Workflow Conclusion

Run jobs conditionally based on the triggering workflow's result:

```go
var Deploy = workflow.Job{
    If: "${{ github.event.workflow_run.conclusion == 'success' }}",
    // ...
}
```

### Accessing Workflow Run Context

Access information about the triggering workflow in steps:

```go
workflow.Step{
    Run: `echo "Workflow: ${{ github.event.workflow_run.name }}"
echo "Conclusion: ${{ github.event.workflow_run.conclusion }}"
echo "Head SHA: ${{ github.event.workflow_run.head_sha }}"`,
}
```

### Downloading Artifacts from Triggering Workflow

Use dawidd6/action-download-artifact for cross-workflow artifact downloads:

```go
dawidd6_download_artifact.DownloadArtifact{
    GitHubToken:       "${{ secrets.GITHUB_TOKEN }}",
    RunID:             "${{ github.event.workflow_run.id }}",
    Name:              "build-artifacts",
    Path:              "./artifacts",
    IfNoArtifactFound: "warn",
}
```

### Checkout at Triggering Commit

Checkout the exact commit that triggered the CI workflow:

```go
checkout.Checkout{
    Ref: "${{ github.event.workflow_run.head_sha }}",
}
```

## Workflow Run Event Types

| Type | Description |
|------|-------------|
| `completed` | Workflow run has completed (success, failure, or cancelled) |
| `requested` | Workflow run has been requested |
| `in_progress` | Workflow run is in progress |

## Workflow Run Conclusion Values

| Value | Description |
|-------|-------------|
| `success` | All jobs completed successfully |
| `failure` | At least one job failed |
| `cancelled` | Workflow was cancelled |
| `skipped` | Workflow was skipped |
| `neutral` | Neutral conclusion |
| `timed_out` | Workflow timed out |
| `action_required` | Action is required |

## Common Use Cases

1. **Deploy after CI passes** - Deploy only when all tests succeed
2. **Send notifications** - Alert on workflow success or failure
3. **Trigger dependent workflows** - Chain workflows together
4. **Process artifacts** - Download and process build artifacts
5. **Update dashboards** - Update status dashboards after builds
