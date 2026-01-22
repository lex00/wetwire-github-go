<picture>
  <source media="(prefers-color-scheme: dark)" srcset="../../docs/wetwire-dark.svg">
  <img src="../../docs/wetwire-light.svg" width="100" height="67">
</picture>

A complete example demonstrating how to implement manual approval gates between staging and production deployments using wetwire-github-go.

## Features Demonstrated

- **Explicit approval gate job** that blocks production deployment until approved
- **Sequential deployment pipeline** with job dependencies (build -> staging -> approve -> production)
- **Environment protection rules** for manual approval enforcement
- **Separate staging and production environments** with distinct configurations
- **Manual triggering** with `workflow_dispatch`

## Project Structure

```
approval-gates-workflow/
├── go.mod                    # Module with replace directive
├── README.md                 # This file
├── CLAUDE.md                 # AI assistant context
└── workflows/
    ├── workflows.go          # Workflow declarations
    ├── jobs.go               # Job definitions with environments
    ├── triggers.go           # Trigger configurations
    └── steps.go              # Step sequences
```

## Usage

### Generate YAML

```bash
cd examples/approval-gates-workflow
go mod tidy
wetwire-github build .
```

This generates `.github/workflows/approval-gates.yml`.

### View Generated YAML

```bash
cat .github/workflows/approval-gates.yml
```

### Validate with actionlint

```bash
wetwire-github validate .github/workflows/approval-gates.yml
```

## Key Patterns

### Approval Gate Job

The approval gate is a minimal job that uses the production environment to trigger GitHub's approval flow:

```go
var ApprovalGateEnvironment = workflow.Environment{
    Name: "production",
}

var ApproveProduction = workflow.Job{
    Name:        "Approve Production Deployment",
    RunsOn:      "ubuntu-latest",
    Needs:       []any{DeployStaging},
    Environment: &ApprovalGateEnvironment,
    Steps:       ApprovalSteps,
}
```

### Sequential Pipeline with Explicit Gate

The pipeline enforces a strict order with an approval checkpoint:

```go
var BuildJob = workflow.Job{...}                           // Step 1: Build
var DeployStaging = workflow.Job{Needs: []any{BuildJob}}   // Step 2: Deploy staging
var ApproveProduction = workflow.Job{Needs: []any{DeployStaging}}  // Step 3: Wait for approval
var DeployProduction = workflow.Job{Needs: []any{ApproveProduction}}  // Step 4: Deploy production
```

## GitHub Environment Configuration

To enable the approval gate functionality:

1. Go to **Settings** > **Environments**
2. Create a `staging` environment (optional protection rules)
3. Create a `production` environment with:
   - **Required reviewers** - Add team members who must approve production deployments
   - **Wait timer** - Optional delay before deployment starts (e.g., 5 minutes)
   - **Deployment branches** - Restrict to `main` branch only

## Pipeline Flow

```
┌─────────┐     ┌──────────────────┐     ┌─────────────────────┐     ┌─────────────────────┐
│  Build  │ --> │ Deploy to Staging│ --> │ Approve Production  │ --> │ Deploy to Production│
└─────────┘     └──────────────────┘     └─────────────────────┘     └─────────────────────┘
                                                    │
                                          ┌─────────▼─────────┐
                                          │ Manual Approval   │
                                          │ Required (GitHub) │
                                          └───────────────────┘
```

When a push to `main` occurs:
1. **Build** job compiles and tests the application
2. **Deploy to Staging** deploys to staging environment
3. **Approve Production** waits for manual approval via GitHub UI
4. **Deploy to Production** runs only after approval is granted
