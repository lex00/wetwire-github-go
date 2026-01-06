package template

import (
	"github.com/lex00/wetwire-github-go/internal/discover"
	"github.com/lex00/wetwire-github-go/internal/runner"
	"github.com/lex00/wetwire-github-go/internal/serialize"
	"github.com/lex00/wetwire-github-go/templates"
)

// DiscussionTemplateBuildResult contains the result of building DiscussionTemplate templates.
type DiscussionTemplateBuildResult struct {
	// Templates contains the assembled templates with YAML output
	Templates []BuiltDiscussionTemplate

	// Errors contains any non-fatal errors encountered
	Errors []string
}

// BuiltDiscussionTemplate represents a DiscussionTemplate ready for output.
type BuiltDiscussionTemplate struct {
	// Name is the template variable name
	Name string

	// Template is the assembled DiscussionTemplate
	Template *templates.DiscussionTemplate

	// YAML is the serialized YAML output
	YAML []byte
}

// BuildDiscussionTemplates assembles DiscussionTemplate templates from discovery and extraction results.
func (b *Builder) BuildDiscussionTemplates(discovered *discover.DiscussionTemplateDiscoveryResult, extracted *runner.DiscussionTemplateExtractionResult) (*DiscussionTemplateBuildResult, error) {
	result := &DiscussionTemplateBuildResult{
		Templates: []BuiltDiscussionTemplate{},
		Errors:    []string{},
	}

	// Process each template
	for _, dt := range discovered.Templates {
		// Find the extracted template data
		var templateData map[string]any
		for _, et := range extracted.Templates {
			if et.Name == dt.Name {
				templateData = et.Data
				break
			}
		}

		if templateData == nil {
			result.Errors = append(result.Errors, "template "+dt.Name+": extraction data not found")
			continue
		}

		// Reconstruct the DiscussionTemplate from the map
		tmpl := b.reconstructDiscussionTemplate(templateData)

		// Serialize to YAML
		yaml, err := serialize.DiscussionTemplateToYAML(tmpl)
		if err != nil {
			result.Errors = append(result.Errors, "template "+dt.Name+": "+err.Error())
			continue
		}

		result.Templates = append(result.Templates, BuiltDiscussionTemplate{
			Name:     dt.Name,
			Template: tmpl,
			YAML:     yaml,
		})
	}

	return result, nil
}

// reconstructDiscussionTemplate reconstructs a DiscussionTemplate from a map.
func (b *Builder) reconstructDiscussionTemplate(data map[string]any) *templates.DiscussionTemplate {
	tmpl := &templates.DiscussionTemplate{}

	if v, ok := data["Title"].(string); ok {
		tmpl.Title = v
	}
	if v, ok := data["Description"].(string); ok {
		tmpl.Description = v
	}
	if v, ok := data["Labels"].([]any); ok {
		tmpl.Labels = anySliceToStrings(v)
	} else if v, ok := data["Labels"].([]string); ok {
		tmpl.Labels = v
	}

	// Reuse the form element reconstruction from issue templates
	if body, ok := data["Body"].([]any); ok {
		tmpl.Body = b.reconstructFormElements(body)
	} else if body, ok := data["Body"].([]templates.FormElement); ok {
		tmpl.Body = body
	}

	return tmpl
}
