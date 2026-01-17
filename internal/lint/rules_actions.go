package lint

import (
	"fmt"
	"go/ast"
	"go/format"
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

		typeName := getTypeName(lit.Type)
		if typeName != "workflow.Step" && typeName != "Step" {
			return true
		}

		for _, elt := range lit.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			key, ok := kv.Key.(*ast.Ident)
			if !ok || key.Name != "Uses" {
				continue
			}
			if bl, ok := kv.Value.(*ast.BasicLit); ok {
				pos := fset.Position(kv.Pos())
				usesVal := strings.Trim(bl.Value, `"'`)
				actionName := parseActionRef(usesVal)
				_, fixable := knownActions[actionName]
				issues = append(issues, LintIssue{
					File:     path,
					Line:     pos.Line,
					Column:   pos.Column,
					Severity: SeverityWarning,
					Message:  "Use typed action wrapper instead of raw uses: string",
					Rule:     r.ID(),
					Fixable:  fixable,
				})
			}
		}
		return true
	})

	return issues
}

// Fix implements the Fixer interface for WAG001.
func (r *WAG001) Fix(fset *token.FileSet, file *ast.File, path string, src []byte, issue LintIssue) ([]byte, error) {
	var targetLit *ast.CompositeLit
	var usesValue string

	ast.Inspect(file, func(n ast.Node) bool {
		lit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}

		typeName := getTypeName(lit.Type)
		if typeName != "workflow.Step" && typeName != "Step" {
			return true
		}

		pos := fset.Position(lit.Pos())
		for _, elt := range lit.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			key, ok := kv.Key.(*ast.Ident)
			if !ok || key.Name != "Uses" {
				continue
			}
			kvPos := fset.Position(kv.Pos())
			if kvPos.Line == issue.Line {
				targetLit = lit
				if bl, ok := kv.Value.(*ast.BasicLit); ok {
					usesValue = strings.Trim(bl.Value, `"'`)
				}
				return false
			}
		}
		_ = pos
		return true
	})

	if targetLit == nil || usesValue == "" {
		return nil, fmt.Errorf("could not find target step")
	}

	actionName := parseActionRef(usesValue)
	info, ok := knownActions[actionName]
	if !ok {
		return nil, fmt.Errorf("unknown action: %s", actionName)
	}

	replacement := fmt.Sprintf("%s.%s{}", info.pkg, info.typ)

	startPos := fset.Position(targetLit.Pos())
	endPos := fset.Position(targetLit.End())

	result := make([]byte, 0, len(src))
	result = append(result, src[:startPos.Offset]...)
	result = append(result, []byte(replacement)...)
	result = append(result, src[endPos.Offset:]...)

	result = addImportIfNeeded(result, info.importPath, info.pkg)

	formatted, err := format.Source(result)
	if err != nil {
		return result, nil
	}

	return formatted, nil
}

// Ensure WAG001 implements Fixer
var _ Fixer = (*WAG001)(nil)

// WAG010 flags missing recommended action inputs.
type WAG010 struct{}

func (r *WAG010) ID() string          { return "WAG010" }
func (r *WAG010) Description() string { return "Flag missing recommended action inputs" }

func (r *WAG010) Check(fset *token.FileSet, file *ast.File, path string) []LintIssue {
	var issues []LintIssue

	ast.Inspect(file, func(n ast.Node) bool {
		lit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}

		typeName := getTypeName(lit.Type)
		recommended, exists := recommendedInputs[typeName]
		if !exists {
			return true
		}

		setFields := make(map[string]bool)
		for _, elt := range lit.Elts {
			if kv, ok := elt.(*ast.KeyValueExpr); ok {
				if key, ok := kv.Key.(*ast.Ident); ok {
					setFields[key.Name] = true
				}
			}
		}

		hasRecommended := false
		for _, field := range recommended {
			if setFields[field] {
				hasRecommended = true
				break
			}
		}

		if !hasRecommended && len(recommended) > 0 {
			pos := fset.Position(lit.Pos())
			issues = append(issues, LintIssue{
				File:     path,
				Line:     pos.Line,
				Column:   pos.Column,
				Severity: SeverityWarning,
				Message:  fmt.Sprintf("%s: consider setting %s", typeName, strings.Join(recommended, " or ")),
				Rule:     r.ID(),
				Fixable:  false,
			})
		}
		return true
	})

	return issues
}

// WAG012 warns about deprecated action versions.
type WAG012 struct{}

func (r *WAG012) ID() string          { return "WAG012" }
func (r *WAG012) Description() string { return "Warn about deprecated action versions" }

func (r *WAG012) Check(fset *token.FileSet, file *ast.File, path string) []LintIssue {
	var issues []LintIssue

	ast.Inspect(file, func(n ast.Node) bool {
		bl, ok := n.(*ast.BasicLit)
		if !ok || bl.Kind != token.STRING {
			return true
		}

		val := strings.Trim(bl.Value, `"'`)

		if !strings.Contains(val, "/") || !strings.Contains(val, "@") {
			return true
		}

		parts := strings.Split(val, "@")
		if len(parts) != 2 {
			return true
		}

		actionName := parts[0]
		version := parts[1]

		if info, exists := deprecatedVersions[actionName]; exists {
			for _, deprecated := range info.deprecated {
				if version == deprecated {
					pos := fset.Position(bl.Pos())
					issues = append(issues, LintIssue{
						File:     path,
						Line:     pos.Line,
						Column:   pos.Column,
						Severity: SeverityWarning,
						Message:  fmt.Sprintf("Action %s@%s is deprecated, use %s@%s", actionName, version, actionName, info.recommended),
						Rule:     r.ID(),
						Fixable:  false,
					})
					break
				}
			}
		}
		return true
	})

	return issues
}

// WAG015 suggests caching for setup actions.
type WAG015 struct{}

func (r *WAG015) ID() string { return "WAG015" }
func (r *WAG015) Description() string {
	return "Suggest caching for setup actions (setup-go, setup-node, setup-python)"
}

func (r *WAG015) Check(fset *token.FileSet, file *ast.File, path string) []LintIssue {
	var issues []LintIssue

	type stepsInfo struct {
		hasCache     bool
		setupActions []struct {
			typeName string
			pos      token.Position
		}
	}

	stepsArrays := make(map[string]*stepsInfo)

	ast.Inspect(file, func(n ast.Node) bool {
		comp, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}

		arrType, ok := comp.Type.(*ast.ArrayType)
		if !ok {
			return true
		}
		elemIdent, ok := arrType.Elt.(*ast.Ident)
		if !ok || elemIdent.Name != "any" {
			return true
		}

		info := &stepsInfo{}
		for _, elt := range comp.Elts {
			switch e := elt.(type) {
			case *ast.CompositeLit:
				typeName := getTypeName(e.Type)
				if typeName == "cache.Cache" || typeName == "Cache" {
					info.hasCache = true
				}
				if _, isSetup := setupActionsNeedingCache[typeName]; isSetup {
					info.setupActions = append(info.setupActions, struct {
						typeName string
						pos      token.Position
					}{typeName: typeName, pos: fset.Position(e.Pos())})
				}
			}
		}

		if !info.hasCache {
			for _, setup := range info.setupActions {
				displayName := setupActionsNeedingCache[setup.typeName]
				issues = append(issues, LintIssue{
					File:     path,
					Line:     setup.pos.Line,
					Column:   setup.pos.Column,
					Severity: SeverityWarning,
					Message:  fmt.Sprintf("Consider adding cache action for %s to improve build performance", displayName),
					Rule:     r.ID(),
					Fixable:  false,
				})
			}
		}

		return true
	})

	_ = stepsArrays

	return issues
}
