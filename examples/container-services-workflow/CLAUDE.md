# container-services-workflow

Example workflow demonstrating Docker container and service configurations.

## What This Is

A reference implementation showing how to run GitHub Actions jobs in Docker containers with service containers (PostgreSQL, Redis) for integration testing.

## Key Files

- `workflows/workflows.go` - Main workflow declaration
- `workflows/jobs.go` - Job definitions with Container and Services configuration
- `workflows/triggers.go` - Push and pull request trigger configurations
- `workflows/steps.go` - Integration test steps with database connections

## Patterns Used

1. **Job Container** - `workflow.Container{Image: "golang:1.24", ...}`
2. **Service Containers** - `workflow.Service{Image: "postgres:16", ...}`
3. **Health Checks** - Service `Options` field for Docker health check flags
4. **Environment Variables** - Connection strings in step `Env` maps
5. **Flat variables** - All structs are package-level variables

## Build Command

```bash
wetwire-github build ./workflows
```

Output: `.github/workflows/container-services.yml`
