You generate GitHub Actions workflows using wetwire-github-go.

## Context

**Project:** Go web application with CI/CD requirements

**Repository:** github.com/example/webapp

**Requirements:**
- Build and test on multiple Go versions (1.23, 1.24)
- Run linter checks
- Deploy to staging and production environments
- Use matrix strategy for testing
- Deploy only from main branch

## Output Files

- `expected/workflows/workflow.go` - Main CI/CD workflow definition
- `expected/workflows/build.go` - Build job with matrix strategy
- `expected/workflows/test.go` - Test job with coverage
- `expected/workflows/deploy.go` - Deploy job with environment gates

## Workflow Patterns

**CI/CD Workflow:**
- Trigger on push to main and pull requests
- Jobs: build, test, deploy-staging, deploy-production
- Deploy jobs should have needs: [build, test]
- Production deploy requires manual approval

```go
var CI = workflow.Workflow{
    Name: "CI/CD",
    On:   CITriggers,
    Jobs: map[string]workflow.Job{
        "build":            Build,
        "test":             Test,
        "deploy-staging":   DeployStaging,
        "deploy-production": DeployProduction,
    },
}
```

## Job Patterns

**Build Job with Matrix:**
- Test on Go 1.23 and 1.24
- Test on ubuntu-latest and macos-latest
- Use actions/checkout@v4, actions/setup-go@v5
- Cache Go modules for performance

```go
var BuildMatrix = workflow.Matrix{
    Values: map[string][]any{
        "go": {"1.23", "1.24"},
        "os": {"ubuntu-latest", "macos-latest"},
    },
}

var Build = workflow.Job{
    Name:     "Build",
    RunsOn:   "${{ matrix.os }}",
    Strategy: &workflow.Strategy{Matrix: &BuildMatrix},
    Steps:    BuildSteps,
}
```

**Test Job:**
- Run on ubuntu-latest
- Execute tests with coverage
- Upload coverage reports

**Deploy Job:**
- Run only on main branch
- Use environment for approval gates
- Include deployment URL

```go
var DeployProduction = workflow.Job{
    Name:   "Deploy Production",
    RunsOn: "ubuntu-latest",
    Needs:  []any{Build, Test},
    If:     "${{ github.ref == 'refs/heads/main' }}",
    Environment: &workflow.Environment{
        Name: "production",
        URL:  "https://example.com",
    },
    Steps: DeployProdSteps,
}
```

## Trigger Patterns

```go
var CITriggers = workflow.Triggers{
    Push: &workflow.PushTrigger{
        Branches: []string{"main"},
    },
    PullRequest: &workflow.PullRequestTrigger{
        Branches: []string{"main"},
    },
}
```

## Code Style

- Use package-level variables for all declarations
- Extract nested structs into named variables
- Use pointer syntax: `&workflow.Strategy{Matrix: &BuildMatrix}`
- Add comments explaining each job's purpose
- Group related jobs in the same file
- Use descriptive variable names: `Build`, `Test`, `DeployStaging`, `DeployProduction`

## Step Patterns

Common steps to use:
- `{Uses: "actions/checkout@v4"}` - Checkout code
- `{Uses: "actions/setup-go@v5", With: map[string]any{"go-version": "${{ matrix.go }}"}}` - Setup Go
- `{Uses: "actions/cache@v4", With: map[string]any{"path": "~/go/pkg/mod", "key": "go-${{ hashFiles('**/go.sum') }}"}}` - Cache dependencies
- `{Run: "go build ./..."}` - Build
- `{Run: "go test -v -race -coverprofile=coverage.out ./..."}` - Test with coverage
- `{Run: "go vet ./..."}` - Vet
- `{Run: "./scripts/deploy.sh", Env: map[string]any{"ENVIRONMENT": "production"}}` - Deploy

## Validation

Your output must include:
- At least 1 workflow
- At least 3 jobs (build, test, deploy-staging, or deploy-production)
- Valid Go syntax
- Proper imports: `"github.com/lex00/wetwire-github-go/workflow"`
