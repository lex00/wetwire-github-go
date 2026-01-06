package discover

import (
	"strings"
	"testing"
)

func TestNewDependencyGraph(t *testing.T) {
	jobs := []DiscoveredJob{
		{Name: "Build", Dependencies: []string{}},
		{Name: "Test", Dependencies: []string{"Build"}},
		{Name: "Deploy", Dependencies: []string{"Build", "Test"}},
	}

	g := NewDependencyGraph(jobs)

	if len(g.Nodes) != 3 {
		t.Errorf("len(g.Nodes) = %d, want 3", len(g.Nodes))
	}

	if len(g.Edges) != 3 {
		t.Errorf("len(g.Edges) = %d, want 3", len(g.Edges))
	}

	if len(g.Edges["Deploy"]) != 2 {
		t.Errorf("len(g.Edges[Deploy]) = %d, want 2", len(g.Edges["Deploy"]))
	}
}

func TestDependencyGraph_TopologicalSort(t *testing.T) {
	jobs := []DiscoveredJob{
		{Name: "Build", Dependencies: []string{}},
		{Name: "Test", Dependencies: []string{"Build"}},
		{Name: "Deploy", Dependencies: []string{"Test"}},
	}

	g := NewDependencyGraph(jobs)
	sorted, err := g.TopologicalSort()
	if err != nil {
		t.Fatalf("TopologicalSort() error = %v", err)
	}

	// Build should come before Test, Test before Deploy
	buildIdx := indexOf(sorted, "Build")
	testIdx := indexOf(sorted, "Test")
	deployIdx := indexOf(sorted, "Deploy")

	if buildIdx > testIdx {
		t.Error("Build should come before Test")
	}
	if testIdx > deployIdx {
		t.Error("Test should come before Deploy")
	}
}

func TestDependencyGraph_TopologicalSort_CycleDetection(t *testing.T) {
	jobs := []DiscoveredJob{
		{Name: "A", Dependencies: []string{"B"}},
		{Name: "B", Dependencies: []string{"A"}},
	}

	g := NewDependencyGraph(jobs)
	_, err := g.TopologicalSort()
	if err == nil {
		t.Error("TopologicalSort() expected error for cycle")
	}
}

func TestDependencyGraph_DetectCycles(t *testing.T) {
	// Test with a cycle
	jobs := []DiscoveredJob{
		{Name: "A", Dependencies: []string{"B"}},
		{Name: "B", Dependencies: []string{"C"}},
		{Name: "C", Dependencies: []string{"A"}},
	}

	g := NewDependencyGraph(jobs)
	cycles := g.DetectCycles()

	if len(cycles) == 0 {
		t.Error("DetectCycles() expected to find cycles")
	}
}

func TestDependencyGraph_DetectCycles_NoCycle(t *testing.T) {
	jobs := []DiscoveredJob{
		{Name: "A", Dependencies: []string{}},
		{Name: "B", Dependencies: []string{"A"}},
		{Name: "C", Dependencies: []string{"B"}},
	}

	g := NewDependencyGraph(jobs)
	cycles := g.DetectCycles()

	if len(cycles) != 0 {
		t.Errorf("DetectCycles() expected no cycles, got %d", len(cycles))
	}
}

func TestDependencyGraph_GetDependents(t *testing.T) {
	jobs := []DiscoveredJob{
		{Name: "Build", Dependencies: []string{}},
		{Name: "Test", Dependencies: []string{"Build"}},
		{Name: "Lint", Dependencies: []string{"Build"}},
		{Name: "Deploy", Dependencies: []string{"Test"}},
	}

	g := NewDependencyGraph(jobs)
	dependents := g.GetDependents("Build")

	if len(dependents) != 2 {
		t.Errorf("len(GetDependents(Build)) = %d, want 2", len(dependents))
	}

	// Check that Test and Lint are dependents
	hasTest := false
	hasLint := false
	for _, d := range dependents {
		if d == "Test" {
			hasTest = true
		}
		if d == "Lint" {
			hasLint = true
		}
	}

	if !hasTest || !hasLint {
		t.Errorf("GetDependents(Build) = %v, want [Lint, Test]", dependents)
	}
}

func TestDependencyGraph_GetAllDependencies(t *testing.T) {
	jobs := []DiscoveredJob{
		{Name: "Build", Dependencies: []string{}},
		{Name: "Test", Dependencies: []string{"Build"}},
		{Name: "Deploy", Dependencies: []string{"Test"}},
	}

	g := NewDependencyGraph(jobs)
	deps := g.GetAllDependencies("Deploy")

	if len(deps) != 2 {
		t.Errorf("len(GetAllDependencies(Deploy)) = %d, want 2", len(deps))
	}

	// Should include both Build and Test
	hasBuild := false
	hasTest := false
	for _, d := range deps {
		if d == "Build" {
			hasBuild = true
		}
		if d == "Test" {
			hasTest = true
		}
	}

	if !hasBuild || !hasTest {
		t.Errorf("GetAllDependencies(Deploy) = %v, want [Build, Test]", deps)
	}
}

func TestDependencyGraph_ToDOT(t *testing.T) {
	jobs := []DiscoveredJob{
		{Name: "Build", Dependencies: []string{}},
		{Name: "Test", Dependencies: []string{"Build"}},
	}

	g := NewDependencyGraph(jobs)
	dot := g.ToDOT("workflow")

	expectedStrings := []string{
		"digraph workflow {",
		`"Build"`,
		`"Test"`,
		`"Test" -> "Build"`,
		"}",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(dot, expected) {
			t.Errorf("ToDOT() missing %q\n\nGenerated:\n%s", expected, dot)
		}
	}
}

// Helper function
func indexOf(slice []string, item string) int {
	for i, s := range slice {
		if s == item {
			return i
		}
	}
	return -1
}
