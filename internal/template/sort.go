// Package template provides workflow template building with dependency ordering.
package template

import (
	"fmt"
	"sort"
)

// Graph represents a directed acyclic graph for topological sorting.
type Graph struct {
	// Nodes is the set of all node names
	Nodes map[string]bool

	// Edges maps node -> nodes it depends on
	Edges map[string][]string
}

// NewGraph creates a new empty Graph.
func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string]bool),
		Edges: make(map[string][]string),
	}
}

// AddNode adds a node to the graph.
func (g *Graph) AddNode(name string) {
	g.Nodes[name] = true
	if _, exists := g.Edges[name]; !exists {
		g.Edges[name] = []string{}
	}
}

// AddEdge adds a dependency edge: node depends on dependency.
func (g *Graph) AddEdge(node, dependency string) {
	g.AddNode(node)
	g.AddNode(dependency)
	g.Edges[node] = append(g.Edges[node], dependency)
}

// TopologicalSortKahn performs topological sort using Kahn's algorithm.
// Returns nodes in dependency order (dependencies first).
// Returns error if a cycle is detected.
func (g *Graph) TopologicalSortKahn() ([]string, error) {
	// Calculate in-degrees (how many nodes depend on each node)
	// Note: We're tracking "what depends on this" not "what this depends on"
	inDegree := make(map[string]int)
	dependents := make(map[string][]string) // node -> nodes that depend on it

	for node := range g.Nodes {
		inDegree[node] = 0
	}

	for node, deps := range g.Edges {
		for _, dep := range deps {
			dependents[dep] = append(dependents[dep], node)
			inDegree[node]++
		}
	}

	// Find all nodes with no dependencies (in-degree 0)
	var queue []string
	for node := range g.Nodes {
		if inDegree[node] == 0 {
			queue = append(queue, node)
		}
	}

	// Sort for deterministic output
	sort.Strings(queue)

	var result []string
	for len(queue) > 0 {
		// Take first node from queue
		node := queue[0]
		queue = queue[1:]
		result = append(result, node)

		// For each node that depends on this node
		deps := dependents[node]
		sort.Strings(deps)
		for _, dependent := range deps {
			inDegree[dependent]--
			if inDegree[dependent] == 0 {
				queue = append(queue, dependent)
				sort.Strings(queue)
			}
		}
	}

	// If we didn't process all nodes, there's a cycle
	if len(result) != len(g.Nodes) {
		cycle := g.findCycle()
		return nil, fmt.Errorf("cycle detected: %v", cycle)
	}

	return result, nil
}

// findCycle finds a cycle in the graph (for error reporting).
func (g *Graph) findCycle() []string {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	var path []string

	var dfs func(node string) []string
	dfs = func(node string) []string {
		visited[node] = true
		recStack[node] = true
		path = append(path, node)

		for _, dep := range g.Edges[node] {
			if !visited[dep] {
				if cycle := dfs(dep); cycle != nil {
					return cycle
				}
			} else if recStack[dep] {
				// Found cycle - extract the cycle portion
				cycleStart := 0
				for i, n := range path {
					if n == dep {
						cycleStart = i
						break
					}
				}
				cycle := make([]string, len(path)-cycleStart)
				copy(cycle, path[cycleStart:])
				return append(cycle, dep)
			}
		}

		path = path[:len(path)-1]
		recStack[node] = false
		return nil
	}

	// Get sorted node names for deterministic output
	nodes := make([]string, 0, len(g.Nodes))
	for node := range g.Nodes {
		nodes = append(nodes, node)
	}
	sort.Strings(nodes)

	for _, node := range nodes {
		if !visited[node] {
			if cycle := dfs(node); cycle != nil {
				return cycle
			}
		}
	}

	return nil
}

// DetectCycles returns all cycles found in the graph.
func (g *Graph) DetectCycles() [][]string {
	var cycles [][]string
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	var path []string

	var detectCycle func(node string) bool
	detectCycle = func(node string) bool {
		visited[node] = true
		recStack[node] = true
		path = append(path, node)

		for _, dep := range g.Edges[node] {
			if !visited[dep] {
				if detectCycle(dep) {
					return true
				}
			} else if recStack[dep] {
				// Found cycle
				cycleStart := 0
				for i, n := range path {
					if n == dep {
						cycleStart = i
						break
					}
				}
				cycle := make([]string, len(path)-cycleStart)
				copy(cycle, path[cycleStart:])
				cycle = append(cycle, dep)
				cycles = append(cycles, cycle)
				return true
			}
		}

		path = path[:len(path)-1]
		recStack[node] = false
		return false
	}

	// Get sorted node names for deterministic output
	nodes := make([]string, 0, len(g.Nodes))
	for node := range g.Nodes {
		nodes = append(nodes, node)
	}
	sort.Strings(nodes)

	for _, node := range nodes {
		if !visited[node] {
			detectCycle(node)
		}
	}

	return cycles
}
