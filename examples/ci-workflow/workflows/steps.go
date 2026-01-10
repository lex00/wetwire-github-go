package workflows

import (
	"github.com/lex00/wetwire-github-go/actions/cache"
	"github.com/lex00/wetwire-github-go/actions/checkout"
	"github.com/lex00/wetwire-github-go/actions/setup_go"
	"github.com/lex00/wetwire-github-go/workflow"
)

// BuildSteps are the steps for the build job.
var BuildSteps = []any{
	checkout.Checkout{},
	setup_go.SetupGo{
		GoVersion: "${{ matrix.go }}",
	},
	cache.Cache{
		Path:        "~/go/pkg/mod",
		Key:         "go-mod-${{ runner.os }}-${{ hashFiles('**/go.sum') }}",
		RestoreKeys: "go-mod-${{ runner.os }}-",
	},
	workflow.Step{
		Name: "Build",
		Run:  "go build ./...",
	},
	workflow.Step{
		Name: "Test",
		Run:  "go test -v ./...",
	},
}

// LintSteps are the steps for the lint job.
var LintSteps = []any{
	checkout.Checkout{},
	setup_go.SetupGo{
		GoVersion: "1.24",
	},
	workflow.Step{
		Name: "Install golangci-lint",
		Run:  "go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest",
	},
	workflow.Step{
		Name: "Run linter",
		Run:  "golangci-lint run ./...",
	},
}
