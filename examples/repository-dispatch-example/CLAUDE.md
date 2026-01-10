# repository-dispatch-example

Example workflow using repository_dispatch trigger for API-triggered deployments.

## What This Is

A reference implementation showing how to declare GitHub Actions workflows that respond to repository dispatch events. This example creates an API-triggered deployment workflow with event type routing.

## Key Files

- `workflows/workflows.go` - Main workflow declaration
- `workflows/jobs.go` - Validate and Deploy job definitions with conditional execution
- `workflows/triggers.go` - RepositoryDispatch trigger with event type filtering
- `workflows/steps.go` - Step sequences with client_payload access

## Patterns Used

1. **Flat variables** - All structs are package-level variables, not nested
2. **Typed wrappers** - Uses `checkout.Checkout{}` and other typed actions
3. **RepositoryDispatch trigger** - Responds to API-triggered events
4. **Event type filtering** - Multiple event types (deploy, deploy-staging, deploy-production)
5. **Client payload access** - Uses `github.event.client_payload.*` for custom parameters
6. **Conditional jobs** - Different jobs run based on `github.event.action`
7. **Job outputs** - Passes validated payload data between jobs

## Build Command

```bash
wetwire-github build ./workflows
```

Output: `.github/workflows/api-triggered-deploy.yml`
