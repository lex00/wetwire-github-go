package cosign_installer

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestCosignInstaller_Action(t *testing.T) {
	a := CosignInstaller{}
	if got := a.Action(); got != "sigstore/cosign-installer@v3" {
		t.Errorf("Action() = %q, want %q", got, "sigstore/cosign-installer@v3")
	}
}

func TestCosignInstaller_Inputs_Empty(t *testing.T) {
	a := CosignInstaller{}
	inputs := a.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty CosignInstaller.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestCosignInstaller_Inputs_CosignRelease(t *testing.T) {
	a := CosignInstaller{
		CosignRelease: "v2.2.0",
	}

	inputs := a.Inputs()

	if inputs["cosign-release"] != "v2.2.0" {
		t.Errorf("inputs[cosign-release] = %v, want %q", inputs["cosign-release"], "v2.2.0")
	}
}

func TestCosignInstaller_Inputs_InstallDir(t *testing.T) {
	a := CosignInstaller{
		InstallDir: "/usr/local/bin",
	}

	inputs := a.Inputs()

	if inputs["install-dir"] != "/usr/local/bin" {
		t.Errorf("inputs[install-dir] = %v, want %q", inputs["install-dir"], "/usr/local/bin")
	}
}

func TestCosignInstaller_Inputs_AllFields(t *testing.T) {
	a := CosignInstaller{
		CosignRelease: "v2.1.1",
		InstallDir:    "$HOME/.local/bin",
	}

	inputs := a.Inputs()

	if inputs["cosign-release"] != "v2.1.1" {
		t.Errorf("inputs[cosign-release] = %v, want %q", inputs["cosign-release"], "v2.1.1")
	}
	if inputs["install-dir"] != "$HOME/.local/bin" {
		t.Errorf("inputs[install-dir] = %v, want %q", inputs["install-dir"], "$HOME/.local/bin")
	}
	if len(inputs) != 2 {
		t.Errorf("inputs has %d entries, want 2", len(inputs))
	}
}

func TestCosignInstaller_ImplementsStepAction(t *testing.T) {
	var _ workflow.StepAction = CosignInstaller{}
}
