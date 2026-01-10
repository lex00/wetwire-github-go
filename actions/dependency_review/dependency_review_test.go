package dependency_review

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestDependencyReview_Action(t *testing.T) {
	dr := DependencyReview{}
	if got := dr.Action(); got != "actions/dependency-review-action@v4" {
		t.Errorf("Action() = %q, want %q", got, "actions/dependency-review-action@v4")
	}
}

func TestDependencyReview_Inputs(t *testing.T) {
	dr := DependencyReview{
		FailOnSeverity: "high",
		FailOnScopes:   "runtime",
		AllowLicenses:  "MIT,Apache-2.0",
		DenyLicenses:   "GPL-3.0",
	}

	inputs := dr.Inputs()

	if inputs["fail-on-severity"] != "high" {
		t.Errorf("inputs[fail-on-severity] = %v, want %q", inputs["fail-on-severity"], "high")
	}

	if inputs["fail-on-scopes"] != "runtime" {
		t.Errorf("inputs[fail-on-scopes] = %v, want %q", inputs["fail-on-scopes"], "runtime")
	}

	if inputs["allow-licenses"] != "MIT,Apache-2.0" {
		t.Errorf("inputs[allow-licenses] = %v, want %q", inputs["allow-licenses"], "MIT,Apache-2.0")
	}

	if inputs["deny-licenses"] != "GPL-3.0" {
		t.Errorf("inputs[deny-licenses] = %v, want %q", inputs["deny-licenses"], "GPL-3.0")
	}
}

func TestDependencyReview_Inputs_Empty(t *testing.T) {
	dr := DependencyReview{}
	inputs := dr.Inputs()

	// Empty DependencyReview should have no inputs
	if len(inputs) != 0 {
		t.Errorf("empty DependencyReview.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestDependencyReview_Inputs_BoolFields(t *testing.T) {
	dr := DependencyReview{
		CommentSummaryInPR:      true,
		WarnOnly:                true,
		LicenseCheck:            true,
		VulnerabilityCheck:      true,
		RetryOnSnapshotWarnings: true,
	}

	inputs := dr.Inputs()

	if inputs["comment-summary-in-pr"] != true {
		t.Errorf("inputs[comment-summary-in-pr] = %v, want true", inputs["comment-summary-in-pr"])
	}

	if inputs["warn-only"] != true {
		t.Errorf("inputs[warn-only] = %v, want true", inputs["warn-only"])
	}

	if inputs["license-check"] != true {
		t.Errorf("inputs[license-check] = %v, want true", inputs["license-check"])
	}

	if inputs["vulnerability-check"] != true {
		t.Errorf("inputs[vulnerability-check] = %v, want true", inputs["vulnerability-check"])
	}

	if inputs["retry-on-snapshot-warnings"] != true {
		t.Errorf("inputs[retry-on-snapshot-warnings] = %v, want true", inputs["retry-on-snapshot-warnings"])
	}
}

func TestDependencyReview_Inputs_AllFields(t *testing.T) {
	dr := DependencyReview{
		FailOnSeverity:          "critical",
		FailOnScopes:            "runtime,development",
		AllowLicenses:           "MIT,Apache-2.0,BSD-3-Clause",
		DenyLicenses:            "GPL-3.0,AGPL-3.0",
		AllowGHSAs:              "GHSA-xxxx-yyyy-zzzz",
		ConfigFile:              ".github/dependency-review-config.yml",
		BaseRef:                 "main",
		HeadRef:                 "feature-branch",
		CommentSummaryInPR:      true,
		WarnOnly:                true,
		LicenseCheck:            true,
		VulnerabilityCheck:      true,
		RetryOnSnapshotWarnings: true,
	}

	inputs := dr.Inputs()

	if inputs["fail-on-severity"] != "critical" {
		t.Errorf("inputs[fail-on-severity] = %v, want %q", inputs["fail-on-severity"], "critical")
	}

	if inputs["fail-on-scopes"] != "runtime,development" {
		t.Errorf("inputs[fail-on-scopes] = %v, want %q", inputs["fail-on-scopes"], "runtime,development")
	}

	if inputs["allow-licenses"] != "MIT,Apache-2.0,BSD-3-Clause" {
		t.Errorf("inputs[allow-licenses] = %v, want %q", inputs["allow-licenses"], "MIT,Apache-2.0,BSD-3-Clause")
	}

	if inputs["deny-licenses"] != "GPL-3.0,AGPL-3.0" {
		t.Errorf("inputs[deny-licenses] = %v, want %q", inputs["deny-licenses"], "GPL-3.0,AGPL-3.0")
	}

	if inputs["allow-ghsas"] != "GHSA-xxxx-yyyy-zzzz" {
		t.Errorf("inputs[allow-ghsas] = %v, want %q", inputs["allow-ghsas"], "GHSA-xxxx-yyyy-zzzz")
	}

	if inputs["config-file"] != ".github/dependency-review-config.yml" {
		t.Errorf("inputs[config-file] = %v, want %q", inputs["config-file"], ".github/dependency-review-config.yml")
	}

	if inputs["base-ref"] != "main" {
		t.Errorf("inputs[base-ref] = %v, want %q", inputs["base-ref"], "main")
	}

	if inputs["head-ref"] != "feature-branch" {
		t.Errorf("inputs[head-ref] = %v, want %q", inputs["head-ref"], "feature-branch")
	}

	if inputs["comment-summary-in-pr"] != true {
		t.Errorf("inputs[comment-summary-in-pr] = %v, want true", inputs["comment-summary-in-pr"])
	}

	if inputs["warn-only"] != true {
		t.Errorf("inputs[warn-only] = %v, want true", inputs["warn-only"])
	}

	if inputs["license-check"] != true {
		t.Errorf("inputs[license-check] = %v, want true", inputs["license-check"])
	}

	if inputs["vulnerability-check"] != true {
		t.Errorf("inputs[vulnerability-check] = %v, want true", inputs["vulnerability-check"])
	}

	if inputs["retry-on-snapshot-warnings"] != true {
		t.Errorf("inputs[retry-on-snapshot-warnings] = %v, want true", inputs["retry-on-snapshot-warnings"])
	}
}

func TestDependencyReview_Inputs_SeverityLevels(t *testing.T) {
	tests := []struct {
		name     string
		severity string
	}{
		{"low", "low"},
		{"moderate", "moderate"},
		{"high", "high"},
		{"critical", "critical"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dr := DependencyReview{
				FailOnSeverity: tt.severity,
			}

			inputs := dr.Inputs()

			if inputs["fail-on-severity"] != tt.severity {
				t.Errorf("inputs[fail-on-severity] = %v, want %q", inputs["fail-on-severity"], tt.severity)
			}
		})
	}
}

func TestDependencyReview_Inputs_Scopes(t *testing.T) {
	tests := []struct {
		name   string
		scopes string
	}{
		{"development", "development"},
		{"runtime", "runtime"},
		{"unknown", "unknown"},
		{"multiple", "development,runtime"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dr := DependencyReview{
				FailOnScopes: tt.scopes,
			}

			inputs := dr.Inputs()

			if inputs["fail-on-scopes"] != tt.scopes {
				t.Errorf("inputs[fail-on-scopes] = %v, want %q", inputs["fail-on-scopes"], tt.scopes)
			}
		})
	}
}

func TestDependencyReview_Inputs_MinimalConfig(t *testing.T) {
	dr := DependencyReview{
		FailOnSeverity: "high",
	}

	inputs := dr.Inputs()

	if inputs["fail-on-severity"] != "high" {
		t.Errorf("inputs[fail-on-severity] = %v, want %q", inputs["fail-on-severity"], "high")
	}

	if len(inputs) != 1 {
		t.Errorf("minimal DependencyReview should have 1 input entry, got %d", len(inputs))
	}
}

func TestDependencyReview_Inputs_ConfigFile(t *testing.T) {
	dr := DependencyReview{
		ConfigFile: ".github/dependency-review.yml",
	}

	inputs := dr.Inputs()

	if inputs["config-file"] != ".github/dependency-review.yml" {
		t.Errorf("inputs[config-file] = %v, want %q", inputs["config-file"], ".github/dependency-review.yml")
	}
}

func TestDependencyReview_Inputs_RefComparison(t *testing.T) {
	dr := DependencyReview{
		BaseRef: "main",
		HeadRef: "feature-branch",
	}

	inputs := dr.Inputs()

	if inputs["base-ref"] != "main" {
		t.Errorf("inputs[base-ref] = %v, want %q", inputs["base-ref"], "main")
	}

	if inputs["head-ref"] != "feature-branch" {
		t.Errorf("inputs[head-ref] = %v, want %q", inputs["head-ref"], "feature-branch")
	}
}

func TestDependencyReview_Inputs_WarnOnlyMode(t *testing.T) {
	dr := DependencyReview{
		WarnOnly: true,
	}

	inputs := dr.Inputs()

	if inputs["warn-only"] != true {
		t.Errorf("inputs[warn-only] = %v, want true", inputs["warn-only"])
	}

	if len(inputs) != 1 {
		t.Errorf("warn-only DependencyReview should have 1 input entry, got %d", len(inputs))
	}
}

func TestDependencyReview_ImplementsStepAction(t *testing.T) {
	var _ workflow.StepAction = DependencyReview{}
}
