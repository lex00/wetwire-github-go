package importer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Scaffold generates project scaffold files.
type Scaffold struct {
	// ModulePath is the Go module path
	ModulePath string
	// ProjectName is the name of the project
	ProjectName string
}

// NewScaffold creates a new Scaffold.
func NewScaffold(modulePath, projectName string) *Scaffold {
	return &Scaffold{
		ModulePath:  modulePath,
		ProjectName: projectName,
	}
}

// ScaffoldFiles contains the generated scaffold files.
type ScaffoldFiles struct {
	// Files maps filenames to content
	Files map[string]string
}

// Generate generates scaffold files.
func (s *Scaffold) Generate() *ScaffoldFiles {
	files := &ScaffoldFiles{
		Files: make(map[string]string),
	}

	files.Files["go.mod"] = s.generateGoMod()
	files.Files["cmd/main.go"] = s.generateMain()
	files.Files["README.md"] = s.generateReadme()
	files.Files["CLAUDE.md"] = s.generateClaudeMD()
	files.Files[".gitignore"] = s.generateGitignore()

	return files
}

// generateGoMod generates go.mod content.
func (s *Scaffold) generateGoMod() string {
	return fmt.Sprintf(`module %s

go 1.23

require github.com/lex00/wetwire-github-go v0.0.0
`, s.ModulePath)
}

// generateMain generates cmd/main.go content.
func (s *Scaffold) generateMain() string {
	return fmt.Sprintf(`package main

import (
	"fmt"

	// Import workflows to ensure they compile
	_ "%s/workflows"
)

func main() {
	fmt.Println("wetwire-github project: %s")
	fmt.Println("")
	fmt.Println("Build workflows:")
	fmt.Println("  wetwire-github build .")
	fmt.Println("")
	fmt.Println("Validate generated YAML:")
	fmt.Println("  wetwire-github validate .github/workflows/*.yml")
}
`, s.ModulePath, s.ProjectName)
}

// generateReadme generates README.md content.
func (s *Scaffold) generateReadme() string {
	return fmt.Sprintf("# %s\n\n"+
		"GitHub Actions workflows written in Go using [wetwire-github-go](https://github.com/lex00/wetwire-github-go).\n\n"+
		"## Building\n\n"+
		"Generate YAML files:\n\n"+
		"```bash\n"+
		"wetwire-github build .\n"+
		"```\n\n"+
		"This outputs to `.github/workflows/`.\n\n"+
		"## Validating\n\n"+
		"Validate generated YAML:\n\n"+
		"```bash\n"+
		"wetwire-github validate .github/workflows/*.yml\n"+
		"```\n\n"+
		"## Structure\n\n"+
		"- `workflows.go` - Workflow definitions\n"+
		"- `triggers.go` - Trigger configurations\n"+
		"- `jobs.go` - Job definitions\n"+
		"- `steps.go` - Step lists\n\n"+
		"## Linting\n\n"+
		"Check Go code for best practices:\n\n"+
		"```bash\n"+
		"wetwire-github lint .\n"+
		"```\n", s.ProjectName)
}

// generateClaudeMD generates CLAUDE.md content.
func (s *Scaffold) generateClaudeMD() string {
	return fmt.Sprintf("# %s\n\n"+
		"GitHub Actions workflows using wetwire-github-go.\n\n"+
		"## Syntax\n\n"+
		"Use struct literals for all declarations:\n\n"+
		"```go\n"+
		"var CI = workflow.Workflow{\n"+
		"    Name: \"CI\",\n"+
		"    On:   CITriggers,\n"+
		"    Jobs: map[string]workflow.Job{\n"+
		"        \"build\": Build,\n"+
		"    },\n"+
		"}\n"+
		"```\n\n"+
		"## Building\n\n"+
		"```bash\n"+
		"wetwire-github build .\n"+
		"```\n", s.ProjectName)
}

// generateGitignore generates .gitignore content.
func (s *Scaffold) generateGitignore() string {
	return `# Binaries
*.exe
*.dll
*.so
*.dylib

# Test binary
*.test

# Output of the go coverage tool
*.out

# Dependency directories
vendor/

# IDE
.idea/
.vscode/
*.swp
*.swo
`
}

// WriteScaffold writes scaffold files to disk.
func WriteScaffold(outputDir string, files *ScaffoldFiles) error {
	for filename, content := range files.Files {
		path := filepath.Join(outputDir, filename)

		// Create directory if needed
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("creating directory %s: %w", dir, err)
		}

		// Write file
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("writing %s: %w", filename, err)
		}
	}

	return nil
}

// WriteGeneratedCode writes generated Go code to disk.
func WriteGeneratedCode(outputDir string, code *GeneratedCode) error {
	for filename, content := range code.Files {
		path := filepath.Join(outputDir, filename)

		// Skip empty files
		if strings.TrimSpace(content) == "" {
			continue
		}

		// Write file
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("writing %s: %w", filename, err)
		}
	}

	return nil
}
