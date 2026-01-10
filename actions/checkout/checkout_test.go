package checkout

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestCheckout_Action(t *testing.T) {
	c := Checkout{}
	if got := c.Action(); got != "actions/checkout@v4" {
		t.Errorf("Action() = %q, want %q", got, "actions/checkout@v4")
	}
}

func TestCheckout_Inputs(t *testing.T) {
	c := Checkout{
		Repository: "owner/repo",
		Ref:        "main",
		FetchDepth: 1,
		Submodules: "recursive",
	}

	inputs := c.Inputs()

	if inputs["repository"] != "owner/repo" {
		t.Errorf("inputs[repository] = %v, want %q", inputs["repository"], "owner/repo")
	}

	if inputs["ref"] != "main" {
		t.Errorf("inputs[ref] = %v, want %q", inputs["ref"], "main")
	}

	if inputs["fetch-depth"] != 1 {
		t.Errorf("inputs[fetch-depth] = %v, want 1", inputs["fetch-depth"])
	}

	if inputs["submodules"] != "recursive" {
		t.Errorf("inputs[submodules] = %v, want %q", inputs["submodules"], "recursive")
	}
}

func TestCheckout_Inputs_Empty(t *testing.T) {
	c := Checkout{}
	inputs := c.Inputs()

	// Empty checkout should have no inputs
	if len(inputs) != 0 {
		t.Errorf("empty Checkout.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestCheckout_Inputs_BoolFields(t *testing.T) {
	c := Checkout{
		Clean:              true,
		LFS:                true,
		PersistCredentials: true,
	}

	inputs := c.Inputs()

	if inputs["clean"] != true {
		t.Errorf("inputs[clean] = %v, want true", inputs["clean"])
	}

	if inputs["lfs"] != true {
		t.Errorf("inputs[lfs] = %v, want true", inputs["lfs"])
	}

	if inputs["persist-credentials"] != true {
		t.Errorf("inputs[persist-credentials] = %v, want true", inputs["persist-credentials"])
	}
}

func TestCheckout_ImplementsStepAction(t *testing.T) {
	c := Checkout{}
	// Verify Checkout implements StepAction interface
	var _ workflow.StepAction = c
}
