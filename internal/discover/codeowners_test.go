package discover

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverer_DiscoverCodeowners(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "codeowners.go")
	content := `package main

import (
	"github.com/lex00/wetwire-github-go/codeowners"
)

var CodeOwners = codeowners.Owners{
	Rules: []codeowners.Rule{
		{Pattern: "*", Owners: []string{"@default-team"}},
		{Pattern: "/docs/", Owners: []string{"@docs-team"}},
	},
}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverCodeowners(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverCodeowners() error = %v", err)
	}

	if len(result.Configs) != 1 {
		t.Errorf("len(result.Configs) = %d, want 1", len(result.Configs))
	}

	if result.Configs[0].Name != "CodeOwners" {
		t.Errorf("Name = %q, want %q", result.Configs[0].Name, "CodeOwners")
	}
}

func TestDiscoverer_DiscoverCodeowners_Empty(t *testing.T) {
	tmpDir := t.TempDir()

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
	result, err := d.DiscoverCodeowners(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverCodeowners() error = %v", err)
	}

	if len(result.Configs) != 0 {
		t.Errorf("Expected 0 configs, got %d", len(result.Configs))
	}
}

func TestDiscoverer_DiscoverCodeowners_SkipsVendor(t *testing.T) {
	tmpDir := t.TempDir()

	vendorDir := filepath.Join(tmpDir, "vendor")
	if err := os.MkdirAll(vendorDir, 0755); err != nil {
		t.Fatal(err)
	}

	vendorFile := filepath.Join(vendorDir, "codeowners.go")
	content := `package vendor

import "github.com/lex00/wetwire-github-go/codeowners"

var VendorOwners = codeowners.Owners{}
`
	if err := os.WriteFile(vendorFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverCodeowners(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverCodeowners() error = %v", err)
	}

	if len(result.Configs) != 0 {
		t.Errorf("Expected no configs from vendor, got %d", len(result.Configs))
	}
}

func TestDiscoverer_DiscoverCodeowners_SkipsTestFiles(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "codeowners_test.go")
	content := `package main

import "github.com/lex00/wetwire-github-go/codeowners"

var TestOwners = codeowners.Owners{}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverCodeowners(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverCodeowners() error = %v", err)
	}

	if len(result.Configs) != 0 {
		t.Errorf("Expected no configs from test files, got %d", len(result.Configs))
	}
}

func TestDiscoverer_DiscoverCodeowners_InferredType(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "codeowners.go")
	content := `package main

import (
	"github.com/lex00/wetwire-github-go/codeowners"
)

var InferredOwners = codeowners.Owners{
	Rules: []codeowners.Rule{
		{Pattern: "*", Owners: []string{"@team"}},
	},
}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverCodeowners(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverCodeowners() error = %v", err)
	}

	if len(result.Configs) != 1 {
		t.Errorf("len(result.Configs) = %d, want 1", len(result.Configs))
	}

	if result.Configs[0].Name != "InferredOwners" {
		t.Errorf("Name = %q, want %q", result.Configs[0].Name, "InferredOwners")
	}
}

func TestDiscoverer_DiscoverCodeowners_MultipleConfigs(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "codeowners.go")
	content := `package main

import (
	"github.com/lex00/wetwire-github-go/codeowners"
)

var MainCodeOwners = codeowners.Owners{
	Rules: []codeowners.Rule{
		{Pattern: "*", Owners: []string{"@main-team"}},
	},
}

var DocsCodeOwners = codeowners.Owners{
	Rules: []codeowners.Rule{
		{Pattern: "/docs/", Owners: []string{"@docs-team"}},
	},
}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverCodeowners(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverCodeowners() error = %v", err)
	}

	if len(result.Configs) != 2 {
		t.Errorf("len(result.Configs) = %d, want 2", len(result.Configs))
	}

	names := make(map[string]bool)
	for _, cfg := range result.Configs {
		names[cfg.Name] = true
	}

	if !names["MainCodeOwners"] {
		t.Error("Expected MainCodeOwners to be discovered")
	}
	if !names["DocsCodeOwners"] {
		t.Error("Expected DocsCodeOwners to be discovered")
	}
}
