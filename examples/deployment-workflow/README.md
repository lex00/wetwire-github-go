<picture>
  <source media="(prefers-color-scheme: dark)" srcset="../../docs/wetwire-dark.svg">
  <img src="../../docs/wetwire-light.svg" width="100" height="67">
</picture>

A complete example demonstrating how to define multi-environment deployment workflows with wetwire-github-go.

## Features Demonstrated

- **Multiple environments** (staging, production) with `workflow.Environment`
- **Sequential deployment pipeline** with job dependencies (build -> staging -> production)
- **Manual approval gates** through GitHub environment protection rules
- **Environment-scoped secrets** using `workflow.Secrets.Get()`
- **Manual triggering** with `workflow_dispatch` and environment selection
- **Typed action wrappers** for checkout and setup-go

## Project Structure

```
deployment-workflow/
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
cd examples/deployment-workflow
go mod tidy
wetwire-github build .
```

This generates `.github/workflows/deploy.yml`.

### View Generated YAML

```bash
cat .github/workflows/deploy.yml
```

### Validate with actionlint

```bash
wetwire-github validate .github/workflows/deploy.yml
```

### Local Development

When developing wetwire-github-go locally, add a replace directive to go.mod:

```go
replace github.com/lex00/wetwire-github-go => ../..
```

Then run `go mod tidy` before building.

## Key Patterns

### Environment Configuration

Define deployment environments with names and URLs:

```go
var ProductionEnvironment = workflow.Environment{
    Name: "production",
    URL:  "https://example.com",
}

var DeployProduction = workflow.Job{
    Name:        "Deploy to Production",
    RunsOn:      "ubuntu-latest",
    Environment: &ProductionEnvironment,
    Steps:       DeployProductionSteps,
}
```

### Environment-Scoped Secrets

Access environment-specific secrets:

```go
var DeployStep = workflow.Step{
    Name: "Deploy",
    Env: map[string]any{
        "DEPLOY_TOKEN": workflow.Secrets.Get("PRODUCTION_DEPLOY_TOKEN"),
    },
    Run: "deploy.sh",
}
```

### Job Dependencies

Create a sequential pipeline where each stage depends on the previous:

```go
var Build = workflow.Job{
    Name:   "Build",
    RunsOn: "ubuntu-latest",
    Steps:  BuildSteps,
}

var DeployStaging = workflow.Job{
    Name:   "Deploy to Staging",
    Needs:  []any{Build},  // Runs after Build
    Steps:  DeployStagingSteps,
}

var DeployProduction = workflow.Job{
    Name:   "Deploy to Production",
    Needs:  []any{DeployStaging},  // Runs after Staging
    Steps:  DeployProductionSteps,
}
```

### Manual Dispatch with Environment Selection

Allow manual triggering with environment choice:

```go
var DeployDispatchInputs = map[string]workflow.WorkflowInput{
    "environment": {
        Description: "Target environment for deployment",
        Required:    true,
        Type:        "choice",
        Options:     []string{"staging", "production"},
        Default:     "staging",
    },
}

var DeployDispatch = workflow.WorkflowDispatchTrigger{
    Inputs: DeployDispatchInputs,
}
```

### Flat Variable Structure

Extract all nested structs to package-level variables for clarity:

```go
// Separate variables for each component
var DeployPush = workflow.PushTrigger{Branches: []string{"main"}}
var DeployTriggers = workflow.Triggers{
    Push:             &DeployPush,
    WorkflowDispatch: &DeployDispatch,
}
```

## GitHub Environment Configuration

For the full deployment experience with manual approval gates, configure environments in your GitHub repository:

1. Go to **Settings** > **Environments**
2. Create `staging` and `production` environments
3. For `production`, add:
   - **Required reviewers** - Team members who must approve deployments
   - **Wait timer** - Optional delay before deployment starts
   - **Deployment branches** - Restrict to `main` branch only
4. Add environment-specific secrets:
   - `STAGING_DEPLOY_TOKEN` for staging
   - `PRODUCTION_DEPLOY_TOKEN` for production

## Pipeline Flow

```
┌─────────┐     ┌──────────────────┐     ┌─────────────────────┐
│  Build  │ --> │ Deploy to Staging │ --> │ Deploy to Production │
└─────────┘     └──────────────────┘     └─────────────────────┘
                                                    │
                                          ┌─────────▼─────────┐
                                          │ Manual Approval   │
                                          │ (via GitHub UI)   │
                                          └───────────────────┘
```

When a push to `main` occurs:
1. **Build** job compiles and tests the application
2. **Deploy to Staging** runs after build succeeds
3. **Deploy to Production** waits for staging and manual approval (if configured)
