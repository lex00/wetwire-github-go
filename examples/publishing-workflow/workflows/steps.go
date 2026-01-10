package workflows

import (
	"github.com/lex00/wetwire-github-go/actions/checkout"
	"github.com/lex00/wetwire-github-go/actions/docker_build_push"
	"github.com/lex00/wetwire-github-go/actions/docker_login"
	"github.com/lex00/wetwire-github-go/actions/docker_setup_buildx"
	"github.com/lex00/wetwire-github-go/actions/gh_release"
	"github.com/lex00/wetwire-github-go/actions/setup_go"
	"github.com/lex00/wetwire-github-go/actions/upload_artifact"
	"github.com/lex00/wetwire-github-go/workflow"
)

// ============================================================================
// Docker Publish Steps
// ============================================================================

// CheckoutForDocker checks out the repository for Docker builds.
var CheckoutForDocker = checkout.Checkout{}

// SetupBuildx configures Docker Buildx for multi-platform builds.
var SetupBuildx = docker_setup_buildx.DockerSetupBuildx{}

// GHCRLogin logs into GitHub Container Registry.
var GHCRLogin = docker_login.DockerLogin{
	Registry: "ghcr.io",
	Username: workflow.GitHub.Actor().String(),
	Password: workflow.Secrets.GITHUB_TOKEN().String(),
}

// DockerBuildAndPush builds and pushes the Docker image with version tags.
var DockerBuildAndPush = docker_build_push.DockerBuildPush{
	Context:   ".",
	Push:      true,
	Platforms: "linux/amd64,linux/arm64",
	Tags:      "ghcr.io/${{ github.repository }}:${{ github.ref_name }}\nghcr.io/${{ github.repository }}:latest",
	CacheFrom: "type=gha",
	CacheTo:   "type=gha,mode=max",
	Labels:    "org.opencontainers.image.source=${{ github.server_url }}/${{ github.repository }}",
}

// DockerBuildPrerelease builds Docker image for prereleases (no latest tag).
var DockerBuildPrerelease = docker_build_push.DockerBuildPush{
	Context:   ".",
	Push:      true,
	Platforms: "linux/amd64,linux/arm64",
	Tags:      "ghcr.io/${{ github.repository }}:${{ github.ref_name }}",
	CacheFrom: "type=gha",
	CacheTo:   "type=gha,mode=max",
}

// DockerPushStep pushes stable releases with the latest tag.
var DockerPushStep = workflow.Step{
	Name: "Build and Push Docker Image",
	If:   "!contains(github.ref_name, '-')",
	Uses: DockerBuildAndPush.Action(),
	With: DockerBuildAndPush.Inputs(),
}

// DockerPushPrereleaseStep pushes prereleases without the latest tag.
var DockerPushPrereleaseStep = workflow.Step{
	Name: "Build and Push Docker Image (Prerelease)",
	If:   "contains(github.ref_name, '-')",
	Uses: DockerBuildPrerelease.Action(),
	With: DockerBuildPrerelease.Inputs(),
}

// DockerPublishSteps are the steps for the Docker publish job.
var DockerPublishSteps = []any{
	CheckoutForDocker,
	SetupBuildx,
	GHCRLogin,
	DockerPushStep,
	DockerPushPrereleaseStep,
}

// ============================================================================
// Create Release Steps
// ============================================================================

// CheckoutWithHistory checks out with full history for changelog generation.
var CheckoutWithHistory = checkout.Checkout{
	FetchDepth: 0,
}

// GenerateChangelog generates a changelog from git commits.
var GenerateChangelog = workflow.Step{
	Name: "Generate Changelog",
	ID:   "changelog",
	Run: `echo "## Changes" > CHANGELOG.md
git log $(git describe --tags --abbrev=0 HEAD^)..HEAD --pretty=format:"- %s" >> CHANGELOG.md 2>/dev/null || echo "- Initial release" >> CHANGELOG.md
echo "" >> CHANGELOG.md`,
}

// CreateGHRelease creates a GitHub release with auto-generated notes.
var CreateGHRelease = gh_release.GHRelease{
	GenerateReleaseNotes: true,
	BodyPath:             "CHANGELOG.md",
}

// CreatePrereleaseStep creates a prerelease for tags containing a hyphen.
var CreatePrereleaseStep = workflow.Step{
	Name: "Create Prerelease",
	If:   "contains(github.ref_name, '-')",
	Uses: CreateGHRelease.Action(),
	With: map[string]any{
		"generate_release_notes": true,
		"body_path":              "CHANGELOG.md",
		"prerelease":             true,
	},
}

// CreateStableReleaseStep creates a stable release for non-prerelease tags.
var CreateStableReleaseStep = workflow.Step{
	Name: "Create Release",
	If:   "!contains(github.ref_name, '-')",
	Uses: CreateGHRelease.Action(),
	With: CreateGHRelease.Inputs(),
}

// CreateReleaseSteps are the steps for the create release job.
var CreateReleaseSteps = []any{
	CheckoutWithHistory,
	GenerateChangelog,
	CreatePrereleaseStep,
	CreateStableReleaseStep,
}

// ============================================================================
// Build Artifacts Steps (Release workflow)
// ============================================================================

// CheckoutForBuild checks out the repository for building.
var CheckoutForBuild = checkout.Checkout{}

// SetupGoForBuild sets up Go for building binaries.
var SetupGoForBuild = setup_go.SetupGo{
	GoVersion: "1.23",
}

// BuildBinary builds the Go binary for the target platform.
var BuildBinary = workflow.Step{
	Name: "Build Binary",
	ID:   "build",
	Run:  `go build -ldflags="-s -w -X main.version=${{ github.event.release.tag_name }}" -o dist/app-${{ matrix.goos }}-${{ matrix.goarch }}${{ matrix.goos == 'windows' && '.exe' || '' }} ./cmd/app`,
	Env: workflow.Env{
		"GOOS":        "${{ matrix.goos }}",
		"GOARCH":      "${{ matrix.goarch }}",
		"CGO_ENABLED": "0",
	},
}

// UploadBinaryArtifact uploads the built binary as an artifact.
var UploadBinaryArtifact = upload_artifact.UploadArtifact{
	Name: "app-${{ matrix.goos }}-${{ matrix.goarch }}",
	Path: "dist/",
}

// UploadToRelease uploads artifacts to the GitHub release.
var UploadToRelease = workflow.Step{
	Name: "Upload to Release",
	Uses: "softprops/action-gh-release@v2",
	With: map[string]any{
		"files": "dist/*",
	},
	Env: workflow.Env{
		"GITHUB_TOKEN": workflow.Secrets.GITHUB_TOKEN().String(),
	},
}

// BuildArtifactSteps are the steps for building and uploading artifacts.
var BuildArtifactSteps = []any{
	CheckoutForBuild,
	SetupGoForBuild,
	BuildBinary,
	UploadBinaryArtifact,
	UploadToRelease,
}

// ============================================================================
// Notify Steps
// ============================================================================

// PrintReleaseInfo prints release information for notification.
var PrintReleaseInfo = workflow.Step{
	Name: "Print Release Info",
	Run: `echo "Release ${{ github.event.release.tag_name }} published!"
echo "Release URL: ${{ github.event.release.html_url }}"
echo "Artifacts have been uploaded to the release."`,
}

// SlackNotification sends a notification to Slack (example using curl).
var SlackNotification = workflow.Step{
	Name: "Send Slack Notification",
	If:   workflow.Secrets.Get("SLACK_WEBHOOK_URL").String() + " != ''",
	Run: `curl -X POST -H 'Content-type: application/json' \
  --data '{"text":"Release ${{ github.event.release.tag_name }} is now available!\n${{ github.event.release.html_url }}"}' \
  ${{ secrets.SLACK_WEBHOOK_URL }}`,
}

// NotifySteps are the steps for the notification job.
var NotifySteps = []any{
	PrintReleaseInfo,
	SlackNotification,
}
