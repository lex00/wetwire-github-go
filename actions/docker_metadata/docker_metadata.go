// Package docker_metadata provides a typed wrapper for docker/metadata-action.
package docker_metadata

// DockerMetadata wraps the docker/metadata-action@v5 action.
// GitHub Action to extract metadata (tags, labels) from Git reference and GitHub events for Docker.
type DockerMetadata struct {
	// Where to get context data. Allowed options are: workflow (default), git.
	Context string `yaml:"context,omitempty"`

	// List of Docker images to use as base name for tags (newline-delimited).
	Images string `yaml:"images,omitempty"`

	// List of tags as key-value pair attributes (newline-delimited).
	Tags string `yaml:"tags,omitempty"`

	// Flavor to apply (newline-delimited).
	Flavor string `yaml:"flavor,omitempty"`

	// List of custom labels (newline-delimited).
	Labels string `yaml:"labels,omitempty"`

	// List of custom annotations (newline-delimited).
	Annotations string `yaml:"annotations,omitempty"`

	// Separator to use for tags output (default \n).
	SepTags string `yaml:"sep-tags,omitempty"`

	// Separator to use for labels output (default \n).
	SepLabels string `yaml:"sep-labels,omitempty"`

	// Separator to use for annotations output (default \n).
	SepAnnotations string `yaml:"sep-annotations,omitempty"`

	// Bake target name (default docker-metadata-action).
	BakeTarget string `yaml:"bake-target,omitempty"`
}

// Action returns the action reference.
func (a DockerMetadata) Action() string {
	return "docker/metadata-action@v5"
}

// Inputs returns the action inputs as a map.
func (a DockerMetadata) Inputs() map[string]any {
	with := make(map[string]any)

	if a.Context != "" {
		with["context"] = a.Context
	}
	if a.Images != "" {
		with["images"] = a.Images
	}
	if a.Tags != "" {
		with["tags"] = a.Tags
	}
	if a.Flavor != "" {
		with["flavor"] = a.Flavor
	}
	if a.Labels != "" {
		with["labels"] = a.Labels
	}
	if a.Annotations != "" {
		with["annotations"] = a.Annotations
	}
	if a.SepTags != "" {
		with["sep-tags"] = a.SepTags
	}
	if a.SepLabels != "" {
		with["sep-labels"] = a.SepLabels
	}
	if a.SepAnnotations != "" {
		with["sep-annotations"] = a.SepAnnotations
	}
	if a.BakeTarget != "" {
		with["bake-target"] = a.BakeTarget
	}

	return with
}
