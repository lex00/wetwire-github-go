package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var AnalyzeSteps = []any{
	workflow.Step{
		Name: "Checkout repository",
		Uses: "actions/checkout@8e8c483db84b4bee98b60c0593521ed34d9990e8",
		With: map[string]any{
			"persist-credentials": false,
		},
	},
	workflow.Step{
		Name: "Initialize CodeQL",
		Uses: "github/codeql-action/init@5d4e8d1aca955e8d8589aabd499c5cae939e33c7",
		With: map[string]any{
			"languages": "${{ matrix.language }}",
		},
	},
	workflow.Step{
		Name: "Autobuild",
		Uses: "github/codeql-action/autobuild@5d4e8d1aca955e8d8589aabd499c5cae939e33c7",
	},
	workflow.Step{
		Name: "Perform CodeQL Analysis",
		Uses: "github/codeql-action/analyze@5d4e8d1aca955e8d8589aabd499c5cae939e33c7",
	},
}
