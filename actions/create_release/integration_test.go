package create_release

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

// TestCreateRelease_Integration verifies the action can be used in a workflow step.
func TestCreateRelease_Integration(t *testing.T) {
	action := CreateRelease{
		TagName:     "v1.0.0",
		ReleaseName: "Release 1.0.0",
		Body:        "Initial release",
		Draft:       false,
		Prerelease:  false,
	}

	// Convert to Step using ToStep
	step := workflow.ToStep(action)

	if step.Uses != "actions/create-release@v1" {
		t.Errorf("Uses = %q, want %q", step.Uses, "actions/create-release@v1")
	}

	if step.With["tag_name"] != "v1.0.0" {
		t.Errorf("With[tag_name] = %v, want %q", step.With["tag_name"], "v1.0.0")
	}

	if step.With["release_name"] != "Release 1.0.0" {
		t.Errorf("With[release_name] = %v, want %q", step.With["release_name"], "Release 1.0.0")
	}

	if step.With["body"] != "Initial release" {
		t.Errorf("With[body] = %v, want %q", step.With["body"], "Initial release")
	}
}

// TestCreateRelease_UsedInStepsSlice verifies the action can be used directly in []any{} steps.
func TestCreateRelease_UsedInStepsSlice(t *testing.T) {
	steps := []any{
		CreateRelease{
			TagName:     "v2.0.0",
			ReleaseName: "Major Release 2.0.0",
		},
	}

	if len(steps) != 1 {
		t.Fatalf("Expected 1 step, got %d", len(steps))
	}

	action, ok := steps[0].(CreateRelease)
	if !ok {
		t.Fatal("Step is not a CreateRelease")
	}

	if action.TagName != "v2.0.0" {
		t.Errorf("TagName = %q, want %q", action.TagName, "v2.0.0")
	}
}

// TestCreateRelease_WithStepID demonstrates using the action with a step ID.
func TestCreateRelease_WithStepID(t *testing.T) {
	action := CreateRelease{
		TagName:     "v1.0.0",
		ReleaseName: "Release 1.0.0",
	}

	step := workflow.ToStep(action)
	step.ID = "create_release"

	if step.ID != "create_release" {
		t.Errorf("ID = %q, want %q", step.ID, "create_release")
	}

	// Verify outputs can be referenced
	uploadURL := step.Output("upload_url")
	if uploadURL.StepID != "create_release" {
		t.Errorf("OutputRef.StepID = %q, want %q", uploadURL.StepID, "create_release")
	}
	if uploadURL.Output != "upload_url" {
		t.Errorf("OutputRef.Output = %q, want %q", uploadURL.Output, "upload_url")
	}
}
