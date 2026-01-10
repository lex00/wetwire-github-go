// Package docker_build_push provides a typed wrapper for docker/build-push-action.
package docker_build_push

// DockerBuildPush wraps the docker/build-push-action@v6 action.
// Build and push Docker images with Buildx.
type DockerBuildPush struct {
	// Build context. Path to the Dockerfile context.
	Context string `yaml:"context,omitempty"`

	// Path to the Dockerfile.
	File string `yaml:"file,omitempty"`

	// Push the image to the registry.
	Push bool `yaml:"push,omitempty"`

	// Load the image into the Docker daemon.
	Load bool `yaml:"load,omitempty"`

	// List of tags for the image (newline-delimited).
	Tags string `yaml:"tags,omitempty"`

	// List of build-time variables (newline-delimited).
	BuildArgs string `yaml:"build-args,omitempty"`

	// List of target platforms for build (comma-separated).
	Platforms string `yaml:"platforms,omitempty"`

	// External cache sources (e.g., type=gha).
	CacheFrom string `yaml:"cache-from,omitempty"`

	// Cache export destinations (e.g., type=gha,mode=max).
	CacheTo string `yaml:"cache-to,omitempty"`

	// Target stage to build.
	Target string `yaml:"target,omitempty"`

	// Do not use cache when building the image.
	NoCache bool `yaml:"no-cache,omitempty"`

	// Always attempt to pull all referenced images.
	Pull bool `yaml:"pull,omitempty"`

	// List of secrets to expose to the build (newline-delimited).
	Secrets string `yaml:"secrets,omitempty"`

	// List of metadata labels for the image (newline-delimited).
	Labels string `yaml:"labels,omitempty"`

	// List of output destinations.
	Outputs string `yaml:"outputs,omitempty"`

	// Generate provenance attestation. Can be "true", "false", or "mode=max".
	Provenance string `yaml:"provenance,omitempty"`

	// Generate SBOM attestation. Can be "true" or "false".
	SBOM string `yaml:"sbom,omitempty"`
}

// Action returns the action reference.
func (a DockerBuildPush) Action() string {
	return "docker/build-push-action@v6"
}

// Inputs returns the action inputs as a map.
func (a DockerBuildPush) Inputs() map[string]any {
	with := make(map[string]any)

	if a.Context != "" {
		with["context"] = a.Context
	}
	if a.File != "" {
		with["file"] = a.File
	}
	if a.Push {
		with["push"] = a.Push
	}
	if a.Load {
		with["load"] = a.Load
	}
	if a.Tags != "" {
		with["tags"] = a.Tags
	}
	if a.BuildArgs != "" {
		with["build-args"] = a.BuildArgs
	}
	if a.Platforms != "" {
		with["platforms"] = a.Platforms
	}
	if a.CacheFrom != "" {
		with["cache-from"] = a.CacheFrom
	}
	if a.CacheTo != "" {
		with["cache-to"] = a.CacheTo
	}
	if a.Target != "" {
		with["target"] = a.Target
	}
	if a.NoCache {
		with["no-cache"] = a.NoCache
	}
	if a.Pull {
		with["pull"] = a.Pull
	}
	if a.Secrets != "" {
		with["secrets"] = a.Secrets
	}
	if a.Labels != "" {
		with["labels"] = a.Labels
	}
	if a.Outputs != "" {
		with["outputs"] = a.Outputs
	}
	if a.Provenance != "" {
		with["provenance"] = a.Provenance
	}
	if a.SBOM != "" {
		with["sbom"] = a.SBOM
	}

	return with
}
