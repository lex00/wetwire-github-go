package importer

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Parser parses GitHub Actions YAML into intermediate representation.
type Parser struct{}

// NewParser creates a new Parser.
func NewParser() *Parser {
	return &Parser{}
}

// ParseFile parses a workflow YAML file from disk.
func (p *Parser) ParseFile(path string) (*IRWorkflow, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}
	return p.Parse(content)
}

// Parse parses workflow YAML content.
func (p *Parser) Parse(content []byte) (*IRWorkflow, error) {
	// First, parse into a raw map to handle flexible trigger syntax
	var rawMap map[string]any
	if err := yaml.Unmarshal(content, &rawMap); err != nil {
		return nil, fmt.Errorf("parsing yaml: %w", err)
	}

	// Now parse into the structured type
	var workflow IRWorkflow
	if err := yaml.Unmarshal(content, &workflow); err != nil {
		return nil, fmt.Errorf("parsing workflow: %w", err)
	}

	// Handle flexible matrix syntax
	if jobs, ok := rawMap["jobs"].(map[string]any); ok {
		for jobName, jobData := range jobs {
			if jobMap, ok := jobData.(map[string]any); ok {
				if strategy, ok := jobMap["strategy"].(map[string]any); ok {
					if matrix, ok := strategy["matrix"].(map[string]any); ok {
						if workflow.Jobs[jobName] != nil && workflow.Jobs[jobName].Strategy != nil {
							p.parseMatrix(workflow.Jobs[jobName].Strategy.Matrix, matrix)
						}
					}
				}
			}
		}
	}

	return &workflow, nil
}


// parseMatrix extracts matrix values from the raw matrix map.
func (p *Parser) parseMatrix(matrix *IRMatrix, raw map[string]any) {
	if matrix == nil {
		return
	}
	matrix.Raw = raw
	matrix.Values = make(map[string][]any)

	for key, value := range raw {
		// Skip special keys
		if key == "include" || key == "exclude" {
			continue
		}
		// Extract array values
		if arr, ok := value.([]any); ok {
			matrix.Values[key] = arr
		}
	}
}

// ParseWorkflow is a convenience function for parsing a workflow.
func ParseWorkflow(content []byte) (*IRWorkflow, error) {
	return NewParser().Parse(content)
}

// ParseWorkflowFile is a convenience function for parsing a workflow file.
func ParseWorkflowFile(path string) (*IRWorkflow, error) {
	return NewParser().ParseFile(path)
}

// ReferenceGraph tracks references between workflow elements.
type ReferenceGraph struct {
	// JobDependencies maps job names to their dependencies
	JobDependencies map[string][]string
	// StepOutputs maps step IDs to their output names
	StepOutputs map[string][]string
	// UsedActions lists all action references
	UsedActions []string
}

// BuildReferenceGraph builds a reference graph from a workflow.
func BuildReferenceGraph(workflow *IRWorkflow) *ReferenceGraph {
	graph := &ReferenceGraph{
		JobDependencies: make(map[string][]string),
		StepOutputs:     make(map[string][]string),
		UsedActions:     []string{},
	}

	for jobName, job := range workflow.Jobs {
		// Track job dependencies
		graph.JobDependencies[jobName] = job.GetNeeds()

		// Track step outputs and actions
		for _, step := range job.Steps {
			if step.ID != "" {
				// We don't know output names from YAML alone
				graph.StepOutputs[step.ID] = []string{}
			}
			if step.Uses != "" {
				graph.UsedActions = append(graph.UsedActions, step.Uses)
			}
		}

		// Track reusable workflow
		if job.Uses != "" {
			graph.UsedActions = append(graph.UsedActions, job.Uses)
		}
	}

	return graph
}
