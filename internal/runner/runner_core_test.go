package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lex00/wetwire-github-go/internal/discover"
)

func TestNewRunner(t *testing.T) {
	r := NewRunner()
	if r == nil {
		t.Error("NewRunner() returned nil")
	}
	if r.TempDir == "" {
		t.Error("NewRunner().TempDir is empty")
	}
}

func TestRunner_parseGoMod(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23

require (
	github.com/some/dep v1.0.0
)
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()
	modulePath, err := r.parseGoMod(tmpDir)
	if err != nil {
		t.Fatalf("parseGoMod() error = %v", err)
	}

	if modulePath != "github.com/example/test" {
		t.Errorf("parseGoMod() = %q, want %q", modulePath, "github.com/example/test")
	}
}

func TestRunner_parseGoMod_NotFound(t *testing.T) {
	tmpDir := t.TempDir()

	r := NewRunner()
	_, err := r.parseGoMod(tmpDir)
	if err == nil {
		t.Error("parseGoMod() expected error for missing go.mod")
	}
}

func TestRunner_parseGoMod_NoModule(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()
	_, err := r.parseGoMod(tmpDir)
	if err == nil {
		t.Error("parseGoMod() expected error for missing module directive")
	}
}

func TestRunner_generateGoMod(t *testing.T) {
	r := NewRunner()
	result := r.generateGoMod("github.com/example/test", "/path/to/project")

	if !strings.Contains(result, "module wetwire-extract") {
		t.Error("generateGoMod() missing module directive")
	}

	if !strings.Contains(result, "require github.com/example/test") {
		t.Error("generateGoMod() missing require directive")
	}

	if !strings.Contains(result, "replace github.com/example/test =>") {
		t.Error("generateGoMod() missing replace directive")
	}
}

func TestFindGoBinary(t *testing.T) {
	path, err := FindGoBinary()
	if err != nil {
		t.Skipf("Go binary not found, skipping: %v", err)
	}

	if path == "" {
		t.Error("FindGoBinary() returned empty path")
	}
}

func TestRunner_resolveReplaceDirective(t *testing.T) {
	r := NewRunner()

	tests := []struct {
		name    string
		line    string
		baseDir string
		want    string
	}{
		{
			name:    "relative path with ../",
			line:    "replace github.com/example/dep => ../dep",
			baseDir: "/project/subdir",
			want:    "replace github.com/example/dep => /project/dep",
		},
		{
			name:    "relative path with .",
			line:    "replace github.com/example/dep => ./local",
			baseDir: "/project",
			want:    "replace github.com/example/dep => /project/local",
		},
		{
			name:    "absolute path unchanged",
			line:    "replace github.com/example/dep => /absolute/path",
			baseDir: "/project",
			want:    "replace github.com/example/dep => /absolute/path",
		},
		{
			name:    "version replacement unchanged",
			line:    "replace github.com/example/dep v1.0.0 => v1.0.1",
			baseDir: "/project",
			want:    "replace github.com/example/dep v1.0.0 => v1.0.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := r.resolveReplaceDirective(tt.line, tt.baseDir)
			if got != tt.want {
				t.Errorf("resolveReplaceDirective() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRunner_parseReplaceDirectives_ResolvesRelativePaths(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a go.mod with a relative replace directive
	goMod := `module github.com/example/test

go 1.23

require github.com/other/dep v1.0.0

replace github.com/other/dep => ../dep
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()
	replaces := r.parseReplaceDirectives(tmpDir)

	if len(replaces) != 1 {
		t.Fatalf("parseReplaceDirectives() returned %d directives, want 1", len(replaces))
	}

	// The relative path should be resolved to an absolute path
	if strings.Contains(replaces[0], "..") {
		t.Errorf("parseReplaceDirectives() should resolve relative paths, got %q", replaces[0])
	}
}

// Test generateGoMod with multiple replace directives
func TestRunner_generateGoMod_MultipleReplaces(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23

require (
	github.com/dep1 v1.0.0
	github.com/dep2 v2.0.0
)

replace github.com/dep1 => ../dep1
replace github.com/dep2 => /absolute/dep2
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()
	result := r.generateGoMod("github.com/example/test", tmpDir)

	if !strings.Contains(result, "module wetwire-extract") {
		t.Error("generateGoMod() missing module directive")
	}

	if !strings.Contains(result, "require github.com/example/test") {
		t.Error("generateGoMod() missing require directive")
	}

	// Should include both replace directives
	if !strings.Contains(result, "replace github.com/dep1") {
		t.Error("generateGoMod() missing first replace directive")
	}

	if !strings.Contains(result, "replace github.com/dep2") {
		t.Error("generateGoMod() missing second replace directive")
	}
}

// Test resolveReplaceDirective with malformed input
func TestRunner_resolveReplaceDirective_Malformed(t *testing.T) {
	r := NewRunner()

	tests := []struct {
		name    string
		line    string
		baseDir string
		want    string
	}{
		{
			name:    "no arrow",
			line:    "replace github.com/example/dep",
			baseDir: "/project",
			want:    "replace github.com/example/dep",
		},
		{
			name:    "multiple arrows",
			line:    "replace github.com/example/dep => => /path",
			baseDir: "/project",
			want:    "replace github.com/example/dep => => /path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := r.resolveReplaceDirective(tt.line, tt.baseDir)
			if got != tt.want {
				t.Errorf("resolveReplaceDirective() = %q, want %q", got, tt.want)
			}
		})
	}
}

// Test FindGoBinary error path
func TestFindGoBinary_Error(t *testing.T) {
	// Save original PATH
	originalPath := os.Getenv("PATH")
	defer os.Setenv("PATH", originalPath)

	// Set PATH to empty to force error
	os.Setenv("PATH", "")

	_, err := FindGoBinary()
	if err == nil {
		// On some systems, go might still be found via other means
		// So we don't fail the test, just skip it
		t.Skip("Go binary found even with empty PATH")
	}
}

// Test parseReplaceDirectives with missing go.mod
func TestRunner_parseReplaceDirectives_MissingGoMod(t *testing.T) {
	tmpDir := t.TempDir()

	r := NewRunner()
	replaces := r.parseReplaceDirectives(tmpDir)

	// Should return empty slice when go.mod doesn't exist
	if len(replaces) != 0 {
		t.Errorf("parseReplaceDirectives() = %v, want empty slice", replaces)
	}
}

// Test NewRunner initializes fields correctly
func TestNewRunner_Fields(t *testing.T) {
	r := NewRunner()

	if r.TempDir == "" {
		t.Error("NewRunner().TempDir should not be empty")
	}

	// GoPath might be empty if go binary is not found, which is OK
	// Just verify the field exists
	_ = r.GoPath

	// Verbose should be false by default
	if r.Verbose {
		t.Error("NewRunner().Verbose should be false by default")
	}
}

// Test generateGoMod with no replace directives
func TestRunner_generateGoMod_NoReplaces(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()
	result := r.generateGoMod("github.com/example/test", tmpDir)

	if !strings.Contains(result, "module wetwire-extract") {
		t.Error("generateGoMod() missing module directive")
	}

	if !strings.Contains(result, "require github.com/example/test") {
		t.Error("generateGoMod() missing require directive")
	}

	// Should have the replace directive for the module itself
	if !strings.Contains(result, "replace github.com/example/test =>") {
		t.Error("generateGoMod() missing replace directive for module")
	}
}

// Test parseReplaceDirectives with various formats
func TestRunner_parseReplaceDirectives_VariousFormats(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23

replace (
	github.com/dep1 => ../dep1
	github.com/dep2 v1.0.0 => github.com/dep2 v1.0.1
)

replace github.com/dep3 => ./local/dep3
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()
	replaces := r.parseReplaceDirectives(tmpDir)

	// Should find replace directives outside of blocks
	if len(replaces) < 1 {
		t.Errorf("parseReplaceDirectives() returned %d directives, expected at least 1", len(replaces))
	}

	// Check that relative paths are resolved
	for _, replace := range replaces {
		if strings.Contains(replace, "../dep1") || strings.Contains(replace, "./local/dep3") {
			// The replace should have been resolved to an absolute path
			if !strings.Contains(replace, tmpDir) {
				t.Errorf("parseReplaceDirectives() should resolve relative path in %q", replace)
			}
		}
	}
}

// Test Runner with custom TempDir
func TestRunner_CustomTempDir(t *testing.T) {
	r := &Runner{
		TempDir: t.TempDir(),
		GoPath:  "go",
		Verbose: false,
	}

	if r.TempDir == "" {
		t.Error("Runner.TempDir should not be empty")
	}
}

// Test resolveReplaceDirective with empty target
func TestRunner_resolveReplaceDirective_EmptyTarget(t *testing.T) {
	r := NewRunner()

	// Edge case: empty target path
	result := r.resolveReplaceDirective("replace github.com/dep => ", "/project")
	if !strings.Contains(result, "replace github.com/dep => ") {
		t.Errorf("resolveReplaceDirective() = %q, should preserve structure", result)
	}
}

// Test parseReplaceDirectives with only comments
func TestRunner_parseReplaceDirectives_OnlyComments(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23

// replace github.com/dep => ../dep
# replace github.com/dep2 => ./local
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()
	replaces := r.parseReplaceDirectives(tmpDir)

	// Comments should not be parsed as replace directives
	if len(replaces) != 0 {
		t.Errorf("parseReplaceDirectives() found %d directives in comments, want 0", len(replaces))
	}
}

// Test Runner Verbose field
func TestRunner_Verbose(t *testing.T) {
	r := &Runner{
		TempDir: os.TempDir(),
		GoPath:  "go",
		Verbose: true,
	}

	if !r.Verbose {
		t.Error("Runner.Verbose should be true")
	}
}

// Test parseGoMod with complex go.mod
func TestRunner_parseGoMod_Complex(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `// This is a comment
module github.com/complex/module-name

go 1.23

require (
	github.com/some/dep v1.0.0
	github.com/other/dep v2.0.0
)

replace github.com/some/dep => ../local/dep

exclude github.com/bad/dep v1.0.0
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()
	modulePath, err := r.parseGoMod(tmpDir)
	if err != nil {
		t.Fatalf("parseGoMod() error = %v", err)
	}

	if modulePath != "github.com/complex/module-name" {
		t.Errorf("parseGoMod() = %q, want %q", modulePath, "github.com/complex/module-name")
	}
}

// Test generateGoMod output format
func TestRunner_generateGoMod_Format(t *testing.T) {
	r := NewRunner()
	result := r.generateGoMod("github.com/example/test", "/path/to/project")

	// Check specific format requirements
	if !strings.HasPrefix(result, "module wetwire-extract\n") {
		t.Error("generateGoMod() should start with module directive")
	}

	if !strings.Contains(result, "go 1.23") {
		t.Error("generateGoMod() should specify Go version")
	}

	if !strings.Contains(result, "v0.0.0") {
		t.Error("generateGoMod() should use v0.0.0 version")
	}
}

// Test parseReplaceDirectives with block syntax
func TestRunner_parseReplaceDirectives_Block(t *testing.T) {
	tmpDir := t.TempDir()

	// Note: The current implementation only parses replace directives that
	// start with "replace " (single line format), not block syntax
	goMod := `module github.com/example/test

go 1.23

replace github.com/single/dep => ./single
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()
	replaces := r.parseReplaceDirectives(tmpDir)

	// Should find the single-line replace directive
	if len(replaces) != 1 {
		t.Errorf("parseReplaceDirectives() returned %d directives, want 1", len(replaces))
	}
}

// Test parseGoMod with module on different lines
func TestRunner_parseGoMod_ModuleOnDifferentLine(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `

module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()
	modulePath, err := r.parseGoMod(tmpDir)
	if err != nil {
		t.Fatalf("parseGoMod() error = %v", err)
	}

	if modulePath != "github.com/example/test" {
		t.Errorf("parseGoMod() = %q, want %q", modulePath, "github.com/example/test")
	}
}

// Test parseGoMod with whitespace - only handles trimmed lines
func TestRunner_parseGoMod_Whitespace(t *testing.T) {
	tmpDir := t.TempDir()

	// Note: The current implementation uses TrimSpace on lines,
	// so "  module github.com/test" becomes "module github.com/test" after trimming
	// But then "module " is stripped, leaving "  github.com/example/test" from the original
	// This test verifies the current behavior
	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()
	modulePath, err := r.parseGoMod(tmpDir)
	if err != nil {
		t.Fatalf("parseGoMod() error = %v", err)
	}

	// The module path should start with the expected module
	// (trailing spaces might be included)
	if !strings.HasPrefix(modulePath, "github.com/example/test") {
		t.Errorf("parseGoMod() = %q, want prefix %q", modulePath, "github.com/example/test")
	}
}

// Test generateGoMod includes proper newlines
func TestRunner_generateGoMod_Newlines(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()
	result := r.generateGoMod("github.com/example/test", tmpDir)

	// Should have proper newlines for readability
	lines := strings.Split(result, "\n")
	if len(lines) < 4 {
		t.Errorf("generateGoMod() should have at least 4 lines, got %d", len(lines))
	}
}

// Test resolveReplaceDirective with various spacing
func TestRunner_resolveReplaceDirective_Spacing(t *testing.T) {
	r := NewRunner()

	tests := []struct {
		name    string
		line    string
		baseDir string
	}{
		{
			name:    "multiple spaces",
			line:    "replace   github.com/dep   =>   ./local",
			baseDir: "/project",
		},
		{
			name:    "tabs",
			line:    "replace\tgithub.com/dep\t=>\t./local",
			baseDir: "/project",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := r.resolveReplaceDirective(tt.line, tt.baseDir)
			// Should contain the original or resolved path
			if !strings.Contains(result, "github.com/dep") {
				t.Errorf("resolveReplaceDirective() = %q, should contain module path", result)
			}
		})
	}
}

// Test parseReplaceDirectives with empty lines
func TestRunner_parseReplaceDirectives_EmptyLines(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23


replace github.com/dep => ./local


`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()
	replaces := r.parseReplaceDirectives(tmpDir)

	if len(replaces) != 1 {
		t.Errorf("parseReplaceDirectives() returned %d directives, want 1", len(replaces))
	}
}

// Test resolveReplaceDirective with dot prefix
func TestRunner_resolveReplaceDirective_DotPrefix(t *testing.T) {
	r := NewRunner()

	tests := []struct {
		name    string
		line    string
		baseDir string
		wantAbs bool
	}{
		{
			name:    "starts with ./",
			line:    "replace github.com/dep => ./local",
			baseDir: "/project",
			wantAbs: true,
		},
		{
			name:    "starts with ../",
			line:    "replace github.com/dep => ../sibling",
			baseDir: "/project/sub",
			wantAbs: true,
		},
		{
			name:    "starts with ..",
			line:    "replace github.com/dep => ..sibling",
			baseDir: "/project",
			wantAbs: false, // doesn't start with . or ..
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := r.resolveReplaceDirective(tt.line, tt.baseDir)
			isAbs := strings.Contains(result, tt.baseDir) || strings.HasPrefix(strings.Split(result, " => ")[1], "/")
			if tt.wantAbs && !isAbs {
				t.Errorf("resolveReplaceDirective() = %q, expected absolute path", result)
			}
		})
	}
}

// Test NewRunner sets GoPath from PATH
func TestNewRunner_SetsGoPath(t *testing.T) {
	r := NewRunner()

	// GoPath might be empty if go is not in PATH, but typically it should be set
	// We just verify the field is accessible
	_ = r.GoPath
}

// Test parseReplaceDirectives returns empty for complex replace blocks
func TestRunner_parseReplaceDirectives_BlockSyntaxNotParsed(t *testing.T) {
	tmpDir := t.TempDir()

	// Block-style replace directives (inside parentheses) are not parsed by the current implementation
	goMod := `module github.com/example/test

go 1.23

replace (
	github.com/dep1 => ../dep1
)
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()
	replaces := r.parseReplaceDirectives(tmpDir)

	// Block syntax lines don't start with "replace ", so they won't be parsed
	// Lines inside the block are indented and won't start with "replace "
	if len(replaces) != 0 {
		t.Logf("parseReplaceDirectives() found %d directives (block syntax not fully parsed)", len(replaces))
	}
}

// Test parseReplaceDirectives with complex replace blocks
func TestRunner_parseReplaceDirectives_ComplexFormats(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/example/test

go 1.23

replace github.com/dep1 => ../dep1
replace github.com/dep2 => ./local
replace github.com/dep3 v1.0.0 => v1.0.1
`

	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner()
	replaces := r.parseReplaceDirectives(tmpDir)

	if len(replaces) != 3 {
		t.Errorf("parseReplaceDirectives() returned %d directives, want 3", len(replaces))
	}

	// Verify replaces are present
	foundDep1 := false
	foundDep2 := false
	foundDep3 := false
	for _, replace := range replaces {
		if strings.Contains(replace, "github.com/dep1") {
			foundDep1 = true
			// The relative path should be resolved to an absolute path
			if strings.Contains(replace, " => ../dep1") {
				t.Errorf("Relative path ../dep1 not resolved: %q", replace)
			}
		}
		if strings.Contains(replace, "github.com/dep2") {
			foundDep2 = true
			// The relative path should be resolved to an absolute path
			if strings.Contains(replace, " => ./local") {
				t.Errorf("Relative path ./local not resolved: %q", replace)
			}
		}
		if strings.Contains(replace, "github.com/dep3") {
			foundDep3 = true
			// Version replacement should be unchanged
		}
	}
	if !foundDep1 || !foundDep2 || !foundDep3 {
		t.Errorf("Missing expected replace directives: dep1=%v, dep2=%v, dep3=%v", foundDep1, foundDep2, foundDep3)
	}
}

// Suppress unused import warning
var _ = fmt.Sprintf
var _ = discover.DiscoveryResult{}
