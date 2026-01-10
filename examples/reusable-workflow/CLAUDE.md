# reusable-workflow

Example reusable workflows using wetwire-github-go.

## What This Is

A reference implementation showing how to declare reusable GitHub Actions workflows using typed Go structs. This example creates a reusable build workflow with typed inputs and outputs, and a caller workflow that invokes it.

## Key Files

- `workflows/workflows.go` - Reusable and caller workflow declarations
- `workflows/jobs.go` - Build job with outputs and caller job using `uses`
- `workflows/triggers.go` - workflow_call trigger with inputs, outputs, and secrets
- `workflows/steps.go` - Step sequences for the build job

## Patterns Used

1. **Flat variables** - All structs are package-level variables, not nested
2. **Typed wrappers** - Uses `checkout.Checkout{}`, `setup_go.SetupGo{}`
3. **workflow_call trigger** - Defines reusable workflow with typed inputs/outputs
4. **Job outputs** - Passes values from job steps to workflow outputs
5. **Caller workflow** - Uses `uses` field to call reusable workflow
6. **Secrets handling** - Both explicit secrets and `secrets: inherit`

## Build Command

```bash
wetwire-github build ./workflows
```

Output:
- `.github/workflows/build-reusable.yml`
- `.github/workflows/ci-caller.yml`
