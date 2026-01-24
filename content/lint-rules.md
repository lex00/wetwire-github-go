---
title: "Lint Rules"
---
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="./wetwire-dark.svg">
  <img src="./wetwire-light.svg" width="100" height="67">
</picture>

wetwire-github-go includes 20 lint rules to enforce best practices for GitHub Actions workflow declarations.

## Quick Reference

| Rule | Description | Severity | Auto-fix |
|------|-------------|----------|----------|
| WAG001 | Use typed action wrappers | warning | Yes |
| WAG002 | Use condition builders | warning | No |
| WAG003 | Use secrets context | error | No |
| WAG004 | Use matrix builder | info | No |
| WAG005 | Extract inline structs | info | No |
| WAG006 | Detect duplicate workflow names | error | No |
| WAG007 | Flag oversized files | warning | No |
| WAG008 | Avoid hardcoded expressions | info | No |
| WAG009 | Validate matrix dimensions | error | No |
| WAG010 | Flag missing recommended inputs | warning | No |
| WAG011 | Detect unreachable jobs | error | No |
| WAG012 | Warn about deprecated versions | warning | No |
| WAG013 | Avoid pointer assignments | error | No |
| WAG014 | Jobs should have timeout | warning | No |
| WAG015 | Suggest caching for setup actions | warning | No |
| WAG016 | Validate concurrency settings | warning | No |
| WAG017 | Suggest explicit permissions | info | No |
| WAG018 | Detect dangerous pull_request_target | warning | No |
| WAG019 | Detect circular dependencies | error | No |
| WAG020 | Detect hardcoded secrets | error | No |

## Rule Details

### WAG001: Use Typed Action Wrappers

**Description:** Use typed action wrappers instead of raw `uses:` strings for better type safety and IDE support.

**Severity:** warning
**Auto-fix:** Yes

#### Bad
```go
var Step = workflow.Step{Uses: "actions/checkout@v4"}
```

#### Good
```go
var Step = checkout.Checkout{}
```

---

### WAG002: Use Condition Builders

**Description:** Use condition builders instead of raw expression strings for the `If` field.

**Severity:** warning
**Auto-fix:** No

#### Bad
```go
var Step = workflow.Step{
    If: "${{ github.ref == 'refs/heads/main' }}",
}
```

#### Good
```go
var Step = workflow.Step{
    If: workflow.Branch("main"),
}
```

---

### WAG003: Use Secrets Context

**Description:** Detect hardcoded GitHub token patterns. Use the secrets context instead.

**Severity:** error
**Auto-fix:** No

#### Bad
```go
var token = "ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
```

#### Good
```go
var Step = workflow.Step{
    Env: workflow.Env{
        "TOKEN": workflow.Secrets.Get("MY_TOKEN"),
    },
}
```

---

### WAG004: Use Matrix Builder

**Description:** Extract inline matrix definitions to named variables for better readability.

**Severity:** info
**Auto-fix:** No

#### Bad
```go
var Job = workflow.Job{
    Strategy: workflow.Strategy{
        Matrix: workflow.Matrix{
            Values: map[string][]any{"os": {"ubuntu", "macos"}},
        },
    },
}
```

#### Good
```go
var BuildMatrix = workflow.Matrix{
    Values: map[string][]any{"os": {"ubuntu", "macos"}},
}

var Job = workflow.Job{
    Strategy: workflow.Strategy{Matrix: BuildMatrix},
}
```

---

### WAG005: Extract Inline Structs

**Description:** Deeply nested structs (depth > 2) should be extracted to named variables.

**Severity:** info
**Auto-fix:** No

#### Bad
```go
var Workflow = workflow.Workflow{
    On: workflow.Triggers{
        Push: workflow.PushTrigger{
            Branches: List("main"),
        },
    },
}
```

#### Good
```go
var PushTrigger = workflow.PushTrigger{Branches: List("main")}
var Triggers = workflow.Triggers{Push: PushTrigger}
var Workflow = workflow.Workflow{On: Triggers}
```

---

### WAG006: Detect Duplicate Workflow Names

**Description:** Each workflow must have a unique name within a file.

**Severity:** error
**Auto-fix:** No

#### Bad
```go
var CI1 = workflow.Workflow{Name: "CI"}
var CI2 = workflow.Workflow{Name: "CI"}  // Duplicate!
```

#### Good
```go
var CI = workflow.Workflow{Name: "CI"}
var Release = workflow.Workflow{Name: "Release"}
```

---

### WAG007: Flag Oversized Files

**Description:** Files with more than 10 jobs (configurable) may be too complex and should be split.

**Severity:** warning
**Auto-fix:** No

---

### WAG008: Avoid Hardcoded Expressions

**Description:** Use expression builders instead of hardcoded `${{ ... }}` strings.

**Severity:** info
**Auto-fix:** No

#### Bad
```go
var Value = "${{ always() }}"
```

#### Good
```go
var Value = workflow.Always()
```

---

### WAG009: Validate Matrix Dimensions

**Description:** Matrix dimensions must have at least one value.

**Severity:** error
**Auto-fix:** No

#### Bad
```go
var Matrix = workflow.Matrix{
    Values: map[string][]any{
        "os": {},  // Empty!
    },
}
```

#### Good
```go
var Matrix = workflow.Matrix{
    Values: map[string][]any{
        "os": {"ubuntu-latest"},
    },
}
```

---

### WAG010: Flag Missing Recommended Inputs

**Description:** Some actions have recommended inputs that should be set for proper operation.

**Severity:** warning
**Auto-fix:** No

#### Bad
```go
var Setup = setup_go.SetupGo{}  // Missing GoVersion
```

#### Good
```go
var Setup = setup_go.SetupGo{GoVersion: "1.23"}
```

---

### WAG011: Detect Unreachable Jobs

**Description:** Jobs that depend on undefined jobs will never run.

**Severity:** error
**Auto-fix:** No

#### Bad
```go
var Deploy = workflow.Job{
    Needs: []any{Build, Test},  // Test is not defined!
}
```

#### Good
```go
var Build = workflow.Job{Name: "build"}
var Test = workflow.Job{Name: "test"}
var Deploy = workflow.Job{
    Needs: []any{Build, Test},
}
```

---

### WAG012: Warn About Deprecated Action Versions

**Description:** Using deprecated action versions may cause security or compatibility issues.

**Severity:** warning
**Auto-fix:** No

#### Bad
```go
var Step = workflow.Step{Uses: "actions/checkout@v2"}
```

#### Good
```go
var Step = checkout.Checkout{}  // Uses latest v4
```

---

### WAG013: Avoid Pointer Assignments

**Description:** wetwire uses value semantics. Avoid `&Type{}` pointer assignments.

**Severity:** error
**Auto-fix:** No

#### Bad
```go
var Job = &workflow.Job{Name: "build"}
```

#### Good
```go
var Job = workflow.Job{Name: "build"}
```

---

### WAG014: Jobs Should Have Timeout

**Description:** Jobs without `TimeoutMinutes` can run indefinitely and block resources.

**Severity:** warning
**Auto-fix:** No

#### Bad
```go
var Job = workflow.Job{
    Name:   "build",
    RunsOn: "ubuntu-latest",
}
```

#### Good
```go
var Job = workflow.Job{
    Name:           "build",
    RunsOn:         "ubuntu-latest",
    TimeoutMinutes: 30,
}
```

---

### WAG015: Suggest Caching for Setup Actions

**Description:** Setup actions (setup-go, setup-node, setup-python) should use caching for faster builds.

**Severity:** warning
**Auto-fix:** No

#### Bad
```go
var Steps = []any{
    setup_go.SetupGo{GoVersion: "1.23"},
}
```

#### Good
```go
var Steps = []any{
    setup_go.SetupGo{GoVersion: "1.23"},
    cache.Cache{Path: "~/.cache/go-build", Key: "go-${{ hashFiles('go.sum') }}"},
}
```

---

### WAG016: Validate Concurrency Settings

**Description:** `CancelInProgress: true` requires a `Group` to be effective.

**Severity:** warning
**Auto-fix:** No

#### Bad
```go
var Concurrency = workflow.Concurrency{
    CancelInProgress: true,  // No group!
}
```

#### Good
```go
var Concurrency = workflow.Concurrency{
    Group:            "${{ github.workflow }}-${{ github.ref }}",
    CancelInProgress: true,
}
```

---

### WAG017: Suggest Explicit Permissions

**Description:** Workflows should explicitly declare permissions for security best practices.

**Severity:** info
**Auto-fix:** No

#### Bad
```go
var CI = workflow.Workflow{
    Name: "CI",
    On:   Triggers,
}
```

#### Good
```go
var CI = workflow.Workflow{
    Name:        "CI",
    On:          Triggers,
    Permissions: workflow.Permissions{Contents: "read"},
}
```

---

### WAG018: Detect Dangerous pull_request_target Patterns

**Description:** Using `pull_request_target` with checkout action can be a security risk, as it runs with write permissions on untrusted PR code.

**Severity:** warning
**Auto-fix:** No

#### Bad
```go
var Triggers = workflow.Triggers{
    PullRequestTarget: workflow.PullRequestTargetTrigger{},
}
var Steps = []any{
    checkout.Checkout{},  // Dangerous with pull_request_target!
}
```

#### Good
```go
// Use pull_request instead, or carefully review security implications
var Triggers = workflow.Triggers{
    PullRequest: workflow.PullRequestTrigger{},
}
```

---

### WAG019: Detect Circular Dependencies

**Description:** Circular job dependencies (A -> B -> C -> A) will cause the workflow to fail.

**Severity:** error
**Auto-fix:** No

#### Bad
```go
var JobA = workflow.Job{Needs: []any{JobC}}
var JobB = workflow.Job{Needs: []any{JobA}}
var JobC = workflow.Job{Needs: []any{JobB}}  // Cycle!
```

#### Good
```go
var JobA = workflow.Job{}
var JobB = workflow.Job{Needs: []any{JobA}}
var JobC = workflow.Job{Needs: []any{JobB}}
```

---

### WAG020: Detect Hardcoded Secrets

**Description:** Detects 25+ hardcoded secret patterns including AWS keys, API tokens, private keys, and more.

**Severity:** error
**Auto-fix:** No

**Detected patterns:**
- AWS access keys (AKIA...)
- GitHub tokens (ghp_, ghs_, ghu_, ghr_, gho_, github_pat_)
- Private keys (RSA, EC, DSA, OpenSSH, PGP)
- Stripe keys (sk_live_, sk_test_)
- Slack tokens (xoxb-, xoxp-, xoxa-, xoxr-, xoxs-)
- Google API keys (AIza...)
- Twilio, SendGrid, Mailgun, NPM, PyPI tokens
- JWT tokens
- DigitalOcean, Heroku, Azure credentials

#### Bad
```go
var Config = map[string]string{
    "api_key": "AKIAIOSFODNN7EXAMPLE",
}
```

#### Good
```go
var Step = workflow.Step{
    Env: workflow.Env{
        "AWS_ACCESS_KEY_ID": workflow.Secrets.Get("AWS_ACCESS_KEY_ID"),
    },
}
```

## Usage

### Running the Linter

```bash
wetwire-github lint .
wetwire-github lint ./workflows.go
```

### JSON Output

```bash
wetwire-github lint . --format json
```

### Auto-fix

```bash
wetwire-github lint . --fix
```

Note: Only WAG001 currently supports auto-fix.

## Configuration

The linter uses sensible defaults. Currently, no configuration file is supported, but rule-specific options can be passed programmatically:

```go
// Example: Customize max jobs threshold
linter := linter.NewLinter(
    &linter.WAG007{MaxJobs: 20},
    // ... other rules
)
```
