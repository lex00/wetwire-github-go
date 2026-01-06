// Package dependabot provides types for Dependabot configuration.
package dependabot

// Dependabot represents a Dependabot configuration file.
type Dependabot struct {
	// Version must be 2.
	Version int `yaml:"version"`

	// EnableBetaEcosystems enables beta-level ecosystem support.
	EnableBetaEcosystems bool `yaml:"enable-beta-ecosystems,omitempty"`

	// Updates defines the package ecosystems to update.
	Updates []Update `yaml:"updates"`

	// Registries defines authentication for private registries.
	Registries map[string]Registry `yaml:"registries,omitempty"`
}

// ResourceType returns "dependabot" for interface compliance.
func (d Dependabot) ResourceType() string {
	return "dependabot"
}
