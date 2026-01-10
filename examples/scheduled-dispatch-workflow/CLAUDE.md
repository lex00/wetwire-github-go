# scheduled-dispatch-workflow

Example workflow using scheduled and manual dispatch triggers.

## What This Is

A reference implementation showing how to declare GitHub Actions workflows with schedule and workflow_dispatch triggers using typed Go structs. This example creates a workflow that runs daily on schedule and can be manually triggered with typed inputs.

## Key Files

- `workflows/workflows.go` - Main workflow declaration
- `workflows/jobs.go` - Maintenance and Deploy job definitions
- `workflows/triggers.go` - Schedule and workflow_dispatch trigger configurations
- `workflows/steps.go` - Step sequences for each job

## Patterns Used

1. **Flat variables** - All structs are package-level variables, not nested
2. **Typed wrappers** - Uses `checkout.Checkout{}` and other typed actions
3. **Schedule trigger** - Daily cron pattern for automated maintenance
4. **Workflow dispatch** - Manual trigger with typed inputs (choice, boolean, string)
5. **Conditional jobs** - Jobs that run based on trigger type or input values

## Build Command

```bash
wetwire-github build ./workflows
```

Output: `.github/workflows/scheduled-dispatch.yml`
