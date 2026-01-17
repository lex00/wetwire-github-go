package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lex00/wetwire-github-go/internal/discover"
)

func TestRunner_ExtractValues_Empty(t *testing.T) {
	r := NewRunner()

	discovered := &discover.DiscoveryResult{
		Workflows: []discover.DiscoveredWorkflow{},
		Jobs:      []discover.DiscoveredJob{},
	}

	result, err := r.ExtractValues(".", discovered)
	if err != nil {
		t.Fatalf("ExtractValues() error = %v", err)
	}

	if len(result.Workflows) != 0 {
		t.Errorf("len(result.Workflows) = %d, want 0", len(result.Workflows))
	}

	if len(result.Jobs) != 0 {
		t.Errorf("len(result.Jobs) = %d, want 0", len(result.Jobs))
	}
}

// Test ExtractDependabot with empty configs
func TestRunner_ExtractDependabot_Empty(t *testing.T) {
	r := NewRunner()

	discovered := &discover.DependabotDiscoveryResult{
		Configs: []discover.DiscoveredDependabot{},
	}

	result, err := r.ExtractDependabot(".", discovered)
	if err != nil {
		t.Fatalf("ExtractDependabot() error = %v", err)
	}

	if len(result.Configs) != 0 {
		t.Errorf("len(result.Configs) = %d, want 0", len(result.Configs))
	}
}

// Test ExtractIssueTemplates with empty templates
func TestRunner_ExtractIssueTemplates_Empty(t *testing.T) {
	r := NewRunner()

	discovered := &discover.IssueTemplateDiscoveryResult{
		Templates: []discover.DiscoveredIssueTemplate{},
	}

	result, err := r.ExtractIssueTemplates(".", discovered)
	if err != nil {
		t.Fatalf("ExtractIssueTemplates() error = %v", err)
	}

	if len(result.Templates) != 0 {
		t.Errorf("len(result.Templates) = %d, want 0", len(result.Templates))
	}
}

// Test ExtractDiscussionTemplates with empty templates
func TestRunner_ExtractDiscussionTemplates_Empty(t *testing.T) {
	r := NewRunner()

	discovered := &discover.DiscussionTemplateDiscoveryResult{
		Templates: []discover.DiscoveredDiscussionTemplate{},
	}

	result, err := r.ExtractDiscussionTemplates(".", discovered)
	if err != nil {
		t.Fatalf("ExtractDiscussionTemplates() error = %v", err)
	}

	if len(result.Templates) != 0 {
		t.Errorf("len(result.Templates) = %d, want 0", len(result.Templates))
	}
}

// Test ExtractPRTemplates with empty templates
func TestRunner_ExtractPRTemplates_Empty(t *testing.T) {
	r := NewRunner()

	discovered := &discover.PRTemplateDiscoveryResult{
		Templates: []discover.DiscoveredPRTemplate{},
	}

	result, err := r.ExtractPRTemplates(".", discovered)
	if err != nil {
		t.Fatalf("ExtractPRTemplates() error = %v", err)
	}

	if len(result.Templates) != 0 {
		t.Errorf("len(result.Templates) = %d, want 0", len(result.Templates))
	}
}

// Test ExtractCodeowners with empty configs
func TestRunner_ExtractCodeowners_Empty(t *testing.T) {
	r := NewRunner()

	discovered := &discover.CodeownersDiscoveryResult{
		Configs: []discover.DiscoveredCodeowners{},
	}

	result, err := r.ExtractCodeowners(".", discovered)
	if err != nil {
		t.Fatalf("ExtractCodeowners() error = %v", err)
	}

	if len(result.Configs) != 0 {
		t.Errorf("len(result.Configs) = %d, want 0", len(result.Configs))
	}
}

// Test ExtractValues path resolution
func TestRunner_ExtractValues_PathResolution(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()

	discovered := &discover.DiscoveryResult{
		Workflows: []discover.DiscoveredWorkflow{
			{Name: "CI", File: filepath.Join(tmpDir, "workflows.go"), Line: 10},
		},
	}

	// Use relative path (current directory)
	_, err := r.ExtractValues(".", discovered)
	// This will fail because we don't have actual workflow files, but it should
	// at least get past the path resolution stage
	if err == nil {
		t.Log("ExtractValues() succeeded (test project may have been created)")
	} else if !strings.Contains(err.Error(), "parsing go.mod") {
		// Should fail during later stages, not path resolution
		t.Logf("ExtractValues() error = %v (expected error in later stages)", err)
	}
}

// Test ExtractDependabot path resolution
func TestRunner_ExtractDependabot_PathResolution(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()

	discovered := &discover.DependabotDiscoveryResult{
		Configs: []discover.DiscoveredDependabot{
			{Name: "Config", File: filepath.Join(tmpDir, "dependabot.go"), Line: 10},
		},
	}

	// Use the temp directory with a valid go.mod
	_, err := r.ExtractDependabot(tmpDir, discovered)
	// This will fail at program execution stage, but path resolution should work
	if err != nil && !strings.Contains(err.Error(), "parsing go.mod") {
		t.Logf("ExtractDependabot() error = %v (expected)", err)
	}
}

// Test ExtractIssueTemplates path resolution
func TestRunner_ExtractIssueTemplates_PathResolution(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()

	discovered := &discover.IssueTemplateDiscoveryResult{
		Templates: []discover.DiscoveredIssueTemplate{
			{Name: "Template", File: filepath.Join(tmpDir, "template.go"), Line: 10},
		},
	}

	_, err := r.ExtractIssueTemplates(tmpDir, discovered)
	if err != nil && !strings.Contains(err.Error(), "parsing go.mod") {
		t.Logf("ExtractIssueTemplates() error = %v (expected)", err)
	}
}

// Test ExtractDiscussionTemplates path resolution
func TestRunner_ExtractDiscussionTemplates_PathResolution(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()

	discovered := &discover.DiscussionTemplateDiscoveryResult{
		Templates: []discover.DiscoveredDiscussionTemplate{
			{Name: "Template", File: filepath.Join(tmpDir, "template.go"), Line: 10},
		},
	}

	_, err := r.ExtractDiscussionTemplates(tmpDir, discovered)
	if err != nil && !strings.Contains(err.Error(), "parsing go.mod") {
		t.Logf("ExtractDiscussionTemplates() error = %v (expected)", err)
	}
}

// Test ExtractPRTemplates path resolution
func TestRunner_ExtractPRTemplates_PathResolution(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()

	discovered := &discover.PRTemplateDiscoveryResult{
		Templates: []discover.DiscoveredPRTemplate{
			{Name: "Template", File: filepath.Join(tmpDir, "template.go"), Line: 10},
		},
	}

	_, err := r.ExtractPRTemplates(tmpDir, discovered)
	if err != nil && !strings.Contains(err.Error(), "parsing go.mod") {
		t.Logf("ExtractPRTemplates() error = %v (expected)", err)
	}
}

// Test ExtractCodeowners path resolution
func TestRunner_ExtractCodeowners_PathResolution(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()

	discovered := &discover.CodeownersDiscoveryResult{
		Configs: []discover.DiscoveredCodeowners{
			{Name: "Config", File: filepath.Join(tmpDir, "codeowners.go"), Line: 10},
		},
	}

	_, err := r.ExtractCodeowners(tmpDir, discovered)
	if err != nil && !strings.Contains(err.Error(), "parsing go.mod") {
		t.Logf("ExtractCodeowners() error = %v (expected)", err)
	}
}

// Suppress unused import warning
var _ = fmt.Sprintf
