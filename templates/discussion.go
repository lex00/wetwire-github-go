package templates

// DiscussionTemplate represents a GitHub Discussion Form template.
type DiscussionTemplate struct {
	// Title is the template title (required, min 1 char).
	Title string `yaml:"title"`

	// Description is the template overview (required, min 1 char).
	Description string `yaml:"description"`

	// Labels are auto-applied labels (string array or comma-delimited string).
	Labels []string `yaml:"labels,omitempty"`

	// Body contains the form fields (required, min 1 item).
	// Reuses the same FormElement types as IssueTemplate.
	Body []FormElement `yaml:"body"`
}

// ResourceType returns "discussion-template" for interface compliance.
func (t DiscussionTemplate) ResourceType() string {
	return "discussion-template"
}
