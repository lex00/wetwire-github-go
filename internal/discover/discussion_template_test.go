package discover

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverer_DiscoverDiscussionTemplates(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "templates.go")
	content := `package main

import (
	"github.com/lex00/wetwire-github-go/templates"
)

var AnnouncementTemplate = templates.DiscussionTemplate{
	Name:  "Announcement",
	About: "Share announcements with the community",
	Body: []templates.DiscussionTemplateField{
		{Type: "markdown", Attributes: map[string]any{"value": "Share your announcement"}},
		{Type: "textarea", ID: "announcement", Attributes: map[string]any{"label": "Announcement"}},
	},
}

var IdeaTemplate = templates.DiscussionTemplate{
	Name:  "Idea",
	About: "Share your ideas",
}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverDiscussionTemplates(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverDiscussionTemplates() error = %v", err)
	}

	if len(result.Templates) != 2 {
		t.Errorf("len(result.Templates) = %d, want 2", len(result.Templates))
	}

	names := make(map[string]bool)
	for _, tmpl := range result.Templates {
		names[tmpl.Name] = true
	}

	if !names["AnnouncementTemplate"] {
		t.Error("Expected AnnouncementTemplate to be discovered")
	}
	if !names["IdeaTemplate"] {
		t.Error("Expected IdeaTemplate to be discovered")
	}
}

func TestDiscoverer_DiscoverDiscussionTemplates_Empty(t *testing.T) {
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
	result, err := d.DiscoverDiscussionTemplates(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverDiscussionTemplates() error = %v", err)
	}

	if len(result.Templates) != 0 {
		t.Errorf("Expected 0 templates, got %d", len(result.Templates))
	}
}

func TestDiscoverer_DiscoverDiscussionTemplates_SkipsVendor(t *testing.T) {
	tmpDir := t.TempDir()

	vendorDir := filepath.Join(tmpDir, "vendor")
	if err := os.MkdirAll(vendorDir, 0755); err != nil {
		t.Fatal(err)
	}

	vendorFile := filepath.Join(vendorDir, "templates.go")
	content := `package vendor

import "github.com/lex00/wetwire-github-go/templates"

var VendorTemplate = templates.DiscussionTemplate{Name: "vendor"}
`
	if err := os.WriteFile(vendorFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverDiscussionTemplates(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverDiscussionTemplates() error = %v", err)
	}

	if len(result.Templates) != 0 {
		t.Errorf("Expected no templates from vendor, got %d", len(result.Templates))
	}
}

func TestDiscoverer_DiscoverDiscussionTemplates_SkipsTestFiles(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "templates_test.go")
	content := `package main

import "github.com/lex00/wetwire-github-go/templates"

var TestTemplate = templates.DiscussionTemplate{Name: "test"}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverDiscussionTemplates(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverDiscussionTemplates() error = %v", err)
	}

	if len(result.Templates) != 0 {
		t.Errorf("Expected no templates from test files, got %d", len(result.Templates))
	}
}

func TestDiscoverer_DiscoverDiscussionTemplates_SkipsHiddenDirs(t *testing.T) {
	tmpDir := t.TempDir()

	hiddenDir := filepath.Join(tmpDir, ".hidden")
	if err := os.MkdirAll(hiddenDir, 0755); err != nil {
		t.Fatal(err)
	}

	hiddenFile := filepath.Join(hiddenDir, "templates.go")
	content := `package hidden

import "github.com/lex00/wetwire-github-go/templates"

var HiddenTemplate = templates.DiscussionTemplate{Name: "hidden"}
`
	if err := os.WriteFile(hiddenFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverDiscussionTemplates(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverDiscussionTemplates() error = %v", err)
	}

	if len(result.Templates) != 0 {
		t.Errorf("Expected no templates from hidden dir, got %d", len(result.Templates))
	}
}

func TestDiscoverer_DiscoverDiscussionTemplates_SkipsTestdata(t *testing.T) {
	tmpDir := t.TempDir()

	testdataDir := filepath.Join(tmpDir, "testdata")
	if err := os.MkdirAll(testdataDir, 0755); err != nil {
		t.Fatal(err)
	}

	testdataFile := filepath.Join(testdataDir, "templates.go")
	content := `package testdata

import "github.com/lex00/wetwire-github-go/templates"

var TestdataTemplate = templates.DiscussionTemplate{Name: "testdata"}
`
	if err := os.WriteFile(testdataFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverDiscussionTemplates(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverDiscussionTemplates() error = %v", err)
	}

	if len(result.Templates) != 0 {
		t.Errorf("Expected no templates from testdata, got %d", len(result.Templates))
	}
}

func TestDiscoverer_DiscoverDiscussionTemplates_InferredType(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "templates.go")
	content := `package main

import (
	"github.com/lex00/wetwire-github-go/templates"
)

var InferredTemplate = templates.DiscussionTemplate{
	Name:  "inferred",
	About: "Type inferred from value",
}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverDiscussionTemplates(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverDiscussionTemplates() error = %v", err)
	}

	if len(result.Templates) != 1 {
		t.Errorf("len(result.Templates) = %d, want 1", len(result.Templates))
	}

	if result.Templates[0].Name != "InferredTemplate" {
		t.Errorf("Name = %q, want %q", result.Templates[0].Name, "InferredTemplate")
	}
}

func TestDiscoverer_DiscoverDiscussionTemplates_ExplicitType(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "templates.go")
	content := `package main

import (
	"github.com/lex00/wetwire-github-go/templates"
)

var ExplicitTemplate templates.DiscussionTemplate = templates.DiscussionTemplate{
	Name:  "explicit",
	About: "Explicit type annotation",
}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverDiscussionTemplates(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverDiscussionTemplates() error = %v", err)
	}

	if len(result.Templates) != 1 {
		t.Errorf("len(result.Templates) = %d, want 1", len(result.Templates))
	}

	if result.Templates[0].Name != "ExplicitTemplate" {
		t.Errorf("Name = %q, want %q", result.Templates[0].Name, "ExplicitTemplate")
	}
}

func TestDiscoverer_DiscoverDiscussionTemplates_UnqualifiedType(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "templates.go")
	content := `package main

import (
	. "github.com/lex00/wetwire-github-go/templates"
)

var DotImportTemplate DiscussionTemplate = DiscussionTemplate{
	Name:  "dotimport",
	About: "Dot import",
}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverDiscussionTemplates(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverDiscussionTemplates() error = %v", err)
	}

	if len(result.Templates) != 1 {
		t.Errorf("len(result.Templates) = %d, want 1", len(result.Templates))
	}

	if result.Templates[0].Name != "DotImportTemplate" {
		t.Errorf("Name = %q, want %q", result.Templates[0].Name, "DotImportTemplate")
	}
}

func TestDiscoverer_DiscoverDiscussionTemplates_ParseError(t *testing.T) {
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
	result, err := d.DiscoverDiscussionTemplates(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverDiscussionTemplates() error = %v", err)
	}

	if len(result.Errors) == 0 {
		t.Error("Expected parse errors to be recorded")
	}
}

func TestDiscoverer_DiscoverDiscussionTemplates_NoImport(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "other.go")
	content := `package main

type DiscussionTemplate struct {
	Name string
}

var MyTemplate = DiscussionTemplate{Name: "local"}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverDiscussionTemplates(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverDiscussionTemplates() error = %v", err)
	}

	// Should not find any templates since templates package is not imported
	if len(result.Templates) != 0 {
		t.Errorf("Expected no templates without templates import, got %d", len(result.Templates))
	}
}
