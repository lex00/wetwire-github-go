// Package runner provides value extraction from Go declarations.
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

// Runner extracts values from Go declarations by executing a temporary program.
type Runner struct {
	// TempDir is the directory for temporary files
	TempDir string
	// GoPath is the path to the Go binary
	GoPath string
	// Verbose enables verbose logging
	Verbose bool
}

// NewRunner creates a new Runner.
func NewRunner() *Runner {
	goPath, _ := exec.LookPath("go")
	return &Runner{
		TempDir: os.TempDir(),
		GoPath:  goPath,
	}
}

// ExtractedWorkflow contains the extracted values for a workflow.
type ExtractedWorkflow struct {
	Name string                 `json:"name"`
	Data map[string]any         `json:"data"`
}

// ExtractedJob contains the extracted values for a job.
type ExtractedJob struct {
	Name string         `json:"name"`
	Data map[string]any `json:"data"`
}

// ExtractionResult contains all extracted values.
type ExtractionResult struct {
	Workflows []ExtractedWorkflow `json:"workflows"`
	Jobs      []ExtractedJob      `json:"jobs"`
	Error     string              `json:"error,omitempty"`
}

// ExtractValues extracts values from discovered workflows and jobs.
func (r *Runner) ExtractValues(dir string, discovered *discover.DiscoveryResult) (*ExtractionResult, error) {
	if len(discovered.Workflows) == 0 && len(discovered.Jobs) == 0 {
		return &ExtractionResult{
			Workflows: []ExtractedWorkflow{},
			Jobs:      []ExtractedJob{},
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
	program, err := r.generateProgram(modulePath, absDir, discovered)
	if err != nil {
		return nil, fmt.Errorf("generating program: %w", err)
	}

	// Create temp directory for the program
	tempDir, err := os.MkdirTemp(r.TempDir, "wetwire-extract-*")
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
	var result ExtractionResult
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("parsing output: %w\nOutput: %s", err, output)
	}

	return &result, nil
}

// parseGoMod extracts the module path from go.mod.
func (r *Runner) parseGoMod(dir string) (string, error) {
	goModPath := filepath.Join(dir, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return "", err
	}

	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimPrefix(line, "module "), nil
		}
	}

	return "", fmt.Errorf("module directive not found in go.mod")
}

// generateGoMod creates a go.mod with a replace directive.
func (r *Runner) generateGoMod(modulePath, dir string) string {
	absDir, _ := filepath.Abs(dir)

	// Start with base go.mod
	var sb strings.Builder
	sb.WriteString("module wetwire-extract\n\n")
	sb.WriteString("go 1.23\n\n")
	sb.WriteString(fmt.Sprintf("require %s v0.0.0\n\n", modulePath))
	sb.WriteString(fmt.Sprintf("replace %s => %s\n", modulePath, absDir))

	// Copy replace directives from user's go.mod
	replaces := r.parseReplaceDirectives(dir)
	for _, replace := range replaces {
		sb.WriteString(replace + "\n")
	}

	return sb.String()
}

// parseReplaceDirectives extracts replace directives from go.mod.
func (r *Runner) parseReplaceDirectives(dir string) []string {
	var replaces []string

	goModPath := filepath.Join(dir, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return replaces
	}

	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "replace ") {
			replaces = append(replaces, line)
		}
	}

	return replaces
}

// generateProgram creates the extraction program source code.
func (r *Runner) generateProgram(modulePath, baseDir string, discovered *discover.DiscoveryResult) (string, error) {
	var sb strings.Builder

	sb.WriteString(`package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

`)

	// Import the user's package
	// We need to determine the package path from the discovered files
	packages := make(map[string]bool)
	for _, w := range discovered.Workflows {
		pkgPath := r.getPackagePath(modulePath, baseDir, w.File)
		packages[pkgPath] = true
	}
	for _, j := range discovered.Jobs {
		pkgPath := r.getPackagePath(modulePath, baseDir, j.File)
		packages[pkgPath] = true
	}

	for pkgPath := range packages {
		alias := r.pkgAlias(pkgPath)
		sb.WriteString(fmt.Sprintf("\t%s %q\n", alias, pkgPath))
	}

	sb.WriteString(`)

type ExtractionResult struct {
	Workflows []ExtractedWorkflow ` + "`json:\"workflows\"`" + `
	Jobs      []ExtractedJob      ` + "`json:\"jobs\"`" + `
}

type ExtractedWorkflow struct {
	Name string                 ` + "`json:\"name\"`" + `
	Data map[string]any         ` + "`json:\"data\"`" + `
}

type ExtractedJob struct {
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
	result := ExtractionResult{
		Workflows: []ExtractedWorkflow{},
		Jobs:      []ExtractedJob{},
	}

`)

	// Add workflow extractions
	for _, w := range discovered.Workflows {
		alias := r.pkgAlias(r.getPackagePath(modulePath, baseDir, w.File))
		sb.WriteString(fmt.Sprintf("\tresult.Workflows = append(result.Workflows, ExtractedWorkflow{Name: %q, Data: toMap(%s.%s)})\n",
			w.Name, alias, w.Name))
	}

	// Add job extractions
	for _, j := range discovered.Jobs {
		alias := r.pkgAlias(r.getPackagePath(modulePath, baseDir, j.File))
		sb.WriteString(fmt.Sprintf("\tresult.Jobs = append(result.Jobs, ExtractedJob{Name: %q, Data: toMap(%s.%s)})\n",
			j.Name, alias, j.Name))
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

// getPackagePath determines the import path for a source file.
func (r *Runner) getPackagePath(modulePath, baseDir, file string) string {
	// Make file path relative to baseDir
	relFile, err := filepath.Rel(baseDir, file)
	if err != nil {
		// Fall back to just the module path if we can't compute relative path
		return modulePath
	}

	// Get directory containing the file
	dir := filepath.Dir(relFile)

	// If the file is in the root, return the module path
	if dir == "." || dir == "" {
		return modulePath
	}

	// Otherwise, append the relative path
	return modulePath + "/" + filepath.ToSlash(dir)
}

// pkgAlias creates a safe import alias for a package path.
func (r *Runner) pkgAlias(pkgPath string) string {
	// Use the last component as the alias
	parts := strings.Split(pkgPath, "/")
	alias := parts[len(parts)-1]
	// Replace hyphens with underscores
	alias = strings.ReplaceAll(alias, "-", "_")
	return alias
}

// FindGoBinary locates the Go binary.
func FindGoBinary() (string, error) {
	path, err := exec.LookPath("go")
	if err != nil {
		return "", fmt.Errorf("go binary not found: %w", err)
	}
	return path, nil
}
