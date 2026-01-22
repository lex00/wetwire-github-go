<picture>
  <source media="(prefers-color-scheme: dark)" srcset="../../../../docs/wetwire-dark.svg">
  <img src="../../../../docs/wetwire-light.svg" width="100" height="67">
</picture>

This project contains GitHub Actions workflow definitions using wetwire-github.

## Workflow Overview

**Name:** CI/CD

**Triggers:**
- Push to `main` branch
- Pull requests to `main` branch

**Jobs:**

1. **Build** - Matrix build and test
   - Go versions: 1.23, 1.24
   - Operating systems: ubuntu-latest, macos-latest
   - Caches Go modules for performance
   - Runs on all triggers

2. **Test** - Run tests with coverage
   - Runs on ubuntu-latest with Go 1.23
   - Generates coverage report
   - Displays coverage summary
   - Caches Go modules for performance

3. **Deploy Staging** - Deploy to staging environment
   - Requires Build and Test jobs to pass
   - Only runs on pushes to `main` branch
   - Auto-deploys without approval

4. **Deploy Production** - Deploy to production with approval
   - Requires Build and Test jobs to pass
   - Only runs on pushes to `main` branch
   - **Requires manual approval** via GitHub environment gate
   - Environment: `production`
   - URL: https://example.com

## Project Structure

```
.
├── go.mod              # Module definition
├── workflows.go        # Main workflow declaration
├── jobs.go            # Job definitions (Build, Test, Deploy)
├── triggers.go        # Trigger configurations
├── helpers.go         # Helper functions (List)
├── cmd/main.go        # Usage instructions
└── README.md          # This file
```

## Building the Workflow

Generate the GitHub Actions YAML:

```bash
wetwire-github build .
```

This creates `.github/workflows/cicd.yml` in your repository.

## Validation

Validate the generated YAML with actionlint:

```bash
wetwire-github validate
```

## Linting

Check code against wetwire-github lint rules:

```bash
wetwire-github lint .
```

## Key Features

- **Type-safe workflow definitions** - All resources are Go structs
- **Flat, declarative syntax** - No function calls, no pointers
- **Direct cross-references** - Jobs reference each other by variable name
- **Matrix strategy** - Test across multiple Go versions and operating systems
- **Caching** - Go module cache for faster builds
- **Environment gates** - Production deployment requires manual approval
- **Conditional execution** - Deploy jobs only run on main branch

## GitHub Environment Setup

To enable the production deployment approval:

1. Go to your repository Settings → Environments
2. Create an environment named `production`
3. Enable "Required reviewers" and add reviewers
4. Optionally set deployment branch rules (e.g., only `main`)

When the workflow runs, production deployment will wait for manual approval from the designated reviewers.

## Customization

To modify the workflow:

1. Edit the Go files (`workflows.go`, `jobs.go`, `triggers.go`)
2. Run `wetwire-github build .` to regenerate YAML
3. Commit both Go files and generated YAML to version control

## Example Deployment Steps

The deployment steps are placeholders. Replace them with your actual deployment commands:

```go
var DeployStagingSteps = []any{
    checkout.Checkout{},
    workflow.Step{
        Name: "Deploy to staging",
        Run:  "./scripts/deploy.sh staging",
        Env: workflow.Env{
            "DEPLOY_TOKEN": workflow.Secrets.Get("STAGING_DEPLOY_TOKEN"),
        },
    },
}
```
