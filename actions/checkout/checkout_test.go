package checkout

import (
	"testing"
)

func TestCheckout_Action(t *testing.T) {
	c := Checkout{}
	if got := c.Action(); got != "actions/checkout@v4" {
		t.Errorf("Action() = %q, want %q", got, "actions/checkout@v4")
	}
}

func TestCheckout_ToStep(t *testing.T) {
	c := Checkout{
		Repository: "owner/repo",
		Ref:        "main",
		FetchDepth: 1,
		Submodules: "recursive",
	}

	step := c.ToStep()

	if step.Uses != "actions/checkout@v4" {
		t.Errorf("step.Uses = %q, want %q", step.Uses, "actions/checkout@v4")
	}

	if step.With["repository"] != "owner/repo" {
		t.Errorf("step.With[repository] = %v, want %q", step.With["repository"], "owner/repo")
	}

	if step.With["ref"] != "main" {
		t.Errorf("step.With[ref] = %v, want %q", step.With["ref"], "main")
	}

	if step.With["fetch-depth"] != 1 {
		t.Errorf("step.With[fetch-depth] = %v, want 1", step.With["fetch-depth"])
	}

	if step.With["submodules"] != "recursive" {
		t.Errorf("step.With[submodules] = %v, want %q", step.With["submodules"], "recursive")
	}
}

func TestCheckout_ToStep_EmptyWithMaps(t *testing.T) {
	c := Checkout{}
	step := c.ToStep()

	// Empty checkout should have no with values
	if len(step.With) != 0 {
		t.Errorf("empty Checkout.ToStep() has %d with entries, want 0", len(step.With))
	}
}

func TestCheckout_ToStep_BoolFields(t *testing.T) {
	c := Checkout{
		Clean:           true,
		LFS:             true,
		PersistCredentials: true,
	}

	step := c.ToStep()

	if step.With["clean"] != true {
		t.Errorf("step.With[clean] = %v, want true", step.With["clean"])
	}

	if step.With["lfs"] != true {
		t.Errorf("step.With[lfs] = %v, want true", step.With["lfs"])
	}

	if step.With["persist-credentials"] != true {
		t.Errorf("step.With[persist-credentials] = %v, want true", step.With["persist-credentials"])
	}
}
