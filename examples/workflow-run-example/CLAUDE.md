# workflow-run-example

Example workflow using workflow_run trigger to respond to CI completion.

## What This Is

A reference implementation showing how to declare GitHub Actions workflows that trigger when another workflow completes. This example creates a deployment workflow that runs after CI passes.

## Key Files

- `workflows/workflows.go` - Main workflow declaration
- `workflows/jobs.go` - Deploy and Notify job definitions with conditional execution
- `workflows/triggers.go` - WorkflowRun trigger configuration
- `workflows/steps.go` - Step sequences with artifact download and context access

## Patterns Used

1. **Flat variables** - All structs are package-level variables, not nested
2. **Typed wrappers** - Uses `checkout.Checkout{}`, `dawidd6_download_artifact.DownloadArtifact{}`
3. **WorkflowRun trigger** - Responds to "CI" workflow completion on main branch
4. **Conditional jobs** - Deploy only runs when triggering workflow succeeded
5. **Context access** - Uses `github.event.workflow_run.*` for triggering workflow info

## Build Command

```bash
wetwire-github build ./workflows
```

Output: `.github/workflows/deploy-after-ci.yml`
