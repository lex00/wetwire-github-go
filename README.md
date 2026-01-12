# wetwire-github-go

[![CI](https://github.com/lex00/wetwire-github-go/actions/workflows/ci.yml/badge.svg)](https://github.com/lex00/wetwire-github-go/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/lex00/wetwire-github-go/branch/main/graph/badge.svg)](https://codecov.io/gh/lex00/wetwire-github-go)
[![Go Reference](https://pkg.go.dev/badge/github.com/lex00/wetwire-github-go.svg)](https://pkg.go.dev/github.com/lex00/wetwire-github-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/lex00/wetwire-github-go)](https://goreportcard.com/report/github.com/lex00/wetwire-github-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

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

// Workflow declaration
var CI = workflow.Workflow{
    Name: "CI",
    On:   CITriggers,
}

// Triggers as flat variable
var CITriggers = workflow.Triggers{
    Push:        &CIPush,
    PullRequest: &CIPullRequest,
}

var CIPush = workflow.PushTrigger{Branches: []string{"main"}}
var CIPullRequest = workflow.PullRequestTrigger{Branches: []string{"main"}}

// Job declaration
var Build = workflow.Job{
    Name:   "build",
    RunsOn: "ubuntu-latest",
    Steps:  BuildSteps,
}

var BuildSteps = []any{
    checkout.Checkout{},
    setup_go.SetupGo{GoVersion: "1.23"},
    workflow.Step{Run: "go build ./..."},
    workflow.Step{Run: "go test ./..."},
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
// Declare variables
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
}

// Expression contexts
var ConditionalStep = workflow.Step{
    If:  workflow.Branch("main").And(workflow.Push()),
    Run: "deploy.sh",
    Env: workflow.Env{
        "TOKEN": workflow.Secrets.Get("DEPLOY_TOKEN"),
    },
}

// Matrix configuration
var BuildMatrix = workflow.Matrix{
    Values: map[string][]any{"go": {"1.22", "1.23"}},
}

var BuildStrategy = workflow.Strategy{
    Matrix: &BuildMatrix,
}

var MatrixJob = workflow.Job{
    RunsOn:   "ubuntu-latest",
    Strategy: &BuildStrategy,
}
```

The CLI discovers declarations via **AST parsing** — no registration required.

## Helpers

```go
// []string{} for typed slices
Branches: []string{"main", "develop"}

// []any{} for mixed-type slices
Needs: []any{BuildJob, TestJob}

// Env type alias
Env: workflow.Env{"TOKEN": workflow.Secrets.Get("TOKEN")}
```

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
├── go.mod                    # Module with wetwire-github-go dependency
├── README.md                 # Generated docs
├── cmd/main.go               # Usage instructions
└── workflows/
    ├── workflows.go          # Workflow declarations
    ├── jobs.go               # Job declarations
    ├── triggers.go           # Trigger configurations
    └── steps.go              # Step declarations
```

All nested structs become flat variables.

## Scope

| Config Type | Output | Schema |
|-------------|--------|--------|
| **GitHub Actions** | `.github/workflows/*.yml` | workflow schema |
| **Dependabot** | `.github/dependabot.yml` | dependabot-2.0 |
| **Issue Templates** | `.github/ISSUE_TEMPLATE/*.yml` | issue-forms |
| **Discussion Templates** | `.github/DISCUSSION_TEMPLATE/*.yml` | discussion-forms |

## Status

Under development. See [Roadmap](https://github.com/lex00/wetwire-github-go/issues/32) for implementation plan and feature matrix.

For the wetwire pattern, see the [Wetwire Specification](https://github.com/lex00/wetwire/blob/main/docs/WETWIRE_SPEC.md).

## License

MIT
