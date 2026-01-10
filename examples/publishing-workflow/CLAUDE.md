# publishing-workflow

Example publishing workflow using wetwire-github-go.

## What This Is

A reference implementation showing how to declare publishing workflows using typed Go structs. This example creates workflows for Docker image publishing, Go module tagging, GitHub releases, and multi-platform artifact builds.

## Key Files

- `workflows/workflows.go` - Main workflow declarations (Publish and Release)
- `workflows/jobs.go` - Job definitions with permissions
- `workflows/triggers.go` - Tag push and release triggers
- `workflows/steps.go` - Step sequences using typed action wrappers

## Patterns Used

1. **Flat variables** - All structs are package-level variables, not nested
2. **Typed wrappers** - Uses `checkout.Checkout{}`, `setup_go.SetupGo{}`, `docker_login.DockerLogin{}`, `docker_build_push.DockerBuildPush{}`, `gh_release.GHRelease{}`, `upload_artifact.UploadArtifact{}`
3. **Expression contexts** - Uses `workflow.Secrets`, `workflow.GitHub`, `workflow.TagPrefix()`
4. **Conditional steps** - Tag-based conditions using `workflow.TagPrefix()`

## Build Command

```bash
wetwire-github build ./workflows
```

Output: `.github/workflows/publish.yml` and `.github/workflows/release.yml`
