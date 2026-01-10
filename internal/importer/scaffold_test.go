package importer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewScaffold(t *testing.T) {
	s := NewScaffold("github.com/user/repo", "MyProject")
	if s == nil {
		t.Fatal("NewScaffold() returned nil")
	}
	if s.ModulePath != "github.com/user/repo" {
		t.Errorf("ModulePath = %q, want %q", s.ModulePath, "github.com/user/repo")
	}
	if s.ProjectName != "MyProject" {
		t.Errorf("ProjectName = %q, want %q", s.ProjectName, "MyProject")
	}
}

func TestScaffold_Generate(t *testing.T) {
	s := &Scaffold{
		ModulePath:  "github.com/example/test",
		ProjectName: "TestProject",
	}

	files := s.Generate()
	if files == nil {
		t.Fatal("Generate() returned nil")
	}

	expectedFiles := []string{
		"go.mod",
		"cmd/main.go",
		"README.md",
		"CLAUDE.md",
		".gitignore",
	}

	for _, filename := range expectedFiles {
		if _, ok := files.Files[filename]; !ok {
			t.Errorf("Generate() missing expected file: %s", filename)
		}
	}
}

func TestScaffold_GenerateGoMod(t *testing.T) {
	s := &Scaffold{
		ModulePath:  "github.com/test/mod",
		ProjectName: "Test",
	}

	content := s.generateGoMod()

	if !strings.Contains(content, "module github.com/test/mod") {
		t.Error("Missing module path")
	}
	if !strings.Contains(content, "go 1.23") {
		t.Error("Missing go version")
	}
	if !strings.Contains(content, "require github.com/lex00/wetwire-github-go") {
		t.Error("Missing wetwire-github-go dependency")
	}
}

func TestScaffold_GenerateMain(t *testing.T) {
	s := &Scaffold{
		ModulePath:  "github.com/test/project",
		ProjectName: "MyApp",
	}

	content := s.generateMain()

	if !strings.Contains(content, "package main") {
		t.Error("Missing package main")
	}
	if !strings.Contains(content, `_ "github.com/test/project/workflows"`) {
		t.Error("Missing workflows import")
	}
	if !strings.Contains(content, "wetwire-github project: MyApp") {
		t.Error("Missing project name in output")
	}
	if !strings.Contains(content, "wetwire-github build .") {
		t.Error("Missing build command")
	}
}

func TestScaffold_GenerateReadme(t *testing.T) {
	s := &Scaffold{
		ModulePath:  "github.com/test/repo",
		ProjectName: "TestApp",
	}

	content := s.generateReadme()

	if !strings.Contains(content, "# TestApp") {
		t.Error("Missing project title")
	}
	if !strings.Contains(content, "wetwire-github build .") {
		t.Error("Missing build command")
	}
	if !strings.Contains(content, "wetwire-github validate") {
		t.Error("Missing validate command")
	}
	if !strings.Contains(content, "wetwire-github lint") {
		t.Error("Missing lint command")
	}
	if !strings.Contains(content, "## Building") {
		t.Error("Missing Building section")
	}
	if !strings.Contains(content, "## Structure") {
		t.Error("Missing Structure section")
	}
}

func TestScaffold_GenerateClaudeMD(t *testing.T) {
	s := &Scaffold{
		ModulePath:  "github.com/test/repo",
		ProjectName: "TestApp",
	}

	content := s.generateClaudeMD()

	if !strings.Contains(content, "# TestApp") {
		t.Error("Missing project title")
	}
	if !strings.Contains(content, "## Syntax") {
		t.Error("Missing Syntax section")
	}
	if !strings.Contains(content, "var CI = workflow.Workflow{") {
		t.Error("Missing example code")
	}
	if !strings.Contains(content, "wetwire-github build .") {
		t.Error("Missing build command")
	}
}

func TestScaffold_GenerateGitignore(t *testing.T) {
	s := &Scaffold{
		ModulePath:  "github.com/test/repo",
		ProjectName: "Test",
	}

	content := s.generateGitignore()

	if !strings.Contains(content, "*.exe") {
		t.Error("Missing exe pattern")
	}
	if !strings.Contains(content, "*.test") {
		t.Error("Missing test pattern")
	}
	if !strings.Contains(content, "vendor/") {
		t.Error("Missing vendor directory")
	}
	if !strings.Contains(content, ".idea/") {
		t.Error("Missing IDE directory")
	}
	if !strings.Contains(content, "# Binaries") {
		t.Error("Missing comment headers")
	}
}

func TestWriteScaffold(t *testing.T) {
	tmpDir := t.TempDir()

	files := &ScaffoldFiles{
		Files: map[string]string{
			"go.mod":      "module test\n\ngo 1.23\n",
			"cmd/main.go": "package main\n\nfunc main() {}\n",
			"README.md":   "# Test\n",
		},
	}

	err := WriteScaffold(tmpDir, files)
	if err != nil {
		t.Fatalf("WriteScaffold() error = %v", err)
	}

	// Verify files were written
	goModPath := filepath.Join(tmpDir, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		t.Error("go.mod was not created")
	}

	mainPath := filepath.Join(tmpDir, "cmd/main.go")
	if _, err := os.Stat(mainPath); os.IsNotExist(err) {
		t.Error("cmd/main.go was not created")
	}

	// Verify content
	content, err := os.ReadFile(goModPath)
	if err != nil {
		t.Fatalf("Failed to read go.mod: %v", err)
	}
	if string(content) != "module test\n\ngo 1.23\n" {
		t.Errorf("go.mod content mismatch")
	}
}

func TestWriteScaffold_Error(t *testing.T) {
	// Try to write to a directory we can't access
	files := &ScaffoldFiles{
		Files: map[string]string{
			"test.go": "package test\n",
		},
	}

	err := WriteScaffold("/nonexistent/path/that/should/not/exist", files)
	if err == nil {
		t.Error("WriteScaffold() expected error for invalid path")
	}
}

func TestWriteGeneratedCode(t *testing.T) {
	tmpDir := t.TempDir()

	code := &GeneratedCode{
		Files: map[string]string{
			"workflows.go": "package workflows\n\nvar CI = workflow.Workflow{}\n",
			"jobs.go":      "package workflows\n\nvar Build = workflow.Job{}\n",
			"empty.go":     "   \n  \n",
		},
	}

	err := WriteGeneratedCode(tmpDir, code)
	if err != nil {
		t.Fatalf("WriteGeneratedCode() error = %v", err)
	}

	// Verify files were written
	workflowsPath := filepath.Join(tmpDir, "workflows.go")
	if _, err := os.Stat(workflowsPath); os.IsNotExist(err) {
		t.Error("workflows.go was not created")
	}

	// Verify empty file was skipped
	emptyPath := filepath.Join(tmpDir, "empty.go")
	if _, err := os.Stat(emptyPath); !os.IsNotExist(err) {
		t.Error("empty.go should not have been created")
	}

	// Verify content
	content, err := os.ReadFile(workflowsPath)
	if err != nil {
		t.Fatalf("Failed to read workflows.go: %v", err)
	}
	expected := "package workflows\n\nvar CI = workflow.Workflow{}\n"
	if string(content) != expected {
		t.Errorf("workflows.go content = %q, want %q", string(content), expected)
	}
}

func TestWriteGeneratedCode_Error(t *testing.T) {
	code := &GeneratedCode{
		Files: map[string]string{
			"test.go": "package test\n",
		},
	}

	err := WriteGeneratedCode("/nonexistent/path/that/should/not/exist", code)
	if err == nil {
		t.Error("WriteGeneratedCode() expected error for invalid path")
	}
}

func TestScaffold_IntegrationGenerate(t *testing.T) {
	tmpDir := t.TempDir()

	s := NewScaffold("github.com/example/integration", "IntegrationTest")
	files := s.Generate()

	err := WriteScaffold(tmpDir, files)
	if err != nil {
		t.Fatalf("WriteScaffold() error = %v", err)
	}

	// Verify all expected files exist
	expectedFiles := map[string]bool{
		"go.mod":      false,
		"cmd/main.go": false,
		"README.md":   false,
		"CLAUDE.md":   false,
		".gitignore":  false,
	}

	for filename := range expectedFiles {
		path := filepath.Join(tmpDir, filename)
		if _, err := os.Stat(path); err == nil {
			expectedFiles[filename] = true
		}
	}

	for filename, exists := range expectedFiles {
		if !exists {
			t.Errorf("Expected file %s was not created", filename)
		}
	}

	// Verify the content integrity of go.mod
	goModPath := filepath.Join(tmpDir, "go.mod")
	goModContent, err := os.ReadFile(goModPath)
	if err != nil {
		t.Fatalf("Failed to read go.mod: %v", err)
	}
	if !strings.Contains(string(goModContent), "github.com/example/integration") {
		t.Error("go.mod does not contain correct module path")
	}
}
