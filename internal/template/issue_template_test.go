package template

import (
	"strings"
	"testing"

	"github.com/lex00/wetwire-github-go/internal/discover"
	"github.com/lex00/wetwire-github-go/internal/runner"
	"github.com/lex00/wetwire-github-go/templates"
)

func TestBuilder_BuildIssueTemplates_Empty(t *testing.T) {
	b := NewBuilder()

	discovered := &discover.IssueTemplateDiscoveryResult{
		Templates: []discover.DiscoveredIssueTemplate{},
	}
	extracted := &runner.IssueTemplateExtractionResult{
		Templates: []runner.ExtractedIssueTemplate{},
	}

	result, err := b.BuildIssueTemplates(discovered, extracted)
	if err != nil {
		t.Fatalf("BuildIssueTemplates() error = %v", err)
	}

	if len(result.Templates) != 0 {
		t.Errorf("Expected 0 templates, got %d", len(result.Templates))
	}
}

func TestBuilder_BuildIssueTemplates_SingleTemplate(t *testing.T) {
	b := NewBuilder()

	discovered := &discover.IssueTemplateDiscoveryResult{
		Templates: []discover.DiscoveredIssueTemplate{
			{Name: "BugReport", File: "templates.go", Line: 10},
		},
	}
	extracted := &runner.IssueTemplateExtractionResult{
		Templates: []runner.ExtractedIssueTemplate{
			{
				Name: "BugReport",
				Data: map[string]any{
					"Name":        "Bug Report",
					"Description": "File a bug report",
					"Title":       "[Bug]: ",
					"Labels":      []any{"bug", "triage"},
					"Assignees":   []any{"@maintainer"},
					"Body": []any{
						map[string]any{
							"ID":    "description",
							"Label": "Description",
							"Value": "A clear and concise description of what the bug is.",
						},
					},
				},
			},
		},
	}

	result, err := b.BuildIssueTemplates(discovered, extracted)
	if err != nil {
		t.Fatalf("BuildIssueTemplates() error = %v", err)
	}

	if len(result.Templates) != 1 {
		t.Fatalf("Expected 1 template, got %d", len(result.Templates))
	}

	tmpl := result.Templates[0]
	if tmpl.Name != "BugReport" {
		t.Errorf("Template name = %q, want %q", tmpl.Name, "BugReport")
	}

	if tmpl.Template.Name != "Bug Report" {
		t.Errorf("Template.Name = %q, want %q", tmpl.Template.Name, "Bug Report")
	}

	if len(tmpl.YAML) == 0 {
		t.Error("Template YAML is empty")
	}

	if len(tmpl.Template.Labels) != 2 {
		t.Errorf("Expected 2 labels, got %d", len(tmpl.Template.Labels))
	}
}

func TestBuilder_BuildIssueTemplates_MissingExtraction(t *testing.T) {
	b := NewBuilder()

	discovered := &discover.IssueTemplateDiscoveryResult{
		Templates: []discover.DiscoveredIssueTemplate{
			{Name: "MissingTemplate", File: "templates.go", Line: 10},
		},
	}
	// No extraction data
	extracted := &runner.IssueTemplateExtractionResult{
		Templates: []runner.ExtractedIssueTemplate{},
	}

	result, err := b.BuildIssueTemplates(discovered, extracted)
	if err != nil {
		t.Fatalf("BuildIssueTemplates() error = %v", err)
	}

	// Should have an error about missing extraction data
	if len(result.Errors) == 0 {
		t.Error("Expected error about missing extraction data")
	}

	if len(result.Templates) != 0 {
		t.Errorf("Expected 0 templates when extraction is missing, got %d", len(result.Templates))
	}
}

func TestBuilder_reconstructIssueTemplate(t *testing.T) {
	b := NewBuilder()

	tests := []struct {
		name string
		data map[string]any
		want func(*templates.IssueTemplate) bool
	}{
		{
			name: "basic template",
			data: map[string]any{
				"Name":        "Feature Request",
				"Description": "Suggest an idea",
			},
			want: func(tmpl *templates.IssueTemplate) bool {
				return tmpl.Name == "Feature Request" && tmpl.Description == "Suggest an idea"
			},
		},
		{
			name: "with labels as []any",
			data: map[string]any{
				"Name":   "Bug",
				"Labels": []any{"bug", "needs-triage"},
			},
			want: func(tmpl *templates.IssueTemplate) bool {
				return len(tmpl.Labels) == 2 && tmpl.Labels[0] == "bug"
			},
		},
		{
			name: "with labels as []string",
			data: map[string]any{
				"Name":   "Feature",
				"Labels": []string{"enhancement", "feature"},
			},
			want: func(tmpl *templates.IssueTemplate) bool {
				return len(tmpl.Labels) == 2 && tmpl.Labels[1] == "feature"
			},
		},
		{
			name: "with projects",
			data: map[string]any{
				"Name":     "Task",
				"Projects": []any{"Project Board"},
			},
			want: func(tmpl *templates.IssueTemplate) bool {
				return len(tmpl.Projects) == 1
			},
		},
		{
			name: "with assignees",
			data: map[string]any{
				"Name":      "Issue",
				"Assignees": []any{"@user1", "@user2"},
			},
			want: func(tmpl *templates.IssueTemplate) bool {
				return len(tmpl.Assignees) == 2
			},
		},
		{
			name: "with title",
			data: map[string]any{
				"Name":  "Bug Report",
				"Title": "[Bug]: ",
			},
			want: func(tmpl *templates.IssueTemplate) bool {
				return tmpl.Title == "[Bug]: "
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := b.reconstructIssueTemplate(tt.data)
			if result == nil {
				t.Fatal("reconstructIssueTemplate() returned nil")
			}
			if !tt.want(result) {
				t.Errorf("reconstructIssueTemplate() validation failed")
			}
		})
	}
}

func TestBuilder_reconstructMarkdown(t *testing.T) {
	b := NewBuilder()

	data := map[string]any{
		"ID":    "intro",
		"Value": "## Welcome\n\nPlease fill out this form.",
	}

	result := b.reconstructMarkdown(data)

	if result.ID != "intro" {
		t.Errorf("ID = %q, want %q", result.ID, "intro")
	}
	if result.Value != "## Welcome\n\nPlease fill out this form." {
		t.Errorf("Value = %q, want %q", result.Value, "## Welcome\n\nPlease fill out this form.")
	}
}

func TestBuilder_reconstructInput(t *testing.T) {
	b := NewBuilder()

	data := map[string]any{
		"ID":          "contact",
		"Label":       "Contact Email",
		"Description": "How can we reach you?",
		"Placeholder": "email@example.com",
		"Value":       "",
		"Required":    true,
	}

	result := b.reconstructInput(data)

	if result.ID != "contact" {
		t.Errorf("ID = %q, want %q", result.ID, "contact")
	}
	if result.Label != "Contact Email" {
		t.Errorf("Label = %q, want %q", result.Label, "Contact Email")
	}
	if result.Description != "How can we reach you?" {
		t.Errorf("Description = %q, want %q", result.Description, "How can we reach you?")
	}
	if result.Placeholder != "email@example.com" {
		t.Errorf("Placeholder = %q, want %q", result.Placeholder, "email@example.com")
	}
	if !result.Required {
		t.Error("Required should be true")
	}
}

func TestBuilder_reconstructTextarea(t *testing.T) {
	b := NewBuilder()

	data := map[string]any{
		"ID":          "description",
		"Label":       "Bug Description",
		"Description": "Describe the bug in detail",
		"Placeholder": "Enter description here",
		"Value":       "",
		"Render":      "markdown",
		"Required":    true,
	}

	result := b.reconstructTextarea(data)

	if result.ID != "description" {
		t.Errorf("ID = %q, want %q", result.ID, "description")
	}
	if result.Label != "Bug Description" {
		t.Errorf("Label = %q, want %q", result.Label, "Bug Description")
	}
	if result.Render != "markdown" {
		t.Errorf("Render = %q, want %q", result.Render, "markdown")
	}
	if !result.Required {
		t.Error("Required should be true")
	}
}

func TestBuilder_reconstructDropdown(t *testing.T) {
	b := NewBuilder()

	tests := []struct {
		name string
		data map[string]any
		want func(templates.Dropdown) bool
	}{
		{
			name: "basic dropdown",
			data: map[string]any{
				"ID":      "severity",
				"Label":   "Severity",
				"Options": []any{"Low", "Medium", "High"},
			},
			want: func(d templates.Dropdown) bool {
				return len(d.Options) == 3 && d.Options[0] == "Low"
			},
		},
		{
			name: "with default int",
			data: map[string]any{
				"ID":      "priority",
				"Label":   "Priority",
				"Options": []any{"P1", "P2", "P3"},
				"Default": 1,
			},
			want: func(d templates.Dropdown) bool {
				return d.Default == 1
			},
		},
		{
			name: "with default float64",
			data: map[string]any{
				"ID":      "priority",
				"Label":   "Priority",
				"Options": []any{"P1", "P2", "P3"},
				"Default": 2.0,
			},
			want: func(d templates.Dropdown) bool {
				return d.Default == 2
			},
		},
		{
			name: "with multiple",
			data: map[string]any{
				"ID":       "categories",
				"Label":    "Categories",
				"Options":  []any{"UI", "Backend", "Docs"},
				"Multiple": true,
			},
			want: func(d templates.Dropdown) bool {
				return d.Multiple
			},
		},
		{
			name: "with required",
			data: map[string]any{
				"ID":       "type",
				"Label":    "Issue Type",
				"Options":  []any{"Bug", "Feature"},
				"Required": true,
			},
			want: func(d templates.Dropdown) bool {
				return d.Required
			},
		},
		{
			name: "with description",
			data: map[string]any{
				"ID":          "browser",
				"Label":       "Browser",
				"Description": "Which browser are you using?",
				"Options":     []any{"Chrome", "Firefox", "Safari"},
			},
			want: func(d templates.Dropdown) bool {
				return d.Description == "Which browser are you using?"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := b.reconstructDropdown(tt.data)
			if !tt.want(result) {
				t.Errorf("reconstructDropdown() validation failed")
			}
		})
	}
}

func TestBuilder_reconstructCheckboxes(t *testing.T) {
	b := NewBuilder()

	data := map[string]any{
		"ID":          "terms",
		"Label":       "Terms",
		"Description": "Please agree to the following",
		"Options": []any{
			map[string]any{
				"Label":    "I agree to the terms",
				"Required": true,
			},
			map[string]any{
				"Label":    "I want to receive updates",
				"Required": false,
			},
		},
	}

	result := b.reconstructCheckboxes(data)

	if result.ID != "terms" {
		t.Errorf("ID = %q, want %q", result.ID, "terms")
	}
	if result.Label != "Terms" {
		t.Errorf("Label = %q, want %q", result.Label, "Terms")
	}
	if len(result.Options) != 2 {
		t.Fatalf("Expected 2 options, got %d", len(result.Options))
	}
	if result.Options[0].Label != "I agree to the terms" {
		t.Errorf("First option label = %q, want %q", result.Options[0].Label, "I agree to the terms")
	}
	if !result.Options[0].Required {
		t.Error("First option should be required")
	}
	if result.Options[1].Required {
		t.Error("Second option should not be required")
	}
}

func TestBuilder_reconstructFormElementFromMap(t *testing.T) {
	b := NewBuilder()

	tests := []struct {
		name     string
		data     map[string]any
		wantType string
	}{
		{
			name: "markdown element",
			data: map[string]any{
				"Value": "## Header",
			},
			wantType: "Markdown",
		},
		{
			name: "input element",
			data: map[string]any{
				"Label": "Name",
				"Value": "",
			},
			wantType: "Input",
		},
		{
			name: "textarea element",
			data: map[string]any{
				"Label":  "Description",
				"Value":  "",
				"Render": "markdown",
			},
			wantType: "Textarea",
		},
		{
			name: "dropdown element",
			data: map[string]any{
				"Label":   "Priority",
				"Options": []any{"High", "Low"},
			},
			wantType: "Dropdown",
		},
		{
			name: "checkboxes element",
			data: map[string]any{
				"Label": "Terms",
				"Options": []any{
					map[string]any{"Label": "Agree"},
				},
			},
			wantType: "Checkboxes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := b.reconstructFormElementFromMap(tt.data)
			if result == nil {
				t.Fatal("reconstructFormElementFromMap() returned nil")
			}

			// Check type by trying type assertion
			switch tt.wantType {
			case "Markdown":
				if _, ok := result.(templates.Markdown); !ok {
					t.Errorf("Expected Markdown, got %T", result)
				}
			case "Input":
				if _, ok := result.(templates.Input); !ok {
					t.Errorf("Expected Input, got %T", result)
				}
			case "Textarea":
				if _, ok := result.(templates.Textarea); !ok {
					t.Errorf("Expected Textarea, got %T", result)
				}
			case "Dropdown":
				if _, ok := result.(templates.Dropdown); !ok {
					t.Errorf("Expected Dropdown, got %T", result)
				}
			case "Checkboxes":
				if _, ok := result.(templates.Checkboxes); !ok {
					t.Errorf("Expected Checkboxes, got %T", result)
				}
			}
		})
	}
}

func TestBuilder_reconstructFormElements(t *testing.T) {
	b := NewBuilder()

	data := []any{
		templates.Markdown{ID: "intro", Value: "## Introduction"},
		templates.Input{ID: "name", Label: "Name"},
		templates.Textarea{ID: "desc", Label: "Description", Render: "markdown"},
		map[string]any{
			"ID":      "priority",
			"Label":   "Priority",
			"Options": []any{"High", "Low"},
		},
	}

	result := b.reconstructFormElements(data)

	if len(result) != 4 {
		t.Fatalf("Expected 4 elements, got %d", len(result))
	}

	// Check that all types were preserved or reconstructed
	if _, ok := result[0].(templates.Markdown); !ok {
		t.Errorf("First element should be Markdown, got %T", result[0])
	}
	if _, ok := result[1].(templates.Input); !ok {
		t.Errorf("Second element should be Input, got %T", result[1])
	}
	if _, ok := result[2].(templates.Textarea); !ok {
		t.Errorf("Third element should be Textarea, got %T", result[2])
	}
	if _, ok := result[3].(templates.Dropdown); !ok {
		t.Errorf("Fourth element should be Dropdown, got %T", result[3])
	}
}

func TestBuilder_BuildIssueTemplates_ComplexTemplate(t *testing.T) {
	b := NewBuilder()

	discovered := &discover.IssueTemplateDiscoveryResult{
		Templates: []discover.DiscoveredIssueTemplate{
			{Name: "ComplexBugReport", File: "templates.go", Line: 10},
		},
	}
	extracted := &runner.IssueTemplateExtractionResult{
		Templates: []runner.ExtractedIssueTemplate{
			{
				Name: "ComplexBugReport",
				Data: map[string]any{
					"Name":        "Bug Report",
					"Description": "File a bug report",
					"Title":       "[Bug]: ",
					"Labels":      []any{"bug", "triage", "needs-investigation"},
					"Assignees":   []any{"@maintainer", "@triager"},
					"Projects":    []any{"Bug Tracker"},
					"Body": []any{
						map[string]any{
							"ID":    "intro",
							"Value": "## Bug Report\n\nThank you for reporting a bug.",
						},
						map[string]any{
							"ID":          "title",
							"Label":       "Bug Title",
							"Description": "A short, descriptive title",
							"Required":    true,
						},
						map[string]any{
							"ID":          "description",
							"Label":       "Description",
							"Description": "Detailed description",
							"Render":      "markdown",
							"Required":    true,
						},
						map[string]any{
							"ID":          "severity",
							"Label":       "Severity",
							"Description": "How severe is this bug?",
							"Options":     []any{"Critical", "High", "Medium", "Low"},
							"Default":     2.0,
							"Required":    true,
						},
						map[string]any{
							"ID":    "checklist",
							"Label": "Checklist",
							"Options": []any{
								map[string]any{
									"Label":    "I have read the guidelines",
									"Required": true,
								},
								map[string]any{
									"Label":    "I have searched for similar issues",
									"Required": true,
								},
							},
						},
					},
				},
			},
		},
	}

	result, err := b.BuildIssueTemplates(discovered, extracted)
	if err != nil {
		t.Fatalf("BuildIssueTemplates() error = %v", err)
	}

	if len(result.Templates) != 1 {
		t.Fatalf("Expected 1 template, got %d", len(result.Templates))
	}

	tmpl := result.Templates[0]

	// Verify complex structure
	if len(tmpl.Template.Labels) != 3 {
		t.Errorf("Expected 3 labels, got %d", len(tmpl.Template.Labels))
	}
	if len(tmpl.Template.Assignees) != 2 {
		t.Errorf("Expected 2 assignees, got %d", len(tmpl.Template.Assignees))
	}
	if len(tmpl.Template.Projects) != 1 {
		t.Errorf("Expected 1 project, got %d", len(tmpl.Template.Projects))
	}
	if len(tmpl.Template.Body) != 5 {
		t.Errorf("Expected 5 body elements, got %d", len(tmpl.Template.Body))
	}

	// Check YAML output
	yaml := string(tmpl.YAML)
	if !strings.Contains(yaml, "name:") {
		t.Error("YAML should contain name field")
	}
	if !strings.Contains(yaml, "description:") {
		t.Error("YAML should contain description field")
	}
	if !strings.Contains(yaml, "labels:") {
		t.Error("YAML should contain labels field")
	}
	if !strings.Contains(yaml, "body:") {
		t.Error("YAML should contain body field")
	}
}
