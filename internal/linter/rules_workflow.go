package linter

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

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

// WAG009 validates matrix dimension values are not empty.
type WAG009 struct{}

func (r *WAG009) ID() string          { return "WAG009" }
func (r *WAG009) Description() string { return "Validate matrix dimension values" }

func (r *WAG009) Check(fset *token.FileSet, file *ast.File, path string) []LintIssue {
	var issues []LintIssue

	ast.Inspect(file, func(n ast.Node) bool {
		lit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}

		typeName := getTypeName(lit.Type)
		if typeName != "workflow.Matrix" && typeName != "Matrix" {
			return true
		}

		for _, elt := range lit.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			key, ok := kv.Key.(*ast.Ident)
			if !ok || key.Name != "Values" {
				continue
			}

			if mapLit, ok := kv.Value.(*ast.CompositeLit); ok {
				for _, mapElt := range mapLit.Elts {
					if mapKV, ok := mapElt.(*ast.KeyValueExpr); ok {
						if sliceLit, ok := mapKV.Value.(*ast.CompositeLit); ok {
							if len(sliceLit.Elts) == 0 {
								pos := fset.Position(mapKV.Pos())
								var dimName string
								if bl, ok := mapKV.Key.(*ast.BasicLit); ok {
									dimName = strings.Trim(bl.Value, `"'`)
								}
								issues = append(issues, LintIssue{
									File:     path,
									Line:     pos.Line,
									Column:   pos.Column,
									Severity: "error",
									Message:  fmt.Sprintf("Empty matrix dimension %q - must have at least one value", dimName),
									Rule:     r.ID(),
									Fixable:  false,
								})
							}
						}
					}
				}
			}
		}
		return true
	})

	return issues
}

// WAG014 checks for jobs without TimeoutMinutes set.
type WAG014 struct{}

func (r *WAG014) ID() string          { return "WAG014" }
func (r *WAG014) Description() string { return "Jobs should have timeout-minutes set" }

func (r *WAG014) Check(fset *token.FileSet, file *ast.File, path string) []LintIssue {
	var issues []LintIssue

	ast.Inspect(file, func(n ast.Node) bool {
		lit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}

		typeName := getTypeName(lit.Type)
		if typeName != "workflow.Job" && typeName != "Job" {
			return true
		}

		hasTimeout := false
		for _, elt := range lit.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			key, ok := kv.Key.(*ast.Ident)
			if !ok {
				continue
			}
			if key.Name == "TimeoutMinutes" {
				hasTimeout = true
				break
			}
		}

		if !hasTimeout {
			pos := fset.Position(lit.Pos())
			issues = append(issues, LintIssue{
				File:     path,
				Line:     pos.Line,
				Column:   pos.Column,
				Severity: "warning",
				Message:  "Job missing TimeoutMinutes - consider adding a timeout (e.g., 30 minutes)",
				Rule:     r.ID(),
				Fixable:  false,
			})
		}
		return true
	})

	return issues
}

// WAG016 validates concurrency settings.
type WAG016 struct{}

func (r *WAG016) ID() string          { return "WAG016" }
func (r *WAG016) Description() string { return "Validate concurrency settings" }

func (r *WAG016) Check(fset *token.FileSet, file *ast.File, path string) []LintIssue {
	var issues []LintIssue

	ast.Inspect(file, func(n ast.Node) bool {
		lit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}

		typeName := getTypeName(lit.Type)
		if typeName != "workflow.Concurrency" && typeName != "Concurrency" {
			return true
		}

		hasGroup := false
		hasCancelInProgress := false
		var cancelPos token.Position

		for _, elt := range lit.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			key, ok := kv.Key.(*ast.Ident)
			if !ok {
				continue
			}

			switch key.Name {
			case "Group":
				if bl, ok := kv.Value.(*ast.BasicLit); ok && bl.Kind == token.STRING {
					val := strings.Trim(bl.Value, `"'`)
					if val != "" {
						hasGroup = true
					}
				} else if _, ok := kv.Value.(*ast.Ident); ok {
					hasGroup = true
				}
			case "CancelInProgress":
				if ident, ok := kv.Value.(*ast.Ident); ok && ident.Name == "true" {
					hasCancelInProgress = true
					cancelPos = fset.Position(kv.Pos())
				}
			}
		}

		if hasCancelInProgress && !hasGroup {
			issues = append(issues, LintIssue{
				File:     path,
				Line:     cancelPos.Line,
				Column:   cancelPos.Column,
				Severity: "warning",
				Message:  "CancelInProgress is set without a Group - define a concurrency group",
				Rule:     r.ID(),
				Fixable:  false,
			})
		}

		return true
	})

	return issues
}
