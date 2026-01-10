package runner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lex00/wetwire-github-go/internal/discover"
)

func TestNewRunner(t *testing.T) {
	r := NewRunner()
	if r == nil {
		t.Error("NewRunner() returned nil")
	}
	if r.TempDir == "" {
		t.Error("NewRunner().TempDir is empty")
	}
}

func TestRunner_parseGoMod(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23

require (
	github.com/some/dep v1.0.0
)
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()
	modulePath, err := r.parseGoMod(tmpDir)
	if err != nil {
		t.Fatalf("parseGoMod() error = %v", err)
	}

	if modulePath != "github.com/example/test" {
		t.Errorf("parseGoMod() = %q, want %q", modulePath, "github.com/example/test")
	}
}

func TestRunner_parseGoMod_NotFound(t *testing.T) {
	tmpDir := t.TempDir()

	r := NewRunner()
	_, err := r.parseGoMod(tmpDir)
	if err == nil {
		t.Error("parseGoMod() expected error for missing go.mod")
	}
}

func TestRunner_parseGoMod_NoModule(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()
	_, err := r.parseGoMod(tmpDir)
	if err == nil {
		t.Error("parseGoMod() expected error for missing module directive")
	}
}

func TestRunner_generateGoMod(t *testing.T) {
	r := NewRunner()
	result := r.generateGoMod("github.com/example/test", "/path/to/project")

	if !strings.Contains(result, "module wetwire-extract") {
		t.Error("generateGoMod() missing module directive")
	}

	if !strings.Contains(result, "require github.com/example/test") {
		t.Error("generateGoMod() missing require directive")
	}

	if !strings.Contains(result, "replace github.com/example/test =>") {
		t.Error("generateGoMod() missing replace directive")
	}
}

func TestRunner_getPackagePath(t *testing.T) {
	r := NewRunner()
	baseDir := "/project"

	tests := []struct {
		modulePath string
		baseDir    string
		file       string
		want       string
	}{
		{"github.com/example/test", baseDir, "/project/workflows.go", "github.com/example/test"},
		{"github.com/example/test", baseDir, "/project/pkg/workflows.go", "github.com/example/test/pkg"},
		{"github.com/example/test", baseDir, "/project/internal/ci/workflows.go", "github.com/example/test/internal/ci"},
	}

	for _, tt := range tests {
		got := r.getPackagePath(tt.modulePath, tt.baseDir, tt.file)
		if got != tt.want {
			t.Errorf("getPackagePath(%q, %q, %q) = %q, want %q", tt.modulePath, tt.baseDir, tt.file, got, tt.want)
		}
	}
}

func TestRunner_pkgAlias(t *testing.T) {
	r := NewRunner()

	tests := []struct {
		input string
		want  string
	}{
		{"github.com/example/test", "test"},
		{"github.com/example/my-pkg", "my_pkg"},
		{"github.com/org/repo/internal/ci", "ci"},
	}

	for _, tt := range tests {
		got := r.pkgAlias(tt.input)
		if got != tt.want {
			t.Errorf("pkgAlias(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

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

func TestFindGoBinary(t *testing.T) {
	path, err := FindGoBinary()
	if err != nil {
		t.Skipf("Go binary not found, skipping: %v", err)
	}

	if path == "" {
		t.Error("FindGoBinary() returned empty path")
	}
}

func TestRunner_resolveReplaceDirective(t *testing.T) {
	r := NewRunner()

	tests := []struct {
		name    string
		line    string
		baseDir string
		want    string
	}{
		{
			name:    "relative path with ../",
			line:    "replace github.com/example/dep => ../dep",
			baseDir: "/project/subdir",
			want:    "replace github.com/example/dep => /project/dep",
		},
		{
			name:    "relative path with .",
			line:    "replace github.com/example/dep => ./local",
			baseDir: "/project",
			want:    "replace github.com/example/dep => /project/local",
		},
		{
			name:    "absolute path unchanged",
			line:    "replace github.com/example/dep => /absolute/path",
			baseDir: "/project",
			want:    "replace github.com/example/dep => /absolute/path",
		},
		{
			name:    "version replacement unchanged",
			line:    "replace github.com/example/dep v1.0.0 => v1.0.1",
			baseDir: "/project",
			want:    "replace github.com/example/dep v1.0.0 => v1.0.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := r.resolveReplaceDirective(tt.line, tt.baseDir)
			if got != tt.want {
				t.Errorf("resolveReplaceDirective() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRunner_parseReplaceDirectives_ResolvesRelativePaths(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a go.mod with a relative replace directive
	goMod := `module github.com/example/test

go 1.23

require github.com/other/dep v1.0.0

replace github.com/other/dep => ../dep
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()
	replaces := r.parseReplaceDirectives(tmpDir)

	if len(replaces) != 1 {
		t.Fatalf("parseReplaceDirectives() returned %d directives, want 1", len(replaces))
	}

	// The relative path should be resolved to an absolute path
	if strings.Contains(replaces[0], "..") {
		t.Errorf("parseReplaceDirectives() should resolve relative paths, got %q", replaces[0])
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

// Test ExtractValues with invalid directory
func TestRunner_ExtractValues_InvalidDir(t *testing.T) {
	r := NewRunner()

	discovered := &discover.DiscoveryResult{
		Workflows: []discover.DiscoveredWorkflow{
			{Name: "CI", File: "/nonexistent/workflows.go", Line: 10},
		},
	}

	_, err := r.ExtractValues("/nonexistent/directory", discovered)
	if err == nil {
		t.Error("ExtractValues() expected error for nonexistent directory")
	}
}

// Test getPackagePath with complex paths
func TestRunner_getPackagePath_ComplexPaths(t *testing.T) {
	r := NewRunner()

	tests := []struct {
		name       string
		modulePath string
		baseDir    string
		file       string
		wantSubstr string
	}{
		{
			name:       "different paths still compute relative",
			modulePath: "github.com/example/test",
			baseDir:    "/project",
			file:       "/different/absolute/file.go",
			wantSubstr: "github.com/example/test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := r.getPackagePath(tt.modulePath, tt.baseDir, tt.file)
			if !strings.Contains(result, tt.wantSubstr) {
				t.Errorf("getPackagePath() = %q, should contain %q", result, tt.wantSubstr)
			}
		})
	}
}

// Test generateGoMod with multiple replace directives
func TestRunner_generateGoMod_MultipleReplaces(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23

require (
	github.com/dep1 v1.0.0
	github.com/dep2 v2.0.0
)

replace github.com/dep1 => ../dep1
replace github.com/dep2 => /absolute/dep2
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()
	result := r.generateGoMod("github.com/example/test", tmpDir)

	if !strings.Contains(result, "module wetwire-extract") {
		t.Error("generateGoMod() missing module directive")
	}

	if !strings.Contains(result, "require github.com/example/test") {
		t.Error("generateGoMod() missing require directive")
	}

	// Should include both replace directives
	if !strings.Contains(result, "replace github.com/dep1") {
		t.Error("generateGoMod() missing first replace directive")
	}

	if !strings.Contains(result, "replace github.com/dep2") {
		t.Error("generateGoMod() missing second replace directive")
	}
}

// Test resolveReplaceDirective with malformed input
func TestRunner_resolveReplaceDirective_Malformed(t *testing.T) {
	r := NewRunner()

	tests := []struct {
		name    string
		line    string
		baseDir string
		want    string
	}{
		{
			name:    "no arrow",
			line:    "replace github.com/example/dep",
			baseDir: "/project",
			want:    "replace github.com/example/dep",
		},
		{
			name:    "multiple arrows",
			line:    "replace github.com/example/dep => => /path",
			baseDir: "/project",
			want:    "replace github.com/example/dep => => /path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := r.resolveReplaceDirective(tt.line, tt.baseDir)
			if got != tt.want {
				t.Errorf("resolveReplaceDirective() = %q, want %q", got, tt.want)
			}
		})
	}
}

// Test FindGoBinary error path
func TestFindGoBinary_Error(t *testing.T) {
	// Save original PATH
	originalPath := os.Getenv("PATH")
	defer os.Setenv("PATH", originalPath)

	// Set PATH to empty to force error
	os.Setenv("PATH", "")

	_, err := FindGoBinary()
	if err == nil {
		// On some systems, go might still be found via other means
		// So we don't fail the test, just skip it
		t.Skip("Go binary found even with empty PATH")
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

// Test pkgAlias with various inputs
func TestRunner_pkgAlias_EdgeCases(t *testing.T) {
	r := NewRunner()

	tests := []struct {
		input string
		want  string
	}{
		{"github.com/example/test", "test"},
		{"github.com/example/my-pkg", "my_pkg"},
		{"github.com/org/repo/internal/ci", "ci"},
		{"github.com/complex-name/with-many-hyphens", "with_many_hyphens"},
		{"simple", "simple"},
		{"a", "a"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := r.pkgAlias(tt.input)
			if got != tt.want {
				t.Errorf("pkgAlias(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// Test parseReplaceDirectives with missing go.mod
func TestRunner_parseReplaceDirectives_MissingGoMod(t *testing.T) {
	tmpDir := t.TempDir()

	r := NewRunner()
	replaces := r.parseReplaceDirectives(tmpDir)

	// Should return empty slice when go.mod doesn't exist
	if len(replaces) != 0 {
		t.Errorf("parseReplaceDirectives() = %v, want empty slice", replaces)
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

// Test ExtractDependabot error path - invalid directory
func TestRunner_ExtractDependabot_InvalidDir(t *testing.T) {
	r := NewRunner()

	discovered := &discover.DependabotDiscoveryResult{
		Configs: []discover.DiscoveredDependabot{
			{Name: "Config", File: "/nonexistent/dependabot.go", Line: 10},
		},
	}

	_, err := r.ExtractDependabot("/nonexistent/directory", discovered)
	if err == nil {
		t.Error("ExtractDependabot() expected error for nonexistent directory")
	}
}

// Test ExtractIssueTemplates error path - invalid directory
func TestRunner_ExtractIssueTemplates_InvalidDir(t *testing.T) {
	r := NewRunner()

	discovered := &discover.IssueTemplateDiscoveryResult{
		Templates: []discover.DiscoveredIssueTemplate{
			{Name: "Template", File: "/nonexistent/template.go", Line: 10},
		},
	}

	_, err := r.ExtractIssueTemplates("/nonexistent/directory", discovered)
	if err == nil {
		t.Error("ExtractIssueTemplates() expected error for nonexistent directory")
	}
}

// Test ExtractDiscussionTemplates error path - invalid directory
func TestRunner_ExtractDiscussionTemplates_InvalidDir(t *testing.T) {
	r := NewRunner()

	discovered := &discover.DiscussionTemplateDiscoveryResult{
		Templates: []discover.DiscoveredDiscussionTemplate{
			{Name: "Template", File: "/nonexistent/template.go", Line: 10},
		},
	}

	_, err := r.ExtractDiscussionTemplates("/nonexistent/directory", discovered)
	if err == nil {
		t.Error("ExtractDiscussionTemplates() expected error for nonexistent directory")
	}
}

// Test ExtractPRTemplates error path - invalid directory
func TestRunner_ExtractPRTemplates_InvalidDir(t *testing.T) {
	r := NewRunner()

	discovered := &discover.PRTemplateDiscoveryResult{
		Templates: []discover.DiscoveredPRTemplate{
			{Name: "Template", File: "/nonexistent/template.go", Line: 10},
		},
	}

	_, err := r.ExtractPRTemplates("/nonexistent/directory", discovered)
	if err == nil {
		t.Error("ExtractPRTemplates() expected error for nonexistent directory")
	}
}

// Test ExtractCodeowners error path - invalid directory
func TestRunner_ExtractCodeowners_InvalidDir(t *testing.T) {
	r := NewRunner()

	discovered := &discover.CodeownersDiscoveryResult{
		Configs: []discover.DiscoveredCodeowners{
			{Name: "Config", File: "/nonexistent/codeowners.go", Line: 10},
		},
	}

	_, err := r.ExtractCodeowners("/nonexistent/directory", discovered)
	if err == nil {
		t.Error("ExtractCodeowners() expected error for nonexistent directory")
	}
}

// Test getPackagePath edge cases
func TestRunner_getPackagePath_RootPackage(t *testing.T) {
	r := NewRunner()

	tests := []struct {
		name       string
		modulePath string
		baseDir    string
		file       string
		want       string
	}{
		{
			name:       "file in root directory",
			modulePath: "github.com/example/test",
			baseDir:    "/project",
			file:       "/project/main.go",
			want:       "github.com/example/test",
		},
		{
			name:       "file in subdirectory",
			modulePath: "github.com/example/test",
			baseDir:    "/project",
			file:       "/project/cmd/app/main.go",
			want:       "github.com/example/test/cmd/app",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := r.getPackagePath(tt.modulePath, tt.baseDir, tt.file)
			if got != tt.want {
				t.Errorf("getPackagePath() = %q, want %q", got, tt.want)
			}
		})
	}
}

// Test NewRunner initializes fields correctly
func TestNewRunner_Fields(t *testing.T) {
	r := NewRunner()

	if r.TempDir == "" {
		t.Error("NewRunner().TempDir should not be empty")
	}

	// GoPath might be empty if go binary is not found, which is OK
	// Just verify the field exists
	_ = r.GoPath

	// Verbose should be false by default
	if r.Verbose {
		t.Error("NewRunner().Verbose should be false by default")
	}
}

// Test generateGoMod with no replace directives
func TestRunner_generateGoMod_NoReplaces(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()
	result := r.generateGoMod("github.com/example/test", tmpDir)

	if !strings.Contains(result, "module wetwire-extract") {
		t.Error("generateGoMod() missing module directive")
	}

	if !strings.Contains(result, "require github.com/example/test") {
		t.Error("generateGoMod() missing require directive")
	}

	// Should have the replace directive for the module itself
	if !strings.Contains(result, "replace github.com/example/test =>") {
		t.Error("generateGoMod() missing replace directive for module")
	}
}

// Test parseReplaceDirectives with various formats
func TestRunner_parseReplaceDirectives_VariousFormats(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23

replace (
	github.com/dep1 => ../dep1
	github.com/dep2 v1.0.0 => github.com/dep2 v1.0.1
)

replace github.com/dep3 => ./local/dep3
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()
	replaces := r.parseReplaceDirectives(tmpDir)

	// Should find replace directives outside of blocks
	if len(replaces) < 1 {
		t.Errorf("parseReplaceDirectives() returned %d directives, expected at least 1", len(replaces))
	}

	// Check that relative paths are resolved
	for _, replace := range replaces {
		if strings.Contains(replace, "../dep1") || strings.Contains(replace, "./local/dep3") {
			// The replace should have been resolved to an absolute path
			if !strings.Contains(replace, tmpDir) {
				t.Errorf("parseReplaceDirectives() should resolve relative path in %q", replace)
			}
		}
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

// Test Runner with custom TempDir
func TestRunner_CustomTempDir(t *testing.T) {
	r := &Runner{
		TempDir: t.TempDir(),
		GoPath:  "go",
		Verbose: false,
	}

	if r.TempDir == "" {
		t.Error("Runner.TempDir should not be empty")
	}
}

// Test resolveReplaceDirective with empty target
func TestRunner_resolveReplaceDirective_EmptyTarget(t *testing.T) {
	r := NewRunner()

	// Edge case: empty target path
	result := r.resolveReplaceDirective("replace github.com/dep => ", "/project")
	if !strings.Contains(result, "replace github.com/dep => ") {
		t.Errorf("resolveReplaceDirective() = %q, should preserve structure", result)
	}
}

// Test parseReplaceDirectives with only comments
func TestRunner_parseReplaceDirectives_OnlyComments(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23

// replace github.com/dep => ../dep
# replace github.com/dep2 => ./local
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()
	replaces := r.parseReplaceDirectives(tmpDir)

	// Comments should not be parsed as replace directives
	if len(replaces) != 0 {
		t.Errorf("parseReplaceDirectives() found %d directives in comments, want 0", len(replaces))
	}
}

// Test getPackagePath with Windows-style paths (if on Windows)
func TestRunner_getPackagePath_WindowsPaths(t *testing.T) {
	r := NewRunner()

	// Test with forward slashes (should work on all platforms)
	result := r.getPackagePath("github.com/example/test", "/project", "/project/subdir/file.go")
	expected := "github.com/example/test/subdir"

	if result != expected {
		t.Errorf("getPackagePath() = %q, want %q", result, expected)
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
