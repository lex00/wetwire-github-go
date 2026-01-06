package template

import (
	"reflect"
	"testing"
)

func TestNewGraph(t *testing.T) {
	g := NewGraph()
	if g == nil {
		t.Fatal("NewGraph() returned nil")
	}
	if g.Nodes == nil {
		t.Error("NewGraph().Nodes is nil")
	}
	if g.Edges == nil {
		t.Error("NewGraph().Edges is nil")
	}
}

func TestGraph_AddNode(t *testing.T) {
	g := NewGraph()
	g.AddNode("A")
	g.AddNode("B")

	if !g.Nodes["A"] {
		t.Error("Node A not added")
	}
	if !g.Nodes["B"] {
		t.Error("Node B not added")
	}
}

func TestGraph_AddEdge(t *testing.T) {
	g := NewGraph()
	g.AddEdge("A", "B") // A depends on B

	if !g.Nodes["A"] {
		t.Error("Node A not added")
	}
	if !g.Nodes["B"] {
		t.Error("Node B not added")
	}
	if len(g.Edges["A"]) != 1 || g.Edges["A"][0] != "B" {
		t.Errorf("Edge A->B not added correctly: %v", g.Edges["A"])
	}
}

func TestGraph_TopologicalSortKahn_Simple(t *testing.T) {
	g := NewGraph()
	// B depends on A (A must come first)
	g.AddNode("A")
	g.AddEdge("B", "A")

	result, err := g.TopologicalSortKahn()
	if err != nil {
		t.Fatalf("TopologicalSortKahn() error = %v", err)
	}

	// A should come before B
	if len(result) != 2 {
		t.Fatalf("Expected 2 nodes, got %d", len(result))
	}
	if result[0] != "A" {
		t.Errorf("Expected A first, got %s", result[0])
	}
	if result[1] != "B" {
		t.Errorf("Expected B second, got %s", result[1])
	}
}

func TestGraph_TopologicalSortKahn_Chain(t *testing.T) {
	g := NewGraph()
	// C depends on B, B depends on A
	g.AddEdge("B", "A")
	g.AddEdge("C", "B")

	result, err := g.TopologicalSortKahn()
	if err != nil {
		t.Fatalf("TopologicalSortKahn() error = %v", err)
	}

	expected := []string{"A", "B", "C"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("TopologicalSortKahn() = %v, want %v", result, expected)
	}
}

func TestGraph_TopologicalSortKahn_Diamond(t *testing.T) {
	g := NewGraph()
	// Diamond pattern: D depends on B and C, B and C depend on A
	g.AddEdge("B", "A")
	g.AddEdge("C", "A")
	g.AddEdge("D", "B")
	g.AddEdge("D", "C")

	result, err := g.TopologicalSortKahn()
	if err != nil {
		t.Fatalf("TopologicalSortKahn() error = %v", err)
	}

	// A must come first, D must come last, B and C in between
	if len(result) != 4 {
		t.Fatalf("Expected 4 nodes, got %d", len(result))
	}
	if result[0] != "A" {
		t.Errorf("Expected A first, got %s", result[0])
	}
	if result[3] != "D" {
		t.Errorf("Expected D last, got %s", result[3])
	}

	// B and C should be in positions 1 and 2 (order doesn't matter but should be deterministic)
	middle := []string{result[1], result[2]}
	expectedMiddle := []string{"B", "C"}
	if !reflect.DeepEqual(middle, expectedMiddle) {
		t.Errorf("Expected B and C in middle, got %v", middle)
	}
}

func TestGraph_TopologicalSortKahn_MultipleDependencies(t *testing.T) {
	g := NewGraph()
	// Deploy depends on Build and Test
	// Test depends on Build
	g.AddEdge("Test", "Build")
	g.AddEdge("Deploy", "Build")
	g.AddEdge("Deploy", "Test")

	result, err := g.TopologicalSortKahn()
	if err != nil {
		t.Fatalf("TopologicalSortKahn() error = %v", err)
	}

	// Build must come first, Deploy must come last
	if result[0] != "Build" {
		t.Errorf("Expected Build first, got %s", result[0])
	}
	if result[len(result)-1] != "Deploy" {
		t.Errorf("Expected Deploy last, got %s", result[len(result)-1])
	}
}

func TestGraph_TopologicalSortKahn_NoDependencies(t *testing.T) {
	g := NewGraph()
	g.AddNode("A")
	g.AddNode("B")
	g.AddNode("C")

	result, err := g.TopologicalSortKahn()
	if err != nil {
		t.Fatalf("TopologicalSortKahn() error = %v", err)
	}

	// All nodes present, order is alphabetical (deterministic)
	expected := []string{"A", "B", "C"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("TopologicalSortKahn() = %v, want %v", result, expected)
	}
}

func TestGraph_TopologicalSortKahn_Cycle(t *testing.T) {
	g := NewGraph()
	// A -> B -> C -> A (cycle)
	g.AddEdge("A", "B")
	g.AddEdge("B", "C")
	g.AddEdge("C", "A")

	_, err := g.TopologicalSortKahn()
	if err == nil {
		t.Error("TopologicalSortKahn() expected error for cycle")
	}
}

func TestGraph_DetectCycles_NoCycle(t *testing.T) {
	g := NewGraph()
	g.AddEdge("B", "A")
	g.AddEdge("C", "B")

	cycles := g.DetectCycles()
	if len(cycles) != 0 {
		t.Errorf("DetectCycles() found cycles in acyclic graph: %v", cycles)
	}
}

func TestGraph_DetectCycles_SimpleCycle(t *testing.T) {
	g := NewGraph()
	// A -> B -> A (cycle)
	g.AddEdge("A", "B")
	g.AddEdge("B", "A")

	cycles := g.DetectCycles()
	if len(cycles) == 0 {
		t.Error("DetectCycles() did not find cycle")
	}
}

func TestGraph_DetectCycles_SelfLoop(t *testing.T) {
	g := NewGraph()
	g.AddEdge("A", "A")

	cycles := g.DetectCycles()
	if len(cycles) == 0 {
		t.Error("DetectCycles() did not find self-loop")
	}
}

func TestGraph_Empty(t *testing.T) {
	g := NewGraph()

	result, err := g.TopologicalSortKahn()
	if err != nil {
		t.Fatalf("TopologicalSortKahn() error = %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected empty result, got %v", result)
	}

	cycles := g.DetectCycles()
	if len(cycles) != 0 {
		t.Errorf("Expected no cycles, got %v", cycles)
	}
}
