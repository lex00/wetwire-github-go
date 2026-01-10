package download_artifact

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestDownloadArtifact_Action(t *testing.T) {
	a := DownloadArtifact{}
	if got := a.Action(); got != "actions/download-artifact@v4" {
		t.Errorf("Action() = %q, want %q", got, "actions/download-artifact@v4")
	}
}

func TestDownloadArtifact_Inputs(t *testing.T) {
	a := DownloadArtifact{
		Name: "build-artifacts",
		Path: "dist/",
	}

	inputs := a.Inputs()

	if a.Action() != "actions/download-artifact@v4" {
		t.Errorf("Action() = %q, want %q", a.Action(), "actions/download-artifact@v4")
	}

	if inputs["name"] != "build-artifacts" {
		t.Errorf("inputs[name] = %v, want %q", inputs["name"], "build-artifacts")
	}

	if inputs["path"] != "dist/" {
		t.Errorf("inputs[path] = %v, want %q", inputs["path"], "dist/")
	}
}

func TestDownloadArtifact_Inputs_Empty(t *testing.T) {
	a := DownloadArtifact{}
	inputs := a.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty DownloadArtifact.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestDownloadArtifact_Inputs_AllFields(t *testing.T) {
	a := DownloadArtifact{
		Name:          "test-artifact",
		Path:          "output/",
		Pattern:       "*.zip",
		MergeMultiple: true,
		GithubToken:   "token",
		Repository:    "owner/repo",
		RunID:         "12345",
	}

	inputs := a.Inputs()

	if inputs["name"] != "test-artifact" {
		t.Errorf("inputs[name] = %v, want %q", inputs["name"], "test-artifact")
	}

	if inputs["path"] != "output/" {
		t.Errorf("inputs[path] = %v, want %q", inputs["path"], "output/")
	}

	if inputs["pattern"] != "*.zip" {
		t.Errorf("inputs[pattern] = %v, want %q", inputs["pattern"], "*.zip")
	}

	if inputs["merge-multiple"] != true {
		t.Errorf("inputs[merge-multiple] = %v, want true", inputs["merge-multiple"])
	}

	if inputs["github-token"] != "token" {
		t.Errorf("inputs[github-token] = %v, want %q", inputs["github-token"], "token")
	}

	if inputs["repository"] != "owner/repo" {
		t.Errorf("inputs[repository] = %v, want %q", inputs["repository"], "owner/repo")
	}

	if inputs["run-id"] != "12345" {
		t.Errorf("inputs[run-id] = %v, want %q", inputs["run-id"], "12345")
	}
}

func TestDownloadArtifact_ImplementsStepAction(t *testing.T) {
	var _ workflow.StepAction = DownloadArtifact{}
}
