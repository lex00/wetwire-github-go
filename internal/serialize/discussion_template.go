package serialize

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/lex00/wetwire-github-go/templates"
)

// DiscussionTemplateToYAML serializes a DiscussionTemplate to YAML bytes.
func DiscussionTemplateToYAML(t *templates.DiscussionTemplate) ([]byte, error) {
	m, err := discussionTemplateToMap(t)
	if err != nil {
		return nil, fmt.Errorf("converting discussion template to map: %w", err)
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

// discussionTemplateToMap converts a DiscussionTemplate to a map for YAML serialization.
func discussionTemplateToMap(t *templates.DiscussionTemplate) (map[string]any, error) {
	m := make(map[string]any)

	m["title"] = t.Title
	m["description"] = t.Description

	if len(t.Labels) > 0 {
		m["labels"] = t.Labels
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
