package actions_rs_toolchain

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestToolchain_Action(t *testing.T) {
	a := Toolchain{}
	if got := a.Action(); got != "actions-rs/toolchain@v1" {
		t.Errorf("Action() = %q, want %q", got, "actions-rs/toolchain@v1")
	}
}

func TestToolchain_Inputs_Empty(t *testing.T) {
	a := Toolchain{}
	inputs := a.Inputs()

	if a.Action() != "actions-rs/toolchain@v1" {
		t.Errorf("Action() = %q, want %q", a.Action(), "actions-rs/toolchain@v1")
	}

	if _, ok := inputs["toolchain"]; ok {
		t.Error("Empty toolchain should not be in inputs")
	}
	if _, ok := inputs["target"]; ok {
		t.Error("Empty target should not be in inputs")
	}
	if _, ok := inputs["default"]; ok {
		t.Error("Empty default should not be in inputs")
	}
	if _, ok := inputs["override"]; ok {
		t.Error("Empty override should not be in inputs")
	}
	if _, ok := inputs["profile"]; ok {
		t.Error("Empty profile should not be in inputs")
	}
	if _, ok := inputs["components"]; ok {
		t.Error("Empty components should not be in inputs")
	}
}

func TestToolchain_Inputs_WithToolchain(t *testing.T) {
	a := Toolchain{
		ToolchainName: "stable",
	}

	inputs := a.Inputs()

	if inputs["toolchain"] != "stable" {
		t.Errorf("inputs[toolchain] = %v, want %q", inputs["toolchain"], "stable")
	}
}

func TestToolchain_Inputs_WithTarget(t *testing.T) {
	a := Toolchain{
		ToolchainName: "stable",
		Target:        "wasm32-unknown-unknown",
	}

	inputs := a.Inputs()

	if inputs["toolchain"] != "stable" {
		t.Errorf("toolchain = %v, want %q", inputs["toolchain"], "stable")
	}
	if inputs["target"] != "wasm32-unknown-unknown" {
		t.Errorf("target = %v, want %q", inputs["target"], "wasm32-unknown-unknown")
	}
}

func TestToolchain_Inputs_WithDefaultTrue(t *testing.T) {
	a := Toolchain{
		ToolchainName: "nightly",
		Default:       true,
	}

	inputs := a.Inputs()

	if inputs["default"] != true {
		t.Errorf("default = %v, want true", inputs["default"])
	}
}

func TestToolchain_Inputs_WithDefaultFalse(t *testing.T) {
	a := Toolchain{
		ToolchainName: "stable",
		Default:       false,
	}

	inputs := a.Inputs()

	if _, ok := inputs["default"]; ok {
		t.Error("default=false should not be in inputs (omitted when false)")
	}
}

func TestToolchain_Inputs_WithOverrideTrue(t *testing.T) {
	a := Toolchain{
		ToolchainName: "stable",
		Override:      true,
	}

	inputs := a.Inputs()

	if inputs["override"] != true {
		t.Errorf("override = %v, want true", inputs["override"])
	}
}

func TestToolchain_Inputs_WithOverrideFalse(t *testing.T) {
	a := Toolchain{
		ToolchainName: "stable",
		Override:      false,
	}

	inputs := a.Inputs()

	if _, ok := inputs["override"]; ok {
		t.Error("override=false should not be in inputs (omitted when false)")
	}
}

func TestToolchain_Inputs_WithProfile(t *testing.T) {
	a := Toolchain{
		ToolchainName: "stable",
		Profile:       "minimal",
	}

	inputs := a.Inputs()

	if inputs["profile"] != "minimal" {
		t.Errorf("profile = %v, want %q", inputs["profile"], "minimal")
	}
}

func TestToolchain_Inputs_WithComponents(t *testing.T) {
	a := Toolchain{
		ToolchainName: "nightly",
		Components:    "rustfmt, clippy",
	}

	inputs := a.Inputs()

	if inputs["components"] != "rustfmt, clippy" {
		t.Errorf("components = %v, want %q", inputs["components"], "rustfmt, clippy")
	}
}

func TestToolchain_Inputs_AllFields(t *testing.T) {
	a := Toolchain{
		ToolchainName: "1.70.0",
		Target:        "wasm32-unknown-unknown",
		Default:       true,
		Override:      true,
		Profile:       "default",
		Components:    "rustfmt, clippy, llvm-tools-preview",
	}

	inputs := a.Inputs()

	if inputs["toolchain"] != "1.70.0" {
		t.Errorf("toolchain = %v, want %q", inputs["toolchain"], "1.70.0")
	}
	if inputs["target"] != "wasm32-unknown-unknown" {
		t.Errorf("target = %v, want expected value", inputs["target"])
	}
	if inputs["default"] != true {
		t.Errorf("default = %v, want true", inputs["default"])
	}
	if inputs["override"] != true {
		t.Errorf("override = %v, want true", inputs["override"])
	}
	if inputs["profile"] != "default" {
		t.Errorf("profile = %v, want %q", inputs["profile"], "default")
	}
	if inputs["components"] != "rustfmt, clippy, llvm-tools-preview" {
		t.Errorf("components = %v, want expected value", inputs["components"])
	}
}

func TestToolchain_Helpers_Stable(t *testing.T) {
	a := Stable()
	if a.ToolchainName != "stable" {
		t.Errorf("ToolchainName = %q, want %q", a.ToolchainName, "stable")
	}
}

func TestToolchain_Helpers_Nightly(t *testing.T) {
	a := Nightly()
	if a.ToolchainName != "nightly" {
		t.Errorf("ToolchainName = %q, want %q", a.ToolchainName, "nightly")
	}
}

func TestToolchain_Helpers_Beta(t *testing.T) {
	a := Beta()
	if a.ToolchainName != "beta" {
		t.Errorf("ToolchainName = %q, want %q", a.ToolchainName, "beta")
	}
}

func TestToolchain_ImplementsStepAction(t *testing.T) {
	var _ workflow.StepAction = Toolchain{}
}
