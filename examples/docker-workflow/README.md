# Docker Workflow Example

A complete example demonstrating how to define Docker build and push workflows with wetwire-github-go.

## Features Demonstrated

- **Docker build workflow** with GHCR (GitHub Container Registry) integration
- **Conditional push** - builds on PRs without pushing, pushes on main branch
- **Typed action wrappers** for docker/login-action, docker/setup-buildx-action, docker/build-push-action
- **GitHub secrets** integration for registry authentication

## Project Structure

```
docker-workflow/
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
cd examples/docker-workflow
go mod tidy
wetwire-github build .
```

This generates `.github/workflows/docker.yml`.

### View Generated YAML

```bash
cat .github/workflows/docker.yml
```

### Validate with actionlint

```bash
wetwire-github validate .github/workflows/docker.yml
```

### Local Development

When developing wetwire-github-go locally, add a replace directive to go.mod:

```go
replace github.com/lex00/wetwire-github-go => ../..
```

Then run `go mod tidy` before building.

## Key Patterns

### Docker Action Wrappers

Instead of raw `uses:` strings, use typed wrappers directly in `[]any{}` slices:

```go
docker_setup_buildx.DockerSetupBuildx{}
docker_login.DockerLogin{
    Registry: "ghcr.io",
    Username: "${{ github.actor }}",
    Password: "${{ secrets.GITHUB_TOKEN }}",
}
docker_build_push.DockerBuildPush{
    Context: ".",
    Push:    true,
    Tags:    "ghcr.io/${{ github.repository }}:latest",
}
```

### Conditional Push Logic

Use GitHub expressions to control push behavior:

```go
// Build job pushes only on main branch
var BuildDocker = docker_build_push.DockerBuildPush{
    Push: "${{ github.ref == 'refs/heads/main' }}" != "",  // Expression for conditional
}
```

### Flat Variable Structure

Extract all nested structs to package-level variables for clarity:

```go
// Separate variables for each component
var DockerPush = workflow.PushTrigger{...}
var DockerPullRequest = workflow.PullRequestTrigger{...}
var DockerTriggers = workflow.Triggers{Push: &DockerPush, PullRequest: &DockerPullRequest}
```
