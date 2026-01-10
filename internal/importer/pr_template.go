package importer

import (
	"fmt"
	"sort"
	"strings"
)

// IRPRTemplate represents a parsed PR template.
type IRPRTemplate struct {
	Name    string // Template name (from filename or explicit)
	Content string // Raw markdown content
}

// ParsePRTemplateContent parses PR template content (plain markdown).
func ParsePRTemplateContent(name, content string) (*IRPRTemplate, error) {
	return &IRPRTemplate{
		Name:    name,
		Content: content,
	}, nil
}

// PRTemplateCodeGenerator generates Go code from parsed PR templates.
type PRTemplateCodeGenerator struct {
	PackageName string
}

// PRTemplateGeneratedCode contains the generated Go code.
type PRTemplateGeneratedCode struct {
	Files     map[string]string // filename -> content
	Templates int
}

// Generate generates Go code from PR templates.
func (g *PRTemplateCodeGenerator) Generate(templates map[string]*IRPRTemplate) (*PRTemplateGeneratedCode, error) {
	result := &PRTemplateGeneratedCode{
		Files:     make(map[string]string),
		Templates: len(templates),
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("package %s\n\n", g.PackageName))
	sb.WriteString("import (\n")
	sb.WriteString("\t\"github.com/lex00/wetwire-github-go/templates\"\n")
	sb.WriteString(")\n\n")

	// Sort for deterministic output
	var names []string
	for name := range templates {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		template := templates[name]
		sb.WriteString(g.generatePRTemplate(name, template))
		sb.WriteString("\n")
	}

	result.Files["templates.go"] = sb.String()
	return result, nil
}

func (g *PRTemplateCodeGenerator) generatePRTemplate(varName string, template *IRPRTemplate) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("var %s = templates.PRTemplate{\n", varName))

	if template.Name != "" {
		sb.WriteString(fmt.Sprintf("\tName: %q,\n", template.Name))
	}

	sb.WriteString(fmt.Sprintf("\tContent: %s,\n", formatPRTemplateContent(template.Content)))
	sb.WriteString("}\n")

	return sb.String()
}

// formatPRTemplateContent formats multiline content as a raw string literal.
func formatPRTemplateContent(content string) string {
	// Use backticks for multiline strings if content doesn't contain backticks
	if !strings.Contains(content, "`") {
		return "`" + content + "`"
	}

	// Fall back to quoted string with escaping
	return fmt.Sprintf("%q", content)
}
