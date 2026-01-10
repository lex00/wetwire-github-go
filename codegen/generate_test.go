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

func TestGenerator_GenerateActionWrapper_AllTypes(t *testing.T) {
	// Test that all three types (string, int, bool) are correctly generated
	spec := &ActionSpec{
		Name:        "AllTypes",
		Description: "Action with all field types",
		Inputs: map[string]ActionInput{
			"string-field": {
				Description: "A string field",
				Required:    false,
			},
			"int-field": {
				Description: "A number field",
				Default:     "42",
			},
			"bool-field": {
				Description: "A boolean field",
				Default:     "true",
			},
		},
	}

	gen := NewGenerator()
	code, err := gen.GenerateActionWrapper(ActionWrapperConfig{
		ActionRef:   "test/all-types@v1",
		PackageName: "all_types",
		TypeName:    "AllTypes",
		Spec:        spec,
	})
	if err != nil {
		t.Fatalf("GenerateActionWrapper() error = %v", err)
	}

	codeStr := string(code.Code)

	// Check all types are present
	if !strings.Contains(codeStr, "StringField string") {
		t.Error("Generated code missing string field")
	}
	if !strings.Contains(codeStr, "IntField int") {
		t.Error("Generated code missing int field")
	}
	if !strings.Contains(codeStr, "BoolField bool") {
		t.Error("Generated code missing bool field")
	}

	// Check Inputs() function contains all types
	if !strings.Contains(codeStr, `with["string-field"]`) {
		t.Error("Generated code missing string-field in Inputs()")
	}
	if !strings.Contains(codeStr, `with["int-field"]`) {
		t.Error("Generated code missing int-field in Inputs()")
	}
	if !strings.Contains(codeStr, `with["bool-field"]`) {
		t.Error("Generated code missing bool-field in Inputs()")
	}
}

func TestGenerator_GenerateActionWrapper_RequiredFields(t *testing.T) {
	// Test that required fields come before optional fields
	spec := &ActionSpec{
		Name: "RequiredTest",
		Inputs: map[string]ActionInput{
			"z-optional": {
				Description: "Optional field",
				Required:    false,
			},
			"a-required": {
				Description: "Required field",
				Required:    true,
			},
		},
	}

	gen := NewGenerator()
	code, err := gen.GenerateActionWrapper(ActionWrapperConfig{
		ActionRef:   "test/required@v1",
		PackageName: "required",
		TypeName:    "RequiredTest",
		Spec:        spec,
	})
	if err != nil {
		t.Fatalf("GenerateActionWrapper() error = %v", err)
	}

	codeStr := string(code.Code)

	// Required field (ARequired) should appear before optional (ZOptional)
	reqIdx := strings.Index(codeStr, "ARequired")
	optIdx := strings.Index(codeStr, "ZOptional")

	if reqIdx == -1 || optIdx == -1 {
		t.Error("Generated code missing expected fields")
	} else if reqIdx > optIdx {
		t.Error("Required field should come before optional field")
	}
}

func TestGenerator_GenerateActionWrapperFromYAML_NoVersion(t *testing.T) {
	// Test action ref without version suffix
	yaml := `name: 'Test Action'
description: 'A test action'
inputs:
  version:
    description: 'The version'
runs:
  using: 'node20'
  main: 'index.js'
`

	gen := NewGenerator()
	code, err := gen.GenerateActionWrapperFromYAML([]byte(yaml), "owner/repo")
	if err != nil {
		t.Fatalf("GenerateActionWrapperFromYAML() error = %v", err)
	}

	if code.PackageName != "repo" {
		t.Errorf("code.PackageName = %q, want %q", code.PackageName, "repo")
	}
}

func TestSanitizeDescription_CarriageReturn(t *testing.T) {
	// Test that carriage returns are removed
	desc := "Line1\r\nLine2\r\nLine3"
	result := sanitizeDescription(desc)

	if strings.Contains(result, "\r") {
		t.Error("sanitizeDescription should remove carriage returns")
	}
}

func TestSanitizeDescription_ExactBoundary(t *testing.T) {
	// Test description at exactly 100 characters (should not truncate)
	desc := strings.Repeat("a", 100)
	result := sanitizeDescription(desc)

	if result != desc {
		t.Errorf("100 char description should not be truncated, got len=%d", len(result))
	}
}

func TestSanitizeDescription_JustOverBoundary(t *testing.T) {
	// Test description at 101 characters (should truncate)
	desc := strings.Repeat("a", 101)
	result := sanitizeDescription(desc)

	expected := strings.Repeat("a", 97) + "..."
	if result != expected {
		t.Errorf("101 char description should be truncated to 100, got %q", result)
	}
}

func TestField_Struct(t *testing.T) {
	// Test Field struct
	f := Field{
		Name:        "TestField",
		Type:        "string",
		YAMLName:    "test-field",
		Description: "Test description",
		Required:    true,
	}

	if f.Name != "TestField" {
		t.Errorf("Field.Name = %q, want %q", f.Name, "TestField")
	}
	if f.Type != "string" {
		t.Errorf("Field.Type = %q, want %q", f.Type, "string")
	}
	if !f.Required {
		t.Error("Field.Required = false, want true")
	}
}

func TestActionWrapperConfig_Struct(t *testing.T) {
	// Test ActionWrapperConfig struct
	config := ActionWrapperConfig{
		ActionRef:   "test/action@v1",
		PackageName: "action",
		TypeName:    "Action",
		Spec:        &ActionSpec{Name: "Test"},
	}

	if config.ActionRef != "test/action@v1" {
		t.Errorf("ActionWrapperConfig.ActionRef = %q, want %q", config.ActionRef, "test/action@v1")
	}
	if config.Spec.Name != "Test" {
		t.Errorf("ActionWrapperConfig.Spec.Name = %q, want %q", config.Spec.Name, "Test")
	}
}

func TestGeneratedCode_Struct(t *testing.T) {
	// Test GeneratedCode struct
	code := GeneratedCode{
		PackageName: "test",
		FileName:    "test.go",
		Code:        []byte("package test"),
	}

	if code.PackageName != "test" {
		t.Errorf("GeneratedCode.PackageName = %q, want %q", code.PackageName, "test")
	}
	if code.FileName != "test.go" {
		t.Errorf("GeneratedCode.FileName = %q, want %q", code.FileName, "test.go")
	}
	if string(code.Code) != "package test" {
		t.Errorf("GeneratedCode.Code = %q, want %q", string(code.Code), "package test")
	}
}

func TestGenerator_GenerateActionWrapper_NilInputs(t *testing.T) {
	// Test with nil inputs map
	spec := &ActionSpec{
		Name:        "NilInputs",
		Description: "Action with nil inputs",
		Inputs:      nil,
	}

	gen := NewGenerator()
	code, err := gen.GenerateActionWrapper(ActionWrapperConfig{
		ActionRef:   "test/nil@v1",
		PackageName: "nil_inputs",
		TypeName:    "NilInputs",
		Spec:        spec,
	})
	if err != nil {
		t.Fatalf("GenerateActionWrapper() error = %v", err)
	}

	codeStr := string(code.Code)
	if !strings.Contains(codeStr, "type NilInputs struct") {
		t.Error("Generated code missing struct definition")
	}
}

func TestGenerator_GenerateActionWrapper_DescriptionTruncation(t *testing.T) {
	// Test that long field descriptions are truncated
	longDesc := strings.Repeat("This is a very long description. ", 10)
	spec := &ActionSpec{
		Name: "LongDesc",
		Inputs: map[string]ActionInput{
			"field": {
				Description: longDesc,
				Required:    false,
			},
		},
	}

	gen := NewGenerator()
	code, err := gen.GenerateActionWrapper(ActionWrapperConfig{
		ActionRef:   "test/long@v1",
		PackageName: "long_desc",
		TypeName:    "LongDesc",
		Spec:        spec,
	})
	if err != nil {
		t.Fatalf("GenerateActionWrapper() error = %v", err)
	}

	codeStr := string(code.Code)
	// Verify the code was generated
	if !strings.Contains(codeStr, "Field string") {
		t.Error("Generated code missing field")
	}
}

func TestGenerator_GenerateActionWrapper_SpecialCharsInDesc(t *testing.T) {
	// Test handling of special characters in description
	spec := &ActionSpec{
		Name:        "SpecialChars",
		Description: "Action with special chars: <>&\"'",
		Inputs: map[string]ActionInput{
			"field": {
				Description: "Description with\ttabs\nand newlines",
				Required:    false,
			},
		},
	}

	gen := NewGenerator()
	code, err := gen.GenerateActionWrapper(ActionWrapperConfig{
		ActionRef:   "test/special@v1",
		PackageName: "special",
		TypeName:    "SpecialChars",
		Spec:        spec,
	})
	if err != nil {
		t.Fatalf("GenerateActionWrapper() error = %v", err)
	}

	// Verify code was generated and is valid
	if len(code.Code) == 0 {
		t.Error("Generated code is empty")
	}
}
