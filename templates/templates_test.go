package templates

import (
	"testing"
)

func TestIssueTemplate_ResourceType(t *testing.T) {
	it := IssueTemplate{}
	if got := it.ResourceType(); got != "issue-template" {
		t.Errorf("ResourceType() = %q, want %q", got, "issue-template")
	}
}

func TestIssueTemplate_Basic(t *testing.T) {
	it := IssueTemplate{
		Name:        "Bug Report",
		Description: "Report a bug in the project",
		Title:       "[Bug]: ",
		Labels:      []string{"bug", "triage"},
		Assignees:   []string{"maintainer1"},
		Body: []FormElement{
			Input{
				ID:       "title",
				Label:    "Bug Title",
				Required: true,
			},
		},
	}

	if it.Name != "Bug Report" {
		t.Errorf("Name = %q, want %q", it.Name, "Bug Report")
	}

	if it.Description != "Report a bug in the project" {
		t.Errorf("Description = %q, want %q", it.Description, "Report a bug in the project")
	}

	if it.Title != "[Bug]: " {
		t.Errorf("Title = %q, want %q", it.Title, "[Bug]: ")
	}

	if len(it.Labels) != 2 {
		t.Errorf("len(Labels) = %d, want 2", len(it.Labels))
	}

	if len(it.Assignees) != 1 {
		t.Errorf("len(Assignees) = %d, want 1", len(it.Assignees))
	}

	if len(it.Body) != 1 {
		t.Errorf("len(Body) = %d, want 1", len(it.Body))
	}
}

func TestIssueTemplate_WithProjects(t *testing.T) {
	it := IssueTemplate{
		Name:        "Feature Request",
		Description: "Request a new feature",
		Projects:    []string{"org/1", "org/2"},
		Body: []FormElement{
			Textarea{
				ID:    "description",
				Label: "Feature Description",
			},
		},
	}

	if len(it.Projects) != 2 {
		t.Errorf("len(Projects) = %d, want 2", len(it.Projects))
	}

	if it.Projects[0] != "org/1" {
		t.Errorf("Projects[0] = %q, want %q", it.Projects[0], "org/1")
	}
}

func TestDiscussionTemplate_ResourceType(t *testing.T) {
	dt := DiscussionTemplate{}
	if got := dt.ResourceType(); got != "discussion-template" {
		t.Errorf("ResourceType() = %q, want %q", got, "discussion-template")
	}
}

func TestDiscussionTemplate_Basic(t *testing.T) {
	dt := DiscussionTemplate{
		Title:       "General Discussion",
		Description: "Start a general discussion",
		Labels:      []string{"discussion"},
		Body: []FormElement{
			Textarea{
				ID:    "content",
				Label: "Discussion Content",
			},
		},
	}

	if dt.Title != "General Discussion" {
		t.Errorf("Title = %q, want %q", dt.Title, "General Discussion")
	}

	if dt.Description != "Start a general discussion" {
		t.Errorf("Description = %q, want %q", dt.Description, "Start a general discussion")
	}

	if len(dt.Labels) != 1 {
		t.Errorf("len(Labels) = %d, want 1", len(dt.Labels))
	}

	if len(dt.Body) != 1 {
		t.Errorf("len(Body) = %d, want 1", len(dt.Body))
	}
}

func TestMarkdown_ElementType(t *testing.T) {
	m := Markdown{}
	if got := m.ElementType(); got != "markdown" {
		t.Errorf("ElementType() = %q, want %q", got, "markdown")
	}
}

func TestMarkdown_Fields(t *testing.T) {
	m := Markdown{
		ID:    "welcome",
		Value: "## Welcome\n\nThank you for reporting!",
	}

	if m.ID != "welcome" {
		t.Errorf("ID = %q, want %q", m.ID, "welcome")
	}

	if m.Value != "## Welcome\n\nThank you for reporting!" {
		t.Errorf("Value = %q, want %q", m.Value, "## Welcome\n\nThank you for reporting!")
	}
}

func TestInput_ElementType(t *testing.T) {
	i := Input{}
	if got := i.ElementType(); got != "input" {
		t.Errorf("ElementType() = %q, want %q", got, "input")
	}
}

func TestInput_AllFields(t *testing.T) {
	i := Input{
		ID:          "email",
		Label:       "Email Address",
		Description: "Your email for follow-up",
		Placeholder: "user@example.com",
		Value:       "default@example.com",
		Required:    true,
	}

	if i.ID != "email" {
		t.Errorf("ID = %q, want %q", i.ID, "email")
	}

	if i.Label != "Email Address" {
		t.Errorf("Label = %q, want %q", i.Label, "Email Address")
	}

	if i.Description != "Your email for follow-up" {
		t.Errorf("Description = %q, want %q", i.Description, "Your email for follow-up")
	}

	if i.Placeholder != "user@example.com" {
		t.Errorf("Placeholder = %q, want %q", i.Placeholder, "user@example.com")
	}

	if i.Value != "default@example.com" {
		t.Errorf("Value = %q, want %q", i.Value, "default@example.com")
	}

	if !i.Required {
		t.Error("Required should be true")
	}
}

func TestTextarea_ElementType(t *testing.T) {
	ta := Textarea{}
	if got := ta.ElementType(); got != "textarea" {
		t.Errorf("ElementType() = %q, want %q", got, "textarea")
	}
}

func TestTextarea_AllFields(t *testing.T) {
	ta := Textarea{
		ID:          "code",
		Label:       "Code Sample",
		Description: "Paste your code here",
		Placeholder: "func main() { ... }",
		Value:       "// Default code",
		Render:      "go",
		Required:    true,
	}

	if ta.ID != "code" {
		t.Errorf("ID = %q, want %q", ta.ID, "code")
	}

	if ta.Label != "Code Sample" {
		t.Errorf("Label = %q, want %q", ta.Label, "Code Sample")
	}

	if ta.Description != "Paste your code here" {
		t.Errorf("Description = %q, want %q", ta.Description, "Paste your code here")
	}

	if ta.Render != "go" {
		t.Errorf("Render = %q, want %q", ta.Render, "go")
	}

	if !ta.Required {
		t.Error("Required should be true")
	}
}

func TestDropdown_ElementType(t *testing.T) {
	d := Dropdown{}
	if got := d.ElementType(); got != "dropdown" {
		t.Errorf("ElementType() = %q, want %q", got, "dropdown")
	}
}

func TestDropdown_AllFields(t *testing.T) {
	d := Dropdown{
		ID:          "severity",
		Label:       "Bug Severity",
		Description: "How severe is this bug?",
		Options:     []string{"Low", "Medium", "High", "Critical"},
		Multiple:    false,
		Default:     1,
		Required:    true,
	}

	if d.ID != "severity" {
		t.Errorf("ID = %q, want %q", d.ID, "severity")
	}

	if d.Label != "Bug Severity" {
		t.Errorf("Label = %q, want %q", d.Label, "Bug Severity")
	}

	if len(d.Options) != 4 {
		t.Errorf("len(Options) = %d, want 4", len(d.Options))
	}

	if d.Default != 1 {
		t.Errorf("Default = %d, want 1", d.Default)
	}

	if d.Multiple {
		t.Error("Multiple should be false")
	}

	if !d.Required {
		t.Error("Required should be true")
	}
}

func TestDropdown_Multiple(t *testing.T) {
	d := Dropdown{
		ID:       "platforms",
		Label:    "Affected Platforms",
		Options:  []string{"Windows", "macOS", "Linux", "iOS", "Android"},
		Multiple: true,
	}

	if !d.Multiple {
		t.Error("Multiple should be true")
	}
}

func TestCheckboxes_ElementType(t *testing.T) {
	c := Checkboxes{}
	if got := c.ElementType(); got != "checkboxes" {
		t.Errorf("ElementType() = %q, want %q", got, "checkboxes")
	}
}

func TestCheckboxes_AllFields(t *testing.T) {
	c := Checkboxes{
		ID:          "agreements",
		Label:       "Agreements",
		Description: "Please confirm the following",
		Options: []CheckboxOption{
			{Label: "I have searched for existing issues", Required: true},
			{Label: "I have read the contributing guidelines", Required: true},
			{Label: "I want to help fix this", Required: false},
		},
	}

	if c.ID != "agreements" {
		t.Errorf("ID = %q, want %q", c.ID, "agreements")
	}

	if c.Label != "Agreements" {
		t.Errorf("Label = %q, want %q", c.Label, "Agreements")
	}

	if len(c.Options) != 3 {
		t.Errorf("len(Options) = %d, want 3", len(c.Options))
	}

	if !c.Options[0].Required {
		t.Error("Options[0].Required should be true")
	}

	if c.Options[2].Required {
		t.Error("Options[2].Required should be false")
	}
}

func TestCheckboxOption(t *testing.T) {
	co := CheckboxOption{
		Label:    "I agree to the terms",
		Required: true,
	}

	if co.Label != "I agree to the terms" {
		t.Errorf("Label = %q, want %q", co.Label, "I agree to the terms")
	}

	if !co.Required {
		t.Error("Required should be true")
	}
}

func TestFormElement_Interface(t *testing.T) {
	elements := []FormElement{
		Markdown{Value: "# Hello"},
		Input{Label: "Name"},
		Textarea{Label: "Description"},
		Dropdown{Label: "Options"},
		Checkboxes{Label: "Checkboxes"},
	}

	expectedTypes := []string{"markdown", "input", "textarea", "dropdown", "checkboxes"}

	for i, elem := range elements {
		if got := elem.ElementType(); got != expectedTypes[i] {
			t.Errorf("elements[%d].ElementType() = %q, want %q", i, got, expectedTypes[i])
		}
	}
}

func TestIssueTemplate_ComplexBody(t *testing.T) {
	it := IssueTemplate{
		Name:        "Bug Report",
		Description: "Report a bug",
		Body: []FormElement{
			Markdown{Value: "## Bug Report\n\nPlease fill out the form below."},
			Input{ID: "title", Label: "Bug Title", Required: true},
			Textarea{ID: "description", Label: "Description", Required: true},
			Dropdown{ID: "severity", Label: "Severity", Options: []string{"Low", "Medium", "High"}},
			Checkboxes{
				ID:    "confirm",
				Label: "Confirmation",
				Options: []CheckboxOption{
					{Label: "I have searched for similar issues", Required: true},
				},
			},
		},
	}

	if len(it.Body) != 5 {
		t.Errorf("len(Body) = %d, want 5", len(it.Body))
	}

	// Verify element types
	if it.Body[0].ElementType() != "markdown" {
		t.Errorf("Body[0].ElementType() = %q, want %q", it.Body[0].ElementType(), "markdown")
	}

	if it.Body[1].ElementType() != "input" {
		t.Errorf("Body[1].ElementType() = %q, want %q", it.Body[1].ElementType(), "input")
	}

	if it.Body[2].ElementType() != "textarea" {
		t.Errorf("Body[2].ElementType() = %q, want %q", it.Body[2].ElementType(), "textarea")
	}

	if it.Body[3].ElementType() != "dropdown" {
		t.Errorf("Body[3].ElementType() = %q, want %q", it.Body[3].ElementType(), "dropdown")
	}

	if it.Body[4].ElementType() != "checkboxes" {
		t.Errorf("Body[4].ElementType() = %q, want %q", it.Body[4].ElementType(), "checkboxes")
	}
}

func TestPRTemplate_ResourceType(t *testing.T) {
	pr := PRTemplate{}
	if got := pr.ResourceType(); got != "pr-template" {
		t.Errorf("ResourceType() = %q, want %q", got, "pr-template")
	}
}

func TestPRTemplate_Filename_Default(t *testing.T) {
	tests := []struct {
		name     string
		template PRTemplate
		want     string
	}{
		{
			name:     "empty name",
			template: PRTemplate{},
			want:     "PULL_REQUEST_TEMPLATE.md",
		},
		{
			name:     "default name",
			template: PRTemplate{Name: "default"},
			want:     "PULL_REQUEST_TEMPLATE.md",
		},
		{
			name:     "named template",
			template: PRTemplate{Name: "feature"},
			want:     "PULL_REQUEST_TEMPLATE/feature.md",
		},
		{
			name:     "named template with spaces",
			template: PRTemplate{Name: "bug-fix"},
			want:     "PULL_REQUEST_TEMPLATE/bug-fix.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.template.Filename(); got != tt.want {
				t.Errorf("Filename() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPRTemplate_Content(t *testing.T) {
	pr := PRTemplate{
		Name: "feature",
		Content: `## Description

Please include a summary of the change.

## Test Plan

How was this tested?
`,
	}

	if pr.Name != "feature" {
		t.Errorf("Name = %q, want %q", pr.Name, "feature")
	}

	if pr.Content == "" {
		t.Error("Content should not be empty")
	}

	if len(pr.Content) < 50 {
		t.Errorf("Content too short: %d chars", len(pr.Content))
	}
}
