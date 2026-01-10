package workflows

import (
	"github.com/lex00/wetwire-github-go/actions/cache"
	"github.com/lex00/wetwire-github-go/actions/checkout"
	"github.com/lex00/wetwire-github-go/actions/setup_go"
	"github.com/lex00/wetwire-github-go/workflow"
)

// BuildSteps are the steps for the build job.
var BuildSteps = []workflow.Step{
	checkout.Checkout{}.ToStep(),
	setup_go.SetupGo{
		GoVersion: "${{ matrix.go }}",
	}.ToStep(),
	cache.Cache{
		Path:        "~/go/pkg/mod",
		Key:         "go-mod-${{ runner.os }}-${{ hashFiles('**/go.sum') }}",
		RestoreKeys: "go-mod-${{ runner.os }}-",
	}.ToStep(),
	{
		Name: "Build",
		Run:  "go build ./...",
	},
	{
		Name: "Test",
		Run:  "go test -v ./...",
	},
}

// LintSteps are the steps for the lint job.
var LintSteps = []workflow.Step{
	checkout.Checkout{}.ToStep(),
	setup_go.SetupGo{
		GoVersion: "1.24",
	}.ToStep(),
	{
		Name: "Install golangci-lint",
		Run:  "go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest",
	},
	{
		Name: "Run linter",
		Run:  "golangci-lint run ./...",
	},
}
