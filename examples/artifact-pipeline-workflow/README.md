# Artifact Pipeline Workflow Example

A complete example demonstrating how to define multi-stage GitHub Actions pipelines with artifact passing using wetwire-github-go.

## Features Demonstrated

- **Multi-stage pipeline** with Build, Test, and Release jobs
- **Artifact passing** between jobs using `upload_artifact` and `download_artifact` wrappers
- **Job dependencies** using `Needs` to enforce execution order
- **Conditional release** that only runs on version tags (v*)
- **Cross-platform builds** compiling for Linux, macOS, and Windows

## Project Structure

```
artifact-pipeline-workflow/
├── go.mod                    # Module with replace directive
├── README.md                 # This file
├── CLAUDE.md                 # AI assistant context
└── workflows/
    ├── workflows.go          # Main workflow declaration
    ├── jobs.go               # Build, Test, Release job definitions
    ├── triggers.go           # Push trigger for main and tags
    └── steps.go              # Step sequences with artifact actions
```

## Pipeline Flow

```
┌─────────┐     ┌─────────┐     ┌─────────────┐
│  Build  │────>│  Test   │────>│   Release   │
│         │     │         │     │ (tags only) │
└─────────┘     └─────────┘     └─────────────┘
     │               │                 │
     │  Upload       │  Download       │  Download
     │  binaries     │  binaries       │  binaries
     └───────────────┴─────────────────┘
```

1. **Build Job**: Compiles binaries for multiple platforms, uploads as artifact
2. **Test Job**: Downloads binaries, verifies them, runs test suite
3. **Release Job**: Downloads binaries, creates GitHub release (only on v* tags)

## Usage

### Generate YAML

```bash
cd examples/artifact-pipeline-workflow
go mod tidy
wetwire-github build .
```

This generates `.github/workflows/artifact-pipeline.yml`.

### View Generated YAML

```bash
cat .github/workflows/artifact-pipeline.yml
```

### Validate with actionlint

```bash
wetwire-github validate .github/workflows/artifact-pipeline.yml
```

## Key Patterns

### Artifact Upload

Use the typed `upload_artifact.UploadArtifact` wrapper to upload build outputs:

```go
import "github.com/lex00/wetwire-github-go/actions/upload_artifact"

var UploadBinaries = upload_artifact.UploadArtifact{
    Name:          "binaries",
    Path:          "dist/*",
    RetentionDays: 7,
}
```

### Artifact Download

Use the typed `download_artifact.DownloadArtifact` wrapper to retrieve artifacts:

```go
import "github.com/lex00/wetwire-github-go/actions/download_artifact"

var DownloadBinaries = download_artifact.DownloadArtifact{
    Name: "binaries",
    Path: "dist",
}
```

### Job Dependencies

Use `Needs` to establish job execution order:

```go
var Test = workflow.Job{
    Name:   "Test",
    RunsOn: "ubuntu-latest",
    Needs:  []any{Build},  // Wait for Build to complete
    Steps:  TestSteps,
}

var Release = workflow.Job{
    Name:   "Release",
    RunsOn: "ubuntu-latest",
    Needs:  []any{Build, Test},  // Wait for both
    Steps:  ReleaseSteps,
}
```

### Conditional Release on Tags

Use expression helpers for tag-based conditions:

```go
import "github.com/lex00/wetwire-github-go/workflow"

var ReleaseCondition = workflow.StartsWith(
    workflow.GitHub.Ref(),
    workflow.Expression("'refs/tags/v'"),
)

var Release = workflow.Job{
    Name:   "Release",
    If:     ReleaseCondition,
    // ...
}
```

### Push Triggers for Branches and Tags

Configure triggers for both branches and tags:

```go
var PipelinePush = workflow.PushTrigger{
    Branches: []string{"main"},
    Tags:     []string{"v*"},
}
```

## Triggering the Pipeline

### On Push to Main

When code is pushed to the main branch:
- Build job runs and uploads binaries
- Test job downloads binaries and runs tests
- Release job is **skipped** (no tag)

### On Version Tag

When a version tag is pushed (e.g., `git tag v1.0.0 && git push --tags`):
- Build job runs and uploads binaries
- Test job downloads binaries and runs tests
- Release job runs and creates a GitHub release with binaries

## Artifact Retention

By default, artifacts are retained for the repository's default period (usually 90 days). This example uses `RetentionDays: 7` to minimize storage usage for temporary build artifacts.

## Release Permissions

The Release job requires write permissions to create GitHub releases:

```go
var ReleasePermissions = workflow.Permissions{
    Contents: "write",
}
```
