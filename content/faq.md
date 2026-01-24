---
title: "FAQ"
---

Frequently asked questions about wetwire-github-go.

---

## Getting Started

<details>
<summary>How do I install wetwire-github?</summary>

```bash
go install github.com/lex00/wetwire-github-go/cmd/wetwire-github@latest
```

See [Quick Start]({{< relref "/quick-start" >}}) for complete installation instructions.
</details>

<details>
<summary>How do I create a new project?</summary>

```bash
wetwire-github init my-workflows
cd my-workflows
```
</details>

<details>
<summary>How do I generate workflow YAML?</summary>

```bash
wetwire-github build ./my-workflows
# Outputs to .github/workflows/
```
</details>

<details>
<summary>What's the recommended project structure?</summary>

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

Keep related workflows together and extract shared patterns into dedicated files.
</details>

---

## Syntax

<details>
<summary>How do I declare a workflow?</summary>

```go
import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
    Name: "CI",
    On:   CITriggers,
}
```
</details>

<details>
<summary>How do I reference another job?</summary>

Use direct variable references in the `Needs` field:

```go
var Deploy = workflow.Job{
    Needs: []any{Build, Test},  // Direct variable references
}
```
</details>

<details>
<summary>How do I use action wrappers?</summary>

Import typed action wrappers and use them directly in `[]any{}` slices:

```go
import "github.com/lex00/wetwire-github-go/actions/checkout"

var CheckoutStep = checkout.Checkout{
    FetchDepth: 0,
    Submodules: "recursive",
}
```
</details>

<details>
<summary>How do I access secrets and matrix values?</summary>

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
</details>

<details>
<summary>How do I manage workflow secrets?</summary>

Always use the `workflow.Secrets` context instead of hardcoding values:

```go
var DeployStep = workflow.Step{
    Name: "Deploy",
    Run:  "deploy.sh",
    Env: workflow.Env{
        "DEPLOY_TOKEN": workflow.Secrets.Get("DEPLOY_TOKEN"),
        "GH_TOKEN":     workflow.Secrets.GITHUB_TOKEN(),
    },
}
```

For sensitive deployments, use GitHub Environments to scope secrets:

```go
var ProductionDeploy = workflow.Job{
    Name:        "Deploy to Production",
    RunsOn:      "ubuntu-latest",
    Environment: "production",
    Steps:       DeploySteps,
}
```

See [Security Patterns]({{< relref "/security-patterns" >}}) for more details.
</details>

---

## Lint Rules

<details>
<summary>What do WAG rules check?</summary>

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

See [Lint Rules]({{< relref "/lint-rules" >}}) for the complete reference.
</details>

<details>
<summary>How do I auto-fix lint issues?</summary>

```bash
wetwire-github lint --fix ./my-workflows
```
</details>

<details>
<summary>How does the linter help catch errors?</summary>

The linter runs static analysis on your Go code to detect:

- **Security issues**: Hardcoded secrets (WAG003, WAG020), dangerous trigger patterns (WAG018)
- **Type safety**: Raw action strings that should use typed wrappers (WAG001)
- **Best practices**: Missing timeouts (WAG014), missing permissions (WAG017)
- **Logic errors**: Circular dependencies (WAG019), unreachable jobs (WAG011)

Run the linter as part of your CI workflow:

```bash
wetwire-github lint .
wetwire-github build .
wetwire-github validate .github/workflows/*.yml
```
</details>

---

## Import

<details>
<summary>How do I convert an existing workflow?</summary>

```bash
wetwire-github import .github/workflows/ci.yml -o my-workflows/
```
</details>

<details>
<summary>Can I import existing workflow YAML files?</summary>

Yes! The `import` command converts YAML workflows to typed Go code:

```bash
# Import a single workflow
wetwire-github import .github/workflows/ci.yml -o my-workflows/

# Import with single-file output
wetwire-github import ci.yml --single-file -o my-workflows/

# Import to existing project (skip go.mod, README)
wetwire-github import ci.yml --no-scaffold -o existing-project/
```

The importer:
1. Parses the YAML into an intermediate representation
2. Flattens nested structures to named variables
3. Maps known actions to typed wrappers
4. Converts expressions to type-safe builders

See [Import Workflow]({{< relref "/import-workflow" >}}) for detailed documentation.
</details>

<details>
<summary>Import produced code with errors?</summary>

Import is best-effort. Complex workflows may need manual cleanup:

1. Run `wetwire-github lint --fix` to apply automatic fixes
2. Review and manually fix remaining issues
3. Check import logs for unsupported features
</details>

---

## Reusable Workflows

<details>
<summary>How do I handle reusable workflows?</summary>

Create reusable workflows with `workflow_call` triggers:

```go
var ReusableDeploy = workflow.Workflow{
    Name: "Reusable Deploy",
    On: workflow.Triggers{
        WorkflowCall: &workflow.WorkflowCallTrigger{
            Inputs: map[string]workflow.WorkflowInput{
                "environment": {
                    Type:        "string",
                    Required:    true,
                    Description: "Deployment environment",
                },
            },
            Secrets: map[string]workflow.WorkflowSecret{
                "deploy-token": {Required: true},
            },
        },
    },
    Jobs: map[string]workflow.Job{"deploy": DeployJob},
}
```

Call the reusable workflow from another workflow:

```go
var CallDeploy = workflow.Job{
    Uses: "./.github/workflows/deploy.yml",
    With: workflow.With{
        "environment": "production",
    },
    Secrets: "inherit",
}
```
</details>

---

## Config Types

<details>
<summary>What config types are supported?</summary>

| Config Type | Output Location | Status |
|-------------|-----------------|--------|
| GitHub Actions | `.github/workflows/*.yml` | Implemented |
| Dependabot | `.github/dependabot.yml` | Implemented |
| Issue Templates | `.github/ISSUE_TEMPLATE/*.yml` | Implemented |
| Discussion Templates | `.github/DISCUSSION_TEMPLATE/*.yml` | Implemented |
</details>

<details>
<summary>How do I generate Dependabot config?</summary>

```bash
wetwire-github build --type dependabot ./my-config
```
</details>

---

## Matrix Configuration

<details>
<summary>How do I define a build matrix?</summary>

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
</details>

---

## Troubleshooting

<details>
<summary>ModuleNotFoundError</summary>

Ensure the CLI is installed:

```bash
go install github.com/lex00/wetwire-github-go/cmd/wetwire-github@latest
```
</details>

<details>
<summary>Build produces empty output</summary>

Check that:
1. Workflows are declared as package-level `var` with struct literals
2. The package path is correct in the build command
3. Variables use `workflow.Workflow` or `workflow.Job` types
</details>

<details>
<summary>actionlint validation errors</summary>

The `validate` command uses actionlint to check generated YAML:

```bash
wetwire-github validate .github/workflows/ci.yml
```

Fix issues based on actionlint messages, then rebuild.
</details>

---

## Resources

- [CLI Reference]({{< relref "/cli" >}})
- [Quick Start]({{< relref "/quick-start" >}})
- [Examples]({{< relref "/examples" >}})
- [Lint Rules]({{< relref "/lint-rules" >}})
