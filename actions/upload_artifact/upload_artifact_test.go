package upload_artifact

import (
	"testing"
)

func TestUploadArtifact_Action(t *testing.T) {
	a := UploadArtifact{}
	if got := a.Action(); got != "actions/upload-artifact@v4" {
		t.Errorf("Action() = %q, want %q", got, "actions/upload-artifact@v4")
	}
}

func TestUploadArtifact_ToStep(t *testing.T) {
	a := UploadArtifact{
		Name: "build-artifacts",
		Path: "dist/",
	}

	step := a.ToStep()

	if step.Uses != "actions/upload-artifact@v4" {
		t.Errorf("step.Uses = %q, want %q", step.Uses, "actions/upload-artifact@v4")
	}

	if step.With["name"] != "build-artifacts" {
		t.Errorf("step.With[name] = %v, want %q", step.With["name"], "build-artifacts")
	}

	if step.With["path"] != "dist/" {
		t.Errorf("step.With[path] = %v, want %q", step.With["path"], "dist/")
	}
}

func TestUploadArtifact_ToStep_Empty(t *testing.T) {
	a := UploadArtifact{}
	step := a.ToStep()

	if len(step.With) != 0 {
		t.Errorf("empty UploadArtifact.ToStep() has %d with entries, want 0", len(step.With))
	}
}

func TestUploadArtifact_ToStep_AllFields(t *testing.T) {
	a := UploadArtifact{
		Name:               "test-artifact",
		Path:               "output/",
		IfNoFilesFound:     "error",
		RetentionDays:      7,
		CompressionLevel:   9,
		Overwrite:          true,
		IncludeHiddenFiles: true,
	}

	step := a.ToStep()

	if step.With["name"] != "test-artifact" {
		t.Errorf("step.With[name] = %v, want %q", step.With["name"], "test-artifact")
	}

	if step.With["path"] != "output/" {
		t.Errorf("step.With[path] = %v, want %q", step.With["path"], "output/")
	}

	if step.With["if-no-files-found"] != "error" {
		t.Errorf("step.With[if-no-files-found] = %v, want %q", step.With["if-no-files-found"], "error")
	}

	if step.With["retention-days"] != 7 {
		t.Errorf("step.With[retention-days] = %v, want 7", step.With["retention-days"])
	}

	if step.With["compression-level"] != 9 {
		t.Errorf("step.With[compression-level] = %v, want 9", step.With["compression-level"])
	}

	if step.With["overwrite"] != true {
		t.Errorf("step.With[overwrite] = %v, want true", step.With["overwrite"])
	}

	if step.With["include-hidden-files"] != true {
		t.Errorf("step.With[include-hidden-files] = %v, want true", step.With["include-hidden-files"])
	}
}
