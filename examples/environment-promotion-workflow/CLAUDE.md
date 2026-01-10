# environment-promotion-workflow

Example workflow demonstrating environment promotion from dev to staging to production.

## What This Is

A reference implementation showing a progressive promotion pipeline where code moves through dev, staging, and production environments. Each promotion stage uses GitHub environment protection rules to control the flow.

## Key Files

- `workflows/workflows.go` - Main workflow declaration with three-stage promotion
- `workflows/jobs.go` - Dev, staging, and production deployment job definitions
- `workflows/triggers.go` - Push on main and workflow_dispatch with environment input
- `workflows/steps.go` - Step sequences with environment-specific configurations

## Patterns Used

1. **Progressive promotion** - Dev -> Staging -> Production pipeline
2. **Environment inputs** - Uses `workflow_dispatch` with environment selection input
3. **Environment protection** - Each environment can have its own approval rules
4. **Environment-specific config** - Different configurations per deployment target
5. **Flat variables** - All structs are package-level variables

## Build Command

```bash
wetwire-github build ./workflows
```

Output: `.github/workflows/environment-promotion.yml`
