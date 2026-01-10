package runner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lex00/wetwire-github-go/internal/discover"
)

// createBenchmarkProject creates a realistic project structure for benchmarking.
func createBenchmarkProject(b *testing.B) string {
	b.Helper()
	tmpDir := b.TempDir()

	// Create go.mod
	goMod := `module github.com/example/benchmark

go 1.23

require github.com/lex00/wetwire-github-go v0.0.0
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		b.Fatal(err)
	}

	// Create workflow file
	workflowContent := `package main

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var CI = workflow.Workflow{
	Name: "CI",
	Jobs: []any{Build, Test},
}

var Build = workflow.Job{
	Name:   "build",
	RunsOn: "ubuntu-latest",
}

var Test = workflow.Job{
	Name:   "test",
	RunsOn: "ubuntu-latest",
	Needs:  []any{Build},
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "workflows.go"), []byte(workflowContent), 0644); err != nil {
		b.Fatal(err)
	}

	return tmpDir
}

// BenchmarkExtractValues benchmarks value extraction from Go declarations.
// Note: This benchmark is skipped in short mode as it requires actual Go execution.
func BenchmarkExtractValues(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping benchmark in short mode")
	}

	// This benchmark is expensive because it actually runs Go code
	// In practice, we benchmark the setup and program generation instead

	r := NewRunner()
	discovered := &discover.DiscoveryResult{
		Workflows: []discover.DiscoveredWorkflow{},
		Jobs:      []discover.DiscoveredJob{},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result, err := r.ExtractValues(".", discovered)
		if err != nil {
			b.Fatal(err)
		}
		// Empty discovery result should return quickly
		_ = result
	}
}

// BenchmarkGenerateProgram benchmarks extraction program generation.
func BenchmarkGenerateProgram(b *testing.B) {
	r := NewRunner()
	baseDir := "/project"
	modulePath := "github.com/example/test"

	discovered := &discover.DiscoveryResult{
		Workflows: []discover.DiscoveredWorkflow{
			{Name: "CI", File: "/project/workflows.go", Line: 10, Jobs: []string{"Build", "Test"}},
			{Name: "Release", File: "/project/workflows.go", Line: 20, Jobs: []string{"Build", "Test", "Publish"}},
		},
		Jobs: []discover.DiscoveredJob{
			{Name: "Build", File: "/project/jobs.go", Line: 5, Dependencies: []string{}},
			{Name: "Test", File: "/project/jobs.go", Line: 15, Dependencies: []string{"Build"}},
			{Name: "Publish", File: "/project/jobs.go", Line: 25, Dependencies: []string{"Build", "Test"}},
			{Name: "Deploy", File: "/project/jobs.go", Line: 35, Dependencies: []string{"Publish"}},
		},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		program, err := r.generateProgram(modulePath, baseDir, discovered)
		if err != nil {
			b.Fatal(err)
		}
		if len(program) == 0 {
			b.Fatal("expected program output")
		}
	}
}

// BenchmarkGenerateProgramLarge benchmarks program generation with many resources.
func BenchmarkGenerateProgramLarge(b *testing.B) {
	r := NewRunner()
	baseDir := "/project"
	modulePath := "github.com/example/test"

	// Create a large discovery result
	workflows := make([]discover.DiscoveredWorkflow, 10)
	for i := 0; i < 10; i++ {
		workflows[i] = discover.DiscoveredWorkflow{
			Name: "Workflow" + string(rune('A'+i)),
			File: "/project/workflows.go",
			Line: i * 10,
			Jobs: []string{"Job1", "Job2", "Job3"},
		}
	}

	jobs := make([]discover.DiscoveredJob, 30)
	for i := 0; i < 30; i++ {
		jobs[i] = discover.DiscoveredJob{
			Name:         "Job" + string(rune('A'+i)),
			File:         "/project/jobs.go",
			Line:         i * 5,
			Dependencies: []string{},
		}
		if i > 0 {
			jobs[i].Dependencies = []string{"Job" + string(rune('A'+i-1))}
		}
	}

	discovered := &discover.DiscoveryResult{
		Workflows: workflows,
		Jobs:      jobs,
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		program, err := r.generateProgram(modulePath, baseDir, discovered)
		if err != nil {
			b.Fatal(err)
		}
		if len(program) == 0 {
			b.Fatal("expected program output")
		}
	}
}

// BenchmarkParseGoMod benchmarks go.mod parsing.
func BenchmarkParseGoMod(b *testing.B) {
	tmpDir := b.TempDir()

	goMod := `module github.com/example/test

go 1.23

require (
	github.com/some/dep v1.0.0
	github.com/another/dep v2.0.0
)

replace github.com/some/dep => ./local/dep
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		b.Fatal(err)
	}

	r := NewRunner()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		modulePath, err := r.parseGoMod(tmpDir)
		if err != nil {
			b.Fatal(err)
		}
		if modulePath == "" {
			b.Fatal("expected module path")
		}
	}
}

// BenchmarkGenerateGoMod benchmarks go.mod generation.
func BenchmarkGenerateGoMod(b *testing.B) {
	r := NewRunner()
	modulePath := "github.com/example/test"
	dir := "/project"

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		goMod := r.generateGoMod(modulePath, dir)
		if len(goMod) == 0 {
			b.Fatal("expected go.mod output")
		}
	}
}

// BenchmarkParseReplaceDirectives benchmarks parsing replace directives.
func BenchmarkParseReplaceDirectives(b *testing.B) {
	tmpDir := b.TempDir()

	goMod := `module github.com/example/test

go 1.23

replace github.com/dep1 => ./local/dep1
replace github.com/dep2 => ../dep2
replace github.com/dep3 => /absolute/path/dep3
replace github.com/dep4 v1.0.0 => v1.0.1
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		b.Fatal(err)
	}

	r := NewRunner()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		replaces := r.parseReplaceDirectives(tmpDir)
		if len(replaces) == 0 {
			b.Fatal("expected replace directives")
		}
	}
}

// BenchmarkGetPackagePath benchmarks package path computation.
func BenchmarkGetPackagePath(b *testing.B) {
	r := NewRunner()
	modulePath := "github.com/example/test"
	baseDir := "/project"

	files := []string{
		"/project/workflows.go",
		"/project/pkg/workflows.go",
		"/project/internal/ci/workflows.go",
		"/project/cmd/main/main.go",
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, file := range files {
			path := r.getPackagePath(modulePath, baseDir, file)
			if path == "" {
				b.Fatal("expected package path")
			}
		}
	}
}

// BenchmarkPkgAlias benchmarks package alias generation.
func BenchmarkPkgAlias(b *testing.B) {
	r := NewRunner()

	paths := []string{
		"github.com/example/test",
		"github.com/example/my-pkg",
		"github.com/org/repo/internal/ci",
		"github.com/org/repo/pkg/workflows",
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, path := range paths {
			alias := r.pkgAlias(path)
			if alias == "" {
				b.Fatal("expected alias")
			}
		}
	}
}

// BenchmarkResolveReplaceDirective benchmarks resolving replace directives.
func BenchmarkResolveReplaceDirective(b *testing.B) {
	r := NewRunner()
	baseDir := "/project"

	directives := []string{
		"replace github.com/dep => ../dep",
		"replace github.com/dep => ./local",
		"replace github.com/dep => /absolute/path",
		"replace github.com/dep v1.0.0 => v1.0.1",
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, directive := range directives {
			resolved := r.resolveReplaceDirective(directive, baseDir)
			if resolved == "" {
				b.Fatal("expected resolved directive")
			}
		}
	}
}

// BenchmarkNewRunner benchmarks runner creation.
func BenchmarkNewRunner(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r := NewRunner()
		if r == nil {
			b.Fatal("expected runner")
		}
	}
}
