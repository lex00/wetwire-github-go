package codegen

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Intermediate Representation Types

// ActionSpec represents a parsed action.yml file.
type ActionSpec struct {
	Name        string              `yaml:"name" json:"name"`
	Description string              `yaml:"description" json:"description"`
	Author      string              `yaml:"author,omitempty" json:"author,omitempty"`
	Inputs      map[string]ActionInput  `yaml:"inputs,omitempty" json:"inputs,omitempty"`
	Outputs     map[string]ActionOutput `yaml:"outputs,omitempty" json:"outputs,omitempty"`
	Runs        ActionRuns          `yaml:"runs" json:"runs"`
	Branding    *ActionBranding     `yaml:"branding,omitempty" json:"branding,omitempty"`
}

// ActionInput represents an input parameter for an action.
type ActionInput struct {
	Description        string `yaml:"description" json:"description"`
	Required           bool   `yaml:"required,omitempty" json:"required,omitempty"`
	Default            string `yaml:"default,omitempty" json:"default,omitempty"`
	DeprecationMessage string `yaml:"deprecationMessage,omitempty" json:"deprecation_message,omitempty"`
}

// ActionOutput represents an output from an action.
type ActionOutput struct {
	Description string `yaml:"description" json:"description"`
	Value       string `yaml:"value,omitempty" json:"value,omitempty"`
}

// ActionRuns defines how the action runs.
type ActionRuns struct {
	Using          string            `yaml:"using" json:"using"`
	Main           string            `yaml:"main,omitempty" json:"main,omitempty"`
	Pre            string            `yaml:"pre,omitempty" json:"pre,omitempty"`
	PreIf          string            `yaml:"pre-if,omitempty" json:"pre_if,omitempty"`
	Post           string            `yaml:"post,omitempty" json:"post,omitempty"`
	PostIf         string            `yaml:"post-if,omitempty" json:"post_if,omitempty"`
	Image          string            `yaml:"image,omitempty" json:"image,omitempty"`
	Env            map[string]string `yaml:"env,omitempty" json:"env,omitempty"`
	Args           []string          `yaml:"args,omitempty" json:"args,omitempty"`
	Steps          []any             `yaml:"steps,omitempty" json:"steps,omitempty"` // For composite actions
}

// ActionBranding defines the branding for the action in GitHub marketplace.
type ActionBranding struct {
	Icon  string `yaml:"icon,omitempty" json:"icon,omitempty"`
	Color string `yaml:"color,omitempty" json:"color,omitempty"`
}

// WorkflowSchema represents a parsed JSON schema for workflows.
type WorkflowSchema struct {
	Schema      string                       `json:"$schema,omitempty"`
	ID          string                       `json:"$id,omitempty"`
	Title       string                       `json:"title,omitempty"`
	Description string                       `json:"description,omitempty"`
	Type        string                       `json:"type,omitempty"`
	Properties  map[string]SchemaProperty    `json:"properties,omitempty"`
	Definitions map[string]SchemaDefinition  `json:"definitions,omitempty"`
	Required    []string                     `json:"required,omitempty"`
}

// SchemaProperty represents a property in a JSON schema.
type SchemaProperty struct {
	Type        any                       `json:"type,omitempty"` // string or []string
	Description string                       `json:"description,omitempty"`
	Enum        []string                     `json:"enum,omitempty"`
	Default     any                          `json:"default,omitempty"`
	Ref         string                       `json:"$ref,omitempty"`
	Items       *SchemaProperty              `json:"items,omitempty"`
	Properties  map[string]SchemaProperty    `json:"properties,omitempty"`
	Required    []string                     `json:"required,omitempty"`
	OneOf       []SchemaProperty             `json:"oneOf,omitempty"`
	AnyOf       []SchemaProperty             `json:"anyOf,omitempty"`
	AllOf       []SchemaProperty             `json:"allOf,omitempty"`
	Pattern     string                       `json:"pattern,omitempty"`
	MinItems    *int                         `json:"minItems,omitempty"`
	MaxItems    *int                         `json:"maxItems,omitempty"`
	MinLength   *int                         `json:"minLength,omitempty"`
	MaxLength   *int                         `json:"maxLength,omitempty"`
	Minimum     *float64                     `json:"minimum,omitempty"`
	Maximum     *float64                     `json:"maximum,omitempty"`
}

// SchemaDefinition represents a definition in a JSON schema.
type SchemaDefinition struct {
	SchemaProperty
	Title string `json:"title,omitempty"`
}

// Parsers

// ParseActionYAML parses an action.yml file content into an ActionSpec.
func ParseActionYAML(data []byte) (*ActionSpec, error) {
	var action ActionSpec
	if err := yaml.Unmarshal(data, &action); err != nil {
		return nil, fmt.Errorf("parsing action.yml: %w", err)
	}
	return &action, nil
}

// ParseWorkflowSchema parses a JSON schema for workflows.
func ParseWorkflowSchema(data []byte) (*WorkflowSchema, error) {
	var schema WorkflowSchema
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, fmt.Errorf("parsing workflow schema: %w", err)
	}
	return &schema, nil
}

// Utility functions for working with parsed schemas

// GetInputNames returns a sorted list of input names from an action spec.
func (a *ActionSpec) GetInputNames() []string {
	names := make([]string, 0, len(a.Inputs))
	for name := range a.Inputs {
		names = append(names, name)
	}
	return names
}

// GetRequiredInputs returns a list of required input names.
func (a *ActionSpec) GetRequiredInputs() []string {
	var required []string
	for name, input := range a.Inputs {
		if input.Required {
			required = append(required, name)
		}
	}
	return required
}

// GetOutputNames returns a list of output names from an action spec.
func (a *ActionSpec) GetOutputNames() []string {
	names := make([]string, 0, len(a.Outputs))
	for name := range a.Outputs {
		names = append(names, name)
	}
	return names
}

// IsCompositeAction returns true if this is a composite action.
func (a *ActionSpec) IsCompositeAction() bool {
	return a.Runs.Using == "composite"
}

// IsJavaScriptAction returns true if this is a JavaScript/Node action.
func (a *ActionSpec) IsJavaScriptAction() bool {
	return strings.HasPrefix(a.Runs.Using, "node")
}

// IsDockerAction returns true if this is a Docker action.
func (a *ActionSpec) IsDockerAction() bool {
	return a.Runs.Using == "docker"
}

// GetGoFieldName converts a kebab-case or snake_case input name to Go field name.
func GetGoFieldName(name string) string {
	// Split on hyphens and underscores
	parts := strings.FieldsFunc(name, func(r rune) bool {
		return r == '-' || r == '_'
	})

	// Title case each part
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
		}
	}

	return strings.Join(parts, "")
}

// GetGoTypeName converts an action name to a Go type name.
func GetGoTypeName(name string) string {
	// Remove version suffix like @v4
	if idx := strings.Index(name, "@"); idx != -1 {
		name = name[:idx]
	}

	// Convert to title case
	return GetGoFieldName(name)
}

// ResolveRef resolves a $ref in a JSON schema.
func (s *WorkflowSchema) ResolveRef(ref string) (*SchemaDefinition, error) {
	// Refs are typically like "#/definitions/normalJob"
	if !strings.HasPrefix(ref, "#/definitions/") {
		return nil, fmt.Errorf("unsupported ref format: %s", ref)
	}

	defName := strings.TrimPrefix(ref, "#/definitions/")
	def, ok := s.Definitions[defName]
	if !ok {
		return nil, fmt.Errorf("definition not found: %s", defName)
	}

	return &def, nil
}

// GetPropertyType returns a string representation of the property type.
func (p *SchemaProperty) GetPropertyType() string {
	switch t := p.Type.(type) {
	case string:
		return t
	case []any:
		types := make([]string, len(t))
		for i, v := range t {
			types[i] = fmt.Sprint(v)
		}
		return strings.Join(types, "|")
	default:
		if p.Ref != "" {
			return "ref:" + p.Ref
		}
		return "unknown"
	}
}
