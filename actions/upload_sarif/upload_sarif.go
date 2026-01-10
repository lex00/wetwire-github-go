// Package upload_sarif provides a typed wrapper for github/codeql-action/upload-sarif.
package upload_sarif

// UploadSarif wraps the github/codeql-action/upload-sarif@v3 action.
// Upload a SARIF file to GitHub Security to display alerts.
type UploadSarif struct {
	// Path to the SARIF file to upload.
	SarifFile string `yaml:"sarif_file,omitempty"`

	// Path to the checkout location.
	CheckoutPath string `yaml:"checkout_path,omitempty"`

	// The ref where results were found.
	Ref string `yaml:"ref,omitempty"`

	// The sha of the commit where results were found.
	Sha string `yaml:"sha,omitempty"`

	// Category for the results.
	Category string `yaml:"category,omitempty"`

	// GitHub token for API access.
	Token string `yaml:"token,omitempty"`

	// Wait for processing to complete.
	WaitForProcessing bool `yaml:"wait-for-processing,omitempty"`
}

// Action returns the action reference.
func (a UploadSarif) Action() string {
	return "github/codeql-action/upload-sarif@v3"
}

// Inputs returns the action inputs as a map.
func (a UploadSarif) Inputs() map[string]any {
	with := make(map[string]any)

	if a.SarifFile != "" {
		with["sarif_file"] = a.SarifFile
	}
	if a.CheckoutPath != "" {
		with["checkout_path"] = a.CheckoutPath
	}
	if a.Ref != "" {
		with["ref"] = a.Ref
	}
	if a.Sha != "" {
		with["sha"] = a.Sha
	}
	if a.Category != "" {
		with["category"] = a.Category
	}
	if a.Token != "" {
		with["token"] = a.Token
	}
	if a.WaitForProcessing {
		with["wait-for-processing"] = a.WaitForProcessing
	}

	return with
}
