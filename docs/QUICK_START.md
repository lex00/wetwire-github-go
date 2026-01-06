# Quick Start

Get started with wetwire-github-go in 5 minutes.

## Installation

```bash
go install github.com/lex00/wetwire-github-go/cmd/wetwire-github@latest
```

## Create a New Project

```bash
wetwire-github init my-workflows
cd my-workflows
```

This creates:

```
my-workflows/
├── go.mod
├── workflows.go      # Your workflow declarations
├── jobs.go           # Job declarations
├── triggers.go       # Trigger configurations
└── cmd/main.go       # Usage instructions
```

## Define a Workflow

Edit `workflows.go`:

```go
package myworkflows

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

var CITriggers = workflow.Triggers{
    Push:        workflow.PushTrigger{Branches: List("main")},
    PullRequest: workflow.PullRequestTrigger{Branches: List("main")},
}
```

Edit `jobs.go`:

```go
package myworkflows

import (
    "github.com/lex00/wetwire-github-go/workflow"
    "github.com/lex00/wetwire-github-go/actions/checkout"
    "github.com/lex00/wetwire-github-go/actions/setup_go"
)

var Build = workflow.Job{
    Name:   "build",
    RunsOn: "ubuntu-latest",
    Steps:  BuildSteps,
}

var BuildSteps = List(
    checkout.Checkout{}.ToStep(),
    setup_go.SetupGo{GoVersion: "1.23"}.ToStep(),
    workflow.Step{Run: "go build ./..."},
    workflow.Step{Run: "go test ./..."},
)
```

## Build YAML

```bash
wetwire-github build .
```

Output in `.github/workflows/ci.yml`:

```yaml
name: CI
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  build:
    name: build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - run: go build ./...
      - run: go test ./...
```

## Import Existing Workflow

Convert an existing YAML workflow to Go:

```bash
wetwire-github import .github/workflows/ci.yml -o my-workflows/
```

## Validate

Check generated YAML with actionlint:

```bash
wetwire-github validate .github/workflows/ci.yml
```

## Next Steps

- See [CLI Reference](CLI.md) for all commands
- See [Import Workflow](IMPORT_WORKFLOW.md) for import details
- See [FAQ](FAQ.md) for common questions
