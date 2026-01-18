// Package discover provides AST-based resource discovery for wetwire workflows.
package discover

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	coreast "github.com/lex00/wetwire-core-go/ast"
)

// DiscoveredWorkflow represents a workflow found by AST parsing.
type DiscoveredWorkflow struct {
	Name string   // Variable name
	File string   // Source file path
	Line int      // Line number
	Jobs []string // Job variable names in this workflow
}

// DiscoveredJob represents a job found by AST parsing.
type DiscoveredJob struct {
	Name         string   // Variable name
	File         string   // Source file path
	Line         int      // Line number
	Dependencies []string // Referenced job names (Needs field)
}

// DiscoveryResult contains all discovered resources.
type DiscoveryResult struct {
	Workflows []DiscoveredWorkflow
	Jobs      []DiscoveredJob
	Errors    []string
}

// Discoverer finds workflow resources in Go source files.
type Discoverer struct {
	fset *token.FileSet
}

// NewDiscoverer creates a new Discoverer.
func NewDiscoverer() *Discoverer {
	return &Discoverer{
		fset: token.NewFileSet(),
	}
}

// Discover finds all workflow resources in the given directory.
func (d *Discoverer) Discover(dir string) (*DiscoveryResult, error) {
	result := &DiscoveryResult{
		Workflows: []DiscoveredWorkflow{},
		Jobs:      []DiscoveredJob{},
		Errors:    []string{},
	}

	// Walk the directory tree using coreast.WalkGoFiles
	opts := coreast.ParseOptions{
		SkipTests:   true,
		SkipVendor:  true,
		SkipHidden:  true,
		ExcludeDirs: []string{"testdata"},
	}

	err := coreast.WalkGoFiles(dir, opts, func(path string) error {
		// Parse the file
		file, err := parser.ParseFile(d.fset, path, nil, parser.ParseComments)
		if err != nil {
			result.Errors = append(result.Errors, err.Error())
			return nil
		}

		// Find workflow and job variables
		d.processFile(file, path, result)

		return nil
	})

	return result, err
}

// processFile processes a single Go file to find workflow resources.
func (d *Discoverer) processFile(file *ast.File, path string, result *DiscoveryResult) {
	// Check if this file imports the workflow package
	if !d.hasWorkflowImport(file) {
		return
	}

	// Look for package-level variable declarations
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			continue
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			for i, name := range valueSpec.Names {
				// Check if this is a workflow or job
				if valueSpec.Type != nil {
					typeName := d.getTypeName(valueSpec.Type)
					pos := d.fset.Position(name.Pos())

					if typeName == "workflow.Workflow" || typeName == "Workflow" {
						workflow := DiscoveredWorkflow{
							Name: name.Name,
							File: path,
							Line: pos.Line,
							Jobs: []string{},
						}
						// Try to extract jobs from the value
						if i < len(valueSpec.Values) {
							workflow.Jobs = d.extractJobRefs(valueSpec.Values[i])
						}
						result.Workflows = append(result.Workflows, workflow)
					}

					if typeName == "workflow.Job" || typeName == "Job" {
						job := DiscoveredJob{
							Name:         name.Name,
							File:         path,
							Line:         pos.Line,
							Dependencies: []string{},
						}
						// Try to extract dependencies from the value
						if i < len(valueSpec.Values) {
							job.Dependencies = d.extractDependencies(valueSpec.Values[i])
						}
						result.Jobs = append(result.Jobs, job)
					}
				}

				// If no explicit type, check the value
				if valueSpec.Type == nil && len(valueSpec.Values) > i {
					value := valueSpec.Values[i]
					typeName := d.inferTypeFromValue(value)
					pos := d.fset.Position(name.Pos())

					if typeName == "workflow.Workflow" || typeName == "Workflow" {
						workflow := DiscoveredWorkflow{
							Name: name.Name,
							File: path,
							Line: pos.Line,
							Jobs: d.extractJobRefs(value),
						}
						result.Workflows = append(result.Workflows, workflow)
					}

					if typeName == "workflow.Job" || typeName == "Job" {
						job := DiscoveredJob{
							Name:         name.Name,
							File:         path,
							Line:         pos.Line,
							Dependencies: d.extractDependencies(value),
						}
						result.Jobs = append(result.Jobs, job)
					}
				}
			}
		}
	}
}

// hasWorkflowImport checks if the file imports the workflow package.
func (d *Discoverer) hasWorkflowImport(file *ast.File) bool {
	for _, imp := range file.Imports {
		if imp.Path != nil {
			path := strings.Trim(imp.Path.Value, `"`)
			if strings.HasSuffix(path, "/workflow") || path == "workflow" {
				return true
			}
		}
	}
	return false
}

// getTypeName extracts the type name from a type expression.
func (d *Discoverer) getTypeName(expr ast.Expr) string {
	typeName, pkgName := coreast.ExtractTypeName(expr)
	if pkgName != "" {
		return pkgName + "." + typeName
	}
	return typeName
}

// inferTypeFromValue tries to infer the type from a composite literal.
func (d *Discoverer) inferTypeFromValue(expr ast.Expr) string {
	typeName, pkgName := coreast.InferTypeFromValue(expr)
	if pkgName != "" {
		return pkgName + "." + typeName
	}
	return typeName
}

// extractJobRefs extracts job references from a workflow value.
func (d *Discoverer) extractJobRefs(expr ast.Expr) []string {
	var refs []string

	lit, ok := expr.(*ast.CompositeLit)
	if !ok {
		return refs
	}

	for _, elt := range lit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		// Look for Jobs field
		key, ok := kv.Key.(*ast.Ident)
		if !ok || key.Name != "Jobs" {
			continue
		}

		// Extract identifiers from the Jobs value
		refs = d.extractIdentifiers(kv.Value)
	}

	return refs
}

// extractDependencies extracts job dependencies from a job value.
func (d *Discoverer) extractDependencies(expr ast.Expr) []string {
	var deps []string

	lit, ok := expr.(*ast.CompositeLit)
	if !ok {
		return deps
	}

	for _, elt := range lit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		// Look for Needs field
		key, ok := kv.Key.(*ast.Ident)
		if !ok || key.Name != "Needs" {
			continue
		}

		// Extract identifiers from the Needs value
		deps = d.extractIdentifiers(kv.Value)
	}

	return deps
}

// extractIdentifiers extracts all identifiers from an expression.
func (d *Discoverer) extractIdentifiers(expr ast.Expr) []string {
	var ids []string

	ast.Inspect(expr, func(n ast.Node) bool {
		if ident, ok := n.(*ast.Ident); ok {
			// Skip built-in identifiers and type names
			if !isBuiltinIdent(ident.Name) {
				ids = append(ids, ident.Name)
			}
		}
		return true
	})

	return ids
}

// isBuiltinIdent checks if an identifier is a built-in.
func isBuiltinIdent(name string) bool {
	// Check Go built-ins using core AST package
	if coreast.IsBuiltinIdent(name) {
		return true
	}

	// GitHub-specific type identifiers to skip
	githubTypes := map[string]bool{
		"workflow": true, "Job": true, "Workflow": true,
	}
	return githubTypes[name]
}
