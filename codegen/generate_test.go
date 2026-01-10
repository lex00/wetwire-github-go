package codegen

import (
	"strings"
	"testing"
)

func TestGenerator_GenerateActionWrapper(t *testing.T) {
	spec := &ActionSpec{
		Name:        "Checkout",
		Description: "Checkout a Git repository at a particular version",
		Inputs: map[string]ActionInput{
			"repository": {
				Description: "Repository name with owner",
				Required:    false,
			},
			"ref": {
				Description: "The branch, tag or SHA to checkout",
				Required:    false,
			},
			"token": {
				Description: "Personal access token",
				Required:    true,
			},
			"fetch-depth": {
				Description: "Number of commits to fetch",
				Default:     "1",
			},
		},
	}

	gen := NewGenerator()
	code, err := gen.GenerateActionWrapper(ActionWrapperConfig{
		ActionRef:   "actions/checkout@v4",
		PackageName: "checkout",
		TypeName:    "Checkout",
		Spec:        spec,
	})
	if err != nil {
		t.Fatalf("GenerateActionWrapper() error = %v", err)
	}

	// Check package name
	if code.PackageName != "checkout" {
		t.Errorf("code.PackageName = %q, want %q", code.PackageName, "checkout")
	}

	// Check file name
	if code.FileName != "checkout.go" {
		t.Errorf("code.FileName = %q, want %q", code.FileName, "checkout.go")
	}

	// Check generated code contains expected elements
	codeStr := string(code.Code)

	expectedStrings := []string{
		"package checkout",
		"type Checkout struct",
		`func (a Checkout) Action() string`,
		`return "actions/checkout@v4"`,
		"func (a Checkout) Inputs() map[string]any",
		"Token string", // Required field
		"FetchDepth int",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(codeStr, expected) {
			t.Errorf("Generated code missing %q\n\nGenerated:\n%s", expected, codeStr)
		}
	}
}

func TestGenerator_GenerateActionWrapperFromYAML(t *testing.T) {
	yaml := `name: 'Setup Go environment'
description: 'Setup a Go environment and add it to PATH'
inputs:
  go-version:
    description: 'The Go version to download (if necessary) and use'
    required: false
  go-version-file:
    description: 'Path to the go.mod file'
    required: false
  cache:
    description: 'Used to specify whether caching is needed'
    default: true
runs:
  using: 'node20'
  main: 'dist/setup/index.js'
`

	gen := NewGenerator()
	code, err := gen.GenerateActionWrapperFromYAML([]byte(yaml), "actions/setup-go@v5")
	if err != nil {
		t.Fatalf("GenerateActionWrapperFromYAML() error = %v", err)
	}

	if code.PackageName != "setup_go" {
		t.Errorf("code.PackageName = %q, want %q", code.PackageName, "setup_go")
	}

	codeStr := string(code.Code)

	expectedStrings := []string{
		"package setup_go",
		"type SetupGo struct",
		`return "actions/setup-go@v5"`,
		"GoVersion string",
		"GoVersionFile string",
		"Cache bool",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(codeStr, expected) {
			t.Errorf("Generated code missing %q\n\nGenerated:\n%s", expected, codeStr)
		}
	}
}

func TestInferGoType(t *testing.T) {
	tests := []struct {
		name  string
		input ActionInput
		want  string
	}{
		{
			name:  "default true",
			input: ActionInput{Default: "true"},
			want:  "bool",
		},
		{
			name:  "default false",
			input: ActionInput{Default: "false"},
			want:  "bool",
		},
		{
			name:  "default number",
			input: ActionInput{Default: "1"},
			want:  "int",
		},
		{
			name:  "description with number",
			input: ActionInput{Description: "Number of retries"},
			want:  "int",
		},
		{
			name:  "description with count",
			input: ActionInput{Description: "The count of items"},
			want:  "int",
		},
		{
			name:  "description with depth",
			input: ActionInput{Description: "Fetch depth"},
			want:  "int",
		},
		{
			name:  "description with timeout",
			input: ActionInput{Description: "Connection timeout in seconds"},
			want:  "int",
		},
		{
			name:  "description with whether",
			input: ActionInput{Description: "Whether to use cache"},
			want:  "bool",
		},
		{
			name:  "description with enable",
			input: ActionInput{Description: "Enable strict mode"},
			want:  "bool",
		},
		{
			name:  "description with disable",
			input: ActionInput{Description: "Disable logging"},
			want:  "bool",
		},
		{
			name:  "default string",
			input: ActionInput{Default: "hello"},
			want:  "string",
		},
		{
			name:  "no hints",
			input: ActionInput{Description: "Some description"},
			want:  "string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := inferGoType(tt.input)
			if got != tt.want {
				t.Errorf("inferGoType() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIsNumericDefault(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"1", true},
		{"123", true},
		{"0", true},
		{"", false},
		{"abc", false},
		{"1.5", false},
		{"-1", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := isNumericDefault(tt.input)
			if got != tt.want {
				t.Errorf("isNumericDefault(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestSanitizeDescription(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple",
			input: "A simple description",
			want:  "A simple description",
		},
		{
			name:  "newlines",
			input: "Line 1\nLine 2\nLine 3",
			want:  "Line 1 Line 2 Line 3",
		},
		{
			name:  "extra spaces",
			input: "Too   many    spaces",
			want:  "Too many spaces",
		},
		{
			name:  "long description",
			input: strings.Repeat("a", 200),
			want:  strings.Repeat("a", 97) + "...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeDescription(tt.input)
			if got != tt.want {
				t.Errorf("sanitizeDescription() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGenerator_GenerateActionWrapperFromYAML_InvalidRef(t *testing.T) {
	gen := NewGenerator()
	_, err := gen.GenerateActionWrapperFromYAML([]byte("name: test\nruns:\n  using: node20\n  main: index.js"), "invalid")
	if err == nil {
		t.Error("GenerateActionWrapperFromYAML() expected error for invalid ref")
	}
}

func TestGenerator_GenerateActionWrapperFromYAML_InvalidYAML(t *testing.T) {
	gen := NewGenerator()
	_, err := gen.GenerateActionWrapperFromYAML([]byte("invalid: yaml: ["), "actions/test@v1")
	if err == nil {
		t.Error("GenerateActionWrapperFromYAML() expected error for invalid YAML")
	}
}

func TestNewGenerator(t *testing.T) {
	gen := NewGenerator()
	if gen == nil {
		t.Error("NewGenerator() returned nil")
	}
}

func TestGenerator_EmptyInputs(t *testing.T) {
	spec := &ActionSpec{
		Name:        "Empty",
		Description: "An action with no inputs",
		Inputs:      map[string]ActionInput{},
	}

	gen := NewGenerator()
	code, err := gen.GenerateActionWrapper(ActionWrapperConfig{
		ActionRef:   "test/empty@v1",
		PackageName: "empty",
		TypeName:    "Empty",
		Spec:        spec,
	})
	if err != nil {
		t.Fatalf("GenerateActionWrapper() error = %v", err)
	}

	codeStr := string(code.Code)
	if !strings.Contains(codeStr, "type Empty struct {") {
		t.Errorf("Generated code missing empty struct")
	}
}
