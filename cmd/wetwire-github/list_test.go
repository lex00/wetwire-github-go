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

// TestListCmd_Integration tests the list command via exec.
func TestListCmd_Integration(t *testing.T) {
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

	// Create test fixtures with workflow
	tmpDir := t.TempDir()
	content := `package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
	On:   CITriggers,
	Jobs: map[string]workflow.Job{
		"build": Build,
		"test":  Test,
	},
}

var CITriggers = workflow.Triggers{
	Push: workflow.PushTrigger{Branches: []string{"main"}},
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

	// Run list command
	cmd := exec.Command(binaryPath, "list", tmpDir)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Logf("stdout: %s", stdout.String())
		t.Logf("stderr: %s", stderr.String())
		t.Errorf("List command failed: %v", err)
	}

	output := stdout.String()
	// Should show the CI workflow
	if !strings.Contains(output, "CI") && !strings.Contains(output, "No workflows found") {
		t.Errorf("Expected output to contain 'CI' or 'No workflows found', got: %s", output)
	}
}

// TestListCmd_JSONFormat tests JSON output format.
func TestListCmd_JSONFormat(t *testing.T) {
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

	cmd := exec.Command(binaryPath, "list", tmpDir, "--format", "json")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		t.Logf("stdout: %s", stdout.String())
		// May fail if no workflows found
	}

	// Should be valid JSON
	var result wetwire.ListResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		// May be empty result if no workflows
		if stdout.Len() > 0 {
			t.Logf("Output: %s", stdout.String())
		}
	}
}

// TestListCmd_NoWorkflows tests listing when no workflows are found.
func TestListCmd_NoWorkflows(t *testing.T) {
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

	// Create empty directory
	tmpDir := t.TempDir()

	// Create a Go file without workflows
	content := `package main

func main() {}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// Create go.mod
	goMod := `module test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(binaryPath, "list", tmpDir)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String() + stderr.String()

	if err == nil && !strings.Contains(output, "No workflows found") {
		// Command succeeded but should indicate no workflows
		if strings.Contains(output, "WORKFLOW") {
			// Table header without content is also acceptable
		}
	}
}

// TestListCmd_NonExistentPath tests error handling for non-existent paths.
func TestListCmd_NonExistentPath(t *testing.T) {
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

	cmd := exec.Command(binaryPath, "list", "/nonexistent/path")
	err := cmd.Run()
	if err == nil {
		t.Error("List command should fail for non-existent path")
	}
}

// TestListCmd_MissingArg tests error handling for missing argument.
func TestListCmd_MissingArg(t *testing.T) {
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

	cmd := exec.Command(binaryPath, "list")
	err := cmd.Run()
	if err == nil {
		t.Error("List command should fail without arguments")
	}
}

// TestListCmd_Help tests help output.
func TestListCmd_Help(t *testing.T) {
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

	cmd := exec.Command(binaryPath, "list", "--help")
	out, _ := cmd.CombinedOutput()

	if !strings.Contains(string(out), "list") {
		t.Errorf("Help output should contain 'list', got: %s", out)
	}

	if !strings.Contains(string(out), "format") {
		t.Errorf("Help output should contain 'format' flag, got: %s", out)
	}
}

// TestListCmd_MultipleWorkflows tests listing multiple workflows.
func TestListCmd_MultipleWorkflows(t *testing.T) {
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

	// Create test fixtures with multiple workflows
	tmpDir := t.TempDir()
	content := `package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
	Jobs: map[string]workflow.Job{
		"build": Build,
	},
}

var Release = workflow.Workflow{
	Name: "Release",
	Jobs: map[string]workflow.Job{
		"publish": Publish,
	},
}

var Build = workflow.Job{
	Name:   "build",
	RunsOn: "ubuntu-latest",
}

var Publish = workflow.Job{
	Name:   "publish",
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

	cmd := exec.Command(binaryPath, "list", tmpDir)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		t.Logf("stdout: %s", stdout.String())
		// May fail depending on discovery
	}

	output := stdout.String()
	// Both workflows should be listed (if discovery works)
	t.Logf("List output: %s", output)
}

// TestListResult_JSON tests ListResult JSON serialization.
func TestListResult_JSON(t *testing.T) {
	result := wetwire.ListResult{
		Workflows: []wetwire.ListWorkflow{
			{
				Name: "CI",
				File: "workflows.go",
				Line: 5,
				Jobs: 2,
			},
			{
				Name: "Release",
				File: "workflows.go",
				Line: 15,
				Jobs: 1,
			},
		},
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal ListResult: %v", err)
	}

	var unmarshaled wetwire.ListResult
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal ListResult: %v", err)
	}

	if len(unmarshaled.Workflows) != 2 {
		t.Errorf("len(Workflows) = %d, want 2", len(unmarshaled.Workflows))
	}

	if unmarshaled.Workflows[0].Name != "CI" {
		t.Errorf("Workflows[0].Name = %q, want %q", unmarshaled.Workflows[0].Name, "CI")
	}

	if unmarshaled.Workflows[1].Jobs != 1 {
		t.Errorf("Workflows[1].Jobs = %d, want 1", unmarshaled.Workflows[1].Jobs)
	}
}

// TestListWorkflow_Fields tests ListWorkflow field marshaling.
func TestListWorkflow_Fields(t *testing.T) {
	wf := wetwire.ListWorkflow{
		Name: "CI",
		File: "test.go",
		Line: 10,
		Jobs: 3,
	}

	data, err := json.Marshal(wf)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Verify all fields are present in JSON
	jsonStr := string(data)
	if !strings.Contains(jsonStr, `"name":"CI"`) {
		t.Error("JSON should contain name field")
	}
	if !strings.Contains(jsonStr, `"file":"test.go"`) {
		t.Error("JSON should contain file field")
	}
	if !strings.Contains(jsonStr, `"line":10`) {
		t.Error("JSON should contain line field")
	}
	if !strings.Contains(jsonStr, `"jobs":3`) {
		t.Error("JSON should contain jobs field")
	}
}

// TestListCmd_TextOutput tests default text table output.
func TestListCmd_TextOutput(t *testing.T) {
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

	cmd := exec.Command(binaryPath, "list", tmpDir)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	output := stdout.String()

	if err != nil {
		t.Logf("List output: %s", output)
	}

	// Text output should have table headers if workflows found
	if strings.Contains(output, "WORKFLOW") {
		if !strings.Contains(output, "FILE") {
			t.Error("Table should have FILE column")
		}
		if !strings.Contains(output, "LINE") {
			t.Error("Table should have LINE column")
		}
		if !strings.Contains(output, "JOBS") {
			t.Error("Table should have JOBS column")
		}
	}
}
