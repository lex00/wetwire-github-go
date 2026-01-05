# wetwire-github-go

Generate GitHub YAML configurations from typed Go declarations.

## Syntax: Simple, Flat, Declarative

All resources are Go struct literals. No function calls, no pointers, no registration.

### Workflows

```go
var CI = workflow.Workflow{
    Name: "CI",
    On:   CITriggers,
}

var CITriggers = workflow.Triggers{
    Push:        CIPush,
    PullRequest: CIPullRequest,
}

var CIPush = workflow.PushTrigger{Branches: []string{"main"}}
var CIPullRequest = workflow.PullRequestTrigger{Branches: []string{"main"}}
```

### Jobs

```go
var Build = workflow.Job{
    Name:   "build",
    RunsOn: "ubuntu-latest",
    Steps:  BuildSteps,
}

var BuildSteps = []workflow.Step{
    checkout.Checkout{}.ToStep(),
    setup_go.SetupGo{GoVersion: "1.23"}.ToStep(),
    {Run: "go build ./..."},
    {Run: "go test ./..."},
}
```

### Cross-References

Variables reference each other directly:

```go
var Deploy = workflow.Job{
    Needs: []any{Build, Test},  // References other jobs
    Steps: DeploySteps,
}
```

### Expression Contexts

```go
var ConditionalStep = workflow.Step{
    If:  workflow.Branch("main"),
    Run: "deploy.sh",
    Env: workflow.Env{
        "TOKEN": workflow.Secrets.Get("DEPLOY_TOKEN"),
    },
}
```

### Matrix

```go
var BuildMatrix = workflow.Matrix{
    Values: map[string][]any{
        "go":   {"1.22", "1.23"},
        "os":   {"ubuntu-latest", "macos-latest"},
    },
}

var BuildStrategy = workflow.Strategy{
    Matrix: BuildMatrix,
}

var MatrixJob = workflow.Job{
    RunsOn:   workflow.Matrix.Get("os"),
    Strategy: BuildStrategy,
}
```

### Action Wrappers

Type-safe wrappers for popular actions:

```go
checkout.Checkout{FetchDepth: 0, Submodules: "recursive"}.ToStep()
setup_go.SetupGo{GoVersion: "1.23"}.ToStep()
cache.Cache{Path: "~/.cache/go-build", Key: "go-cache"}.ToStep()
```

## Key Principles

1. **Flat variables** — Extract all nested structs into named variables
2. **No pointers** — Never use `&` or `*` in declarations
3. **Direct references** — Variables reference each other by name
4. **Struct literals only** — No function calls except `.ToStep()` for actions

## Build

```bash
wetwire-github build .
# Outputs .github/workflows/*.yml
```

## Project Structure

```
my-ci/
├── go.mod
├── workflows.go    # Workflow declarations
├── jobs.go         # Job declarations
├── triggers.go     # Trigger configurations
└── cmd/main.go     # Usage instructions
```
