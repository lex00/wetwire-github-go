package linter

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

// WAG001 checks for raw uses: strings instead of typed action wrappers.
type WAG001 struct{}

func (r *WAG001) ID() string          { return "WAG001" }
func (r *WAG001) Description() string { return "Use typed action wrappers instead of raw uses: strings" }

func (r *WAG001) Check(fset *token.FileSet, file *ast.File, path string) []LintIssue {
	var issues []LintIssue

	ast.Inspect(file, func(n ast.Node) bool {
		lit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}

		// Check if this is a workflow.Step
		typeName := getTypeName(lit.Type)
		if typeName != "workflow.Step" && typeName != "Step" {
			return true
		}

		// Look for Uses field with a string literal
		for _, elt := range lit.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			key, ok := kv.Key.(*ast.Ident)
			if !ok || key.Name != "Uses" {
				continue
			}
			if _, ok := kv.Value.(*ast.BasicLit); ok {
				pos := fset.Position(kv.Pos())
				issues = append(issues, LintIssue{
					File:     path,
					Line:     pos.Line,
					Column:   pos.Column,
					Severity: "warning",
					Message:  "Use typed action wrapper instead of raw uses: string",
					Rule:     r.ID(),
					Fixable:  false,
				})
			}
		}
		return true
	})

	return issues
}

// WAG002 checks for raw expression strings instead of condition builders.
type WAG002 struct{}

func (r *WAG002) ID() string          { return "WAG002" }
func (r *WAG002) Description() string { return "Use condition builders instead of raw expression strings" }

func (r *WAG002) Check(fset *token.FileSet, file *ast.File, path string) []LintIssue {
	var issues []LintIssue

	ast.Inspect(file, func(n ast.Node) bool {
		lit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}

		// Look for If field with expression string
		for _, elt := range lit.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			key, ok := kv.Key.(*ast.Ident)
			if !ok || key.Name != "If" {
				continue
			}
			if bl, ok := kv.Value.(*ast.BasicLit); ok && bl.Kind == token.STRING {
				val := strings.Trim(bl.Value, `"'`)
				if strings.Contains(val, "${{") {
					pos := fset.Position(kv.Pos())
					issues = append(issues, LintIssue{
						File:     path,
						Line:     pos.Line,
						Column:   pos.Column,
						Severity: "warning",
						Message:  "Use condition builder instead of raw expression string",
						Rule:     r.ID(),
						Fixable:  false,
					})
				}
			}
		}
		return true
	})

	return issues
}

// WAG003 checks for hardcoded secrets instead of using the secrets context.
type WAG003 struct{}

func (r *WAG003) ID() string          { return "WAG003" }
func (r *WAG003) Description() string { return "Use secrets context instead of hardcoded strings" }

func (r *WAG003) Check(fset *token.FileSet, file *ast.File, path string) []LintIssue {
	var issues []LintIssue

	secretPatterns := []string{
		"ghp_", // GitHub personal access token
		"ghs_", // GitHub server token
		"ghu_", // GitHub user token
		"ghr_", // GitHub refresh token
		"github_pat_",
	}

	ast.Inspect(file, func(n ast.Node) bool {
		bl, ok := n.(*ast.BasicLit)
		if !ok || bl.Kind != token.STRING {
			return true
		}

		val := strings.Trim(bl.Value, `"'`)
		for _, pattern := range secretPatterns {
			if strings.Contains(val, pattern) {
				pos := fset.Position(bl.Pos())
				issues = append(issues, LintIssue{
					File:     path,
					Line:     pos.Line,
					Column:   pos.Column,
					Severity: "error",
					Message:  "Hardcoded secret detected - use secrets context",
					Rule:     r.ID(),
					Fixable:  false,
				})
				break
			}
		}
		return true
	})

	return issues
}

// WAG004 checks for inline matrix maps instead of using the matrix builder.
type WAG004 struct{}

func (r *WAG004) ID() string          { return "WAG004" }
func (r *WAG004) Description() string { return "Use matrix builder instead of inline maps" }

func (r *WAG004) Check(fset *token.FileSet, file *ast.File, path string) []LintIssue {
	var issues []LintIssue

	ast.Inspect(file, func(n ast.Node) bool {
		lit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}

		// Check if this is a Strategy field with inline matrix
		typeName := getTypeName(lit.Type)
		if typeName != "workflow.Strategy" && typeName != "Strategy" {
			return true
		}

		for _, elt := range lit.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			key, ok := kv.Key.(*ast.Ident)
			if !ok || key.Name != "Matrix" {
				continue
			}
			// Check if Matrix is an inline CompositeLit
			if innerLit, ok := kv.Value.(*ast.CompositeLit); ok {
				// Check if it's defining Values inline
				for _, innerElt := range innerLit.Elts {
					if innerKV, ok := innerElt.(*ast.KeyValueExpr); ok {
						if innerKey, ok := innerKV.Key.(*ast.Ident); ok && innerKey.Name == "Values" {
							pos := fset.Position(kv.Pos())
							issues = append(issues, LintIssue{
								File:     path,
								Line:     pos.Line,
								Column:   pos.Column,
								Severity: "info",
								Message:  "Consider extracting matrix to a named variable",
								Rule:     r.ID(),
								Fixable:  false,
							})
						}
					}
				}
			}
		}
		return true
	})

	return issues
}

// WAG005 checks for inline struct definitions that should be extracted.
type WAG005 struct{}

func (r *WAG005) ID() string          { return "WAG005" }
func (r *WAG005) Description() string { return "Extract inline structs to named variables" }

func (r *WAG005) Check(fset *token.FileSet, file *ast.File, path string) []LintIssue {
	var issues []LintIssue
	nestingDepth := 0
	const maxNesting = 2

	var checkNesting func(n ast.Node) bool
	checkNesting = func(n ast.Node) bool {
		lit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}

		nestingDepth++
		defer func() { nestingDepth-- }()

		if nestingDepth > maxNesting {
			pos := fset.Position(lit.Pos())
			issues = append(issues, LintIssue{
				File:     path,
				Line:     pos.Line,
				Column:   pos.Column,
				Severity: "info",
				Message:  fmt.Sprintf("Deeply nested struct (depth %d) - consider extracting to named variable", nestingDepth),
				Rule:     r.ID(),
				Fixable:  false,
			})
		}

		for _, elt := range lit.Elts {
			if kv, ok := elt.(*ast.KeyValueExpr); ok {
				ast.Inspect(kv.Value, checkNesting)
			}
		}

		return false // Don't recurse further, we handle it manually
	}

	ast.Inspect(file, checkNesting)
	return issues
}

// WAG006 checks for duplicate workflow names.
type WAG006 struct{}

func (r *WAG006) ID() string          { return "WAG006" }
func (r *WAG006) Description() string { return "Detect duplicate workflow names" }

func (r *WAG006) Check(fset *token.FileSet, file *ast.File, path string) []LintIssue {
	var issues []LintIssue
	workflowNames := make(map[string]token.Pos)

	ast.Inspect(file, func(n ast.Node) bool {
		lit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}

		typeName := getTypeName(lit.Type)
		if typeName != "workflow.Workflow" && typeName != "Workflow" {
			return true
		}

		// Extract Name field
		for _, elt := range lit.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			key, ok := kv.Key.(*ast.Ident)
			if !ok || key.Name != "Name" {
				continue
			}
			if bl, ok := kv.Value.(*ast.BasicLit); ok && bl.Kind == token.STRING {
				name := strings.Trim(bl.Value, `"'`)
				if prevPos, exists := workflowNames[name]; exists {
					pos := fset.Position(kv.Pos())
					prevPosInfo := fset.Position(prevPos)
					issues = append(issues, LintIssue{
						File:     path,
						Line:     pos.Line,
						Column:   pos.Column,
						Severity: "error",
						Message:  fmt.Sprintf("Duplicate workflow name %q (first defined at line %d)", name, prevPosInfo.Line),
						Rule:     r.ID(),
						Fixable:  false,
					})
				} else {
					workflowNames[name] = kv.Pos()
				}
			}
		}
		return true
	})

	return issues
}

// WAG007 checks for files with too many jobs.
type WAG007 struct {
	MaxJobs int
}

func (r *WAG007) ID() string          { return "WAG007" }
func (r *WAG007) Description() string { return "Flag oversized files (>N jobs)" }

func (r *WAG007) Check(fset *token.FileSet, file *ast.File, path string) []LintIssue {
	var issues []LintIssue
	jobCount := 0

	maxJobs := r.MaxJobs
	if maxJobs == 0 {
		maxJobs = 10
	}

	ast.Inspect(file, func(n ast.Node) bool {
		lit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}

		typeName := getTypeName(lit.Type)
		if typeName == "workflow.Job" || typeName == "Job" {
			jobCount++
		}
		return true
	})

	if jobCount > maxJobs {
		issues = append(issues, LintIssue{
			File:     path,
			Line:     1,
			Column:   1,
			Severity: "warning",
			Message:  fmt.Sprintf("File contains %d jobs (max recommended: %d)", jobCount, maxJobs),
			Rule:     r.ID(),
			Fixable:  false,
		})
	}

	return issues
}

// WAG008 checks for hardcoded expression strings.
type WAG008 struct{}

func (r *WAG008) ID() string          { return "WAG008" }
func (r *WAG008) Description() string { return "Avoid hardcoded expression strings" }

func (r *WAG008) Check(fset *token.FileSet, file *ast.File, path string) []LintIssue {
	var issues []LintIssue

	ast.Inspect(file, func(n ast.Node) bool {
		bl, ok := n.(*ast.BasicLit)
		if !ok || bl.Kind != token.STRING {
			return true
		}

		val := strings.Trim(bl.Value, `"'`)
		// Check for GitHub expression syntax
		if strings.HasPrefix(val, "${{") && strings.HasSuffix(val, "}}") {
			// Check for common expressions that should use builders
			expContent := strings.TrimPrefix(strings.TrimSuffix(val, "}}"), "${{")
			expContent = strings.TrimSpace(expContent)

			// Skip simple variable references like ${{ github.token }}
			if strings.HasPrefix(expContent, "github.") ||
				strings.HasPrefix(expContent, "secrets.") ||
				strings.HasPrefix(expContent, "matrix.") ||
				strings.HasPrefix(expContent, "steps.") ||
				strings.HasPrefix(expContent, "needs.") ||
				strings.HasPrefix(expContent, "inputs.") ||
				strings.HasPrefix(expContent, "env.") {
				return true
			}

			pos := fset.Position(bl.Pos())
			issues = append(issues, LintIssue{
				File:     path,
				Line:     pos.Line,
				Column:   pos.Column,
				Severity: "info",
				Message:  "Consider using expression builder instead of hardcoded expression",
				Rule:     r.ID(),
				Fixable:  false,
			})
		}
		return true
	})

	return issues
}

// Helper function to get type name from AST expression
func getTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		if x, ok := t.X.(*ast.Ident); ok {
			return x.Name + "." + t.Sel.Name
		}
	}
	return ""
}
