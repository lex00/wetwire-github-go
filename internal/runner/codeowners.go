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

// ExtractedCodeownersRule contains the extracted values for a Rule.
type ExtractedCodeownersRule struct {
	Pattern string   `json:"pattern"`
	Owners  []string `json:"owners"`
	Comment string   `json:"comment,omitempty"`
}

// ExtractedCodeowners contains the extracted values for a Codeowners config.
type ExtractedCodeowners struct {
	Name  string                    `json:"name"`
	Rules []ExtractedCodeownersRule `json:"rules"`
}

// CodeownersExtractionResult contains all extracted Codeowners configs.
type CodeownersExtractionResult struct {
	Configs []ExtractedCodeowners `json:"configs"`
	Error   string                `json:"error,omitempty"`
}

// ExtractCodeowners extracts values from discovered Codeowners configs.
func (r *Runner) ExtractCodeowners(dir string, discovered *discover.CodeownersDiscoveryResult) (*CodeownersExtractionResult, error) {
	if len(discovered.Configs) == 0 {
		return &CodeownersExtractionResult{
			Configs: []ExtractedCodeowners{},
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
	program, err := r.generateCodeownersProgram(modulePath, absDir, discovered)
	if err != nil {
		return nil, fmt.Errorf("generating program: %w", err)
	}

	// Create temp directory for the program
	tempDir, err := os.MkdirTemp(r.TempDir, "wetwire-extract-codeowners-*")
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
	var result CodeownersExtractionResult
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("parsing output: %w\n%s", err, output)
	}

	return &result, nil
}

// generateCodeownersProgram generates a Go program that extracts Codeowners values.
func (r *Runner) generateCodeownersProgram(modulePath, baseDir string, discovered *discover.CodeownersDiscoveryResult) (string, error) {
	var imports []string
	var vars []string

	// Group configs by package
	packages := make(map[string][]discover.DiscoveredCodeowners)
	for _, cfg := range discovered.Configs {
		pkgPath := r.getPackagePath(modulePath, baseDir, cfg.File)
		packages[pkgPath] = append(packages[pkgPath], cfg)
	}

	// Generate imports and variable references
	for pkgPath, configs := range packages {
		alias := r.pkgAlias(pkgPath)
		imports = append(imports, fmt.Sprintf(`%s "%s"`, alias, pkgPath))

		for _, cfg := range configs {
			vars = append(vars, fmt.Sprintf(`extractConfig("%s", %s.%s.Rules)`,
				cfg.Name, alias, cfg.Name))
		}
	}

	program := `package main

import (
	"encoding/json"
	"fmt"
	` + strings.Join(imports, "\n\t") + `
	"github.com/lex00/wetwire-github-go/codeowners"
)

type ExtractedCodeownersRule struct {
	Pattern string   ` + "`json:\"pattern\"`" + `
	Owners  []string ` + "`json:\"owners\"`" + `
	Comment string   ` + "`json:\"comment,omitempty\"`" + `
}

type ExtractedCodeowners struct {
	Name  string                    ` + "`json:\"name\"`" + `
	Rules []ExtractedCodeownersRule ` + "`json:\"rules\"`" + `
}

type CodeownersExtractionResult struct {
	Configs []ExtractedCodeowners ` + "`json:\"configs\"`" + `
}

func extractConfig(name string, rules []codeowners.Rule) ExtractedCodeowners {
	extracted := ExtractedCodeowners{
		Name:  name,
		Rules: make([]ExtractedCodeownersRule, len(rules)),
	}
	for i, rule := range rules {
		extracted.Rules[i] = ExtractedCodeownersRule{
			Pattern: rule.Pattern,
			Owners:  rule.Owners,
			Comment: rule.Comment,
		}
	}
	return extracted
}

func main() {
	result := CodeownersExtractionResult{
		Configs: []ExtractedCodeowners{
			` + strings.Join(vars, ",\n\t\t\t") + `
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
