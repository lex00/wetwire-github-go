package template

import (
	"github.com/lex00/wetwire-github-go/internal/discover"
	"github.com/lex00/wetwire-github-go/internal/runner"
	"github.com/lex00/wetwire-github-go/templates"
)

// PRTemplateBuildResult contains the result of building PRTemplate templates.
type PRTemplateBuildResult struct {
	// Templates contains the assembled templates with markdown output
	Templates []BuiltPRTemplate

	// Errors contains any non-fatal errors encountered
	Errors []string
}

// BuiltPRTemplate represents a PRTemplate ready for output.
type BuiltPRTemplate struct {
	// Name is the template variable name
	Name string

	// Template is the assembled PRTemplate
	Template *templates.PRTemplate

	// Filename is the output filename (e.g., "PULL_REQUEST_TEMPLATE.md")
	Filename string

	// Content is the markdown content
	Content []byte
}

// BuildPRTemplates assembles PRTemplate templates from discovery and extraction results.
func (b *Builder) BuildPRTemplates(discovered *discover.PRTemplateDiscoveryResult, extracted *runner.PRTemplateExtractionResult) (*PRTemplateBuildResult, error) {
	result := &PRTemplateBuildResult{
		Templates: []BuiltPRTemplate{},
		Errors:    []string{},
	}

	// Process each template
	for _, dt := range discovered.Templates {
		// Find the extracted template data
		var content string
		var found bool
		for _, et := range extracted.Templates {
			if et.Name == dt.Name {
				content = et.Content
				found = true
				break
			}
		}

		if !found {
			result.Errors = append(result.Errors, "template "+dt.Name+": extraction data not found")
			continue
		}

		// Create the PRTemplate
		tmpl := &templates.PRTemplate{
			Name:    dt.Name,
			Content: content,
		}

		result.Templates = append(result.Templates, BuiltPRTemplate{
			Name:     dt.Name,
			Template: tmpl,
			Filename: tmpl.Filename(),
			Content:  []byte(content),
		})
	}

	return result, nil
}
