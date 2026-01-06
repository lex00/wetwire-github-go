package discover

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// DiscoveredDependabot represents a Dependabot config found by AST parsing.
type DiscoveredDependabot struct {
	Name string // Variable name
	File string // Source file path
	Line int    // Line number
}

// DependabotDiscoveryResult contains discovered Dependabot configs.
type DependabotDiscoveryResult struct {
	Configs []DiscoveredDependabot
	Errors  []string
}

// DiscoverDependabot finds all Dependabot configs in the given directory.
func (d *Discoverer) DiscoverDependabot(dir string) (*DependabotDiscoveryResult, error) {
	result := &DependabotDiscoveryResult{
		Configs: []DiscoveredDependabot{},
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

		// Find dependabot variables
		d.processDependabotFile(file, path, result)

		return nil
	})

	return result, err
}

// processDependabotFile processes a single Go file to find Dependabot configs.
func (d *Discoverer) processDependabotFile(file *ast.File, path string, result *DependabotDiscoveryResult) {
	// Check if this file imports the dependabot package
	if !d.hasDependabotImport(file) {
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
				// Check if this is a Dependabot config
				if valueSpec.Type != nil {
					typeName := d.getTypeName(valueSpec.Type)
					if typeName == "dependabot.Dependabot" || typeName == "Dependabot" {
						pos := d.fset.Position(name.Pos())
						result.Configs = append(result.Configs, DiscoveredDependabot{
							Name: name.Name,
							File: path,
							Line: pos.Line,
						})
					}
				}

				// If no explicit type, check the value
				if valueSpec.Type == nil && len(valueSpec.Values) > i {
					typeName := d.inferTypeFromValue(valueSpec.Values[i])
					if typeName == "dependabot.Dependabot" || typeName == "Dependabot" {
						pos := d.fset.Position(name.Pos())
						result.Configs = append(result.Configs, DiscoveredDependabot{
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

// hasDependabotImport checks if the file imports the dependabot package.
func (d *Discoverer) hasDependabotImport(file *ast.File) bool {
	for _, imp := range file.Imports {
		if imp.Path != nil {
			path := strings.Trim(imp.Path.Value, `"`)
			if strings.HasSuffix(path, "/dependabot") || path == "dependabot" {
				return true
			}
		}
	}
	return false
}
