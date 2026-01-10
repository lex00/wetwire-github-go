package workflows

import (
	"github.com/lex00/wetwire-github-go/actions/checkout"
	"github.com/lex00/wetwire-github-go/workflow"
)

// IntegrationTestSteps are steps that connect to PostgreSQL and Redis services.
var IntegrationTestSteps = []any{
	checkout.Checkout{},
	workflow.Step{
		Name: "Wait for PostgreSQL",
		Run: `until pg_isready -h postgres -p 5432 -U testuser; do
  echo "Waiting for PostgreSQL..."
  sleep 2
done`,
	},
	workflow.Step{
		Name: "Wait for Redis",
		Run: `until redis-cli -h redis -p 6379 ping; do
  echo "Waiting for Redis..."
  sleep 2
done`,
	},
	workflow.Step{
		Name: "Run database migrations",
		Run:  "go run ./cmd/migrate up",
		Env: workflow.Env{
			"DATABASE_URL": "postgres://testuser:testpass@postgres:5432/testdb?sslmode=disable",
		},
	},
	workflow.Step{
		Name: "Run integration tests",
		Run:  "go test -v -tags=integration ./...",
		Env: workflow.Env{
			"DATABASE_URL": "postgres://testuser:testpass@postgres:5432/testdb?sslmode=disable",
			"REDIS_URL":    "redis://redis:6379",
			"TEST_ENV":     "ci",
		},
	},
	workflow.Step{
		Name: "Generate coverage report",
		Run:  "go test -coverprofile=coverage.out -tags=integration ./... && go tool cover -html=coverage.out -o coverage.html",
		Env: workflow.Env{
			"DATABASE_URL": "postgres://testuser:testpass@postgres:5432/testdb?sslmode=disable",
			"REDIS_URL":    "redis://redis:6379",
		},
	},
}

// UnitTestSteps are steps for running unit tests (no services required).
var UnitTestSteps = []any{
	checkout.Checkout{},
	workflow.Step{
		Name: "Run unit tests",
		Run:  "go test -v -short ./...",
	},
	workflow.Step{
		Name: "Run race detector",
		Run:  "go test -race -short ./...",
	},
}
