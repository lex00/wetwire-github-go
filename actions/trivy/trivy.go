// Package trivy provides a typed wrapper for aquasecurity/trivy-action.
package trivy

// Trivy wraps the aquasecurity/trivy-action@0.28.0 action.
// Scan container images, file systems, repositories for vulnerabilities.
type Trivy struct {
	// Container image to scan (e.g., alpine:3.18)
	ImageRef string `yaml:"image-ref,omitempty"`

	// Type of scan (fs, image, repo, config)
	ScanType string `yaml:"scan-type,omitempty"`

	// Output format (table, json, sarif, cyclonedx, spdx, github)
	Format string `yaml:"format,omitempty"`

	// Severities to report (UNKNOWN,LOW,MEDIUM,HIGH,CRITICAL)
	Severity string `yaml:"severity,omitempty"`

	// Exit code when issues found
	ExitCode int `yaml:"exit-code,omitempty"`

	// Ignore unfixed vulnerabilities
	IgnoreUnfixed bool `yaml:"ignore-unfixed,omitempty"`

	// Vulnerability types to scan (os, library)
	VulnType string `yaml:"vuln-type,omitempty"`

	// Which scanners to run (vuln, config, secret, license)
	Scanners string `yaml:"scanners,omitempty"`

	// Output template
	Template string `yaml:"template,omitempty"`

	// Output file path
	Output string `yaml:"output,omitempty"`
}

// Action returns the action reference.
func (a Trivy) Action() string {
	return "aquasecurity/trivy-action@0.28.0"
}

// Inputs returns the action inputs as a map.
func (a Trivy) Inputs() map[string]any {
	with := make(map[string]any)

	if a.ImageRef != "" {
		with["image-ref"] = a.ImageRef
	}
	if a.ScanType != "" {
		with["scan-type"] = a.ScanType
	}
	if a.Format != "" {
		with["format"] = a.Format
	}
	if a.Severity != "" {
		with["severity"] = a.Severity
	}
	if a.ExitCode != 0 {
		with["exit-code"] = a.ExitCode
	}
	if a.IgnoreUnfixed {
		with["ignore-unfixed"] = a.IgnoreUnfixed
	}
	if a.VulnType != "" {
		with["vuln-type"] = a.VulnType
	}
	if a.Scanners != "" {
		with["scanners"] = a.Scanners
	}
	if a.Template != "" {
		with["template"] = a.Template
	}
	if a.Output != "" {
		with["output"] = a.Output
	}

	return with
}
