// Package linter provides Go code quality rules for wetwire workflow declarations.
package linter

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// LintIssue represents a single lint issue found in Go code.
type LintIssue struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Severity string `json:"severity"` // "error", "warning", "info"
	Message  string `json:"message"`
	Rule     string `json:"rule"`
	Fixable  bool   `json:"fixable"`
}

// Rule is the interface that all linter rules must implement.
type Rule interface {
	// ID returns the unique identifier for this rule (e.g., "WAG001")
	ID() string
	// Description returns a human-readable description of the rule
	Description() string
	// Check analyzes the AST and returns any issues found
	Check(fset *token.FileSet, file *ast.File, path string) []LintIssue
}

// Linter runs rules against Go source code.
type Linter struct {
	rules []Rule
	fset  *token.FileSet
}

// NewLinter creates a new Linter with the specified rules.
func NewLinter(rules ...Rule) *Linter {
	return &Linter{
		rules: rules,
		fset:  token.NewFileSet(),
	}
}

// DefaultLinter creates a linter with all default rules enabled.
func DefaultLinter() *Linter {
	return NewLinter(
		&WAG001{},
		&WAG002{},
		&WAG003{},
		&WAG004{},
		&WAG005{},
		&WAG006{},
		&WAG007{MaxJobs: 10},
		&WAG008{},
	)
}

// LintResult contains the result of linting.
type LintResult struct {
	Success bool        `json:"success"`
	Issues  []LintIssue `json:"issues,omitempty"`
}

// LintFile lints a single Go file.
func (l *Linter) LintFile(path string) (*LintResult, error) {
	file, err := parser.ParseFile(l.fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	return l.lintAST(file, path), nil
}

// LintContent lints Go source code from memory.
func (l *Linter) LintContent(path string, content []byte) (*LintResult, error) {
	file, err := parser.ParseFile(l.fset, path, content, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	return l.lintAST(file, path), nil
}

// LintDir lints all Go files in a directory recursively.
func (l *Linter) LintDir(dir string) (*LintResult, error) {
	result := &LintResult{
		Success: true,
		Issues:  []LintIssue{},
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden directories and vendor
		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") || name == "vendor" || name == "testdata" {
				return filepath.SkipDir
			}
			return nil
		}

		// Only process .go files, skip test files
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		fileResult, err := l.LintFile(path)
		if err != nil {
			// Record the error but continue
			result.Issues = append(result.Issues, LintIssue{
				File:     path,
				Line:     1,
				Column:   1,
				Severity: "error",
				Message:  err.Error(),
				Rule:     "parse-error",
			})
			return nil
		}

		if !fileResult.Success {
			result.Success = false
		}
		result.Issues = append(result.Issues, fileResult.Issues...)

		return nil
	})

	return result, err
}

// lintAST runs all rules against a parsed AST.
func (l *Linter) lintAST(file *ast.File, path string) *LintResult {
	result := &LintResult{
		Success: true,
		Issues:  []LintIssue{},
	}

	for _, rule := range l.rules {
		issues := rule.Check(l.fset, file, path)
		if len(issues) > 0 {
			result.Success = false
			result.Issues = append(result.Issues, issues...)
		}
	}

	return result
}

// AddRule adds a rule to the linter.
func (l *Linter) AddRule(rule Rule) {
	l.rules = append(l.rules, rule)
}

// Rules returns the list of rules configured in the linter.
func (l *Linter) Rules() []Rule {
	return l.rules
}
