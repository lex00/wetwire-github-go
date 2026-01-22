<picture>
  <source media="(prefers-color-scheme: dark)" srcset="../../docs/wetwire-dark.svg">
  <img src="../../docs/wetwire-light.svg" width="100" height="67">
</picture>

A complete example demonstrating how to define automated release workflows with wetwire-github-go.

## Features Demonstrated

- **Tag-triggered releases** with version tag pattern (v*)
- **Auto-generated release notes** from commits
- **Typed action wrappers** for checkout and softprops/action-gh-release
- **GitHub token permissions** for creating releases

## Project Structure

```
release-workflow/
├── go.mod                    # Module with replace directive
├── README.md                 # This file
├── CLAUDE.md                 # AI assistant context
└── workflows/
    ├── workflows.go          # Workflow declarations
    ├── jobs.go               # Job definitions
    ├── triggers.go           # Trigger configurations
    └── steps.go              # Step sequences
```

## Usage

### Generate YAML

```bash
cd examples/release-workflow
go mod tidy
wetwire-github build .
```

This generates `.github/workflows/release.yml`.

### View Generated YAML

```bash
cat .github/workflows/release.yml
```

### Validate with actionlint

```bash
wetwire-github validate .github/workflows/release.yml
```

### Local Development

When developing wetwire-github-go locally, add a replace directive to go.mod:

```go
replace github.com/lex00/wetwire-github-go => ../..
```

Then run `go mod tidy` before building.

## Key Patterns

### Tag Trigger

Trigger releases only on version tags:

```go
var ReleasePush = workflow.PushTrigger{
    Tags: []string{"v*"},
}
```

### Release Action Wrapper

Use the typed wrapper for softprops/action-gh-release:

```go
gh_release.GHRelease{
    GenerateReleaseNotes: true,
    Draft:                false,
    Prerelease:           false,
}
```

### Flat Variable Structure

Extract all nested structs to package-level variables for clarity:

```go
// Separate variables for each component
var ReleasePush = workflow.PushTrigger{Tags: []string{"v*"}}
var ReleaseTriggers = workflow.Triggers{Push: &ReleasePush}
```
