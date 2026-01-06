package discover

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// DiscoveredDiscussionTemplate represents a DiscussionTemplate found by AST parsing.
type DiscoveredDiscussionTemplate struct {
	Name string // Variable name
	File string // Source file path
	Line int    // Line number
}

// DiscussionTemplateDiscoveryResult contains discovered DiscussionTemplates.
type DiscussionTemplateDiscoveryResult struct {
	Templates []DiscoveredDiscussionTemplate
	Errors    []string
}

// DiscoverDiscussionTemplates finds all DiscussionTemplates in the given directory.
func (d *Discoverer) DiscoverDiscussionTemplates(dir string) (*DiscussionTemplateDiscoveryResult, error) {
	result := &DiscussionTemplateDiscoveryResult{
		Templates: []DiscoveredDiscussionTemplate{},
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

		// Find discussion template variables
		d.processDiscussionTemplateFile(file, path, result)

		return nil
	})

	return result, err
}

// processDiscussionTemplateFile processes a single Go file to find DiscussionTemplates.
func (d *Discoverer) processDiscussionTemplateFile(file *ast.File, path string, result *DiscussionTemplateDiscoveryResult) {
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
				// Check if this is a DiscussionTemplate
				if valueSpec.Type != nil {
					typeName := d.getTypeName(valueSpec.Type)
					if typeName == "templates.DiscussionTemplate" || typeName == "DiscussionTemplate" {
						pos := d.fset.Position(name.Pos())
						result.Templates = append(result.Templates, DiscoveredDiscussionTemplate{
							Name: name.Name,
							File: path,
							Line: pos.Line,
						})
					}
				}

				// If no explicit type, check the value
				if valueSpec.Type == nil && len(valueSpec.Values) > i {
					typeName := d.inferTypeFromValue(valueSpec.Values[i])
					if typeName == "templates.DiscussionTemplate" || typeName == "DiscussionTemplate" {
						pos := d.fset.Position(name.Pos())
						result.Templates = append(result.Templates, DiscoveredDiscussionTemplate{
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
