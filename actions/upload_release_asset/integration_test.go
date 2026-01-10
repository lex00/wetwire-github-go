package upload_release_asset

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

// TestUploadReleaseAsset_Integration verifies the action can be used in a workflow step.
func TestUploadReleaseAsset_Integration(t *testing.T) {
	action := UploadReleaseAsset{
		UploadURL:        "${{ steps.create_release.outputs.upload_url }}",
		AssetPath:        "./dist/myapp.tar.gz",
		AssetName:        "myapp.tar.gz",
		AssetContentType: "application/gzip",
	}

	// Convert to Step using ToStep
	step := workflow.ToStep(action)

	if step.Uses != "actions/upload-release-asset@v1" {
		t.Errorf("Uses = %q, want %q", step.Uses, "actions/upload-release-asset@v1")
	}

	if step.With["upload_url"] != "${{ steps.create_release.outputs.upload_url }}" {
		t.Errorf("With[upload_url] = %v, want expression", step.With["upload_url"])
	}

	if step.With["asset_path"] != "./dist/myapp.tar.gz" {
		t.Errorf("With[asset_path] = %v, want %q", step.With["asset_path"], "./dist/myapp.tar.gz")
	}
}

// TestUploadReleaseAsset_UsedInStepsSlice verifies the action can be used directly in []any{} steps.
func TestUploadReleaseAsset_UsedInStepsSlice(t *testing.T) {
	steps := []any{
		UploadReleaseAsset{
			UploadURL:        "${{ steps.create_release.outputs.upload_url }}",
			AssetPath:        "./myapp",
			AssetName:        "myapp-linux",
			AssetContentType: "application/octet-stream",
		},
	}

	if len(steps) != 1 {
		t.Fatalf("Expected 1 step, got %d", len(steps))
	}

	action, ok := steps[0].(UploadReleaseAsset)
	if !ok {
		t.Fatal("Step is not a UploadReleaseAsset")
	}

	if action.AssetName != "myapp-linux" {
		t.Errorf("AssetName = %q, want %q", action.AssetName, "myapp-linux")
	}
}

// TestUploadReleaseAsset_MultipleAssets demonstrates uploading multiple assets.
func TestUploadReleaseAsset_MultipleAssets(t *testing.T) {
	steps := []any{
		UploadReleaseAsset{
			UploadURL:        "${{ steps.create_release.outputs.upload_url }}",
			AssetPath:        "./dist/myapp.tar.gz",
			AssetName:        "myapp.tar.gz",
			AssetContentType: "application/gzip",
		},
		UploadReleaseAsset{
			UploadURL:        "${{ steps.create_release.outputs.upload_url }}",
			AssetPath:        "./dist/myapp.zip",
			AssetName:        "myapp.zip",
			AssetContentType: "application/zip",
		},
		UploadReleaseAsset{
			UploadURL:        "${{ steps.create_release.outputs.upload_url }}",
			AssetPath:        "./CHANGELOG.md",
			AssetName:        "CHANGELOG.md",
			AssetContentType: "text/plain",
		},
	}

	if len(steps) != 3 {
		t.Fatalf("Expected 3 steps, got %d", len(steps))
	}

	// Verify each asset has the correct content type
	contentTypes := []string{"application/gzip", "application/zip", "text/plain"}
	for i, step := range steps {
		action, ok := step.(UploadReleaseAsset)
		if !ok {
			t.Fatalf("Step %d is not a UploadReleaseAsset", i)
		}
		if action.AssetContentType != contentTypes[i] {
			t.Errorf("Step %d: AssetContentType = %q, want %q", i, action.AssetContentType, contentTypes[i])
		}
	}
}
