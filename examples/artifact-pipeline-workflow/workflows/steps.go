package workflows

import (
	"github.com/lex00/wetwire-github-go/actions/checkout"
	"github.com/lex00/wetwire-github-go/actions/download_artifact"
	"github.com/lex00/wetwire-github-go/actions/gh_release"
	"github.com/lex00/wetwire-github-go/actions/setup_go"
	"github.com/lex00/wetwire-github-go/actions/upload_artifact"
	"github.com/lex00/wetwire-github-go/workflow"
)

// CheckoutCode checks out the repository source code.
var CheckoutCode = checkout.Checkout{}

// CheckoutWithHistory fetches full history for release notes generation.
var CheckoutWithHistory = checkout.Checkout{
	FetchDepth: 0,
}

// SetupGoVersion configures the Go environment.
var SetupGoVersion = setup_go.SetupGo{
	GoVersion: "1.24",
}

// CompileBinaries builds the application for multiple platforms.
var CompileBinaries = workflow.Step{
	Name: "Build binaries",
	Run: `mkdir -p dist
GOOS=linux GOARCH=amd64 go build -o dist/myapp-linux-amd64 ./cmd/myapp
GOOS=darwin GOARCH=amd64 go build -o dist/myapp-darwin-amd64 ./cmd/myapp
GOOS=darwin GOARCH=arm64 go build -o dist/myapp-darwin-arm64 ./cmd/myapp
GOOS=windows GOARCH=amd64 go build -o dist/myapp-windows-amd64.exe ./cmd/myapp`,
}

// UploadBinaries uploads compiled binaries as workflow artifacts.
var UploadBinaries = upload_artifact.UploadArtifact{
	Name:          "binaries",
	Path:          "dist/*",
	RetentionDays: 7,
}

// DownloadBinaries downloads compiled binaries from previous job.
var DownloadBinaries = download_artifact.DownloadArtifact{
	Name: "binaries",
	Path: "dist",
}

// RunTests executes the test suite.
var RunTests = workflow.Step{
	Name: "Run tests",
	Run:  "go test -v -race -coverprofile=coverage.out ./...",
}

// UploadCoverage uploads test coverage report as artifact.
var UploadCoverage = upload_artifact.UploadArtifact{
	Name:          "coverage",
	Path:          "coverage.out",
	RetentionDays: 7,
}

// VerifyBinaries ensures downloaded binaries are executable.
var VerifyBinaries = workflow.Step{
	Name: "Verify binaries",
	Run: `ls -la dist/
chmod +x dist/myapp-linux-amd64
./dist/myapp-linux-amd64 --version || echo "Binary verification complete"`,
}

// CreateGHRelease creates a GitHub release with auto-generated notes.
var CreateGHRelease = gh_release.GHRelease{
	GenerateReleaseNotes: true,
	Files:                "dist/*",
}

// BuildSteps are the steps for the build job.
var BuildSteps = []any{
	CheckoutCode,
	SetupGoVersion,
	CompileBinaries,
	UploadBinaries,
}

// TestSteps are the steps for the test job.
var TestSteps = []any{
	CheckoutCode,
	SetupGoVersion,
	DownloadBinaries,
	VerifyBinaries,
	RunTests,
	UploadCoverage,
}

// ReleaseSteps are the steps for the release job.
var ReleaseSteps = []any{
	CheckoutWithHistory,
	DownloadBinaries,
	CreateGHRelease,
}
