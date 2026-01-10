package setup_rust

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestSetupRust_Action(t *testing.T) {
	a := SetupRust{}
	if got := a.Action(); got != "dtolnay/rust-toolchain@stable" {
		t.Errorf("Action() = %q, want %q", got, "dtolnay/rust-toolchain@stable")
	}
}

func TestSetupRust_Inputs(t *testing.T) {
	a := SetupRust{
		Toolchain: "stable",
	}

	inputs := a.Inputs()

	if a.Action() != "dtolnay/rust-toolchain@stable" {
		t.Errorf("Action() = %q, want %q", a.Action(), "dtolnay/rust-toolchain@stable")
	}

	if inputs["toolchain"] != "stable" {
		t.Errorf("inputs[toolchain] = %v, want %q", inputs["toolchain"], "stable")
	}
}

func TestSetupRust_Inputs_Empty(t *testing.T) {
	a := SetupRust{}
	inputs := a.Inputs()

	if a.Action() != "dtolnay/rust-toolchain@stable" {
		t.Errorf("Action() = %q, want %q", a.Action(), "dtolnay/rust-toolchain@stable")
	}

	if _, ok := inputs["toolchain"]; ok {
		t.Error("Empty toolchain should not be in inputs")
	}
}

func TestSetupRust_Inputs_WithComponents(t *testing.T) {
	a := SetupRust{
		Toolchain:  "nightly",
		Components: "rustfmt, clippy",
	}

	inputs := a.Inputs()

	if inputs["toolchain"] != "nightly" {
		t.Errorf("toolchain = %v, want %q", inputs["toolchain"], "nightly")
	}
	if inputs["components"] != "rustfmt, clippy" {
		t.Errorf("components = %v, want %q", inputs["components"], "rustfmt, clippy")
	}
}

func TestSetupRust_Inputs_WithTargets(t *testing.T) {
	a := SetupRust{
		Toolchain: "stable",
		Targets:   "wasm32-unknown-unknown",
	}

	inputs := a.Inputs()

	if inputs["targets"] != "wasm32-unknown-unknown" {
		t.Errorf("targets = %v, want %q", inputs["targets"], "wasm32-unknown-unknown")
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

func TestSetupRust_Inputs_AllFields(t *testing.T) {
	a := SetupRust{
		Toolchain:  "1.70.0",
		Targets:    "wasm32-unknown-unknown, x86_64-unknown-linux-musl",
		Components: "rustfmt, clippy, llvm-tools-preview",
		Profile:    "default",
	}

	inputs := a.Inputs()

	if inputs["toolchain"] != "1.70.0" {
		t.Errorf("toolchain = %v, want %q", inputs["toolchain"], "1.70.0")
	}
	if inputs["targets"] != "wasm32-unknown-unknown, x86_64-unknown-linux-musl" {
		t.Errorf("targets = %v, want expected value", inputs["targets"])
	}
	if inputs["components"] != "rustfmt, clippy, llvm-tools-preview" {
		t.Errorf("components = %v, want expected value", inputs["components"])
	}
	if inputs["profile"] != "default" {
		t.Errorf("profile = %v, want %q", inputs["profile"], "default")
	}
}

func TestSetupRust_ImplementsStepAction(t *testing.T) {
	var _ workflow.StepAction = SetupRust{}
}
