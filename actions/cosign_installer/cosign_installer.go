// Package cosign_installer provides a typed wrapper for sigstore/cosign-installer.
package cosign_installer

// CosignInstaller wraps the sigstore/cosign-installer@v3 action.
// Install Cosign for signing and verifying container images.
type CosignInstaller struct {
	// Cosign release version to be installed. Examples: "v2.2.0", "v2.1.1"
	CosignRelease string `yaml:"cosign-release,omitempty"`

	// Where to install the cosign binary. Default: $HOME/.cosign
	InstallDir string `yaml:"install-dir,omitempty"`
}

// Action returns the action reference.
func (a CosignInstaller) Action() string {
	return "sigstore/cosign-installer@v3"
}

// Inputs returns the action inputs as a map.
func (a CosignInstaller) Inputs() map[string]any {
	with := make(map[string]any)

	if a.CosignRelease != "" {
		with["cosign-release"] = a.CosignRelease
	}
	if a.InstallDir != "" {
		with["install-dir"] = a.InstallDir
	}

	return with
}
