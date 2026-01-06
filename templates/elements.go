package templates

// FormElement is the interface for all form body elements.
type FormElement interface {
	ElementType() string
}

// Markdown represents a markdown content element.
type Markdown struct {
	// ID is an optional identifier.
	ID string `yaml:"id,omitempty"`

	// Value is the markdown/HTML content (required).
	Value string `yaml:"-"` // Handled in serialization as attributes.value
}

// ElementType returns "markdown".
func (m Markdown) ElementType() string {
	return "markdown"
}

// Input represents a single-line text input field.
type Input struct {
	// ID is an optional identifier.
	ID string `yaml:"id,omitempty"`

	// Label is the field title (required).
	Label string `yaml:"-"` // Handled in serialization

	// Description is optional extended help text.
	Description string `yaml:"-"`

	// Placeholder is an optional input hint.
	Placeholder string `yaml:"-"`

	// Value is the default text.
	Value string `yaml:"-"`

	// Required indicates if the field must be filled.
	Required bool `yaml:"-"`
}

// ElementType returns "input".
func (i Input) ElementType() string {
	return "input"
}

// Textarea represents a multi-line text input field.
type Textarea struct {
	// ID is an optional identifier.
	ID string `yaml:"id,omitempty"`

	// Label is the field title (required).
	Label string `yaml:"-"`

	// Description is optional extended help text.
	Description string `yaml:"-"`

	// Placeholder is an optional input hint.
	Placeholder string `yaml:"-"`

	// Value is the default content.
	Value string `yaml:"-"`

	// Render specifies syntax highlighting language (e.g., "python", "javascript").
	Render string `yaml:"-"`

	// Required indicates if the field must be filled.
	Required bool `yaml:"-"`
}

// ElementType returns "textarea".
func (t Textarea) ElementType() string {
	return "textarea"
}

// Dropdown represents a dropdown selection field.
type Dropdown struct {
	// ID is an optional identifier.
	ID string `yaml:"id,omitempty"`

	// Label is the field title (required).
	Label string `yaml:"-"`

	// Description is optional extended help text.
	Description string `yaml:"-"`

	// Options are the available choices (required, min 1).
	Options []string `yaml:"-"`

	// Multiple allows selecting multiple options.
	Multiple bool `yaml:"-"`

	// Default is the index of the pre-selected option.
	Default int `yaml:"-"`

	// Required indicates if a selection must be made.
	Required bool `yaml:"-"`
}

// ElementType returns "dropdown".
func (d Dropdown) ElementType() string {
	return "dropdown"
}

// Checkboxes represents a group of checkboxes.
type Checkboxes struct {
	// ID is an optional identifier.
	ID string `yaml:"id,omitempty"`

	// Label is the group title (required).
	Label string `yaml:"-"`

	// Description is optional extended help text.
	Description string `yaml:"-"`

	// Options are the checkbox items (required, min 1).
	Options []CheckboxOption `yaml:"-"`
}

// ElementType returns "checkboxes".
func (c Checkboxes) ElementType() string {
	return "checkboxes"
}

// CheckboxOption represents a single checkbox option.
type CheckboxOption struct {
	// Label is the checkbox text (required).
	Label string `yaml:"label"`

	// Required indicates if this checkbox must be checked.
	Required bool `yaml:"required,omitempty"`
}
