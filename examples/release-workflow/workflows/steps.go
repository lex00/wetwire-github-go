package workflows

import (
	"github.com/lex00/wetwire-github-go/actions/checkout"
	"github.com/lex00/wetwire-github-go/actions/gh_release"
	"github.com/lex00/wetwire-github-go/workflow"
)

// CheckoutWithHistory fetches full history for changelog generation.
var CheckoutWithHistory = checkout.Checkout{
	FetchDepth: 0,
}

// CreateGHRelease creates a GitHub release with auto-generated notes.
var CreateGHRelease = gh_release.GHRelease{
	GenerateReleaseNotes: true,
}

// GenerateChangelog generates a changelog from git commits.
var GenerateChangelog = workflow.Step{
	Name: "Generate Changelog",
	ID:   "changelog",
	Run: `echo "## Changes" > CHANGELOG.md
git log $(git describe --tags --abbrev=0 HEAD^)..HEAD --pretty=format:"- %s" >> CHANGELOG.md || echo "- Initial release" >> CHANGELOG.md
echo "" >> CHANGELOG.md`,
}

// ReleaseSteps are the steps for the release job.
var ReleaseSteps = []any{
	CheckoutWithHistory,
	GenerateChangelog,
	CreateGHRelease,
}
