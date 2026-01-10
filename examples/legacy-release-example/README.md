# Legacy Release Actions Example

This example demonstrates how to use the `actions/create-release` and `actions/upload-release-asset` action wrappers.

## Important Note

Both `actions/create-release` and `actions/upload-release-asset` are **deprecated** and unmaintained by GitHub.

For new projects, we recommend using modern alternatives:
- `softprops/action-gh-release` (see `gh_release.GHRelease` wrapper)
- `ncipollo/release-action`

This example is provided for legacy compatibility and migration scenarios.

## Example Workflow

```go
package main

import (
	"github.com/lex00/wetwire-github-go/actions/checkout"
	"github.com/lex00/wetwire-github-go/actions/create_release"
	"github.com/lex00/wetwire-github-go/actions/upload_release_asset"
	"github.com/lex00/wetwire-github-go/workflow"
)

// Workflow triggered on version tags
var ReleaseWorkflow = workflow.Workflow{
	Name: "Legacy Release",
	On:   ReleaseTriggers,
	Jobs: map[string]any{
		"release": ReleaseJob,
	},
}

var ReleaseTriggers = workflow.Triggers{
	Push: ReleasePush,
}

var ReleasePush = workflow.PushTrigger{
	Tags: []string{"v*"},
}

var ReleaseJob = workflow.Job{
	Name:   "Create Release",
	RunsOn: "ubuntu-latest",
	Steps:  ReleaseSteps,
}

var ReleaseSteps = []any{
	// Checkout code
	checkout.Checkout{},

	// Build artifacts (example)
	workflow.Step{
		Run: "go build -o myapp .",
	},

	// Create the release
	workflow.Step{
		ID:   "create_release",
		Uses: create_release.CreateRelease{
			TagName:     "${{ github.ref_name }}",
			ReleaseName: "Release ${{ github.ref_name }}",
			Body:        "See the changelog for details.",
			Draft:       false,
			Prerelease:  false,
		}.Action(),
		With: create_release.CreateRelease{
			TagName:     "${{ github.ref_name }}",
			ReleaseName: "Release ${{ github.ref_name }}",
			Body:        "See the changelog for details.",
			Draft:       false,
			Prerelease:  false,
		}.Inputs(),
	},

	// Upload binary asset
	upload_release_asset.UploadReleaseAsset{
		UploadURL:        "${{ steps.create_release.outputs.upload_url }}",
		AssetPath:        "./myapp",
		AssetName:        "myapp-linux-amd64",
		AssetContentType: "application/octet-stream",
	},
}
```

## Usage Patterns

### Creating a Release

```go
create_release.CreateRelease{
	TagName:     "v1.0.0",
	ReleaseName: "Version 1.0.0",
	Body:        "Release notes here",
	Draft:       false,
	Prerelease:  false,
}
```

### Uploading Release Assets

Upload binaries:

```go
upload_release_asset.UploadReleaseAsset{
	UploadURL:        "${{ steps.create_release.outputs.upload_url }}",
	AssetPath:        "./bin/myapp",
	AssetName:        "myapp-linux-amd64",
	AssetContentType: "application/octet-stream",
}
```

Upload archives:

```go
upload_release_asset.UploadReleaseAsset{
	UploadURL:        "${{ steps.create_release.outputs.upload_url }}",
	AssetPath:        "./dist/myapp.tar.gz",
	AssetName:        "myapp-v1.0.0.tar.gz",
	AssetContentType: "application/gzip",
}
```

## Common Content Types

- `application/octet-stream` - Binary files
- `application/gzip` - .tar.gz, .gz files
- `application/zip` - .zip files
- `application/json` - .json files
- `text/plain` - .txt files

## Migration to Modern Alternative

Modern approach using `gh_release.GHRelease`:

```go
import "github.com/lex00/wetwire-github-go/actions/gh_release"

gh_release.GHRelease{
	TagName:              "v1.0.0",
	Name:                 "Version 1.0.0",
	Body:                 "Release notes",
	Files:                "dist/*\nbin/*",
	GenerateReleaseNotes: true,
	Draft:                false,
	Prerelease:           false,
}
```

The `gh_release.GHRelease` wrapper is simpler and handles both release creation and asset uploads in a single step.
