package discover

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverer_DiscoverPRTemplates(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "templates.go")
	content := `package main

import (
	"github.com/lex00/wetwire-github-go/templates"
)

var DefaultPRTemplate = templates.PRTemplate{
	Content: "## Description\n\nPlease describe your changes.\n",
}

var FeaturePRTemplate = templates.PRTemplate{
	Name:    "feature",
	Content: "## Feature\n\nDescribe the new feature.\n",
}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverPRTemplates(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverPRTemplates() error = %v", err)
	}

	if len(result.Templates) != 2 {
		t.Errorf("len(result.Templates) = %d, want 2", len(result.Templates))
	}

	// Check template names
	names := make(map[string]bool)
	for _, tmpl := range result.Templates {
		names[tmpl.Name] = true
	}

	if !names["DefaultPRTemplate"] {
		t.Error("Expected DefaultPRTemplate to be discovered")
	}
	if !names["FeaturePRTemplate"] {
		t.Error("Expected FeaturePRTemplate to be discovered")
	}
}

func TestDiscoverer_DiscoverPRTemplates_Empty(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a Go file without PR templates
	testFile := filepath.Join(tmpDir, "other.go")
	content := `package main

import "fmt"

func main() {
	fmt.Println("hello")
}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverPRTemplates(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverPRTemplates() error = %v", err)
	}

	if len(result.Templates) != 0 {
		t.Errorf("Expected 0 templates, got %d", len(result.Templates))
	}
}

func TestDiscoverer_DiscoverPRTemplates_SkipsVendor(t *testing.T) {
	tmpDir := t.TempDir()

	vendorDir := filepath.Join(tmpDir, "vendor")
	if err := os.MkdirAll(vendorDir, 0755); err != nil {
		t.Fatal(err)
	}

	vendorFile := filepath.Join(vendorDir, "templates.go")
	content := `package vendor

import "github.com/lex00/wetwire-github-go/templates"

var VendorTemplate = templates.PRTemplate{Content: "vendor content"}
`
	if err := os.WriteFile(vendorFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverPRTemplates(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverPRTemplates() error = %v", err)
	}

	if len(result.Templates) != 0 {
		t.Errorf("Expected no templates from vendor, got %d", len(result.Templates))
	}
}

func TestDiscoverer_DiscoverPRTemplates_SkipsTestFiles(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "templates_test.go")
	content := `package main

import "github.com/lex00/wetwire-github-go/templates"

var TestTemplate = templates.PRTemplate{Content: "test content"}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverPRTemplates(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverPRTemplates() error = %v", err)
	}

	if len(result.Templates) != 0 {
		t.Errorf("Expected no templates from test files, got %d", len(result.Templates))
	}
}

func TestDiscoverer_DiscoverPRTemplates_InferredType(t *testing.T) {
	tmpDir := t.TempDir()

	// Test discovery with type inference (no explicit type annotation)
	testFile := filepath.Join(tmpDir, "templates.go")
	content := `package main

import (
	"github.com/lex00/wetwire-github-go/templates"
)

var InferredTemplate = templates.PRTemplate{
	Name:    "inferred",
	Content: "## Inferred\n\nType inferred from value.\n",
}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverPRTemplates(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverPRTemplates() error = %v", err)
	}

	if len(result.Templates) != 1 {
		t.Errorf("len(result.Templates) = %d, want 1", len(result.Templates))
	}

	if result.Templates[0].Name != "InferredTemplate" {
		t.Errorf("Template name = %q, want %q", result.Templates[0].Name, "InferredTemplate")
	}
}

func TestDiscoverer_DiscoverPRTemplates_ExplicitType(t *testing.T) {
	tmpDir := t.TempDir()

	// Test discovery with explicit type annotation
	testFile := filepath.Join(tmpDir, "templates.go")
	content := `package main

import (
	"github.com/lex00/wetwire-github-go/templates"
)

var ExplicitTemplate templates.PRTemplate = templates.PRTemplate{
	Name:    "explicit",
	Content: "## Explicit\n\nExplicit type annotation.\n",
}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverPRTemplates(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverPRTemplates() error = %v", err)
	}

	if len(result.Templates) != 1 {
		t.Errorf("len(result.Templates) = %d, want 1", len(result.Templates))
	}

	if result.Templates[0].Name != "ExplicitTemplate" {
		t.Errorf("Template name = %q, want %q", result.Templates[0].Name, "ExplicitTemplate")
	}
}
