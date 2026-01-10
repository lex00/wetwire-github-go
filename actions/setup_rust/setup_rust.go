// Package setup_rust provides a typed wrapper for dtolnay/rust-toolchain.
package setup_rust

// SetupRust wraps the dtolnay/rust-toolchain@stable action.
// Install a Rust toolchain and add it to PATH.
type SetupRust struct {
	// Rust toolchain name to use. Examples: "stable", "nightly", "beta", "1.70.0"
	Toolchain string `yaml:"toolchain,omitempty"`

	// Target triple to add to the toolchain. Examples: "wasm32-unknown-unknown"
	Targets string `yaml:"targets,omitempty"`

	// Additional components to install. Examples: "rustfmt, clippy", "llvm-tools-preview"
	Components string `yaml:"components,omitempty"`

	// Override the default rustup profile (minimal)
	Profile string `yaml:"profile,omitempty"`
}

// Action returns the action reference.
func (a SetupRust) Action() string {
	return "dtolnay/rust-toolchain@stable"
}

// Inputs returns the action inputs as a map.
func (a SetupRust) Inputs() map[string]any {
	with := make(map[string]any)

	if a.Toolchain != "" {
		with["toolchain"] = a.Toolchain
	}
	if a.Targets != "" {
		with["targets"] = a.Targets
	}
	if a.Components != "" {
		with["components"] = a.Components
	}
	if a.Profile != "" {
		with["profile"] = a.Profile
	}

	return with
}

// Stable returns a SetupRust configured for the stable toolchain.
func Stable() SetupRust {
	return SetupRust{Toolchain: "stable"}
}

// Nightly returns a SetupRust configured for the nightly toolchain.
func Nightly() SetupRust {
	return SetupRust{Toolchain: "nightly"}
}

// Beta returns a SetupRust configured for the beta toolchain.
func Beta() SetupRust {
	return SetupRust{Toolchain: "beta"}
}
