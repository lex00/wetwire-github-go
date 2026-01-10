package upload_artifact

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestUploadArtifact_Action(t *testing.T) {
	a := UploadArtifact{}
	if got := a.Action(); got != "actions/upload-artifact@v4" {
		t.Errorf("Action() = %q, want %q", got, "actions/upload-artifact@v4")
	}
}

func TestUploadArtifact_Inputs(t *testing.T) {
	a := UploadArtifact{
		Name: "build-artifacts",
		Path: "dist/",
	}

	inputs := a.Inputs()

	if a.Action() != "actions/upload-artifact@v4" {
		t.Errorf("Action() = %q, want %q", a.Action(), "actions/upload-artifact@v4")
	}

	if inputs["name"] != "build-artifacts" {
		t.Errorf("inputs[name] = %v, want %q", inputs["name"], "build-artifacts")
	}

	if inputs["path"] != "dist/" {
		t.Errorf("inputs[path] = %v, want %q", inputs["path"], "dist/")
	}
}

func TestUploadArtifact_Inputs_Empty(t *testing.T) {
	a := UploadArtifact{}
	inputs := a.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty UploadArtifact.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestUploadArtifact_Inputs_AllFields(t *testing.T) {
	a := UploadArtifact{
		Name:               "test-artifact",
		Path:               "output/",
		IfNoFilesFound:     "error",
		RetentionDays:      7,
		CompressionLevel:   9,
		Overwrite:          true,
		IncludeHiddenFiles: true,
	}

	inputs := a.Inputs()

	if inputs["name"] != "test-artifact" {
		t.Errorf("inputs[name] = %v, want %q", inputs["name"], "test-artifact")
	}

	if inputs["path"] != "output/" {
		t.Errorf("inputs[path] = %v, want %q", inputs["path"], "output/")
	}

	if inputs["if-no-files-found"] != "error" {
		t.Errorf("inputs[if-no-files-found] = %v, want %q", inputs["if-no-files-found"], "error")
	}

	if inputs["retention-days"] != 7 {
		t.Errorf("inputs[retention-days] = %v, want 7", inputs["retention-days"])
	}

	if inputs["compression-level"] != 9 {
		t.Errorf("inputs[compression-level] = %v, want 9", inputs["compression-level"])
	}

	if inputs["overwrite"] != true {
		t.Errorf("inputs[overwrite] = %v, want true", inputs["overwrite"])
	}

	if inputs["include-hidden-files"] != true {
		t.Errorf("inputs[include-hidden-files] = %v, want true", inputs["include-hidden-files"])
	}
}

func TestUploadArtifact_ImplementsStepAction(t *testing.T) {
	var _ workflow.StepAction = UploadArtifact{}
}
