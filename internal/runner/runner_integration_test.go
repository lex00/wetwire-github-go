package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lex00/wetwire-github-go/internal/discover"
)

// Test Runner struct fields can be set independently
func TestRunner_StructFields(t *testing.T) {
	r := &Runner{
		TempDir: "/custom/temp",
		GoPath:  "/custom/go",
		Verbose: true,
	}

	if r.TempDir != "/custom/temp" {
		t.Errorf("Runner.TempDir = %q, want %q", r.TempDir, "/custom/temp")
	}
	if r.GoPath != "/custom/go" {
		t.Errorf("Runner.GoPath = %q, want %q", r.GoPath, "/custom/go")
	}
	if !r.Verbose {
		t.Error("Runner.Verbose should be true")
	}
}

// Test ExtractionResult struct fields
func TestExtractionResult_Fields(t *testing.T) {
	result := ExtractionResult{
		Workflows: []ExtractedWorkflow{
			{Name: "CI", Data: map[string]any{"key": "value"}},
		},
		Jobs: []ExtractedJob{
			{Name: "Build", Data: map[string]any{"step": "compile"}},
		},
		Error: "test error",
	}

	if len(result.Workflows) != 1 {
		t.Errorf("len(Workflows) = %d, want 1", len(result.Workflows))
	}
	if len(result.Jobs) != 1 {
		t.Errorf("len(Jobs) = %d, want 1", len(result.Jobs))
	}
	if result.Error != "test error" {
		t.Errorf("Error = %q, want %q", result.Error, "test error")
	}
}

// Test ExtractedWorkflow struct
func TestExtractedWorkflow_Struct(t *testing.T) {
	w := ExtractedWorkflow{
		Name: "CI",
		Data: map[string]any{"on": "push"},
	}

	if w.Name != "CI" {
		t.Errorf("Name = %q, want %q", w.Name, "CI")
	}
	if w.Data["on"] != "push" {
		t.Errorf("Data[on] = %v, want %q", w.Data["on"], "push")
	}
}

// Test ExtractedJob struct
func TestExtractedJob_Struct(t *testing.T) {
	j := ExtractedJob{
		Name: "Build",
		Data: map[string]any{"runs-on": "ubuntu-latest"},
	}

	if j.Name != "Build" {
		t.Errorf("Name = %q, want %q", j.Name, "Build")
	}
	if j.Data["runs-on"] != "ubuntu-latest" {
		t.Errorf("Data[runs-on] = %v, want %q", j.Data["runs-on"], "ubuntu-latest")
	}
}

// Test DependabotExtractionResult struct
func TestDependabotExtractionResult_Struct(t *testing.T) {
	result := DependabotExtractionResult{
		Configs: []ExtractedDependabot{
			{Name: "Config1", Data: map[string]any{"version": 2}},
		},
		Error: "",
	}

	if len(result.Configs) != 1 {
		t.Errorf("len(Configs) = %d, want 1", len(result.Configs))
	}
}

// Test IssueTemplateExtractionResult struct
func TestIssueTemplateExtractionResult_Struct(t *testing.T) {
	result := IssueTemplateExtractionResult{
		Templates: []ExtractedIssueTemplate{
			{Name: "Bug", Data: map[string]any{"title": "Bug Report"}},
		},
		Error: "",
	}

	if len(result.Templates) != 1 {
		t.Errorf("len(Templates) = %d, want 1", len(result.Templates))
	}
}

// Test DiscussionTemplateExtractionResult struct
func TestDiscussionTemplateExtractionResult_Struct(t *testing.T) {
	result := DiscussionTemplateExtractionResult{
		Templates: []ExtractedDiscussionTemplate{
			{Name: "Announcement", Data: map[string]any{"title": "Announcement"}},
		},
		Error: "",
	}

	if len(result.Templates) != 1 {
		t.Errorf("len(Templates) = %d, want 1", len(result.Templates))
	}
}

// Test PRTemplateExtractionResult struct
func TestPRTemplateExtractionResult_Struct(t *testing.T) {
	result := PRTemplateExtractionResult{
		Templates: []ExtractedPRTemplate{
			{Name: "Default", Content: "## Description"},
		},
		Error: "",
	}

	if len(result.Templates) != 1 {
		t.Errorf("len(Templates) = %d, want 1", len(result.Templates))
	}
}

// Test CodeownersExtractionResult struct
func TestCodeownersExtractionResult_Struct(t *testing.T) {
	result := CodeownersExtractionResult{
		Configs: []ExtractedCodeowners{
			{
				Name: "Main",
				Rules: []ExtractedCodeownersRule{
					{Pattern: "*", Owners: []string{"@team"}, Comment: "Default"},
				},
			},
		},
		Error: "",
	}

	if len(result.Configs) != 1 {
		t.Errorf("len(Configs) = %d, want 1", len(result.Configs))
	}
	if len(result.Configs[0].Rules) != 1 {
		t.Errorf("len(Rules) = %d, want 1", len(result.Configs[0].Rules))
	}
}

// Test ExtractedCodeownersRule struct
func TestExtractedCodeownersRule_Struct(t *testing.T) {
	rule := ExtractedCodeownersRule{
		Pattern: "*.go",
		Owners:  []string{"@go-team", "@backend"},
		Comment: "Go files",
	}

	if rule.Pattern != "*.go" {
		t.Errorf("Pattern = %q, want %q", rule.Pattern, "*.go")
	}
	if len(rule.Owners) != 2 {
		t.Errorf("len(Owners) = %d, want 2", len(rule.Owners))
	}
	if rule.Comment != "Go files" {
		t.Errorf("Comment = %q, want %q", rule.Comment, "Go files")
	}
}

// Test ExtractedDependabot struct
func TestExtractedDependabot_Struct(t *testing.T) {
	d := ExtractedDependabot{
		Name: "DependabotConfig",
		Data: map[string]any{"version": 2},
	}

	if d.Name != "DependabotConfig" {
		t.Errorf("Name = %q, want %q", d.Name, "DependabotConfig")
	}
}

// Test ExtractedIssueTemplate struct
func TestExtractedIssueTemplate_Struct(t *testing.T) {
	tmpl := ExtractedIssueTemplate{
		Name: "BugReport",
		Data: map[string]any{"title": "Bug Report"},
	}

	if tmpl.Name != "BugReport" {
		t.Errorf("Name = %q, want %q", tmpl.Name, "BugReport")
	}
}

// Test ExtractedDiscussionTemplate struct
func TestExtractedDiscussionTemplate_Struct(t *testing.T) {
	tmpl := ExtractedDiscussionTemplate{
		Name: "Announcement",
		Data: map[string]any{"category": "announcements"},
	}

	if tmpl.Name != "Announcement" {
		t.Errorf("Name = %q, want %q", tmpl.Name, "Announcement")
	}
}

// Test ExtractedPRTemplate struct
func TestExtractedPRTemplate_Struct(t *testing.T) {
	tmpl := ExtractedPRTemplate{
		Name:    "DefaultPR",
		Content: "## Description\n\nPlease describe your changes.",
	}

	if tmpl.Name != "DefaultPR" {
		t.Errorf("Name = %q, want %q", tmpl.Name, "DefaultPR")
	}
	if !strings.Contains(tmpl.Content, "Description") {
		t.Error("Content should contain 'Description'")
	}
}

// Test ExtractedCodeowners struct
func TestExtractedCodeowners_Struct(t *testing.T) {
	co := ExtractedCodeowners{
		Name: "MainCodeowners",
		Rules: []ExtractedCodeownersRule{
			{Pattern: "*", Owners: []string{"@default"}},
		},
	}

	if co.Name != "MainCodeowners" {
		t.Errorf("Name = %q, want %q", co.Name, "MainCodeowners")
	}
	if len(co.Rules) != 1 {
		t.Errorf("len(Rules) = %d, want 1", len(co.Rules))
	}
}

// Test RunnerStruct fields
func TestRunnerStruct(t *testing.T) {
	r := &Runner{
		TempDir: "/custom/temp",
		GoPath:  "/custom/go",
		Verbose: true,
	}

	if r.TempDir != "/custom/temp" {
		t.Errorf("TempDir = %q, want %q", r.TempDir, "/custom/temp")
	}
	if r.GoPath != "/custom/go" {
		t.Errorf("GoPath = %q, want %q", r.GoPath, "/custom/go")
	}
	if !r.Verbose {
		t.Error("Verbose = false, want true")
	}
}

// Test ExtractionResult struct
func TestExtractionResultStruct(t *testing.T) {
	result := ExtractionResult{
		Workflows: []ExtractedWorkflow{
			{Name: "CI", Data: map[string]any{"key": "value"}},
		},
		Jobs: []ExtractedJob{
			{Name: "Build", Data: map[string]any{"name": "build"}},
		},
		Error: "some error",
	}

	if len(result.Workflows) != 1 {
		t.Errorf("len(Workflows) = %d, want 1", len(result.Workflows))
	}
	if len(result.Jobs) != 1 {
		t.Errorf("len(Jobs) = %d, want 1", len(result.Jobs))
	}
	if result.Error != "some error" {
		t.Errorf("Error = %q, want %q", result.Error, "some error")
	}
}

// Test ExtractedWorkflow struct
func TestExtractedWorkflowStruct(t *testing.T) {
	w := ExtractedWorkflow{
		Name: "CI",
		Data: map[string]any{"name": "Continuous Integration"},
	}

	if w.Name != "CI" {
		t.Errorf("Name = %q, want %q", w.Name, "CI")
	}
	if w.Data["name"] != "Continuous Integration" {
		t.Errorf("Data[name] = %v, want %q", w.Data["name"], "Continuous Integration")
	}
}

// Test ExtractedJob struct
func TestExtractedJobStruct(t *testing.T) {
	j := ExtractedJob{
		Name: "Build",
		Data: map[string]any{"runs-on": "ubuntu-latest"},
	}

	if j.Name != "Build" {
		t.Errorf("Name = %q, want %q", j.Name, "Build")
	}
	if j.Data["runs-on"] != "ubuntu-latest" {
		t.Errorf("Data[runs-on] = %v, want %q", j.Data["runs-on"], "ubuntu-latest")
	}
}

// Integration test - ExtractValues with real Go code
func TestRunner_ExtractValues_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Get the absolute path to this project's root
	// We're in internal/runner, so we need to go up two levels
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Find the project root (where go.mod is)
	projectRoot := wd
	for {
		if _, err := os.Stat(filepath.Join(projectRoot, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(projectRoot)
		if parent == projectRoot {
			t.Skip("Could not find project root")
		}
		projectRoot = parent
	}

	// Create a test project that imports from this project
	tmpDir := t.TempDir()

	// Create go.mod that references the main project
	goMod := fmt.Sprintf(`module testproject

go 1.23

require github.com/lex00/wetwire-github-go v0.0.0

replace github.com/lex00/wetwire-github-go => %s
`, projectRoot)

	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a simple workflow file
	workflowCode := `package testproject

import "github.com/lex00/wetwire-github-go/workflow"

var TestWorkflow = workflow.Workflow{
	Name: "Test CI",
}

var TestJob = workflow.Job{
	Name:   "test",
	RunsOn: "ubuntu-latest",
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "workflows.go"), []byte(workflowCode), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()

	discovered := &discover.DiscoveryResult{
		Workflows: []discover.DiscoveredWorkflow{
			{Name: "TestWorkflow", File: filepath.Join(tmpDir, "workflows.go"), Line: 5},
		},
		Jobs: []discover.DiscoveredJob{
			{Name: "TestJob", File: filepath.Join(tmpDir, "workflows.go"), Line: 10},
		},
	}

	result, err := r.ExtractValues(tmpDir, discovered)
	if err != nil {
		t.Fatalf("ExtractValues() error = %v", err)
	}

	if len(result.Workflows) != 1 {
		t.Errorf("len(Workflows) = %d, want 1", len(result.Workflows))
	}
	if len(result.Jobs) != 1 {
		t.Errorf("len(Jobs) = %d, want 1", len(result.Jobs))
	}

	if result.Workflows[0].Name != "TestWorkflow" {
		t.Errorf("Workflow name = %q, want %q", result.Workflows[0].Name, "TestWorkflow")
	}
	if result.Jobs[0].Name != "TestJob" {
		t.Errorf("Job name = %q, want %q", result.Jobs[0].Name, "TestJob")
	}
}

// Integration test - ExtractDependabot with real Go code
func TestRunner_ExtractDependabot_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	projectRoot := wd
	for {
		if _, err := os.Stat(filepath.Join(projectRoot, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(projectRoot)
		if parent == projectRoot {
			t.Skip("Could not find project root")
		}
		projectRoot = parent
	}

	tmpDir := t.TempDir()

	goMod := fmt.Sprintf(`module testproject

go 1.23

require github.com/lex00/wetwire-github-go v0.0.0

replace github.com/lex00/wetwire-github-go => %s
`, projectRoot)

	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	dependabotCode := `package testproject

import "github.com/lex00/wetwire-github-go/dependabot"

var TestDependabot = dependabot.Dependabot{
	Version: 2,
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "dependabot.go"), []byte(dependabotCode), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()

	discovered := &discover.DependabotDiscoveryResult{
		Configs: []discover.DiscoveredDependabot{
			{Name: "TestDependabot", File: filepath.Join(tmpDir, "dependabot.go"), Line: 5},
		},
	}

	result, err := r.ExtractDependabot(tmpDir, discovered)
	if err != nil {
		t.Fatalf("ExtractDependabot() error = %v", err)
	}

	if len(result.Configs) != 1 {
		t.Errorf("len(Configs) = %d, want 1", len(result.Configs))
	}

	if result.Configs[0].Name != "TestDependabot" {
		t.Errorf("Config name = %q, want %q", result.Configs[0].Name, "TestDependabot")
	}
}

// Integration test - ExtractIssueTemplates with real Go code
func TestRunner_ExtractIssueTemplates_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	projectRoot := wd
	for {
		if _, err := os.Stat(filepath.Join(projectRoot, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(projectRoot)
		if parent == projectRoot {
			t.Skip("Could not find project root")
		}
		projectRoot = parent
	}

	tmpDir := t.TempDir()

	goMod := fmt.Sprintf(`module testproject

go 1.23

require github.com/lex00/wetwire-github-go v0.0.0

replace github.com/lex00/wetwire-github-go => %s
`, projectRoot)

	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	templateCode := `package testproject

import "github.com/lex00/wetwire-github-go/templates"

var TestIssueTemplate = templates.IssueTemplate{
	Name:        "Bug Report",
	Description: "Report a bug",
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "issue_templates.go"), []byte(templateCode), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()

	discovered := &discover.IssueTemplateDiscoveryResult{
		Templates: []discover.DiscoveredIssueTemplate{
			{Name: "TestIssueTemplate", File: filepath.Join(tmpDir, "issue_templates.go"), Line: 5},
		},
	}

	result, err := r.ExtractIssueTemplates(tmpDir, discovered)
	if err != nil {
		t.Fatalf("ExtractIssueTemplates() error = %v", err)
	}

	if len(result.Templates) != 1 {
		t.Errorf("len(Templates) = %d, want 1", len(result.Templates))
	}

	if result.Templates[0].Name != "TestIssueTemplate" {
		t.Errorf("Template name = %q, want %q", result.Templates[0].Name, "TestIssueTemplate")
	}
}

// Integration test - ExtractDiscussionTemplates with real Go code
func TestRunner_ExtractDiscussionTemplates_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	projectRoot := wd
	for {
		if _, err := os.Stat(filepath.Join(projectRoot, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(projectRoot)
		if parent == projectRoot {
			t.Skip("Could not find project root")
		}
		projectRoot = parent
	}

	tmpDir := t.TempDir()

	goMod := fmt.Sprintf(`module testproject

go 1.23

require github.com/lex00/wetwire-github-go v0.0.0

replace github.com/lex00/wetwire-github-go => %s
`, projectRoot)

	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	discussionTemplateCode := `package testproject

import "github.com/lex00/wetwire-github-go/templates"

var TestDiscussionTemplate = templates.DiscussionTemplate{
	Title:       "Announcement",
	Description: "Make an announcement",
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "discussion_templates.go"), []byte(discussionTemplateCode), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()

	discovered := &discover.DiscussionTemplateDiscoveryResult{
		Templates: []discover.DiscoveredDiscussionTemplate{
			{Name: "TestDiscussionTemplate", File: filepath.Join(tmpDir, "discussion_templates.go"), Line: 5},
		},
	}

	result, err := r.ExtractDiscussionTemplates(tmpDir, discovered)
	if err != nil {
		t.Fatalf("ExtractDiscussionTemplates() error = %v", err)
	}

	if len(result.Templates) != 1 {
		t.Errorf("len(Templates) = %d, want 1", len(result.Templates))
	}

	if result.Templates[0].Name != "TestDiscussionTemplate" {
		t.Errorf("Template name = %q, want %q", result.Templates[0].Name, "TestDiscussionTemplate")
	}
}

// Integration test - ExtractPRTemplates with real Go code
func TestRunner_ExtractPRTemplates_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	projectRoot := wd
	for {
		if _, err := os.Stat(filepath.Join(projectRoot, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(projectRoot)
		if parent == projectRoot {
			t.Skip("Could not find project root")
		}
		projectRoot = parent
	}

	tmpDir := t.TempDir()

	goMod := fmt.Sprintf(`module testproject

go 1.23

require github.com/lex00/wetwire-github-go v0.0.0

replace github.com/lex00/wetwire-github-go => %s
`, projectRoot)

	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	prTemplateCode := "package testproject\n\nimport \"github.com/lex00/wetwire-github-go/templates\"\n\nvar TestPRTemplate = templates.PRTemplate{\n\tContent: \"Description\",\n}\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "pr_templates.go"), []byte(prTemplateCode), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()

	discovered := &discover.PRTemplateDiscoveryResult{
		Templates: []discover.DiscoveredPRTemplate{
			{Name: "TestPRTemplate", File: filepath.Join(tmpDir, "pr_templates.go"), Line: 5},
		},
	}

	result, err := r.ExtractPRTemplates(tmpDir, discovered)
	if err != nil {
		t.Fatalf("ExtractPRTemplates() error = %v", err)
	}

	if len(result.Templates) != 1 {
		t.Errorf("len(Templates) = %d, want 1", len(result.Templates))
	}

	if result.Templates[0].Name != "TestPRTemplate" {
		t.Errorf("Template name = %q, want %q", result.Templates[0].Name, "TestPRTemplate")
	}
}

// Integration test - ExtractCodeowners with real Go code
func TestRunner_ExtractCodeowners_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	projectRoot := wd
	for {
		if _, err := os.Stat(filepath.Join(projectRoot, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(projectRoot)
		if parent == projectRoot {
			t.Skip("Could not find project root")
		}
		projectRoot = parent
	}

	tmpDir := t.TempDir()

	goMod := fmt.Sprintf(`module testproject

go 1.23

require github.com/lex00/wetwire-github-go v0.0.0

replace github.com/lex00/wetwire-github-go => %s
`, projectRoot)

	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	codeownersCode := "package testproject\n\nimport \"github.com/lex00/wetwire-github-go/codeowners\"\n\nvar TestCodeowners = codeowners.Owners{\n\tRules: []codeowners.Rule{\n\t\t{Pattern: \"*\", Owners: []string{\"@default-team\"}},\n\t\t{Pattern: \"*.go\", Owners: []string{\"@go-team\"}},\n\t},\n}\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "codeowners.go"), []byte(codeownersCode), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()

	discovered := &discover.CodeownersDiscoveryResult{
		Configs: []discover.DiscoveredCodeowners{
			{Name: "TestCodeowners", File: filepath.Join(tmpDir, "codeowners.go"), Line: 5},
		},
	}

	result, err := r.ExtractCodeowners(tmpDir, discovered)
	if err != nil {
		t.Fatalf("ExtractCodeowners() error = %v", err)
	}

	if len(result.Configs) != 1 {
		t.Errorf("len(Configs) = %d, want 1", len(result.Configs))
	}

	if result.Configs[0].Name != "TestCodeowners" {
		t.Errorf("Config name = %q, want %q", result.Configs[0].Name, "TestCodeowners")
	}

	if len(result.Configs[0].Rules) != 2 {
		t.Errorf("len(Rules) = %d, want 2", len(result.Configs[0].Rules))
	}
}
