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

// TestInitCmd_Help tests init help output.
func TestInitCmd_Help(t *testing.T) {
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

	cmd := exec.Command(binaryPath, "init", "--help")
	out, _ := cmd.CombinedOutput()

	if !strings.Contains(string(out), "init") {
		t.Errorf("Help output should contain 'init', got: %s", out)
	}

	if !strings.Contains(string(out), "workflow") && !strings.Contains(string(out), "project") {
		t.Errorf("Help output should contain description about workflow project, got: %s", out)
	}

	if !strings.Contains(string(out), "--output") || !strings.Contains(string(out), "-o") {
		t.Errorf("Help output should contain '--output' or '-o' flag, got: %s", out)
	}

	if !strings.Contains(string(out), "--format") {
		t.Errorf("Help output should contain '--format' flag, got: %s", out)
	}
}

// TestInitCmd_InvalidPath tests error handling for invalid output directory.
func TestInitCmd_InvalidPath(t *testing.T) {
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

	// Try to init in a non-existent path
	cmd := exec.Command(binaryPath, "init", "test-project", "-o", "/nonexistent/invalid/path")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		t.Error("Init command should fail for invalid output path")
	}

	// Verify error message is displayed
	if stderr.Len() == 0 {
		t.Error("Expected error message on stderr")
	}
}

// TestInitCmd_BasicInit tests creating a new project.
func TestInitCmd_BasicInit(t *testing.T) {
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

	// Create temp directory for init output
	tmpDir := t.TempDir()
	projectName := "my-test-project"

	cmd := exec.Command(binaryPath, "init", projectName, "-o", tmpDir)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("Init command failed: %v\nstdout: %s\nstderr: %s", err, stdout.String(), stderr.String())
	}

	// Verify output contains expected text
	output := stdout.String()
	if !strings.Contains(output, "Created project") {
		t.Errorf("Output should contain 'Created project', got: %s", output)
	}

	// Verify directory was created
	projectDir := filepath.Join(tmpDir, projectName)
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Errorf("Project directory should exist: %s", projectDir)
	}

	// Verify expected files exist
	expectedFiles := []string{
		"go.mod",
		"README.md",
		"workflows/workflows.go",
		"workflows/triggers.go",
		"workflows/jobs.go",
		"workflows/steps.go",
		"cmd/main.go",
	}

	for _, file := range expectedFiles {
		path := filepath.Join(projectDir, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected file should exist: %s", file)
		}
	}

	// Verify go.mod content
	goModContent, err := os.ReadFile(filepath.Join(projectDir, "go.mod"))
	if err != nil {
		t.Fatalf("Failed to read go.mod: %v", err)
	}
	if !strings.Contains(string(goModContent), projectName) {
		t.Errorf("go.mod should contain project name, got: %s", goModContent)
	}
}

// TestInitCmd_JSONFormat tests JSON output format.
func TestInitCmd_JSONFormat(t *testing.T) {
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

	// Create temp directory for init output
	tmpDir := t.TempDir()
	projectName := "json-test-project"

	cmd := exec.Command(binaryPath, "init", projectName, "-o", tmpDir, "--format", "json")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		t.Fatalf("Init command failed: %v", err)
	}

	// Should be valid JSON
	var result wetwire.InitResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Errorf("Expected valid JSON output, got error: %v\nOutput: %s", err, stdout.String())
	}

	if !result.Success {
		t.Errorf("Expected success to be true, got: %v", result.Success)
	}

	if len(result.Files) == 0 {
		t.Error("Expected files array to be non-empty")
	}

	if result.OutputDir == "" {
		t.Error("Expected output_dir to be set")
	}
}

// TestInitCmd_DirectoryExists tests error handling when directory already exists.
func TestInitCmd_DirectoryExists(t *testing.T) {
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

	// Create temp directory and project directory that already exists
	tmpDir := t.TempDir()
	projectName := "existing-project"
	existingDir := filepath.Join(tmpDir, projectName)
	if err := os.MkdirAll(existingDir, 0755); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(binaryPath, "init", projectName, "-o", tmpDir)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		t.Error("Init command should fail when directory exists")
	}

	// Verify error mentions directory exists
	if !strings.Contains(stderr.String(), "already exists") {
		t.Errorf("Error should mention directory already exists, got: %s", stderr.String())
	}
}

// TestInitCmd_DirectoryExistsJSON tests JSON error output when directory exists.
func TestInitCmd_DirectoryExistsJSON(t *testing.T) {
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

	// Create temp directory and project directory that already exists
	tmpDir := t.TempDir()
	projectName := "existing-project"
	existingDir := filepath.Join(tmpDir, projectName)
	if err := os.MkdirAll(existingDir, 0755); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(binaryPath, "init", projectName, "-o", tmpDir, "--format", "json")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	_ = cmd.Run() // May return error

	// Should still be valid JSON
	var result wetwire.InitResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Errorf("Expected valid JSON output, got error: %v\nOutput: %s", err, stdout.String())
	}

	if result.Success {
		t.Error("Expected success to be false for existing directory")
	}

	if !strings.Contains(result.Error, "already exists") {
		t.Errorf("Error should mention directory already exists, got: %s", result.Error)
	}
}

// TestInitCmd_MissingArg tests error handling for missing project name.
func TestInitCmd_MissingArg(t *testing.T) {
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

	cmd := exec.Command(binaryPath, "init")
	err := cmd.Run()
	if err == nil {
		t.Error("Init command should fail without project name argument")
	}
}

// TestInitResult_JSON tests InitResult JSON serialization.
func TestInitResult_JSON(t *testing.T) {
	result := wetwire.InitResult{
		Success:   true,
		OutputDir: "/path/to/project",
		Files: []string{
			"/path/to/project/go.mod",
			"/path/to/project/workflows/workflows.go",
		},
		Error: "",
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal InitResult: %v", err)
	}

	var unmarshaled wetwire.InitResult
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal InitResult: %v", err)
	}

	if !unmarshaled.Success {
		t.Error("Expected success to be true")
	}

	if unmarshaled.OutputDir != "/path/to/project" {
		t.Errorf("OutputDir = %q, want %q", unmarshaled.OutputDir, "/path/to/project")
	}

	if len(unmarshaled.Files) != 2 {
		t.Errorf("len(Files) = %d, want 2", len(unmarshaled.Files))
	}
}
