package setup_rust

import (
	"testing"
)

func TestSetupRust_Action(t *testing.T) {
	a := SetupRust{}
	if got := a.Action(); got != "dtolnay/rust-toolchain@stable" {
		t.Errorf("Action() = %q, want %q", got, "dtolnay/rust-toolchain@stable")
	}
}

func TestSetupRust_ToStep(t *testing.T) {
	a := SetupRust{
		Toolchain: "stable",
	}

	step := a.ToStep()

	if step.Uses != "dtolnay/rust-toolchain@stable" {
		t.Errorf("Uses = %q, want %q", step.Uses, "dtolnay/rust-toolchain@stable")
	}

	if step.With["toolchain"] != "stable" {
		t.Errorf("With[toolchain] = %v, want %q", step.With["toolchain"], "stable")
	}
}

func TestSetupRust_ToStep_Empty(t *testing.T) {
	a := SetupRust{}
	step := a.ToStep()

	if step.Uses != "dtolnay/rust-toolchain@stable" {
		t.Errorf("Uses = %q, want %q", step.Uses, "dtolnay/rust-toolchain@stable")
	}

	if _, ok := step.With["toolchain"]; ok {
		t.Error("Empty toolchain should not be in With")
	}
}

func TestSetupRust_ToStep_WithComponents(t *testing.T) {
	a := SetupRust{
		Toolchain:  "nightly",
		Components: "rustfmt, clippy",
	}

	step := a.ToStep()

	if step.With["toolchain"] != "nightly" {
		t.Errorf("toolchain = %v, want %q", step.With["toolchain"], "nightly")
	}
	if step.With["components"] != "rustfmt, clippy" {
		t.Errorf("components = %v, want %q", step.With["components"], "rustfmt, clippy")
	}
}

func TestSetupRust_ToStep_WithTargets(t *testing.T) {
	a := SetupRust{
		Toolchain: "stable",
		Targets:   "wasm32-unknown-unknown",
	}

	step := a.ToStep()

	if step.With["targets"] != "wasm32-unknown-unknown" {
		t.Errorf("targets = %v, want %q", step.With["targets"], "wasm32-unknown-unknown")
	}
}

func TestSetupRust_Helpers(t *testing.T) {
	t.Run("Stable", func(t *testing.T) {
		a := Stable()
		if a.Toolchain != "stable" {
			t.Errorf("Toolchain = %q, want %q", a.Toolchain, "stable")
		}
	})

	t.Run("Nightly", func(t *testing.T) {
		a := Nightly()
		if a.Toolchain != "nightly" {
			t.Errorf("Toolchain = %q, want %q", a.Toolchain, "nightly")
		}
	})

	t.Run("Beta", func(t *testing.T) {
		a := Beta()
		if a.Toolchain != "beta" {
			t.Errorf("Toolchain = %q, want %q", a.Toolchain, "beta")
		}
	})
}

func TestSetupRust_ToStep_AllFields(t *testing.T) {
	a := SetupRust{
		Toolchain:  "1.70.0",
		Targets:    "wasm32-unknown-unknown, x86_64-unknown-linux-musl",
		Components: "rustfmt, clippy, llvm-tools-preview",
		Profile:    "default",
	}

	step := a.ToStep()

	if step.With["toolchain"] != "1.70.0" {
		t.Errorf("toolchain = %v, want %q", step.With["toolchain"], "1.70.0")
	}
	if step.With["targets"] != "wasm32-unknown-unknown, x86_64-unknown-linux-musl" {
		t.Errorf("targets = %v, want expected value", step.With["targets"])
	}
	if step.With["components"] != "rustfmt, clippy, llvm-tools-preview" {
		t.Errorf("components = %v, want expected value", step.With["components"])
	}
	if step.With["profile"] != "default" {
		t.Errorf("profile = %v, want %q", step.With["profile"], "default")
	}
}
