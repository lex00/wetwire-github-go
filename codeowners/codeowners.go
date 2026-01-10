// Package codeowners provides types for GitHub CODEOWNERS configuration.
package codeowners

// Owners represents a GitHub CODEOWNERS file.
// CODEOWNERS is used to define who is responsible for code in a repository.
type Owners struct {
	// Rules defines ownership rules in order of precedence.
	// Later rules take precedence over earlier ones for matching files.
	Rules []Rule
}

// ResourceType returns "codeowners" for interface compliance.
func (o Owners) ResourceType() string {
	return "codeowners"
}

// Rule defines ownership for a path pattern.
type Rule struct {
	// Pattern is the file/directory pattern to match.
	// Patterns follow .gitignore syntax:
	// - "*" matches any file
	// - "*.go" matches Go files
	// - "/docs/" matches the docs directory
	// - "src/**/*.ts" matches TypeScript files under src
	Pattern string

	// Owners is a list of users or teams who own matching files.
	// Users are specified as "@username".
	// Teams are specified as "@org/team-name".
	Owners []string

	// Comment is an optional comment to include above the rule.
	Comment string
}
