package wetwire

import "fmt"

// WorkflowResource represents a GitHub workflow resource.
// All resource types (Workflow, Job) implement this interface.
type WorkflowResource interface {
	ResourceType() string // e.g., "workflow", "job"
}

// OutputRef represents a reference to a step output.
// When serialized to YAML, becomes: ${{ steps.step_id.outputs.name }}
type OutputRef struct {
	StepID string
	Output string
}

// String returns the GitHub Actions expression for this output reference.
func (o OutputRef) String() string {
	return fmt.Sprintf("${{ steps.%s.outputs.%s }}", o.StepID, o.Output)
}

// DiscoveredWorkflow represents a workflow found by AST parsing.
type DiscoveredWorkflow struct {
	Name string   // Variable name
	File string   // Source file path
	Line int      // Line number
	Jobs []string // Job variable names in this workflow
}

// DiscoveredJob represents a job found by AST parsing.
type DiscoveredJob struct {
	Name         string   // Variable name
	File         string   // Source file path
	Line         int      // Line number
	Dependencies []string // Referenced job names (Needs field)
}

// Result types for CLI JSON output.

// BuildResult contains the result of a build operation.
type BuildResult struct {
	Success   bool     `json:"success"`
	Workflows []string `json:"workflows,omitempty"`
	Files     []string `json:"files,omitempty"`
	Errors    []string `json:"errors,omitempty"`
}

// LintResult contains the result of a lint operation.
type LintResult struct {
	Success bool        `json:"success"`
	Issues  []LintIssue `json:"issues,omitempty"`
}

// LintIssue represents a single lint issue.
type LintIssue struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Severity string `json:"severity"` // "error", "warning", "info"
	Message  string `json:"message"`
	Rule     string `json:"rule"`
	Fixable  bool   `json:"fixable"`
}

// ValidateResult contains the result of a validate operation.
type ValidateResult struct {
	Success  bool              `json:"success"`
	Errors   []ValidationError `json:"errors,omitempty"`
	Warnings []string          `json:"warnings,omitempty"`
}

// ValidationError represents a validation error from actionlint.
type ValidationError struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Column  int    `json:"column"`
	Message string `json:"message"`
	RuleID  string `json:"rule_id,omitempty"`
}

// ListResult contains the result of a list operation.
type ListResult struct {
	Workflows []ListWorkflow `json:"workflows"`
}

// ListWorkflow represents a workflow in list output.
type ListWorkflow struct {
	Name string `json:"name"`
	File string `json:"file"`
	Line int    `json:"line"`
	Jobs int    `json:"jobs"`
}

// ImportResult contains the result of an import operation.
type ImportResult struct {
	Success      bool     `json:"success"`
	OutputDir    string   `json:"output_dir,omitempty"`
	Files        []string `json:"files,omitempty"`
	Workflows    int      `json:"workflows"`
	Jobs         int      `json:"jobs"`
	Steps        int      `json:"steps"`
	Errors       []string `json:"errors,omitempty"`
}

// GraphResult contains the result of a graph operation.
type GraphResult struct {
	Success bool   `json:"success"`
	Format  string `json:"format"` // "dot" or "mermaid"
	Output  string `json:"output"`
	Nodes   int    `json:"nodes"`
	Edges   int    `json:"edges"`
}
