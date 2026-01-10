# security-workflow

Example security workflow using wetwire-github-go.

## What This Is

A reference implementation showing how to declare GitHub Actions security workflows using typed Go structs. This example creates a comprehensive security workflow with CodeQL analysis, Trivy vulnerability scanning, and SLSA build provenance attestation.

## Key Files

- `workflows/workflows.go` - Main security workflow declaration with explicit permissions
- `workflows/jobs.go` - CodeQL, Trivy, and attestation job definitions
- `workflows/triggers.go` - Push and pull request trigger configurations
- `workflows/steps.go` - Step sequences using typed security action wrappers

## Patterns Used

1. **Flat variables** - All structs are package-level variables, not nested
2. **Typed wrappers** - Uses `codeql_init.CodeQLInit{}`, `codeql_analyze.CodeQLAnalyze{}`, `trivy.Trivy{}`, `attest_build_provenance.AttestBuildProvenance{}`
3. **Explicit permissions** - Demonstrates WAG017 compliance with minimal permission scopes per job
4. **Safe triggers** - Demonstrates WAG018 compliance using `pull_request` (not `pull_request_target`) with checkout
5. **Security best practices** - Multiple security tools with results uploaded to GitHub Security

## Build Command

```bash
wetwire-github build ./workflows
```

Output: `.github/workflows/security.yml`

## Security Features

- **CodeQL**: Static code analysis for security vulnerabilities
- **Trivy**: Container and filesystem vulnerability scanning
- **SLSA**: Build provenance attestation for supply chain security
- **Permissions**: Minimal scopes (read/write only what's needed)
- **Safe checkout**: Uses `pull_request` trigger context
