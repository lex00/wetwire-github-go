package setup_helm

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestSetupHelm_Action(t *testing.T) {
	s := SetupHelm{}
	if got := s.Action(); got != "azure/setup-helm@v4" {
		t.Errorf("Action() = %q, want %q", got, "azure/setup-helm@v4")
	}
}

func TestSetupHelm_Inputs_Empty(t *testing.T) {
	s := SetupHelm{}
	inputs := s.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty SetupHelm.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestSetupHelm_Inputs_Version(t *testing.T) {
	s := SetupHelm{
		Version: "v3.14.0",
	}

	inputs := s.Inputs()

	if inputs["version"] != "v3.14.0" {
		t.Errorf("inputs[version] = %v, want %q", inputs["version"], "v3.14.0")
	}
}

func TestSetupHelm_Inputs_Token(t *testing.T) {
	s := SetupHelm{
		Token: "${{ secrets.GITHUB_TOKEN }}",
	}

	inputs := s.Inputs()

	if inputs["token"] != "${{ secrets.GITHUB_TOKEN }}" {
		t.Errorf("inputs[token] = %v, want %q", inputs["token"], "${{ secrets.GITHUB_TOKEN }}")
	}
}

func TestSetupHelm_Inputs_DownloadBaseURL(t *testing.T) {
	s := SetupHelm{
		DownloadBaseURL: "https://custom.helm.sh",
	}

	inputs := s.Inputs()

	if inputs["downloadBaseURL"] != "https://custom.helm.sh" {
		t.Errorf("inputs[downloadBaseURL] = %v, want %q", inputs["downloadBaseURL"], "https://custom.helm.sh")
	}
}

func TestSetupHelm_Inputs_AllFields(t *testing.T) {
	s := SetupHelm{
		Version:         "v3.14.0",
		Token:           "${{ github.token }}",
		DownloadBaseURL: "https://get.helm.sh",
	}

	inputs := s.Inputs()

	expected := map[string]any{
		"version":         "v3.14.0",
		"token":           "${{ github.token }}",
		"downloadBaseURL": "https://get.helm.sh",
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

func TestSetupHelm_ImplementsStepAction(t *testing.T) {
	s := SetupHelm{}
	// Verify SetupHelm implements StepAction interface
	var _ workflow.StepAction = s
}
