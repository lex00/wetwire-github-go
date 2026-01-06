package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	wetwire "github.com/lex00/wetwire-github-go"
	"github.com/lex00/wetwire-github-go/internal/discover"
)

var graphFormat string
var graphDirection string

var graphCmd = &cobra.Command{
	Use:   "graph <path>",
	Short: "Visualize workflow job dependencies as a DAG",
	Long: `Graph generates a visual representation of job dependencies.

Supported formats:
  dot     - Graphviz DOT format
  mermaid - Mermaid diagram format
  json    - JSON representation

Example:
  wetwire-github graph .
  wetwire-github graph . --format dot | dot -Tpng -o workflow.png
  wetwire-github graph . --format mermaid
  wetwire-github graph . --direction LR`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGraph(args[0])
	},
}

func init() {
	graphCmd.Flags().StringVar(&graphFormat, "format", "dot", "output format (dot, mermaid, json)")
	graphCmd.Flags().StringVar(&graphDirection, "direction", "TB", "graph direction (TB, LR)")
}

// runGraph executes the graph command.
func runGraph(path string) error {
	// Resolve absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: resolving path: %v\n", err)
		os.Exit(1)
		return nil
	}

	// Check if path exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "error: path not found: %s\n", path)
		os.Exit(1)
		return nil
	}

	// Discover jobs
	disc := discover.NewDiscoverer()
	discovered, err := disc.Discover(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: discovery failed: %v\n", err)
		os.Exit(1)
		return nil
	}

	// Build graph
	graph := discover.NewDependencyGraph(discovered.Jobs)

	// Count nodes and edges
	nodes := len(graph.Nodes)
	edges := 0
	for _, deps := range graph.Edges {
		edges += len(deps)
	}

	switch graphFormat {
	case "dot":
		fmt.Print(generateDOT(graph, graphDirection))
	case "mermaid":
		fmt.Print(generateMermaid(graph, graphDirection))
	case "json":
		result := wetwire.GraphResult{
			Success: true,
			Format:  "json",
			Output:  generateDOT(graph, graphDirection),
			Nodes:   nodes,
			Edges:   edges,
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	default:
		fmt.Fprintf(os.Stderr, "error: unknown format: %s\n", graphFormat)
		os.Exit(1)
	}

	return nil
}

// generateDOT generates DOT format output.
func generateDOT(graph *discover.DependencyGraph, direction string) string {
	var sb strings.Builder

	sb.WriteString("digraph workflow {\n")
	sb.WriteString(fmt.Sprintf("  rankdir=%s;\n", direction))
	sb.WriteString("  node [shape=box];\n\n")

	// Get sorted node names
	nodeNames := make([]string, 0, len(graph.Nodes))
	for name := range graph.Nodes {
		nodeNames = append(nodeNames, name)
	}
	sort.Strings(nodeNames)

	// Write nodes
	for _, name := range nodeNames {
		sb.WriteString(fmt.Sprintf("  %q;\n", name))
	}

	sb.WriteString("\n")

	// Write edges (node -> dependency means dependency must run first)
	for _, name := range nodeNames {
		deps := graph.Edges[name]
		if len(deps) > 0 {
			sortedDeps := make([]string, len(deps))
			copy(sortedDeps, deps)
			sort.Strings(sortedDeps)
			for _, dep := range sortedDeps {
				// Edge from dependency to dependent (dep runs before name)
				sb.WriteString(fmt.Sprintf("  %q -> %q;\n", dep, name))
			}
		}
	}

	sb.WriteString("}\n")
	return sb.String()
}

// generateMermaid generates Mermaid diagram format output.
func generateMermaid(graph *discover.DependencyGraph, direction string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("graph %s\n", direction))

	// Get sorted node names
	nodeNames := make([]string, 0, len(graph.Nodes))
	for name := range graph.Nodes {
		nodeNames = append(nodeNames, name)
	}
	sort.Strings(nodeNames)

	// Write edges
	hasEdges := false
	for _, name := range nodeNames {
		deps := graph.Edges[name]
		if len(deps) > 0 {
			sortedDeps := make([]string, len(deps))
			copy(sortedDeps, deps)
			sort.Strings(sortedDeps)
			for _, dep := range sortedDeps {
				sb.WriteString(fmt.Sprintf("    %s --> %s\n", dep, name))
				hasEdges = true
			}
		}
	}

	// If no edges, just list nodes
	if !hasEdges {
		for _, name := range nodeNames {
			sb.WriteString(fmt.Sprintf("    %s\n", name))
		}
	}

	return sb.String()
}
