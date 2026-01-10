# artifact-pipeline-workflow

Example multi-stage artifact pipeline using wetwire-github-go.

## What This Is

A reference implementation showing how to pass artifacts between jobs in a build-test-release pipeline using typed Go structs.

## Key Files

- `workflows/workflows.go` - Main workflow declaration
- `workflows/jobs.go` - Build, Test, Release job definitions with Needs dependencies
- `workflows/triggers.go` - Push triggers for main branch and version tags
- `workflows/steps.go` - Step sequences with upload/download artifact actions

## Patterns Used

1. **Flat variables** - All structs are package-level variables, not nested
2. **Typed wrappers** - Uses `upload_artifact.UploadArtifact{}`, `download_artifact.DownloadArtifact{}`
3. **Job dependencies** - Uses `Needs` to enforce Build -> Test -> Release order
4. **Conditional execution** - Uses `workflow.StartsWith()` for tag-based release condition
5. **Permissions** - Uses explicit `Permissions` for release job

## Build Command

```bash
wetwire-github build ./workflows
```

Output: `.github/workflows/artifact-pipeline.yml`
