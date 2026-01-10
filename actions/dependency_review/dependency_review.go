// Package dependency_review provides a typed wrapper for actions/dependency-review-action.
package dependency_review

// DependencyReview wraps the actions/dependency-review-action@v4 action.
// Review pull requests for dependency changes and identify security vulnerabilities.
type DependencyReview struct {
	// Severity level to fail on (low, moderate, high, critical)
	FailOnSeverity string `yaml:"fail-on-severity,omitempty"`

	// Scopes to fail on (development, runtime, unknown)
	FailOnScopes string `yaml:"fail-on-scopes,omitempty"`

	// Allowed licenses (comma-separated SPDX identifiers)
	AllowLicenses string `yaml:"allow-licenses,omitempty"`

	// Denied licenses (comma-separated SPDX identifiers)
	DenyLicenses string `yaml:"deny-licenses,omitempty"`

	// Allowed GitHub Security Advisory IDs (comma-separated)
	AllowGHSAs string `yaml:"allow-ghsas,omitempty"`

	// Path to configuration file
	ConfigFile string `yaml:"config-file,omitempty"`

	// Base ref for comparison
	BaseRef string `yaml:"base-ref,omitempty"`

	// Head ref for comparison
	HeadRef string `yaml:"head-ref,omitempty"`

	// Post summary as PR comment (default: false)
	CommentSummaryInPR bool `yaml:"comment-summary-in-pr,omitempty"`

	// Warn instead of fail (default: false)
	WarnOnly bool `yaml:"warn-only,omitempty"`

	// Enable license checking (default: true)
	LicenseCheck bool `yaml:"license-check,omitempty"`

	// Enable vulnerability checking (default: true)
	VulnerabilityCheck bool `yaml:"vulnerability-check,omitempty"`

	// Retry on snapshot warnings (default: false)
	RetryOnSnapshotWarnings bool `yaml:"retry-on-snapshot-warnings,omitempty"`
}

// Action returns the action reference.
func (a DependencyReview) Action() string {
	return "actions/dependency-review-action@v4"
}

// Inputs returns the action inputs as a map.
func (a DependencyReview) Inputs() map[string]any {
	with := make(map[string]any)

	if a.FailOnSeverity != "" {
		with["fail-on-severity"] = a.FailOnSeverity
	}
	if a.FailOnScopes != "" {
		with["fail-on-scopes"] = a.FailOnScopes
	}
	if a.AllowLicenses != "" {
		with["allow-licenses"] = a.AllowLicenses
	}
	if a.DenyLicenses != "" {
		with["deny-licenses"] = a.DenyLicenses
	}
	if a.AllowGHSAs != "" {
		with["allow-ghsas"] = a.AllowGHSAs
	}
	if a.ConfigFile != "" {
		with["config-file"] = a.ConfigFile
	}
	if a.BaseRef != "" {
		with["base-ref"] = a.BaseRef
	}
	if a.HeadRef != "" {
		with["head-ref"] = a.HeadRef
	}
	if a.CommentSummaryInPR {
		with["comment-summary-in-pr"] = a.CommentSummaryInPR
	}
	if a.WarnOnly {
		with["warn-only"] = a.WarnOnly
	}
	if a.LicenseCheck {
		with["license-check"] = a.LicenseCheck
	}
	if a.VulnerabilityCheck {
		with["vulnerability-check"] = a.VulnerabilityCheck
	}
	if a.RetryOnSnapshotWarnings {
		with["retry-on-snapshot-warnings"] = a.RetryOnSnapshotWarnings
	}

	return with
}
