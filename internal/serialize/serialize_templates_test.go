package serialize_test

import (
	"strings"
	"testing"

	"github.com/lex00/wetwire-github-go/codeowners"
	"github.com/lex00/wetwire-github-go/internal/serialize"
	"github.com/lex00/wetwire-github-go/templates"
)

// ===== Issue Template Tests =====

// TestBasicIssueTemplate tests basic issue template serialization.
func TestBasicIssueTemplate(t *testing.T) {
	tmpl := &templates.IssueTemplate{
		Name:        "Bug Report",
		Description: "File a bug report",
		Body: []templates.FormElement{
			templates.Markdown{
				Value: "## Bug Report\nPlease describe the bug.",
			},
			templates.Input{
				Label:    "Summary",
				Required: true,
			},
		},
	}

	yaml, err := serialize.IssueTemplateToYAML(tmpl)
	if err != nil {
		t.Fatalf("IssueTemplateToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	if !strings.Contains(yamlStr, "name: Bug Report") {
		t.Errorf("expected name, got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "description: File a bug report") {
		t.Errorf("expected description, got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "type: markdown") {
		t.Errorf("expected markdown type, got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "type: input") {
		t.Errorf("expected input type, got:\n%s", yamlStr)
	}
}

// TestIssueTemplateAllElements tests all form element types.
func TestIssueTemplateAllElements(t *testing.T) {
	tmpl := &templates.IssueTemplate{
		Name:        "Complete Form",
		Description: "Test all form elements",
		Title:       "New Issue",
		Labels:      []string{"bug", "needs-triage"},
		Assignees:   []string{"@maintainer"},
		Projects:    []string{"org/project/1"},
		Body: []templates.FormElement{
			templates.Markdown{
				ID:    "intro",
				Value: "## Welcome",
			},
			templates.Input{
				ID:          "summary",
				Label:       "Summary",
				Description: "Brief description",
				Placeholder: "Enter summary...",
				Value:       "Default value",
				Required:    true,
			},
			templates.Textarea{
				ID:          "details",
				Label:       "Details",
				Description: "Detailed description",
				Placeholder: "Enter details...",
				Value:       "Default details",
				Render:      "markdown",
				Required:    true,
			},
			templates.Dropdown{
				ID:          "severity",
				Label:       "Severity",
				Description: "Select severity",
				Options:     []string{"Low", "Medium", "High", "Critical"},
				Multiple:    false,
				Default:     2,
				Required:    true,
			},
			templates.Checkboxes{
				ID:          "checklist",
				Label:       "Pre-submission checklist",
				Description: "Please check all boxes",
				Options: []templates.CheckboxOption{
					{Label: "I have searched existing issues", Required: true},
					{Label: "I have read the documentation", Required: true},
					{Label: "I agree to the code of conduct", Required: false},
				},
			},
		},
	}

	yaml, err := serialize.IssueTemplateToYAML(tmpl)
	if err != nil {
		t.Fatalf("IssueTemplateToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	expectedFields := []string{
		"name: Complete Form",
		"description: Test all form elements",
		"title: New Issue",
		"labels:",
		"assignees:",
		"projects:",
		"type: markdown",
		"type: input",
		"type: textarea",
		"type: dropdown",
		"type: checkboxes",
		"render: markdown",
		"placeholder:",
		"options:",
		"default: 2",
		"validations:",
		"required: true",
	}

	for _, field := range expectedFields {
		if !strings.Contains(yamlStr, field) {
			t.Errorf("expected field %q, got:\n%s", field, yamlStr)
		}
	}
}

// TestIssueTemplateMinimal tests minimal issue template.
func TestIssueTemplateMinimal(t *testing.T) {
	tmpl := &templates.IssueTemplate{
		Name:        "Simple",
		Description: "Simple template",
		Body: []templates.FormElement{
			templates.Input{
				Label: "Title",
			},
		},
	}

	yaml, err := serialize.IssueTemplateToYAML(tmpl)
	if err != nil {
		t.Fatalf("IssueTemplateToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	// Should not have optional fields
	if strings.Contains(yamlStr, "title:") {
		t.Errorf("should not have title field, got:\n%s", yamlStr)
	}
	if strings.Contains(yamlStr, "labels:") {
		t.Errorf("should not have labels field, got:\n%s", yamlStr)
	}
}

// ===== Discussion Template Tests =====

// TestBasicDiscussionTemplate tests basic discussion template serialization.
func TestBasicDiscussionTemplate(t *testing.T) {
	tmpl := &templates.DiscussionTemplate{
		Title:       "General Discussion",
		Description: "Start a general discussion",
		Body: []templates.FormElement{
			templates.Markdown{
				Value: "## Discussion Topic",
			},
			templates.Textarea{
				Label:    "Description",
				Required: true,
			},
		},
	}

	yaml, err := serialize.DiscussionTemplateToYAML(tmpl)
	if err != nil {
		t.Fatalf("DiscussionTemplateToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	if !strings.Contains(yamlStr, "title: General Discussion") {
		t.Errorf("expected title, got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "description: Start a general discussion") {
		t.Errorf("expected description, got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "type: markdown") {
		t.Errorf("expected markdown type, got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "type: textarea") {
		t.Errorf("expected textarea type, got:\n%s", yamlStr)
	}
}

// TestDiscussionTemplateWithLabels tests discussion template with labels.
func TestDiscussionTemplateWithLabels(t *testing.T) {
	tmpl := &templates.DiscussionTemplate{
		Title:       "Feature Request",
		Description: "Request a new feature",
		Labels:      []string{"enhancement", "discussion"},
		Body: []templates.FormElement{
			templates.Input{
				Label:    "Feature Name",
				Required: true,
			},
		},
	}

	yaml, err := serialize.DiscussionTemplateToYAML(tmpl)
	if err != nil {
		t.Fatalf("DiscussionTemplateToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	if !strings.Contains(yamlStr, "labels:") {
		t.Errorf("expected labels field, got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "enhancement") {
		t.Errorf("expected enhancement label, got:\n%s", yamlStr)
	}
}

// ===== CODEOWNERS Tests =====

// TestBasicCodeowners tests basic CODEOWNERS serialization.
func TestBasicCodeowners(t *testing.T) {
	owners := &codeowners.Owners{
		Rules: []codeowners.Rule{
			{
				Pattern: "*",
				Owners:  []string{"@org/default-team"},
			},
			{
				Pattern: "*.go",
				Owners:  []string{"@go-team"},
			},
		},
	}

	text, err := serialize.CodeownersToText(owners)
	if err != nil {
		t.Fatalf("CodeownersToText failed: %v", err)
	}

	textStr := string(text)

	if !strings.Contains(textStr, "* @org/default-team") {
		t.Errorf("expected default rule, got:\n%s", textStr)
	}
	if !strings.Contains(textStr, "*.go @go-team") {
		t.Errorf("expected Go rule, got:\n%s", textStr)
	}
	if !strings.Contains(textStr, "# CODEOWNERS file generated by wetwire-github") {
		t.Errorf("expected header comment, got:\n%s", textStr)
	}
}

// TestCodeownersWithComments tests CODEOWNERS with comments.
func TestCodeownersWithComments(t *testing.T) {
	owners := &codeowners.Owners{
		Rules: []codeowners.Rule{
			{
				Pattern: "*.go",
				Owners:  []string{"@go-team"},
				Comment: "Go source files",
			},
			{
				Pattern: "*.js",
				Owners:  []string{"@js-team"},
				Comment: "JavaScript files",
			},
		},
	}

	text, err := serialize.CodeownersToText(owners)
	if err != nil {
		t.Fatalf("CodeownersToText failed: %v", err)
	}

	textStr := string(text)

	if !strings.Contains(textStr, "# Go source files") {
		t.Errorf("expected Go comment, got:\n%s", textStr)
	}
	if !strings.Contains(textStr, "# JavaScript files") {
		t.Errorf("expected JS comment, got:\n%s", textStr)
	}
}

// TestCodeownersMultipleOwners tests multiple owners per pattern.
func TestCodeownersMultipleOwners(t *testing.T) {
	owners := &codeowners.Owners{
		Rules: []codeowners.Rule{
			{
				Pattern: "/docs/",
				Owners:  []string{"@docs-team", "@user1", "@user2"},
			},
		},
	}

	text, err := serialize.CodeownersToText(owners)
	if err != nil {
		t.Fatalf("CodeownersToText failed: %v", err)
	}

	textStr := string(text)

	if !strings.Contains(textStr, "/docs/ @docs-team @user1 @user2") {
		t.Errorf("expected multiple owners, got:\n%s", textStr)
	}
}

// TestCodeownersEmptyOwners tests pattern without owners.
func TestCodeownersEmptyOwners(t *testing.T) {
	owners := &codeowners.Owners{
		Rules: []codeowners.Rule{
			{
				Pattern: "*.tmp",
				Owners:  []string{},
				Comment: "Temporary files",
			},
		},
	}

	text, err := serialize.CodeownersToText(owners)
	if err != nil {
		t.Fatalf("CodeownersToText failed: %v", err)
	}

	textStr := string(text)

	if !strings.Contains(textStr, "# Temporary files") {
		t.Errorf("expected comment, got:\n%s", textStr)
	}
	if !strings.Contains(textStr, "*.tmp") {
		t.Errorf("expected pattern, got:\n%s", textStr)
	}
}

// TestCodeownersRulesToText tests ExtractedCodeownersRule serialization.
func TestCodeownersRulesToText(t *testing.T) {
	rules := []serialize.ExtractedCodeownersRule{
		{
			Pattern: "*.go",
			Owners:  []string{"@go-team"},
			Comment: "Go files",
		},
		{
			Pattern: "/src/**",
			Owners:  []string{"@src-team"},
		},
	}

	text, err := serialize.CodeownersRulesToText(rules)
	if err != nil {
		t.Fatalf("CodeownersRulesToText failed: %v", err)
	}

	textStr := string(text)

	if !strings.Contains(textStr, "*.go @go-team") {
		t.Errorf("expected Go rule, got:\n%s", textStr)
	}
	if !strings.Contains(textStr, "/src/** @src-team") {
		t.Errorf("expected src rule, got:\n%s", textStr)
	}
	if !strings.Contains(textStr, "# Go files") {
		t.Errorf("expected comment, got:\n%s", textStr)
	}
}
