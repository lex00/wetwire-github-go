// Package upload_pages_artifact provides a typed wrapper for actions/upload-pages-artifact.
package upload_pages_artifact

// UploadPagesArtifact wraps the actions/upload-pages-artifact@v3 action.
// Uploads an artifact for GitHub Pages deployment.
type UploadPagesArtifact struct {
	// Path to upload for Pages deployment
	Path string `yaml:"path,omitempty"`

	// Name of the artifact
	Name string `yaml:"name,omitempty"`

	// Number of days to retain the artifact
	RetentionDays int `yaml:"retention-days,omitempty"`
}

// Action returns the action reference.
func (a UploadPagesArtifact) Action() string {
	return "actions/upload-pages-artifact@v3"
}

// Inputs returns the action inputs as a map.
func (a UploadPagesArtifact) Inputs() map[string]any {
	m := make(map[string]any)

	if a.Path != "" {
		m["path"] = a.Path
	}
	if a.Name != "" {
		m["name"] = a.Name
	}
	if a.RetentionDays != 0 {
		m["retention-days"] = a.RetentionDays
	}

	return m
}
