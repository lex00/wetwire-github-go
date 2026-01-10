# CI Workflow Example

A complete example demonstrating how to define GitHub Actions CI workflows with wetwire-github-go.

## Features Demonstrated

- **Workflow declaration** with CI triggers (push and pull request)
- **Matrix strategy** for testing multiple Go versions and OS combinations
- **Typed action wrappers** for checkout, setup-go, and cache
- **Multiple jobs** (build and lint) with their own step sequences

## Project Structure

```
ci-workflow/
├── go.mod                    # Module with replace directive
├── README.md                 # This file
├── CLAUDE.md                 # AI assistant context
└── workflows/
    ├── workflows.go          # Workflow declarations
    ├── jobs.go               # Job definitions with matrix
    ├── triggers.go           # Trigger configurations
    └── steps.go              # Step sequences
```

## Usage

### Generate YAML

```bash
cd examples/ci-workflow
go mod tidy
wetwire-github build .
```

This generates `.github/workflows/ci.yml`.

### View Generated YAML

```bash
cat .github/workflows/ci.yml
```

### Validate with actionlint

```bash
wetwire-github validate .github/workflows/ci.yml
```

### Local Development

When developing wetwire-github-go locally, add a replace directive to go.mod:

```go
replace github.com/lex00/wetwire-github-go => ../..
```

Then run `go mod tidy` before building.

## Key Patterns

### Typed Action Wrappers

Instead of raw `uses:` strings, use typed wrappers:

```go
checkout.Checkout{}.ToStep()
setup_go.SetupGo{GoVersion: "1.24"}.ToStep()
cache.Cache{Path: "~/go/pkg/mod", Key: "..."}.ToStep()
```

### Matrix Strategy

Define matrix values as a map and reference with expressions:

```go
var BuildMatrix = workflow.Matrix{
    Values: map[string][]any{
        "go": {"1.23", "1.24"},
        "os": {"ubuntu-latest", "macos-latest"},
    },
}

var Build = workflow.Job{
    RunsOn: "${{ matrix.os }}",
    Strategy: &BuildStrategy,
}
```

### Flat Variable Structure

Extract all nested structs to package-level variables for clarity:

```go
// Separate variables for each component
var CIPush = workflow.PushTrigger{...}
var CIPullRequest = workflow.PullRequestTrigger{...}
var CITriggers = workflow.Triggers{Push: &CIPush, PullRequest: &CIPullRequest}
```
