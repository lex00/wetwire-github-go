package upload_pages_artifact

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestUploadPagesArtifact_Action(t *testing.T) {
	u := UploadPagesArtifact{}
	if got := u.Action(); got != "actions/upload-pages-artifact@v3" {
		t.Errorf("Action() = %q, want %q", got, "actions/upload-pages-artifact@v3")
	}
}

func TestUploadPagesArtifact_Inputs_Empty(t *testing.T) {
	u := UploadPagesArtifact{}
	inputs := u.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty UploadPagesArtifact.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestUploadPagesArtifact_Inputs_AllFields(t *testing.T) {
	u := UploadPagesArtifact{
		Path:          "./dist",
		Name:          "github-pages",
		RetentionDays: 7,
	}

	inputs := u.Inputs()

	expected := map[string]any{
		"path":           "./dist",
		"name":           "github-pages",
		"retention-days": 7,
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

func TestUploadPagesArtifact_Inputs_Path(t *testing.T) {
	u := UploadPagesArtifact{
		Path: "./build",
	}

	inputs := u.Inputs()

	if inputs["path"] != "./build" {
		t.Errorf("inputs[path] = %v, want %q", inputs["path"], "./build")
	}
	if len(inputs) != 1 {
		t.Errorf("inputs has %d entries, want 1", len(inputs))
	}
}

func TestUploadPagesArtifact_Inputs_Name(t *testing.T) {
	u := UploadPagesArtifact{
		Name: "my-pages",
	}

	inputs := u.Inputs()

	if inputs["name"] != "my-pages" {
		t.Errorf("inputs[name] = %v, want %q", inputs["name"], "my-pages")
	}
	if len(inputs) != 1 {
		t.Errorf("inputs has %d entries, want 1", len(inputs))
	}
}

func TestUploadPagesArtifact_Inputs_RetentionDays(t *testing.T) {
	u := UploadPagesArtifact{
		RetentionDays: 14,
	}

	inputs := u.Inputs()

	if inputs["retention-days"] != 14 {
		t.Errorf("inputs[retention-days] = %v, want %d", inputs["retention-days"], 14)
	}
	if len(inputs) != 1 {
		t.Errorf("inputs has %d entries, want 1", len(inputs))
	}
}

func TestUploadPagesArtifact_ImplementsStepAction(t *testing.T) {
	u := UploadPagesArtifact{}
	// Verify UploadPagesArtifact implements StepAction interface
	var _ workflow.StepAction = u
}

func TestUploadPagesArtifact_InSteps(t *testing.T) {
	// Test that UploadPagesArtifact can be used in a steps slice
	steps := []any{
		UploadPagesArtifact{
			Path: "./public",
			Name: "github-pages",
		},
	}

	if len(steps) != 1 {
		t.Errorf("steps has %d entries, want 1", len(steps))
	}

	upa, ok := steps[0].(UploadPagesArtifact)
	if !ok {
		t.Fatal("steps[0] is not UploadPagesArtifact")
	}

	if upa.Path != "./public" {
		t.Errorf("Path = %q, want %q", upa.Path, "./public")
	}
	if upa.Name != "github-pages" {
		t.Errorf("Name = %q, want %q", upa.Name, "github-pages")
	}
}

func TestUploadPagesArtifact_Inputs_ZeroRetentionDays(t *testing.T) {
	// Test that RetentionDays = 0 is not included
	u := UploadPagesArtifact{
		RetentionDays: 0,
	}

	inputs := u.Inputs()

	if _, exists := inputs["retention-days"]; exists {
		t.Errorf("inputs[retention-days] should not exist for RetentionDays=0")
	}
}

func TestUploadPagesArtifact_CommonPaths(t *testing.T) {
	// Test common build output paths
	paths := []string{"./dist", "./build", "./public", "./_site", "./out"}

	for _, path := range paths {
		u := UploadPagesArtifact{
			Path: path,
		}
		inputs := u.Inputs()
		if inputs["path"] != path {
			t.Errorf("inputs[path] = %v, want %q", inputs["path"], path)
		}
	}
}
