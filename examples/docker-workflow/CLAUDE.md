# docker-workflow

Example Docker build and push workflow using wetwire-github-go.

## What This Is

A reference implementation showing how to declare Docker build workflows using typed Go structs. This example creates a workflow that builds and pushes Docker images to GitHub Container Registry (GHCR).

## Key Files

- `workflows/workflows.go` - Main Docker workflow declaration
- `workflows/jobs.go` - Build job definition with conditional push
- `workflows/triggers.go` - Push and pull request trigger configurations
- `workflows/steps.go` - Step sequences using typed Docker action wrappers

## Patterns Used

1. **Flat variables** - All structs are package-level variables, not nested
2. **Typed wrappers** - Uses `docker_login.DockerLogin{}`, `docker_setup_buildx.DockerSetupBuildx{}`, `docker_build_push.DockerBuildPush{}`
3. **Conditional push** - Only pushes to GHCR on main branch
4. **GitHub expressions** - Uses `${{ github.repository }}` and `${{ secrets.GITHUB_TOKEN }}`

## Build Command

```bash
wetwire-github build ./workflows
```

Output: `.github/workflows/docker.yml`
