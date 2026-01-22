<picture>
  <source media="(prefers-color-scheme: dark)" srcset="../../docs/wetwire-dark.svg">
  <img src="../../docs/wetwire-light.svg" width="100" height="67">
</picture>

A complete example demonstrating how to define GitHub Actions workflows triggered by API calls using wetwire-github-go.

## Features Demonstrated

- **RepositoryDispatchTrigger** with event_type filtering
- **Event type routing** for different deployment targets:
  - `deploy` - Generic deployment with validation
  - `deploy-staging` - Staging-specific deployment
  - `deploy-production` - Production with approval requirement
- **Client payload access** via `github.event.client_payload`
- **Conditional jobs** based on `github.event.action`
- **Job outputs** for passing validated data between jobs

## Project Structure

```
repository-dispatch-example/
├── go.mod                    # Module with replace directive
├── README.md                 # This file
├── CLAUDE.md                 # AI assistant context
└── workflows/
    ├── workflows.go          # Workflow declarations
    ├── jobs.go               # Validate and Deploy job definitions
    ├── triggers.go           # RepositoryDispatch trigger configuration
    └── steps.go              # Step sequences with payload handling
```

## Usage

### Generate YAML

```bash
cd examples/repository-dispatch-example
go mod tidy
wetwire-github build ./workflows
```

This generates `.github/workflows/api-triggered-deploy.yml`.

### View Generated YAML

```bash
cat .github/workflows/api-triggered-deploy.yml
```

### Validate with actionlint

```bash
wetwire-github validate .github/workflows/api-triggered-deploy.yml
```

## Triggering the Workflow

### Using GitHub CLI

```bash
# Generic deploy
gh api repos/{owner}/{repo}/dispatches \
  --method POST \
  -f event_type=deploy \
  -f client_payload='{"environment":"staging","version":"v1.2.3","ref":"main"}'

# Deploy to staging
gh api repos/{owner}/{repo}/dispatches \
  --method POST \
  -f event_type=deploy-staging \
  -f client_payload='{"version":"v1.2.3","features":"new-ui"}'

# Deploy to production (requires approval)
gh api repos/{owner}/{repo}/dispatches \
  --method POST \
  -f event_type=deploy-production \
  -f client_payload='{"version":"v1.2.3","approved_by":"admin-user"}'
```

### Using curl

```bash
curl -X POST \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer $GITHUB_TOKEN" \
  https://api.github.com/repos/{owner}/{repo}/dispatches \
  -d '{"event_type":"deploy","client_payload":{"environment":"staging","version":"v1.2.3"}}'
```

## Key Patterns

### RepositoryDispatch Trigger

Configure a workflow to respond to API dispatch events:

```go
var DeployDispatch = workflow.RepositoryDispatchTrigger{
    Types: []string{"deploy", "deploy-staging", "deploy-production"},
}

var DispatchTriggers = workflow.Triggers{
    RepositoryDispatch: &DeployDispatch,
}
```

### Conditional Jobs by Event Type

Route to different jobs based on the event type:

```go
var DeployStaging = workflow.Job{
    If: "${{ github.event.action == 'deploy-staging' }}",
    // ...
}

var DeployProduction = workflow.Job{
    If: "${{ github.event.action == 'deploy-production' }}",
    // ...
}
```

### Accessing Client Payload

Access custom parameters sent with the dispatch event:

```go
workflow.Step{
    Run: `echo "Environment: ${{ github.event.client_payload.environment }}"
echo "Version: ${{ github.event.client_payload.version }}"
echo "Ref: ${{ github.event.client_payload.ref }}"`,
}
```

### Payload with Default Values

Use GitHub expression syntax for defaults:

```go
checkout.Checkout{
    Ref: "${{ github.event.client_payload.ref || 'main' }}",
}
```

### Job Outputs for Validated Data

Pass validated payload data between jobs:

```go
var ValidateOutputs = map[string]string{
    "environment": "${{ steps.validate.outputs.environment }}",
    "version":     "${{ steps.validate.outputs.version }}",
}

var Validate = workflow.Job{
    Outputs: ValidateOutputs,
    Steps:   ValidateSteps,
}
```

## Client Payload Examples

### Minimal Deployment

```json
{
  "environment": "staging"
}
```

### Full Deployment Configuration

```json
{
  "environment": "production",
  "version": "v2.0.0",
  "ref": "release/2.0",
  "approved_by": "admin-user",
  "features": ["new-ui", "api-v2"],
  "rollback_version": "v1.9.0"
}
```

## Event Context Reference

| Expression | Description |
|------------|-------------|
| `github.event.action` | The event_type from the dispatch |
| `github.event.client_payload` | Full client payload object |
| `github.event.client_payload.<key>` | Specific payload field |
| `github.event.sender.login` | User who triggered the dispatch |
| `github.event.repository.full_name` | Repository (owner/repo) |

## Common Use Cases

1. **External CI/CD integration** - Trigger deployments from external systems
2. **Scheduled deployments** - External scheduler triggers specific deploy times
3. **Cross-repository workflows** - One repo triggers workflows in another
4. **ChatOps** - Slack/Discord bots trigger deployments
5. **Custom automation** - Scripts trigger workflows with parameters
