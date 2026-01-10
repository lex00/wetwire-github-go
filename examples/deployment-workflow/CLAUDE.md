# deployment-workflow

Example multi-environment deployment workflow using wetwire-github-go.

## What This Is

A reference implementation showing how to declare deployment workflows with environment-specific configurations using typed Go structs. This example creates a sequential deployment pipeline: build, deploy to staging, then deploy to production.

## Key Files

- `workflows/workflows.go` - Main deployment workflow declaration
- `workflows/jobs.go` - Build and deployment job definitions with environments
- `workflows/triggers.go` - Push and workflow_dispatch triggers
- `workflows/steps.go` - Step sequences using typed action wrappers

## Patterns Used

1. **Flat variables** - All structs are package-level variables, not nested
2. **Typed wrappers** - Uses `checkout.Checkout{}`, `setup_go.SetupGo{}`
3. **Environment configuration** - Uses `workflow.Environment{}` on jobs
4. **Environment secrets** - Uses `workflow.Secrets.Get()` for environment-scoped secrets
5. **Job dependencies** - Uses `Needs` for sequential deployment pipeline
6. **Manual dispatch** - Uses `workflow_dispatch` for manual triggering

## Build Command

```bash
wetwire-github build ./workflows
```

Output: `.github/workflows/deploy.yml`
