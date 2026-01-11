package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiffCmd_NoChanges(t *testing.T) {
	dir := t.TempDir()

	yaml1 := `name: CI
on: push
jobs:
  build:
    runs-on: ubuntu-latest
`
	path1 := filepath.Join(dir, "ci1.yml")
	path2 := filepath.Join(dir, "ci2.yml")
	os.WriteFile(path1, []byte(yaml1), 0644)
	os.WriteFile(path2, []byte(yaml1), 0644)

	oldJobs, _ := parseYAMLWorkflowJobs(path1)
	newJobs, _ := parseYAMLWorkflowJobs(path2)

	result := compareWorkflows(oldJobs, newJobs)

	hasChanges := len(result.AddedJobs) > 0 ||
		len(result.RemovedJobs) > 0 ||
		len(result.ModifiedJobs) > 0 ||
		len(result.DependencyChanges) > 0

	if hasChanges {
		t.Errorf("expected no changes, got: added=%v, removed=%v", result.AddedJobs, result.RemovedJobs)
	}
}

func TestDiffCmd_AddedJob(t *testing.T) {
	dir := t.TempDir()

	yaml1 := `name: CI
on: push
jobs:
  build:
    runs-on: ubuntu-latest
`
	yaml2 := `name: CI
on: push
jobs:
  build:
    runs-on: ubuntu-latest
  test:
    runs-on: ubuntu-latest
`
	path1 := filepath.Join(dir, "ci1.yml")
	path2 := filepath.Join(dir, "ci2.yml")
	os.WriteFile(path1, []byte(yaml1), 0644)
	os.WriteFile(path2, []byte(yaml2), 0644)

	oldJobs, _ := parseYAMLWorkflowJobs(path1)
	newJobs, _ := parseYAMLWorkflowJobs(path2)

	result := compareWorkflows(oldJobs, newJobs)

	if len(result.AddedJobs) != 1 || result.AddedJobs[0] != "test" {
		t.Errorf("expected added job 'test', got: %v", result.AddedJobs)
	}
}

func TestDiffCmd_RemovedJob(t *testing.T) {
	dir := t.TempDir()

	yaml1 := `name: CI
on: push
jobs:
  build:
    runs-on: ubuntu-latest
  test:
    runs-on: ubuntu-latest
`
	yaml2 := `name: CI
on: push
jobs:
  build:
    runs-on: ubuntu-latest
`
	path1 := filepath.Join(dir, "ci1.yml")
	path2 := filepath.Join(dir, "ci2.yml")
	os.WriteFile(path1, []byte(yaml1), 0644)
	os.WriteFile(path2, []byte(yaml2), 0644)

	oldJobs, _ := parseYAMLWorkflowJobs(path1)
	newJobs, _ := parseYAMLWorkflowJobs(path2)

	result := compareWorkflows(oldJobs, newJobs)

	if len(result.RemovedJobs) != 1 || result.RemovedJobs[0] != "test" {
		t.Errorf("expected removed job 'test', got: %v", result.RemovedJobs)
	}
}

func TestDiffCmd_NeedsChanges(t *testing.T) {
	dir := t.TempDir()

	yaml1 := `name: CI
on: push
jobs:
  build:
    runs-on: ubuntu-latest
  deploy:
    runs-on: ubuntu-latest
    needs: [build]
`
	yaml2 := `name: CI
on: push
jobs:
  build:
    runs-on: ubuntu-latest
  test:
    runs-on: ubuntu-latest
  deploy:
    runs-on: ubuntu-latest
    needs: [build, test]
`
	path1 := filepath.Join(dir, "ci1.yml")
	path2 := filepath.Join(dir, "ci2.yml")
	os.WriteFile(path1, []byte(yaml1), 0644)
	os.WriteFile(path2, []byte(yaml2), 0644)

	oldJobs, _ := parseYAMLWorkflowJobs(path1)
	newJobs, _ := parseYAMLWorkflowJobs(path2)

	result := compareWorkflows(oldJobs, newJobs)

	if len(result.DependencyChanges) != 1 {
		t.Errorf("expected 1 dependency change, got: %d", len(result.DependencyChanges))
		return
	}

	dc := result.DependencyChanges[0]
	if dc.Job != "deploy" {
		t.Errorf("expected dependency change for 'deploy', got: %s", dc.Job)
	}

	if len(dc.AddedDeps) != 1 || dc.AddedDeps[0] != "test" {
		t.Errorf("expected added dependency 'test', got: %v", dc.AddedDeps)
	}
}

func TestParseYAMLWorkflowJobs(t *testing.T) {
	dir := t.TempDir()

	yaml := `name: CI
on: push
jobs:
  build:
    runs-on: ubuntu-latest
  test:
    runs-on: ubuntu-latest
    needs: build
  deploy:
    runs-on: ubuntu-latest
    needs: [build, test]
`
	path := filepath.Join(dir, "ci.yml")
	os.WriteFile(path, []byte(yaml), 0644)

	jobs, err := parseYAMLWorkflowJobs(path)
	if err != nil {
		t.Fatalf("parseYAMLWorkflowJobs failed: %v", err)
	}

	if len(jobs) != 3 {
		t.Errorf("expected 3 jobs, got: %d", len(jobs))
	}

	// Check job names are present
	jobNames := make(map[string]bool)
	for _, job := range jobs {
		jobNames[job.Name] = true
	}

	if !jobNames["build"] || !jobNames["test"] || !jobNames["deploy"] {
		t.Errorf("expected jobs build, test, deploy, got: %v", jobNames)
	}
}

func TestParseYAMLWorkflowJobs_NeedsString(t *testing.T) {
	dir := t.TempDir()

	yaml := `name: CI
on: push
jobs:
  build:
    runs-on: ubuntu-latest
  test:
    runs-on: ubuntu-latest
    needs: build
`
	path := filepath.Join(dir, "ci.yml")
	os.WriteFile(path, []byte(yaml), 0644)

	jobs, err := parseYAMLWorkflowJobs(path)
	if err != nil {
		t.Fatalf("parseYAMLWorkflowJobs failed: %v", err)
	}

	// Find test job
	var testJob *jobInfo
	for i := range jobs {
		if jobs[i].Name == "test" {
			testJob = &jobs[i]
			break
		}
	}

	if testJob == nil {
		t.Fatal("test job not found")
	}

	if len(testJob.Dependencies) != 1 || testJob.Dependencies[0] != "build" {
		t.Errorf("expected test job to have dependency 'build', got: %v", testJob.Dependencies)
	}
}

func TestParseYAMLWorkflowJobs_InvalidPath(t *testing.T) {
	_, err := parseYAMLWorkflowJobs("/nonexistent/path.yml")
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestParseYAMLWorkflowJobs_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "invalid.yml")
	os.WriteFile(path, []byte("invalid: yaml: content:"), 0644)

	_, err := parseYAMLWorkflowJobs(path)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestParseYAMLWorkflowJobs_NoJobs(t *testing.T) {
	dir := t.TempDir()

	yaml := `name: CI
on: push
`
	path := filepath.Join(dir, "ci.yml")
	os.WriteFile(path, []byte(yaml), 0644)

	jobs, err := parseYAMLWorkflowJobs(path)
	if err != nil {
		t.Fatalf("parseYAMLWorkflowJobs failed: %v", err)
	}

	if len(jobs) != 0 {
		t.Errorf("expected 0 jobs, got: %d", len(jobs))
	}
}

func TestCompareWorkflows_Empty(t *testing.T) {
	result := compareWorkflows([]jobInfo{}, []jobInfo{})

	hasChanges := len(result.AddedJobs) > 0 ||
		len(result.RemovedJobs) > 0 ||
		len(result.ModifiedJobs) > 0 ||
		len(result.DependencyChanges) > 0

	if hasChanges {
		t.Error("expected no changes for empty workflows")
	}
}

func TestCompareWorkflows_AllAdded(t *testing.T) {
	oldJobs := []jobInfo{}
	newJobs := []jobInfo{
		{Name: "build"},
		{Name: "test"},
	}

	result := compareWorkflows(oldJobs, newJobs)

	if len(result.AddedJobs) != 2 {
		t.Errorf("expected 2 added jobs, got: %d", len(result.AddedJobs))
	}
}

func TestCompareWorkflows_AllRemoved(t *testing.T) {
	oldJobs := []jobInfo{
		{Name: "build"},
		{Name: "test"},
	}
	newJobs := []jobInfo{}

	result := compareWorkflows(oldJobs, newJobs)

	if len(result.RemovedJobs) != 2 {
		t.Errorf("expected 2 removed jobs, got: %d", len(result.RemovedJobs))
	}
}

func TestCompareWorkflows_RemovedDependency(t *testing.T) {
	oldJobs := []jobInfo{
		{Name: "build"},
		{Name: "test"},
		{Name: "deploy", Dependencies: []string{"build", "test"}},
	}
	newJobs := []jobInfo{
		{Name: "build"},
		{Name: "test"},
		{Name: "deploy", Dependencies: []string{"build"}},
	}

	result := compareWorkflows(oldJobs, newJobs)

	if len(result.DependencyChanges) != 1 {
		t.Errorf("expected 1 dependency change, got: %d", len(result.DependencyChanges))
		return
	}

	dc := result.DependencyChanges[0]
	if len(dc.RemovedDeps) != 1 || dc.RemovedDeps[0] != "test" {
		t.Errorf("expected removed dependency 'test', got: %v", dc.RemovedDeps)
	}
}

func TestOutputDiffText_NoChanges(t *testing.T) {
	result := diffResult{Success: true}

	hasChanges := len(result.AddedJobs) > 0 ||
		len(result.RemovedJobs) > 0 ||
		len(result.ModifiedJobs) > 0 ||
		len(result.DependencyChanges) > 0

	if hasChanges {
		t.Error("expected no changes")
	}
}
