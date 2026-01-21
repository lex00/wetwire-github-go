// Package differ provides semantic comparison of GitHub Actions workflows.
package differ

import (
	"fmt"
	"os"
	"reflect"
	"sort"

	coredomain "github.com/lex00/wetwire-core-go/domain"
	"gopkg.in/yaml.v3"
)

// WorkflowDiffer implements coredomain.Differ for GitHub Actions workflows.
type WorkflowDiffer struct{}

// Compile-time check that WorkflowDiffer implements Differ.
var _ coredomain.Differ = (*WorkflowDiffer)(nil)

// New creates a new workflow differ.
func New() *WorkflowDiffer {
	return &WorkflowDiffer{}
}

// Workflow represents a GitHub Actions workflow.
type Workflow struct {
	Name        string                 `yaml:"name"`
	On          interface{}            `yaml:"on"`
	Env         map[string]string      `yaml:"env"`
	Permissions interface{}            `yaml:"permissions"`
	Concurrency interface{}            `yaml:"concurrency"`
	Jobs        map[string]Job         `yaml:"jobs"`
	raw         map[string]interface{} `yaml:"-"`
}

// Job represents a job in a workflow.
type Job struct {
	Name        string                 `yaml:"name"`
	RunsOn      interface{}            `yaml:"runs-on"`
	Needs       interface{}            `yaml:"needs"`
	If          string                 `yaml:"if"`
	Strategy    interface{}            `yaml:"strategy"`
	Env         map[string]string      `yaml:"env"`
	Environment interface{}            `yaml:"environment"`
	Permissions interface{}            `yaml:"permissions"`
	Steps       []Step                 `yaml:"steps"`
	Outputs     map[string]string      `yaml:"outputs"`
	raw         map[string]interface{} `yaml:"-"`
}

// Step represents a step in a job.
type Step struct {
	Name            string            `yaml:"name"`
	ID              string            `yaml:"id"`
	Uses            string            `yaml:"uses"`
	Run             string            `yaml:"run"`
	With            map[string]string `yaml:"with"`
	Env             map[string]string `yaml:"env"`
	If              string            `yaml:"if"`
	ContinueOnError bool              `yaml:"continue-on-error"`
	WorkingDir      string            `yaml:"working-directory"`
}

// Diff compares two GitHub Actions workflow files and returns differences.
func (d *WorkflowDiffer) Diff(ctx *coredomain.Context, file1, file2 string, opts coredomain.DiffOpts) (*coredomain.DiffResult, error) {
	wf1, err := loadWorkflow(file1)
	if err != nil {
		return nil, fmt.Errorf("failed to load %s: %w", file1, err)
	}

	wf2, err := loadWorkflow(file2)
	if err != nil {
		return nil, fmt.Errorf("failed to load %s: %w", file2, err)
	}

	return compareWorkflows(wf1, wf2, opts)
}

// loadWorkflow loads a workflow from a YAML file.
func loadWorkflow(path string) (*Workflow, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var wf Workflow
	if err := yaml.Unmarshal(data, &wf); err != nil {
		return nil, err
	}

	// Also parse raw for detailed comparison
	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	wf.raw = raw

	return &wf, nil
}

// compareWorkflows compares two workflows and returns differences.
func compareWorkflows(wf1, wf2 *Workflow, opts coredomain.DiffOpts) (*coredomain.DiffResult, error) {
	result := &coredomain.DiffResult{}

	// Compare workflow-level properties
	workflowChanges := compareWorkflowProperties(wf1, wf2, opts)
	if len(workflowChanges) > 0 {
		result.Entries = append(result.Entries, coredomain.DiffEntry{
			Resource: wf1.Name,
			Type:     "workflow",
			Action:   "modified",
			Changes:  workflowChanges,
		})
	}

	// Compare jobs
	jobs1 := wf1.Jobs
	jobs2 := wf2.Jobs

	// Find added jobs
	for name, job := range jobs2 {
		if _, exists := jobs1[name]; !exists {
			result.Entries = append(result.Entries, coredomain.DiffEntry{
				Resource: fmt.Sprintf("job:%s", name),
				Type:     "job",
				Action:   "added",
				Changes:  []string{fmt.Sprintf("runs-on: %v", job.RunsOn)},
			})
		}
	}

	// Find removed jobs
	for name := range jobs1 {
		if _, exists := jobs2[name]; !exists {
			result.Entries = append(result.Entries, coredomain.DiffEntry{
				Resource: fmt.Sprintf("job:%s", name),
				Type:     "job",
				Action:   "removed",
			})
		}
	}

	// Find modified jobs
	for name, job1 := range jobs1 {
		if job2, exists := jobs2[name]; exists {
			changes := compareJobs(job1, job2, opts)
			if len(changes) > 0 {
				result.Entries = append(result.Entries, coredomain.DiffEntry{
					Resource: fmt.Sprintf("job:%s", name),
					Type:     "job",
					Action:   "modified",
					Changes:  changes,
				})
			}
		}
	}

	// Sort entries for consistent output
	sort.Slice(result.Entries, func(i, j int) bool {
		if result.Entries[i].Action != result.Entries[j].Action {
			order := map[string]int{"added": 0, "modified": 1, "removed": 2}
			return order[result.Entries[i].Action] < order[result.Entries[j].Action]
		}
		return result.Entries[i].Resource < result.Entries[j].Resource
	})

	// Calculate summary
	for _, e := range result.Entries {
		switch e.Action {
		case "added":
			result.Summary.Added++
		case "removed":
			result.Summary.Removed++
		case "modified":
			result.Summary.Modified++
		}
	}
	result.Summary.Total = result.Summary.Added + result.Summary.Removed + result.Summary.Modified

	return result, nil
}

// compareWorkflowProperties compares workflow-level properties.
func compareWorkflowProperties(wf1, wf2 *Workflow, opts coredomain.DiffOpts) []string {
	var changes []string

	if wf1.Name != wf2.Name {
		changes = append(changes, fmt.Sprintf("name: %q → %q", wf1.Name, wf2.Name))
	}

	if !deepEqual(wf1.On, wf2.On, opts) {
		changes = append(changes, "triggers changed")
	}

	if !deepEqual(wf1.Permissions, wf2.Permissions, opts) {
		changes = append(changes, "permissions changed")
	}

	if !deepEqual(wf1.Env, wf2.Env, opts) {
		changes = append(changes, "env changed")
	}

	if !deepEqual(wf1.Concurrency, wf2.Concurrency, opts) {
		changes = append(changes, "concurrency changed")
	}

	return changes
}

// compareJobs compares two jobs and returns changes.
func compareJobs(j1, j2 Job, opts coredomain.DiffOpts) []string {
	var changes []string

	if j1.Name != j2.Name {
		changes = append(changes, fmt.Sprintf("name: %q → %q", j1.Name, j2.Name))
	}

	if !deepEqual(j1.RunsOn, j2.RunsOn, opts) {
		changes = append(changes, fmt.Sprintf("runs-on: %v → %v", j1.RunsOn, j2.RunsOn))
	}

	if !deepEqual(j1.Needs, j2.Needs, opts) {
		changes = append(changes, "needs changed")
	}

	if j1.If != j2.If {
		changes = append(changes, fmt.Sprintf("if: %q → %q", j1.If, j2.If))
	}

	if !deepEqual(j1.Strategy, j2.Strategy, opts) {
		changes = append(changes, "strategy changed")
	}

	if !deepEqual(j1.Env, j2.Env, opts) {
		changes = append(changes, "env changed")
	}

	if !deepEqual(j1.Environment, j2.Environment, opts) {
		changes = append(changes, "environment changed")
	}

	if !deepEqual(j1.Permissions, j2.Permissions, opts) {
		changes = append(changes, "permissions changed")
	}

	// Compare steps
	stepChanges := compareSteps(j1.Steps, j2.Steps, opts)
	changes = append(changes, stepChanges...)

	if !deepEqual(j1.Outputs, j2.Outputs, opts) {
		changes = append(changes, "outputs changed")
	}

	return changes
}

// compareSteps compares two lists of steps.
func compareSteps(steps1, steps2 []Step, opts coredomain.DiffOpts) []string {
	var changes []string

	if len(steps1) != len(steps2) {
		changes = append(changes, fmt.Sprintf("steps count: %d → %d", len(steps1), len(steps2)))
	}

	// Compare step by step up to the shorter length
	minLen := len(steps1)
	if len(steps2) < minLen {
		minLen = len(steps2)
	}

	for i := 0; i < minLen; i++ {
		s1, s2 := steps1[i], steps2[i]
		stepChanges := compareStep(s1, s2, i, opts)
		changes = append(changes, stepChanges...)
	}

	// Report added steps
	for i := minLen; i < len(steps2); i++ {
		s := steps2[i]
		identifier := stepIdentifier(s)
		changes = append(changes, fmt.Sprintf("step[%d]: %s added", i, identifier))
	}

	// Report removed steps
	for i := minLen; i < len(steps1); i++ {
		s := steps1[i]
		identifier := stepIdentifier(s)
		changes = append(changes, fmt.Sprintf("step[%d]: %s removed", i, identifier))
	}

	return changes
}

// compareStep compares two steps at the same position.
func compareStep(s1, s2 Step, index int, opts coredomain.DiffOpts) []string {
	var changes []string

	prefix := fmt.Sprintf("step[%d]", index)

	if s1.Name != s2.Name {
		changes = append(changes, fmt.Sprintf("%s.name: %q → %q", prefix, s1.Name, s2.Name))
	}

	if s1.ID != s2.ID {
		changes = append(changes, fmt.Sprintf("%s.id: %q → %q", prefix, s1.ID, s2.ID))
	}

	if s1.Uses != s2.Uses {
		changes = append(changes, fmt.Sprintf("%s.uses: %q → %q", prefix, s1.Uses, s2.Uses))
	}

	if s1.Run != s2.Run {
		changes = append(changes, fmt.Sprintf("%s.run changed", prefix))
	}

	if !deepEqual(s1.With, s2.With, opts) {
		changes = append(changes, fmt.Sprintf("%s.with changed", prefix))
	}

	if !deepEqual(s1.Env, s2.Env, opts) {
		changes = append(changes, fmt.Sprintf("%s.env changed", prefix))
	}

	if s1.If != s2.If {
		changes = append(changes, fmt.Sprintf("%s.if: %q → %q", prefix, s1.If, s2.If))
	}

	if s1.ContinueOnError != s2.ContinueOnError {
		changes = append(changes, fmt.Sprintf("%s.continue-on-error: %v → %v", prefix, s1.ContinueOnError, s2.ContinueOnError))
	}

	if s1.WorkingDir != s2.WorkingDir {
		changes = append(changes, fmt.Sprintf("%s.working-directory: %q → %q", prefix, s1.WorkingDir, s2.WorkingDir))
	}

	return changes
}

// stepIdentifier returns a human-readable identifier for a step.
func stepIdentifier(s Step) string {
	if s.Name != "" {
		return fmt.Sprintf("%q", s.Name)
	}
	if s.Uses != "" {
		return s.Uses
	}
	if s.Run != "" {
		// Truncate long run commands
		run := s.Run
		if len(run) > 30 {
			run = run[:27] + "..."
		}
		return fmt.Sprintf("run:%q", run)
	}
	return "unnamed"
}

// deepEqual compares two values with optional order ignoring.
func deepEqual(a, b interface{}, opts coredomain.DiffOpts) bool {
	if opts.IgnoreOrder {
		a = normalizeValue(a)
		b = normalizeValue(b)
	}
	return reflect.DeepEqual(a, b)
}

// normalizeValue normalizes a value for comparison.
func normalizeValue(v interface{}) interface{} {
	switch val := v.(type) {
	case []interface{}:
		result := make([]interface{}, len(val))
		copy(result, val)
		return result
	case map[string]interface{}:
		result := make(map[string]interface{})
		for k, v := range val {
			result[k] = normalizeValue(v)
		}
		return result
	default:
		return v
	}
}
