# Container Services Workflow Example

A complete example demonstrating how to run GitHub Actions jobs in Docker containers with service containers for database and cache testing.

## Features Demonstrated

- **Job Container** - Running all job steps inside a Docker container (`golang:1.24`)
- **Service Containers** - PostgreSQL and Redis service containers for integration testing
- **Health Checks** - Service container health check configuration
- **Environment Variables** - Database connection strings and service configuration
- **Integration Testing Patterns** - Waiting for services, running migrations, executing tests

## Project Structure

```
container-services-workflow/
├── go.mod                    # Module with replace directive
├── README.md                 # This file
├── CLAUDE.md                 # AI assistant context
└── workflows/
    ├── workflows.go          # Main workflow declaration
    ├── jobs.go               # Job definitions with container and services
    ├── triggers.go           # Push and pull request triggers
    └── steps.go              # Step sequences for integration tests
```

## Usage

### Generate YAML

```bash
cd examples/container-services-workflow
go mod tidy
wetwire-github build .
```

This generates `.github/workflows/container-services.yml`.

### View Generated YAML

```bash
cat .github/workflows/container-services.yml
```

### Validate with actionlint

```bash
wetwire-github validate .github/workflows/container-services.yml
```

## Key Patterns

### Job Container

Run all job steps in a container image:

```go
var JobContainer = workflow.Container{
    Image: "golang:1.24",
    Env: workflow.Env{
        "CGO_ENABLED": "0",
        "GOOS":        "linux",
    },
}

var IntegrationTest = workflow.Job{
    Container: &JobContainer,
    // ...
}
```

### Service Containers

Define service containers with health checks:

```go
var PostgresService = workflow.Service{
    Image: "postgres:16",
    Env: workflow.Env{
        "POSTGRES_USER":     "testuser",
        "POSTGRES_PASSWORD": "testpass",
        "POSTGRES_DB":       "testdb",
    },
    Ports: []any{"5432:5432"},
    Options: "--health-cmd pg_isready " +
        "--health-interval 10s " +
        "--health-timeout 5s " +
        "--health-retries 5",
}
```

### Service Map in Jobs

Attach services to jobs using a map:

```go
var IntegrationTest = workflow.Job{
    Services: map[string]workflow.Service{
        "postgres": PostgresService,
        "redis":    RedisService,
    },
    // ...
}
```

### Connection Environment Variables

Pass service connection strings to test steps:

```go
workflow.Step{
    Name: "Run integration tests",
    Run:  "go test -v -tags=integration ./...",
    Env: workflow.Env{
        "DATABASE_URL": "postgres://testuser:testpass@postgres:5432/testdb?sslmode=disable",
        "REDIS_URL":    "redis://redis:6379",
    },
}
```

## Service Container Details

### PostgreSQL Service

- **Image**: `postgres:16`
- **Port**: 5432
- **Hostname**: `postgres` (service name)
- **Health Check**: `pg_isready` command
- **Default Database**: `testdb`

### Redis Service

- **Image**: `redis:7-alpine`
- **Port**: 6379
- **Hostname**: `redis` (service name)
- **Health Check**: `redis-cli ping` command

## Container vs Non-Container Jobs

This example includes two job types:

1. **Integration Tests** - Runs in a container with service containers attached
2. **Unit Tests** - Runs in a container without services (faster feedback)

Unit tests can run in parallel with integration tests since they do not depend on external services.

## Health Check Options

The `Options` field supports Docker health check flags:

| Flag | Description |
|------|-------------|
| `--health-cmd` | Command to check health |
| `--health-interval` | Time between checks |
| `--health-timeout` | Timeout for each check |
| `--health-retries` | Number of retries before unhealthy |

## Notes

- Service containers are only available when the job runs on a Linux runner
- The job container and service containers share a Docker network
- Service hostnames match the keys in the `Services` map
- Ports are mapped as `host:container` format
