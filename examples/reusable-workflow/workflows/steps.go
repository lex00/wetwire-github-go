package workflows

import (
	"github.com/lex00/wetwire-github-go/actions/checkout"
	"github.com/lex00/wetwire-github-go/actions/setup_go"
	"github.com/lex00/wetwire-github-go/actions/upload_artifact"
	"github.com/lex00/wetwire-github-go/workflow"
)

// BuildSteps are the steps for the build job in the reusable workflow.
var BuildSteps = []any{
	checkout.Checkout{},
	setup_go.SetupGo{
		GoVersion: "${{ inputs.go_version }}",
	},
	workflow.Step{
		Name: "Build binary",
		ID:   "build",
		Run: `VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
GOOS=${{ inputs.build_target }} go build -ldflags="-X main.version=$VERSION" -o app ./...
echo "artifact=app-${{ inputs.build_target }}-$VERSION" >> $GITHUB_OUTPUT
echo "version=$VERSION" >> $GITHUB_OUTPUT`,
	},
	workflow.Step{
		Name: "Run tests",
		If:   "${{ inputs.run_tests }}",
		Run:  "go test -v ./...",
	},
	upload_artifact.UploadArtifact{
		Name: "${{ steps.build.outputs.artifact }}",
		Path: "app",
	},
}

// UseOutputSteps demonstrate using outputs from the reusable workflow.
var UseOutputSteps = []any{
	workflow.Step{
		Name: "Display build info",
		Run: `echo "Artifact: ${{ needs.call-build.outputs.artifact_name }}"
echo "Version: ${{ needs.call-build.outputs.build_version }}"`,
	},
}
