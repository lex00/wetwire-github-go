# Environment Promotion Workflow Example

A complete example demonstrating how to implement environment promotion (dev -> staging -> production) using wetwire-github-go.

## Features Demonstrated

- **Three-tier environment promotion** (dev, staging, production)
- **Automatic dev deployment** on push to main
- **Environment-gated promotions** using GitHub environment protection rules
- **Workflow dispatch with environment input** for manual promotion control
- **Environment-specific configuration** for each deployment stage

## Project Structure

```
environment-promotion-workflow/
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
cd examples/environment-promotion-workflow
go mod tidy
wetwire-github build .
```

This generates `.github/workflows/environment-promotion.yml`.

### View Generated YAML

```bash
cat .github/workflows/environment-promotion.yml
```

### Validate with actionlint

```bash
wetwire-github validate .github/workflows/environment-promotion.yml
```

## Key Patterns

### Environment Definitions

Each environment has its own configuration:

```go
var DevEnvironment = workflow.Environment{
    Name: "dev",
    URL:  "https://dev.example.com",
}

var StagingEnvironment = workflow.Environment{
    Name: "staging",
    URL:  "https://staging.example.com",
}

var ProductionEnvironment = workflow.Environment{
    Name: "production",
    URL:  "https://example.com",
}
```

### Progressive Promotion Chain

Jobs form a promotion chain where each stage depends on the previous:

```go
var DeployDev = workflow.Job{...}  // Automatic on push
var PromoteStaging = workflow.Job{Needs: []any{DeployDev}, Environment: &StagingEnvironment}
var PromoteProduction = workflow.Job{Needs: []any{PromoteStaging}, Environment: &ProductionEnvironment}
```

### Workflow Dispatch with Environment Input

Manual triggering allows environment selection:

```go
var PromotionDispatchInputs = map[string]workflow.WorkflowInput{
    "environment": {
        Description: "Target environment for promotion",
        Required:    true,
        Type:        "choice",
        Options:     []string{"dev", "staging", "production"},
        Default:     "dev",
    },
}
```

### Environment-Specific Steps

Each environment uses tailored deployment configuration:

```go
var DeployDevStep = workflow.Step{
    Name: "Deploy to dev",
    Env: map[string]any{
        "DEPLOY_TOKEN": workflow.Secrets.Get("DEV_DEPLOY_TOKEN"),
        "ENVIRONMENT":  "dev",
        "API_URL":      "https://api-dev.example.com",
    },
    Run: "deploy.sh",
}
```

## GitHub Environment Configuration

To enable environment-gated promotions:

1. Go to **Settings** > **Environments**
2. Create environments: `dev`, `staging`, `production`
3. Configure protection rules:
   - **dev**: No protection (automatic deployment)
   - **staging**: Optional reviewers or wait timer
   - **production**: Required reviewers and deployment branch restrictions
4. Add environment-specific secrets for each environment

## Pipeline Flow

```
┌─────────────┐     ┌────────────────────┐     ┌─────────────────────────┐
│ Deploy Dev  │ --> │ Promote to Staging │ --> │ Promote to Production   │
│ (automatic) │     │ (staging env gate) │     │ (production env gate)   │
└─────────────┘     └────────────────────┘     └─────────────────────────┘
       │                      │                            │
       ▼                      ▼                            ▼
  No approval         Optional approval           Required approval
  required            (if configured)             (recommended)
```

When a push to `main` occurs:
1. **Deploy Dev** automatically deploys to development environment
2. **Promote to Staging** waits for staging environment approval (if configured)
3. **Promote to Production** waits for production environment approval
