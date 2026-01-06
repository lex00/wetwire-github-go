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

// ExtractedIssueTemplate contains the extracted values for an IssueTemplate.
type ExtractedIssueTemplate struct {
	Name string         `json:"name"`
	Data map[string]any `json:"data"`
}

// IssueTemplateExtractionResult contains all extracted IssueTemplates.
type IssueTemplateExtractionResult struct {
	Templates []ExtractedIssueTemplate `json:"templates"`
	Error     string                   `json:"error,omitempty"`
}

// ExtractIssueTemplates extracts values from discovered IssueTemplates.
func (r *Runner) ExtractIssueTemplates(dir string, discovered *discover.IssueTemplateDiscoveryResult) (*IssueTemplateExtractionResult, error) {
	if len(discovered.Templates) == 0 {
		return &IssueTemplateExtractionResult{
			Templates: []ExtractedIssueTemplate{},
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
	program, err := r.generateIssueTemplateProgram(modulePath, absDir, discovered)
	if err != nil {
		return nil, fmt.Errorf("generating program: %w", err)
	}

	// Create temp directory for the program
	tempDir, err := os.MkdirTemp(r.TempDir, "wetwire-extract-issue-template-*")
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
	var result IssueTemplateExtractionResult
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("parsing output: %w\nOutput: %s", err, output)
	}

	return &result, nil
}

// generateIssueTemplateProgram creates the extraction program for IssueTemplates.
func (r *Runner) generateIssueTemplateProgram(modulePath, baseDir string, discovered *discover.IssueTemplateDiscoveryResult) (string, error) {
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
	for _, t := range discovered.Templates {
		pkgPath := r.getPackagePath(modulePath, baseDir, t.File)
		packages[pkgPath] = true
	}

	for pkgPath := range packages {
		alias := r.pkgAlias(pkgPath)
		sb.WriteString(fmt.Sprintf("\t%s %q\n", alias, pkgPath))
	}

	sb.WriteString(`)

type IssueTemplateExtractionResult struct {
	Templates []ExtractedIssueTemplate ` + "`json:\"templates\"`" + `
}

type ExtractedIssueTemplate struct {
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
	result := IssueTemplateExtractionResult{
		Templates: []ExtractedIssueTemplate{},
	}

`)

	// Add template extractions
	for _, t := range discovered.Templates {
		alias := r.pkgAlias(r.getPackagePath(modulePath, baseDir, t.File))
		sb.WriteString(fmt.Sprintf("\tresult.Templates = append(result.Templates, ExtractedIssueTemplate{Name: %q, Data: toMap(%s.%s)})\n",
			t.Name, alias, t.Name))
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
