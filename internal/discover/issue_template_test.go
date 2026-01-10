package discover

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverer_DiscoverIssueTemplates(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "templates.go")
	content := `package main

import (
	"github.com/lex00/wetwire-github-go/templates"
)

var BugReportTemplate = templates.IssueTemplate{
	Name:        "Bug Report",
	Description: "Report a bug",
	Body: []templates.IssueTemplateField{
		{Type: "markdown", Attributes: map[string]any{"value": "Thanks for reporting!"}},
		{Type: "textarea", ID: "description", Attributes: map[string]any{"label": "Description"}},
	},
}

var FeatureRequestTemplate = templates.IssueTemplate{
	Name:        "Feature Request",
	Description: "Request a new feature",
}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverIssueTemplates(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverIssueTemplates() error = %v", err)
	}

	if len(result.Templates) != 2 {
		t.Errorf("len(result.Templates) = %d, want 2", len(result.Templates))
	}

	names := make(map[string]bool)
	for _, tmpl := range result.Templates {
		names[tmpl.Name] = true
	}

	if !names["BugReportTemplate"] {
		t.Error("Expected BugReportTemplate to be discovered")
	}
	if !names["FeatureRequestTemplate"] {
		t.Error("Expected FeatureRequestTemplate to be discovered")
	}
}

func TestDiscoverer_DiscoverIssueTemplates_Empty(t *testing.T) {
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
	result, err := d.DiscoverIssueTemplates(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverIssueTemplates() error = %v", err)
	}

	if len(result.Templates) != 0 {
		t.Errorf("Expected 0 templates, got %d", len(result.Templates))
	}
}

func TestDiscoverer_DiscoverIssueTemplates_SkipsVendor(t *testing.T) {
	tmpDir := t.TempDir()

	vendorDir := filepath.Join(tmpDir, "vendor")
	if err := os.MkdirAll(vendorDir, 0755); err != nil {
		t.Fatal(err)
	}

	vendorFile := filepath.Join(vendorDir, "templates.go")
	content := `package vendor

import "github.com/lex00/wetwire-github-go/templates"

var VendorTemplate = templates.IssueTemplate{Name: "vendor"}
`
	if err := os.WriteFile(vendorFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverIssueTemplates(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverIssueTemplates() error = %v", err)
	}

	if len(result.Templates) != 0 {
		t.Errorf("Expected no templates from vendor, got %d", len(result.Templates))
	}
}

func TestDiscoverer_DiscoverIssueTemplates_SkipsTestFiles(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "templates_test.go")
	content := `package main

import "github.com/lex00/wetwire-github-go/templates"

var TestTemplate = templates.IssueTemplate{Name: "test"}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverIssueTemplates(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverIssueTemplates() error = %v", err)
	}

	if len(result.Templates) != 0 {
		t.Errorf("Expected no templates from test files, got %d", len(result.Templates))
	}
}

func TestDiscoverer_DiscoverIssueTemplates_SkipsHiddenDirs(t *testing.T) {
	tmpDir := t.TempDir()

	hiddenDir := filepath.Join(tmpDir, ".hidden")
	if err := os.MkdirAll(hiddenDir, 0755); err != nil {
		t.Fatal(err)
	}

	hiddenFile := filepath.Join(hiddenDir, "templates.go")
	content := `package hidden

import "github.com/lex00/wetwire-github-go/templates"

var HiddenTemplate = templates.IssueTemplate{Name: "hidden"}
`
	if err := os.WriteFile(hiddenFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverIssueTemplates(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverIssueTemplates() error = %v", err)
	}

	if len(result.Templates) != 0 {
		t.Errorf("Expected no templates from hidden dir, got %d", len(result.Templates))
	}
}

func TestDiscoverer_DiscoverIssueTemplates_SkipsTestdata(t *testing.T) {
	tmpDir := t.TempDir()

	testdataDir := filepath.Join(tmpDir, "testdata")
	if err := os.MkdirAll(testdataDir, 0755); err != nil {
		t.Fatal(err)
	}

	testdataFile := filepath.Join(testdataDir, "templates.go")
	content := `package testdata

import "github.com/lex00/wetwire-github-go/templates"

var TestdataTemplate = templates.IssueTemplate{Name: "testdata"}
`
	if err := os.WriteFile(testdataFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverIssueTemplates(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverIssueTemplates() error = %v", err)
	}

	if len(result.Templates) != 0 {
		t.Errorf("Expected no templates from testdata, got %d", len(result.Templates))
	}
}

func TestDiscoverer_DiscoverIssueTemplates_InferredType(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "templates.go")
	content := `package main

import (
	"github.com/lex00/wetwire-github-go/templates"
)

var InferredTemplate = templates.IssueTemplate{
	Name:        "inferred",
	Description: "Type inferred from value",
}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverIssueTemplates(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverIssueTemplates() error = %v", err)
	}

	if len(result.Templates) != 1 {
		t.Errorf("len(result.Templates) = %d, want 1", len(result.Templates))
	}

	if result.Templates[0].Name != "InferredTemplate" {
		t.Errorf("Name = %q, want %q", result.Templates[0].Name, "InferredTemplate")
	}
}

func TestDiscoverer_DiscoverIssueTemplates_ExplicitType(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "templates.go")
	content := `package main

import (
	"github.com/lex00/wetwire-github-go/templates"
)

var ExplicitTemplate templates.IssueTemplate = templates.IssueTemplate{
	Name:        "explicit",
	Description: "Explicit type annotation",
}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverIssueTemplates(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverIssueTemplates() error = %v", err)
	}

	if len(result.Templates) != 1 {
		t.Errorf("len(result.Templates) = %d, want 1", len(result.Templates))
	}

	if result.Templates[0].Name != "ExplicitTemplate" {
		t.Errorf("Name = %q, want %q", result.Templates[0].Name, "ExplicitTemplate")
	}
}

func TestDiscoverer_DiscoverIssueTemplates_UnqualifiedType(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "templates.go")
	content := `package main

import (
	. "github.com/lex00/wetwire-github-go/templates"
)

var DotImportTemplate IssueTemplate = IssueTemplate{
	Name:        "dotimport",
	Description: "Dot import",
}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverIssueTemplates(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverIssueTemplates() error = %v", err)
	}

	if len(result.Templates) != 1 {
		t.Errorf("len(result.Templates) = %d, want 1", len(result.Templates))
	}

	if result.Templates[0].Name != "DotImportTemplate" {
		t.Errorf("Name = %q, want %q", result.Templates[0].Name, "DotImportTemplate")
	}
}

func TestDiscoverer_DiscoverIssueTemplates_ParseError(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "invalid.go")
	content := `package main

import "github.com/lex00/wetwire-github-go/templates"

var InvalidSyntax = // missing value
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverIssueTemplates(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverIssueTemplates() error = %v", err)
	}

	if len(result.Errors) == 0 {
		t.Error("Expected parse errors to be recorded")
	}
}

func TestDiscoverer_DiscoverIssueTemplates_NoImport(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "other.go")
	content := `package main

type IssueTemplate struct {
	Name string
}

var MyTemplate = IssueTemplate{Name: "local"}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverIssueTemplates(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverIssueTemplates() error = %v", err)
	}

	// Should not find any templates since templates package is not imported
	if len(result.Templates) != 0 {
		t.Errorf("Expected no templates without templates import, got %d", len(result.Templates))
	}
}
