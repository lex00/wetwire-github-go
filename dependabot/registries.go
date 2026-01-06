package dependabot

// Registry defines authentication for a private registry.
type Registry struct {
	// Type is the registry type.
	// Values: "composer-repository", "docker-registry", "git", "hex-organization",
	// "hex-repository", "maven-repository", "npm-registry", "nuget-feed",
	// "python-index", "rubygems-server", "terraform-registry".
	Type string `yaml:"type"`

	// URL is the registry endpoint.
	URL string `yaml:"url,omitempty"`

	// Username for authentication.
	Username string `yaml:"username,omitempty"`

	// Password for authentication.
	Password string `yaml:"password,omitempty"`

	// Token for token-based authentication.
	Token string `yaml:"token,omitempty"`

	// Key for key-based authentication (Hex).
	Key string `yaml:"key,omitempty"`

	// Organization for Hex organization.
	Organization string `yaml:"organization,omitempty"`

	// ReplacesBase indicates this replaces the base index (Python).
	ReplacesBase bool `yaml:"replaces-base,omitempty"`
}
