package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var ConsistencyChecksSteps = []any{
	workflow.Step{
		Name: "Fetch source code",
		Uses: "actions/checkout@8e8c483db84b4bee98b60c0593521ed34d9990e8",
		With: map[string]any{
			"fetch-depth": 0,
		},
	},
	workflow.Step{
		ID:   "go",
		Name: "Determine Go version",
		Uses: "./.github/actions/go-version",
	},
	workflow.Step{
		Name: "Install Go toolchain",
		Uses: "actions/setup-go@4dc6199c7b1a012772edbd06daecab0f50c9053c",
		With: map[string]any{
			"cache-dependency-path": "go.sum",
			"go-version":            "${{ steps.go.outputs.version }}",
		},
	},
	workflow.Step{
		Name: "go.mod and go.sum consistency check",
		Run: `make syncdeps
CHANGED="$(git status --porcelain)"
if [[ -n "$CHANGED" ]]; then
  git diff
  echo >&2 "ERROR: go.mod/go.sum files are not up-to-date. Run 'make syncdeps' and then commit the updated files."
  echo >&2 $'Affected files:\n'"$CHANGED"
  exit 1
fi
`,
	},
	workflow.Step{
		Name: "Cache protobuf tools",
		Uses: "actions/cache@9255dc7a253b0ccc959486e2bca901246202afeb",
		With: map[string]any{
			"key":          "protobuf-tools-${{ hashFiles('tools/protobuf-compile/protobuf-compile.go') }}",
			"path":         "tools/protobuf-compile/.workdir",
			"restore-keys": "protobuf-tools-\n",
		},
	},
	workflow.Step{
		Name: "Code consistency checks",
		Run: `make fmtcheck importscheck vetcheck copyright generate staticcheck exhaustive protobuf
CHANGED="$(git status --porcelain)"
if [[ -n "$CHANGED" ]]; then
  git diff
  echo >&2 "ERROR: Generated files are inconsistent. Run 'make generate' and 'make protobuf' locally and then commit the updated files."
  echo >&2 $'Affected files:\n'"$CHANGED"
  exit 1
fi
`,
	},
}

var E2eTestsSteps = []any{
	workflow.Step{
		Name: "Fetch source code",
		Uses: "actions/checkout@8e8c483db84b4bee98b60c0593521ed34d9990e8",
	},
	workflow.Step{
		ID:   "go",
		Name: "Determine Go version",
		Uses: "./.github/actions/go-version",
	},
	workflow.Step{
		Name: "Install Go toolchain",
		Uses: "actions/setup-go@4dc6199c7b1a012772edbd06daecab0f50c9053c",
		With: map[string]any{
			"cache-dependency-path": "go.sum",
			"go-version":            "${{ steps.go.outputs.version }}",
		},
	},
	workflow.Step{
		Name: "End-to-end tests",
		Run: `TF_ACC=1 go test -v ./internal/command/e2etest
`,
	},
}

var RaceTestsSteps = []any{
	workflow.Step{
		Name: "Fetch source code",
		Uses: "actions/checkout@8e8c483db84b4bee98b60c0593521ed34d9990e8",
	},
	workflow.Step{
		ID:   "go",
		Name: "Determine Go version",
		Uses: "./.github/actions/go-version",
	},
	workflow.Step{
		Name: "Install Go toolchain",
		Uses: "actions/setup-go@4dc6199c7b1a012772edbd06daecab0f50c9053c",
		With: map[string]any{
			"cache-dependency-path": "go.sum",
			"go-version":            "${{ steps.go.outputs.version }}",
		},
	},
	workflow.Step{
		Name: "Race detector",
		Run: `go test -race ./internal/terraform ./internal/command ./internal/states
`,
	},
}

var UnitTestsSteps = []any{
	workflow.Step{
		Name: "Fetch source code",
		Uses: "actions/checkout@8e8c483db84b4bee98b60c0593521ed34d9990e8",
	},
	workflow.Step{
		ID:   "go",
		Name: "Determine Go version",
		Uses: "./.github/actions/go-version",
	},
	workflow.Step{
		Name: "Install Go toolchain",
		Uses: "actions/setup-go@4dc6199c7b1a012772edbd06daecab0f50c9053c",
		With: map[string]any{
			"cache-dependency-path": "go.sum",
			"go-version":            "${{ steps.go.outputs.version }}",
		},
	},
	workflow.Step{
		Name: "Unit tests",
		Run: `# We run tests for all packages from all modules in this repository.
for dir in $(go list -m -f '{{.Dir}}' github.com/hashicorp/terraform/...); do
    (cd $dir && go test -cover "./...")
done
`,
	},
}
