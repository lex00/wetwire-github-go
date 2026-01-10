# Reusable Workflow Example

A complete example demonstrating how to define reusable GitHub Actions workflows with wetwire-github-go.

## Features Demonstrated

- **Reusable workflow** with workflow_call trigger
- **Typed inputs** for passing parameters to reusable workflows
- **Typed outputs** for returning values from reusable workflows
- **Secrets handling** for secure credential passing
- **Caller workflow** that invokes the reusable workflow

## Project Structure

```
reusable-workflow/
├── go.mod                    # Module with replace directive
├── README.md                 # This file
├── CLAUDE.md                 # AI assistant context
└── workflows/
    ├── workflows.go          # Workflow declarations (reusable + caller)
    ├── jobs.go               # Job definitions with outputs
    ├── triggers.go           # Trigger configurations including workflow_call
    └── steps.go              # Step sequences
```

## Usage

### Generate YAML

```bash
cd examples/reusable-workflow
go mod tidy
wetwire-github build .
```

This generates:
- `.github/workflows/build-reusable.yml` - The reusable workflow
- `.github/workflows/ci-caller.yml` - The caller workflow

### View Generated YAML

```bash
cat .github/workflows/build-reusable.yml
cat .github/workflows/ci-caller.yml
```

### Validate with actionlint

```bash
wetwire-github validate .github/workflows/*.yml
```

### Local Development

When developing wetwire-github-go locally, add a replace directive to go.mod:

```go
replace github.com/lex00/wetwire-github-go => ../..
```

Then run `go mod tidy` before building.

## Key Patterns

### Reusable Workflow with workflow_call

Define a workflow that can be called by other workflows:

```go
var ReusableWorkflowCall = workflow.WorkflowCallTrigger{
    Inputs: map[string]workflow.WorkflowInput{
        "go_version": {
            Description: "Go version to use",
            Required:    true,
            Type:        "string",
        },
    },
    Outputs: map[string]workflow.WorkflowOutput{
        "artifact_name": {
            Description: "Name of the built artifact",
            Value:       "${{ jobs.build.outputs.artifact }}",
        },
    },
    Secrets: map[string]workflow.WorkflowSecret{
        "deploy_token": {
            Description: "Token for deployment",
            Required:    false,
        },
    },
}
```

### Calling a Reusable Workflow

Use the `uses` field with a path to another workflow:

```go
var CallBuild = workflow.Job{
    Name: "Call Build",
    Uses: "./.github/workflows/build-reusable.yml",
    With: map[string]any{
        "go_version": "1.24",
    },
    Secrets: "inherit",
}
```

### Job Outputs

Define outputs from a job to pass to other jobs or workflow outputs:

```go
var Build = workflow.Job{
    Name:   "Build",
    RunsOn: "ubuntu-latest",
    Outputs: map[string]any{
        "artifact": "${{ steps.build.outputs.name }}",
    },
    Steps: BuildSteps,
}
```

### Flat Variable Structure

Extract all nested structs to package-level variables for clarity:

```go
// Separate variables for each component
var ReusableInputs = map[string]workflow.WorkflowInput{...}
var ReusableOutputs = map[string]workflow.WorkflowOutput{...}
var ReusableWorkflowCall = workflow.WorkflowCallTrigger{
    Inputs:  ReusableInputs,
    Outputs: ReusableOutputs,
}
```
