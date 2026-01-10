package workflows

import (
	"github.com/lex00/wetwire-github-go/actions/cache"
	"github.com/lex00/wetwire-github-go/actions/checkout"
	"github.com/lex00/wetwire-github-go/actions/setup_go"
	"github.com/lex00/wetwire-github-go/workflow"
)

// SetupGoMatrix sets up Go with the matrix version.
var SetupGoMatrix = setup_go.SetupGo{
	GoVersion: "${{ matrix.go }}",
}

// GoModCache caches Go module dependencies.
var GoModCache = cache.Cache{
	Path:        "~/go/pkg/mod",
	Key:         "go-mod-${{ runner.os }}-${{ matrix.go }}-${{ hashFiles('**/go.sum') }}",
	RestoreKeys: "go-mod-${{ runner.os }}-${{ matrix.go }}-",
}

// GoBuildCache caches Go build artifacts.
var GoBuildCache = cache.Cache{
	Path:        "~/.cache/go-build",
	Key:         "go-build-${{ runner.os }}-${{ matrix.go }}-${{ hashFiles('**/*.go') }}",
	RestoreKeys: "go-build-${{ runner.os }}-${{ matrix.go }}-",
}

// BuildStep compiles the project.
var BuildStep = workflow.Step{
	Name: "Build",
	Run:  "go build ./...",
}

// TestStep runs all tests with verbose output.
var TestStep = workflow.Step{
	Name: "Test",
	Run:  "go test -v -race ./...",
}

// TestSteps are the steps for the matrix test job.
var TestSteps = []any{
	checkout.Checkout{},
	SetupGoMatrix,
	GoModCache,
	GoBuildCache,
	BuildStep,
	TestStep,
}
