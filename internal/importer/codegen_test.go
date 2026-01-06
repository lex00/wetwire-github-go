package importer

import "testing"

func TestToVarName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Basic cases
		{"CI", "CI"},
		{"build", "Build"},
		{"my-workflow", "MyWorkflow"},
		{"my_workflow", "MyWorkflow"},
		{"my workflow", "MyWorkflow"},

		// Special characters that need sanitization
		{"C/C++ CI", "CCppCI"},
		{"C++ Build", "CppBuild"},
		{"Node.js CI", "NodeJsCI"},
		{"iOS Build", "IOSBuild"},
		{"D CI", "DCI"},

		// Parentheses
		{"Build (Linux)", "BuildLinux"},

		// Multiple special chars
		{"C/C++/Objective-C", "CCppObjectiveC"},

		// Reserved words
		{"type", "TypeJob"},
		{"go", "GoJob"},
		{"import", "ImportJob"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := toVarName(tt.input)
			if result != tt.expected {
				t.Errorf("toVarName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToVarNameProducesValidIdentifier(t *testing.T) {
	// These are real workflow names from actions/starter-workflows
	names := []string{
		"C/C++ CI",
		"D",
		"iOS",
		"Objective-C Xcode",
		".NET",
		"Node.js",
		"R",
	}

	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			result := toVarName(name)

			// Check it's not empty
			if result == "" {
				t.Errorf("toVarName(%q) returned empty string", name)
				return
			}

			// Check first char is letter (Go identifier requirement)
			first := rune(result[0])
			if !((first >= 'A' && first <= 'Z') || (first >= 'a' && first <= 'z')) {
				t.Errorf("toVarName(%q) = %q, first char is not a letter", name, result)
			}

			// Check all chars are valid Go identifier chars
			for i, r := range result {
				valid := (r >= 'A' && r <= 'Z') ||
				         (r >= 'a' && r <= 'z') ||
				         (r >= '0' && r <= '9' && i > 0) ||
				         r == '_'
				if !valid {
					t.Errorf("toVarName(%q) = %q, char %q at pos %d is invalid", name, result, string(r), i)
				}
			}
		})
	}
}
