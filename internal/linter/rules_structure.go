package linter

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

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
			if innerLit, ok := kv.Value.(*ast.CompositeLit); ok {
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

		return false
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
		if strings.HasPrefix(val, "${{") && strings.HasSuffix(val, "}}") {
			expContent := strings.TrimPrefix(strings.TrimSuffix(val, "}}"), "${{")
			expContent = strings.TrimSpace(expContent)

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

// WAG013 checks for pointer assignments (&Type{}) in workflow declarations.
type WAG013 struct{}

func (r *WAG013) ID() string          { return "WAG013" }
func (r *WAG013) Description() string { return "Avoid pointer assignments (&Type{}) - use value types" }

func (r *WAG013) Check(fset *token.FileSet, file *ast.File, path string) []LintIssue {
	var issues []LintIssue

	ast.Inspect(file, func(n ast.Node) bool {
		unary, ok := n.(*ast.UnaryExpr)
		if !ok || unary.Op != token.AND {
			return true
		}

		comp, ok := unary.X.(*ast.CompositeLit)
		if !ok {
			return true
		}

		typeName := getTypeName(comp.Type)
		if isWorkflowType(typeName) {
			pos := fset.Position(unary.Pos())
			issues = append(issues, LintIssue{
				File:     path,
				Line:     pos.Line,
				Column:   pos.Column,
				Severity: "error",
				Message:  fmt.Sprintf("Avoid pointer assignment &%s{} - use value type instead", typeName),
				Rule:     r.ID(),
				Fixable:  false,
			})
		}
		return true
	})

	return issues
}
