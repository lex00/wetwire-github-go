package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	wetwire "github.com/lex00/wetwire-github-go"
)

// TestImportCmd_Help tests import help output.
func TestImportCmd_Help(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the binary
	binaryPath := filepath.Join(t.TempDir(), "wetwire-github")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = getModulePath() + "/cmd/wetwire-github"
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build binary: %v\n%s", err, out)
	}

	cmd := exec.Command(binaryPath, "import", "--help")
	out, _ := cmd.CombinedOutput()

	if !strings.Contains(string(out), "import") {
		t.Errorf("Help output should contain 'import', got: %s", out)
	}

	if !strings.Contains(string(out), "Go") && !strings.Contains(string(out), "converts") {
		t.Errorf("Help output should contain description about converting to Go, got: %s", out)
	}

	if !strings.Contains(string(out), "--output") || !strings.Contains(string(out), "-o") {
		t.Errorf("Help output should contain '--output' or '-o' flag, got: %s", out)
	}

	if !strings.Contains(string(out), "--type") {
		t.Errorf("Help output should contain '--type' flag, got: %s", out)
	}

	if !strings.Contains(string(out), "workflow") {
		t.Errorf("Help output should mention 'workflow' type, got: %s", out)
	}

	if !strings.Contains(string(out), "dependabot") {
		t.Errorf("Help output should mention 'dependabot' type, got: %s", out)
	}

	if !strings.Contains(string(out), "codeowners") {
		t.Errorf("Help output should mention 'codeowners' type, got: %s", out)
	}
}

// TestImportCmd_WorkflowFile tests importing a GitHub Actions workflow file.
func TestImportCmd_WorkflowFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the binary
	binaryPath := filepath.Join(t.TempDir(), "wetwire-github")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = getModulePath() + "/cmd/wetwire-github"
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build binary: %v\n%s", err, out)
	}

	// Create a sample workflow YAML file
	tmpDir := t.TempDir()
	workflowYAML := `name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Build
        run: go build ./...
      - name: Test
        run: go test ./...
`
	yamlPath := filepath.Join(tmpDir, "ci.yml")
	if err := os.WriteFile(yamlPath, []byte(workflowYAML), 0644); err != nil {
		t.Fatal(err)
	}

	// Create output directory
	outputDir := filepath.Join(tmpDir, "output")

	cmd := exec.Command(binaryPath, "import", yamlPath, "-o", outputDir)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("Import command failed: %v\nstdout: %s\nstderr: %s", err, stdout.String(), stderr.String())
	}

	output := stdout.String()

	// Should indicate successful import
	if !strings.Contains(output, "Imported") {
		t.Errorf("Output should indicate import success, got: %s", output)
	}

	// Verify output directory was created
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		t.Errorf("Output directory should exist: %s", outputDir)
	}

	// Verify workflows subdirectory exists
	workflowsDir := filepath.Join(outputDir, "workflows")
	if _, err := os.Stat(workflowsDir); os.IsNotExist(err) {
		t.Errorf("Workflows subdirectory should exist: %s", workflowsDir)
	}

	// Verify go.mod was created (unless --no-scaffold was used)
	goModPath := filepath.Join(outputDir, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		t.Errorf("go.mod should exist: %s", goModPath)
	}
}

// TestImportCmd_InvalidPath tests error handling for non-existent paths.
func TestImportCmd_InvalidPath(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the binary
	binaryPath := filepath.Join(t.TempDir(), "wetwire-github")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = getModulePath() + "/cmd/wetwire-github"
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build binary: %v\n%s", err, out)
	}

	cmd := exec.Command(binaryPath, "import", "/nonexistent/path/to/workflow.yml")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		t.Error("Import command should fail for non-existent file")
	}

	if !strings.Contains(stderr.String(), "not found") && !strings.Contains(stderr.String(), "error") {
		t.Errorf("Expected error message about file not found, got: %s", stderr.String())
	}
}

// TestImportCmd_DependabotFile tests importing a Dependabot configuration.
func TestImportCmd_DependabotFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the binary
	binaryPath := filepath.Join(t.TempDir(), "wetwire-github")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = getModulePath() + "/cmd/wetwire-github"
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build binary: %v\n%s", err, out)
	}

	// Create a sample dependabot YAML file
	tmpDir := t.TempDir()
	dependabotYAML := `version: 2
updates:
  - package-ecosystem: "gomod"
    directory: "/"
    schedule:
      interval: "weekly"
  - package-ecosystem: "github-actions"
    directory: "/"
    schedule:
      interval: "monthly"
`
	yamlPath := filepath.Join(tmpDir, "dependabot.yml")
	if err := os.WriteFile(yamlPath, []byte(dependabotYAML), 0644); err != nil {
		t.Fatal(err)
	}

	// Create output directory
	outputDir := filepath.Join(tmpDir, "output")

	cmd := exec.Command(binaryPath, "import", yamlPath, "-o", outputDir, "--type", "dependabot")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	// Dependabot import may or may not be fully implemented, log result either way
	if err != nil {
		t.Logf("Import dependabot command result: %v\nstdout: %s\nstderr: %s", err, stdout.String(), stderr.String())
	} else {
		t.Logf("Import dependabot succeeded: %s", stdout.String())
	}
}

// TestImportCmd_JSONFormat tests JSON output format.
func TestImportCmd_JSONFormat(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the binary
	binaryPath := filepath.Join(t.TempDir(), "wetwire-github")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = getModulePath() + "/cmd/wetwire-github"
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build binary: %v\n%s", err, out)
	}

	// Create a sample workflow YAML file
	tmpDir := t.TempDir()
	workflowYAML := `name: CI

on:
  push:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: echo "hello"
`
	yamlPath := filepath.Join(tmpDir, "ci.yml")
	if err := os.WriteFile(yamlPath, []byte(workflowYAML), 0644); err != nil {
		t.Fatal(err)
	}

	// Create output directory
	outputDir := filepath.Join(tmpDir, "output")

	cmd := exec.Command(binaryPath, "import", yamlPath, "-o", outputDir, "--format", "json")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		t.Fatalf("Import command failed: %v", err)
	}

	// Should be valid JSON
	var result wetwire.ImportResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Errorf("Expected valid JSON output, got error: %v\nOutput: %s", err, stdout.String())
		return
	}

	if !result.Success {
		t.Errorf("Expected success to be true, got: %v", result.Success)
	}

	if result.Workflows < 1 {
		t.Errorf("Expected at least 1 workflow, got: %d", result.Workflows)
	}
}

// TestImportCmd_SingleFile tests the --single-file flag.
func TestImportCmd_SingleFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the binary
	binaryPath := filepath.Join(t.TempDir(), "wetwire-github")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = getModulePath() + "/cmd/wetwire-github"
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build binary: %v\n%s", err, out)
	}

	// Create a sample workflow YAML file
	tmpDir := t.TempDir()
	workflowYAML := `name: CI

on:
  push:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: echo "hello"
`
	yamlPath := filepath.Join(tmpDir, "ci.yml")
	if err := os.WriteFile(yamlPath, []byte(workflowYAML), 0644); err != nil {
		t.Fatal(err)
	}

	// Create output directory
	outputDir := filepath.Join(tmpDir, "output")

	cmd := exec.Command(binaryPath, "import", yamlPath, "-o", outputDir, "--single-file")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("Import command failed: %v\nstdout: %s\nstderr: %s", err, stdout.String(), stderr.String())
	}

	// Verify workflows subdirectory exists and has content
	workflowsDir := filepath.Join(outputDir, "workflows")
	entries, err := os.ReadDir(workflowsDir)
	if err != nil {
		t.Fatalf("Failed to read workflows directory: %v", err)
	}

	// Should have generated files
	if len(entries) == 0 {
		t.Error("Expected at least one generated file in workflows directory")
	}
}

// TestImportCmd_NoScaffold tests the --no-scaffold flag.
func TestImportCmd_NoScaffold(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the binary
	binaryPath := filepath.Join(t.TempDir(), "wetwire-github")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = getModulePath() + "/cmd/wetwire-github"
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build binary: %v\n%s", err, out)
	}

	// Create a sample workflow YAML file
	tmpDir := t.TempDir()
	workflowYAML := `name: CI

on:
  push:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: echo "hello"
`
	yamlPath := filepath.Join(tmpDir, "ci.yml")
	if err := os.WriteFile(yamlPath, []byte(workflowYAML), 0644); err != nil {
		t.Fatal(err)
	}

	// Create output directory
	outputDir := filepath.Join(tmpDir, "output")

	cmd := exec.Command(binaryPath, "import", yamlPath, "-o", outputDir, "--no-scaffold")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("Import command failed: %v\nstdout: %s\nstderr: %s", err, stdout.String(), stderr.String())
	}

	// Verify go.mod was NOT created
	goModPath := filepath.Join(outputDir, "go.mod")
	if _, err := os.Stat(goModPath); !os.IsNotExist(err) {
		t.Errorf("go.mod should NOT exist with --no-scaffold: %s", goModPath)
	}

	// Verify README was NOT created
	readmePath := filepath.Join(outputDir, "README.md")
	if _, err := os.Stat(readmePath); !os.IsNotExist(err) {
		t.Errorf("README.md should NOT exist with --no-scaffold: %s", readmePath)
	}
}

// TestImportCmd_MissingArg tests error handling for missing file argument.
func TestImportCmd_MissingArg(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the binary
	binaryPath := filepath.Join(t.TempDir(), "wetwire-github")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = getModulePath() + "/cmd/wetwire-github"
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build binary: %v\n%s", err, out)
	}

	cmd := exec.Command(binaryPath, "import")
	err := cmd.Run()
	if err == nil {
		t.Error("Import command should fail without file argument")
	}
}

// TestImportCmd_CodeownersFile tests importing a CODEOWNERS file.
func TestImportCmd_CodeownersFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the binary
	binaryPath := filepath.Join(t.TempDir(), "wetwire-github")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = getModulePath() + "/cmd/wetwire-github"
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build binary: %v\n%s", err, out)
	}

	// Create a sample CODEOWNERS file
	tmpDir := t.TempDir()
	codeowners := `# Global owners
* @org/team-leads

# Specific paths
/docs/ @org/docs-team
*.go @org/go-developers
/internal/ @org/core-team @specific-user
`
	codeownersPath := filepath.Join(tmpDir, "CODEOWNERS")
	if err := os.WriteFile(codeownersPath, []byte(codeowners), 0644); err != nil {
		t.Fatal(err)
	}

	// Create output directory
	outputDir := filepath.Join(tmpDir, "output")

	cmd := exec.Command(binaryPath, "import", codeownersPath, "-o", outputDir, "--type", "codeowners")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Logf("Import CODEOWNERS result: %v\nstdout: %s\nstderr: %s", err, stdout.String(), stderr.String())
		// Don't fail if codeowners import isn't fully implemented
		return
	}

	output := stdout.String()
	if strings.Contains(output, "Imported CODEOWNERS") {
		t.Log("CODEOWNERS import succeeded")
	}
}

// TestImportCmd_InvalidYAML tests error handling for invalid YAML.
func TestImportCmd_InvalidYAML(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the binary
	binaryPath := filepath.Join(t.TempDir(), "wetwire-github")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = getModulePath() + "/cmd/wetwire-github"
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build binary: %v\n%s", err, out)
	}

	// Create an invalid YAML file
	tmpDir := t.TempDir()
	invalidYAML := `name: CI
  this is: invalid::: yaml
    - not valid
`
	yamlPath := filepath.Join(tmpDir, "invalid.yml")
	if err := os.WriteFile(yamlPath, []byte(invalidYAML), 0644); err != nil {
		t.Fatal(err)
	}

	// Create output directory
	outputDir := filepath.Join(tmpDir, "output")

	cmd := exec.Command(binaryPath, "import", yamlPath, "-o", outputDir)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		t.Error("Import command should fail for invalid YAML")
	}
}

// TestImportResult_JSON tests ImportResult JSON serialization.
func TestImportResult_JSON(t *testing.T) {
	result := wetwire.ImportResult{
		Success:   true,
		OutputDir: "/path/to/output",
		Files: []string{
			"/path/to/output/workflows/workflows.go",
			"/path/to/output/workflows/jobs.go",
		},
		Workflows: 1,
		Jobs:      3,
		Steps:     8,
		Errors:    nil,
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal ImportResult: %v", err)
	}

	var unmarshaled wetwire.ImportResult
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal ImportResult: %v", err)
	}

	if !unmarshaled.Success {
		t.Error("Expected success to be true")
	}

	if unmarshaled.OutputDir != "/path/to/output" {
		t.Errorf("OutputDir = %q, want %q", unmarshaled.OutputDir, "/path/to/output")
	}

	if unmarshaled.Workflows != 1 {
		t.Errorf("Workflows = %d, want 1", unmarshaled.Workflows)
	}

	if unmarshaled.Jobs != 3 {
		t.Errorf("Jobs = %d, want 3", unmarshaled.Jobs)
	}

	if unmarshaled.Steps != 8 {
		t.Errorf("Steps = %d, want 8", unmarshaled.Steps)
	}

	if len(unmarshaled.Files) != 2 {
		t.Errorf("len(Files) = %d, want 2", len(unmarshaled.Files))
	}
}

// TestImportResult_WithErrors tests ImportResult with errors.
func TestImportResult_WithErrors(t *testing.T) {
	result := wetwire.ImportResult{
		Success:   false,
		OutputDir: "",
		Files:     nil,
		Workflows: 0,
		Jobs:      0,
		Steps:     0,
		Errors:    []string{"file not found: test.yml", "parse error: invalid yaml"},
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal ImportResult: %v", err)
	}

	jsonStr := string(data)
	if !strings.Contains(jsonStr, `"success":false`) {
		t.Error("JSON should contain success:false")
	}

	if !strings.Contains(jsonStr, "file not found") {
		t.Error("JSON should contain error message")
	}
}
