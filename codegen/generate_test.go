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

func TestGenerator_GenerateActionWrapper_AllRequiredFields(t *testing.T) {
	// Test with all required fields (no optional fields)
	spec := &ActionSpec{
		Name:        "AllRequired",
		Description: "Action with all required fields",
		Inputs: map[string]ActionInput{
			"alpha": {
				Description: "First required field",
				Required:    true,
			},
			"beta": {
				Description: "Second required field",
				Required:    true,
			},
			"gamma": {
				Description: "Third required field",
				Required:    true,
			},
		},
	}

	gen := NewGenerator()
	code, err := gen.GenerateActionWrapper(ActionWrapperConfig{
		ActionRef:   "test/all-required@v1",
		PackageName: "all_required",
		TypeName:    "AllRequired",
		Spec:        spec,
	})
	if err != nil {
		t.Fatalf("GenerateActionWrapper() error = %v", err)
	}

	codeStr := string(code.Code)

	// Verify all fields are present and alphabetically sorted
	alphaIdx := strings.Index(codeStr, "Alpha string")
	betaIdx := strings.Index(codeStr, "Beta string")
	gammaIdx := strings.Index(codeStr, "Gamma string")

	if alphaIdx == -1 || betaIdx == -1 || gammaIdx == -1 {
		t.Error("Generated code missing expected fields")
	}

	// Since all are required, they should be alphabetically sorted
	if alphaIdx > betaIdx || betaIdx > gammaIdx {
		t.Error("Required fields should be alphabetically sorted")
	}
}

func TestGenerator_GenerateActionWrapper_NoRequiredFields(t *testing.T) {
	// Test with no required fields (all optional)
	spec := &ActionSpec{
		Name:        "NoRequired",
		Description: "Action with no required fields",
		Inputs: map[string]ActionInput{
			"zebra": {
				Description: "Optional field z",
				Required:    false,
			},
			"alpha": {
				Description: "Optional field a",
				Required:    false,
			},
		},
	}

	gen := NewGenerator()
	code, err := gen.GenerateActionWrapper(ActionWrapperConfig{
		ActionRef:   "test/no-required@v1",
		PackageName: "no_required",
		TypeName:    "NoRequired",
		Spec:        spec,
	})
	if err != nil {
		t.Fatalf("GenerateActionWrapper() error = %v", err)
	}

	codeStr := string(code.Code)

	// Verify fields are alphabetically sorted (since all are optional)
	alphaIdx := strings.Index(codeStr, "Alpha string")
	zebraIdx := strings.Index(codeStr, "Zebra string")

	if alphaIdx == -1 || zebraIdx == -1 {
		t.Error("Generated code missing expected fields")
	}

	if alphaIdx > zebraIdx {
		t.Error("Optional fields should be alphabetically sorted")
	}
}

func TestSanitizeDescription_WithSpecialCharacters(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "with backticks",
			input: "Use `command` for execution",
			want:  "Use `command` for execution",
		},
		{
			name:  "with quotes",
			input: `Say "hello" world`,
			want:  `Say "hello" world`,
		},
		{
			name:  "with angle brackets",
			input: "Use <input> tag",
			want:  "Use <input> tag",
		},
		{
			name:  "mixed special chars",
			input: "Path: /usr/bin/*",
			want:  "Path: /usr/bin/*",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "only whitespace",
			input: "   \t\n   ",
			want:  "",
		},
		{
			name:  "unicode characters",
			input: "Hello \u00e9\u00e8\u00ea world",
			want:  "Hello \u00e9\u00e8\u00ea world",
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

func TestInferGoType_AdditionalCases(t *testing.T) {
	tests := []struct {
		name  string
		input ActionInput
		want  string
	}{
		{
			name:  "numeric string 123",
			input: ActionInput{Default: "123"},
			want:  "int",
		},
		{
			name:  "numeric string 0",
			input: ActionInput{Default: "0"},
			want:  "int",
		},
		{
			name:  "float string 1.5",
			input: ActionInput{Default: "1.5"},
			want:  "string", // floats treated as strings
		},
		{
			name:  "negative number",
			input: ActionInput{Default: "-5"},
			want:  "string", // negative numbers treated as strings
		},
		{
			name:  "mixed alphanumeric",
			input: ActionInput{Default: "v1.2.3"},
			want:  "string",
		},
		{
			name:  "empty default with number description",
			input: ActionInput{Description: "The retry count value"},
			want:  "int",
		},
		{
			name:  "empty default with timeout description",
			input: ActionInput{Description: "timeout value in ms"},
			want:  "int",
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

func TestGenerator_GenerateActionWrapper_MixedRequiredOptional(t *testing.T) {
	// Test with mix of required and optional fields
	spec := &ActionSpec{
		Name: "MixedFields",
		Inputs: map[string]ActionInput{
			"z-optional-2": {
				Description: "Last optional",
				Required:    false,
			},
			"a-optional-1": {
				Description: "First optional",
				Required:    false,
			},
			"z-required-2": {
				Description: "Last required",
				Required:    true,
			},
			"a-required-1": {
				Description: "First required",
				Required:    true,
			},
		},
	}

	gen := NewGenerator()
	code, err := gen.GenerateActionWrapper(ActionWrapperConfig{
		ActionRef:   "test/mixed@v1",
		PackageName: "mixed",
		TypeName:    "MixedFields",
		Spec:        spec,
	})
	if err != nil {
		t.Fatalf("GenerateActionWrapper() error = %v", err)
	}

	codeStr := string(code.Code)

	// Get positions of all fields
	aReq1 := strings.Index(codeStr, "ARequired1 string")
	zReq2 := strings.Index(codeStr, "ZRequired2 string")
	aOpt1 := strings.Index(codeStr, "AOptional1 string")
	zOpt2 := strings.Index(codeStr, "ZOptional2 string")

	// Required fields should come before optional
	if aReq1 > aOpt1 || aReq1 > zOpt2 {
		t.Error("Required field ARequired1 should come before optional fields")
	}
	if zReq2 > aOpt1 || zReq2 > zOpt2 {
		t.Error("Required field ZRequired2 should come before optional fields")
	}

	// Within required, alphabetically sorted
	if aReq1 > zReq2 {
		t.Error("Required fields should be alphabetically sorted")
	}

	// Within optional, alphabetically sorted
	if aOpt1 > zOpt2 {
		t.Error("Optional fields should be alphabetically sorted")
	}
}
