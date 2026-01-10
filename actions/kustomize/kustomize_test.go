package kustomize

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestKustomize_Action(t *testing.T) {
	k := Kustomize{}
	if got := k.Action(); got != "stefanprodan/kustomize-action@master" {
		t.Errorf("Action() = %q, want %q", got, "stefanprodan/kustomize-action@master")
	}
}

func TestKustomize_Inputs_Empty(t *testing.T) {
	k := Kustomize{}
	inputs := k.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty Kustomize.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestKustomize_Inputs_Kustomization(t *testing.T) {
	k := Kustomize{
		Kustomization: "./deploy/overlays/production",
	}

	inputs := k.Inputs()

	if inputs["kustomization"] != "./deploy/overlays/production" {
		t.Errorf("inputs[kustomization] = %v, want %q", inputs["kustomization"], "./deploy/overlays/production")
	}
}

func TestKustomize_Inputs_KustomizationBase(t *testing.T) {
	k := Kustomize{
		Kustomization: "./deploy/base",
	}

	inputs := k.Inputs()

	if inputs["kustomization"] != "./deploy/base" {
		t.Errorf("inputs[kustomization] = %v, want %q", inputs["kustomization"], "./deploy/base")
	}
}

func TestKustomize_Inputs_KustomizationAbsolute(t *testing.T) {
	k := Kustomize{
		Kustomization: "/path/to/kustomization",
	}

	inputs := k.Inputs()

	if inputs["kustomization"] != "/path/to/kustomization" {
		t.Errorf("inputs[kustomization] = %v, want %q", inputs["kustomization"], "/path/to/kustomization")
	}
}

func TestKustomize_ImplementsStepAction(t *testing.T) {
	k := Kustomize{}
	// Verify Kustomize implements StepAction interface
	var _ workflow.StepAction = k
}
