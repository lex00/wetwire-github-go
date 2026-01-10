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

// Fixer is an optional interface that rules can implement to provide auto-fix capability.
type Fixer interface {
	// Fix attempts to fix an issue and returns the modified source code.
	// The issue parameter contains information about what to fix.
	// Returns the fixed source code, or nil if the fix could not be applied.
	Fix(fset *token.FileSet, file *ast.File, path string, src []byte, issue LintIssue) ([]byte, error)
}

// FixResult contains the result of a fix operation.
type FixResult struct {
	Content    []byte `json:"-"`
	FixedCount int    `json:"fixed_count"`
	Issues     []LintIssue `json:"issues,omitempty"` // Remaining unfixed issues
}

// FixDirResult contains the result of fixing a directory.
type FixDirResult struct {
	Files      []string `json:"files"`
	TotalFixed int      `json:"total_fixed"`
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
		&WAG009{},
		&WAG010{},
		&WAG011{},
		&WAG012{},
		&WAG013{},
		&WAG014{},
		&WAG015{},
		&WAG016{},
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

// Fix applies fixes to Go source code from memory.
func (l *Linter) Fix(path string, content []byte) (*FixResult, error) {
	// Parse the source
	file, err := parser.ParseFile(l.fset, path, content, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	result := &FixResult{
		Content:    content,
		FixedCount: 0,
		Issues:     []LintIssue{},
	}

	// Collect all issues first
	var allIssues []LintIssue
	for _, rule := range l.rules {
		issues := rule.Check(l.fset, file, path)
		allIssues = append(allIssues, issues...)
	}

	// Try to fix each issue
	currentContent := content
	for _, issue := range allIssues {
		if !issue.Fixable {
			result.Issues = append(result.Issues, issue)
			continue
		}

		// Find the rule that generated this issue
		for _, rule := range l.rules {
			if rule.ID() != issue.Rule {
				continue
			}

			// Check if rule implements Fixer
			fixer, ok := rule.(Fixer)
			if !ok {
				result.Issues = append(result.Issues, issue)
				continue
			}

			// Re-parse with current content
			currentFile, err := parser.ParseFile(l.fset, path, currentContent, parser.ParseComments)
			if err != nil {
				result.Issues = append(result.Issues, issue)
				continue
			}

			// Apply fix
			fixed, err := fixer.Fix(l.fset, currentFile, path, currentContent, issue)
			if err != nil || fixed == nil {
				result.Issues = append(result.Issues, issue)
				continue
			}

			currentContent = fixed
			result.FixedCount++
			break
		}
	}

	result.Content = currentContent
	return result, nil
}

// FixFile applies fixes to a Go file and writes the result back.
func (l *Linter) FixFile(path string) (*FixResult, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	result, err := l.Fix(path, content)
	if err != nil {
		return nil, err
	}

	// Write back if changes were made
	if result.FixedCount > 0 {
		if err := os.WriteFile(path, result.Content, 0644); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// FixDir applies fixes to all Go files in a directory.
func (l *Linter) FixDir(dir string) (*FixDirResult, error) {
	result := &FixDirResult{
		Files:      []string{},
		TotalFixed: 0,
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-Go files
		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") || name == "vendor" || name == "testdata" {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		fileResult, err := l.FixFile(path)
		if err != nil {
			// Continue on error
			return nil
		}

		if fileResult.FixedCount > 0 {
			result.Files = append(result.Files, path)
			result.TotalFixed += fileResult.FixedCount
		}

		return nil
	})

	return result, err
}
