package discover

import (
	"fmt"
	"sort"
	"strings"
)

// DependencyGraph represents the dependency relationships between jobs.
type DependencyGraph struct {
	// Nodes maps job names to their discovered job info
	Nodes map[string]*DiscoveredJob

	// Edges maps job names to their dependencies (jobs they depend on)
	Edges map[string][]string
}

// NewDependencyGraph creates a new dependency graph from discovered jobs.
func NewDependencyGraph(jobs []DiscoveredJob) *DependencyGraph {
	g := &DependencyGraph{
		Nodes: make(map[string]*DiscoveredJob),
		Edges: make(map[string][]string),
	}

	for i := range jobs {
		job := &jobs[i]
		g.Nodes[job.Name] = job
		g.Edges[job.Name] = job.Dependencies
	}

	return g
}

// TopologicalSort returns the jobs in dependency order.
// Jobs with no dependencies come first, then jobs that depend on them, etc.
func (g *DependencyGraph) TopologicalSort() ([]string, error) {
	// Track visited nodes and nodes in current path (for cycle detection)
	visited := make(map[string]bool)
	inPath := make(map[string]bool)
	result := []string{}

	var visit func(name string) error
	visit = func(name string) error {
		if inPath[name] {
			return fmt.Errorf("cycle detected: job %q depends on itself", name)
		}
		if visited[name] {
			return nil
		}

		inPath[name] = true

		// Visit dependencies first
		for _, dep := range g.Edges[name] {
			if err := visit(dep); err != nil {
				return err
			}
		}

		inPath[name] = false
		visited[name] = true
		result = append(result, name)

		return nil
	}

	// Get sorted node names for deterministic output
	names := make([]string, 0, len(g.Nodes))
	for name := range g.Nodes {
		names = append(names, name)
	}
	sort.Strings(names)

	// Visit all nodes
	for _, name := range names {
		if err := visit(name); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// DetectCycles returns any cycles found in the dependency graph.
func (g *DependencyGraph) DetectCycles() [][]string {
	var cycles [][]string
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	path := []string{}

	var detectCycle func(name string) bool
	detectCycle = func(name string) bool {
		visited[name] = true
		recStack[name] = true
		path = append(path, name)

		for _, dep := range g.Edges[name] {
			if !visited[dep] {
				if detectCycle(dep) {
					return true
				}
			} else if recStack[dep] {
				// Found a cycle - extract the cycle path
				cycleStart := 0
				for i, n := range path {
					if n == dep {
						cycleStart = i
						break
					}
				}
				cycle := make([]string, len(path)-cycleStart)
				copy(cycle, path[cycleStart:])
				cycle = append(cycle, dep) // Complete the cycle
				cycles = append(cycles, cycle)
				return true
			}
		}

		path = path[:len(path)-1]
		recStack[name] = false
		return false
	}

	// Get sorted node names for deterministic output
	names := make([]string, 0, len(g.Nodes))
	for name := range g.Nodes {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		if !visited[name] {
			detectCycle(name)
		}
	}

	return cycles
}

// GetDependents returns all jobs that depend on the given job.
func (g *DependencyGraph) GetDependents(jobName string) []string {
	var dependents []string

	for name, deps := range g.Edges {
		for _, dep := range deps {
			if dep == jobName {
				dependents = append(dependents, name)
				break
			}
		}
	}

	sort.Strings(dependents)
	return dependents
}

// GetAllDependencies returns all direct and transitive dependencies of a job.
func (g *DependencyGraph) GetAllDependencies(jobName string) []string {
	visited := make(map[string]bool)
	var result []string

	var collect func(name string)
	collect = func(name string) {
		for _, dep := range g.Edges[name] {
			if !visited[dep] {
				visited[dep] = true
				result = append(result, dep)
				collect(dep)
			}
		}
	}

	collect(jobName)
	sort.Strings(result)
	return result
}

// ToDOT generates a DOT graph representation.
func (g *DependencyGraph) ToDOT(graphName string) string {
	var sb strings.Builder

	sb.WriteString("digraph ")
	sb.WriteString(graphName)
	sb.WriteString(" {\n")
	sb.WriteString("  rankdir=TB;\n")
	sb.WriteString("  node [shape=box];\n\n")

	// Get sorted node names
	names := make([]string, 0, len(g.Nodes))
	for name := range g.Nodes {
		names = append(names, name)
	}
	sort.Strings(names)

	// Write nodes
	for _, name := range names {
		sb.WriteString(fmt.Sprintf("  %q;\n", name))
	}

	sb.WriteString("\n")

	// Write edges
	for _, name := range names {
		for _, dep := range g.Edges[name] {
			sb.WriteString(fmt.Sprintf("  %q -> %q;\n", name, dep))
		}
	}

	sb.WriteString("}\n")

	return sb.String()
}
