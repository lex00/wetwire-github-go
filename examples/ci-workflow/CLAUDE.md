# ci-workflow

Example CI workflow using wetwire-github-go.

## What This Is

A reference implementation showing how to declare GitHub Actions workflows using typed Go structs. This example creates a CI workflow with build/test and lint jobs.

## Key Files

- `workflows/workflows.go` - Main CI workflow declaration
- `workflows/jobs.go` - Build and Lint job definitions with matrix strategy
- `workflows/triggers.go` - Push and pull request trigger configurations
- `workflows/steps.go` - Step sequences using typed action wrappers

## Patterns Used

1. **Flat variables** - All structs are package-level variables, not nested
2. **Typed wrappers** - Uses `checkout.Checkout{}`, `setup_go.SetupGo{}`, `cache.Cache{}`
3. **Matrix strategy** - Tests multiple Go versions and OS combinations
4. **Job references** - Jobs reference step slices by variable name

## Build Command

```bash
wetwire-github build ./workflows
```

Output: `.github/workflows/ci.yml`
