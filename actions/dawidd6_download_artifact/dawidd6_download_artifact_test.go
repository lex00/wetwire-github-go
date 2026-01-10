package dawidd6_download_artifact

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestDownloadArtifact_Action(t *testing.T) {
	a := DownloadArtifact{}
	if got := a.Action(); got != "dawidd6/action-download-artifact@v6" {
		t.Errorf("Action() = %q, want %q", got, "dawidd6/action-download-artifact@v6")
	}
}

func TestDownloadArtifact_Inputs_Empty(t *testing.T) {
	a := DownloadArtifact{}
	inputs := a.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty DownloadArtifact.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestDownloadArtifact_ImplementsStepAction(t *testing.T) {
	var _ workflow.StepAction = DownloadArtifact{}
}

func TestDownloadArtifact_Inputs_GitHubToken(t *testing.T) {
	a := DownloadArtifact{
		GitHubToken: "${{ secrets.GITHUB_TOKEN }}",
	}

	inputs := a.Inputs()

	if inputs["github_token"] != "${{ secrets.GITHUB_TOKEN }}" {
		t.Errorf("inputs[github_token] = %v, want %q", inputs["github_token"], "${{ secrets.GITHUB_TOKEN }}")
	}
}

func TestDownloadArtifact_Inputs_Workflow(t *testing.T) {
	a := DownloadArtifact{
		Workflow: "build.yml",
	}

	inputs := a.Inputs()

	if inputs["workflow"] != "build.yml" {
		t.Errorf("inputs[workflow] = %v, want %q", inputs["workflow"], "build.yml")
	}
}

func TestDownloadArtifact_Inputs_Name(t *testing.T) {
	a := DownloadArtifact{
		Name: "my-artifact",
	}

	inputs := a.Inputs()

	if inputs["name"] != "my-artifact" {
		t.Errorf("inputs[name] = %v, want %q", inputs["name"], "my-artifact")
	}
}

func TestDownloadArtifact_Inputs_Path(t *testing.T) {
	a := DownloadArtifact{
		Path: "./artifacts",
	}

	inputs := a.Inputs()

	if inputs["path"] != "./artifacts" {
		t.Errorf("inputs[path] = %v, want %q", inputs["path"], "./artifacts")
	}
}

func TestDownloadArtifact_Inputs_Branch(t *testing.T) {
	a := DownloadArtifact{
		Branch: "main",
	}

	inputs := a.Inputs()

	if inputs["branch"] != "main" {
		t.Errorf("inputs[branch] = %v, want %q", inputs["branch"], "main")
	}
}

func TestDownloadArtifact_Inputs_Repo(t *testing.T) {
	a := DownloadArtifact{
		Repo: "owner/repo",
	}

	inputs := a.Inputs()

	if inputs["repo"] != "owner/repo" {
		t.Errorf("inputs[repo] = %v, want %q", inputs["repo"], "owner/repo")
	}
}

func TestDownloadArtifact_Inputs_RunID(t *testing.T) {
	a := DownloadArtifact{
		RunID: "123456789",
	}

	inputs := a.Inputs()

	if inputs["run_id"] != "123456789" {
		t.Errorf("inputs[run_id] = %v, want %q", inputs["run_id"], "123456789")
	}
}

func TestDownloadArtifact_Inputs_RunNumber(t *testing.T) {
	a := DownloadArtifact{
		RunNumber: "42",
	}

	inputs := a.Inputs()

	if inputs["run_number"] != "42" {
		t.Errorf("inputs[run_number] = %v, want %q", inputs["run_number"], "42")
	}
}

func TestDownloadArtifact_Inputs_IfNoArtifactFound(t *testing.T) {
	a := DownloadArtifact{
		IfNoArtifactFound: "warn",
	}

	inputs := a.Inputs()

	if inputs["if_no_artifact_found"] != "warn" {
		t.Errorf("inputs[if_no_artifact_found] = %v, want %q", inputs["if_no_artifact_found"], "warn")
	}
}

func TestDownloadArtifact_Inputs_AllowForks(t *testing.T) {
	a := DownloadArtifact{
		AllowForks: true,
	}

	inputs := a.Inputs()

	if inputs["allow_forks"] != true {
		t.Errorf("inputs[allow_forks] = %v, want true", inputs["allow_forks"])
	}
}

func TestDownloadArtifact_Inputs_CheckArtifacts(t *testing.T) {
	a := DownloadArtifact{
		CheckArtifacts: true,
	}

	inputs := a.Inputs()

	if inputs["check_artifacts"] != true {
		t.Errorf("inputs[check_artifacts] = %v, want true", inputs["check_artifacts"])
	}
}

func TestDownloadArtifact_Inputs_SearchArtifacts(t *testing.T) {
	a := DownloadArtifact{
		SearchArtifacts: true,
	}

	inputs := a.Inputs()

	if inputs["search_artifacts"] != true {
		t.Errorf("inputs[search_artifacts] = %v, want true", inputs["search_artifacts"])
	}
}

func TestDownloadArtifact_Inputs_AllFields(t *testing.T) {
	a := DownloadArtifact{
		GitHubToken:       "${{ secrets.GITHUB_TOKEN }}",
		Workflow:          "ci.yml",
		Name:              "build-output",
		Path:              "./downloads",
		Branch:            "develop",
		Repo:              "myorg/myrepo",
		RunID:             "987654321",
		RunNumber:         "100",
		IfNoArtifactFound: "error",
		AllowForks:        true,
		CheckArtifacts:    true,
		SearchArtifacts:   true,
	}

	inputs := a.Inputs()

	expected := map[string]any{
		"github_token":        "${{ secrets.GITHUB_TOKEN }}",
		"workflow":            "ci.yml",
		"name":                "build-output",
		"path":                "./downloads",
		"branch":              "develop",
		"repo":                "myorg/myrepo",
		"run_id":              "987654321",
		"run_number":          "100",
		"if_no_artifact_found": "error",
		"allow_forks":         true,
		"check_artifacts":     true,
		"search_artifacts":    true,
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

func TestDownloadArtifact_Inputs_FalseBoolFields(t *testing.T) {
	a := DownloadArtifact{
		AllowForks:      false,
		CheckArtifacts:  false,
		SearchArtifacts: false,
	}

	inputs := a.Inputs()

	if len(inputs) != 0 {
		t.Errorf("inputs for false bools has %d entries, want 0. Got: %v", len(inputs), inputs)
	}
}
