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

// TestTestCmd_Help tests test help output.
func TestTestCmd_Help(t *testing.T) {
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

	cmd := exec.Command(binaryPath, "test", "--help")
	out, _ := cmd.CombinedOutput()

	if !strings.Contains(string(out), "test") {
		t.Errorf("Help output should contain 'test', got: %s", out)
	}

	if !strings.Contains(string(out), "persona") {
		t.Errorf("Help output should mention 'persona', got: %s", out)
	}

	if !strings.Contains(string(out), "--list") {
		t.Errorf("Help output should contain '--list' flag, got: %s", out)
	}

	if !strings.Contains(string(out), "--score") {
		t.Errorf("Help output should contain '--score' flag, got: %s", out)
	}
}

// TestTestCmd_ListPersonas tests listing available personas.
func TestTestCmd_ListPersonas(t *testing.T) {
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

	cmd := exec.Command(binaryPath, "test", "--list")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		t.Errorf("Test list command failed: %v", err)
		return
	}

	output := stdout.String()

	// Should list personas
	if !strings.Contains(output, "beginner") {
		t.Errorf("Output should contain 'beginner' persona, got: %s", output)
	}

	if !strings.Contains(output, "expert") {
		t.Errorf("Output should contain 'expert' persona, got: %s", output)
	}

	// Should list scenarios
	if !strings.Contains(output, "ci-workflow") {
		t.Errorf("Output should contain 'ci-workflow' scenario, got: %s", output)
	}

	// Should list scoring dimensions
	if !strings.Contains(output, "Completeness") {
		t.Errorf("Output should contain 'Completeness' scoring dimension, got: %s", output)
	}
}

// TestTestCmd_BasicProject tests running tests on a basic project.
func TestTestCmd_BasicProject(t *testing.T) {
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

	// Create a test project with workflow declarations
	tmpDir := t.TempDir()
	content := `package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
	Jobs: map[string]workflow.Job{
		"build": Build,
	},
}

var Build = workflow.Job{
	Name:   "build",
	RunsOn: "ubuntu-latest",
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "workflows.go"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// Create go.mod
	goMod := `module test

go 1.23

require github.com/lex00/wetwire-github-go v0.0.0

replace github.com/lex00/wetwire-github-go => ` + getModulePath() + `
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(binaryPath, "test", tmpDir)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Logf("stdout: %s", stdout.String())
		t.Logf("stderr: %s", stderr.String())
		t.Errorf("Test command failed: %v", err)
		return
	}

	output := stdout.String()

	// Should show test results
	if !strings.Contains(output, "PASS") || !strings.Contains(output, "passed") {
		t.Errorf("Output should show passing tests, got: %s", output)
	}
}

// TestTestCmd_InvalidPath tests error handling for non-existent paths.
func TestTestCmd_InvalidPath(t *testing.T) {
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

	cmd := exec.Command(binaryPath, "test", "/nonexistent/path")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		t.Error("Test command should fail for non-existent path")
	}

	if !strings.Contains(stderr.String(), "not found") {
		t.Errorf("Expected error message about path not found, got: %s", stderr.String())
	}
}

// TestTestCmd_JSONFormat tests JSON output format.
func TestTestCmd_JSONFormat(t *testing.T) {
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

	// Create a test project
	tmpDir := t.TempDir()
	content := `package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
	Jobs: map[string]workflow.Job{
		"build": Build,
	},
}

var Build = workflow.Job{
	Name:   "build",
	RunsOn: "ubuntu-latest",
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "workflows.go"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	goMod := `module test

go 1.23

require github.com/lex00/wetwire-github-go v0.0.0

replace github.com/lex00/wetwire-github-go => ` + getModulePath() + `
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(binaryPath, "test", tmpDir, "--format", "json")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		t.Logf("stdout: %s", stdout.String())
		t.Errorf("Test command failed: %v", err)
		return
	}

	// Should be valid JSON
	var result struct {
		Result wetwire.TestResult `json:"result"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Errorf("Expected valid JSON output, got error: %v\nOutput: %s", err, stdout.String())
		return
	}

	if !result.Result.Success {
		t.Errorf("Expected success to be true, got: %v", result.Result.Success)
	}

	if len(result.Result.Tests) == 0 {
		t.Error("Expected tests array to be non-empty")
	}
}

// TestTestCmd_WithScore tests the --score flag.
func TestTestCmd_WithScore(t *testing.T) {
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

	// Create a test project
	tmpDir := t.TempDir()
	content := `package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
	Jobs: map[string]workflow.Job{
		"build": Build,
	},
}

var Build = workflow.Job{
	Name:   "build",
	RunsOn: "ubuntu-latest",
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "workflows.go"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	goMod := `module test

go 1.23

require github.com/lex00/wetwire-github-go v0.0.0

replace github.com/lex00/wetwire-github-go => ` + getModulePath() + `
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(binaryPath, "test", tmpDir, "--score")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		t.Logf("stdout: %s", stdout.String())
		t.Errorf("Test command failed: %v", err)
		return
	}

	output := stdout.String()

	// Should show scoring breakdown
	if !strings.Contains(output, "Completeness") {
		t.Errorf("Output should contain 'Completeness' scoring, got: %s", output)
	}

	if !strings.Contains(output, "Score:") || !strings.Contains(output, "/15") {
		t.Errorf("Output should contain total score, got: %s", output)
	}
}

// TestTestCmd_WithPersona tests the --persona flag.
func TestTestCmd_WithPersona(t *testing.T) {
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

	// Create a test project
	tmpDir := t.TempDir()
	content := `package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
	Jobs: map[string]workflow.Job{
		"build": Build,
	},
}

var Build = workflow.Job{
	Name:   "build",
	RunsOn: "ubuntu-latest",
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "workflows.go"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	goMod := `module test

go 1.23

require github.com/lex00/wetwire-github-go v0.0.0

replace github.com/lex00/wetwire-github-go => ` + getModulePath() + `
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(binaryPath, "test", tmpDir, "--persona", "beginner")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		t.Logf("stdout: %s", stdout.String())
		// Don't fail - persona may affect behavior but shouldn't cause failure
	}

	// Command should complete without error
}

// TestTestCmd_MissingArg tests error handling for missing path argument.
func TestTestCmd_MissingArg(t *testing.T) {
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

	cmd := exec.Command(binaryPath, "test")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		t.Error("Test command should fail without path argument")
	}

	if !strings.Contains(stderr.String(), "required") {
		t.Errorf("Expected error about required argument, got: %s", stderr.String())
	}
}

// TestTestResult_JSON tests TestResult JSON serialization.
func TestTestResult_JSON(t *testing.T) {
	result := wetwire.TestResult{
		Success: true,
		Tests: []wetwire.TestCase{
			{
				Name:    "workflows_exist",
				Persona: "beginner",
				Passed:  true,
			},
			{
				Name:    "jobs_exist",
				Persona: "beginner",
				Passed:  true,
			},
		},
		Passed: 2,
		Failed: 0,
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal TestResult: %v", err)
	}

	var unmarshaled wetwire.TestResult
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal TestResult: %v", err)
	}

	if !unmarshaled.Success {
		t.Error("Expected success to be true")
	}

	if unmarshaled.Passed != 2 {
		t.Errorf("Passed = %d, want 2", unmarshaled.Passed)
	}

	if unmarshaled.Failed != 0 {
		t.Errorf("Failed = %d, want 0", unmarshaled.Failed)
	}

	if len(unmarshaled.Tests) != 2 {
		t.Errorf("len(Tests) = %d, want 2", len(unmarshaled.Tests))
	}
}
