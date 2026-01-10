// Package upload_release_asset provides a typed wrapper for actions/upload-release-asset.
package upload_release_asset

// UploadReleaseAsset wraps the actions/upload-release-asset@v1 action.
// Upload a release asset to an existing release in your repository.
//
// Note: This action is deprecated by GitHub and unmaintained.
// Consider using softprops/action-gh-release or ncipollo/release-action instead.
type UploadReleaseAsset struct {
	// UploadURL is the URL for uploading assets to the release (required).
	// This typically comes from the create-release action's upload_url output.
	UploadURL string `yaml:"upload_url,omitempty"`

	// AssetPath is the path to the asset you want to upload (required).
	AssetPath string `yaml:"asset_path,omitempty"`

	// AssetName is the name of the asset you want to upload (required).
	AssetName string `yaml:"asset_name,omitempty"`

	// AssetContentType is the content-type of the asset you want to upload (required).
	// See https://www.iana.org/assignments/media-types/media-types.xhtml for supported types.
	// Common values: application/zip, application/gzip, application/octet-stream, text/plain
	AssetContentType string `yaml:"asset_content_type,omitempty"`
}

// Action returns the action reference.
func (a UploadReleaseAsset) Action() string {
	return "actions/upload-release-asset@v1"
}

// Inputs returns the action inputs as a map.
func (a UploadReleaseAsset) Inputs() map[string]any {
	with := make(map[string]any)

	if a.UploadURL != "" {
		with["upload_url"] = a.UploadURL
	}
	if a.AssetPath != "" {
		with["asset_path"] = a.AssetPath
	}
	if a.AssetName != "" {
		with["asset_name"] = a.AssetName
	}
	if a.AssetContentType != "" {
		with["asset_content_type"] = a.AssetContentType
	}

	return with
}
