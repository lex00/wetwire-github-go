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

// Workflow declaration - no function calls, just struct initialization
var CI = workflow.Workflow{
    Name: "CI",
    On: workflow.Triggers{
        Push:        &workflow.PushTrigger{Branches: []string{"main"}},
        PullRequest: &workflow.PullRequestTrigger{Branches: []string{"main"}},
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

// Expression contexts
var ConditionalStep = workflow.Step{
    If:  workflow.Branch("main").And(workflow.Push()),
    Run: "deploy.sh",
    Env: map[string]any{
        "TOKEN": workflow.Secrets.Get("DEPLOY_TOKEN"),
    },
}
```

The CLI discovers declarations via **AST parsing** — no registration required.

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
