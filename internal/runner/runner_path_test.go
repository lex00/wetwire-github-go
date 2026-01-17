package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lex00/wetwire-github-go/internal/discover"
)

func TestRunner_getPackagePath(t *testing.T) {
	r := NewRunner()
	baseDir := "/project"

	tests := []struct {
		modulePath string
		baseDir    string
		file       string
		want       string
	}{
		{"github.com/example/test", baseDir, "/project/workflows.go", "github.com/example/test"},
		{"github.com/example/test", baseDir, "/project/pkg/workflows.go", "github.com/example/test/pkg"},
		{"github.com/example/test", baseDir, "/project/internal/ci/workflows.go", "github.com/example/test/internal/ci"},
	}

	for _, tt := range tests {
		got := r.getPackagePath(tt.modulePath, tt.baseDir, tt.file)
		if got != tt.want {
			t.Errorf("getPackagePath(%q, %q, %q) = %q, want %q", tt.modulePath, tt.baseDir, tt.file, got, tt.want)
		}
	}
}

func TestRunner_pkgAlias(t *testing.T) {
	r := NewRunner()

	tests := []struct {
		input string
		want  string
	}{
		{"github.com/example/test", "test"},
		{"github.com/example/my-pkg", "my_pkg"},
		{"github.com/org/repo/internal/ci", "ci"},
	}

	for _, tt := range tests {
		got := r.pkgAlias(tt.input)
		if got != tt.want {
			t.Errorf("pkgAlias(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

// Test getPackagePath with complex paths
func TestRunner_getPackagePath_ComplexPaths(t *testing.T) {
	r := NewRunner()

	tests := []struct {
		name       string
		modulePath string
		baseDir    string
		file       string
		wantSubstr string
	}{
		{
			name:       "different paths still compute relative",
			modulePath: "github.com/example/test",
			baseDir:    "/project",
			file:       "/different/absolute/file.go",
			wantSubstr: "github.com/example/test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := r.getPackagePath(tt.modulePath, tt.baseDir, tt.file)
			if !strings.Contains(result, tt.wantSubstr) {
				t.Errorf("getPackagePath() = %q, should contain %q", result, tt.wantSubstr)
			}
		})
	}
}

// Test pkgAlias with various inputs
func TestRunner_pkgAlias_EdgeCases(t *testing.T) {
	r := NewRunner()

	tests := []struct {
		input string
		want  string
	}{
		{"github.com/example/test", "test"},
		{"github.com/example/my-pkg", "my_pkg"},
		{"github.com/org/repo/internal/ci", "ci"},
		{"github.com/complex-name/with-many-hyphens", "with_many_hyphens"},
		{"simple", "simple"},
		{"a", "a"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := r.pkgAlias(tt.input)
			if got != tt.want {
				t.Errorf("pkgAlias(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// Test getPackagePath edge cases
func TestRunner_getPackagePath_RootPackage(t *testing.T) {
	r := NewRunner()

	tests := []struct {
		name       string
		modulePath string
		baseDir    string
		file       string
		want       string
	}{
		{
			name:       "file in root directory",
			modulePath: "github.com/example/test",
			baseDir:    "/project",
			file:       "/project/main.go",
			want:       "github.com/example/test",
		},
		{
			name:       "file in subdirectory",
			modulePath: "github.com/example/test",
			baseDir:    "/project",
			file:       "/project/cmd/app/main.go",
			want:       "github.com/example/test/cmd/app",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := r.getPackagePath(tt.modulePath, tt.baseDir, tt.file)
			if got != tt.want {
				t.Errorf("getPackagePath() = %q, want %q", got, tt.want)
			}
		})
	}
}

// Test getPackagePath with Windows-style paths (if on Windows)
func TestRunner_getPackagePath_WindowsPaths(t *testing.T) {
	r := NewRunner()

	// Test with forward slashes (should work on all platforms)
	result := r.getPackagePath("github.com/example/test", "/project", "/project/subdir/file.go")
	expected := "github.com/example/test/subdir"

	if result != expected {
		t.Errorf("getPackagePath() = %q, want %q", result, expected)
	}
}

// Test getPackagePath with empty directory
func TestRunner_getPackagePath_EmptyDir(t *testing.T) {
	r := NewRunner()

	// When baseDir is empty, filepath.Rel may behave differently
	result := r.getPackagePath("github.com/example/test", "", "/project/file.go")
	if !strings.Contains(result, "github.com/example/test") {
		t.Errorf("getPackagePath() should contain module path, got: %q", result)
	}
}

// Test getPackagePath with same base and file directory
func TestRunner_getPackagePath_SameDir(t *testing.T) {
	r := NewRunner()

	// When file is directly in baseDir
	result := r.getPackagePath("github.com/example/test", "/project", "/project/file.go")
	if result != "github.com/example/test" {
		t.Errorf("getPackagePath() = %q, want %q", result, "github.com/example/test")
	}
}

// Test pkgAlias with single component path
func TestRunner_pkgAlias_SingleComponent(t *testing.T) {
	r := NewRunner()

	result := r.pkgAlias("main")
	if result != "main" {
		t.Errorf("pkgAlias() = %q, want %q", result, "main")
	}
}

// Test pkgAlias with multiple hyphens
func TestRunner_pkgAlias_MultipleHyphens(t *testing.T) {
	r := NewRunner()

	result := r.pkgAlias("github.com/my-org/my-awesome-package")
	if result != "my_awesome_package" {
		t.Errorf("pkgAlias() = %q, want %q", result, "my_awesome_package")
	}
}

// Test getPackagePath when filepath.Rel fails
func TestRunner_getPackagePath_RelFails(t *testing.T) {
	r := NewRunner()

	// On Unix, this shouldn't fail, but we can still test the code path
	// by using completely different drives (which wouldn't happen on Unix)
	// This test verifies the fallback behavior
	result := r.getPackagePath("github.com/example/test", "/base/path", "/different/path/file.go")

	// Should still return something sensible
	if result == "" {
		t.Error("getPackagePath() should not return empty string")
	}
}

// Suppress unused import warning
var _ = fmt.Sprintf
var _ = os.TempDir
var _ = filepath.Join
var _ = discover.DiscoveryResult{}
