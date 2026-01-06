// Package workflow provides typed Go declarations for GitHub Actions workflows.
package workflow

// Env is a shorthand for map[string]any used in environment variables.
type Env = map[string]any

// With is a shorthand for map[string]any used in action inputs.
type With = map[string]any

// List creates a typed slice from items.
// Use this instead of slice literals for cleaner declarations.
//
// Example:
//
//	Branches: List("main", "develop")
func List[T any](items ...T) []T {
	return items
}

// Strings creates a []string slice from items.
// Convenience wrapper around List for string slices.
func Strings(items ...string) []string {
	return items
}

// Ptr returns a pointer to the value.
// Use this when a field requires a pointer type.
func Ptr[T any](v T) *T {
	return &v
}
