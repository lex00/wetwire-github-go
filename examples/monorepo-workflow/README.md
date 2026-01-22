<picture>
  <source media="(prefers-color-scheme: dark)" srcset="../../docs/wetwire-dark.svg">
  <img src="../../docs/wetwire-light.svg" width="100" height="67">
</picture>

A complete example demonstrating how to define GitHub Actions workflows for monorepo CI with wetwire-github-go.

## Features Demonstrated

- **Path filter triggers** on PushTrigger and PullRequestTrigger to only run on relevant changes
- **Change detection** using dorny/paths-filter action to determine which services changed
- **Conditional job execution** based on detected changes
- **Parallel service builds** for API (Go), Web (Node.js), and Shared (Go) components
- **Job outputs** for passing change detection results to downstream jobs

## Project Structure

```
monorepo-workflow/
├── go.mod                    # Module with replace directive
├── README.md                 # This file
├── CLAUDE.md                 # AI assistant context
└── workflows/
    ├── workflows.go          # Main workflow declaration
    ├── jobs.go               # Job definitions with conditions
    ├── triggers.go           # Trigger configurations with path filters
    └── steps.go              # Step sequences for each service
```

## Assumed Monorepo Structure

This workflow assumes your monorepo is structured as:

```
your-monorepo/
├── services/
│   ├── api/                  # Go API service
│   │   ├── go.mod
│   │   └── main.go
│   └── web/                  # Node.js web frontend
│       ├── package.json
│       └── src/
├── shared/                   # Shared Go library
│   ├── go.mod
│   └── lib.go
└── .github/
    └── workflows/
        └── monorepo.yml      # Generated workflow
```

## Usage

### Generate YAML

```bash
cd examples/monorepo-workflow
go mod tidy
wetwire-github build .
```

This generates `.github/workflows/monorepo.yml`.

### View Generated YAML

```bash
cat .github/workflows/monorepo.yml
```

### Validate with actionlint

```bash
wetwire-github validate .github/workflows/monorepo.yml
```

## Key Patterns

### Path Filters on Triggers

Limit workflow runs to relevant file changes:

```go
var MonorepoPush = workflow.PushTrigger{
    Branches: []string{"main"},
    Paths: []string{
        "services/api/**",
        "services/web/**",
        "shared/**",
    },
}
```

### Change Detection with dorny/paths-filter

Use the popular paths-filter action to detect which directories changed:

```go
var PathsFilterStep = workflow.Step{
    ID:   "changes",
    Name: "Detect changed paths",
    Uses: "dorny/paths-filter@v3",
    With: map[string]any{
        "filters": `
api:
  - 'services/api/**'
web:
  - 'services/web/**'
shared:
  - 'shared/**'
`,
    },
}
```

### Job Outputs for Change Detection

Export change detection results for use by other jobs:

```go
var DetectChanges = workflow.Job{
    Name:   "Detect Changes",
    RunsOn: "ubuntu-latest",
    Outputs: map[string]any{
        "api":    workflow.Steps.Get("changes", "api").String(),
        "web":    workflow.Steps.Get("changes", "web").String(),
        "shared": workflow.Steps.Get("changes", "shared").String(),
    },
    Steps: DetectChangesSteps,
}
```

### Conditional Job Execution

Run jobs only when specific services changed (including shared dependencies):

```go
var APICondition = workflow.Needs.Get("detect-changes", "api").
    Or(workflow.Needs.Get("detect-changes", "shared"))

var BuildAPI = workflow.Job{
    Name:   "Build API",
    RunsOn: "ubuntu-latest",
    Needs:  []any{DetectChanges},
    If:     APICondition.String(),
    Steps:  APIBuildSteps,
}
```

### Service-Specific Build Steps

Each service has its own build configuration with appropriate tooling:

```go
// Go API service
var APIBuildSteps = []any{
    checkout.Checkout{},
    setup_go.SetupGo{GoVersion: "1.24"},
    workflow.Step{
        Name:             "Build API",
        Run:              "go build ./...",
        WorkingDirectory: "services/api",
    },
}

// Node.js Web service
var WebBuildSteps = []any{
    checkout.Checkout{},
    setup_node.SetupNode{NodeVersion: "20", Cache: "npm"},
    workflow.Step{
        Name:             "Install dependencies",
        Run:              "npm ci",
        WorkingDirectory: "services/web",
    },
}
```

## Expression Helpers Used

- `workflow.Steps.Get(stepID, output)` - Reference step outputs
- `workflow.Needs.Get(jobID, output)` - Reference job outputs
- `Expression.Or()` - Combine conditions with OR logic
- `Expression.String()` - Convert expression for use in If conditions

## Benefits of Monorepo CI

1. **Faster builds** - Only build what changed
2. **Reduced costs** - Fewer compute minutes wasted
3. **Clear dependencies** - Shared library changes trigger dependent services
4. **Parallel execution** - All service builds run concurrently
