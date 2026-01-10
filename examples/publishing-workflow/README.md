# Publishing Workflow Example

A comprehensive example demonstrating how to define publishing workflows with wetwire-github-go for Docker images, GitHub releases, and multi-platform artifact distribution.

## Features Demonstrated

### Publish Workflow (triggers on tag push v*)

- **Docker image build and push** to GitHub Container Registry (GHCR)
- **Multi-platform Docker builds** (linux/amd64, linux/arm64)
- **GitHub release creation** with auto-generated notes
- **Conditional steps** based on tag pattern (stable vs prerelease)
- **Go module tagging** compatible with Go module versioning

### Release Workflow (triggers on release published)

- **Multi-platform Go binary builds** using matrix strategy
- **Asset uploads** to existing GitHub release
- **Notification/announcement steps** (Slack webhook example)
- **Job dependencies** with `needs` clause

## Project Structure

```
publishing-workflow/
├── go.mod                    # Module with replace directive
├── README.md                 # This file
├── CLAUDE.md                 # AI assistant context
└── workflows/
    ├── workflows.go          # Workflow declarations (Publish, Release)
    ├── jobs.go               # Job definitions with matrix strategy
    ├── triggers.go           # Tag push and release triggers
    └── steps.go              # Step sequences using typed wrappers
```

## Workflows

### 1. Publish Workflow

Triggered when a version tag (e.g., `v1.0.0`, `v2.0.0-beta.1`) is pushed.

**Jobs:**
- `docker` - Builds and pushes Docker image to GHCR
- `release` - Creates GitHub release with changelog

**Conditional Logic:**
- Stable releases (`v1.0.0`) get the `latest` tag on Docker images
- Prereleases (`v1.0.0-beta.1`) are marked as prerelease on GitHub

### 2. Release Workflow

Triggered when a GitHub release is published.

**Jobs:**
- `build-artifacts` - Builds Go binaries for multiple platforms using matrix
- `notify` - Sends notifications after artifacts are ready

**Matrix Strategy:**
```go
goos:   ["linux", "darwin", "windows"]
goarch: ["amd64", "arm64"]
// Excludes windows/arm64
```

## Usage

### Generate YAML

```bash
cd examples/publishing-workflow
go mod tidy
wetwire-github build ./workflows
```

This generates:
- `.github/workflows/publish.yml`
- `.github/workflows/release.yml`

### View Generated YAML

```bash
cat .github/workflows/publish.yml
cat .github/workflows/release.yml
```

### Validate with actionlint

```bash
wetwire-github validate .github/workflows/*.yml
```

## Key Patterns

### Expression Contexts

Access GitHub Actions contexts with typed helpers:

```go
// Secrets context
workflow.Secrets.GITHUB_TOKEN()          // secrets.GITHUB_TOKEN
workflow.Secrets.Get("DEPLOY_TOKEN")     // secrets.DEPLOY_TOKEN

// GitHub context
workflow.GitHub.Actor()                  // github.actor
workflow.GitHub.RefName()                // github.ref_name
workflow.GitHub.Repository()             // github.repository
workflow.GitHub.Event("release.tag_name") // github.event.release.tag_name
```

### Tag Conditions

Use conditions to differentiate between stable and prerelease tags:

```go
// Stable release (no hyphen in tag)
If: "!contains(github.ref_name, '-')"

// Prerelease (contains hyphen like v1.0.0-beta.1)
If: "contains(github.ref_name, '-')"
```

### Docker Action Wrappers

Type-safe wrappers for Docker actions:

```go
docker_login.DockerLogin{
    Registry: "ghcr.io",
    Username: workflow.GitHub.Actor().String(),
    Password: workflow.Secrets.GITHUB_TOKEN().String(),
}

docker_build_push.DockerBuildPush{
    Context:   ".",
    Push:      true,
    Platforms: "linux/amd64,linux/arm64",
    Tags:      "ghcr.io/${{ github.repository }}:${{ github.ref_name }}",
}
```

### Release Action Wrapper

Create releases with typed inputs:

```go
gh_release.GHRelease{
    GenerateReleaseNotes: true,
    BodyPath:             "CHANGELOG.md",
    Prerelease:           false,
}
```

### Matrix Strategy

Define multi-platform builds:

```go
var BuildMatrix = workflow.Matrix{
    Values: map[string][]any{
        "goos":   {"linux", "darwin", "windows"},
        "goarch": {"amd64", "arm64"},
    },
    Exclude: []map[string]any{
        {"goos": "windows", "goarch": "arm64"},
    },
}
```

### Artifact Upload

Upload build artifacts:

```go
upload_artifact.UploadArtifact{
    Name: "app-${{ matrix.goos }}-${{ matrix.goarch }}",
    Path: "dist/",
}
```

### Job Dependencies

Chain jobs with the `Needs` field:

```go
var Notify = workflow.Job{
    Name:   "Notify",
    RunsOn: "ubuntu-latest",
    Needs:  []any{BuildArtifacts},  // Waits for BuildArtifacts job
    Steps:  NotifySteps,
}
```

## Customization

### Adding More Platforms

Extend the matrix in `jobs.go`:

```go
Values: map[string][]any{
    "goos":   {"linux", "darwin", "windows", "freebsd"},
    "goarch": {"amd64", "arm64", "386"},
},
```

### Custom Registry

Replace GHCR with Docker Hub or another registry:

```go
var DockerHubLogin = docker_login.DockerLogin{
    Username: workflow.Secrets.Get("DOCKERHUB_USERNAME").String(),
    Password: workflow.Secrets.Get("DOCKERHUB_TOKEN").String(),
}
```

### Additional Notifications

Add more notification channels in `steps.go`:

```go
var DiscordNotification = workflow.Step{
    Name: "Send Discord Notification",
    If:   workflow.Secrets.Get("DISCORD_WEBHOOK").String() + " != ''",
    Run:  `curl -X POST -H "Content-Type: application/json" ...`,
}
```

## Local Development

When developing wetwire-github-go locally, add a replace directive to go.mod:

```go
replace github.com/lex00/wetwire-github-go => ../..
```

Then run `go mod tidy` before building.
