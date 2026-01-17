package lint

import (
	"os"
	"path/filepath"
	"testing"
)

// Realistic Go workflow code for benchmarking
var benchmarkGoCode = []byte(`package main

import (
	"github.com/lex00/wetwire-github-go/workflow"
	"github.com/lex00/wetwire-github-go/actions/checkout"
	"github.com/lex00/wetwire-github-go/actions/setup_go"
	"github.com/lex00/wetwire-github-go/actions/cache"
)

// CI is the main CI workflow
var CI = workflow.Workflow{
	Name: "CI",
	On:   CITriggers,
	Jobs: map[string]workflow.Job{
		"build": Build,
		"test":  Test,
		"lint":  Lint,
	},
}

var CITriggers = workflow.Triggers{
	Push:        CIPush,
	PullRequest: CIPullRequest,
}

var CIPush = workflow.PushTrigger{
	Branches: []string{"main", "develop"},
}

var CIPullRequest = workflow.PullRequestTrigger{
	Branches: []string{"main"},
}

var Build = workflow.Job{
	Name:           "build",
	RunsOn:         "ubuntu-latest",
	TimeoutMinutes: 30,
	Steps:          BuildSteps,
}

var BuildSteps = []any{
	checkout.Checkout{},
	setup_go.SetupGo{GoVersion: "1.23"},
	cache.Cache{
		Path: "~/go/pkg/mod",
		Key:  "go-mod-${{ hashFiles('go.sum') }}",
	},
	workflow.Step{Run: "go build ./..."},
}

var Test = workflow.Job{
	Name:           "test",
	RunsOn:         "ubuntu-latest",
	TimeoutMinutes: 30,
	Needs:          []any{Build},
	Steps:          TestSteps,
}

var TestSteps = []any{
	checkout.Checkout{},
	setup_go.SetupGo{GoVersion: "1.23"},
	workflow.Step{Run: "go test -v ./..."},
}

var Lint = workflow.Job{
	Name:           "lint",
	RunsOn:         "ubuntu-latest",
	TimeoutMinutes: 15,
	Steps:          LintSteps,
}

var LintSteps = []any{
	checkout.Checkout{},
	workflow.Step{
		Name: "Run golangci-lint",
		Uses: "golangci/golangci-lint-action@v4",
	},
}
`)

// Code with various lint issues for benchmarking
var benchmarkCodeWithIssues = []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

// Has WAG001 issue: raw uses string
var CheckoutStep = workflow.Step{Uses: "actions/checkout@v4"}

// Has WAG002 issue: raw expression string
var ConditionalStep = workflow.Step{
	If: "${{ github.ref == 'refs/heads/main' }}",
}

// Has WAG003 issue: hardcoded secret
var token = "ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

// Has WAG004 issue: inline matrix
var MatrixJob = workflow.Job{
	Strategy: workflow.Strategy{
		Matrix: workflow.Matrix{
			Values: map[string][]any{
				"os": {"ubuntu-latest", "macos-latest"},
			},
		},
	},
}

// Has WAG006 issue: duplicate workflow names
var CI1 = workflow.Workflow{Name: "CI"}
var CI2 = workflow.Workflow{Name: "CI"}

// Has WAG012 issue: deprecated action version
var OldCheckout = workflow.Step{Uses: "actions/checkout@v2"}

// Has WAG013 issue: pointer assignment
var PointerWorkflow = &workflow.Workflow{
	Name: "Pointer",
}

// Has WAG014 issue: missing timeout
var NoTimeoutJob = workflow.Job{
	Name:   "no-timeout",
	RunsOn: "ubuntu-latest",
}
`)

// Large file with many declarations for benchmarking
func generateLargeCode() []byte {
	code := `package main

import "github.com/lex00/wetwire-github-go/workflow"

`
	for i := 0; i < 50; i++ {
		code += `var Workflow` + string(rune('A'+i%26)) + string(rune('0'+i/26)) + ` = workflow.Workflow{
	Name: "Workflow ` + string(rune('A'+i%26)) + `",
	Jobs: map[string]workflow.Job{
		"build": Job` + string(rune('A'+i%26)) + string(rune('0'+i/26)) + `,
	},
}

var Job` + string(rune('A'+i%26)) + string(rune('0'+i/26)) + ` = workflow.Job{
	Name:           "job-` + string(rune('a'+i%26)) + `",
	RunsOn:         "ubuntu-latest",
	TimeoutMinutes: 30,
	Steps: []any{
		workflow.Step{Run: "echo hello"},
	},
}

`
	}
	return []byte(code)
}

// BenchmarkLintFile benchmarks linting a single file.
func BenchmarkLintFile(b *testing.B) {
	tmpDir := b.TempDir()
	testFile := filepath.Join(tmpDir, "workflows.go")
	if err := os.WriteFile(testFile, benchmarkGoCode, 0644); err != nil {
		b.Fatal(err)
	}

	l := DefaultLinter()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result, err := l.LintFile(testFile)
		if err != nil {
			b.Fatal(err)
		}
		_ = result
	}
}

// BenchmarkAllRules benchmarks all 16 rules against a file.
func BenchmarkAllRules(b *testing.B) {
	l := DefaultLinter()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result, err := l.LintContent("test.go", benchmarkGoCode)
		if err != nil {
			b.Fatal(err)
		}
		_ = result
	}
}

// BenchmarkLintContent benchmarks linting from memory.
func BenchmarkLintContent(b *testing.B) {
	l := DefaultLinter()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result, err := l.LintContent("test.go", benchmarkGoCode)
		if err != nil {
			b.Fatal(err)
		}
		_ = result
	}
}

// BenchmarkLintCodeWithIssues benchmarks linting code with known issues.
func BenchmarkLintCodeWithIssues(b *testing.B) {
	l := DefaultLinter()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result, err := l.LintContent("test.go", benchmarkCodeWithIssues)
		if err != nil {
			b.Fatal(err)
		}
		if result.Success {
			b.Fatal("expected issues to be found")
		}
	}
}

// BenchmarkLintLargeFile benchmarks linting a large file.
func BenchmarkLintLargeFile(b *testing.B) {
	largeCode := generateLargeCode()
	l := DefaultLinter()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result, err := l.LintContent("test.go", largeCode)
		if err != nil {
			b.Fatal(err)
		}
		_ = result
	}
}

// BenchmarkLintDir benchmarks linting a directory.
func BenchmarkLintDir(b *testing.B) {
	tmpDir := b.TempDir()

	// Create multiple files
	for i := 0; i < 5; i++ {
		filename := filepath.Join(tmpDir, "workflow"+string(rune('a'+i))+".go")
		if err := os.WriteFile(filename, benchmarkGoCode, 0644); err != nil {
			b.Fatal(err)
		}
	}

	l := DefaultLinter()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result, err := l.LintDir(tmpDir)
		if err != nil {
			b.Fatal(err)
		}
		_ = result
	}
}

// BenchmarkWAG001 benchmarks the WAG001 rule specifically.
func BenchmarkWAG001(b *testing.B) {
	l := NewLinter(&WAG001{})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result, err := l.LintContent("test.go", benchmarkCodeWithIssues)
		if err != nil {
			b.Fatal(err)
		}
		_ = result
	}
}

// BenchmarkWAG002 benchmarks the WAG002 rule specifically.
func BenchmarkWAG002(b *testing.B) {
	l := NewLinter(&WAG002{})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result, err := l.LintContent("test.go", benchmarkCodeWithIssues)
		if err != nil {
			b.Fatal(err)
		}
		_ = result
	}
}

// BenchmarkWAG003 benchmarks the WAG003 rule (secret detection).
func BenchmarkWAG003(b *testing.B) {
	l := NewLinter(&WAG003{})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result, err := l.LintContent("test.go", benchmarkCodeWithIssues)
		if err != nil {
			b.Fatal(err)
		}
		_ = result
	}
}

// BenchmarkWAG011 benchmarks the WAG011 rule (undefined dependencies).
func BenchmarkWAG011(b *testing.B) {
	code := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var Build = workflow.Job{
	Name:   "build",
	RunsOn: "ubuntu-latest",
}

var Test = workflow.Job{
	Name:   "test",
	RunsOn: "ubuntu-latest",
	Needs:  []any{Build},
}

var Deploy = workflow.Job{
	Name:   "deploy",
	RunsOn: "ubuntu-latest",
	Needs:  []any{Build, Test},
}
`)

	l := NewLinter(&WAG011{})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result, err := l.LintContent("test.go", code)
		if err != nil {
			b.Fatal(err)
		}
		_ = result
	}
}

// BenchmarkNewLinter benchmarks linter creation.
func BenchmarkNewLinter(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l := NewLinter(&WAG001{}, &WAG002{}, &WAG003{})
		if l == nil {
			b.Fatal("expected linter")
		}
	}
}

// BenchmarkDefaultLinter benchmarks default linter creation.
func BenchmarkDefaultLinter(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l := DefaultLinter()
		if l == nil {
			b.Fatal("expected linter")
		}
	}
}

// BenchmarkFix benchmarks fixing issues.
func BenchmarkFix(b *testing.B) {
	code := []byte(`package main

import "github.com/lex00/wetwire-github-go/workflow"

var CheckoutStep = workflow.Step{Uses: "actions/checkout@v4"}
var SetupGoStep = workflow.Step{Uses: "actions/setup-go@v5"}
var CacheStep = workflow.Step{Uses: "actions/cache@v4"}
`)

	l := NewLinter(&WAG001{})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result, err := l.Fix("test.go", code)
		if err != nil {
			b.Fatal(err)
		}
		_ = result
	}
}

// BenchmarkSingleRule benchmarks a single rule to measure rule overhead.
func BenchmarkSingleRule(b *testing.B) {
	l := NewLinter(&WAG006{})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result, err := l.LintContent("test.go", benchmarkGoCode)
		if err != nil {
			b.Fatal(err)
		}
		_ = result
	}
}

// BenchmarkAllRulesCleanCode benchmarks all rules against clean code.
func BenchmarkAllRulesCleanCode(b *testing.B) {
	// Clean code that should pass all rules
	cleanCode := []byte(`package main

import (
	"github.com/lex00/wetwire-github-go/workflow"
	"github.com/lex00/wetwire-github-go/actions/checkout"
	"github.com/lex00/wetwire-github-go/actions/setup_go"
)

var CI = workflow.Workflow{
	Name: "CI",
	Jobs: map[string]workflow.Job{
		"build": Build,
	},
}

var Build = workflow.Job{
	Name:           "build",
	RunsOn:         "ubuntu-latest",
	TimeoutMinutes: 30,
	Steps:          BuildSteps,
}

var BuildSteps = []any{
	checkout.Checkout{},
	setup_go.SetupGo{GoVersion: "1.23"},
	workflow.Step{Run: "go build ./..."},
}
`)

	l := DefaultLinter()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result, err := l.LintContent("test.go", cleanCode)
		if err != nil {
			b.Fatal(err)
		}
		_ = result
	}
}
