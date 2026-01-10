package hugo

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestHugo_Action(t *testing.T) {
	h := Hugo{}
	if got := h.Action(); got != "peaceiris/actions-hugo@v3" {
		t.Errorf("Action() = %q, want %q", got, "peaceiris/actions-hugo@v3")
	}
}

func TestHugo_Inputs_Empty(t *testing.T) {
	h := Hugo{}
	inputs := h.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty Hugo.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestHugo_Inputs_HugoVersion(t *testing.T) {
	h := Hugo{
		HugoVersion: "0.123.0",
	}

	inputs := h.Inputs()

	if inputs["hugo-version"] != "0.123.0" {
		t.Errorf("inputs[hugo-version] = %v, want %q", inputs["hugo-version"], "0.123.0")
	}
}

func TestHugo_Inputs_Extended(t *testing.T) {
	h := Hugo{
		Extended: true,
	}

	inputs := h.Inputs()

	if inputs["extended"] != true {
		t.Errorf("inputs[extended] = %v, want true", inputs["extended"])
	}
}

func TestHugo_Inputs_ExtendedFalse(t *testing.T) {
	h := Hugo{
		Extended: false,
	}

	inputs := h.Inputs()

	if _, exists := inputs["extended"]; exists {
		t.Errorf("inputs[extended] should not exist when Extended is false")
	}
}

func TestHugo_Inputs_GitHubToken(t *testing.T) {
	h := Hugo{
		GitHubToken: "${{ secrets.GITHUB_TOKEN }}",
	}

	inputs := h.Inputs()

	if inputs["github-token"] != "${{ secrets.GITHUB_TOKEN }}" {
		t.Errorf("inputs[github-token] = %v, want %q", inputs["github-token"], "${{ secrets.GITHUB_TOKEN }}")
	}
}

func TestHugo_Inputs_AllFields(t *testing.T) {
	h := Hugo{
		HugoVersion: "0.123.0",
		Extended:    true,
		GitHubToken: "${{ secrets.GITHUB_TOKEN }}",
	}

	inputs := h.Inputs()

	expected := map[string]any{
		"hugo-version": "0.123.0",
		"extended":     true,
		"github-token": "${{ secrets.GITHUB_TOKEN }}",
	}

	if len(inputs) != len(expected) {
		t.Errorf("inputs has %d entries, want %d", len(inputs), len(expected))
	}

	for key, want := range expected {
		if got := inputs[key]; got != want {
			t.Errorf("inputs[%q] = %v, want %v", key, got, want)
		}
	}
}

func TestHugo_ImplementsStepAction(t *testing.T) {
	h := Hugo{}
	// Verify Hugo implements StepAction interface
	var _ workflow.StepAction = h
}
