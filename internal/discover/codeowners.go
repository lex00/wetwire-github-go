package discover

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// DiscoveredCodeowners represents a Codeowners configuration found by AST parsing.
type DiscoveredCodeowners struct {
	Name string // Variable name
	File string // Source file path
	Line int    // Line number
}

// CodeownersDiscoveryResult contains discovered Codeowners configurations.
type CodeownersDiscoveryResult struct {
	Configs []DiscoveredCodeowners
	Errors  []string
}

// DiscoverCodeowners finds all Codeowners configurations in the given directory.
func (d *Discoverer) DiscoverCodeowners(dir string) (*CodeownersDiscoveryResult, error) {
	result := &CodeownersDiscoveryResult{
		Configs: []DiscoveredCodeowners{},
		Errors:  []string{},
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

		// Find codeowners variables
		d.processCodeownersFile(file, path, result)

		return nil
	})

	return result, err
}

// processCodeownersFile processes a single Go file to find Codeowners configs.
func (d *Discoverer) processCodeownersFile(file *ast.File, path string, result *CodeownersDiscoveryResult) {
	// Check if this file imports the codeowners package
	if !d.hasCodeownersImport(file) {
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
				// Check if this is a codeowners.Owners
				if valueSpec.Type != nil {
					typeName := d.getTypeName(valueSpec.Type)
					if typeName == "codeowners.Owners" || typeName == "Owners" {
						pos := d.fset.Position(name.Pos())
						result.Configs = append(result.Configs, DiscoveredCodeowners{
							Name: name.Name,
							File: path,
							Line: pos.Line,
						})
					}
				}

				// If no explicit type, check the value
				if valueSpec.Type == nil && len(valueSpec.Values) > i {
					typeName := d.inferTypeFromValue(valueSpec.Values[i])
					if typeName == "codeowners.Owners" || typeName == "Owners" {
						pos := d.fset.Position(name.Pos())
						result.Configs = append(result.Configs, DiscoveredCodeowners{
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

// hasCodeownersImport checks if the file imports the codeowners package.
func (d *Discoverer) hasCodeownersImport(file *ast.File) bool {
	for _, imp := range file.Imports {
		if imp.Path != nil {
			path := strings.Trim(imp.Path.Value, `"`)
			if strings.HasSuffix(path, "/codeowners") || path == "codeowners" {
				return true
			}
		}
	}
	return false
}
