---
title: "Developers"
---

Comprehensive guide for developers working on wetwire-github-go.

## Table of Contents

- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Running Tests](#running-tests)
- [Adding Action Wrappers](#adding-action-wrappers)
- [Adding Lint Rules](#adding-lint-rules)
- [Contributing](#contributing)
- [Dependencies](#dependencies)

---

## Development Setup

### Prerequisites

- **Go 1.23+** (required)
- **git** (version control)

### Clone and Setup

```bash
# Clone repository
git clone https://github.com/lex00/wetwire-github-go.git
cd wetwire-github-go

# Download dependencies
go mod download

# Build CLI
go build -o wetwire-github ./cmd/wetwire-github

# Verify installation
./wetwire-github version
```

### Running Tests

```bash
# Run all tests
go test -v ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test -v ./internal/lint/...
```

---

## Project Structure

```
wetwire-github-go/
├── cmd/wetwire-github/           # CLI application
│   ├── main.go                   # Entry point, command registration
│   ├── build.go                  # build command
│   ├── validate.go               # validate command (actionlint)
│   ├── list.go                   # list command
│   ├── graph.go                  # graph command
│   ├── lint.go                   # lint command (WAG001-WAG012)
│   ├── init.go                   # init command
│   ├── import.go                 # import command
│   ├── design.go                 # design command (AI-assisted)
│   ├── test.go                   # test command (persona testing)
│   ├── mcp.go                    # MCP server for IDE integration
│   └── version.go                # version handling
│
├── internal/
│   ├── discover/                 # AST-based resource discovery
│   │   ├── discover.go           # Parse Go source for var declarations
│   │   ├── workflow.go           # Workflow/Job discovery
│   │   ├── dependabot.go         # Dependabot config discovery
│   │   ├── templates.go          # Issue/Discussion template discovery
│   │   ├── pr_template.go        # PR template discovery
│   │   ├── codeowners.go         # CODEOWNERS discovery
│   │   └── graph.go              # Dependency graph building
│   ├── template/                 # GitHub config builder
│   │   ├── builder.go            # Build YAML from discovered resources
│   │   ├── sort.go               # Topological sort for dependencies
│   │   ├── dependabot.go         # Dependabot config builder
│   │   ├── issue_template.go     # Issue template builder
│   │   ├── pr_template.go        # PR template builder
│   │   └── codeowners.go         # CODEOWNERS builder
│   ├── linter/                   # Lint rules (WAG001-WAG012)
│   │   ├── linter.go             # Linter engine
│   │   └── rules.go              # Rule implementations
│   ├── importer/                 # YAML-to-Go code generator
│   │   ├── ir.go                 # Intermediate representation
│   │   ├── parser.go             # YAML parser
│   │   └── codegen.go            # Go code generator
│   ├── runner/                   # Go code execution for value extraction
│   ├── serialize/                # YAML serialization
│   ├── validation/               # actionlint integration
│   ├── agent/                    # AI agent integration
│   ├── personas/                 # Developer personas for testing
│   └── scoring/                  # 5-dimension scoring system
│
├── workflow/                     # Core workflow types
│   ├── workflow.go               # Workflow, Job, Step types
│   ├── triggers.go               # Push, PR, Schedule triggers
│   ├── expressions.go            # Expression contexts (secrets, env, etc.)
│   ├── matrix.go                 # Matrix strategy types
│   ├── conditions.go             # Conditional expressions
│   └── helpers.go                # Utility functions (List, etc.)
│
├── actions/                      # Type-safe action wrappers
│   ├── checkout/                 # actions/checkout
│   ├── setup_go/                 # actions/setup-go
│   ├── setup_node/               # actions/setup-node
│   ├── setup_python/             # actions/setup-python
│   ├── setup_java/               # actions/setup-java
│   ├── setup_dotnet/             # actions/setup-dotnet
│   ├── setup_ruby/               # actions/setup-ruby
│   ├── setup_rust/               # dtolnay/rust-toolchain
│   ├── cache/                    # actions/cache
│   ├── upload_artifact/          # actions/upload-artifact
│   ├── download_artifact/        # actions/download-artifact
│   ├── gh_release/               # softprops/action-gh-release
│   ├── github_script/            # actions/github-script
│   ├── docker_login/             # docker/login-action
│   ├── docker_build_push/        # docker/build-push-action
│   ├── docker_setup_buildx/      # docker/setup-buildx-action
│   └── codecov/                  # codecov/codecov-action
│
├── dependabot/                   # Dependabot configuration types
├── templates/                    # Issue/Discussion template types
├── codeowners/                   # CODEOWNERS types
├── codegen/                      # Action wrapper code generation
│
├── contracts.go                  # Core types (OutputRef, etc.)
├── go.mod                        # Module definition
├── go.sum                        # Dependency checksums
└── docs/                         # Documentation
```

---

## Running Tests

```bash
# Run all tests
go test -v ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -v ./internal/lint/... -run TestWAG001

# Run with race detection
go test -race ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## Adding Action Wrappers

Action wrappers provide type-safe interfaces for GitHub Actions.

### Step 1: Create Package

```bash
mkdir -p actions/my_action
```

### Step 2: Implement Wrapper

Create `actions/my_action/my_action.go`:

```go
// Package my_action provides a typed wrapper for owner/my-action.
package my_action

import (
    "github.com/lex00/wetwire-github-go/workflow"
)

// MyAction wraps the owner/my-action@v1 action.
type MyAction struct {
    // Input description
    InputName string `yaml:"input-name,omitempty"`

    // Boolean input
    EnableFeature bool `yaml:"enable-feature,omitempty"`
}

// Action returns the action reference.
func (a MyAction) Action() string {
    return "owner/my-action@v1"
}

// ToStep converts this action to a workflow step.
func (a MyAction) ToStep() workflow.Step {
    with := make(workflow.With)

    if a.InputName != "" {
        with["input-name"] = a.InputName
    }
    if a.EnableFeature {
        with["enable-feature"] = a.EnableFeature
    }

    return workflow.Step{
        Uses: a.Action(),
        With: with,
    }
}
```

### Step 3: Add Tests

Create `actions/my_action/my_action_test.go`:

```go
package my_action

import "testing"

func TestMyAction_Action(t *testing.T) {
    a := MyAction{}
    if got := a.Action(); got != "owner/my-action@v1" {
        t.Errorf("Action() = %q, want %q", got, "owner/my-action@v1")
    }
}

func TestMyAction_Inputs(t *testing.T) {
    a := MyAction{
        InputName: "value",
        EnableFeature: true,
    }

    if a.InputName != "value" {
        t.Errorf("InputName = %v, want %q", a.InputName, "value")
    }
    if !a.EnableFeature {
        t.Errorf("EnableFeature = %v, want true", a.EnableFeature)
    }
}
```

### Step 4: Update Documentation

1. Add to `docs/ROADMAP.md` Action Wrappers table
2. Add to `CHANGELOG.md` under Unreleased

---

## Adding Lint Rules

Lint rules enforce best practices for workflow declarations.

### Step 1: Choose Rule ID

Rules follow the pattern `WAG<NNN>`:
- WAG = Wetwire Actions GitHub
- NNN = 3-digit rule number

Check `internal/lint/rules.go` for the next available number.

### Step 2: Implement Rule

Add to `internal/lint/rules.go`:

```go
// checkWAG013 checks for [description].
func checkWAG013(issues *[]Issue, decl DiscoveredDecl) {
    // Implementation
    if /* condition detected */ {
        *issues = append(*issues, Issue{
            Rule:     "WAG013",
            Severity: "warning", // or "error", "info"
            Message:  "Description of the issue",
            File:     decl.File,
            Line:     decl.Line,
            Fixable:  false,
        })
    }
}
```

### Step 3: Register Rule

Add to the `Lint()` function in `internal/lint/linter.go`:

```go
checkWAG013(&issues, decl)
```

### Step 4: Add Tests

Add to `internal/lint/linter_test.go`:

```go
func TestWAG013(t *testing.T) {
    src := `package test
    var BadExample = workflow.Workflow{...}
    `
    issues, err := lintSource(src)
    if err != nil {
        t.Fatal(err)
    }

    found := false
    for _, issue := range issues {
        if issue.Rule == "WAG013" {
            found = true
            break
        }
    }
    if !found {
        t.Error("expected WAG013 issue not found")
    }
}
```

### Step 5: Update Documentation

1. Add to `docs/ROADMAP.md` Lint Rules table
2. Add to `CHANGELOG.md` under Unreleased
3. Update `docs/CLI.md` lint command section

---

## Contributing

### Code Style

- Run `go fmt` before committing
- Run `go vet` to catch issues
- Follow existing patterns in the codebase

### Commit Messages

Use conventional commit format:

```
feat: add action wrapper for owner/action
fix: correct YAML serialization for matrix
docs: update CLI reference
test: add tests for WAG012
```

### Pull Request Process

1. Create feature branch from `main`
2. Write tests first (TDD)
3. Implement feature
4. Ensure all tests pass: `go test ./...`
5. Ensure no vet issues: `go vet ./...`
6. Update CHANGELOG.md
7. Update ROADMAP.md if applicable
8. Create PR with description

---

## Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/spf13/cobra` | CLI framework |
| `gopkg.in/yaml.v3` | YAML parsing/generation |
| `github.com/rhysd/actionlint` | GitHub Actions validation |
| `github.com/modelcontextprotocol/go-sdk` | MCP server for IDE integration |
| `github.com/anthropics/anthropic-sdk-go` | Claude API for design command |

---

## See Also

- [Quick Start](QUICK_START.md) - Getting started
- [CLI Reference](CLI.md) - CLI commands
- [Internals](INTERNALS.md) - Architecture details
- [FAQ](FAQ.md) - Common questions
