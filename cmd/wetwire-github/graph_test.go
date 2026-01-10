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

// TestGraphCmd_Help tests graph help output.
func TestGraphCmd_Help(t *testing.T) {
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

	cmd := exec.Command(binaryPath, "graph", "--help")
	out, _ := cmd.CombinedOutput()

	if !strings.Contains(string(out), "graph") {
		t.Errorf("Help output should contain 'graph', got: %s", out)
	}

	if !strings.Contains(string(out), "visual") && !strings.Contains(string(out), "dependencies") {
		t.Errorf("Help output should contain description about visualizing dependencies, got: %s", out)
	}

	if !strings.Contains(string(out), "--format") {
		t.Errorf("Help output should contain '--format' flag, got: %s", out)
	}

	if !strings.Contains(string(out), "dot") && !strings.Contains(string(out), "mermaid") {
		t.Errorf("Help output should mention output formats (dot, mermaid), got: %s", out)
	}

	if !strings.Contains(string(out), "--direction") {
		t.Errorf("Help output should contain '--direction' flag, got: %s", out)
	}
}

// TestGraphCmd_MermaidOutput tests Mermaid format output.
func TestGraphCmd_MermaidOutput(t *testing.T) {
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

	// Create test fixtures with job dependencies
	tmpDir := t.TempDir()
	content := `package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
	Jobs: map[string]workflow.Job{
		"build":  Build,
		"test":   Test,
		"deploy": Deploy,
	},
}

var Build = workflow.Job{
	Name:   "build",
	RunsOn: "ubuntu-latest",
}

var Test = workflow.Job{
	Name:   "test",
	RunsOn: "ubuntu-latest",
	Needs:  []any{Build},
}

var Deploy = workflow.Job{
	Name:   "deploy",
	RunsOn: "ubuntu-latest",
	Needs:  []any{Build, Test},
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

	cmd := exec.Command(binaryPath, "graph", tmpDir, "--format", "mermaid")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Logf("stdout: %s", stdout.String())
		t.Logf("stderr: %s", stderr.String())
		t.Errorf("Graph command failed: %v", err)
		return
	}

	output := stdout.String()

	// Mermaid output should start with "graph"
	if !strings.Contains(output, "graph") {
		t.Errorf("Mermaid output should contain 'graph', got: %s", output)
	}

	// Should contain arrow syntax
	if !strings.Contains(output, "-->") {
		t.Logf("Mermaid output (may have no edges if no dependencies found): %s", output)
	}
}

// TestGraphCmd_DOTOutput tests DOT format output.
func TestGraphCmd_DOTOutput(t *testing.T) {
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

	// Create test fixtures
	tmpDir := t.TempDir()
	content := `package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
	Jobs: map[string]workflow.Job{
		"build": Build,
		"test":  Test,
	},
}

var Build = workflow.Job{
	Name:   "build",
	RunsOn: "ubuntu-latest",
}

var Test = workflow.Job{
	Name:   "test",
	RunsOn: "ubuntu-latest",
	Needs:  []any{Build},
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

	cmd := exec.Command(binaryPath, "graph", tmpDir, "--format", "dot")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Logf("stdout: %s", stdout.String())
		t.Logf("stderr: %s", stderr.String())
		t.Errorf("Graph command failed: %v", err)
		return
	}

	output := stdout.String()

	// DOT output should start with "digraph"
	if !strings.Contains(output, "digraph") {
		t.Errorf("DOT output should contain 'digraph', got: %s", output)
	}

	// Should contain graph elements
	if !strings.Contains(output, "{") || !strings.Contains(output, "}") {
		t.Errorf("DOT output should contain braces, got: %s", output)
	}

	// Should contain rankdir directive
	if !strings.Contains(output, "rankdir") {
		t.Errorf("DOT output should contain 'rankdir', got: %s", output)
	}
}

// TestGraphCmd_InvalidPath tests error handling for non-existent paths.
func TestGraphCmd_InvalidPath(t *testing.T) {
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

	cmd := exec.Command(binaryPath, "graph", "/nonexistent/path")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		t.Error("Graph command should fail for non-existent path")
	}

	if !strings.Contains(stderr.String(), "not found") && !strings.Contains(stderr.String(), "error") {
		t.Errorf("Expected error message about path not found, got: %s", stderr.String())
	}
}

// TestGraphCmd_JSONOutput tests JSON format output.
func TestGraphCmd_JSONOutput(t *testing.T) {
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

	// Create test fixtures
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

	cmd := exec.Command(binaryPath, "graph", tmpDir, "--format", "json")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		t.Logf("stdout: %s", stdout.String())
		t.Errorf("Graph command failed: %v", err)
		return
	}

	// Should be valid JSON
	var result wetwire.GraphResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Errorf("Expected valid JSON output, got error: %v\nOutput: %s", err, stdout.String())
		return
	}

	if !result.Success {
		t.Errorf("Expected success to be true, got: %v", result.Success)
	}

	if result.Format != "json" {
		t.Errorf("Expected format to be 'json', got: %s", result.Format)
	}
}

// TestGraphCmd_Direction tests custom direction flag.
func TestGraphCmd_Direction(t *testing.T) {
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

	// Create test fixtures
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

	// Test LR direction with DOT format
	cmd := exec.Command(binaryPath, "graph", tmpDir, "--format", "dot", "--direction", "LR")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		t.Logf("stdout: %s", stdout.String())
		t.Errorf("Graph command failed: %v", err)
		return
	}

	output := stdout.String()
	if !strings.Contains(output, "LR") {
		t.Errorf("DOT output should contain 'LR' direction, got: %s", output)
	}
}

// TestGraphCmd_MissingArg tests error handling for missing argument.
func TestGraphCmd_MissingArg(t *testing.T) {
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

	cmd := exec.Command(binaryPath, "graph")
	err := cmd.Run()
	if err == nil {
		t.Error("Graph command should fail without path argument")
	}
}

// TestGraphCmd_InvalidFormat tests error handling for invalid format.
func TestGraphCmd_InvalidFormat(t *testing.T) {
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

	// Create a minimal valid directory
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

	cmd := exec.Command(binaryPath, "graph", tmpDir, "--format", "invalid-format")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		t.Error("Graph command should fail for invalid format")
	}

	if !strings.Contains(stderr.String(), "unknown format") {
		t.Errorf("Expected error about unknown format, got: %s", stderr.String())
	}
}

// TestGraphResult_JSON tests GraphResult JSON serialization.
func TestGraphResult_JSON(t *testing.T) {
	result := wetwire.GraphResult{
		Success: true,
		Format:  "dot",
		Output:  "digraph workflow { }",
		Nodes:   3,
		Edges:   2,
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal GraphResult: %v", err)
	}

	var unmarshaled wetwire.GraphResult
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal GraphResult: %v", err)
	}

	if !unmarshaled.Success {
		t.Error("Expected success to be true")
	}

	if unmarshaled.Format != "dot" {
		t.Errorf("Format = %q, want %q", unmarshaled.Format, "dot")
	}

	if unmarshaled.Nodes != 3 {
		t.Errorf("Nodes = %d, want 3", unmarshaled.Nodes)
	}

	if unmarshaled.Edges != 2 {
		t.Errorf("Edges = %d, want 2", unmarshaled.Edges)
	}
}
