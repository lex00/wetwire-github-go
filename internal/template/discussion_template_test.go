package template

import (
	"strings"
	"testing"

	"github.com/lex00/wetwire-github-go/internal/discover"
	"github.com/lex00/wetwire-github-go/internal/runner"
	"github.com/lex00/wetwire-github-go/templates"
)

func TestBuilder_BuildDiscussionTemplates_Empty(t *testing.T) {
	b := NewBuilder()

	discovered := &discover.DiscussionTemplateDiscoveryResult{
		Templates: []discover.DiscoveredDiscussionTemplate{},
	}
	extracted := &runner.DiscussionTemplateExtractionResult{
		Templates: []runner.ExtractedDiscussionTemplate{},
	}

	result, err := b.BuildDiscussionTemplates(discovered, extracted)
	if err != nil {
		t.Fatalf("BuildDiscussionTemplates() error = %v", err)
	}

	if len(result.Templates) != 0 {
		t.Errorf("Expected 0 templates, got %d", len(result.Templates))
	}
}

func TestBuilder_BuildDiscussionTemplates_SingleTemplate(t *testing.T) {
	b := NewBuilder()

	discovered := &discover.DiscussionTemplateDiscoveryResult{
		Templates: []discover.DiscoveredDiscussionTemplate{
			{Name: "IdeasTemplate", File: "templates.go", Line: 10},
		},
	}
	extracted := &runner.DiscussionTemplateExtractionResult{
		Templates: []runner.ExtractedDiscussionTemplate{
			{
				Name: "IdeasTemplate",
				Data: map[string]any{
					"Title":       "Share Your Ideas",
					"Description": "A place to share new ideas",
					"Labels":      []any{"idea", "discussion"},
					"Body": []any{
						map[string]any{
							"ID":    "intro",
							"Value": "## Welcome!\n\nShare your ideas here.",
						},
					},
				},
			},
		},
	}

	result, err := b.BuildDiscussionTemplates(discovered, extracted)
	if err != nil {
		t.Fatalf("BuildDiscussionTemplates() error = %v", err)
	}

	if len(result.Templates) != 1 {
		t.Fatalf("Expected 1 template, got %d", len(result.Templates))
	}

	tmpl := result.Templates[0]
	if tmpl.Name != "IdeasTemplate" {
		t.Errorf("Template name = %q, want %q", tmpl.Name, "IdeasTemplate")
	}

	if tmpl.Template.Title != "Share Your Ideas" {
		t.Errorf("Template.Title = %q, want %q", tmpl.Template.Title, "Share Your Ideas")
	}

	if len(tmpl.YAML) == 0 {
		t.Error("Template YAML is empty")
	}

	if len(tmpl.Template.Labels) != 2 {
		t.Errorf("Expected 2 labels, got %d", len(tmpl.Template.Labels))
	}
}

func TestBuilder_BuildDiscussionTemplates_MissingExtraction(t *testing.T) {
	b := NewBuilder()

	discovered := &discover.DiscussionTemplateDiscoveryResult{
		Templates: []discover.DiscoveredDiscussionTemplate{
			{Name: "MissingTemplate", File: "templates.go", Line: 10},
		},
	}
	// No extraction data
	extracted := &runner.DiscussionTemplateExtractionResult{
		Templates: []runner.ExtractedDiscussionTemplate{},
	}

	result, err := b.BuildDiscussionTemplates(discovered, extracted)
	if err != nil {
		t.Fatalf("BuildDiscussionTemplates() error = %v", err)
	}

	// Should have an error about missing extraction data
	if len(result.Errors) == 0 {
		t.Error("Expected error about missing extraction data")
	}

	if len(result.Templates) != 0 {
		t.Errorf("Expected 0 templates when extraction is missing, got %d", len(result.Templates))
	}
}

func TestBuilder_reconstructDiscussionTemplate(t *testing.T) {
	b := NewBuilder()

	tests := []struct {
		name string
		data map[string]any
		want func(*templates.DiscussionTemplate) bool
	}{
		{
			name: "basic template",
			data: map[string]any{
				"Title":       "Q&A",
				"Description": "Ask questions",
			},
			want: func(tmpl *templates.DiscussionTemplate) bool {
				return tmpl.Title == "Q&A" && tmpl.Description == "Ask questions"
			},
		},
		{
			name: "with labels as []any",
			data: map[string]any{
				"Title":  "Announcement",
				"Labels": []any{"announcement", "important"},
			},
			want: func(tmpl *templates.DiscussionTemplate) bool {
				return len(tmpl.Labels) == 2 && tmpl.Labels[0] == "announcement"
			},
		},
		{
			name: "with labels as []string",
			data: map[string]any{
				"Title":  "Feature Request",
				"Labels": []string{"feature", "request"},
			},
			want: func(tmpl *templates.DiscussionTemplate) bool {
				return len(tmpl.Labels) == 2 && tmpl.Labels[1] == "request"
			},
		},
		{
			name: "with body elements",
			data: map[string]any{
				"Title":       "Ideas",
				"Description": "Share ideas",
				"Body": []any{
					map[string]any{
						"ID":    "intro",
						"Value": "## Share your idea",
					},
				},
			},
			want: func(tmpl *templates.DiscussionTemplate) bool {
				return len(tmpl.Body) == 1
			},
		},
		{
			name: "with FormElement body",
			data: map[string]any{
				"Title": "Feedback",
				"Body": []templates.FormElement{
					templates.Markdown{ID: "intro", Value: "Intro"},
				},
			},
			want: func(tmpl *templates.DiscussionTemplate) bool {
				return len(tmpl.Body) == 1
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := b.reconstructDiscussionTemplate(tt.data)
			if result == nil {
				t.Fatal("reconstructDiscussionTemplate() returned nil")
			}
			if !tt.want(result) {
				t.Errorf("reconstructDiscussionTemplate() validation failed")
			}
		})
	}
}

func TestBuilder_BuildDiscussionTemplates_MultipleTemplates(t *testing.T) {
	b := NewBuilder()

	discovered := &discover.DiscussionTemplateDiscoveryResult{
		Templates: []discover.DiscoveredDiscussionTemplate{
			{Name: "Ideas", File: "templates.go", Line: 10},
			{Name: "QA", File: "templates.go", Line: 20},
			{Name: "Announcements", File: "templates.go", Line: 30},
		},
	}
	extracted := &runner.DiscussionTemplateExtractionResult{
		Templates: []runner.ExtractedDiscussionTemplate{
			{
				Name: "Ideas",
				Data: map[string]any{
					"Title":       "Ideas",
					"Description": "Share ideas",
					"Labels":      []any{"idea"},
				},
			},
			{
				Name: "QA",
				Data: map[string]any{
					"Title":       "Q&A",
					"Description": "Ask questions",
					"Labels":      []any{"question"},
				},
			},
			{
				Name: "Announcements",
				Data: map[string]any{
					"Title":       "Announcements",
					"Description": "Important announcements",
					"Labels":      []any{"announcement"},
				},
			},
		},
	}

	result, err := b.BuildDiscussionTemplates(discovered, extracted)
	if err != nil {
		t.Fatalf("BuildDiscussionTemplates() error = %v", err)
	}

	if len(result.Templates) != 3 {
		t.Fatalf("Expected 3 templates, got %d", len(result.Templates))
	}
}

func TestBuilder_BuildDiscussionTemplates_ComplexTemplate(t *testing.T) {
	b := NewBuilder()

	discovered := &discover.DiscussionTemplateDiscoveryResult{
		Templates: []discover.DiscoveredDiscussionTemplate{
			{Name: "ComplexIdeas", File: "templates.go", Line: 10},
		},
	}
	extracted := &runner.DiscussionTemplateExtractionResult{
		Templates: []runner.ExtractedDiscussionTemplate{
			{
				Name: "ComplexIdeas",
				Data: map[string]any{
					"Title":       "Feature Ideas",
					"Description": "Propose new features",
					"Labels":      []any{"idea", "feature", "enhancement"},
					"Body": []any{
						map[string]any{
							"ID":    "welcome",
							"Value": "## Welcome!\n\nThank you for your idea.",
						},
						map[string]any{
							"ID":          "title",
							"Label":       "Idea Title",
							"Description": "A short, catchy title",
							"Required":    true,
						},
						map[string]any{
							"ID":          "description",
							"Label":       "Description",
							"Description": "Describe your idea in detail",
							"Render":      "markdown",
							"Required":    true,
						},
						map[string]any{
							"ID":          "category",
							"Label":       "Category",
							"Description": "Which area does this affect?",
							"Options":     []any{"UI/UX", "Backend", "API", "Documentation", "Other"},
							"Multiple":    true,
						},
						map[string]any{
							"ID":    "agreement",
							"Label": "Agreement",
							"Options": []any{
								map[string]any{
									"Label":    "I have searched for similar ideas",
									"Required": true,
								},
							},
						},
					},
				},
			},
		},
	}

	result, err := b.BuildDiscussionTemplates(discovered, extracted)
	if err != nil {
		t.Fatalf("BuildDiscussionTemplates() error = %v", err)
	}

	if len(result.Templates) != 1 {
		t.Fatalf("Expected 1 template, got %d", len(result.Templates))
	}

	tmpl := result.Templates[0]

	// Verify complex structure
	if len(tmpl.Template.Labels) != 3 {
		t.Errorf("Expected 3 labels, got %d", len(tmpl.Template.Labels))
	}
	if len(tmpl.Template.Body) != 5 {
		t.Errorf("Expected 5 body elements, got %d", len(tmpl.Template.Body))
	}

	// Check YAML output
	yaml := string(tmpl.YAML)
	if !strings.Contains(yaml, "title:") {
		t.Error("YAML should contain title field")
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

func TestBuilder_BuildDiscussionTemplates_YAMLOutput(t *testing.T) {
	b := NewBuilder()

	discovered := &discover.DiscussionTemplateDiscoveryResult{
		Templates: []discover.DiscoveredDiscussionTemplate{
			{Name: "TestTemplate", File: "templates.go", Line: 10},
		},
	}
	extracted := &runner.DiscussionTemplateExtractionResult{
		Templates: []runner.ExtractedDiscussionTemplate{
			{
				Name: "TestTemplate",
				Data: map[string]any{
					"Title":       "Test Discussion",
					"Description": "Test template",
					"Labels":      []any{"test"},
					"Body": []any{
						map[string]any{
							"ID":    "intro",
							"Value": "Test content",
						},
					},
				},
			},
		},
	}

	result, err := b.BuildDiscussionTemplates(discovered, extracted)
	if err != nil {
		t.Fatalf("BuildDiscussionTemplates() error = %v", err)
	}

	if len(result.Templates) != 1 {
		t.Fatalf("Expected 1 template, got %d", len(result.Templates))
	}

	yaml := string(result.Templates[0].YAML)

	// Verify YAML contains expected fields
	expectedFields := []string{"title:", "description:", "labels:", "body:"}
	for _, field := range expectedFields {
		if !strings.Contains(yaml, field) {
			t.Errorf("YAML missing expected field %q", field)
		}
	}
}
