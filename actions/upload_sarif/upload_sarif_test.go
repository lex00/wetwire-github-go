package upload_sarif

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestUploadSarif_Action(t *testing.T) {
	a := UploadSarif{}
	if got := a.Action(); got != "github/codeql-action/upload-sarif@v3" {
		t.Errorf("Action() = %q, want %q", got, "github/codeql-action/upload-sarif@v3")
	}
}

func TestUploadSarif_Inputs(t *testing.T) {
	a := UploadSarif{
		SarifFile: "results.sarif",
		Category:  "security",
	}

	inputs := a.Inputs()

	if inputs["sarif_file"] != "results.sarif" {
		t.Errorf("inputs[sarif_file] = %v, want results.sarif", inputs["sarif_file"])
	}
	if inputs["category"] != "security" {
		t.Errorf("inputs[category] = %v, want security", inputs["category"])
	}
}

func TestUploadSarif_Inputs_Empty(t *testing.T) {
	a := UploadSarif{}
	inputs := a.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty UploadSarif.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestUploadSarif_Inputs_CheckoutPath(t *testing.T) {
	a := UploadSarif{
		CheckoutPath: "/home/runner/work/repo",
	}

	inputs := a.Inputs()

	if inputs["checkout_path"] != "/home/runner/work/repo" {
		t.Errorf("inputs[checkout_path] = %v, want /home/runner/work/repo", inputs["checkout_path"])
	}
}

func TestUploadSarif_Inputs_RefAndSha(t *testing.T) {
	a := UploadSarif{
		Ref: "refs/heads/main",
		Sha: "abc123def456",
	}

	inputs := a.Inputs()

	if inputs["ref"] != "refs/heads/main" {
		t.Errorf("inputs[ref] = %v, want refs/heads/main", inputs["ref"])
	}
	if inputs["sha"] != "abc123def456" {
		t.Errorf("inputs[sha] = %v, want abc123def456", inputs["sha"])
	}
}

func TestUploadSarif_Inputs_Token(t *testing.T) {
	a := UploadSarif{
		Token: "${{ secrets.GITHUB_TOKEN }}",
	}

	inputs := a.Inputs()

	if inputs["token"] != "${{ secrets.GITHUB_TOKEN }}" {
		t.Errorf("inputs[token] = %v, want secret reference", inputs["token"])
	}
}

func TestUploadSarif_Inputs_WaitForProcessing(t *testing.T) {
	a := UploadSarif{
		WaitForProcessing: true,
	}

	inputs := a.Inputs()

	if inputs["wait-for-processing"] != true {
		t.Errorf("inputs[wait-for-processing] = %v, want true", inputs["wait-for-processing"])
	}
}

func TestUploadSarif_Inputs_WaitForProcessingFalse(t *testing.T) {
	a := UploadSarif{
		WaitForProcessing: false,
	}

	inputs := a.Inputs()

	if _, ok := inputs["wait-for-processing"]; ok {
		t.Error("wait-for-processing=false should not be in inputs")
	}
}

func TestUploadSarif_Inputs_AllFields(t *testing.T) {
	a := UploadSarif{
		SarifFile:         "./security/results.sarif",
		CheckoutPath:      "/workspace/repo",
		Ref:               "refs/pull/123/merge",
		Sha:               "fedcba987654321",
		Category:          "trivy-scan",
		Token:             "${{ secrets.SECURITY_TOKEN }}",
		WaitForProcessing: true,
	}

	inputs := a.Inputs()

	expected := map[string]any{
		"sarif_file":          "./security/results.sarif",
		"checkout_path":       "/workspace/repo",
		"ref":                 "refs/pull/123/merge",
		"sha":                 "fedcba987654321",
		"category":            "trivy-scan",
		"token":               "${{ secrets.SECURITY_TOKEN }}",
		"wait-for-processing": true,
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

func TestUploadSarif_ImplementsStepAction(t *testing.T) {
	a := UploadSarif{}
	var _ workflow.StepAction = a
}
