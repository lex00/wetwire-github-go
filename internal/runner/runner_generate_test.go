package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lex00/wetwire-github-go/internal/discover"
)

func TestRunner_generateProgram(t *testing.T) {
	r := NewRunner()
	baseDir := "/project"

	discovered := &discover.DiscoveryResult{
		Workflows: []discover.DiscoveredWorkflow{
			{Name: "CI", File: "/project/workflows.go", Line: 10},
		},
		Jobs: []discover.DiscoveredJob{
			{Name: "Build", File: "/project/jobs.go", Line: 5},
		},
	}

	program, err := r.generateProgram("github.com/example/test", baseDir, discovered)
	if err != nil {
		t.Fatalf("generateProgram() error = %v", err)
	}

	expectedStrings := []string{
		"package main",
		"encoding/json",
		"ExtractionResult",
		"ExtractedWorkflow",
		"ExtractedJob",
		"toMap",
		"json.Marshal",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(program, expected) {
			t.Errorf("generateProgram() missing %q\n\nGenerated:\n%s", expected, program)
		}
	}
}

// Test generateDependabotProgram
func TestRunner_generateDependabotProgram(t *testing.T) {
	r := NewRunner()
	baseDir := "/project"

	discovered := &discover.DependabotDiscoveryResult{
		Configs: []discover.DiscoveredDependabot{
			{Name: "DependabotConfig", File: "/project/dependabot.go", Line: 10},
		},
	}

	program, err := r.generateDependabotProgram("github.com/example/test", baseDir, discovered)
	if err != nil {
		t.Fatalf("generateDependabotProgram() error = %v", err)
	}

	expectedStrings := []string{
		"package main",
		"encoding/json",
		"DependabotExtractionResult",
		"ExtractedDependabot",
		"toMap",
		"json.Marshal",
		"DependabotConfig",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(program, expected) {
			t.Errorf("generateDependabotProgram() missing %q\n\nGenerated:\n%s", expected, program)
		}
	}
}

// Test generateIssueTemplateProgram
func TestRunner_generateIssueTemplateProgram(t *testing.T) {
	r := NewRunner()
	baseDir := "/project"

	discovered := &discover.IssueTemplateDiscoveryResult{
		Templates: []discover.DiscoveredIssueTemplate{
			{Name: "BugReport", File: "/project/issue_templates.go", Line: 10},
		},
	}

	program, err := r.generateIssueTemplateProgram("github.com/example/test", baseDir, discovered)
	if err != nil {
		t.Fatalf("generateIssueTemplateProgram() error = %v", err)
	}

	expectedStrings := []string{
		"package main",
		"encoding/json",
		"IssueTemplateExtractionResult",
		"ExtractedIssueTemplate",
		"toMap",
		"json.Marshal",
		"BugReport",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(program, expected) {
			t.Errorf("generateIssueTemplateProgram() missing %q\n\nGenerated:\n%s", expected, program)
		}
	}
}

// Test generateDiscussionTemplateProgram
func TestRunner_generateDiscussionTemplateProgram(t *testing.T) {
	r := NewRunner()
	baseDir := "/project"

	discovered := &discover.DiscussionTemplateDiscoveryResult{
		Templates: []discover.DiscoveredDiscussionTemplate{
			{Name: "Announcement", File: "/project/discussion_templates.go", Line: 10},
		},
	}

	program, err := r.generateDiscussionTemplateProgram("github.com/example/test", baseDir, discovered)
	if err != nil {
		t.Fatalf("generateDiscussionTemplateProgram() error = %v", err)
	}

	expectedStrings := []string{
		"package main",
		"encoding/json",
		"DiscussionTemplateExtractionResult",
		"ExtractedDiscussionTemplate",
		"toMap",
		"json.Marshal",
		"Announcement",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(program, expected) {
			t.Errorf("generateDiscussionTemplateProgram() missing %q\n\nGenerated:\n%s", expected, program)
		}
	}
}

// Test generatePRTemplateProgram
func TestRunner_generatePRTemplateProgram(t *testing.T) {
	r := NewRunner()
	baseDir := "/project"

	discovered := &discover.PRTemplateDiscoveryResult{
		Templates: []discover.DiscoveredPRTemplate{
			{Name: "DefaultPR", File: "/project/pr_templates.go", Line: 10},
		},
	}

	program, err := r.generatePRTemplateProgram("github.com/example/test", baseDir, discovered)
	if err != nil {
		t.Fatalf("generatePRTemplateProgram() error = %v", err)
	}

	expectedStrings := []string{
		"package main",
		"encoding/json",
		"PRTemplateExtractionResult",
		"ExtractedPRTemplate",
		"DefaultPR",
		"Content",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(program, expected) {
			t.Errorf("generatePRTemplateProgram() missing %q\n\nGenerated:\n%s", expected, program)
		}
	}
}

// Test generateCodeownersProgram
func TestRunner_generateCodeownersProgram(t *testing.T) {
	r := NewRunner()
	baseDir := "/project"

	discovered := &discover.CodeownersDiscoveryResult{
		Configs: []discover.DiscoveredCodeowners{
			{Name: "DefaultCodeowners", File: "/project/codeowners.go", Line: 10},
		},
	}

	program, err := r.generateCodeownersProgram("github.com/example/test", baseDir, discovered)
	if err != nil {
		t.Fatalf("generateCodeownersProgram() error = %v", err)
	}

	expectedStrings := []string{
		"package main",
		"encoding/json",
		"CodeownersExtractionResult",
		"ExtractedCodeowners",
		"DefaultCodeowners",
		"extractConfig",
		"Rules",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(program, expected) {
			t.Errorf("generateCodeownersProgram() missing %q\n\nGenerated:\n%s", expected, program)
		}
	}
}

// Test generateProgram with multiple packages
func TestRunner_generateProgram_MultiplePackages(t *testing.T) {
	r := NewRunner()
	baseDir := "/project"

	discovered := &discover.DiscoveryResult{
		Workflows: []discover.DiscoveredWorkflow{
			{Name: "CI", File: "/project/workflows.go", Line: 10},
			{Name: "Deploy", File: "/project/internal/deploy/workflows.go", Line: 5},
		},
		Jobs: []discover.DiscoveredJob{
			{Name: "Build", File: "/project/jobs.go", Line: 5},
			{Name: "Test", File: "/project/pkg/test/jobs.go", Line: 3},
		},
	}

	program, err := r.generateProgram("github.com/example/test", baseDir, discovered)
	if err != nil {
		t.Fatalf("generateProgram() error = %v", err)
	}

	// Should import multiple packages
	if !strings.Contains(program, "github.com/example/test") {
		t.Error("generateProgram() missing root package import")
	}

	if !strings.Contains(program, "github.com/example/test/internal/deploy") {
		t.Error("generateProgram() missing internal/deploy package import")
	}

	if !strings.Contains(program, "github.com/example/test/pkg/test") {
		t.Error("generateProgram() missing pkg/test package import")
	}

	// Should reference all workflows and jobs
	for _, expected := range []string{"CI", "Deploy", "Build", "Test"} {
		if !strings.Contains(program, expected) {
			t.Errorf("generateProgram() missing reference to %q", expected)
		}
	}
}

// Test generateDependabotProgram with multiple configs from different packages
func TestRunner_generateDependabotProgram_MultiplePackages(t *testing.T) {
	r := NewRunner()
	baseDir := "/project"

	discovered := &discover.DependabotDiscoveryResult{
		Configs: []discover.DiscoveredDependabot{
			{Name: "Config1", File: "/project/dependabot.go", Line: 10},
			{Name: "Config2", File: "/project/internal/ci/dependabot.go", Line: 5},
		},
	}

	program, err := r.generateDependabotProgram("github.com/example/test", baseDir, discovered)
	if err != nil {
		t.Fatalf("generateDependabotProgram() error = %v", err)
	}

	if !strings.Contains(program, "Config1") {
		t.Error("generateDependabotProgram() missing Config1")
	}

	if !strings.Contains(program, "Config2") {
		t.Error("generateDependabotProgram() missing Config2")
	}

	if !strings.Contains(program, "github.com/example/test/internal/ci") {
		t.Error("generateDependabotProgram() missing internal/ci package import")
	}
}

// Test generateIssueTemplateProgram with multiple templates
func TestRunner_generateIssueTemplateProgram_MultipleTemplates(t *testing.T) {
	r := NewRunner()
	baseDir := "/project"

	discovered := &discover.IssueTemplateDiscoveryResult{
		Templates: []discover.DiscoveredIssueTemplate{
			{Name: "BugReport", File: "/project/templates.go", Line: 10},
			{Name: "FeatureRequest", File: "/project/templates.go", Line: 20},
		},
	}

	program, err := r.generateIssueTemplateProgram("github.com/example/test", baseDir, discovered)
	if err != nil {
		t.Fatalf("generateIssueTemplateProgram() error = %v", err)
	}

	if !strings.Contains(program, "BugReport") {
		t.Error("generateIssueTemplateProgram() missing BugReport")
	}

	if !strings.Contains(program, "FeatureRequest") {
		t.Error("generateIssueTemplateProgram() missing FeatureRequest")
	}
}

// Test generatePRTemplateProgram with multiple templates from different packages
func TestRunner_generatePRTemplateProgram_MultiplePackages(t *testing.T) {
	r := NewRunner()
	baseDir := "/project"

	discovered := &discover.PRTemplateDiscoveryResult{
		Templates: []discover.DiscoveredPRTemplate{
			{Name: "DefaultPR", File: "/project/pr.go", Line: 10},
			{Name: "HotfixPR", File: "/project/pkg/templates/pr.go", Line: 5},
		},
	}

	program, err := r.generatePRTemplateProgram("github.com/example/test", baseDir, discovered)
	if err != nil {
		t.Fatalf("generatePRTemplateProgram() error = %v", err)
	}

	if !strings.Contains(program, "DefaultPR") {
		t.Error("generatePRTemplateProgram() missing DefaultPR")
	}

	if !strings.Contains(program, "HotfixPR") {
		t.Error("generatePRTemplateProgram() missing HotfixPR")
	}

	if !strings.Contains(program, "github.com/example/test/pkg/templates") {
		t.Error("generatePRTemplateProgram() missing pkg/templates package import")
	}
}

// Test generateCodeownersProgram with multiple configs
func TestRunner_generateCodeownersProgram_MultipleConfigs(t *testing.T) {
	r := NewRunner()
	baseDir := "/project"

	discovered := &discover.CodeownersDiscoveryResult{
		Configs: []discover.DiscoveredCodeowners{
			{Name: "MainCodeowners", File: "/project/codeowners.go", Line: 10},
			{Name: "SubCodeowners", File: "/project/internal/codeowners.go", Line: 5},
		},
	}

	program, err := r.generateCodeownersProgram("github.com/example/test", baseDir, discovered)
	if err != nil {
		t.Fatalf("generateCodeownersProgram() error = %v", err)
	}

	if !strings.Contains(program, "MainCodeowners") {
		t.Error("generateCodeownersProgram() missing MainCodeowners")
	}

	if !strings.Contains(program, "SubCodeowners") {
		t.Error("generateCodeownersProgram() missing SubCodeowners")
	}

	if !strings.Contains(program, "github.com/lex00/wetwire-github-go/codeowners") {
		t.Error("generateCodeownersProgram() missing codeowners package import")
	}
}

// Test generateProgram with same package for workflows and jobs
func TestRunner_generateProgram_SamePackage(t *testing.T) {
	r := NewRunner()
	baseDir := "/project"

	discovered := &discover.DiscoveryResult{
		Workflows: []discover.DiscoveredWorkflow{
			{Name: "CI", File: "/project/workflows.go", Line: 10},
		},
		Jobs: []discover.DiscoveredJob{
			{Name: "Build", File: "/project/workflows.go", Line: 20},
		},
	}

	program, err := r.generateProgram("github.com/example/test", baseDir, discovered)
	if err != nil {
		t.Fatalf("generateProgram() error = %v", err)
	}

	// Should only import the package once
	importCount := strings.Count(program, `"github.com/example/test"`)
	if importCount != 1 {
		t.Errorf("generateProgram() imported package %d times, want 1", importCount)
	}
}

// Test generateDependabotProgram with same package for multiple configs
func TestRunner_generateDependabotProgram_SamePackage(t *testing.T) {
	r := NewRunner()
	baseDir := "/project"

	discovered := &discover.DependabotDiscoveryResult{
		Configs: []discover.DiscoveredDependabot{
			{Name: "Config1", File: "/project/dependabot.go", Line: 10},
			{Name: "Config2", File: "/project/dependabot.go", Line: 20},
		},
	}

	program, err := r.generateDependabotProgram("github.com/example/test", baseDir, discovered)
	if err != nil {
		t.Fatalf("generateDependabotProgram() error = %v", err)
	}

	// Should only import the package once
	importCount := strings.Count(program, `"github.com/example/test"`)
	if importCount != 1 {
		t.Errorf("generateDependabotProgram() imported package %d times, want 1", importCount)
	}

	// Should reference both configs
	if !strings.Contains(program, "Config1") || !strings.Contains(program, "Config2") {
		t.Error("generateDependabotProgram() should reference both configs")
	}
}

// Test generateIssueTemplateProgram with same package for multiple templates
func TestRunner_generateIssueTemplateProgram_SamePackage(t *testing.T) {
	r := NewRunner()
	baseDir := "/project"

	discovered := &discover.IssueTemplateDiscoveryResult{
		Templates: []discover.DiscoveredIssueTemplate{
			{Name: "Bug", File: "/project/templates.go", Line: 10},
			{Name: "Feature", File: "/project/templates.go", Line: 20},
		},
	}

	program, err := r.generateIssueTemplateProgram("github.com/example/test", baseDir, discovered)
	if err != nil {
		t.Fatalf("generateIssueTemplateProgram() error = %v", err)
	}

	// Should only import the package once
	importCount := strings.Count(program, `"github.com/example/test"`)
	if importCount != 1 {
		t.Errorf("generateIssueTemplateProgram() imported package %d times, want 1", importCount)
	}
}

// Test generateDiscussionTemplateProgram with same package
func TestRunner_generateDiscussionTemplateProgram_SamePackage(t *testing.T) {
	r := NewRunner()
	baseDir := "/project"

	discovered := &discover.DiscussionTemplateDiscoveryResult{
		Templates: []discover.DiscoveredDiscussionTemplate{
			{Name: "Announcement", File: "/project/templates.go", Line: 10},
			{Name: "Question", File: "/project/templates.go", Line: 20},
		},
	}

	program, err := r.generateDiscussionTemplateProgram("github.com/example/test", baseDir, discovered)
	if err != nil {
		t.Fatalf("generateDiscussionTemplateProgram() error = %v", err)
	}

	// Should only import the package once
	importCount := strings.Count(program, `"github.com/example/test"`)
	if importCount != 1 {
		t.Errorf("generateDiscussionTemplateProgram() imported package %d times, want 1", importCount)
	}
}

// Test generatePRTemplateProgram with same package
func TestRunner_generatePRTemplateProgram_SamePackage(t *testing.T) {
	r := NewRunner()
	baseDir := "/project"

	discovered := &discover.PRTemplateDiscoveryResult{
		Templates: []discover.DiscoveredPRTemplate{
			{Name: "Default", File: "/project/templates.go", Line: 10},
			{Name: "Hotfix", File: "/project/templates.go", Line: 20},
		},
	}

	program, err := r.generatePRTemplateProgram("github.com/example/test", baseDir, discovered)
	if err != nil {
		t.Fatalf("generatePRTemplateProgram() error = %v", err)
	}

	// Should only import the package once
	importCount := strings.Count(program, `"github.com/example/test"`)
	if importCount != 1 {
		t.Errorf("generatePRTemplateProgram() imported package %d times, want 1", importCount)
	}
}

// Test generateCodeownersProgram with same package
func TestRunner_generateCodeownersProgram_SamePackage(t *testing.T) {
	r := NewRunner()
	baseDir := "/project"

	discovered := &discover.CodeownersDiscoveryResult{
		Configs: []discover.DiscoveredCodeowners{
			{Name: "Main", File: "/project/codeowners.go", Line: 10},
			{Name: "Secondary", File: "/project/codeowners.go", Line: 20},
		},
	}

	program, err := r.generateCodeownersProgram("github.com/example/test", baseDir, discovered)
	if err != nil {
		t.Fatalf("generateCodeownersProgram() error = %v", err)
	}

	// Should only import the package once
	importCount := strings.Count(program, `"github.com/example/test"`)
	if importCount != 1 {
		t.Errorf("generateCodeownersProgram() imported package %d times, want 1", importCount)
	}
}

// Test generateProgram with deeply nested packages
func TestRunner_generateProgram_DeeplyNested(t *testing.T) {
	r := NewRunner()
	baseDir := "/project"

	discovered := &discover.DiscoveryResult{
		Workflows: []discover.DiscoveredWorkflow{
			{Name: "CI", File: "/project/a/b/c/d/workflows.go", Line: 10},
		},
	}

	program, err := r.generateProgram("github.com/example/test", baseDir, discovered)
	if err != nil {
		t.Fatalf("generateProgram() error = %v", err)
	}

	if !strings.Contains(program, "github.com/example/test/a/b/c/d") {
		t.Error("generateProgram() should handle deeply nested packages")
	}
}

// Test generateDiscussionTemplateProgram with multiple packages
func TestRunner_generateDiscussionTemplateProgram_MultiplePackages(t *testing.T) {
	r := NewRunner()
	baseDir := "/project"

	discovered := &discover.DiscussionTemplateDiscoveryResult{
		Templates: []discover.DiscoveredDiscussionTemplate{
			{Name: "Announcement", File: "/project/templates.go", Line: 10},
			{Name: "Question", File: "/project/internal/discussion/templates.go", Line: 5},
		},
	}

	program, err := r.generateDiscussionTemplateProgram("github.com/example/test", baseDir, discovered)
	if err != nil {
		t.Fatalf("generateDiscussionTemplateProgram() error = %v", err)
	}

	if !strings.Contains(program, "Announcement") {
		t.Error("generateDiscussionTemplateProgram() missing Announcement")
	}

	if !strings.Contains(program, "Question") {
		t.Error("generateDiscussionTemplateProgram() missing Question")
	}

	if !strings.Contains(program, "github.com/example/test/internal/discussion") {
		t.Error("generateDiscussionTemplateProgram() missing internal/discussion package import")
	}
}

// Test generateProgram import ordering is consistent
func TestRunner_generateProgram_ImportConsistency(t *testing.T) {
	r := NewRunner()
	baseDir := "/project"

	discovered := &discover.DiscoveryResult{
		Workflows: []discover.DiscoveredWorkflow{
			{Name: "CI", File: "/project/workflows.go", Line: 10},
		},
		Jobs: []discover.DiscoveredJob{
			{Name: "Build", File: "/project/jobs.go", Line: 5},
		},
	}

	// Generate the same program multiple times
	program1, _ := r.generateProgram("github.com/example/test", baseDir, discovered)
	program2, _ := r.generateProgram("github.com/example/test", baseDir, discovered)

	// Programs should be identical
	if program1 != program2 {
		t.Error("generateProgram() should produce consistent output")
	}
}

// Test generateProgram with empty workflows and jobs
func TestRunner_generateProgram_Empty(t *testing.T) {
	r := NewRunner()
	baseDir := "/project"

	discovered := &discover.DiscoveryResult{
		Workflows: []discover.DiscoveredWorkflow{},
		Jobs:      []discover.DiscoveredJob{},
	}

	program, err := r.generateProgram("github.com/example/test", baseDir, discovered)
	if err != nil {
		t.Fatalf("generateProgram() error = %v", err)
	}

	// Should still generate valid Go code
	if !strings.Contains(program, "package main") {
		t.Error("generateProgram() should contain package main")
	}
}

// Test generateProgram with only workflows
func TestRunner_generateProgram_OnlyWorkflows(t *testing.T) {
	r := NewRunner()
	baseDir := "/project"

	discovered := &discover.DiscoveryResult{
		Workflows: []discover.DiscoveredWorkflow{
			{Name: "CI", File: "/project/workflows.go", Line: 10},
			{Name: "Deploy", File: "/project/workflows.go", Line: 20},
		},
		Jobs: []discover.DiscoveredJob{},
	}

	program, err := r.generateProgram("github.com/example/test", baseDir, discovered)
	if err != nil {
		t.Fatalf("generateProgram() error = %v", err)
	}

	if !strings.Contains(program, "CI") || !strings.Contains(program, "Deploy") {
		t.Error("generateProgram() should contain workflow names")
	}
}

// Test generateProgram with only jobs
func TestRunner_generateProgram_OnlyJobs(t *testing.T) {
	r := NewRunner()
	baseDir := "/project"

	discovered := &discover.DiscoveryResult{
		Workflows: []discover.DiscoveredWorkflow{},
		Jobs: []discover.DiscoveredJob{
			{Name: "Build", File: "/project/jobs.go", Line: 5},
			{Name: "Test", File: "/project/jobs.go", Line: 15},
		},
	}

	program, err := r.generateProgram("github.com/example/test", baseDir, discovered)
	if err != nil {
		t.Fatalf("generateProgram() error = %v", err)
	}

	if !strings.Contains(program, "Build") || !strings.Contains(program, "Test") {
		t.Error("generateProgram() should contain job names")
	}
}

// Test generateDependabotProgram empty slices
func TestRunner_generateDependabotProgram_Empty(t *testing.T) {
	r := NewRunner()
	baseDir := "/project"

	discovered := &discover.DependabotDiscoveryResult{
		Configs: []discover.DiscoveredDependabot{},
	}

	program, err := r.generateDependabotProgram("github.com/example/test", baseDir, discovered)
	if err != nil {
		t.Fatalf("generateDependabotProgram() error = %v", err)
	}

	if !strings.Contains(program, "package main") {
		t.Error("generateDependabotProgram() should contain package main")
	}
}

// Test generateIssueTemplateProgram empty slices
func TestRunner_generateIssueTemplateProgram_Empty(t *testing.T) {
	r := NewRunner()
	baseDir := "/project"

	discovered := &discover.IssueTemplateDiscoveryResult{
		Templates: []discover.DiscoveredIssueTemplate{},
	}

	program, err := r.generateIssueTemplateProgram("github.com/example/test", baseDir, discovered)
	if err != nil {
		t.Fatalf("generateIssueTemplateProgram() error = %v", err)
	}

	if !strings.Contains(program, "package main") {
		t.Error("generateIssueTemplateProgram() should contain package main")
	}
}

// Test generateDiscussionTemplateProgram empty slices
func TestRunner_generateDiscussionTemplateProgram_Empty(t *testing.T) {
	r := NewRunner()
	baseDir := "/project"

	discovered := &discover.DiscussionTemplateDiscoveryResult{
		Templates: []discover.DiscoveredDiscussionTemplate{},
	}

	program, err := r.generateDiscussionTemplateProgram("github.com/example/test", baseDir, discovered)
	if err != nil {
		t.Fatalf("generateDiscussionTemplateProgram() error = %v", err)
	}

	if !strings.Contains(program, "package main") {
		t.Error("generateDiscussionTemplateProgram() should contain package main")
	}
}

// Test generatePRTemplateProgram empty slices
func TestRunner_generatePRTemplateProgram_Empty(t *testing.T) {
	r := NewRunner()
	baseDir := "/project"

	discovered := &discover.PRTemplateDiscoveryResult{
		Templates: []discover.DiscoveredPRTemplate{},
	}

	program, err := r.generatePRTemplateProgram("github.com/example/test", baseDir, discovered)
	if err != nil {
		t.Fatalf("generatePRTemplateProgram() error = %v", err)
	}

	if !strings.Contains(program, "package main") {
		t.Error("generatePRTemplateProgram() should contain package main")
	}
}

// Test generateCodeownersProgram empty slices
func TestRunner_generateCodeownersProgram_Empty(t *testing.T) {
	r := NewRunner()
	baseDir := "/project"

	discovered := &discover.CodeownersDiscoveryResult{
		Configs: []discover.DiscoveredCodeowners{},
	}

	program, err := r.generateCodeownersProgram("github.com/example/test", baseDir, discovered)
	if err != nil {
		t.Fatalf("generateCodeownersProgram() error = %v", err)
	}

	if !strings.Contains(program, "package main") {
		t.Error("generateCodeownersProgram() should contain package main")
	}
}

// Test all generate functions handle nil discovery gracefully
func TestRunner_generateProgram_NilSlices(t *testing.T) {
	r := NewRunner()
	baseDir := "/project"

	// Use actual nil slices
	discovered := &discover.DiscoveryResult{}

	program, err := r.generateProgram("github.com/example/test", baseDir, discovered)
	if err != nil {
		t.Fatalf("generateProgram() error = %v", err)
	}

	if !strings.Contains(program, "package main") {
		t.Error("generateProgram() should generate valid Go code even with nil slices")
	}
}

// Test generateDependabotProgram with nil configs
func TestRunner_generateDependabotProgram_NilSlice(t *testing.T) {
	r := NewRunner()

	discovered := &discover.DependabotDiscoveryResult{}

	program, err := r.generateDependabotProgram("github.com/example/test", "/project", discovered)
	if err != nil {
		t.Fatalf("generateDependabotProgram() error = %v", err)
	}

	if !strings.Contains(program, "package main") {
		t.Error("generateDependabotProgram() should generate valid Go code")
	}
}

// Test generateIssueTemplateProgram with nil templates
func TestRunner_generateIssueTemplateProgram_NilSlice(t *testing.T) {
	r := NewRunner()

	discovered := &discover.IssueTemplateDiscoveryResult{}

	program, err := r.generateIssueTemplateProgram("github.com/example/test", "/project", discovered)
	if err != nil {
		t.Fatalf("generateIssueTemplateProgram() error = %v", err)
	}

	if !strings.Contains(program, "package main") {
		t.Error("generateIssueTemplateProgram() should generate valid Go code")
	}
}

// Test generateDiscussionTemplateProgram with nil templates
func TestRunner_generateDiscussionTemplateProgram_NilSlice(t *testing.T) {
	r := NewRunner()

	discovered := &discover.DiscussionTemplateDiscoveryResult{}

	program, err := r.generateDiscussionTemplateProgram("github.com/example/test", "/project", discovered)
	if err != nil {
		t.Fatalf("generateDiscussionTemplateProgram() error = %v", err)
	}

	if !strings.Contains(program, "package main") {
		t.Error("generateDiscussionTemplateProgram() should generate valid Go code")
	}
}

// Test generatePRTemplateProgram with nil templates
func TestRunner_generatePRTemplateProgram_NilSlice(t *testing.T) {
	r := NewRunner()

	discovered := &discover.PRTemplateDiscoveryResult{}

	program, err := r.generatePRTemplateProgram("github.com/example/test", "/project", discovered)
	if err != nil {
		t.Fatalf("generatePRTemplateProgram() error = %v", err)
	}

	if !strings.Contains(program, "package main") {
		t.Error("generatePRTemplateProgram() should generate valid Go code")
	}
}

// Test generateCodeownersProgram with nil configs
func TestRunner_generateCodeownersProgram_NilSlice(t *testing.T) {
	r := NewRunner()

	discovered := &discover.CodeownersDiscoveryResult{}

	program, err := r.generateCodeownersProgram("github.com/example/test", "/project", discovered)
	if err != nil {
		t.Fatalf("generateCodeownersProgram() error = %v", err)
	}

	if !strings.Contains(program, "package main") {
		t.Error("generateCodeownersProgram() should generate valid Go code")
	}
}

// Suppress unused import warning
var _ = fmt.Sprintf
var _ = os.TempDir
var _ = filepath.Join
