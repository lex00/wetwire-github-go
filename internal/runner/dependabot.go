package runner

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lex00/wetwire-github-go/internal/discover"
)

// ExtractedDependabot contains the extracted values for a Dependabot config.
type ExtractedDependabot struct {
	Name string         `json:"name"`
	Data map[string]any `json:"data"`
}

// DependabotExtractionResult contains all extracted Dependabot configs.
type DependabotExtractionResult struct {
	Configs []ExtractedDependabot `json:"configs"`
	Error   string                `json:"error,omitempty"`
}

// ExtractDependabot extracts values from discovered Dependabot configs.
func (r *Runner) ExtractDependabot(dir string, discovered *discover.DependabotDiscoveryResult) (*DependabotExtractionResult, error) {
	if len(discovered.Configs) == 0 {
		return &DependabotExtractionResult{
			Configs: []ExtractedDependabot{},
		}, nil
	}

	// Get absolute path for consistent path handling
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("resolving path: %w", err)
	}

	// Parse go.mod to get module path
	modulePath, err := r.parseGoMod(absDir)
	if err != nil {
		return nil, fmt.Errorf("parsing go.mod: %w", err)
	}

	// Generate the temporary extraction program
	program, err := r.generateDependabotProgram(modulePath, absDir, discovered)
	if err != nil {
		return nil, fmt.Errorf("generating program: %w", err)
	}

	// Create temp directory for the program
	tempDir, err := os.MkdirTemp(r.TempDir, "wetwire-extract-dependabot-*")
	if err != nil {
		return nil, fmt.Errorf("creating temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Write the program
	programPath := filepath.Join(tempDir, "main.go")
	if err := os.WriteFile(programPath, []byte(program), 0644); err != nil {
		return nil, fmt.Errorf("writing program: %w", err)
	}

	// Write go.mod with replace directive
	goMod := r.generateGoMod(modulePath, dir)
	goModPath := filepath.Join(tempDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte(goMod), 0644); err != nil {
		return nil, fmt.Errorf("writing go.mod: %w", err)
	}

	// Run go mod tidy
	tidyCmd := exec.Command(r.GoPath, "mod", "tidy")
	tidyCmd.Dir = tempDir
	if output, err := tidyCmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("go mod tidy: %w\n%s", err, output)
	}

	// Execute the program
	runCmd := exec.Command(r.GoPath, "run", "main.go")
	runCmd.Dir = tempDir
	output, err := runCmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("running extraction: %w\n%s", err, output)
	}

	// Parse the JSON output
	var result DependabotExtractionResult
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("parsing output: %w\nOutput: %s", err, output)
	}

	return &result, nil
}

// generateDependabotProgram creates the extraction program for Dependabot configs.
func (r *Runner) generateDependabotProgram(modulePath, baseDir string, discovered *discover.DependabotDiscoveryResult) (string, error) {
	var sb strings.Builder

	sb.WriteString(`package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

`)

	// Import the user's package
	packages := make(map[string]bool)
	for _, c := range discovered.Configs {
		pkgPath := r.getPackagePath(modulePath, baseDir, c.File)
		packages[pkgPath] = true
	}

	for pkgPath := range packages {
		alias := r.pkgAlias(pkgPath)
		sb.WriteString(fmt.Sprintf("\t%s %q\n", alias, pkgPath))
	}

	sb.WriteString(`)

type DependabotExtractionResult struct {
	Configs []ExtractedDependabot ` + "`json:\"configs\"`" + `
}

type ExtractedDependabot struct {
	Name string         ` + "`json:\"name\"`" + `
	Data map[string]any ` + "`json:\"data\"`" + `
}

func toMap(v any) map[string]any {
	result := make(map[string]any)
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return result
	}
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath != "" { // unexported
			continue
		}
		result[field.Name] = val.Field(i).Interface()
	}
	return result
}

func main() {
	result := DependabotExtractionResult{
		Configs: []ExtractedDependabot{},
	}

`)

	// Add config extractions
	for _, c := range discovered.Configs {
		alias := r.pkgAlias(r.getPackagePath(modulePath, baseDir, c.File))
		sb.WriteString(fmt.Sprintf("\tresult.Configs = append(result.Configs, ExtractedDependabot{Name: %q, Data: toMap(%s.%s)})\n",
			c.Name, alias, c.Name))
	}

	sb.WriteString(`
	data, err := json.Marshal(result)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling result: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}
`)

	return sb.String(), nil
}
