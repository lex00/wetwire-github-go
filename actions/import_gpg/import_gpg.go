// Package import_gpg provides a typed wrapper for crazy-max/ghaction-import-gpg.
package import_gpg

// ImportGPG wraps the crazy-max/ghaction-import-gpg@v6 action.
// Import a GPG key for signing commits, tags, and pushes.
type ImportGPG struct {
	// GPGPrivateKey is the GPG private key exported as an ASCII armored version.
	// Required input.
	GPGPrivateKey string `yaml:"gpg_private_key,omitempty"`

	// Passphrase is the passphrase of the GPG private key.
	Passphrase string `yaml:"passphrase,omitempty"`

	// GitUserSigningkey enables signing key for git.
	GitUserSigningkey bool `yaml:"git_user_signingkey,omitempty"`

	// GitCommitGpgsign enables commit signing.
	GitCommitGpgsign bool `yaml:"git_commit_gpgsign,omitempty"`

	// GitTagGpgsign enables tag signing.
	GitTagGpgsign bool `yaml:"git_tag_gpgsign,omitempty"`

	// GitPushGpgsign enables push signing.
	GitPushGpgsign bool `yaml:"git_push_gpgsign,omitempty"`

	// Fingerprint specifies the fingerprint of the GPG key to use.
	// Useful when you have multiple keys.
	Fingerprint string `yaml:"fingerprint,omitempty"`

	// TrustLevel sets the trust level for the GPG key.
	// Valid values: 1 (unknown), 2 (never), 3 (marginal), 4 (full), 5 (ultimate)
	TrustLevel string `yaml:"trust_level,omitempty"`

	// GitConfigGlobal sets git config globally.
	GitConfigGlobal bool `yaml:"git_config_global,omitempty"`

	// Workdir sets the working directory.
	Workdir string `yaml:"workdir,omitempty"`
}

// Action returns the action reference.
func (a ImportGPG) Action() string {
	return "crazy-max/ghaction-import-gpg@v6"
}

// Inputs returns the action inputs as a map.
func (a ImportGPG) Inputs() map[string]any {
	with := make(map[string]any)

	if a.GPGPrivateKey != "" {
		with["gpg_private_key"] = a.GPGPrivateKey
	}
	if a.Passphrase != "" {
		with["passphrase"] = a.Passphrase
	}
	if a.GitUserSigningkey {
		with["git_user_signingkey"] = a.GitUserSigningkey
	}
	if a.GitCommitGpgsign {
		with["git_commit_gpgsign"] = a.GitCommitGpgsign
	}
	if a.GitTagGpgsign {
		with["git_tag_gpgsign"] = a.GitTagGpgsign
	}
	if a.GitPushGpgsign {
		with["git_push_gpgsign"] = a.GitPushGpgsign
	}
	if a.Fingerprint != "" {
		with["fingerprint"] = a.Fingerprint
	}
	if a.TrustLevel != "" {
		with["trust_level"] = a.TrustLevel
	}
	if a.GitConfigGlobal {
		with["git_config_global"] = a.GitConfigGlobal
	}
	if a.Workdir != "" {
		with["workdir"] = a.Workdir
	}

	return with
}
