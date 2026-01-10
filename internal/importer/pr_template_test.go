package importer

import (
	"testing"
)

func TestParsePRTemplateContent(t *testing.T) {
	content := `# Pull Request

## Description
Please describe your changes.

## Checklist
- [ ] Tests pass
- [ ] Documentation updated
`
	result, err := ParsePRTemplateContent("default", content)
	if err != nil {
		t.Fatalf("ParsePRTemplateContent error: %v", err)
	}

	if result.Name != "default" {
		t.Errorf("Name = %q, want %q", result.Name, "default")
	}

	if result.Content != content {
		t.Errorf("Content mismatch")
	}
}

func TestGeneratePRTemplateCode(t *testing.T) {
	templates := map[string]*IRPRTemplate{
		"DefaultPR": {
			Name:    "default",
			Content: "# PR Template\n\nDescribe changes here.",
		},
	}

	gen := &PRTemplateCodeGenerator{PackageName: "workflows"}
	code, err := gen.Generate(templates)
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	if code.Templates != 1 {
		t.Errorf("Templates = %d, want 1", code.Templates)
	}

	templatesGo, ok := code.Files["templates.go"]
	if !ok {
		t.Fatal("expected templates.go in output")
	}

	// Check for expected content
	if !containsStr(templatesGo, "package workflows") {
		t.Error("expected package declaration")
	}

	if !containsStr(templatesGo, "var DefaultPR = templates.PRTemplate{") {
		t.Error("expected DefaultPR variable declaration")
	}

	if !containsStr(templatesGo, `Name: "default"`) {
		t.Error("expected Name field")
	}
}

func TestGeneratePRTemplateCode_MultipleTemplates(t *testing.T) {
	templates := map[string]*IRPRTemplate{
		"BugfixPR": {
			Name:    "bugfix",
			Content: "# Bugfix\n\nDescribe the bug fix.",
		},
		"FeaturePR": {
			Name:    "feature",
			Content: "# Feature\n\nDescribe the feature.",
		},
	}

	gen := &PRTemplateCodeGenerator{PackageName: "templates"}
	code, err := gen.Generate(templates)
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	if code.Templates != 2 {
		t.Errorf("Templates = %d, want 2", code.Templates)
	}

	templatesGo := code.Files["templates.go"]

	if !containsStr(templatesGo, "var BugfixPR") {
		t.Error("expected BugfixPR variable")
	}

	if !containsStr(templatesGo, "var FeaturePR") {
		t.Error("expected FeaturePR variable")
	}
}

func TestGeneratePRTemplateCode_EmptyName(t *testing.T) {
	templates := map[string]*IRPRTemplate{
		"DefaultPR": {
			Name:    "",
			Content: "# PR Template",
		},
	}

	gen := &PRTemplateCodeGenerator{PackageName: "workflows"}
	code, err := gen.Generate(templates)
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	templatesGo := code.Files["templates.go"]

	// Should not include Name field if empty
	if containsStr(templatesGo, `Name: ""`) {
		t.Error("should not include empty Name field")
	}
}

func TestGeneratePRTemplateCode_MultilineContent(t *testing.T) {
	content := `# Pull Request

## Description
Please include a summary of the changes.

## Type of change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Checklist
- [ ] My code follows the style guidelines
- [ ] I have performed a self-review
- [ ] I have added tests
`
	templates := map[string]*IRPRTemplate{
		"DetailedPR": {
			Name:    "detailed",
			Content: content,
		},
	}

	gen := &PRTemplateCodeGenerator{PackageName: "workflows"}
	code, err := gen.Generate(templates)
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	templatesGo := code.Files["templates.go"]

	// Should handle multiline content properly
	if !containsStr(templatesGo, "Content:") {
		t.Error("expected Content field")
	}
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestGeneratePRTemplateCode_ContentWithBackticks(t *testing.T) {
	content := "# PR Template\n\nUse `backticks` for code.\n\n```go\nfunc main() {}\n```"
	templates := map[string]*IRPRTemplate{
		"CodePR": {
			Name:    "code",
			Content: content,
		},
	}

	gen := &PRTemplateCodeGenerator{PackageName: "workflows"}
	code, err := gen.Generate(templates)
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	templatesGo := code.Files["templates.go"]

	// Content with backticks should be quoted
	if !containsStr(templatesGo, "Content:") {
		t.Error("expected Content field")
	}
	// Should use quoted string, not raw string literal
	if containsStr(templatesGo, "Content: `") {
		t.Error("should not use backtick literal when content has backticks")
	}
}

func TestFormatPRTemplateContent_WithBackticks(t *testing.T) {
	content := "Code: `example`"
	result := formatPRTemplateContent(content)

	// Should use quoted string since content has backticks
	if !containsStr(result, `"`) {
		t.Error("should use quoted string for content with backticks")
	}
}

func TestFormatPRTemplateContent_WithoutBackticks(t *testing.T) {
	content := "No backticks here"
	result := formatPRTemplateContent(content)

	// Should use raw string literal
	if !containsStr(result, "`") {
		t.Error("should use backtick literal for content without backticks")
	}
}

func TestGeneratePRTemplateCode_Empty(t *testing.T) {
	templates := map[string]*IRPRTemplate{}

	gen := &PRTemplateCodeGenerator{PackageName: "workflows"}
	code, err := gen.Generate(templates)
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	if code.Templates != 0 {
		t.Errorf("Templates = %d, want 0", code.Templates)
	}

	templatesGo := code.Files["templates.go"]
	if !containsStr(templatesGo, "package workflows") {
		t.Error("expected package declaration even with no templates")
	}
}
