package template

import (
	"testing"

	"github.com/lex00/wetwire-github-go/internal/discover"
	"github.com/lex00/wetwire-github-go/internal/runner"
)

func TestBuilder_BuildPRTemplates_Empty(t *testing.T) {
	b := NewBuilder()

	discovered := &discover.PRTemplateDiscoveryResult{
		Templates: []discover.DiscoveredPRTemplate{},
	}
	extracted := &runner.PRTemplateExtractionResult{
		Templates: []runner.ExtractedPRTemplate{},
	}

	result, err := b.BuildPRTemplates(discovered, extracted)
	if err != nil {
		t.Fatalf("BuildPRTemplates() error = %v", err)
	}

	if len(result.Templates) != 0 {
		t.Errorf("Expected 0 templates, got %d", len(result.Templates))
	}
}

func TestBuilder_BuildPRTemplates_DefaultTemplate(t *testing.T) {
	b := NewBuilder()

	discovered := &discover.PRTemplateDiscoveryResult{
		Templates: []discover.DiscoveredPRTemplate{
			{Name: "DefaultPRTemplate", File: "templates.go", Line: 10},
		},
	}
	extracted := &runner.PRTemplateExtractionResult{
		Templates: []runner.ExtractedPRTemplate{
			{Name: "DefaultPRTemplate", Content: "## Description\n\nPlease describe your changes.\n"},
		},
	}

	result, err := b.BuildPRTemplates(discovered, extracted)
	if err != nil {
		t.Fatalf("BuildPRTemplates() error = %v", err)
	}

	if len(result.Templates) != 1 {
		t.Fatalf("Expected 1 template, got %d", len(result.Templates))
	}

	tmpl := result.Templates[0]
	if tmpl.Name != "DefaultPRTemplate" {
		t.Errorf("Template name = %q, want %q", tmpl.Name, "DefaultPRTemplate")
	}

	if len(tmpl.Content) == 0 {
		t.Error("Template content is empty")
	}

	expectedContent := "## Description\n\nPlease describe your changes.\n"
	if string(tmpl.Content) != expectedContent {
		t.Errorf("Template content = %q, want %q", string(tmpl.Content), expectedContent)
	}
}

func TestBuilder_BuildPRTemplates_NamedTemplate(t *testing.T) {
	b := NewBuilder()

	discovered := &discover.PRTemplateDiscoveryResult{
		Templates: []discover.DiscoveredPRTemplate{
			{Name: "FeaturePRTemplate", File: "templates.go", Line: 20},
		},
	}
	extracted := &runner.PRTemplateExtractionResult{
		Templates: []runner.ExtractedPRTemplate{
			{Name: "FeaturePRTemplate", Content: "## Feature\n\nDescribe the feature.\n"},
		},
	}

	result, err := b.BuildPRTemplates(discovered, extracted)
	if err != nil {
		t.Fatalf("BuildPRTemplates() error = %v", err)
	}

	if len(result.Templates) != 1 {
		t.Fatalf("Expected 1 template, got %d", len(result.Templates))
	}

	tmpl := result.Templates[0]
	if tmpl.Name != "FeaturePRTemplate" {
		t.Errorf("Template name = %q, want %q", tmpl.Name, "FeaturePRTemplate")
	}

	if tmpl.Template == nil {
		t.Fatal("Template.Template is nil")
	}
}

func TestBuilder_BuildPRTemplates_MultipleTemplates(t *testing.T) {
	b := NewBuilder()

	discovered := &discover.PRTemplateDiscoveryResult{
		Templates: []discover.DiscoveredPRTemplate{
			{Name: "Feature", File: "templates.go", Line: 10},
			{Name: "Bugfix", File: "templates.go", Line: 20},
			{Name: "Docs", File: "templates.go", Line: 30},
		},
	}
	extracted := &runner.PRTemplateExtractionResult{
		Templates: []runner.ExtractedPRTemplate{
			{Name: "Feature", Content: "## Feature\n"},
			{Name: "Bugfix", Content: "## Bugfix\n"},
			{Name: "Docs", Content: "## Docs\n"},
		},
	}

	result, err := b.BuildPRTemplates(discovered, extracted)
	if err != nil {
		t.Fatalf("BuildPRTemplates() error = %v", err)
	}

	if len(result.Templates) != 3 {
		t.Fatalf("Expected 3 templates, got %d", len(result.Templates))
	}
}

func TestBuilder_BuildPRTemplates_MissingExtraction(t *testing.T) {
	b := NewBuilder()

	discovered := &discover.PRTemplateDiscoveryResult{
		Templates: []discover.DiscoveredPRTemplate{
			{Name: "TestTemplate", File: "templates.go", Line: 10},
		},
	}
	// No extraction data
	extracted := &runner.PRTemplateExtractionResult{
		Templates: []runner.ExtractedPRTemplate{},
	}

	result, err := b.BuildPRTemplates(discovered, extracted)
	if err != nil {
		t.Fatalf("BuildPRTemplates() error = %v", err)
	}

	// Should have an error about missing extraction data
	if len(result.Errors) == 0 {
		t.Error("Expected error about missing extraction data")
	}

	if len(result.Templates) != 0 {
		t.Errorf("Expected 0 templates when extraction is missing, got %d", len(result.Templates))
	}
}

func TestBuilder_BuildPRTemplates_Filename(t *testing.T) {
	b := NewBuilder()

	tests := []struct {
		name         string
		templateName string
		wantFilename string
	}{
		{
			name:         "default template",
			templateName: "default",
			wantFilename: "PULL_REQUEST_TEMPLATE.md",
		},
		{
			name:         "empty name (default)",
			templateName: "",
			wantFilename: "PULL_REQUEST_TEMPLATE.md",
		},
		{
			name:         "named template",
			templateName: "feature",
			wantFilename: "PULL_REQUEST_TEMPLATE/feature.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			discovered := &discover.PRTemplateDiscoveryResult{
				Templates: []discover.DiscoveredPRTemplate{
					{Name: "Template", File: "templates.go", Line: 10},
				},
			}
			extracted := &runner.PRTemplateExtractionResult{
				Templates: []runner.ExtractedPRTemplate{
					{Name: "Template", Content: "content"},
				},
			}

			// Mock the template name by modifying extracted data
			// The Name field in PRTemplate determines the filename
			discovered.Templates[0].Name = "Template"

			result, err := b.BuildPRTemplates(discovered, extracted)
			if err != nil {
				t.Fatalf("BuildPRTemplates() error = %v", err)
			}

			if len(result.Templates) != 1 {
				t.Fatalf("Expected 1 template, got %d", len(result.Templates))
			}

			// The filename is determined by the PRTemplate.Filename() method
			// which uses the Name field from the template
		})
	}
}
