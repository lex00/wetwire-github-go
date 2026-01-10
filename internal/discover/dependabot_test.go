package discover

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverer_DiscoverDependabot(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "dependabot.go")
	content := `package main

import (
	"github.com/lex00/wetwire-github-go/dependabot"
)

var DependabotConfig = dependabot.Dependabot{
	Version: 2,
	Updates: []dependabot.Update{
		{
			PackageEcosystem: "gomod",
			Directory:        "/",
			Schedule: dependabot.Schedule{
				Interval: "weekly",
			},
		},
	},
}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverDependabot(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverDependabot() error = %v", err)
	}

	if len(result.Configs) != 1 {
		t.Errorf("len(result.Configs) = %d, want 1", len(result.Configs))
	}

	if result.Configs[0].Name != "DependabotConfig" {
		t.Errorf("Name = %q, want %q", result.Configs[0].Name, "DependabotConfig")
	}
}

func TestDiscoverer_DiscoverDependabot_Empty(t *testing.T) {
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
	result, err := d.DiscoverDependabot(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverDependabot() error = %v", err)
	}

	if len(result.Configs) != 0 {
		t.Errorf("Expected 0 configs, got %d", len(result.Configs))
	}
}

func TestDiscoverer_DiscoverDependabot_SkipsVendor(t *testing.T) {
	tmpDir := t.TempDir()

	vendorDir := filepath.Join(tmpDir, "vendor")
	if err := os.MkdirAll(vendorDir, 0755); err != nil {
		t.Fatal(err)
	}

	vendorFile := filepath.Join(vendorDir, "dependabot.go")
	content := `package vendor

import "github.com/lex00/wetwire-github-go/dependabot"

var VendorConfig = dependabot.Dependabot{Version: 2}
`
	if err := os.WriteFile(vendorFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverDependabot(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverDependabot() error = %v", err)
	}

	if len(result.Configs) != 0 {
		t.Errorf("Expected no configs from vendor, got %d", len(result.Configs))
	}
}

func TestDiscoverer_DiscoverDependabot_SkipsTestFiles(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "dependabot_test.go")
	content := `package main

import "github.com/lex00/wetwire-github-go/dependabot"

var TestConfig = dependabot.Dependabot{Version: 2}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverDependabot(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverDependabot() error = %v", err)
	}

	if len(result.Configs) != 0 {
		t.Errorf("Expected no configs from test files, got %d", len(result.Configs))
	}
}

func TestDiscoverer_DiscoverDependabot_SkipsHiddenDirs(t *testing.T) {
	tmpDir := t.TempDir()

	hiddenDir := filepath.Join(tmpDir, ".hidden")
	if err := os.MkdirAll(hiddenDir, 0755); err != nil {
		t.Fatal(err)
	}

	hiddenFile := filepath.Join(hiddenDir, "dependabot.go")
	content := `package hidden

import "github.com/lex00/wetwire-github-go/dependabot"

var HiddenConfig = dependabot.Dependabot{Version: 2}
`
	if err := os.WriteFile(hiddenFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverDependabot(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverDependabot() error = %v", err)
	}

	if len(result.Configs) != 0 {
		t.Errorf("Expected no configs from hidden dir, got %d", len(result.Configs))
	}
}

func TestDiscoverer_DiscoverDependabot_SkipsTestdata(t *testing.T) {
	tmpDir := t.TempDir()

	testdataDir := filepath.Join(tmpDir, "testdata")
	if err := os.MkdirAll(testdataDir, 0755); err != nil {
		t.Fatal(err)
	}

	testdataFile := filepath.Join(testdataDir, "dependabot.go")
	content := `package testdata

import "github.com/lex00/wetwire-github-go/dependabot"

var TestdataConfig = dependabot.Dependabot{Version: 2}
`
	if err := os.WriteFile(testdataFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverDependabot(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverDependabot() error = %v", err)
	}

	if len(result.Configs) != 0 {
		t.Errorf("Expected no configs from testdata, got %d", len(result.Configs))
	}
}

func TestDiscoverer_DiscoverDependabot_InferredType(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "dependabot.go")
	content := `package main

import (
	"github.com/lex00/wetwire-github-go/dependabot"
)

var InferredConfig = dependabot.Dependabot{
	Version: 2,
	Updates: []dependabot.Update{
		{PackageEcosystem: "npm"},
	},
}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverDependabot(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverDependabot() error = %v", err)
	}

	if len(result.Configs) != 1 {
		t.Errorf("len(result.Configs) = %d, want 1", len(result.Configs))
	}

	if result.Configs[0].Name != "InferredConfig" {
		t.Errorf("Name = %q, want %q", result.Configs[0].Name, "InferredConfig")
	}
}

func TestDiscoverer_DiscoverDependabot_ExplicitType(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "dependabot.go")
	content := `package main

import (
	"github.com/lex00/wetwire-github-go/dependabot"
)

var ExplicitConfig dependabot.Dependabot = dependabot.Dependabot{
	Version: 2,
}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverDependabot(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverDependabot() error = %v", err)
	}

	if len(result.Configs) != 1 {
		t.Errorf("len(result.Configs) = %d, want 1", len(result.Configs))
	}

	if result.Configs[0].Name != "ExplicitConfig" {
		t.Errorf("Name = %q, want %q", result.Configs[0].Name, "ExplicitConfig")
	}
}

func TestDiscoverer_DiscoverDependabot_UnqualifiedType(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "dependabot.go")
	content := `package main

import (
	. "github.com/lex00/wetwire-github-go/dependabot"
)

var DotImportConfig Dependabot = Dependabot{
	Version: 2,
}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverDependabot(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverDependabot() error = %v", err)
	}

	if len(result.Configs) != 1 {
		t.Errorf("len(result.Configs) = %d, want 1", len(result.Configs))
	}

	if result.Configs[0].Name != "DotImportConfig" {
		t.Errorf("Name = %q, want %q", result.Configs[0].Name, "DotImportConfig")
	}
}

func TestDiscoverer_DiscoverDependabot_MultipleConfigs(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "dependabot.go")
	content := `package main

import (
	"github.com/lex00/wetwire-github-go/dependabot"
)

var GoConfig = dependabot.Dependabot{
	Version: 2,
	Updates: []dependabot.Update{
		{PackageEcosystem: "gomod"},
	},
}

var NpmConfig = dependabot.Dependabot{
	Version: 2,
	Updates: []dependabot.Update{
		{PackageEcosystem: "npm"},
	},
}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverDependabot(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverDependabot() error = %v", err)
	}

	if len(result.Configs) != 2 {
		t.Errorf("len(result.Configs) = %d, want 2", len(result.Configs))
	}

	names := make(map[string]bool)
	for _, cfg := range result.Configs {
		names[cfg.Name] = true
	}

	if !names["GoConfig"] {
		t.Error("Expected GoConfig to be discovered")
	}
	if !names["NpmConfig"] {
		t.Error("Expected NpmConfig to be discovered")
	}
}

func TestDiscoverer_DiscoverDependabot_ParseError(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "invalid.go")
	content := `package main

import "github.com/lex00/wetwire-github-go/dependabot"

var InvalidSyntax = // missing value
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverDependabot(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverDependabot() error = %v", err)
	}

	if len(result.Errors) == 0 {
		t.Error("Expected parse errors to be recorded")
	}
}

func TestDiscoverer_DiscoverDependabot_NoImport(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "other.go")
	content := `package main

type Dependabot struct {
	Version int
}

var MyConfig = Dependabot{Version: 2}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	d := NewDiscoverer()
	result, err := d.DiscoverDependabot(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverDependabot() error = %v", err)
	}

	// Should not find any configs since dependabot package is not imported
	if len(result.Configs) != 0 {
		t.Errorf("Expected no configs without dependabot import, got %d", len(result.Configs))
	}
}
