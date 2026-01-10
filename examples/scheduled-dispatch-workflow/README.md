# Scheduled and Dispatch Workflow Example

A complete example demonstrating how to define GitHub Actions workflows with schedule and workflow_dispatch triggers using wetwire-github-go.

## Features Demonstrated

- **Schedule trigger** with daily cron pattern (midnight UTC)
- **Workflow dispatch trigger** with typed inputs:
  - `environment`: Choice input (dev/staging/prod)
  - `dry_run`: Boolean input for safe testing
  - `version`: String input for version specification
- **Conditional jobs** that run based on event type
- **Typed action wrappers** for checkout

## Project Structure

```
scheduled-dispatch-workflow/
├── go.mod                    # Module with replace directive
├── README.md                 # This file
├── CLAUDE.md                 # AI assistant context
└── workflows/
    ├── workflows.go          # Workflow declarations
    ├── jobs.go               # Job definitions
    ├── triggers.go           # Trigger configurations
    └── steps.go              # Step sequences
```

## Usage

### Generate YAML

```bash
cd examples/scheduled-dispatch-workflow
go mod tidy
wetwire-github build ./workflows
```

This generates `.github/workflows/scheduled-dispatch.yml`.

### View Generated YAML

```bash
cat .github/workflows/scheduled-dispatch.yml
```

### Validate with actionlint

```bash
wetwire-github validate .github/workflows/scheduled-dispatch.yml
```

### Local Development

When developing wetwire-github-go locally, add a replace directive to go.mod:

```go
replace github.com/lex00/wetwire-github-go => ../..
```

Then run `go mod tidy` before building.

## Key Patterns

### Schedule Trigger

Define scheduled runs using cron syntax:

```go
var DailySchedule = workflow.ScheduleTrigger{
    Cron: "0 0 * * *", // Daily at midnight UTC
}

var WorkflowTriggers = workflow.Triggers{
    Schedule: []workflow.ScheduleTrigger{DailySchedule},
}
```

### Workflow Dispatch with Inputs

Define manual triggers with typed inputs:

```go
var DispatchInputs = map[string]workflow.WorkflowInput{
    "environment": {
        Description: "Target deployment environment",
        Type:        "choice",
        Options:     []string{"dev", "staging", "prod"},
        Default:     "dev",
    },
    "dry_run": {
        Type:    "boolean",
        Default: false,
    },
    "version": {
        Type:    "string",
        Default: "latest",
    },
}

var ManualDispatch = workflow.WorkflowDispatchTrigger{
    Inputs: DispatchInputs,
}
```

### Conditional Jobs

Run jobs based on trigger type:

```go
var Maintenance = workflow.Job{
    If: "${{ github.event_name == 'schedule' }}",
    // ...
}

var Deploy = workflow.Job{
    If: "${{ github.event_name == 'workflow_dispatch' }}",
    // ...
}
```

### Input References in Steps

Reference dispatch inputs in step commands:

```go
workflow.Step{
    Run: "echo \"Deploying ${{ inputs.version }} to ${{ inputs.environment }}\"",
}
```

## Common Cron Patterns

| Pattern | Description |
|---------|-------------|
| `0 0 * * *` | Daily at midnight |
| `0 */6 * * *` | Every 6 hours |
| `0 0 * * 0` | Weekly on Sunday |
| `0 0 1 * *` | Monthly on the 1st |
| `*/15 * * * *` | Every 15 minutes |
