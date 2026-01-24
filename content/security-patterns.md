---
title: "Security Patterns"
---

This guide covers security best practices for GitHub Actions workflows using wetwire-github-go. It explains how to configure minimal permissions, handle secrets safely, avoid dangerous trigger patterns, and integrate security scanning tools.

**Contents:**
- [Introduction](#introduction)
- [Permission Scoping](#permission-scoping)
- [Secrets Handling](#secrets-handling)
- [Pull Request Target Safety](#pull-request-target-safety)
- [Supply Chain Security](#supply-chain-security)
- [Security Scanning](#security-scanning)
- [Best Practices Checklist](#best-practices-checklist)

---

## Introduction

GitHub Actions workflows run with access to your repository, secrets, and permissions to modify code, packages, and deployments. A misconfigured workflow can expose secrets, allow code injection from malicious pull requests, or give attackers write access to your repository.

wetwire-github-go provides several features to help you write secure workflows:

1. **Type-safe secrets access** - The `workflow.Secrets` context prevents accidental secret exposure
2. **Explicit permissions** - The `workflow.Permissions{}` struct encourages minimal permission scoping
3. **Lint rules** - WAG003, WAG017, and WAG018 detect common security issues
4. **Security action wrappers** - Type-safe wrappers for CodeQL, Trivy, Scorecard, and more

### Related Lint Rules

| Rule | Description |
|------|-------------|
| WAG003 | Detects hardcoded secrets (tokens matching `ghp_`, `ghs_`, `github_pat_`, etc.) |
| WAG017 | Suggests adding explicit `Permissions` field to workflows |
| WAG018 | Detects dangerous `pull_request_target` patterns with checkout actions |

---

## Permission Scoping

By default, GitHub Actions workflows receive a `GITHUB_TOKEN` with broad permissions. The principle of least privilege dictates that workflows should only have the permissions they need.

### Minimal GITHUB_TOKEN Permissions

Always set explicit permissions at the workflow or job level:

```go
import "github.com/lex00/wetwire-github-go/workflow"

var SecurityPermissions = workflow.Permissions{
    Contents:       workflow.PermissionRead,
    SecurityEvents: workflow.PermissionWrite,
}

var Security = workflow.Workflow{
    Name:        "Security",
    On:          SecurityTriggers,
    Permissions: &SecurityPermissions,
    Jobs: map[string]workflow.Job{
        "scan": SecurityScanJob,
    },
}
```

### Per-Job Permission Blocks

For workflows with multiple jobs that have different requirements, set permissions at the job level:

```go
var CodeQLPermissions = workflow.Permissions{
    Actions:        workflow.PermissionRead,
    Contents:       workflow.PermissionRead,
    SecurityEvents: workflow.PermissionWrite,
}

var CodeQL = workflow.Job{
    Name:        "CodeQL Analysis",
    RunsOn:      "ubuntu-latest",
    Permissions: &CodeQLPermissions,
    Steps:       CodeQLSteps,
}
```

### Read-Only vs Write Permissions

Use `workflow.PermissionRead`, `workflow.PermissionWrite`, or `workflow.PermissionNone`:

```go
// Read-only CI workflow
var CIPermissions = workflow.Permissions{
    Contents: workflow.PermissionRead,
}

// Release workflow - needs write access
var ReleasePermissions = workflow.Permissions{
    Contents: workflow.PermissionWrite,
    Packages: workflow.PermissionWrite,
}

// SLSA attestation - needs id-token for OIDC
var AttestPermissions = workflow.Permissions{
    Contents: workflow.PermissionRead,
    IDToken:  workflow.PermissionWrite,
}
```

**Related lint rule**: WAG017 detects workflows missing explicit `Permissions` and suggests adding them.

---

## Secrets Handling

Secrets should never be hardcoded in workflow files. wetwire-github-go provides the `workflow.Secrets` context for type-safe secret access.

### Using workflow.Secrets Context

```go
import "github.com/lex00/wetwire-github-go/workflow"

var DeployStep = workflow.Step{
    Name: "Deploy to Production",
    Run:  "deploy.sh",
    Env: workflow.Env{
        "DEPLOY_TOKEN": workflow.Secrets.Get("DEPLOY_TOKEN"),
        "GH_TOKEN":     workflow.Secrets.GITHUB_TOKEN(),
    },
}
```

### Environment-Scoped Secrets

For sensitive deployments, use GitHub Environments to scope secrets:

```go
var ProductionDeploy = workflow.Job{
    Name:        "Deploy to Production",
    RunsOn:      "ubuntu-latest",
    Environment: "production",  // Secrets scoped to this environment
    Steps: []any{
        workflow.Step{
            Env: workflow.Env{
                "TOKEN": workflow.Secrets.Get("PROD_TOKEN"),
            },
        },
    },
}
```

### Never Hardcode Secrets

The WAG003 lint rule detects hardcoded tokens:

```go
// BAD: Hardcoded token (triggers WAG003 error)
"TOKEN": "ghp_xxxxxxxxxxxxxxxxxxxx"

// GOOD: Use secrets context
"TOKEN": workflow.Secrets.Get("DEPLOY_TOKEN")
```

WAG003 detects patterns: `ghp_`, `ghs_`, `ghu_`, `ghr_`, `github_pat_`

### GitHub App Tokens vs PATs

For elevated permissions or cross-repository access, prefer GitHub App tokens:

```go
import "github.com/lex00/wetwire-github-go/actions/create_github_app_token"

var CreateTokenStep = create_github_app_token.CreateGithubAppToken{
    AppID:      workflow.Secrets.Get("APP_ID"),
    PrivateKey: workflow.Secrets.Get("APP_PRIVATE_KEY"),
}
```

Benefits: fine-grained permissions, automatic expiration (1 hour), audit trail, no personal account dependency.

---

## Pull Request Target Safety

The `pull_request_target` trigger runs workflows in the context of the base branch with write permissions, even for PRs from forks. This is dangerous when combined with checkout actions.

### Why pull_request_target is Dangerous

When a workflow uses `pull_request_target` and checks out PR code, it runs untrusted code with elevated permissions:

```go
// DANGEROUS: Runs untrusted PR code with write permissions
var DangerousTriggers = workflow.Triggers{
    PullRequestTarget: &workflow.PullRequestTargetTrigger{},
}

var DangerousJob = workflow.Job{
    RunsOn: "ubuntu-latest",
    Steps: []any{
        checkout.Checkout{},           // Checks out PR's code
        workflow.Step{Run: "make build"},  // Untrusted code runs with write permissions
    },
}
```

An attacker could modify their PR to exfiltrate secrets, modify repository contents, or push malicious code.

### Safe Patterns for Fork PRs

**Pattern 1: Use `pull_request` trigger** (recommended)

```go
var SafeTriggers = workflow.Triggers{
    PullRequest: &workflow.PullRequestTrigger{
        Branches: []string{"main"},
    },
}
```

**Pattern 2: Only checkout base branch**

```go
var SaferSteps = []any{
    checkout.Checkout{
        Ref: workflow.GitHub.BaseRef(),  // Checkout base, not PR
    },
}
```

**Pattern 3: No checkout - use GitHub API only**

```go
var LabelJob = workflow.Job{
    Steps: []any{
        workflow.Step{
            Run: "gh pr edit $PR --add-label needs-review",
            Env: workflow.Env{
                "GH_TOKEN": workflow.Secrets.GITHUB_TOKEN(),
                "PR":       workflow.GitHub.Event("pull_request.number"),
            },
        },
    },
}
```

### What WAG018 Detects

WAG018 detects `pull_request_target` triggers combined with checkout actions:

```
workflow.go:15:1: warning: Workflow uses pull_request_target with checkout action -
potential security risk. [WAG018]
```

---

## Supply Chain Security

Protect your software supply chain with build provenance, artifact signing, and dependency review.

### SLSA Build Provenance

Generate SLSA provenance attestations using the `attest_build_provenance` wrapper:

```go
import (
    "github.com/lex00/wetwire-github-go/actions/attest_build_provenance"
    "github.com/lex00/wetwire-github-go/workflow"
)

var AttestPermissions = workflow.Permissions{
    Contents: workflow.PermissionRead,
    IDToken:  workflow.PermissionWrite,
}

var BuildAttestSteps = []any{
    checkout.Checkout{},
    setup_go.SetupGo{GoVersion: "1.24"},
    workflow.Step{
        Name: "Build binary",
        Run:  "go build -o myapp ./...",
    },
    attest_build_provenance.AttestBuildProvenance{
        SubjectPath: "myapp",
    },
}
```

For container images:

```go
var ContainerAttestStep = attest_build_provenance.AttestBuildProvenance{
    SubjectName:    "ghcr.io/myorg/myapp",
    SubjectDigest:  "sha256:abc123...",
    PushToRegistry: true,
}
```

### Container Signing with Cosign

Sign container images using the `cosign_installer` wrapper:

```go
import "github.com/lex00/wetwire-github-go/actions/cosign_installer"

var SignContainerSteps = []any{
    cosign_installer.CosignInstaller{CosignRelease: "v2.2.0"},
    workflow.Step{
        Name: "Sign container",
        Run:  "cosign sign --yes ghcr.io/myorg/myapp:latest",
    },
}

var SignPermissions = workflow.Permissions{
    Contents: workflow.PermissionRead,
    Packages: workflow.PermissionWrite,
    IDToken:  workflow.PermissionWrite,  // Required for keyless signing
}
```

### Dependency Review

Block PRs that introduce vulnerable dependencies using the `dependency_review` wrapper:

```go
import "github.com/lex00/wetwire-github-go/actions/dependency_review"

var DependencyReviewSteps = []any{
    checkout.Checkout{},
    dependency_review.DependencyReview{
        FailOnSeverity:     "high",
        DenyLicenses:       "GPL-3.0,AGPL-3.0",
        CommentSummaryInPR: true,
    },
}
```

---

## Security Scanning

Integrate security scanning tools to identify vulnerabilities in code and containers.

### CodeQL Setup

Use the `codeql_init` and `codeql_analyze` wrappers for semantic code analysis:

```go
import (
    "github.com/lex00/wetwire-github-go/actions/codeql_init"
    "github.com/lex00/wetwire-github-go/actions/codeql_analyze"
)

var CodeQLSteps = []any{
    checkout.Checkout{},
    codeql_init.CodeQLInit{
        Languages: "go",
        Queries:   "security-extended",
    },
    workflow.Step{Name: "Build", Run: "go build ./..."},
    codeql_analyze.CodeQLAnalyze{},
}

var CodeQLPermissions = workflow.Permissions{
    Actions:        workflow.PermissionRead,
    Contents:       workflow.PermissionRead,
    SecurityEvents: workflow.PermissionWrite,
}
```

### Trivy Container Scanning

Use the `trivy` wrapper to scan containers and filesystems:

```go
import (
    "github.com/lex00/wetwire-github-go/actions/trivy"
    "github.com/lex00/wetwire-github-go/actions/upload_sarif"
)

var TrivySteps = []any{
    checkout.Checkout{},
    trivy.Trivy{
        ScanType: "fs",
        Format:   "sarif",
        Output:   "trivy-results.sarif",
        Severity: "CRITICAL,HIGH",
    },
    upload_sarif.UploadSarif{
        SarifFile: "trivy-results.sarif",
        Category:  "trivy",
    },
}

var TrivyPermissions = workflow.Permissions{
    Contents:       workflow.PermissionRead,
    SecurityEvents: workflow.PermissionWrite,
}
```

### SARIF Upload

Use the `upload_sarif` wrapper to upload results from any security tool:

```go
import "github.com/lex00/wetwire-github-go/actions/upload_sarif"

var UploadResults = upload_sarif.UploadSarif{
    SarifFile:         "results.sarif",
    Category:          "my-scanner",
    WaitForProcessing: true,
}
```

### OpenSSF Scorecard

Use the `scorecard` wrapper to assess your repository's security posture:

```go
import "github.com/lex00/wetwire-github-go/actions/scorecard"

var ScorecardSteps = []any{
    checkout.Checkout{PersistCredentials: false},
    scorecard.Scorecard{
        ResultsFile:    "scorecard-results.sarif",
        ResultsFormat:  "sarif",
        PublishResults: true,
        RepoToken:      workflow.Secrets.GITHUB_TOKEN(),
    },
    upload_sarif.UploadSarif{
        SarifFile: "scorecard-results.sarif",
        Category:  "scorecard",
    },
}

var ScorecardPermissions = workflow.Permissions{
    Contents:       workflow.PermissionRead,
    SecurityEvents: workflow.PermissionWrite,
    IDToken:        workflow.PermissionWrite,
}
```

### Complete Security Workflow Example

See the [examples/security-workflow/](../examples/security-workflow/) directory for a complete working example combining CodeQL, Trivy, and SLSA attestation.

---

## Best Practices Checklist

### Permissions
- [ ] Workflow has explicit `Permissions` field (WAG017)
- [ ] Each job only has permissions it needs
- [ ] Read-only permissions where writes are not required
- [ ] `id-token: write` only when OIDC is needed

### Secrets
- [ ] No hardcoded tokens in code (WAG003)
- [ ] Secrets accessed via `workflow.Secrets.Get()`
- [ ] Sensitive secrets scoped to environments
- [ ] GitHub App tokens preferred over PATs

### Triggers
- [ ] Avoid `pull_request_target` with checkout (WAG018)
- [ ] Use `pull_request` for PR CI workflows
- [ ] If `pull_request_target` needed, only checkout base branch

### Supply Chain
- [ ] SLSA provenance for release artifacts
- [ ] Container images signed with cosign
- [ ] Dependency review on pull requests
- [ ] Pin action versions to SHA or version tag

### Scanning
- [ ] CodeQL enabled for supported languages
- [ ] Container scanning with Trivy or equivalent
- [ ] SARIF results uploaded to GitHub Security
- [ ] Scorecard running periodically

### General
- [ ] Run `wetwire-github lint` before committing
- [ ] Review WAG003, WAG017, WAG018 warnings
- [ ] Workflows tested with fork PRs

---

## See Also

- [Expression Contexts](EXPRESSIONS.md) - Type-safe access to secrets, matrix, and more
- [Examples](../examples/security-workflow/) - Complete security workflow example
- [FAQ](FAQ.md) - Common questions about lint rules
- [CLI Documentation](CLI.md) - Running lint and validation commands
