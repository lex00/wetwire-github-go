<picture>
  <source media="(prefers-color-scheme: dark)" srcset="./wetwire-dark.svg">
  <img src="./wetwire-light.svg" width="100" height="67">
</picture>

This FAQ covers questions specific to the Go implementation of wetwire for GitHub Actions. For general wetwire questions, see the [central FAQ](https://github.com/lex00/wetwire/blob/main/docs/FAQ.md).

---

## Getting Started

### How do I install wetwire-github?

See [README.md](../README.md) for installation instructions.

### How do I create a new project?

```bash
wetwire-github init my-workflows
cd my-workflows
```

### How do I generate workflow YAML?

```bash
wetwire-github build ./my-workflows
# Outputs to .github/workflows/
```

---

## Syntax

### How do I declare a workflow?

```go
import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
    Name: "CI",
    On:   CITriggers,
}
```

### How do I reference another job?

Use direct variable references in the `Needs` field:

```go
var Deploy = workflow.Job{
    Needs: []any{Build, Test},  // Direct variable references
}
```

### How do I use action wrappers?

Import typed action wrappers and use them directly in `[]any{}` slices:

```go
import "github.com/lex00/wetwire-github-go/actions/checkout"

var CheckoutStep = checkout.Checkout{
    FetchDepth: 0,
    Submodules: "recursive",
}
```

### How do I access secrets and matrix values?

Use expression contexts:

```go
import "github.com/lex00/wetwire-github-go/workflow"

var DeployStep = workflow.Step{
    Env: workflow.Env{
        "TOKEN": workflow.Secrets.Get("DEPLOY_TOKEN"),
    },
    If: workflow.Branch("main"),
}
```

---

## Lint Rules

### What do WAG rules check?

WAG (Wetwire Action GitHub) lint rules enforce best practices:

| Rule | Description |
|------|-------------|
| WAG001 | Use typed action wrappers instead of raw `uses:` strings |
| WAG002 | Use condition builders instead of raw expression strings |
| WAG003 | Use secrets context instead of hardcoded strings |
| WAG004 | Use matrix builder instead of inline maps |
| WAG005 | Extract inline structs to named variables |
| WAG006 | Detect duplicate workflow names |
| WAG007 | Flag oversized files (>N jobs) |
| WAG008 | Avoid hardcoded expression strings |

### How do I auto-fix lint issues?

```bash
wetwire-github lint --fix ./my-workflows
```

---

## Import

### How do I convert an existing workflow?

```bash
wetwire-github import .github/workflows/ci.yml -o my-workflows/
```

### Import produced code with errors?

Import is best-effort. Complex workflows may need manual cleanup:

1. Run `wetwire-github lint --fix` to apply automatic fixes
2. Review and manually fix remaining issues
3. Check import logs for unsupported features

---

## Config Types

### What config types are supported?

| Config Type | Output Location | Status |
|-------------|-----------------|--------|
| GitHub Actions | `.github/workflows/*.yml` | Implemented |
| Dependabot | `.github/dependabot.yml` | Implemented |
| Issue Templates | `.github/ISSUE_TEMPLATE/*.yml` | Implemented |
| Discussion Templates | `.github/DISCUSSION_TEMPLATE/*.yml` | Implemented |

### How do I generate Dependabot config?

```bash
wetwire-github build --type dependabot ./my-config
```

---

## Matrix Configuration

### How do I define a build matrix?

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

---

## Troubleshooting

### ModuleNotFoundError

Ensure the CLI is installed:

```bash
go install github.com/lex00/wetwire-github-go/cmd/wetwire-github@latest
```

### Build produces empty output

Check that:
1. Workflows are declared as package-level `var` with struct literals
2. The package path is correct in the build command
3. Variables use `workflow.Workflow` or `workflow.Job` types

### actionlint validation errors

The `validate` command uses actionlint to check generated YAML:

```bash
wetwire-github validate .github/workflows/ci.yml
```

Fix issues based on actionlint messages, then rebuild.

---

## Resources

- [Wetwire Specification](https://github.com/lex00/wetwire/blob/main/docs/WETWIRE_SPEC.md)
- [CLI Documentation](CLI.md)
- [Quick Start](QUICK_START.md)
