package workflows

import (
	"github.com/lex00/wetwire-github-go/actions/attest_build_provenance"
	"github.com/lex00/wetwire-github-go/actions/checkout"
	"github.com/lex00/wetwire-github-go/actions/codeql_analyze"
	"github.com/lex00/wetwire-github-go/actions/codeql_init"
	"github.com/lex00/wetwire-github-go/actions/setup_go"
	"github.com/lex00/wetwire-github-go/actions/trivy"
	"github.com/lex00/wetwire-github-go/workflow"
)

// CodeQLSteps performs CodeQL security analysis.
var CodeQLSteps = []any{
	checkout.Checkout{},
	codeql_init.CodeQLInit{
		Languages: "go",
		Queries:   "security-extended",
	},
	workflow.Step{
		Name: "Build",
		Run:  "go build ./...",
	},
	codeql_analyze.CodeQLAnalyze{},
}

// TrivySteps scans the repository for vulnerabilities.
var TrivySteps = []any{
	checkout.Checkout{},
	trivy.Trivy{
		ScanType: "fs",
		Format:   "sarif",
		Output:   "trivy-results.sarif",
		Severity: "CRITICAL,HIGH",
	},
	workflow.Step{
		Name: "Upload Trivy results to GitHub Security",
		Uses: "github/codeql-action/upload-sarif@v3",
		With: map[string]any{
			"sarif_file": "trivy-results.sarif",
		},
	},
}

// BuildAttestSteps builds an artifact and generates SLSA attestation.
var BuildAttestSteps = []any{
	checkout.Checkout{},
	setup_go.SetupGo{
		GoVersion: "1.24",
	},
	workflow.Step{
		Name: "Build binary",
		Run:  "go build -o myapp ./...",
	},
	attest_build_provenance.AttestBuildProvenance{
		SubjectPath: "myapp",
	},
}
