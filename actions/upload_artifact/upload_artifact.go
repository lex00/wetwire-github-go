// Package upload_artifact provides a typed wrapper for actions/upload-artifact.
package upload_artifact

// UploadArtifact wraps the actions/upload-artifact@v4 action.
// Upload a build artifact for use in subsequent jobs.
type UploadArtifact struct {
	// Artifact name
	Name string `yaml:"name,omitempty"`

	// A file, directory or wildcard pattern that describes what to upload
	Path string `yaml:"path,omitempty"`

	// The desired behavior if no files are found
	IfNoFilesFound string `yaml:"if-no-files-found,omitempty"`

	// Duration after which artifact will expire in days (0 means default retention)
	RetentionDays int `yaml:"retention-days,omitempty"`

	// The level of compression for Zlib (1-9, 6 is default)
	CompressionLevel int `yaml:"compression-level,omitempty"`

	// If true, an artifact with a matching name will be deleted before a new one is uploaded
	Overwrite bool `yaml:"overwrite,omitempty"`

	// Whether to include hidden files in the artifact
	IncludeHiddenFiles bool `yaml:"include-hidden-files,omitempty"`
}

// Action returns the action reference.
func (a UploadArtifact) Action() string {
	return "actions/upload-artifact@v4"
}

// Inputs returns the action inputs as a map.
func (a UploadArtifact) Inputs() map[string]any {
	with := make(map[string]any)

	if a.Name != "" {
		with["name"] = a.Name
	}
	if a.Path != "" {
		with["path"] = a.Path
	}
	if a.IfNoFilesFound != "" {
		with["if-no-files-found"] = a.IfNoFilesFound
	}
	if a.RetentionDays != 0 {
		with["retention-days"] = a.RetentionDays
	}
	if a.CompressionLevel != 0 {
		with["compression-level"] = a.CompressionLevel
	}
	if a.Overwrite {
		with["overwrite"] = a.Overwrite
	}
	if a.IncludeHiddenFiles {
		with["include-hidden-files"] = a.IncludeHiddenFiles
	}

	return with
}
