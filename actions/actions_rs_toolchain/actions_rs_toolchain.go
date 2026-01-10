// Package actions_rs_toolchain provides a typed wrapper for actions-rs/toolchain.
package actions_rs_toolchain

// Toolchain wraps the actions-rs/toolchain@v1 action.
// Install the Rust toolchain and add it to PATH.
type Toolchain struct {
	// Rust toolchain name. See https://rust-lang.github.io/rustup/concepts/toolchains.html#toolchain-specification
	// If not given, the action will try and install the version specified in the `rust-toolchain` file.
	// Examples: "stable", "nightly", "beta", "1.70.0"
	ToolchainName string `yaml:"toolchain,omitempty"`

	// Target triple to install for this toolchain. Examples: "wasm32-unknown-unknown"
	Target string `yaml:"target,omitempty"`

	// Set installed toolchain as default
	Default bool `yaml:"default,omitempty"`

	// Set installed toolchain as an override for a directory
	Override bool `yaml:"override,omitempty"`

	// Name of the group of components to be installed for a new toolchain.
	// Examples: "minimal", "default", "complete"
	Profile string `yaml:"profile,omitempty"`

	// Comma-separated list of components to be additionally installed for a new toolchain.
	// Examples: "rustfmt, clippy", "llvm-tools-preview"
	Components string `yaml:"components,omitempty"`
}

// Action returns the action reference.
func (a Toolchain) Action() string {
	return "actions-rs/toolchain@v1"
}

// Inputs returns the action inputs as a map.
func (a Toolchain) Inputs() map[string]any {
	with := make(map[string]any)

	if a.ToolchainName != "" {
		with["toolchain"] = a.ToolchainName
	}
	if a.Target != "" {
		with["target"] = a.Target
	}
	if a.Default {
		with["default"] = a.Default
	}
	if a.Override {
		with["override"] = a.Override
	}
	if a.Profile != "" {
		with["profile"] = a.Profile
	}
	if a.Components != "" {
		with["components"] = a.Components
	}

	return with
}

// Stable returns a Toolchain configured for the stable toolchain.
func Stable() Toolchain {
	return Toolchain{ToolchainName: "stable"}
}

// Nightly returns a Toolchain configured for the nightly toolchain.
func Nightly() Toolchain {
	return Toolchain{ToolchainName: "nightly"}
}

// Beta returns a Toolchain configured for the beta toolchain.
func Beta() Toolchain {
	return Toolchain{ToolchainName: "beta"}
}
