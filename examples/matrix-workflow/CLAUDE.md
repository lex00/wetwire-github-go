# matrix-workflow

Example multi-OS/version matrix testing workflow using wetwire-github-go.

## What This Is

A reference implementation showing how to declare matrix strategy workflows using typed Go structs. This example creates a workflow that tests across multiple Go versions and operating systems.

## Key Files

- `workflows/workflows.go` - Main matrix test workflow declaration
- `workflows/jobs.go` - Test job with matrix strategy
- `workflows/triggers.go` - Push and pull request trigger configurations
- `workflows/steps.go` - Step sequences using typed action wrappers

## Patterns Used

1. **Flat variables** - All structs are package-level variables, not nested
2. **Typed wrappers** - Uses `checkout.Checkout{}`, `setup_go.SetupGo{}`
3. **Matrix strategy** - Tests Go 1.22 and 1.23 on ubuntu-latest and macos-latest
4. **Matrix expressions** - Uses `${{ matrix.os }}` and `${{ matrix.go }}`

## Build Command

```bash
wetwire-github build ./workflows
```

Output: `.github/workflows/matrix.yml`
