# Adoption Guide

A practical guide for teams adopting wetwire-github-go for GitHub Actions workflow management.

## Why Adopt wetwire-github-go?

- **Type Safety**: Catch workflow errors at compile time, not runtime
- **IDE Support**: Autocomplete, go-to-definition, refactoring
- **Reusability**: Share workflow patterns across repositories
- **Testability**: Unit test workflow logic before deployment
- **AI-Friendly**: Declarative struct syntax works well with AI code generation

## Getting Started

### Prerequisites

- Go 1.22 or later
- Git
- Existing GitHub repository with workflows (optional)

### Installation

See [README.md](../README.md#installation) for installation instructions.

## Migration Strategies

### Strategy 1: Start Fresh (Recommended for New Projects)

1. Initialize a new workflow project:
   ```bash
   wetwire-github init my-workflows
   cd my-workflows
   ```

2. Define your workflows in Go
3. Generate YAML:
   ```bash
   wetwire-github build .
   ```

### Strategy 2: Gradual Migration (Existing Projects)

Start with one workflow, validate, then expand:

1. **Import one workflow**:
   ```bash
   wetwire-github import .github/workflows/ci.yml -o workflows/
   ```

2. **Review and refine** the generated Go code

3. **Validate** the round-trip:
   ```bash
   wetwire-github build workflows/
   diff .github/workflows/ci.yml output/ci.yml
   ```

4. **Replace** the original with the generated file

5. **Repeat** for other workflows

### Strategy 3: Full Import (Quick Migration)

Import all workflows at once:

```bash
# Create workflow project
mkdir my-workflows && cd my-workflows
go mod init my-workflows

# Import all workflows
wetwire-github import ../.github/workflows/ -o .

# Verify
wetwire-github build .
wetwire-github validate .
```

## Best Practices

### Project Structure

Organize workflow code by concern:

```
workflows/
├── ci.go           # CI workflow
├── release.go      # Release workflow
├── triggers.go     # Shared trigger configurations
├── jobs/
│   ├── build.go    # Build job definitions
│   ├── test.go     # Test job definitions
│   └── deploy.go   # Deploy job definitions
└── shared/
    ├── matrix.go   # Shared matrix strategies
    └── env.go      # Environment configurations
```

### Naming Conventions

```go
// Workflows: PascalCase, descriptive
var CI = workflow.Workflow{...}
var ReleaseOnTag = workflow.Workflow{...}

// Jobs: PascalCase, action-oriented
var Build = workflow.Job{...}
var DeployProduction = workflow.Job{...}

// Triggers: Suffixed with "Trigger" or context
var CIPush = workflow.PushTrigger{...}
var ReleaseTags = workflow.PushTrigger{...}
```

### Reusable Patterns

Extract common patterns into shared variables:

```go
// shared/matrix.go
var GoVersionMatrix = workflow.Matrix{
    Values: map[string][]any{
        "go": {"1.22", "1.23"},
    },
}

// jobs/build.go
var Build = workflow.Job{
    Strategy: workflow.Strategy{Matrix: GoVersionMatrix},
    // ...
}
```

### Using Action Wrappers

Prefer typed wrappers over raw `uses:` strings:

```go
// Good: Type-safe, autocomplete support
var BuildSteps = []any{
    checkout.Checkout{FetchDepth: 0},
    setup_go.SetupGo{GoVersion: "1.23"},
}

// Avoid: Raw strings, no type checking
var BuildSteps = []any{
    workflow.Step{
        Uses: "actions/checkout@v4",
        With: map[string]any{"fetch-depth": 0},
    },
}
```

## CI/CD Integration

### GitHub Actions for Your Workflow Project

Add a workflow to validate your workflow definitions:

```go
var ValidateWorkflows = workflow.Workflow{
    Name: "Validate Workflows",
    On: workflow.Triggers{
        PullRequest: workflow.PullRequestTrigger{},
    },
    Jobs: map[string]workflow.Job{
        "validate": ValidateJob,
    },
}

var ValidateJob = workflow.Job{
    RunsOn: "ubuntu-latest",
    Steps: []any{
        checkout.Checkout{},
        setup_go.SetupGo{GoVersion: "1.23"},
        workflow.Step{Run: "go install github.com/lex00/wetwire-github-go/cmd/wetwire-github@latest"},
        workflow.Step{Run: "wetwire-github lint ."},
        workflow.Step{Run: "wetwire-github build ."},
        workflow.Step{Run: "wetwire-github validate ."},
    },
}
```

### Pre-commit Hooks

Add a pre-commit hook to regenerate YAML:

```bash
#!/bin/bash
# .git/hooks/pre-commit
wetwire-github build workflows/
git add .github/workflows/
```

## Troubleshooting

### Common Issues

**Import produces verbose code**

The importer generates explicit code. Refactor to use:
- Named variables instead of inline structs
- Action wrappers instead of raw `uses:` strings
- Shared triggers and matrix configurations

**Linter reports many issues**

Start with `--fix` for automatic fixes:
```bash
wetwire-github lint . --fix
```

Then address remaining issues manually.

**Generated YAML differs from original**

Some differences are expected:
- Key ordering may differ
- Empty values may be omitted
- Comments are not preserved

Use `wetwire-github validate .` to ensure the YAML is valid.

### Getting Help

- [FAQ](FAQ.md) - Common questions
- [CLI Reference](CLI.md) - Command documentation
- [GitHub Issues](https://github.com/lex00/wetwire-github-go/issues) - Bug reports and feature requests

## Team Adoption Checklist

- [ ] Install wetwire-github-go on all developer machines
- [ ] Import or create initial workflow definitions
- [ ] Add validation workflow to CI
- [ ] Document project structure conventions
- [ ] Set up pre-commit hooks (optional)
- [ ] Train team on basic usage
- [ ] Migrate remaining workflows gradually
