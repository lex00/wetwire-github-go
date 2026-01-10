package fossa

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestFossa_Action(t *testing.T) {
	f := Fossa{}
	if got := f.Action(); got != "fossas/fossa-action@v1" {
		t.Errorf("Action() = %q, want %q", got, "fossas/fossa-action@v1")
	}
}

func TestFossa_Inputs_Empty(t *testing.T) {
	f := Fossa{}
	inputs := f.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty Fossa.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestFossa_Inputs_APIKey(t *testing.T) {
	f := Fossa{
		APIKey: "fossa_api_key_123",
	}

	inputs := f.Inputs()

	if inputs["api-key"] != "fossa_api_key_123" {
		t.Errorf("inputs[api-key] = %v, want %q", inputs["api-key"], "fossa_api_key_123")
	}
}

func TestFossa_Inputs_Branch(t *testing.T) {
	f := Fossa{
		Branch: "feature/new-feature",
	}

	inputs := f.Inputs()

	if inputs["branch"] != "feature/new-feature" {
		t.Errorf("inputs[branch] = %v, want %q", inputs["branch"], "feature/new-feature")
	}
}

func TestFossa_Inputs_Revision(t *testing.T) {
	f := Fossa{
		Revision: "abc123def456",
	}

	inputs := f.Inputs()

	if inputs["revision"] != "abc123def456" {
		t.Errorf("inputs[revision] = %v, want %q", inputs["revision"], "abc123def456")
	}
}

func TestFossa_Inputs_Container(t *testing.T) {
	f := Fossa{
		Container: "ghcr.io/fossas/fossa-cli:latest",
	}

	inputs := f.Inputs()

	if inputs["container"] != "ghcr.io/fossas/fossa-cli:latest" {
		t.Errorf("inputs[container] = %v, want %q", inputs["container"], "ghcr.io/fossas/fossa-cli:latest")
	}
}

func TestFossa_ImplementsStepAction(t *testing.T) {
	f := Fossa{}
	// Verify Fossa implements StepAction interface
	var _ workflow.StepAction = f
}

func TestFossa_Inputs_AllFields(t *testing.T) {
	f := Fossa{
		APIKey:    "my-api-key",
		Branch:    "main",
		Revision:  "sha256abc",
		Container: "fossa/cli:v3",
	}

	inputs := f.Inputs()

	expected := map[string]any{
		"api-key":   "my-api-key",
		"branch":    "main",
		"revision":  "sha256abc",
		"container": "fossa/cli:v3",
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

func TestFossa_Inputs_Combined_APIKeyAndBranch(t *testing.T) {
	// Test typical use case: API key + branch override
	f := Fossa{
		APIKey: "secret_key",
		Branch: "develop",
	}

	inputs := f.Inputs()

	if len(inputs) != 2 {
		t.Errorf("inputs has %d entries, want 2", len(inputs))
	}

	if inputs["api-key"] != "secret_key" {
		t.Errorf("inputs[api-key] = %v, want %q", inputs["api-key"], "secret_key")
	}
	if inputs["branch"] != "develop" {
		t.Errorf("inputs[branch] = %v, want %q", inputs["branch"], "develop")
	}
}

func TestFossa_Inputs_Combined_APIKeyAndRevision(t *testing.T) {
	// Test typical use case: API key + revision override
	f := Fossa{
		APIKey:   "secret_key",
		Revision: "v1.0.0",
	}

	inputs := f.Inputs()

	if len(inputs) != 2 {
		t.Errorf("inputs has %d entries, want 2", len(inputs))
	}

	if inputs["api-key"] != "secret_key" {
		t.Errorf("inputs[api-key] = %v, want %q", inputs["api-key"], "secret_key")
	}
	if inputs["revision"] != "v1.0.0" {
		t.Errorf("inputs[revision] = %v, want %q", inputs["revision"], "v1.0.0")
	}
}

func TestFossa_Inputs_OnlyContainer(t *testing.T) {
	// Test using custom container without other options
	f := Fossa{
		Container: "custom-fossa:1.0",
	}

	inputs := f.Inputs()

	if len(inputs) != 1 {
		t.Errorf("inputs has %d entries, want 1", len(inputs))
	}

	if inputs["container"] != "custom-fossa:1.0" {
		t.Errorf("inputs[container] = %v, want %q", inputs["container"], "custom-fossa:1.0")
	}
}
