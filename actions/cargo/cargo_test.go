package cargo

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestCargo_Action(t *testing.T) {
	a := Cargo{}
	if got := a.Action(); got != "actions-rs/cargo@v1" {
		t.Errorf("Action() = %q, want %q", got, "actions-rs/cargo@v1")
	}
}

func TestCargo_Inputs_Empty(t *testing.T) {
	a := Cargo{}
	inputs := a.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty Cargo.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestCargo_Inputs_Command(t *testing.T) {
	a := Cargo{
		Command: "build",
	}

	inputs := a.Inputs()

	if inputs["command"] != "build" {
		t.Errorf("inputs[command] = %v, want %q", inputs["command"], "build")
	}
}

func TestCargo_Inputs_Args(t *testing.T) {
	a := Cargo{
		Command: "build",
		Args:    "--release",
	}

	inputs := a.Inputs()

	if inputs["args"] != "--release" {
		t.Errorf("inputs[args] = %v, want %q", inputs["args"], "--release")
	}
}

func TestCargo_Inputs_UseCross_True(t *testing.T) {
	a := Cargo{
		Command:  "build",
		UseCross: true,
	}

	inputs := a.Inputs()

	if inputs["use-cross"] != true {
		t.Errorf("inputs[use-cross] = %v, want true", inputs["use-cross"])
	}
}

func TestCargo_Inputs_UseCross_False(t *testing.T) {
	a := Cargo{
		Command:  "build",
		UseCross: false,
	}

	inputs := a.Inputs()

	if _, ok := inputs["use-cross"]; ok {
		t.Error("use-cross=false should not be in inputs (omitted when false)")
	}
}

func TestCargo_Inputs_Toolchain(t *testing.T) {
	a := Cargo{
		Command:   "build",
		Toolchain: "nightly",
	}

	inputs := a.Inputs()

	if inputs["toolchain"] != "nightly" {
		t.Errorf("inputs[toolchain] = %v, want %q", inputs["toolchain"], "nightly")
	}
}

func TestCargo_Inputs_AllFields(t *testing.T) {
	a := Cargo{
		Command:   "test",
		Args:      "--all-features --no-fail-fast",
		UseCross:  true,
		Toolchain: "stable",
	}

	inputs := a.Inputs()

	if inputs["command"] != "test" {
		t.Errorf("inputs[command] = %v, want %q", inputs["command"], "test")
	}
	if inputs["args"] != "--all-features --no-fail-fast" {
		t.Errorf("inputs[args] = %v, want %q", inputs["args"], "--all-features --no-fail-fast")
	}
	if inputs["use-cross"] != true {
		t.Errorf("inputs[use-cross] = %v, want true", inputs["use-cross"])
	}
	if inputs["toolchain"] != "stable" {
		t.Errorf("inputs[toolchain] = %v, want %q", inputs["toolchain"], "stable")
	}
	if len(inputs) != 4 {
		t.Errorf("inputs has %d entries, want 4", len(inputs))
	}
}

func TestCargo_Helpers_Build(t *testing.T) {
	a := Build()
	if a.Command != "build" {
		t.Errorf("Command = %q, want %q", a.Command, "build")
	}
}

func TestCargo_Helpers_Test(t *testing.T) {
	a := Test()
	if a.Command != "test" {
		t.Errorf("Command = %q, want %q", a.Command, "test")
	}
}

func TestCargo_Helpers_Check(t *testing.T) {
	a := Check()
	if a.Command != "check" {
		t.Errorf("Command = %q, want %q", a.Command, "check")
	}
}

func TestCargo_Helpers_Clippy(t *testing.T) {
	a := Clippy()
	if a.Command != "clippy" {
		t.Errorf("Command = %q, want %q", a.Command, "clippy")
	}
}

func TestCargo_Helpers_Fmt(t *testing.T) {
	a := Fmt()
	if a.Command != "fmt" {
		t.Errorf("Command = %q, want %q", a.Command, "fmt")
	}
}

func TestCargo_ImplementsStepAction(t *testing.T) {
	var _ workflow.StepAction = Cargo{}
}
