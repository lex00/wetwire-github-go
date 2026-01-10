package discover

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// DiscoveredPRTemplate represents a PRTemplate found by AST parsing.
type DiscoveredPRTemplate struct {
	Name string // Variable name
	File string // Source file path
	Line int    // Line number
}

// PRTemplateDiscoveryResult contains discovered PRTemplates.
type PRTemplateDiscoveryResult struct {
	Templates []DiscoveredPRTemplate
	Errors    []string
}

// DiscoverPRTemplates finds all PRTemplates in the given directory.
func (d *Discoverer) DiscoverPRTemplates(dir string) (*PRTemplateDiscoveryResult, error) {
	result := &PRTemplateDiscoveryResult{
		Templates: []DiscoveredPRTemplate{},
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

		// Find PR template variables
		d.processPRTemplateFile(file, path, result)

		return nil
	})

	return result, err
}

// processPRTemplateFile processes a single Go file to find PRTemplates.
func (d *Discoverer) processPRTemplateFile(file *ast.File, path string, result *PRTemplateDiscoveryResult) {
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
				// Check if this is a PRTemplate
				if valueSpec.Type != nil {
					typeName := d.getTypeName(valueSpec.Type)
					if typeName == "templates.PRTemplate" || typeName == "PRTemplate" {
						pos := d.fset.Position(name.Pos())
						result.Templates = append(result.Templates, DiscoveredPRTemplate{
							Name: name.Name,
							File: path,
							Line: pos.Line,
						})
					}
				}

				// If no explicit type, check the value
				if valueSpec.Type == nil && len(valueSpec.Values) > i {
					typeName := d.inferTypeFromValue(valueSpec.Values[i])
					if typeName == "templates.PRTemplate" || typeName == "PRTemplate" {
						pos := d.fset.Position(name.Pos())
						result.Templates = append(result.Templates, DiscoveredPRTemplate{
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
