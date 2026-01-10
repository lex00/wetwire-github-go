package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// JobContainer runs job steps in a golang container.
var JobContainer = workflow.Container{
	Image: "golang:1.24",
	Env: workflow.Env{
		"CGO_ENABLED": "0",
		"GOOS":        "linux",
	},
}

// PostgresService provides PostgreSQL for integration tests.
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

// RedisService provides Redis for caching tests.
var RedisService = workflow.Service{
	Image: "redis:7-alpine",
	Ports: []any{"6379:6379"},
	Options: "--health-cmd \"redis-cli ping\" " +
		"--health-interval 10s " +
		"--health-timeout 5s " +
		"--health-retries 5",
}

// IntegrationTest runs integration tests with service containers.
var IntegrationTest = workflow.Job{
	Name:      "Integration Tests",
	RunsOn:    "ubuntu-latest",
	Container: &JobContainer,
	Services: map[string]workflow.Service{
		"postgres": PostgresService,
		"redis":    RedisService,
	},
	Steps: IntegrationTestSteps,
}

// UnitTest runs unit tests without containers (faster feedback).
var UnitTest = workflow.Job{
	Name:   "Unit Tests",
	RunsOn: "ubuntu-latest",
	Container: &workflow.Container{
		Image: "golang:1.24",
	},
	Steps: UnitTestSteps,
}
