package templates

// PRTemplate represents a GitHub Pull Request template.
// Unlike Issue/Discussion templates which use YAML forms,
// PR templates are plain Markdown files.
type PRTemplate struct {
	// Name is the template identifier.
	// For the default template, leave empty and Content will be written to PULL_REQUEST_TEMPLATE.md
	// For named templates, the file will be PULL_REQUEST_TEMPLATE/{Name}.md
	Name string

	// Content is the Markdown content of the template.
	Content string
}

// ResourceType returns "pr-template" for interface compliance.
func (t PRTemplate) ResourceType() string {
	return "pr-template"
}

// Filename returns the appropriate filename for this template.
func (t PRTemplate) Filename() string {
	if t.Name == "" || t.Name == "default" {
		return "PULL_REQUEST_TEMPLATE.md"
	}
	return "PULL_REQUEST_TEMPLATE/" + t.Name + ".md"
}
