package workflow_test

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestNewMatrix(t *testing.T) {
	values := map[string][]any{
		"go": {"1.22", "1.23"},
		"os": {"ubuntu-latest", "macos-latest"},
	}

	m := workflow.NewMatrix(values)

	if m == nil {
		t.Fatal("expected matrix to be created")
	}

	if len(m.Values["go"]) != 2 {
		t.Errorf("expected 2 go versions, got %d", len(m.Values["go"]))
	}

	if len(m.Values["os"]) != 2 {
		t.Errorf("expected 2 os values, got %d", len(m.Values["os"]))
	}
}

func TestMatrixWithInclude(t *testing.T) {
	m := workflow.NewMatrix(map[string][]any{
		"go": {"1.22", "1.23"},
	})

	include1 := map[string]any{
		"go": "1.21",
		"experimental": true,
	}
	include2 := map[string]any{
		"go": "1.24-rc",
		"experimental": true,
	}

	m = m.WithInclude(include1, include2)

	if len(m.Include) != 2 {
		t.Errorf("expected 2 includes, got %d", len(m.Include))
	}

	if m.Include[0]["go"] != "1.21" {
		t.Errorf("expected first include to have go=1.21, got %v", m.Include[0]["go"])
	}

	if m.Include[1]["go"] != "1.24-rc" {
		t.Errorf("expected second include to have go=1.24-rc, got %v", m.Include[1]["go"])
	}
}

func TestMatrixWithExclude(t *testing.T) {
	m := workflow.NewMatrix(map[string][]any{
		"go": {"1.22", "1.23"},
		"os": {"ubuntu-latest", "macos-latest", "windows-latest"},
	})

	exclude1 := map[string]any{
		"go": "1.22",
		"os": "windows-latest",
	}
	exclude2 := map[string]any{
		"go": "1.23",
		"os": "macos-latest",
	}

	m = m.WithExclude(exclude1, exclude2)

	if len(m.Exclude) != 2 {
		t.Errorf("expected 2 excludes, got %d", len(m.Exclude))
	}

	if m.Exclude[0]["os"] != "windows-latest" {
		t.Errorf("expected first exclude to have os=windows-latest, got %v", m.Exclude[0]["os"])
	}
}

func TestMatrixChaining(t *testing.T) {
	m := workflow.NewMatrix(map[string][]any{
		"go": {"1.22", "1.23"},
	}).WithInclude(
		map[string]any{"go": "1.21"},
	).WithExclude(
		map[string]any{"go": "1.22"},
	)

	if len(m.Values["go"]) != 2 {
		t.Errorf("expected 2 go versions in values, got %d", len(m.Values["go"]))
	}

	if len(m.Include) != 1 {
		t.Errorf("expected 1 include, got %d", len(m.Include))
	}

	if len(m.Exclude) != 1 {
		t.Errorf("expected 1 exclude, got %d", len(m.Exclude))
	}
}

func TestMatrixInStrategy(t *testing.T) {
	matrix := workflow.NewMatrix(map[string][]any{
		"go": {"1.22", "1.23"},
		"os": {"ubuntu-latest", "macos-latest"},
	})

	strategy := workflow.Strategy{
		Matrix:      matrix,
		FailFast:    workflow.Ptr(false),
		MaxParallel: 4,
	}

	if strategy.Matrix == nil {
		t.Error("expected matrix to be set")
	}

	if *strategy.FailFast {
		t.Error("expected FailFast to be false")
	}

	if strategy.MaxParallel != 4 {
		t.Errorf("expected MaxParallel=4, got %d", strategy.MaxParallel)
	}
}

func TestMatrixContextGet(t *testing.T) {
	expr := workflow.MatrixContext.Get("os")

	expected := "matrix.os"
	if expr.Raw() != expected {
		t.Errorf("expected %q, got %q", expected, expr.Raw())
	}

	// Test with multiple keys
	goExpr := workflow.MatrixContext.Get("go")
	if goExpr.Raw() != "matrix.go" {
		t.Errorf("expected 'matrix.go', got %q", goExpr.Raw())
	}
}
