package wetwire_test

import (
	"os"
	"strings"
	"testing"
)

func TestReadmeBadges(t *testing.T) {
	content, err := os.ReadFile("README.md")
	if err != nil {
		t.Fatalf("failed to read README.md: %v", err)
	}

	readme := string(content)

	// Required badges per WETWIRE_SPEC.md Section 12.4
	requiredBadges := []struct {
		name    string
		pattern string
	}{
		{"CI badge", "github.com/lex00/wetwire-github-go/actions/workflows"},
		{"Go Reference badge", "pkg.go.dev/badge/github.com/lex00/wetwire-github-go"},
		{"Go Report Card badge", "goreportcard.com/badge/github.com/lex00/wetwire-github-go"},
		{"License badge", "License-MIT"},
	}

	for _, badge := range requiredBadges {
		if !strings.Contains(readme, badge.pattern) {
			t.Errorf("README.md missing required badge: %s (expected pattern: %q)", badge.name, badge.pattern)
		}
	}

	// Verify badge order: CI, Go Reference, Go Report, License
	ciPos := strings.Index(readme, "github.com/lex00/wetwire-github-go/actions/workflows")
	goRefPos := strings.Index(readme, "pkg.go.dev/badge/github.com/lex00/wetwire-github-go")
	goReportPos := strings.Index(readme, "goreportcard.com/badge/github.com/lex00/wetwire-github-go")
	licensePos := strings.Index(readme, "License-MIT")

	if ciPos == -1 || goRefPos == -1 || goReportPos == -1 || licensePos == -1 {
		// Skip order check if badges are missing (already reported above)
		return
	}

	if ciPos > goRefPos {
		t.Error("CI badge should appear before Go Reference badge")
	}
	if goRefPos > goReportPos {
		t.Error("Go Reference badge should appear before Go Report Card badge")
	}
	if goReportPos > licensePos {
		t.Error("Go Report Card badge should appear before License badge")
	}
}
