package upload_release_asset

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestUploadReleaseAsset_Action(t *testing.T) {
	a := UploadReleaseAsset{}
	if got := a.Action(); got != "actions/upload-release-asset@v1" {
		t.Errorf("Action() = %q, want %q", got, "actions/upload-release-asset@v1")
	}
}

func TestUploadReleaseAsset_Inputs_Empty(t *testing.T) {
	a := UploadReleaseAsset{}
	inputs := a.Inputs()

	if a.Action() != "actions/upload-release-asset@v1" {
		t.Errorf("Action() = %q, want %q", a.Action(), "actions/upload-release-asset@v1")
	}

	// Empty inputs should not include fields
	if _, ok := inputs["upload_url"]; ok {
		t.Error("Empty upload_url should not be in inputs")
	}
	if _, ok := inputs["asset_path"]; ok {
		t.Error("Empty asset_path should not be in inputs")
	}
	if _, ok := inputs["asset_name"]; ok {
		t.Error("Empty asset_name should not be in inputs")
	}
	if _, ok := inputs["asset_content_type"]; ok {
		t.Error("Empty asset_content_type should not be in inputs")
	}
}

func TestUploadReleaseAsset_Inputs_AllRequired(t *testing.T) {
	a := UploadReleaseAsset{
		UploadURL:        "${{ steps.create_release.outputs.upload_url }}",
		AssetPath:        "./dist/myapp.tar.gz",
		AssetName:        "myapp.tar.gz",
		AssetContentType: "application/gzip",
	}

	inputs := a.Inputs()

	if inputs["upload_url"] != "${{ steps.create_release.outputs.upload_url }}" {
		t.Errorf("upload_url = %v, want expected value", inputs["upload_url"])
	}
	if inputs["asset_path"] != "./dist/myapp.tar.gz" {
		t.Errorf("asset_path = %v, want %q", inputs["asset_path"], "./dist/myapp.tar.gz")
	}
	if inputs["asset_name"] != "myapp.tar.gz" {
		t.Errorf("asset_name = %v, want %q", inputs["asset_name"], "myapp.tar.gz")
	}
	if inputs["asset_content_type"] != "application/gzip" {
		t.Errorf("asset_content_type = %v, want %q", inputs["asset_content_type"], "application/gzip")
	}
}

func TestUploadReleaseAsset_Inputs_UploadURLOnly(t *testing.T) {
	a := UploadReleaseAsset{
		UploadURL: "https://uploads.github.com/repos/owner/repo/releases/1/assets{?name,label}",
	}

	inputs := a.Inputs()

	if inputs["upload_url"] != "https://uploads.github.com/repos/owner/repo/releases/1/assets{?name,label}" {
		t.Errorf("upload_url = %v, want expected value", inputs["upload_url"])
	}

	// Other fields should not be present
	if _, ok := inputs["asset_path"]; ok {
		t.Error("Empty asset_path should not be in inputs")
	}
}

func TestUploadReleaseAsset_Inputs_WithZip(t *testing.T) {
	a := UploadReleaseAsset{
		UploadURL:        "${{ steps.create_release.outputs.upload_url }}",
		AssetPath:        "./dist/myapp.zip",
		AssetName:        "myapp-v1.0.0.zip",
		AssetContentType: "application/zip",
	}

	inputs := a.Inputs()

	if inputs["asset_path"] != "./dist/myapp.zip" {
		t.Errorf("asset_path = %v, want %q", inputs["asset_path"], "./dist/myapp.zip")
	}
	if inputs["asset_name"] != "myapp-v1.0.0.zip" {
		t.Errorf("asset_name = %v, want %q", inputs["asset_name"], "myapp-v1.0.0.zip")
	}
	if inputs["asset_content_type"] != "application/zip" {
		t.Errorf("asset_content_type = %v, want %q", inputs["asset_content_type"], "application/zip")
	}
}

func TestUploadReleaseAsset_Inputs_WithBinary(t *testing.T) {
	a := UploadReleaseAsset{
		UploadURL:        "${{ steps.create_release.outputs.upload_url }}",
		AssetPath:        "./bin/myapp",
		AssetName:        "myapp-linux-amd64",
		AssetContentType: "application/octet-stream",
	}

	inputs := a.Inputs()

	if inputs["asset_path"] != "./bin/myapp" {
		t.Errorf("asset_path = %v, want %q", inputs["asset_path"], "./bin/myapp")
	}
	if inputs["asset_name"] != "myapp-linux-amd64" {
		t.Errorf("asset_name = %v, want %q", inputs["asset_name"], "myapp-linux-amd64")
	}
	if inputs["asset_content_type"] != "application/octet-stream" {
		t.Errorf("asset_content_type = %v, want %q", inputs["asset_content_type"], "application/octet-stream")
	}
}

func TestUploadReleaseAsset_Inputs_WithJSON(t *testing.T) {
	a := UploadReleaseAsset{
		UploadURL:        "${{ steps.create_release.outputs.upload_url }}",
		AssetPath:        "./manifest.json",
		AssetName:        "manifest.json",
		AssetContentType: "application/json",
	}

	inputs := a.Inputs()

	if inputs["asset_content_type"] != "application/json" {
		t.Errorf("asset_content_type = %v, want %q", inputs["asset_content_type"], "application/json")
	}
}

func TestUploadReleaseAsset_Inputs_WithText(t *testing.T) {
	a := UploadReleaseAsset{
		UploadURL:        "${{ steps.create_release.outputs.upload_url }}",
		AssetPath:        "./CHANGELOG.txt",
		AssetName:        "CHANGELOG.txt",
		AssetContentType: "text/plain",
	}

	inputs := a.Inputs()

	if inputs["asset_content_type"] != "text/plain" {
		t.Errorf("asset_content_type = %v, want %q", inputs["asset_content_type"], "text/plain")
	}
}

func TestUploadReleaseAsset_Inputs_ComplexPath(t *testing.T) {
	a := UploadReleaseAsset{
		UploadURL:        "${{ steps.create_release.outputs.upload_url }}",
		AssetPath:        "./dist/releases/v1.0.0/myapp-darwin-arm64.tar.gz",
		AssetName:        "myapp-darwin-arm64-v1.0.0.tar.gz",
		AssetContentType: "application/gzip",
	}

	inputs := a.Inputs()

	if inputs["asset_path"] != "./dist/releases/v1.0.0/myapp-darwin-arm64.tar.gz" {
		t.Errorf("asset_path = %v, want expected value", inputs["asset_path"])
	}
	if inputs["asset_name"] != "myapp-darwin-arm64-v1.0.0.tar.gz" {
		t.Errorf("asset_name = %v, want expected value", inputs["asset_name"])
	}
}

func TestUploadReleaseAsset_Inputs_WithExpression(t *testing.T) {
	a := UploadReleaseAsset{
		UploadURL:        "${{ steps.create_release.outputs.upload_url }}",
		AssetPath:        "${{ env.ASSET_PATH }}",
		AssetName:        "${{ env.ASSET_NAME }}",
		AssetContentType: "application/gzip",
	}

	inputs := a.Inputs()

	if inputs["asset_path"] != "${{ env.ASSET_PATH }}" {
		t.Errorf("asset_path = %v, want expression", inputs["asset_path"])
	}
	if inputs["asset_name"] != "${{ env.ASSET_NAME }}" {
		t.Errorf("asset_name = %v, want expression", inputs["asset_name"])
	}
}

func TestUploadReleaseAsset_ImplementsStepAction(t *testing.T) {
	var _ workflow.StepAction = UploadReleaseAsset{}
}
