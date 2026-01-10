package download_artifact

import (
	"testing"
)

func TestDownloadArtifact_Action(t *testing.T) {
	a := DownloadArtifact{}
	if got := a.Action(); got != "actions/download-artifact@v4" {
		t.Errorf("Action() = %q, want %q", got, "actions/download-artifact@v4")
	}
}

func TestDownloadArtifact_ToStep(t *testing.T) {
	a := DownloadArtifact{
		Name: "build-artifacts",
		Path: "dist/",
	}

	step := a.ToStep()

	if step.Uses != "actions/download-artifact@v4" {
		t.Errorf("step.Uses = %q, want %q", step.Uses, "actions/download-artifact@v4")
	}

	if step.With["name"] != "build-artifacts" {
		t.Errorf("step.With[name] = %v, want %q", step.With["name"], "build-artifacts")
	}

	if step.With["path"] != "dist/" {
		t.Errorf("step.With[path] = %v, want %q", step.With["path"], "dist/")
	}
}

func TestDownloadArtifact_ToStep_Empty(t *testing.T) {
	a := DownloadArtifact{}
	step := a.ToStep()

	if len(step.With) != 0 {
		t.Errorf("empty DownloadArtifact.ToStep() has %d with entries, want 0", len(step.With))
	}
}

func TestDownloadArtifact_ToStep_AllFields(t *testing.T) {
	a := DownloadArtifact{
		Name:          "test-artifact",
		Path:          "output/",
		Pattern:       "*.zip",
		MergeMultiple: true,
		GithubToken:   "token",
		Repository:    "owner/repo",
		RunID:         "12345",
	}

	step := a.ToStep()

	if step.With["name"] != "test-artifact" {
		t.Errorf("step.With[name] = %v, want %q", step.With["name"], "test-artifact")
	}

	if step.With["path"] != "output/" {
		t.Errorf("step.With[path] = %v, want %q", step.With["path"], "output/")
	}

	if step.With["pattern"] != "*.zip" {
		t.Errorf("step.With[pattern] = %v, want %q", step.With["pattern"], "*.zip")
	}

	if step.With["merge-multiple"] != true {
		t.Errorf("step.With[merge-multiple] = %v, want true", step.With["merge-multiple"])
	}

	if step.With["github-token"] != "token" {
		t.Errorf("step.With[github-token] = %v, want %q", step.With["github-token"], "token")
	}

	if step.With["repository"] != "owner/repo" {
		t.Errorf("step.With[repository] = %v, want %q", step.With["repository"], "owner/repo")
	}

	if step.With["run-id"] != "12345" {
		t.Errorf("step.With[run-id] = %v, want %q", step.With["run-id"], "12345")
	}
}
