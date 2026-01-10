package workflows

import (
	"github.com/lex00/wetwire-github-go/actions/checkout"
	"github.com/lex00/wetwire-github-go/actions/setup_go"
	"github.com/lex00/wetwire-github-go/actions/setup_node"
	"github.com/lex00/wetwire-github-go/workflow"
)

// PathsFilterStep uses dorny/paths-filter to detect changed paths.
// Outputs: api, web, shared (boolean strings "true"/"false").
var PathsFilterStep = workflow.Step{
	ID:   "changes",
	Name: "Detect changed paths",
	Uses: "dorny/paths-filter@v3",
	With: map[string]any{
		"filters": `
api:
  - 'services/api/**'
web:
  - 'services/web/**'
shared:
  - 'shared/**'
`,
	},
}

// DetectChangesSteps are the steps for the detect-changes job.
var DetectChangesSteps = []any{
	checkout.Checkout{},
	PathsFilterStep,
}

// APIBuildSteps are the steps for building the API service.
var APIBuildSteps = []any{
	checkout.Checkout{},
	setup_go.SetupGo{
		GoVersion: "1.24",
	},
	workflow.Step{
		Name:             "Build API",
		Run:              "go build ./...",
		WorkingDirectory: "services/api",
	},
	workflow.Step{
		Name:             "Test API",
		Run:              "go test -v ./...",
		WorkingDirectory: "services/api",
	},
}

// WebBuildSteps are the steps for building the Web service.
var WebBuildSteps = []any{
	checkout.Checkout{},
	setup_node.SetupNode{
		NodeVersion: "20",
		Cache:       "npm",
	},
	workflow.Step{
		Name:             "Install dependencies",
		Run:              "npm ci",
		WorkingDirectory: "services/web",
	},
	workflow.Step{
		Name:             "Build Web",
		Run:              "npm run build",
		WorkingDirectory: "services/web",
	},
	workflow.Step{
		Name:             "Test Web",
		Run:              "npm test",
		WorkingDirectory: "services/web",
	},
}

// SharedBuildSteps are the steps for building the Shared library.
var SharedBuildSteps = []any{
	checkout.Checkout{},
	setup_go.SetupGo{
		GoVersion: "1.24",
	},
	workflow.Step{
		Name:             "Build Shared",
		Run:              "go build ./...",
		WorkingDirectory: "shared",
	},
	workflow.Step{
		Name:             "Test Shared",
		Run:              "go test -v ./...",
		WorkingDirectory: "shared",
	},
}
