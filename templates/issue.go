// Package templates provides types for GitHub issue and discussion templates.
package templates

// IssueTemplate represents a GitHub Issue Form template.
type IssueTemplate struct {
	// Name is the template identifier (required, min 1 char).
	Name string `yaml:"name"`

	// Description is the template overview (required, min 1 char).
	Description string `yaml:"description"`

	// Title is the default issue title (optional).
	Title string `yaml:"title,omitempty"`

	// Labels are auto-applied labels (string array or comma-delimited string).
	Labels []string `yaml:"labels,omitempty"`

	// Projects are associated project board paths.
	Projects []string `yaml:"projects,omitempty"`

	// Assignees are default assignees (usernames).
	Assignees []string `yaml:"assignees,omitempty"`

	// Body contains the form fields (required, min 1 item).
	Body []FormElement `yaml:"body"`
}

// ResourceType returns "issue-template" for interface compliance.
func (t IssueTemplate) ResourceType() string {
	return "issue-template"
}
