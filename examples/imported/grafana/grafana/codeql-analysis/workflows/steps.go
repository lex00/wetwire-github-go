package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var AnalyzeSteps = []any{
	workflow.Step{
		Name: "Checkout repository",
		Uses: "actions/checkout@v5",
		If:   "needs.detect-changes.outputs[matrix.language] == 'true'",
		With: map[string]any{
			"fetch-depth":         2,
			"persist-credentials": false,
		},
	},
	workflow.Step{
		Name: "Set go version",
		Uses: "actions/setup-go@44694675825211faa026b3c33043df3e48a5fa00",
		If:   "matrix.language == 'go' && needs.detect-changes.outputs.go == 'true'",
		With: map[string]any{
			"cache":           false,
			"go-version-file": "go.mod",
		},
	},
	workflow.Step{
		Name: "Initialize CodeQL",
		Uses: "github/codeql-action/init@v4",
		If:   "needs.detect-changes.outputs[matrix.language] == 'true'",
		With: map[string]any{
			"languages": "${{ matrix.language }}",
		},
	},
	workflow.Step{
		Name: "Build go files",
		Run: `go mod verify
make build-go
`,
		If: "matrix.language == 'go' && needs.detect-changes.outputs.go == 'true'",
	},
	workflow.Step{
		Name: "Perform CodeQL Analysis",
		Uses: "github/codeql-action/analyze@v4",
	},
}

var DetectChangesSteps = []any{
	workflow.Step{
		Uses: "actions/checkout@v5",
		With: map[string]any{
			"fetch-depth":         2,
			"persist-credentials": true,
		},
	},
	workflow.Step{
		ID:   "detect-changes",
		Name: "Detect changes",
		Uses: "./.github/actions/change-detection",
		With: map[string]any{
			"self": ".github/workflows/codeql-analysis.yml",
		},
	},
}
