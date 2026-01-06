package serialize

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/lex00/wetwire-github-go/templates"
)

// IssueTemplateToYAML serializes an IssueTemplate to YAML bytes.
func IssueTemplateToYAML(t *templates.IssueTemplate) ([]byte, error) {
	m, err := issueTemplateToMap(t)
	if err != nil {
		return nil, fmt.Errorf("converting issue template to map: %w", err)
	}

	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(m); err != nil {
		return nil, fmt.Errorf("encoding YAML: %w", err)
	}
	if err := encoder.Close(); err != nil {
		return nil, fmt.Errorf("closing encoder: %w", err)
	}

	return buf.Bytes(), nil
}

// issueTemplateToMap converts an IssueTemplate to a map for YAML serialization.
func issueTemplateToMap(t *templates.IssueTemplate) (map[string]any, error) {
	m := make(map[string]any)

	m["name"] = t.Name
	m["description"] = t.Description

	if t.Title != "" {
		m["title"] = t.Title
	}
	if len(t.Labels) > 0 {
		m["labels"] = t.Labels
	}
	if len(t.Projects) > 0 {
		m["projects"] = t.Projects
	}
	if len(t.Assignees) > 0 {
		m["assignees"] = t.Assignees
	}

	if len(t.Body) > 0 {
		body := make([]any, len(t.Body))
		for i, elem := range t.Body {
			elemMap, err := formElementToMap(elem)
			if err != nil {
				return nil, fmt.Errorf("serializing body element %d: %w", i, err)
			}
			body[i] = elemMap
		}
		m["body"] = body
	}

	return m, nil
}

// formElementToMap converts a FormElement to a map for YAML serialization.
func formElementToMap(elem templates.FormElement) (map[string]any, error) {
	m := make(map[string]any)
	m["type"] = elem.ElementType()

	switch e := elem.(type) {
	case templates.Markdown:
		if e.ID != "" {
			m["id"] = e.ID
		}
		m["attributes"] = map[string]any{
			"value": e.Value,
		}

	case templates.Input:
		if e.ID != "" {
			m["id"] = e.ID
		}
		attrs := make(map[string]any)
		attrs["label"] = e.Label
		if e.Description != "" {
			attrs["description"] = e.Description
		}
		if e.Placeholder != "" {
			attrs["placeholder"] = e.Placeholder
		}
		if e.Value != "" {
			attrs["value"] = e.Value
		}
		m["attributes"] = attrs
		if e.Required {
			m["validations"] = map[string]any{"required": true}
		}

	case templates.Textarea:
		if e.ID != "" {
			m["id"] = e.ID
		}
		attrs := make(map[string]any)
		attrs["label"] = e.Label
		if e.Description != "" {
			attrs["description"] = e.Description
		}
		if e.Placeholder != "" {
			attrs["placeholder"] = e.Placeholder
		}
		if e.Value != "" {
			attrs["value"] = e.Value
		}
		if e.Render != "" {
			attrs["render"] = e.Render
		}
		m["attributes"] = attrs
		if e.Required {
			m["validations"] = map[string]any{"required": true}
		}

	case templates.Dropdown:
		if e.ID != "" {
			m["id"] = e.ID
		}
		attrs := make(map[string]any)
		attrs["label"] = e.Label
		if e.Description != "" {
			attrs["description"] = e.Description
		}
		attrs["options"] = e.Options
		if e.Multiple {
			attrs["multiple"] = true
		}
		if e.Default > 0 {
			attrs["default"] = e.Default
		}
		m["attributes"] = attrs
		if e.Required {
			m["validations"] = map[string]any{"required": true}
		}

	case templates.Checkboxes:
		if e.ID != "" {
			m["id"] = e.ID
		}
		attrs := make(map[string]any)
		attrs["label"] = e.Label
		if e.Description != "" {
			attrs["description"] = e.Description
		}
		options := make([]map[string]any, len(e.Options))
		for i, opt := range e.Options {
			optMap := map[string]any{"label": opt.Label}
			if opt.Required {
				optMap["required"] = true
			}
			options[i] = optMap
		}
		attrs["options"] = options
		m["attributes"] = attrs

	default:
		return nil, fmt.Errorf("unknown form element type: %T", elem)
	}

	return m, nil
}
