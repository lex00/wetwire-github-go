package linter

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

// WAG003 checks for hardcoded secrets instead of using the secrets context.
type WAG003 struct{}

func (r *WAG003) ID() string          { return "WAG003" }
func (r *WAG003) Description() string { return "Use secrets context instead of hardcoded strings" }

func (r *WAG003) Check(fset *token.FileSet, file *ast.File, path string) []LintIssue {
	var issues []LintIssue

	ghSecretPatterns := []string{
		"ghp_",
		"ghs_",
		"ghu_",
		"ghr_",
		"github_pat_",
	}

	ast.Inspect(file, func(n ast.Node) bool {
		bl, ok := n.(*ast.BasicLit)
		if !ok || bl.Kind != token.STRING {
			return true
		}

		val := strings.Trim(bl.Value, `"'`)
		for _, pattern := range ghSecretPatterns {
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

// WAG017 suggests adding explicit permissions scope to workflows.
type WAG017 struct{}

func (r *WAG017) ID() string          { return "WAG017" }
func (r *WAG017) Description() string { return "Suggest adding explicit permissions scope for security" }

func (r *WAG017) Check(fset *token.FileSet, file *ast.File, path string) []LintIssue {
	var issues []LintIssue

	ast.Inspect(file, func(n ast.Node) bool {
		lit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}

		typeName := getTypeName(lit.Type)
		if typeName != "workflow.Workflow" && typeName != "Workflow" {
			return true
		}

		hasPermissions := false
		for _, elt := range lit.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			key, ok := kv.Key.(*ast.Ident)
			if !ok {
				continue
			}
			if key.Name == "Permissions" {
				hasPermissions = true
				break
			}
		}

		if !hasPermissions {
			pos := fset.Position(lit.Pos())
			issues = append(issues, LintIssue{
				File:     path,
				Line:     pos.Line,
				Column:   pos.Column,
				Severity: "info",
				Message:  "Consider adding explicit Permissions field for security best practices",
				Rule:     r.ID(),
				Fixable:  false,
			})
		}

		return true
	})

	return issues
}

// WAG018 detects dangerous pull_request_target patterns.
type WAG018 struct{}

func (r *WAG018) ID() string { return "WAG018" }
func (r *WAG018) Description() string {
	return "Detect dangerous pull_request_target patterns with checkout actions"
}

func (r *WAG018) Check(fset *token.FileSet, file *ast.File, path string) []LintIssue {
	var issues []LintIssue

	triggersWithPRTarget := make(map[string]bool)

	// First pass: find Triggers with pull_request_target
	ast.Inspect(file, func(n ast.Node) bool {
		genDecl, ok := n.(*ast.GenDecl)
		if !ok {
			return true
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			for i, name := range valueSpec.Names {
				if len(valueSpec.Values) <= i {
					continue
				}

				lit, ok := valueSpec.Values[i].(*ast.CompositeLit)
				if !ok {
					continue
				}

				typeName := getTypeName(lit.Type)
				if typeName == "workflow.Triggers" || typeName == "Triggers" {
					if checkTriggersForPRTarget(lit) {
						triggersWithPRTarget[name.Name] = true
					}
				}
			}
		}
		return true
	})

	// Second pass: find workflows and check them
	ast.Inspect(file, func(n ast.Node) bool {
		genDecl, ok := n.(*ast.GenDecl)
		if !ok {
			return true
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			for i := range valueSpec.Names {
				if len(valueSpec.Values) <= i {
					continue
				}

				lit, ok := valueSpec.Values[i].(*ast.CompositeLit)
				if !ok {
					continue
				}

				typeName := getTypeName(lit.Type)
				if typeName != "workflow.Workflow" && typeName != "Workflow" {
					continue
				}

				hasPRTarget := hasPullRequestTargetWithTracking(lit, triggersWithPRTarget)
				if !hasPRTarget {
					continue
				}

				if hasCheckoutAction(lit) {
					pos := fset.Position(lit.Pos())
					issues = append(issues, LintIssue{
						File:     path,
						Line:     pos.Line,
						Column:   pos.Column,
						Severity: "warning",
						Message:  "Workflow uses pull_request_target with checkout action - potential security risk. Consider using pull_request trigger or reviewing security implications.",
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

// WAG020 detects hardcoded secrets and credentials in code.
type WAG020 struct{}

func (r *WAG020) ID() string { return "WAG020" }
func (r *WAG020) Description() string {
	return "Detect hardcoded secrets, API keys, and credentials"
}

func (r *WAG020) Check(fset *token.FileSet, file *ast.File, path string) []LintIssue {
	var issues []LintIssue

	ast.Inspect(file, func(n ast.Node) bool {
		bl, ok := n.(*ast.BasicLit)
		if !ok || bl.Kind != token.STRING {
			return true
		}

		val := strings.Trim(bl.Value, `"'`+"`")

		if strings.Contains(val, "secrets.") {
			return true
		}

		for _, sp := range secretPatterns {
			if sp.pattern.MatchString(val) {
				pos := fset.Position(bl.Pos())
				issues = append(issues, LintIssue{
					File:     path,
					Line:     pos.Line,
					Column:   pos.Column,
					Severity: "error",
					Message:  fmt.Sprintf("Hardcoded %s detected - use secrets context instead", sp.description),
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
