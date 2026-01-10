// Package cargo provides a typed wrapper for actions-rs/cargo.
package cargo

// Cargo wraps the actions-rs/cargo@v1 action.
// Run cargo commands for Rust projects.
type Cargo struct {
	// Cargo command to run. Examples: "build", "test", "check", "clippy", "fmt"
	Command string `yaml:"command,omitempty"`

	// Arguments to pass to the cargo command. Examples: "--release", "--all-features"
	Args string `yaml:"args,omitempty"`

	// Use cross instead of cargo. Useful for cross-compilation.
	UseCross bool `yaml:"use-cross,omitempty"`

	// Rust toolchain to use. Examples: "stable", "nightly", "1.70.0"
	Toolchain string `yaml:"toolchain,omitempty"`
}

// Action returns the action reference.
func (a Cargo) Action() string {
	return "actions-rs/cargo@v1"
}

// Inputs returns the action inputs as a map.
func (a Cargo) Inputs() map[string]any {
	with := make(map[string]any)

	if a.Command != "" {
		with["command"] = a.Command
	}
	if a.Args != "" {
		with["args"] = a.Args
	}
	if a.UseCross {
		with["use-cross"] = a.UseCross
	}
	if a.Toolchain != "" {
		with["toolchain"] = a.Toolchain
	}

	return with
}

// Build returns a Cargo action configured for cargo build.
func Build() Cargo {
	return Cargo{Command: "build"}
}

// Test returns a Cargo action configured for cargo test.
func Test() Cargo {
	return Cargo{Command: "test"}
}

// Check returns a Cargo action configured for cargo check.
func Check() Cargo {
	return Cargo{Command: "check"}
}

// Clippy returns a Cargo action configured for cargo clippy.
func Clippy() Cargo {
	return Cargo{Command: "clippy"}
}

// Fmt returns a Cargo action configured for cargo fmt.
func Fmt() Cargo {
	return Cargo{Command: "fmt"}
}
