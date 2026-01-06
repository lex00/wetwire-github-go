package discover

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// DiscoveredIssueTemplate represents an IssueTemplate found by AST parsing.
type DiscoveredIssueTemplate struct {
	Name string // Variable name
	File string // Source file path
	Line int    // Line number
}

// IssueTemplateDiscoveryResult contains discovered IssueTemplates.
type IssueTemplateDiscoveryResult struct {
	Templates []DiscoveredIssueTemplate
	Errors    []string
}

// DiscoverIssueTemplates finds all IssueTemplates in the given directory.
func (d *Discoverer) DiscoverIssueTemplates(dir string) (*IssueTemplateDiscoveryResult, error) {
	result := &IssueTemplateDiscoveryResult{
		Templates: []DiscoveredIssueTemplate{},
		Errors:    []string{},
	}

	// Walk the directory tree
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

		// Only process .go files
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Skip test files
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// Parse the file
		file, err := parser.ParseFile(d.fset, path, nil, parser.ParseComments)
		if err != nil {
			result.Errors = append(result.Errors, err.Error())
			return nil
		}

		// Find issue template variables
		d.processIssueTemplateFile(file, path, result)

		return nil
	})

	return result, err
}

// processIssueTemplateFile processes a single Go file to find IssueTemplates.
func (d *Discoverer) processIssueTemplateFile(file *ast.File, path string, result *IssueTemplateDiscoveryResult) {
	// Check if this file imports the templates package
	if !d.hasTemplatesImport(file) {
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
				// Check if this is an IssueTemplate
				if valueSpec.Type != nil {
					typeName := d.getTypeName(valueSpec.Type)
					if typeName == "templates.IssueTemplate" || typeName == "IssueTemplate" {
						pos := d.fset.Position(name.Pos())
						result.Templates = append(result.Templates, DiscoveredIssueTemplate{
							Name: name.Name,
							File: path,
							Line: pos.Line,
						})
					}
				}

				// If no explicit type, check the value
				if valueSpec.Type == nil && len(valueSpec.Values) > i {
					typeName := d.inferTypeFromValue(valueSpec.Values[i])
					if typeName == "templates.IssueTemplate" || typeName == "IssueTemplate" {
						pos := d.fset.Position(name.Pos())
						result.Templates = append(result.Templates, DiscoveredIssueTemplate{
							Name: name.Name,
							File: path,
							Line: pos.Line,
						})
					}
				}
			}
		}
	}
}

// hasTemplatesImport checks if the file imports the templates package.
func (d *Discoverer) hasTemplatesImport(file *ast.File) bool {
	for _, imp := range file.Imports {
		if imp.Path != nil {
			path := strings.Trim(imp.Path.Value, `"`)
			if strings.HasSuffix(path, "/templates") || path == "templates" {
				return true
			}
		}
	}
	return false
}
