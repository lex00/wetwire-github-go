package main

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestDesignCmd_ProviderFlagHelp tests that the --provider flag appears in design help.
func TestDesignCmd_ProviderFlagHelp(t *testing.T) {
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

	cmd := exec.Command(binaryPath, "design", "--help")
	out, _ := cmd.CombinedOutput()

	if !strings.Contains(string(out), "--provider") {
		t.Errorf("Design help output should contain '--provider' flag, got: %s", out)
	}

	if !strings.Contains(string(out), "anthropic") {
		t.Errorf("Design help output should mention 'anthropic' provider, got: %s", out)
	}

	if !strings.Contains(string(out), "kiro") {
		t.Errorf("Design help output should mention 'kiro' provider, got: %s", out)
	}
}

// TestTestCmd_ProviderFlagHelp tests that the --provider flag appears in test help.
func TestTestCmd_ProviderFlagHelp(t *testing.T) {
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

	if !strings.Contains(string(out), "--provider") {
		t.Errorf("Test help output should contain '--provider' flag, got: %s", out)
	}

	if !strings.Contains(string(out), "anthropic") {
		t.Errorf("Test help output should mention 'anthropic' provider, got: %s", out)
	}

	if !strings.Contains(string(out), "kiro") {
		t.Errorf("Test help output should mention 'kiro' provider, got: %s", out)
	}
}

// TestDesignCmd_InvalidProvider tests error handling for invalid provider.
func TestDesignCmd_InvalidProvider(t *testing.T) {
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

	cmd := exec.Command(binaryPath, "design", "--provider", "invalid", "test prompt")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		t.Error("Design command should fail with invalid provider")
	}

	if !strings.Contains(stderr.String(), "unknown provider") {
		t.Errorf("Expected error about unknown provider, got: %s", stderr.String())
	}
}

// TestTestCmd_InvalidProvider tests error handling for invalid provider.
func TestTestCmd_InvalidProvider(t *testing.T) {
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

	cmd := exec.Command(binaryPath, "test", "--provider", "invalid", "test prompt")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		t.Error("Test command should fail with invalid provider")
	}

	if !strings.Contains(stderr.String(), "unknown provider") {
		t.Errorf("Expected error about unknown provider, got: %s", stderr.String())
	}
}

// TestProviderValidation tests the provider validation function.
func TestProviderValidation(t *testing.T) {
	tests := []struct {
		provider string
		valid    bool
	}{
		{"anthropic", true},
		{"kiro", true},
		{"invalid", false},
		{"openai", false},
		{"", false},
	}

	for _, tc := range tests {
		got := isValidProvider(tc.provider)
		if got != tc.valid {
			t.Errorf("isValidProvider(%q) = %v, want %v", tc.provider, got, tc.valid)
		}
	}
}
