# wetwire-github-go

Declarative GitHub YAML configurations using native Go constructs.

## Overview

wetwire-github-go generates GitHub Actions workflows, Dependabot configs, and Issue Templates from typed Go declarations. No YAML required.

```go
package ci

import (
    "github.com/lex00/wetwire-github-go/workflow"
    "github.com/lex00/wetwire-github-go/actions/checkout"
    "github.com/lex00/wetwire-github-go/actions/setup_go"
)

// Flat variables for nested structs (importer generates correct syntax)
var CIPush = &workflow.PushTrigger{Branches: []string{"main"}}
var CIPullRequest = &workflow.PullRequestTrigger{Branches: []string{"main"}}

// Workflow declaration - references flat variables
var CI = workflow.Workflow{
    Name: "CI",
    On: workflow.Triggers{
        Push:        CIPush,
        PullRequest: CIPullRequest,
    },
}

// Job declaration - references workflow automatically via AST discovery
var Build = workflow.Job{
    Name:   "build",
    RunsOn: "ubuntu-latest",
    Steps: []workflow.Step{
        checkout.Checkout{}.ToStep(),
        setup_go.SetupGo{GoVersion: "1.23"}.ToStep(),
        {Run: "go build ./..."},
        {Run: "go test ./..."},
    },
}
```

Build to YAML:

```bash
wetwire-github build ./ci
# Outputs .github/workflows/ci.yml
```

## The "No Parens" Pattern

Resources are declared as Go variables using struct literals — no function calls needed:

```go
// Simple: just declare variables
var MyWorkflow = workflow.Workflow{...}
var BuildJob = workflow.Job{...}
var TestJob = workflow.Job{...}

// Cross-references via direct field access
var DeployJob = workflow.Job{
    Needs: []any{BuildJob, TestJob},  // Automatic dependency resolution
}

// Type-safe action wrappers
var CheckoutStep = checkout.Checkout{
    FetchDepth: 0,
    Submodules: "recursive",
}.ToStep()

// Expression contexts with helper types
var ConditionalStep = workflow.Step{
    If:  workflow.Branch("main").And(workflow.Push()),
    Run: "deploy.sh",
    Env: workflow.Env{
        "TOKEN": workflow.Secrets.Get("DEPLOY_TOKEN"),
    },
}

// Importer generates flat variables with correct syntax
// (users don't decide & — tooling handles it based on field types)
var BuildMatrix = &workflow.Matrix{
    Values: map[string][]any{"go": {"1.22", "1.23"}},
}

var BuildStrategy = &workflow.Strategy{
    FailFast: Ptr(false),
    Matrix:   BuildMatrix,
}

var MatrixJob = workflow.Job{
    RunsOn:   "ubuntu-latest",
    Strategy: BuildStrategy,
}
```

The CLI discovers declarations via **AST parsing** — no registration required.

## Generated Package Structure

Import existing workflows or init a new project:

```bash
wetwire-github import .github/workflows/ci.yml -o my-ci/
# OR
wetwire-github init my-ci/
```

Generated structure:
```
my-ci/
├── go.mod              # Module with wetwire-github-go dependency
├── README.md           # Generated docs
├── CLAUDE.md           # AI assistant context
├── cmd/main.go         # Usage instructions
├── workflows.go        # Workflow declarations
├── jobs.go             # Job declarations
└── triggers.go         # Trigger configurations
```

All nested structs become flat variables. The importer generates correct `&` syntax based on field types — users never decide manually.

## Scope

| Config Type | Output | Schema |
|-------------|--------|--------|
| **GitHub Actions** | `.github/workflows/*.yml` | workflow schema |
| **Dependabot** | `.github/dependabot.yml` | dependabot-2.0 |
| **Issue Templates** | `.github/ISSUE_TEMPLATE/*.yml` | issue-forms |

## Status

Under development. See [docs/PLAN.md](docs/PLAN.md) for implementation plan and [docs/ROADMAP.md](docs/ROADMAP.md) for feature matrix.

## License

MIT
