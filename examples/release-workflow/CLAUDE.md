# release-workflow

Example automated release workflow using wetwire-github-go.

## What This Is

A reference implementation showing how to declare GitHub release workflows using typed Go structs. This example creates a workflow that automatically generates releases when version tags are pushed.

## Key Files

- `workflows/workflows.go` - Main release workflow declaration
- `workflows/jobs.go` - Release job definition
- `workflows/triggers.go` - Push trigger for version tags (v*)
- `workflows/steps.go` - Step sequences using typed action wrappers

## Patterns Used

1. **Flat variables** - All structs are package-level variables, not nested
2. **Typed wrappers** - Uses `checkout.Checkout{}`, `gh_release.GHRelease{}`
3. **Tag triggers** - Triggers on version tags matching `v*` pattern
4. **Auto-generated release notes** - Uses GitHub's release notes generation

## Build Command

```bash
wetwire-github build ./workflows
```

Output: `.github/workflows/release.yml`
