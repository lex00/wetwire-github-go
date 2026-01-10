package main

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestVersionCmd_Output tests version command output.
func TestVersionCmd_Output(t *testing.T) {
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

	cmd := exec.Command(binaryPath, "version")
	out, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Version command failed: %v", err)
	}

	output := string(out)

	if !strings.Contains(output, "wetwire-github") {
		t.Errorf("Output should contain 'wetwire-github', got: %s", output)
	}

	if !strings.Contains(output, "commit") {
		t.Errorf("Output should contain 'commit', got: %s", output)
	}

	if !strings.Contains(output, "built") {
		t.Errorf("Output should contain 'built', got: %s", output)
	}
}

// TestVersionCmd_Help tests version help output.
func TestVersionCmd_Help(t *testing.T) {
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

	cmd := exec.Command(binaryPath, "version", "--help")
	out, _ := cmd.CombinedOutput()

	if !strings.Contains(string(out), "version") {
		t.Errorf("Help output should contain 'version', got: %s", out)
	}

	if !strings.Contains(string(out), "Print") {
		t.Errorf("Help output should contain description, got: %s", out)
	}
}

// TestGetVersion tests the getVersion function.
func TestGetVersion(t *testing.T) {
	v := getVersion()
	if v == "" {
		t.Error("getVersion() should return non-empty string")
	}
}
