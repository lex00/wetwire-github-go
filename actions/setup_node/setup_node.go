// Package setup_node provides a typed wrapper for actions/setup-node.
package setup_node

// SetupNode wraps the actions/setup-node@v4 action.
// Setup a Node.js environment and add it to PATH.
type SetupNode struct {
	// Version Spec of the version to use. Examples: 12.x, 10.15.1, >=10.15.0
	NodeVersion string `yaml:"node-version,omitempty"`

	// File containing the version Spec of the version to use
	NodeVersionFile string `yaml:"node-version-file,omitempty"`

	// Target architecture for Node to use
	Architecture string `yaml:"architecture,omitempty"`

	// Set this option if you want the action to check for latest available version
	CheckLatest bool `yaml:"check-latest,omitempty"`

	// Optional registry to set up for auth
	RegistryURL string `yaml:"registry-url,omitempty"`

	// Optional scope for authenticating against scoped registries
	Scope string `yaml:"scope,omitempty"`

	// Used to pull node distributions from node-versions
	Token string `yaml:"token,omitempty"`

	// Used to specify a package manager for caching in the default directory
	Cache string `yaml:"cache,omitempty"`

	// Used to specify the path to a dependency file
	CacheDependencyPath string `yaml:"cache-dependency-path,omitempty"`

	// Deprecated. Use node-version instead
	AlwaysAuth bool `yaml:"always-auth,omitempty"`
}

// Action returns the action reference.
func (a SetupNode) Action() string {
	return "actions/setup-node@v4"
}

// Inputs returns the action inputs as a map.
func (a SetupNode) Inputs() map[string]any {
	with := make(map[string]any)

	if a.NodeVersion != "" {
		with["node-version"] = a.NodeVersion
	}
	if a.NodeVersionFile != "" {
		with["node-version-file"] = a.NodeVersionFile
	}
	if a.Architecture != "" {
		with["architecture"] = a.Architecture
	}
	if a.CheckLatest {
		with["check-latest"] = a.CheckLatest
	}
	if a.RegistryURL != "" {
		with["registry-url"] = a.RegistryURL
	}
	if a.Scope != "" {
		with["scope"] = a.Scope
	}
	if a.Token != "" {
		with["token"] = a.Token
	}
	if a.Cache != "" {
		with["cache"] = a.Cache
	}
	if a.CacheDependencyPath != "" {
		with["cache-dependency-path"] = a.CacheDependencyPath
	}
	if a.AlwaysAuth {
		with["always-auth"] = a.AlwaysAuth
	}

	return with
}
