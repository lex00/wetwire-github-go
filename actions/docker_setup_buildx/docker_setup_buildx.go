// Package docker_setup_buildx provides a typed wrapper for docker/setup-buildx-action.
package docker_setup_buildx

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

// DockerSetupBuildx wraps the docker/setup-buildx-action@v3 action.
// Set up Docker Buildx for multi-platform builds and advanced features.
type DockerSetupBuildx struct {
	// Buildx version. (e.g., "v0.12.0", "latest")
	Version string `yaml:"version,omitempty"`

	// Driver to use. (e.g., "docker-container", "kubernetes", "remote")
	Driver string `yaml:"driver,omitempty"`

	// Driver options (newline-delimited key=value pairs).
	DriverOpts string `yaml:"driver-opts,omitempty"`

	// Flags for buildkitd daemon.
	BuildkitdFlags string `yaml:"buildkitd-flags,omitempty"`

	// Install buildx as default docker builder.
	Install bool `yaml:"install,omitempty"`

	// Switch to this builder instance.
	Use bool `yaml:"use,omitempty"`

	// Address for a custom builder endpoint.
	Endpoint string `yaml:"endpoint,omitempty"`

	// Fixed platforms for current node.
	Platforms string `yaml:"platforms,omitempty"`

	// BuildKit config file.
	Config string `yaml:"config,omitempty"`

	// Inline BuildKit config.
	ConfigInline string `yaml:"config-inline,omitempty"`

	// Append additional nodes to the builder.
	Append string `yaml:"append,omitempty"`

	// Remove builder when job completes.
	Cleanup bool `yaml:"cleanup,omitempty"`
}

// Action returns the action reference.
func (a DockerSetupBuildx) Action() string {
	return "docker/setup-buildx-action@v3"
}

// ToStep converts this action to a workflow step.
func (a DockerSetupBuildx) ToStep() workflow.Step {
	with := make(workflow.With)

	if a.Version != "" {
		with["version"] = a.Version
	}
	if a.Driver != "" {
		with["driver"] = a.Driver
	}
	if a.DriverOpts != "" {
		with["driver-opts"] = a.DriverOpts
	}
	if a.BuildkitdFlags != "" {
		with["buildkitd-flags"] = a.BuildkitdFlags
	}
	if a.Install {
		with["install"] = a.Install
	}
	if a.Use {
		with["use"] = a.Use
	}
	if a.Endpoint != "" {
		with["endpoint"] = a.Endpoint
	}
	if a.Platforms != "" {
		with["platforms"] = a.Platforms
	}
	if a.Config != "" {
		with["config"] = a.Config
	}
	if a.ConfigInline != "" {
		with["config-inline"] = a.ConfigInline
	}
	if a.Append != "" {
		with["append"] = a.Append
	}
	if a.Cleanup {
		with["cleanup"] = a.Cleanup
	}

	return workflow.Step{
		Uses: a.Action(),
		With: with,
	}
}
