---
title: "Import Workflow"
---
# Import Workflow Guide

Convert existing GitHub Actions YAML workflows to typed Go declarations.

## Basic Usage

```bash
wetwire-github import .github/workflows/ci.yml -o my-workflows/
```

This generates a complete Go project:

```
my-workflows/
├── go.mod                    # Module declaration
├── README.md                 # Generated documentation
├── cmd/main.go               # Usage instructions
└── workflows/
    ├── workflows.go          # Workflow declarations
    ├── jobs.go               # Job declarations
    ├── triggers.go           # Trigger configurations
    └── steps.go              # Step declarations
```

## How Import Works

### 1. Parse YAML

The importer parses the YAML workflow file into an intermediate representation (IR).

### 2. Flatten Nested Structures

All nested structures become flat variables at package scope:

**Input YAML:**
```yaml
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
```

**Generated Go:**
```go
var CIPush = workflow.PushTrigger{Branches: List("main")}
var CIPullRequest = workflow.PullRequestTrigger{Branches: List("main")}

var CITriggers = workflow.Triggers{
    Push:        CIPush,
    PullRequest: CIPullRequest,
}
```

### 3. Generate Variable Names

Variable names are derived from:
- Workflow name → `{WorkflowName}` (e.g., `CI`)
- Jobs → `{JobKey}` (e.g., `Build`, `Test`)
- Steps → `{JobKey}Step{Index}` or step name if provided
- Triggers → `{WorkflowName}Triggers`, `{WorkflowName}Push`, etc.

### 4. Map Action References

Known actions are converted to typed wrappers:

**Input:**
```yaml
- uses: actions/checkout@v4
  with:
    fetch-depth: 0
```

**Output:**
```go
checkout.Checkout{FetchDepth: 0}
```

Unknown actions remain as raw `workflow.Step`:

```go
workflow.Step{
    Uses: "some/unknown-action@v1",
    With: workflow.With{"key": "value"},
}
```

### 5. Handle Expressions

GitHub Actions expressions are converted to typed expression builders:

**Input:**
```yaml
if: github.ref == 'refs/heads/main'
env:
  TOKEN: ${{ secrets.DEPLOY_TOKEN }}
```

**Output:**
```go
workflow.Step{
    If: workflow.Branch("main"),
    Env: workflow.Env{
        "TOKEN": workflow.Secrets.Get("DEPLOY_TOKEN"),
    },
}
```

## Import Options

### `--single-file`

Generate all declarations in a single file instead of splitting by type:

```bash
wetwire-github import ci.yml --single-file -o my-workflows/
```

Output:
```
my-workflows/
├── go.mod
└── ci.go           # All declarations in one file
```

### `--no-scaffold`

Skip generating go.mod, README, cmd/main.go, etc.:

```bash
wetwire-github import ci.yml --no-scaffold -o existing-project/
```

Useful when adding to an existing wetwire project.

### `--type`

Import different configuration types:

```bash
# Dependabot
wetwire-github import .github/dependabot.yml --type dependabot -o my-config/

# Issue templates
wetwire-github import .github/ISSUE_TEMPLATE/bug.yml --type issue-template -o my-templates/

# Discussion templates
wetwire-github import .github/DISCUSSION_TEMPLATE/idea.yml --type discussion-template -o my-templates/
```

## Matrix Import

Matrix configurations are flattened to named variables:

**Input:**
```yaml
strategy:
  matrix:
    go: ['1.22', '1.23']
    os: [ubuntu-latest, macos-latest]
```

**Output:**
```go
var BuildMatrix = workflow.Matrix{
    Values: map[string][]any{
        "go": {"1.22", "1.23"},
        "os": {"ubuntu-latest", "macos-latest"},
    },
}

var BuildStrategy = workflow.Strategy{
    Matrix: BuildMatrix,
}
```

## Reusable Workflow Import

Reusable workflows with `workflow_call` triggers are fully supported:

**Input:**
```yaml
on:
  workflow_call:
    inputs:
      environment:
        type: string
        required: true
    secrets:
      deploy-token:
        required: true
```

**Output:**
```go
var DeployTriggers = workflow.Triggers{
    WorkflowCall: &workflow.WorkflowCallTrigger{
        Inputs: map[string]workflow.WorkflowInput{
            "environment": {Type: "string", Required: true},
        },
        Secrets: map[string]workflow.WorkflowSecret{
            "deploy-token": {Required: true},
        },
    },
}
```

## Local Action Detection

The importer detects local actions (`uses: ./path/to/action`) and preserves them as raw step references:

```go
workflow.Step{
    Uses: "./path/to/action",
    With: workflow.With{"key": "value"},
}
```

## Round-Trip Testing

Verify import/build fidelity:

```bash
# Import workflow
wetwire-github import original.yml -o test/

# Rebuild to YAML
wetwire-github build test/ -o test/.github/workflows/

# Compare (should be semantically equivalent)
diff original.yml test/.github/workflows/original.yml
```

Use the included test script for batch testing:

```bash
./scripts/import_samples.sh
```

## Troubleshooting

### Unknown Action

If an action isn't recognized, it's kept as a raw step. To add support:

1. Check if the action is in the supported list
2. Use `workflow.Step{Uses: "..."}` as a workaround
3. Request support via GitHub issue

### Expression Not Converted

Complex expressions may not convert to builders:

```go
// Complex expression kept as string
If: "${{ github.event.pull_request.head.repo.full_name == github.repository }}"
```

Consider simplifying or using the raw expression.

### Naming Conflicts

If generated names conflict, the importer appends a number:

```go
var Build = workflow.Job{...}
var Build2 = workflow.Job{...}  // Second job named "build"
```

Rename manually for clarity.
