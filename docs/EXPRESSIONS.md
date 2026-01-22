<picture>
  <source media="(prefers-color-scheme: dark)" srcset="./wetwire-dark.svg">
  <img src="./wetwire-light.svg" width="100" height="67">
</picture>

This guide covers GitHub Actions expression contexts in wetwire-github-go. Expression contexts provide type-safe access to runtime values like secrets, matrix dimensions, and job outputs.

**Contents:**
- [Overview](#overview)
- [Secrets Context](#secrets-context)
- [GitHub Context](#github-context)
- [Matrix Expressions](#matrix-expressions)
- [Environment Variables](#environment-variables)
- [Needs Context](#needs-context)
- [Steps Context](#steps-context)
- [Inputs and Vars Contexts](#inputs-and-vars-contexts)
- [Condition Builders](#condition-builders)
- [String Functions](#string-functions)
- [Complex Expressions](#complex-expressions)
- [Security Considerations](#security-considerations)

---

## Overview

GitHub Actions uses `${{ expression }}` syntax to access runtime contexts. wetwire-github-go provides type-safe Go accessors that generate these expressions:

```go
import "github.com/lex00/wetwire-github-go/workflow"

// Type-safe accessor
workflow.Secrets.Get("DEPLOY_TOKEN")
// Generates: ${{ secrets.DEPLOY_TOKEN }}

// Condition builder
workflow.Branch("main")
// Generates: github.ref == 'refs/heads/main'
```

### Expression Type

The `workflow.Expression` type wraps expression strings and provides methods for combining expressions:

```go
expr := workflow.Branch("main")
expr.String() // Returns: ${{ github.ref == 'refs/heads/main' }}
expr.Raw()    // Returns: github.ref == 'refs/heads/main'
```

---

## Secrets Context

Use `workflow.Secrets` to access repository and organization secrets.

### Basic Usage

```go
import "github.com/lex00/wetwire-github-go/workflow"

var DeployStep = workflow.Step{
    Name: "Deploy",
    Run:  "deploy.sh",
    Env: workflow.Env{
        "DEPLOY_TOKEN": workflow.Secrets.Get("DEPLOY_TOKEN"),
        "API_KEY":      workflow.Secrets.Get("API_KEY"),
    },
}
```

### GITHUB_TOKEN

Access the built-in `GITHUB_TOKEN` using the dedicated method:

```go
var TokenStep = workflow.Step{
    Run: "gh api /repos/${{ github.repository }}",
    Env: workflow.Env{
        "GH_TOKEN": workflow.Secrets.GITHUB_TOKEN(),
    },
}
```

### Security Best Practices

> **Best Practice**: Always use `workflow.Secrets.Get()` instead of raw strings. This ensures secrets are properly referenced and makes security audits easier.

```go
// GOOD: Type-safe secret reference
"TOKEN": workflow.Secrets.Get("DEPLOY_TOKEN")

// BAD: Hardcoded secret (triggers WAG003 lint error)
"TOKEN": "ghp_xxxxxxxxxxxxxxxxxxxx"
```

**Related lint rule**: WAG003 detects hardcoded secrets (tokens matching patterns like `ghp_`, `ghs_`, `ghu_`, `github_pat_`).

---

## GitHub Context

Use `workflow.GitHub` to access information about the workflow run and repository.

### Available Properties

| Method | Expression | Description |
|--------|------------|-------------|
| `GitHub.Ref()` | `github.ref` | Full ref (e.g., `refs/heads/main`) |
| `GitHub.RefName()` | `github.ref_name` | Short ref name (e.g., `main`) |
| `GitHub.RefType()` | `github.ref_type` | `branch` or `tag` |
| `GitHub.SHA()` | `github.sha` | Commit SHA |
| `GitHub.Actor()` | `github.actor` | User who triggered the workflow |
| `GitHub.Repository()` | `github.repository` | `owner/repo` format |
| `GitHub.RepositoryOwner()` | `github.repository_owner` | Repository owner |
| `GitHub.EventName()` | `github.event_name` | Triggering event name |
| `GitHub.Workspace()` | `github.workspace` | Workspace directory |
| `GitHub.RunID()` | `github.run_id` | Unique run identifier |
| `GitHub.RunNumber()` | `github.run_number` | Run number for this workflow |
| `GitHub.RunAttempt()` | `github.run_attempt` | Retry attempt number |
| `GitHub.Job()` | `github.job` | Current job ID |
| `GitHub.Token()` | `github.token` | Automatic token (same as GITHUB_TOKEN) |
| `GitHub.HeadRef()` | `github.head_ref` | PR head branch |
| `GitHub.BaseRef()` | `github.base_ref` | PR base branch |
| `GitHub.ServerURL()` | `github.server_url` | GitHub server URL |
| `GitHub.APIURL()` | `github.api_url` | GitHub API URL |
| `GitHub.GraphQLURL()` | `github.graphql_url` | GitHub GraphQL URL |

### Event Data

Access event-specific data using `GitHub.Event()`:

```go
// Access pull request number
workflow.GitHub.Event("pull_request.number")
// Generates: github.event.pull_request.number

// Access issue title
workflow.GitHub.Event("issue.title")
// Generates: github.event.issue.title
```

### Usage Examples

```go
import (
    "github.com/lex00/wetwire-github-go/workflow"
    "github.com/lex00/wetwire-github-go/actions/docker_login"
)

// Use actor for Docker login
var LoginStep = docker_login.DockerLogin{
    Registry: "ghcr.io",
    Username: workflow.GitHub.Actor(),
    Password: workflow.Secrets.Get("GITHUB_TOKEN"),
}

// Tag Docker image with SHA
var BuildStep = workflow.Step{
    Name: "Build",
    Run:  "docker build -t myapp:${{ github.sha }} .",
}

// Conditional based on repository
var RepoStep = workflow.Step{
    Name: "Check Repo",
    Run:  "echo Repository: ${{ github.repository }}",
    Env: workflow.Env{
        "REPO":  workflow.GitHub.Repository(),
        "OWNER": workflow.GitHub.RepositoryOwner(),
    },
}
```

---

## Matrix Expressions

Use `workflow.MatrixContext` to access matrix dimension values in jobs with matrix strategies.

### Defining a Matrix

```go
var BuildMatrix = workflow.Matrix{
    Values: map[string][]any{
        "go":   {"1.22", "1.23"},
        "os":   {"ubuntu-latest", "macos-latest"},
    },
}

var BuildStrategy = workflow.Strategy{
    Matrix: &BuildMatrix,
}
```

### Accessing Matrix Values

```go
import "github.com/lex00/wetwire-github-go/workflow"

var MatrixJob = workflow.Job{
    Name:     "Build",
    RunsOn:   workflow.MatrixContext.Get("os"),  // ${{ matrix.os }}
    Strategy: &BuildStrategy,
    Steps:    MatrixSteps,
}

var MatrixSteps = []any{
    setup_go.SetupGo{
        GoVersion: workflow.MatrixContext.Get("go"),  // ${{ matrix.go }}
    },
}
```

### Matrix Dimension Safety

> **Best Practice**: Always ensure matrix dimensions have at least one value. Empty dimensions cause workflow failures.

```go
// GOOD: Matrix with values
var ValidMatrix = workflow.Matrix{
    Values: map[string][]any{
        "go": {"1.22", "1.23"},  // Has values
    },
}

// BAD: Empty dimension (triggers WAG009 lint error)
var InvalidMatrix = workflow.Matrix{
    Values: map[string][]any{
        "go": {},  // Empty - will fail
    },
}
```

**Related lint rule**: WAG009 validates matrix dimensions have values.

---

## Environment Variables

Use `workflow.EnvContext` to read environment variables in expressions.

### Reading Environment Variables

```go
var CheckStep = workflow.Step{
    Name: "Check CI",
    If:   workflow.EnvContext.Get("CI"),  // ${{ env.CI }}
    Run:  "echo Running in CI",
}
```

### Setting Environment Variables

Use `workflow.Env` to set environment variables for steps:

```go
var DeployStep = workflow.Step{
    Name: "Deploy",
    Run:  "deploy.sh",
    Env: workflow.Env{
        "ENVIRONMENT":  "production",
        "DEPLOY_TOKEN": workflow.Secrets.Get("DEPLOY_TOKEN"),
        "COMMIT_SHA":   workflow.GitHub.SHA(),
    },
}
```

### Workflow-Level Environment

Set environment variables for all jobs in a workflow:

```go
var CI = workflow.Workflow{
    Name: "CI",
    On:   CITriggers,
    Env: workflow.Env{
        "GO_VERSION": "1.23",
        "CI":         "true",
    },
    Jobs: Jobs{"build": Build},
}
```

---

## Needs Context

Use `workflow.Needs` to access outputs from dependent jobs.

### Passing Outputs Between Jobs

```go
// Job that produces outputs
var BuildJob = workflow.Job{
    Name:   "build",
    RunsOn: "ubuntu-latest",
    Outputs: map[string]string{
        "version": "${{ steps.version.outputs.version }}",
    },
    Steps: []any{
        workflow.Step{
            ID:   "version",
            Name: "Get Version",
            Run:  `echo "version=1.2.3" >> $GITHUB_OUTPUT`,
        },
    },
}

// Job that consumes outputs
var DeployJob = workflow.Job{
    Name:   "deploy",
    RunsOn: "ubuntu-latest",
    Needs:  []any{BuildJob},
    Steps: []any{
        workflow.Step{
            Name: "Deploy Version",
            Run:  "deploy.sh",
            Env: workflow.Env{
                "VERSION": workflow.Needs.Get("build", "version"),
            },
        },
    },
}
```

### Checking Job Results

Use `workflow.Needs.Result()` to check if a dependent job succeeded:

```go
var NotifyJob = workflow.Job{
    Name:   "notify",
    RunsOn: "ubuntu-latest",
    Needs:  []any{BuildJob, TestJob},
    If:     workflow.Always(),  // Run even if dependencies failed
    Steps: []any{
        workflow.Step{
            Name: "Check Build Result",
            Run:  "echo Build result: ${{ needs.build.result }}",
            Env: workflow.Env{
                "BUILD_RESULT": workflow.Needs.Result("build"),
                "TEST_RESULT":  workflow.Needs.Result("test"),
            },
        },
    },
}
```

---

## Steps Context

Use `workflow.Steps` to access outputs and status from previous steps.

### Step Outputs

```go
// Step that produces output
var SetVersionStep = workflow.Step{
    ID:   "version",
    Name: "Set Version",
    Run:  `echo "version=1.2.3" >> $GITHUB_OUTPUT`,
}

// Access output from another step
var UseVersionStep = workflow.Step{
    Name: "Use Version",
    Run:  "echo Version: ${{ steps.version.outputs.version }}",
    Env: workflow.Env{
        "VERSION": workflow.Steps.Get("version", "version"),
    },
}
```

### Using Step.Output() Helper

Steps have a convenient `Output()` method:

```go
var VersionStep = workflow.Step{
    ID:   "version",
    Name: "Get Version",
    Run:  `echo "version=1.0.0" >> $GITHUB_OUTPUT`,
}

// Later, reference the output
var DeployStep = workflow.Step{
    Env: workflow.Env{
        "VERSION": VersionStep.Output("version"),  // steps.version.outputs.version
    },
}
```

### Step Status

Check step outcomes for conditional execution:

```go
var RetryStep = workflow.Step{
    Name: "Retry on Failure",
    If:   workflow.Steps.Outcome("build").Raw() + " == 'failure'",
    Run:  "retry-build.sh",
}

// Outcomes: success, failure, cancelled, skipped
// Conclusions: success, failure, cancelled, skipped, neutral
```

---

## Inputs and Vars Contexts

### Workflow Inputs

Use `workflow.Inputs` for `workflow_dispatch` and `workflow_call` inputs:

```go
var ManualTrigger = workflow.WorkflowDispatchTrigger{
    Inputs: map[string]workflow.Input{
        "environment": {
            Description: "Deployment environment",
            Required:    true,
            Default:     "staging",
            Type:        "choice",
            Options:     []string{"staging", "production"},
        },
    },
}

var DeployStep = workflow.Step{
    Name: "Deploy",
    Run:  "deploy.sh",
    Env: workflow.Env{
        "ENVIRONMENT": workflow.Inputs.Get("environment"),
    },
}
```

### Repository Variables

Use `workflow.Vars` for repository and organization variables:

```go
var ConfigStep = workflow.Step{
    Name: "Configure",
    Env: workflow.Env{
        "API_URL":      workflow.Vars.Get("API_URL"),
        "FEATURE_FLAG": workflow.Vars.Get("FEATURE_FLAG"),
    },
}
```

---

## Condition Builders

wetwire-github-go provides builders for common conditions.

### Branch and Tag Conditions

```go
// Run only on main branch
workflow.Branch("main")
// Generates: github.ref == 'refs/heads/main'

// Run only on specific tag
workflow.Tag("v1.0.0")
// Generates: github.ref == 'refs/tags/v1.0.0'

// Run on any v* tag
workflow.TagPrefix("v")
// Generates: startsWith(github.ref, 'refs/tags/v')
```

### Event Conditions

```go
// Check if push event
workflow.Push()
// Generates: github.event_name == 'push'

// Check if pull request
workflow.PullRequest()
// Generates: github.event_name == 'pull_request'
```

### Status Functions

```go
// Always run (even on failure)
workflow.Always()
// Generates: always()

// Run only on failure
workflow.Failure()
// Generates: failure()

// Run only on success (default)
workflow.Success()
// Generates: success()

// Run only if cancelled
workflow.Cancelled()
// Generates: cancelled()
```

### Using Conditions in Jobs and Steps

```go
var DeployJob = workflow.Job{
    Name:   "deploy",
    RunsOn: "ubuntu-latest",
    If:     workflow.Branch("main"),
    Steps:  DeploySteps,
}

var NotifyStep = workflow.Step{
    Name: "Notify on Failure",
    If:   workflow.Failure(),
    Run:  "notify-failure.sh",
}
```

---

## String Functions

Expression helpers for string operations.

### Contains, StartsWith, EndsWith

```go
// Check if ref contains a string
workflow.Contains(workflow.GitHub.Ref(), workflow.Expression("'feature/'"))

// Check if ref starts with prefix
workflow.StartsWith(workflow.GitHub.RefName(), workflow.Expression("'release-'"))

// Check if branch ends with suffix
workflow.EndsWith(workflow.GitHub.RefName(), workflow.Expression("'-rc'"))
```

### Format and Join

```go
// Format a string
workflow.Format("Build {0} on {1}", workflow.GitHub.SHA(), workflow.Runner.OS())

// Join array elements
workflow.Join(workflow.Expression("matrix.os"), ", ")
```

### JSON Functions

```go
// Convert to JSON
workflow.ToJSON(workflow.GitHub.Event("pull_request"))

// Parse JSON from step output
workflow.FromJSON(workflow.Steps.Get("config", "json"))
```

---

## Complex Expressions

### Combining Expressions

Use `.And()`, `.Or()`, and `.Not()` to combine expressions:

```go
// Deploy only on main branch push
condition := workflow.Branch("main").And(workflow.Push())
// Generates: (github.ref == 'refs/heads/main') && (github.event_name == 'push')

// Run on main or develop
condition := workflow.Branch("main").Or(workflow.Branch("develop"))
// Generates: (github.ref == 'refs/heads/main') || (github.ref == 'refs/heads/develop')

// Skip on push events
condition := workflow.Push().Not()
// Generates: !(github.event_name == 'push')
```

### Chaining Multiple Conditions

```go
var ComplexCondition = workflow.Branch("main").
    And(workflow.Push()).
    Or(workflow.Branch("develop").And(workflow.PullRequest()))

// Generates:
// ((github.ref == 'refs/heads/main') && (github.event_name == 'push')) ||
// ((github.ref == 'refs/heads/develop') && (github.event_name == 'pull_request'))
```

### Real-World Examples

**Deploy to staging on main, production on tags:**

```go
var StagingDeploy = workflow.Job{
    Name:   "staging",
    RunsOn: "ubuntu-latest",
    If:     workflow.Branch("main"),
    Steps:  StagingSteps,
}

var ProductionDeploy = workflow.Job{
    Name:   "production",
    RunsOn: "ubuntu-latest",
    If:     workflow.TagPrefix("v"),
    Steps:  ProductionSteps,
}
```

**Notify on any failure:**

```go
var NotifyJob = workflow.Job{
    Name:   "notify",
    RunsOn: "ubuntu-latest",
    Needs:  []any{Build, Test, Deploy},
    If:     workflow.Failure(),
    Steps:  NotifySteps,
}
```

**Skip CI for documentation changes:**

```go
var Build = workflow.Job{
    Name:   "build",
    RunsOn: "ubuntu-latest",
    If:     workflow.Expression("!contains(github.event.head_commit.message, '[skip ci]')"),
    Steps:  BuildSteps,
}
```

---

## Security Considerations

### WAG017: Explicit Permissions

Always set explicit permissions on workflows to follow the principle of least privilege:

```go
var CI = workflow.Workflow{
    Name: "CI",
    On:   CITriggers,
    Permissions: &workflow.Permissions{
        Contents: "read",
        Packages: "write",
    },
    Jobs: Jobs{"build": Build},
}
```

**Related lint rule**: WAG017 suggests adding explicit `Permissions` field.

### WAG018: pull_request_target Safety

The `pull_request_target` trigger runs with write permissions from the base branch. Combined with checkout, this can expose your repository to code injection:

```go
// DANGEROUS: pull_request_target with checkout
var DangerousTriggers = workflow.Triggers{
    PullRequestTarget: &workflow.PullRequestTargetTrigger{},
}

var DangerousJob = workflow.Job{
    Steps: []any{
        checkout.Checkout{},  // Checks out PR code with write permissions!
    },
}
```

**Safe alternatives:**

1. Use `pull_request` trigger when possible (runs with read-only permissions)
2. If using `pull_request_target`, only run trusted code:

```go
// Safer: Only checkout base branch
var SaferSteps = []any{
    checkout.Checkout{
        Ref: workflow.GitHub.BaseRef(),  // Checkout base, not PR
    },
}
```

**Related lint rule**: WAG018 detects `pull_request_target` with checkout actions.

### Expression Injection Prevention

> **Best Practice**: Use type-safe expression builders instead of raw strings to prevent injection vulnerabilities.

```go
// GOOD: Type-safe builders prevent injection
workflow.Secrets.Get("TOKEN")
workflow.GitHub.Actor()
workflow.Branch("main")

// RISKY: Raw strings can be vulnerable if using untrusted input
workflow.Expression(fmt.Sprintf("github.event.issue.title == '%s'", userInput))
```

### Avoid Hardcoded Expressions

Use condition builders instead of raw expression strings:

```go
// GOOD: Condition builder
If: workflow.Branch("main")

// AVOID: Raw expression string (triggers WAG002 lint warning)
If: "${{ github.ref == 'refs/heads/main' }}"
```

**Related lint rules:**
- WAG002 detects raw expression strings in `If` fields
- WAG008 detects hardcoded `${{ }}` expressions

---

## Runner Context

Use `workflow.Runner` to access runner information:

```go
workflow.Runner.OS()        // runner.os (Linux, Windows, macOS)
workflow.Runner.Arch()      // runner.arch (X86, X64, ARM, ARM64)
workflow.Runner.Name()      // runner.name
workflow.Runner.Temp()      // runner.temp (temp directory)
workflow.Runner.ToolCache() // runner.tool_cache (tool cache directory)
```

### Cross-Platform Steps

```go
var InstallStep = workflow.Step{
    Name: "Install Dependencies",
    Run:  "install.sh",
    Env: workflow.Env{
        "RUNNER_OS": workflow.Runner.OS(),
    },
}
```

---

## See Also

- [Quick Start](QUICK_START.md) - Getting started guide
- [Examples](EXAMPLES.md) - Real-world workflow patterns
- [FAQ](FAQ.md) - Common questions about lint rules
- [Internals](INTERNALS.md) - Expression serialization details
