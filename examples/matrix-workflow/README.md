<picture>
  <source media="(prefers-color-scheme: dark)" srcset="../../docs/wetwire-dark.svg">
  <img src="../../docs/wetwire-light.svg" width="100" height="67">
</picture>

A complete example demonstrating how to define multi-language/OS matrix testing workflows with wetwire-github-go.

## Features Demonstrated

- **Matrix strategy** for testing multiple configurations
- **OS matrix** - Tests on ubuntu-latest and macos-latest
- **Version matrix** - Tests Go 1.22 and 1.23
- **Typed action wrappers** for checkout, setup-go, and cache
- **Matrix expressions** for dynamic runner and version selection

## Project Structure

```
matrix-workflow/
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
cd examples/matrix-workflow
go mod tidy
wetwire-github build .
```

This generates `.github/workflows/matrix.yml`.

### View Generated YAML

```bash
cat .github/workflows/matrix.yml
```

### Validate with actionlint

```bash
wetwire-github validate .github/workflows/matrix.yml
```

### Local Development

When developing wetwire-github-go locally, add a replace directive to go.mod:

```go
replace github.com/lex00/wetwire-github-go => ../..
```

Then run `go mod tidy` before building.

## Key Patterns

### Matrix Definition

Define matrix values as a flat variable:

```go
var TestMatrix = workflow.Matrix{
    Values: map[string][]any{
        "go": {"1.22", "1.23"},
        "os": {"ubuntu-latest", "macos-latest"},
    },
}
```

### Strategy Configuration

Reference matrix in strategy:

```go
var TestStrategy = workflow.Strategy{
    Matrix: &TestMatrix,
}
```

### Dynamic Runner Selection

Use matrix expression for runs-on:

```go
var Test = workflow.Job{
    RunsOn:   "${{ matrix.os }}",
    Strategy: &TestStrategy,
    ...
}
```

### Version Selection in Steps

Use matrix expression for tool versions:

```go
setup_go.SetupGo{
    GoVersion: "${{ matrix.go }}",
}
```

### Flat Variable Structure

Extract all nested structs to package-level variables for clarity:

```go
var TestMatrix = workflow.Matrix{...}
var TestStrategy = workflow.Strategy{Matrix: &TestMatrix}
var Test = workflow.Job{Strategy: &TestStrategy, ...}
```
