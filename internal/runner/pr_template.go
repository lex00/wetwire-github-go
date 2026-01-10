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

// ExtractedPRTemplate contains the extracted values for a PRTemplate.
type ExtractedPRTemplate struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

// PRTemplateExtractionResult contains all extracted PRTemplates.
type PRTemplateExtractionResult struct {
	Templates []ExtractedPRTemplate `json:"templates"`
	Error     string                `json:"error,omitempty"`
}

// ExtractPRTemplates extracts values from discovered PRTemplates.
func (r *Runner) ExtractPRTemplates(dir string, discovered *discover.PRTemplateDiscoveryResult) (*PRTemplateExtractionResult, error) {
	if len(discovered.Templates) == 0 {
		return &PRTemplateExtractionResult{
			Templates: []ExtractedPRTemplate{},
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
	program, err := r.generatePRTemplateProgram(modulePath, absDir, discovered)
	if err != nil {
		return nil, fmt.Errorf("generating program: %w", err)
	}

	// Create temp directory for the program
	tempDir, err := os.MkdirTemp(r.TempDir, "wetwire-extract-pr-template-*")
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
		return nil, fmt.Errorf("running program: %w\n%s", err, output)
	}

	// Parse the JSON output
	var result PRTemplateExtractionResult
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("parsing output: %w\n%s", err, output)
	}

	return &result, nil
}

// generatePRTemplateProgram generates a Go program that extracts PRTemplate values.
func (r *Runner) generatePRTemplateProgram(modulePath, baseDir string, discovered *discover.PRTemplateDiscoveryResult) (string, error) {
	var imports []string
	var vars []string

	// Group templates by package
	packages := make(map[string][]discover.DiscoveredPRTemplate)
	for _, tmpl := range discovered.Templates {
		pkgPath := r.getPackagePath(modulePath, baseDir, tmpl.File)
		packages[pkgPath] = append(packages[pkgPath], tmpl)
	}

	// Generate imports and variable references
	for pkgPath, templates := range packages {
		alias := r.pkgAlias(pkgPath)
		imports = append(imports, fmt.Sprintf(`%s "%s"`, alias, pkgPath))

		for _, tmpl := range templates {
			vars = append(vars, fmt.Sprintf(`{Name: "%s", Content: %s.%s.Content}`,
				tmpl.Name, alias, tmpl.Name))
		}
	}

	program := `package main

import (
	"encoding/json"
	"fmt"
	` + strings.Join(imports, "\n\t") + `
)

type ExtractedPRTemplate struct {
	Name    string ` + "`json:\"name\"`" + `
	Content string ` + "`json:\"content\"`" + `
}

type PRTemplateExtractionResult struct {
	Templates []ExtractedPRTemplate ` + "`json:\"templates\"`" + `
}

func main() {
	result := PRTemplateExtractionResult{
		Templates: []ExtractedPRTemplate{
			` + strings.Join(vars, ",\n\t\t\t") + `,
		},
	}

	data, err := json.Marshal(result)
	if err != nil {
		fmt.Printf("{\"error\": \"%s\"}\n", err.Error())
		return
	}
	fmt.Println(string(data))
}
`
	return program, nil
}
