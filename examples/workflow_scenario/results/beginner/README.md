<picture>
  <source media="(prefers-color-scheme: dark)" srcset="../../../../docs/wetwire-dark.svg">
  <img src="../../../../docs/wetwire-light.svg" width="100" height="67">
</picture>

This project contains a complete CI/CD pipeline configuration for a Go web application using wetwire-github-go.

## Features

✅ **Multi-version Testing**
- Tests on Go 1.23 and 1.24
- Runs on Ubuntu and macOS

✅ **Quality Checks**
- Build verification
- Test suite with coverage
- Linting with golangci-lint

✅ **Automated Deployment**
- Auto-deploy to staging on main branch
- Manual approval gate for production
- Environment-specific secrets

## Workflow Structure

### Jobs

1. **Build** - Builds the application on matrix of Go versions and OS
2. **Test** - Runs tests with race detector and coverage
3. **Lint** - Runs golangci-lint checks
4. **Deploy Staging** - Deploys to staging (only after build, test, and lint pass)
5. **Deploy Production** - Deploys to production (requires manual approval)

### Job Dependencies

```
Build ─┐
Test ──┼─→ Deploy Staging ─→ Deploy Production
Lint ──┘
```

## Usage

### Generate Workflow Files

```bash
wetwire-github build .
```

This creates `.github/workflows/ci-cd-pipeline.yml`

### Validate Configuration

```bash
# Lint the Go code
wetwire-github lint .

# Validate generated YAML
wetwire-github validate .
```

### Run via MCP

```bash
wetwire-github mcp
```

## Required Secrets

Configure these secrets in your GitHub repository:

- `STAGING_DEPLOY_TOKEN` - Token for staging deployment
- `PRODUCTION_DEPLOY_TOKEN` - Token for production deployment

## Environment Configuration

The workflow uses GitHub Environments:

- **Staging** - Auto-deploys on main branch push
- **Production** - Requires manual approval, accessible at https://example.com

To set up manual approval:
1. Go to Settings → Environments → production
2. Enable "Required reviewers"
3. Add team members who can approve deployments

## Customization

### Change Go Versions

Edit `jobs.go`:

```go
var BuildMatrix = workflow.Matrix{
    Values: map[string][]any{
        "go": {"1.23", "1.24", "1.25"},  // Add more versions
        "os": {"ubuntu-latest", "macos-latest"},
    },
}
```

### Add More Operating Systems

```go
var BuildMatrix = workflow.Matrix{
    Values: map[string][]any{
        "go": {"1.23", "1.24"},
        "os": {"ubuntu-latest", "macos-latest", "windows-latest"},
    },
}
```

### Modify Deployment Steps

Edit the deployment steps in `jobs.go`:

```go
var DeployProductionSteps = []any{
    checkout.Checkout{},
    workflow.Step{
        Name: "Deploy to Production",
        Run:  "./deploy.sh production",  // Your deployment script
    },
}
```

## Project Structure

```
.
├── go.mod              # Go module definition
├── workflows.go        # Main workflow definition
├── jobs.go             # Job configurations
├── triggers.go         # Workflow triggers
├── helpers.go          # Helper functions
├── cmd/
│   └── main.go         # Usage instructions
└── README.md           # This file
```

## Key Design Principles

Following wetwire-github-go conventions:

- **Flat variables** - All nested structs extracted to named variables
- **No pointers** - Pure struct literals only
- **Direct references** - Jobs reference each other by name
- **Type-safe** - Compile-time checking of workflow configuration

## Next Steps

1. Generate the workflow: `wetwire-github build .`
2. Commit to your repository
3. Configure secrets in GitHub Settings
4. Set up production environment with required reviewers
5. Push to main branch to trigger the workflow

## Learn More

- [wetwire-github-go Documentation](https://github.com/wetwire/wetwire-github-go)
- [GitHub Actions Documentation](https://docs.github.com/actions)
- [Environment Protection Rules](https://docs.github.com/actions/deployment/targeting-different-environments/using-environments-for-deployment)
