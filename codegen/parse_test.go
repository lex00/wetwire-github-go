package codegen

import (
	"testing"
)

func TestParseActionYAML(t *testing.T) {
	yaml := `name: 'Checkout'
description: 'Checkout a Git repository at a particular version'
author: 'GitHub'
inputs:
  repository:
    description: 'Repository name with owner. For example, actions/checkout'
    default: ${{ github.repository }}
  ref:
    description: 'The branch, tag or SHA to checkout.'
    required: false
  token:
    description: 'Personal access token'
    default: ${{ github.token }}
    required: true
  fetch-depth:
    description: 'Number of commits to fetch. 0 indicates all history.'
    default: '1'
outputs:
  ref:
    description: 'The branch, tag or SHA that was checked out'
  commit:
    description: 'The commit SHA that was checked out'
runs:
  using: 'node20'
  main: 'dist/index.js'
  post: 'dist/index.js'
branding:
  icon: 'check-circle'
  color: 'green'
`

	action, err := ParseActionYAML([]byte(yaml))
	if err != nil {
		t.Fatalf("ParseActionYAML() error = %v", err)
	}

	if action.Name != "Checkout" {
		t.Errorf("action.Name = %q, want %q", action.Name, "Checkout")
	}

	if action.Author != "GitHub" {
		t.Errorf("action.Author = %q, want %q", action.Author, "GitHub")
	}

	if len(action.Inputs) != 4 {
		t.Errorf("len(action.Inputs) = %d, want 4", len(action.Inputs))
	}

	if input, ok := action.Inputs["fetch-depth"]; !ok {
		t.Error("action.Inputs missing 'fetch-depth'")
	} else if input.Default != "1" {
		t.Errorf("action.Inputs['fetch-depth'].Default = %q, want %q", input.Default, "1")
	}

	if input, ok := action.Inputs["token"]; !ok {
		t.Error("action.Inputs missing 'token'")
	} else if !input.Required {
		t.Error("action.Inputs['token'].Required = false, want true")
	}

	if len(action.Outputs) != 2 {
		t.Errorf("len(action.Outputs) = %d, want 2", len(action.Outputs))
	}

	if action.Runs.Using != "node20" {
		t.Errorf("action.Runs.Using = %q, want %q", action.Runs.Using, "node20")
	}

	if action.Branding == nil {
		t.Error("action.Branding is nil")
	} else if action.Branding.Icon != "check-circle" {
		t.Errorf("action.Branding.Icon = %q, want %q", action.Branding.Icon, "check-circle")
	}
}

func TestParseActionYAML_Composite(t *testing.T) {
	yaml := `name: 'My Composite Action'
description: 'A composite action'
runs:
  using: 'composite'
  steps:
    - run: echo "Hello"
      shell: bash
`

	action, err := ParseActionYAML([]byte(yaml))
	if err != nil {
		t.Fatalf("ParseActionYAML() error = %v", err)
	}

	if !action.IsCompositeAction() {
		t.Error("IsCompositeAction() = false, want true")
	}

	if action.IsJavaScriptAction() {
		t.Error("IsJavaScriptAction() = true, want false")
	}

	if action.IsDockerAction() {
		t.Error("IsDockerAction() = true, want false")
	}
}

func TestParseActionYAML_Docker(t *testing.T) {
	yaml := `name: 'Docker Action'
description: 'A docker action'
runs:
  using: 'docker'
  image: 'Dockerfile'
`

	action, err := ParseActionYAML([]byte(yaml))
	if err != nil {
		t.Fatalf("ParseActionYAML() error = %v", err)
	}

	if !action.IsDockerAction() {
		t.Error("IsDockerAction() = false, want true")
	}
}

func TestParseWorkflowSchema(t *testing.T) {
	schema := `{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "$id": "https://json.schemastore.org/github-workflow.json",
  "title": "GitHub Workflow",
  "description": "A GitHub Workflow configuration",
  "type": "object",
  "properties": {
    "name": {
      "type": "string",
      "description": "The name of your workflow"
    },
    "on": {
      "type": ["string", "array", "object"],
      "description": "The event that triggers the workflow"
    }
  },
  "definitions": {
    "normalJob": {
      "title": "Normal Job",
      "type": "object",
      "properties": {
        "runs-on": {
          "type": "string"
        }
      }
    }
  },
  "required": ["on"]
}`

	ws, err := ParseWorkflowSchema([]byte(schema))
	if err != nil {
		t.Fatalf("ParseWorkflowSchema() error = %v", err)
	}

	if ws.Title != "GitHub Workflow" {
		t.Errorf("schema.Title = %q, want %q", ws.Title, "GitHub Workflow")
	}

	if ws.Type != "object" {
		t.Errorf("schema.Type = %q, want %q", ws.Type, "object")
	}

	if len(ws.Properties) != 2 {
		t.Errorf("len(schema.Properties) = %d, want 2", len(ws.Properties))
	}

	if prop, ok := ws.Properties["name"]; !ok {
		t.Error("schema.Properties missing 'name'")
	} else if prop.GetPropertyType() != "string" {
		t.Errorf("schema.Properties['name'].GetPropertyType() = %q, want %q", prop.GetPropertyType(), "string")
	}

	if len(ws.Required) != 1 || ws.Required[0] != "on" {
		t.Errorf("schema.Required = %v, want [on]", ws.Required)
	}

	// Test definition resolution
	def, err := ws.ResolveRef("#/definitions/normalJob")
	if err != nil {
		t.Errorf("ResolveRef() error = %v", err)
	}
	if def.Title != "Normal Job" {
		t.Errorf("def.Title = %q, want %q", def.Title, "Normal Job")
	}

	// Test invalid ref
	_, err = ws.ResolveRef("#/invalid/ref")
	if err == nil {
		t.Error("ResolveRef() expected error for invalid ref format")
	}

	_, err = ws.ResolveRef("#/definitions/notFound")
	if err == nil {
		t.Error("ResolveRef() expected error for missing definition")
	}
}

func TestGetGoFieldName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"fetch-depth", "FetchDepth"},
		{"go_version", "GoVersion"},
		{"repository", "Repository"},
		{"pre-if", "PreIf"},
		{"my-long-field-name", "MyLongFieldName"},
		{"UPPERCASE", "Uppercase"},
		{"", ""},
	}

	for _, tt := range tests {
		got := GetGoFieldName(tt.input)
		if got != tt.want {
			t.Errorf("GetGoFieldName(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestGetGoTypeName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"checkout", "Checkout"},
		{"setup-go", "SetupGo"},
		{"checkout@v4", "Checkout"}, // Version suffix stripped
		{"upload-artifact", "UploadArtifact"},
	}

	for _, tt := range tests {
		got := GetGoTypeName(tt.input)
		if got != tt.want {
			t.Errorf("GetGoTypeName(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestActionSpec_GetInputNames(t *testing.T) {
	action := &ActionSpec{
		Inputs: map[string]ActionInput{
			"repository": {},
			"ref":        {},
			"token":      {},
		},
	}

	names := action.GetInputNames()
	if len(names) != 3 {
		t.Errorf("len(GetInputNames()) = %d, want 3", len(names))
	}
}

func TestActionSpec_GetRequiredInputs(t *testing.T) {
	action := &ActionSpec{
		Inputs: map[string]ActionInput{
			"repository": {Required: false},
			"ref":        {Required: false},
			"token":      {Required: true},
		},
	}

	required := action.GetRequiredInputs()
	if len(required) != 1 {
		t.Errorf("len(GetRequiredInputs()) = %d, want 1", len(required))
	}
	if len(required) > 0 && required[0] != "token" {
		t.Errorf("GetRequiredInputs()[0] = %q, want %q", required[0], "token")
	}
}

func TestActionSpec_GetOutputNames(t *testing.T) {
	action := &ActionSpec{
		Outputs: map[string]ActionOutput{
			"ref":    {},
			"commit": {},
		},
	}

	names := action.GetOutputNames()
	if len(names) != 2 {
		t.Errorf("len(GetOutputNames()) = %d, want 2", len(names))
	}
}

func TestSchemaProperty_GetPropertyType(t *testing.T) {
	tests := []struct {
		name string
		prop SchemaProperty
		want string
	}{
		{
			name: "string type",
			prop: SchemaProperty{Type: "string"},
			want: "string",
		},
		{
			name: "array of types",
			prop: SchemaProperty{Type: []any{"string", "number"}},
			want: "string|number",
		},
		{
			name: "ref type",
			prop: SchemaProperty{Ref: "#/definitions/job"},
			want: "ref:#/definitions/job",
		},
		{
			name: "unknown type",
			prop: SchemaProperty{},
			want: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.prop.GetPropertyType()
			if got != tt.want {
				t.Errorf("GetPropertyType() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseActionYAML_Invalid(t *testing.T) {
	_, err := ParseActionYAML([]byte("invalid: yaml: content: ["))
	if err == nil {
		t.Error("ParseActionYAML() expected error for invalid YAML")
	}
}

func TestParseWorkflowSchema_Invalid(t *testing.T) {
	_, err := ParseWorkflowSchema([]byte("invalid json"))
	if err == nil {
		t.Error("ParseWorkflowSchema() expected error for invalid JSON")
	}
}
