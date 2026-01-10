# monorepo-workflow

Example monorepo CI workflow using wetwire-github-go.

## What This Is

A reference implementation showing how to declare GitHub Actions workflows for monorepo CI using typed Go structs. This example detects changes per service and conditionally builds only affected components.

## Key Files

- `workflows/workflows.go` - Main monorepo CI workflow declaration
- `workflows/jobs.go` - Job definitions with conditional execution
- `workflows/triggers.go` - Push and pull request triggers with path filters
- `workflows/steps.go` - Step sequences for change detection and service builds

## Patterns Used

1. **Path filters** - Triggers only fire when relevant paths change
2. **Change detection** - Uses `dorny/paths-filter@v3` to identify changed services
3. **Job outputs** - Exports change detection for conditional job execution
4. **Conditional jobs** - Uses `If` with expression helpers to skip unchanged services
5. **Typed wrappers** - Uses `checkout.Checkout{}`, `setup_go.SetupGo{}`, `setup_node.SetupNode{}`

## Build Command

```bash
wetwire-github build ./workflows
```

Output: `.github/workflows/monorepo.yml`
