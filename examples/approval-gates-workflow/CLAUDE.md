# approval-gates-workflow

Example workflow demonstrating manual approval gates between staging and production deployments.

## What This Is

A reference implementation showing how to enforce manual approval between deployment stages using GitHub environment protection rules. The workflow creates an explicit approval gate job that separates staging from production deployment.

## Key Files

- `workflows/workflows.go` - Main workflow declaration with four-job pipeline
- `workflows/jobs.go` - Build, staging, approval gate, and production job definitions
- `workflows/triggers.go` - Push on main and workflow_dispatch triggers
- `workflows/steps.go` - Step sequences for build and deployment operations

## Patterns Used

1. **Explicit approval gate job** - Separate job using production environment for approval
2. **Job dependencies** - Sequential pipeline: Build -> Staging -> Approve -> Production
3. **Environment protection** - Uses `workflow.Environment{Name: "production"}` for approval gates
4. **Flat variables** - All structs are package-level variables

## Build Command

```bash
wetwire-github build ./workflows
```

Output: `.github/workflows/approval-gates.yml`
