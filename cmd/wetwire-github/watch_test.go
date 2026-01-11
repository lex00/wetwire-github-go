package main

import (
	"testing"
	"time"
)

func TestFormatTimestamp(t *testing.T) {
	now := time.Date(2024, 1, 15, 14, 30, 45, 0, time.UTC)
	result := formatTimestamp(now)
	expected := "14:30:45"
	if result != expected {
		t.Errorf("formatTimestamp(%v) = %q, want %q", now, result, expected)
	}
}

func TestFormatTimestamp_Midnight(t *testing.T) {
	midnight := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	result := formatTimestamp(midnight)
	expected := "00:00:00"
	if result != expected {
		t.Errorf("formatTimestamp(%v) = %q, want %q", midnight, result, expected)
	}
}

func TestFormatTimestamp_EndOfDay(t *testing.T) {
	endOfDay := time.Date(2024, 1, 15, 23, 59, 59, 0, time.UTC)
	result := formatTimestamp(endOfDay)
	expected := "23:59:59"
	if result != expected {
		t.Errorf("formatTimestamp(%v) = %q, want %q", endOfDay, result, expected)
	}
}

func TestShouldProcessEvent(t *testing.T) {
	tests := []struct {
		op     string
		path   string
		expect bool
	}{
		{"CREATE", "workflow.go", true},
		{"WRITE", "jobs.go", true},
		{"REMOVE", "triggers.go", true},
		{"RENAME", "old.go", true},
		{"CHMOD", "workflow.go", false},
		{"CREATE", "readme.md", false},
		{"WRITE", "config.yaml", false},
		{"CREATE", ".hidden.go", true},
		{"WRITE", "/path/to/file.go", true},
		{"CREATE", "test_test.go", true},
		{"WRITE", "file.go.bak", false},
		{"", "file.go", false},
	}

	for _, tt := range tests {
		t.Run(tt.op+"_"+tt.path, func(t *testing.T) {
			result := shouldProcessEvent(tt.op, tt.path)
			if result != tt.expect {
				t.Errorf("shouldProcessEvent(%q, %q) = %v, want %v", tt.op, tt.path, result, tt.expect)
			}
		})
	}
}

func TestIsGoFile(t *testing.T) {
	tests := []struct {
		path   string
		expect bool
	}{
		{"workflow.go", true},
		{"main.go", true},
		{"test_test.go", true},
		{"/path/to/file.go", true},
		{"readme.md", false},
		{"config.yaml", false},
		{"file.go.bak", false},
		{"", false},
		{".hidden.go", true},
		{"GO", false},
		{"file.GO", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := isGoFile(tt.path)
			if result != tt.expect {
				t.Errorf("isGoFile(%q) = %v, want %v", tt.path, result, tt.expect)
			}
		})
	}
}

func TestRunWatchBuild_InvalidPath(t *testing.T) {
	err := runWatchBuild("/nonexistent/path", ".github/workflows")
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestRunWatchLint_InvalidPath(t *testing.T) {
	err := runWatchLint("/nonexistent/path")
	if err == nil {
		t.Error("expected error for invalid path")
	}
}
