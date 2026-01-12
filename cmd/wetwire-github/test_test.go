package main

import (
	"encoding/json"
	"testing"

	wetwire "github.com/lex00/wetwire-github-go"
)

// Note: Tests for the old local validation test command have been removed.
// The new test.go is an AI-assisted persona-based testing command that:
// - Requires ANTHROPIC_API_KEY to run
// - Takes a prompt argument instead of a path
// - Does not have --list, --score, --format flags
// The behavior is covered by manual testing and integration tests.

// TestTestResult_JSON tests TestResult JSON serialization.
func TestTestResult_JSON(t *testing.T) {
	result := wetwire.TestResult{
		Success: true,
		Tests: []wetwire.TestCase{
			{
				Name:    "workflows_exist",
				Persona: "beginner",
				Passed:  true,
			},
			{
				Name:    "jobs_exist",
				Persona: "beginner",
				Passed:  true,
			},
		},
		Passed: 2,
		Failed: 0,
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal TestResult: %v", err)
	}

	var unmarshaled wetwire.TestResult
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal TestResult: %v", err)
	}

	if !unmarshaled.Success {
		t.Error("Expected success to be true")
	}

	if unmarshaled.Passed != 2 {
		t.Errorf("Passed = %d, want 2", unmarshaled.Passed)
	}

	if unmarshaled.Failed != 0 {
		t.Errorf("Failed = %d, want 0", unmarshaled.Failed)
	}

	if len(unmarshaled.Tests) != 2 {
		t.Errorf("len(Tests) = %d, want 2", len(unmarshaled.Tests))
	}
}
