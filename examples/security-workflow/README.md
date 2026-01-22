<picture>
  <source media="(prefers-color-scheme: dark)" srcset="../../docs/wetwire-dark.svg">
  <img src="../../docs/wetwire-light.svg" width="100" height="67">
</picture>

A complete example demonstrating how to define GitHub Actions security workflows with wetwire-github-go.

## Features Demonstrated

- **CodeQL Analysis** - Static code analysis with GitHub's CodeQL
- **Trivy Security Scanning** - Vulnerability scanning for repositories and containers
- **SLSA Build Provenance** - Attestation generation for supply chain security
- **Explicit Permissions** - WAG017 compliance with minimal permission scopes
- **Safe Triggers** - WAG018 compliance using pull_request (not pull_request_target) with checkout

## Project Structure

```
security-workflow/
├── go.mod                    # Module with replace directive
├── README.md                 # This file
└── workflows/
    ├── workflows.go          # Workflow declarations with permissions
    ├── jobs.go               # Job definitions with permission scopes
    ├── triggers.go           # Trigger configurations
    └── steps.go              # Step sequences using security action wrappers
```

## Usage

### Generate YAML

```bash
cd examples/security-workflow
go mod tidy
wetwire-github build .
```

This generates `.github/workflows/security.yml`.

### View Generated YAML

```bash
cat .github/workflows/security.yml
```

### Validate with actionlint

```bash
wetwire-github validate .github/workflows/security.yml
```

### Local Development

When developing wetwire-github-go locally, add a replace directive to go.mod:

```go
replace github.com/lex00/wetwire-github-go => ../..
```

Then run `go mod tidy` before building.

## Key Patterns

### Security Action Wrappers

Instead of raw `uses:` strings, use typed security wrappers:

```go
codeql_init.CodeQLInit{
    Languages: "go",
    Queries: "security-extended",
}
codeql_analyze.CodeQLAnalyze{}
trivy.Trivy{
    ScanType: "fs",
    Format: "sarif",
    Output: "trivy-results.sarif",
}
attest_build_provenance.AttestBuildProvenance{
    SubjectPath: "myapp",
}
```

### Explicit Permissions (WAG017 Compliance)

Define minimal permissions for each job:

```go
var CodeQLPermissions = workflow.Permissions{
    Actions:        workflow.PermissionRead,
    Contents:       workflow.PermissionRead,
    SecurityEvents: workflow.PermissionWrite,
}

var CodeQL = workflow.Job{
    Permissions: &CodeQLPermissions,
    // ...
}
```

### Safe Triggers (WAG018 Compliance)

Use `pull_request` (not `pull_request_target`) with checkout:

```go
var SecurityPullRequest = workflow.PullRequestTrigger{
    Branches: []string{"main"},
}

var SecurityTriggers = workflow.Triggers{
    PullRequest: &SecurityPullRequest,
}
```

This is safe because `pull_request` runs in the context of the PR head, preventing untrusted code execution.

### Flat Variable Structure

Extract all nested structs to package-level variables for clarity:

```go
// Separate variables for each component
var CodeQLSteps = []any{...}
var CodeQLPermissions = workflow.Permissions{...}
var CodeQL = workflow.Job{Steps: CodeQLSteps, Permissions: &CodeQLPermissions}
```

## Security Best Practices

1. **Minimal Permissions** - Each job has only the permissions it needs
2. **Safe Checkouts** - Uses `pull_request` trigger, not `pull_request_target`
3. **Security Scanning** - CodeQL for code analysis, Trivy for vulnerabilities
4. **Supply Chain Security** - SLSA attestation for build provenance
5. **Results Upload** - Scan results uploaded to GitHub Security tab

## Jobs Overview

### CodeQL Analysis

- Initializes CodeQL with Go language support
- Runs security-extended query suite
- Builds the project for accurate analysis
- Uploads results to GitHub Security

### Trivy Security Scan

- Scans the repository filesystem
- Detects CRITICAL and HIGH severity vulnerabilities
- Generates SARIF output
- Uploads results to GitHub Security

### Build with Attestation

- Builds a Go binary
- Generates SLSA build provenance
- Creates signed attestation
- Provides supply chain security guarantees
