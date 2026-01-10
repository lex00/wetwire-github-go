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
├── README.md
├── cmd/main.go               # Usage instructions
└── workflows/
    ├── workflows.go          # Your workflow declarations
    ├── jobs.go               # Job declarations
    ├── triggers.go           # Trigger configurations
    └── steps.go              # Step declarations
```

## Define a Workflow

Edit `workflows/workflows.go`:

```go
package workflows

import (
    "github.com/lex00/wetwire-github-go/workflow"
)

// Workflow declaration
var CI = workflow.Workflow{
    Name: "CI",
    On:   CITriggers,
    Jobs: map[string]workflow.Job{
        "build": Build,
    },
}
```

Edit `workflows/triggers.go`:

```go
package workflows

import (
    "github.com/lex00/wetwire-github-go/workflow"
)

var CIPush = workflow.PushTrigger{Branches: []string{"main"}}
var CIPullRequest = workflow.PullRequestTrigger{Branches: []string{"main"}}

var CITriggers = workflow.Triggers{
    Push:        &CIPush,
    PullRequest: &CIPullRequest,
}
```

Edit `workflows/jobs.go`:

```go
package workflows

import (
    "github.com/lex00/wetwire-github-go/workflow"
)

var Build = workflow.Job{
    RunsOn: "ubuntu-latest",
    Steps:  BuildSteps,
}
```

Edit `workflows/steps.go`:

```go
package workflows

import (
    "github.com/lex00/wetwire-github-go/workflow"
)

var BuildSteps = []workflow.Step{
    {Uses: "actions/checkout@v4"},
    {
        Uses: "actions/setup-go@v5",
        With: map[string]any{"go-version": "1.23"},
    },
    {Run: "go build ./..."},
    {Run: "go test ./..."},
}
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

## Using Typed Action Wrappers

Instead of raw `uses:` strings, use typed action wrappers for better IDE support and type safety.

### Checkout

```go
import "github.com/lex00/wetwire-github-go/actions/checkout"

var CheckoutStep = checkout.Checkout{
    FetchDepth: 0,          // Full history for git operations
    Submodules: "recursive", // Checkout submodules
}.ToStep()
```

### Setup Go

```go
import "github.com/lex00/wetwire-github-go/actions/setup_go"

var SetupGoStep = setup_go.SetupGo{
    GoVersion: "1.23",
    Cache:     true,
}.ToStep()
```

### Setup Node

```go
import "github.com/lex00/wetwire-github-go/actions/setup_node"

var SetupNodeStep = setup_node.SetupNode{
    NodeVersion: "20",
    Cache:       "npm",
}.ToStep()
```

### Setup Python

```go
import "github.com/lex00/wetwire-github-go/actions/setup_python"

var SetupPythonStep = setup_python.SetupPython{
    PythonVersion: "3.12",
    Cache:         "pip",
}.ToStep()
```

### Cache

```go
import "github.com/lex00/wetwire-github-go/actions/cache"

var CacheStep = cache.Cache{
    Path:        "~/.cache/go-build\n~/go/pkg/mod",
    Key:         "go-${{ runner.os }}-${{ hashFiles('**/go.sum') }}",
    RestoreKeys: "go-${{ runner.os }}-",
}.ToStep()
```

### Upload Artifact

```go
import "github.com/lex00/wetwire-github-go/actions/upload_artifact"

var UploadStep = upload_artifact.UploadArtifact{
    Name:          "build-artifacts",
    Path:          "dist/",
    RetentionDays: 7,
}.ToStep()
```

### Download Artifact

```go
import "github.com/lex00/wetwire-github-go/actions/download_artifact"

var DownloadStep = download_artifact.DownloadArtifact{
    Name: "build-artifacts",
    Path: "dist/",
}.ToStep()
```

### Complete Example with Typed Actions

```go
package workflows

import (
    "github.com/lex00/wetwire-github-go/workflow"
    "github.com/lex00/wetwire-github-go/actions/checkout"
    "github.com/lex00/wetwire-github-go/actions/setup_go"
    "github.com/lex00/wetwire-github-go/actions/cache"
)

var BuildSteps = []workflow.Step{
    checkout.Checkout{}.ToStep(),
    setup_go.SetupGo{GoVersion: "1.23"}.ToStep(),
    cache.Cache{
        Path: "~/go/pkg/mod",
        Key:  "go-mod-${{ hashFiles('**/go.sum') }}",
    }.ToStep(),
    {Run: "go build ./..."},
    {Run: "go test ./..."},
}
```

## Next Steps

- See [CLI Reference](CLI.md) for all commands
- See [Import Workflow](IMPORT_WORKFLOW.md) for import details
- See [FAQ](FAQ.md) for common questions
