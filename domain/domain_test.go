package domain

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	coredomain "github.com/lex00/wetwire-core-go/domain"
)

func TestGitHubDomainImplementsInterface(t *testing.T) {
	// Compile-time check that GitHubDomain implements Domain
	var _ coredomain.Domain = (*GitHubDomain)(nil)
}

func TestGitHubDomainImplementsListerDomain(t *testing.T) {
	// Compile-time check that GitHubDomain implements ListerDomain
	var _ coredomain.ListerDomain = (*GitHubDomain)(nil)
}

func TestGitHubDomainImplementsGrapherDomain(t *testing.T) {
	// Compile-time check that GitHubDomain implements GrapherDomain
	var _ coredomain.GrapherDomain = (*GitHubDomain)(nil)
}

func TestGitHubDomainName(t *testing.T) {
	d := &GitHubDomain{}
	if d.Name() != "github" {
		t.Errorf("expected name 'github', got %q", d.Name())
	}
}

func TestGitHubDomainVersion(t *testing.T) {
	d := &GitHubDomain{}
	v := d.Version()
	if v == "" {
		t.Error("version should not be empty")
	}
}

func TestGitHubDomainBuilder(t *testing.T) {
	d := &GitHubDomain{}
	b := d.Builder()
	if b == nil {
		t.Error("builder should not be nil")
	}
}

func TestGitHubDomainLinter(t *testing.T) {
	d := &GitHubDomain{}
	l := d.Linter()
	if l == nil {
		t.Error("linter should not be nil")
	}
}

func TestGitHubDomainInitializer(t *testing.T) {
	d := &GitHubDomain{}
	i := d.Initializer()
	if i == nil {
		t.Error("initializer should not be nil")
	}
}

func TestGitHubDomainValidator(t *testing.T) {
	d := &GitHubDomain{}
	v := d.Validator()
	if v == nil {
		t.Error("validator should not be nil")
	}
}

func TestGitHubDomainLister(t *testing.T) {
	d := &GitHubDomain{}
	l := d.Lister()
	if l == nil {
		t.Error("lister should not be nil")
	}
}

func TestGitHubDomainGrapher(t *testing.T) {
	d := &GitHubDomain{}
	g := d.Grapher()
	if g == nil {
		t.Error("grapher should not be nil")
	}
}

func TestCreateRootCommand(t *testing.T) {
	cmd := CreateRootCommand(&GitHubDomain{})
	if cmd == nil {
		t.Fatal("root command should not be nil")
	}
	if cmd.Use != "wetwire-github" {
		t.Errorf("expected Use 'wetwire-github', got %q", cmd.Use)
	}
}

// TestLintOpts_Fields tests that LintOpts fields are correctly defined
func TestLintOpts_Fields(t *testing.T) {
	opts := LintOpts{
		Format:  "text",
		Fix:     true,
		Disable: []string{"WAG001", "WAG002"},
	}

	if opts.Format != "text" {
		t.Errorf("expected Format 'text', got %q", opts.Format)
	}
	if !opts.Fix {
		t.Error("expected Fix to be true")
	}
	if len(opts.Disable) != 2 {
		t.Errorf("expected 2 disabled rules, got %d", len(opts.Disable))
	}
	if opts.Disable[0] != "WAG001" || opts.Disable[1] != "WAG002" {
		t.Errorf("unexpected disabled rules: %v", opts.Disable)
	}
}

func TestGitHubLinter_Lint_Disable(t *testing.T) {
	// Create a temp directory with a Go file that has lint issues
	tmpDir := t.TempDir()

	// This code uses raw uses: string which triggers WAG001
	code := `package main

import "github.com/lex00/wetwire-github-go/workflow"

var TestStep = workflow.Step{
	Uses: "actions/checkout@v4",
}
`
	filePath := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(filePath, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	linter := &githubLinter{}
	ctx := &Context{}

	// Test without disabling any rules - should find WAG001
	result, err := linter.Lint(ctx, tmpDir, LintOpts{})
	if err != nil {
		t.Fatalf("Lint() error = %v", err)
	}
	if result == nil {
		t.Fatal("Lint() returned nil result")
	}

	// Should have lint issues
	foundWAG001 := false
	for _, e := range result.Errors {
		if e.Code == "WAG001" {
			foundWAG001 = true
			break
		}
	}
	if !foundWAG001 {
		t.Error("Should find WAG001 issue when not disabled")
	}

	// Test with WAG001 disabled - should NOT find WAG001
	result, err = linter.Lint(ctx, tmpDir, LintOpts{
		Disable: []string{"WAG001"},
	})
	if err != nil {
		t.Fatalf("Lint() with Disable error = %v", err)
	}
	if result == nil {
		t.Fatal("Lint() with Disable returned nil result")
	}

	// Should not have WAG001 issues
	for _, e := range result.Errors {
		if e.Code == "WAG001" {
			t.Error("WAG001 should be disabled")
		}
	}
}

func TestGitHubLinter_Lint_Fix(t *testing.T) {
	// Create a temp directory with a Go file that has lint issues
	tmpDir := t.TempDir()

	code := `package main

import "github.com/lex00/wetwire-github-go/workflow"

var TestStep = workflow.Step{
	Uses: "actions/checkout@v4",
}
`
	filePath := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(filePath, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	linter := &githubLinter{}
	ctx := &Context{}

	// Test with Fix=true - should include "auto-fix" or "Fix mode" in message
	result, err := linter.Lint(ctx, tmpDir, LintOpts{
		Fix: true,
	})
	if err != nil {
		t.Fatalf("Lint() with Fix error = %v", err)
	}
	if result == nil {
		t.Fatal("Lint() with Fix returned nil result")
	}

	// If there are issues, the message should mention Fix mode
	if len(result.Errors) > 0 {
		if !strings.Contains(result.Message, "Fix mode") && !strings.Contains(result.Message, "auto-fix") {
			t.Errorf("Expected message to mention Fix mode when Fix=true, got: %s", result.Message)
		}
	}
}

func TestGitHubLinter_Lint_DisableMultipleRules(t *testing.T) {
	// Create a temp directory with a Go file that has multiple lint issues
	tmpDir := t.TempDir()

	// Code that triggers WAG001 (raw uses:) and potentially other rules
	code := `package main

import "github.com/lex00/wetwire-github-go/workflow"

var TestStep = workflow.Step{
	Uses: "actions/checkout@v4",
}
`
	filePath := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(filePath, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	linter := &githubLinter{}
	ctx := &Context{}

	// Test with multiple rules disabled
	result, err := linter.Lint(ctx, tmpDir, LintOpts{
		Disable: []string{"WAG001", "WAG002", "WAG003"},
	})
	if err != nil {
		t.Fatalf("Lint() with multiple Disable error = %v", err)
	}
	if result == nil {
		t.Fatal("Lint() returned nil result")
	}

	// None of the disabled rules should appear
	for _, e := range result.Errors {
		if e.Code == "WAG001" || e.Code == "WAG002" || e.Code == "WAG003" {
			t.Errorf("Rule %s should be disabled but was found", e.Code)
		}
	}
}
