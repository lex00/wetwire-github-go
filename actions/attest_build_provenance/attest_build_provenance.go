// Package attest_build_provenance provides a typed wrapper for actions/attest-build-provenance.
package attest_build_provenance

// AttestBuildProvenance wraps the actions/attest-build-provenance@v1 action.
// Generate signed build provenance attestations for workflow artifacts.
type AttestBuildProvenance struct {
	// Path to the artifact serving as the subject of the attestation.
	// May use wildcards and supports multiple file paths.
	SubjectPath string `yaml:"subject-path,omitempty"`

	// SHA256 digest of the subject for the attestation.
	// Must be in the form "sha256:hex_digest".
	SubjectDigest string `yaml:"subject-digest,omitempty"`

	// Subject name as it should appear in the attestation.
	// Required when using SubjectDigest.
	SubjectName string `yaml:"subject-name,omitempty"`

	// Path to a file containing checksums (digest and name) of subjects.
	SubjectChecksums string `yaml:"subject-checksums,omitempty"`

	// Whether to push the attestation to the image registry.
	// Requires subject-name to be a fully-qualified image name and subject-digest to be set.
	PushToRegistry bool `yaml:"push-to-registry,omitempty"`

	// Whether to create a storage record for the artifact.
	// Requires PushToRegistry to be true.
	CreateStorageRecord bool `yaml:"create-storage-record,omitempty"`

	// Whether to attach a list of generated attestations to the workflow run summary page.
	ShowSummary bool `yaml:"show-summary,omitempty"`

	// The GitHub token used to make authenticated API requests.
	GithubToken string `yaml:"github-token,omitempty"`
}

// Action returns the action reference.
func (a AttestBuildProvenance) Action() string {
	return "actions/attest-build-provenance@v1"
}

// Inputs returns the action inputs as a map.
func (a AttestBuildProvenance) Inputs() map[string]any {
	with := make(map[string]any)

	if a.SubjectPath != "" {
		with["subject-path"] = a.SubjectPath
	}
	if a.SubjectDigest != "" {
		with["subject-digest"] = a.SubjectDigest
	}
	if a.SubjectName != "" {
		with["subject-name"] = a.SubjectName
	}
	if a.SubjectChecksums != "" {
		with["subject-checksums"] = a.SubjectChecksums
	}
	if a.PushToRegistry {
		with["push-to-registry"] = a.PushToRegistry
	}
	if a.CreateStorageRecord {
		with["create-storage-record"] = a.CreateStorageRecord
	}
	if a.ShowSummary {
		with["show-summary"] = a.ShowSummary
	}
	if a.GithubToken != "" {
		with["github-token"] = a.GithubToken
	}

	return with
}
