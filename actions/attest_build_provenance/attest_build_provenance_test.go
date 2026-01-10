package attest_build_provenance

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestAttestBuildProvenance_Action(t *testing.T) {
	a := AttestBuildProvenance{}
	if got := a.Action(); got != "actions/attest-build-provenance@v1" {
		t.Errorf("Action() = %q, want %q", got, "actions/attest-build-provenance@v1")
	}
}

func TestAttestBuildProvenance_Inputs(t *testing.T) {
	a := AttestBuildProvenance{
		SubjectPath: "dist/*.tar.gz",
		GithubToken: "${{ github.token }}",
	}

	inputs := a.Inputs()

	if inputs["subject-path"] != "dist/*.tar.gz" {
		t.Errorf("inputs[subject-path] = %v, want %q", inputs["subject-path"], "dist/*.tar.gz")
	}

	if inputs["github-token"] != "${{ github.token }}" {
		t.Errorf("inputs[github-token] = %v, want %q", inputs["github-token"], "${{ github.token }}")
	}
}

func TestAttestBuildProvenance_Inputs_Empty(t *testing.T) {
	a := AttestBuildProvenance{}
	inputs := a.Inputs()

	// Empty attestation should have no inputs
	if len(inputs) != 0 {
		t.Errorf("empty AttestBuildProvenance.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestAttestBuildProvenance_Inputs_BoolFields(t *testing.T) {
	a := AttestBuildProvenance{
		PushToRegistry:      true,
		CreateStorageRecord: true,
		ShowSummary:         true,
	}

	inputs := a.Inputs()

	if inputs["push-to-registry"] != true {
		t.Errorf("inputs[push-to-registry] = %v, want true", inputs["push-to-registry"])
	}

	if inputs["create-storage-record"] != true {
		t.Errorf("inputs[create-storage-record] = %v, want true", inputs["create-storage-record"])
	}

	if inputs["show-summary"] != true {
		t.Errorf("inputs[show-summary] = %v, want true", inputs["show-summary"])
	}
}

func TestAttestBuildProvenance_ImplementsStepAction(t *testing.T) {
	a := AttestBuildProvenance{}
	// Verify AttestBuildProvenance implements StepAction interface
	var _ workflow.StepAction = a
}

func TestAttestBuildProvenance_Inputs_SubjectDigest(t *testing.T) {
	a := AttestBuildProvenance{
		SubjectDigest: "sha256:abcd1234",
		SubjectName:   "my-artifact",
	}

	inputs := a.Inputs()

	if inputs["subject-digest"] != "sha256:abcd1234" {
		t.Errorf("inputs[subject-digest] = %v, want %q", inputs["subject-digest"], "sha256:abcd1234")
	}

	if inputs["subject-name"] != "my-artifact" {
		t.Errorf("inputs[subject-name] = %v, want %q", inputs["subject-name"], "my-artifact")
	}
}

func TestAttestBuildProvenance_Inputs_SubjectChecksums(t *testing.T) {
	a := AttestBuildProvenance{
		SubjectChecksums: "checksums.txt",
	}

	inputs := a.Inputs()

	if inputs["subject-checksums"] != "checksums.txt" {
		t.Errorf("inputs[subject-checksums] = %v, want %q", inputs["subject-checksums"], "checksums.txt")
	}
}

func TestAttestBuildProvenance_Inputs_AllFields(t *testing.T) {
	a := AttestBuildProvenance{
		SubjectPath:         "build/*.bin",
		SubjectDigest:       "sha256:1234abcd",
		SubjectName:         "my-binary",
		SubjectChecksums:    "SHA256SUMS",
		PushToRegistry:      true,
		CreateStorageRecord: false,
		ShowSummary:         false,
		GithubToken:         "ghp_token123",
	}

	inputs := a.Inputs()

	// Verify all non-empty fields are present
	expected := map[string]any{
		"subject-path":      "build/*.bin",
		"subject-digest":    "sha256:1234abcd",
		"subject-name":      "my-binary",
		"subject-checksums": "SHA256SUMS",
		"push-to-registry":  true,
		"github-token":      "ghp_token123",
	}

	if len(inputs) != len(expected) {
		t.Errorf("inputs has %d entries, want %d. Got: %v", len(inputs), len(expected), inputs)
	}

	for key, want := range expected {
		if got := inputs[key]; got != want {
			t.Errorf("inputs[%q] = %v, want %v", key, got, want)
		}
	}
}

func TestAttestBuildProvenance_Inputs_FalseBoolFields(t *testing.T) {
	// Test that false boolean values are not included in inputs
	a := AttestBuildProvenance{
		PushToRegistry:      false,
		CreateStorageRecord: false,
		ShowSummary:         false,
	}

	inputs := a.Inputs()

	// None of these should be in the inputs map
	if len(inputs) != 0 {
		t.Errorf("inputs for false bools has %d entries, want 0. Got: %v", len(inputs), inputs)
	}
}

func TestAttestBuildProvenance_Inputs_SubjectPathOnly(t *testing.T) {
	a := AttestBuildProvenance{
		SubjectPath: "dist/myapp-*.tar.gz",
	}

	inputs := a.Inputs()

	if len(inputs) != 1 {
		t.Errorf("inputs has %d entries, want 1. Got: %v", len(inputs), inputs)
	}

	if inputs["subject-path"] != "dist/myapp-*.tar.gz" {
		t.Errorf("inputs[subject-path] = %v, want %q", inputs["subject-path"], "dist/myapp-*.tar.gz")
	}
}

func TestAttestBuildProvenance_Inputs_SubjectDigestWithName(t *testing.T) {
	a := AttestBuildProvenance{
		SubjectDigest: "sha256:fedcba9876543210",
		SubjectName:   "container-image:v1.2.3",
	}

	inputs := a.Inputs()

	if len(inputs) != 2 {
		t.Errorf("inputs has %d entries, want 2. Got: %v", len(inputs), inputs)
	}

	if inputs["subject-digest"] != "sha256:fedcba9876543210" {
		t.Errorf("inputs[subject-digest] = %v, want %q", inputs["subject-digest"], "sha256:fedcba9876543210")
	}

	if inputs["subject-name"] != "container-image:v1.2.3" {
		t.Errorf("inputs[subject-name] = %v, want %q", inputs["subject-name"], "container-image:v1.2.3")
	}
}

func TestAttestBuildProvenance_Inputs_PushToRegistryWithToken(t *testing.T) {
	a := AttestBuildProvenance{
		SubjectDigest:  "sha256:abc123",
		SubjectName:    "ghcr.io/owner/image:latest",
		PushToRegistry: true,
		GithubToken:    "${{ secrets.GITHUB_TOKEN }}",
	}

	inputs := a.Inputs()

	if inputs["subject-digest"] != "sha256:abc123" {
		t.Errorf("inputs[subject-digest] = %v, want %q", inputs["subject-digest"], "sha256:abc123")
	}

	if inputs["subject-name"] != "ghcr.io/owner/image:latest" {
		t.Errorf("inputs[subject-name] = %v, want %q", inputs["subject-name"], "ghcr.io/owner/image:latest")
	}

	if inputs["push-to-registry"] != true {
		t.Errorf("inputs[push-to-registry] = %v, want true", inputs["push-to-registry"])
	}

	if inputs["github-token"] != "${{ secrets.GITHUB_TOKEN }}" {
		t.Errorf("inputs[github-token] = %v, want %q", inputs["github-token"], "${{ secrets.GITHUB_TOKEN }}")
	}
}
