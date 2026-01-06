package dependabot

// Group defines dependency grouping rules.
type Group struct {
	// Patterns are glob patterns for dependency names.
	Patterns []string `yaml:"patterns,omitempty"`

	// DependencyType matches dependencies by type.
	// Values: "direct", "indirect", "all", "production", "development".
	DependencyType string `yaml:"dependency-type,omitempty"`

	// UpdateTypes specifies update types for this group.
	// Values: "minor", "patch".
	UpdateTypes []string `yaml:"update-types,omitempty"`

	// ExcludePatterns are glob patterns for dependencies to exclude.
	ExcludePatterns []string `yaml:"exclude-patterns,omitempty"`

	// AppliesTo specifies which update types this applies to.
	// Values: "version-updates", "security-updates".
	AppliesTo string `yaml:"applies-to,omitempty"`
}
