// Package codeql_analyze provides a typed wrapper for github/codeql-action/analyze.
package codeql_analyze

// CodeQLAnalyze wraps the github/codeql-action/analyze@v3 action.
// Analyze code with CodeQL and upload results to GitHub Security.
type CodeQLAnalyze struct {
	// Category to add to the SARIF file.
	Category string `yaml:"category,omitempty"`

	// Directory for SARIF file output.
	Output string `yaml:"output,omitempty"`

	// Whether to upload SARIF to GitHub.
	Upload bool `yaml:"upload,omitempty"`

	// Whether to upload CodeQL database.
	UploadDatabase bool `yaml:"upload-database,omitempty"`

	// Path to the repo checkout.
	CheckoutPath string `yaml:"checkout-path,omitempty"`

	// RAM limit in MB.
	RAM string `yaml:"ram,omitempty"`

	// Number of threads.
	Threads string `yaml:"threads,omitempty"`
}

// Action returns the action reference.
func (a CodeQLAnalyze) Action() string {
	return "github/codeql-action/analyze@v3"
}

// Inputs returns the action inputs as a map.
func (a CodeQLAnalyze) Inputs() map[string]any {
	with := make(map[string]any)

	if a.Category != "" {
		with["category"] = a.Category
	}
	if a.Output != "" {
		with["output"] = a.Output
	}
	if a.Upload {
		with["upload"] = a.Upload
	}
	if a.UploadDatabase {
		with["upload-database"] = a.UploadDatabase
	}
	if a.CheckoutPath != "" {
		with["checkout-path"] = a.CheckoutPath
	}
	if a.RAM != "" {
		with["ram"] = a.RAM
	}
	if a.Threads != "" {
		with["threads"] = a.Threads
	}

	return with
}
