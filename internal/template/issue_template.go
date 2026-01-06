package template

import (
	"github.com/lex00/wetwire-github-go/internal/discover"
	"github.com/lex00/wetwire-github-go/internal/runner"
	"github.com/lex00/wetwire-github-go/internal/serialize"
	"github.com/lex00/wetwire-github-go/templates"
)

// IssueTemplateBuildResult contains the result of building IssueTemplate templates.
type IssueTemplateBuildResult struct {
	// Templates contains the assembled templates with YAML output
	Templates []BuiltIssueTemplate

	// Errors contains any non-fatal errors encountered
	Errors []string
}

// BuiltIssueTemplate represents an IssueTemplate ready for output.
type BuiltIssueTemplate struct {
	// Name is the template variable name
	Name string

	// Template is the assembled IssueTemplate
	Template *templates.IssueTemplate

	// YAML is the serialized YAML output
	YAML []byte
}

// BuildIssueTemplates assembles IssueTemplate templates from discovery and extraction results.
func (b *Builder) BuildIssueTemplates(discovered *discover.IssueTemplateDiscoveryResult, extracted *runner.IssueTemplateExtractionResult) (*IssueTemplateBuildResult, error) {
	result := &IssueTemplateBuildResult{
		Templates: []BuiltIssueTemplate{},
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

		// Reconstruct the IssueTemplate from the map
		tmpl := b.reconstructIssueTemplate(templateData)

		// Serialize to YAML
		yaml, err := serialize.IssueTemplateToYAML(tmpl)
		if err != nil {
			result.Errors = append(result.Errors, "template "+dt.Name+": "+err.Error())
			continue
		}

		result.Templates = append(result.Templates, BuiltIssueTemplate{
			Name:     dt.Name,
			Template: tmpl,
			YAML:     yaml,
		})
	}

	return result, nil
}

// reconstructIssueTemplate reconstructs an IssueTemplate from a map.
func (b *Builder) reconstructIssueTemplate(data map[string]any) *templates.IssueTemplate {
	tmpl := &templates.IssueTemplate{}

	if v, ok := data["Name"].(string); ok {
		tmpl.Name = v
	}
	if v, ok := data["Description"].(string); ok {
		tmpl.Description = v
	}
	if v, ok := data["Title"].(string); ok {
		tmpl.Title = v
	}
	if v, ok := data["Labels"].([]any); ok {
		tmpl.Labels = anySliceToStrings(v)
	} else if v, ok := data["Labels"].([]string); ok {
		tmpl.Labels = v
	}
	if v, ok := data["Projects"].([]any); ok {
		tmpl.Projects = anySliceToStrings(v)
	} else if v, ok := data["Projects"].([]string); ok {
		tmpl.Projects = v
	}
	if v, ok := data["Assignees"].([]any); ok {
		tmpl.Assignees = anySliceToStrings(v)
	} else if v, ok := data["Assignees"].([]string); ok {
		tmpl.Assignees = v
	}

	if body, ok := data["Body"].([]any); ok {
		tmpl.Body = b.reconstructFormElements(body)
	} else if body, ok := data["Body"].([]templates.FormElement); ok {
		tmpl.Body = body
	}

	return tmpl
}

// reconstructFormElements reconstructs FormElements from a slice.
func (b *Builder) reconstructFormElements(data []any) []templates.FormElement {
	var elements []templates.FormElement

	for _, item := range data {
		if elem, ok := item.(templates.FormElement); ok {
			elements = append(elements, elem)
		} else if m, ok := item.(map[string]any); ok {
			// Try to determine element type from the map
			// This is a fallback for when type information is lost in JSON
			elem := b.reconstructFormElementFromMap(m)
			if elem != nil {
				elements = append(elements, elem)
			}
		} else if elem, ok := item.(templates.Markdown); ok {
			elements = append(elements, elem)
		} else if elem, ok := item.(templates.Input); ok {
			elements = append(elements, elem)
		} else if elem, ok := item.(templates.Textarea); ok {
			elements = append(elements, elem)
		} else if elem, ok := item.(templates.Dropdown); ok {
			elements = append(elements, elem)
		} else if elem, ok := item.(templates.Checkboxes); ok {
			elements = append(elements, elem)
		}
	}

	return elements
}

// reconstructFormElementFromMap attempts to reconstruct a FormElement from a generic map.
func (b *Builder) reconstructFormElementFromMap(data map[string]any) templates.FormElement {
	// Check for specific fields to determine the type
	if _, hasValue := data["Value"]; hasValue {
		// Could be Markdown, Input, or Textarea
		if _, hasLabel := data["Label"]; hasLabel {
			// Input or Textarea
			if _, hasRender := data["Render"]; hasRender {
				return b.reconstructTextarea(data)
			}
			return b.reconstructInput(data)
		}
		// Markdown
		return b.reconstructMarkdown(data)
	}
	if _, hasOptions := data["Options"]; hasOptions {
		// Dropdown or Checkboxes
		if opts, ok := data["Options"].([]any); ok && len(opts) > 0 {
			// Check if options are CheckboxOptions or strings
			if _, isMap := opts[0].(map[string]any); isMap {
				return b.reconstructCheckboxes(data)
			}
		}
		return b.reconstructDropdown(data)
	}
	if _, hasLabel := data["Label"]; hasLabel {
		// Could be Input, Textarea, or Dropdown without options yet
		if _, hasRender := data["Render"]; hasRender {
			return b.reconstructTextarea(data)
		}
		return b.reconstructInput(data)
	}

	return nil
}

// reconstructMarkdown reconstructs a Markdown element from a map.
func (b *Builder) reconstructMarkdown(data map[string]any) templates.Markdown {
	elem := templates.Markdown{}
	if v, ok := data["ID"].(string); ok {
		elem.ID = v
	}
	if v, ok := data["Value"].(string); ok {
		elem.Value = v
	}
	return elem
}

// reconstructInput reconstructs an Input element from a map.
func (b *Builder) reconstructInput(data map[string]any) templates.Input {
	elem := templates.Input{}
	if v, ok := data["ID"].(string); ok {
		elem.ID = v
	}
	if v, ok := data["Label"].(string); ok {
		elem.Label = v
	}
	if v, ok := data["Description"].(string); ok {
		elem.Description = v
	}
	if v, ok := data["Placeholder"].(string); ok {
		elem.Placeholder = v
	}
	if v, ok := data["Value"].(string); ok {
		elem.Value = v
	}
	if v, ok := data["Required"].(bool); ok {
		elem.Required = v
	}
	return elem
}

// reconstructTextarea reconstructs a Textarea element from a map.
func (b *Builder) reconstructTextarea(data map[string]any) templates.Textarea {
	elem := templates.Textarea{}
	if v, ok := data["ID"].(string); ok {
		elem.ID = v
	}
	if v, ok := data["Label"].(string); ok {
		elem.Label = v
	}
	if v, ok := data["Description"].(string); ok {
		elem.Description = v
	}
	if v, ok := data["Placeholder"].(string); ok {
		elem.Placeholder = v
	}
	if v, ok := data["Value"].(string); ok {
		elem.Value = v
	}
	if v, ok := data["Render"].(string); ok {
		elem.Render = v
	}
	if v, ok := data["Required"].(bool); ok {
		elem.Required = v
	}
	return elem
}

// reconstructDropdown reconstructs a Dropdown element from a map.
func (b *Builder) reconstructDropdown(data map[string]any) templates.Dropdown {
	elem := templates.Dropdown{}
	if v, ok := data["ID"].(string); ok {
		elem.ID = v
	}
	if v, ok := data["Label"].(string); ok {
		elem.Label = v
	}
	if v, ok := data["Description"].(string); ok {
		elem.Description = v
	}
	if v, ok := data["Options"].([]any); ok {
		elem.Options = anySliceToStrings(v)
	} else if v, ok := data["Options"].([]string); ok {
		elem.Options = v
	}
	if v, ok := data["Multiple"].(bool); ok {
		elem.Multiple = v
	}
	if v, ok := data["Default"].(int); ok {
		elem.Default = v
	} else if v, ok := data["Default"].(float64); ok {
		elem.Default = int(v)
	}
	if v, ok := data["Required"].(bool); ok {
		elem.Required = v
	}
	return elem
}

// reconstructCheckboxes reconstructs a Checkboxes element from a map.
func (b *Builder) reconstructCheckboxes(data map[string]any) templates.Checkboxes {
	elem := templates.Checkboxes{}
	if v, ok := data["ID"].(string); ok {
		elem.ID = v
	}
	if v, ok := data["Label"].(string); ok {
		elem.Label = v
	}
	if v, ok := data["Description"].(string); ok {
		elem.Description = v
	}
	if opts, ok := data["Options"].([]any); ok {
		for _, opt := range opts {
			if optMap, ok := opt.(map[string]any); ok {
				option := templates.CheckboxOption{}
				if v, ok := optMap["Label"].(string); ok {
					option.Label = v
				}
				if v, ok := optMap["Required"].(bool); ok {
					option.Required = v
				}
				elem.Options = append(elem.Options, option)
			}
		}
	}
	return elem
}
