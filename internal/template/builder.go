package template

import (
	"fmt"
	"sort"

	"github.com/lex00/wetwire-github-go/internal/discover"
	"github.com/lex00/wetwire-github-go/internal/runner"
	"github.com/lex00/wetwire-github-go/internal/serialize"
	"github.com/lex00/wetwire-github-go/workflow"
)

// Builder assembles workflow templates from discovered resources.
type Builder struct {
	// Verbose enables verbose output
	Verbose bool
}

// NewBuilder creates a new Builder.
func NewBuilder() *Builder {
	return &Builder{}
}

// BuildResult contains the result of building workflow templates.
type BuildResult struct {
	// Workflows contains the assembled workflows with YAML output
	Workflows []BuiltWorkflow

	// Errors contains any non-fatal errors encountered
	Errors []string
}

// BuiltWorkflow represents a workflow ready for output.
type BuiltWorkflow struct {
	// Name is the workflow variable name
	Name string

	// Workflow is the assembled workflow with jobs
	Workflow *workflow.Workflow

	// YAML is the serialized YAML output
	YAML []byte

	// Jobs lists the job names in dependency order
	Jobs []string
}

// Build assembles workflow templates from discovery and extraction results.
func (b *Builder) Build(discovered *discover.DiscoveryResult, extracted *runner.ExtractionResult) (*BuildResult, error) {
	result := &BuildResult{
		Workflows: []BuiltWorkflow{},
		Errors:    []string{},
	}

	// Build a map of jobs by name for quick lookup
	jobMap := make(map[string]*runner.ExtractedJob)
	for i := range extracted.Jobs {
		job := &extracted.Jobs[i]
		jobMap[job.Name] = job
	}

	// Build job dependency graph
	graph := NewGraph()
	jobDeps := make(map[string][]string)

	for _, dj := range discovered.Jobs {
		graph.AddNode(dj.Name)
		jobDeps[dj.Name] = dj.Dependencies
		for _, dep := range dj.Dependencies {
			graph.AddEdge(dj.Name, dep)
		}
	}

	// Detect cycles
	cycles := graph.DetectCycles()
	if len(cycles) > 0 {
		return nil, fmt.Errorf("dependency cycles detected: %v", cycles)
	}

	// Topologically sort jobs
	sortedJobs, err := graph.TopologicalSortKahn()
	if err != nil {
		return nil, fmt.Errorf("sorting jobs: %w", err)
	}

	// Process each workflow
	for _, dw := range discovered.Workflows {
		// Find the extracted workflow data
		var workflowData map[string]any
		for _, ew := range extracted.Workflows {
			if ew.Name == dw.Name {
				workflowData = ew.Data
				break
			}
		}

		if workflowData == nil {
			result.Errors = append(result.Errors, fmt.Sprintf("workflow %s: extraction data not found", dw.Name))
			continue
		}

		// Build the workflow
		wf, err := b.buildWorkflow(dw, workflowData, dw.Jobs, jobMap, sortedJobs)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("workflow %s: %v", dw.Name, err))
			continue
		}

		// Serialize to YAML
		yaml, err := serialize.ToYAML(wf)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("workflow %s: serialization failed: %v", dw.Name, err))
			continue
		}

		// Get ordered jobs for this workflow
		orderedJobs := b.filterAndOrderJobs(dw.Jobs, sortedJobs)

		result.Workflows = append(result.Workflows, BuiltWorkflow{
			Name:     dw.Name,
			Workflow: wf,
			YAML:     yaml,
			Jobs:     orderedJobs,
		})
	}

	return result, nil
}

// buildWorkflow assembles a workflow with its jobs.
func (b *Builder) buildWorkflow(dw discover.DiscoveredWorkflow, data map[string]any, jobNames []string, jobMap map[string]*runner.ExtractedJob, sortedJobs []string) (*workflow.Workflow, error) {
	wf := &workflow.Workflow{}

	// Set name
	if name, ok := data["Name"].(string); ok && name != "" {
		wf.Name = name
	}

	// Set triggers - handle both direct type and map reconstruction
	if on, ok := data["On"].(workflow.Triggers); ok {
		wf.On = on
	} else if onMap, ok := data["On"].(map[string]any); ok {
		wf.On = b.reconstructTriggers(onMap)
	}

	// Set env
	if env, ok := data["Env"].(map[string]any); ok {
		wf.Env = env
	}

	// Set defaults
	if defaults, ok := data["Defaults"].(*workflow.WorkflowDefaults); ok {
		wf.Defaults = defaults
	}

	// Set concurrency
	if concurrency, ok := data["Concurrency"].(*workflow.Concurrency); ok {
		wf.Concurrency = concurrency
	}

	// Set permissions
	if permissions, ok := data["Permissions"].(*workflow.Permissions); ok {
		wf.Permissions = permissions
	}

	// Add jobs in dependency order
	wf.Jobs = make(map[string]workflow.Job)
	orderedJobs := b.filterAndOrderJobs(jobNames, sortedJobs)

	for _, jobName := range orderedJobs {
		job, ok := jobMap[jobName]
		if !ok {
			continue
		}

		wfJob, err := b.buildJob(job)
		if err != nil {
			return nil, fmt.Errorf("building job %s: %w", jobName, err)
		}

		// Use the job's Name field for the YAML key, or fall back to variable name
		yamlKey := wfJob.Name
		if yamlKey == "" {
			yamlKey = jobName
		}
		wf.Jobs[yamlKey] = *wfJob
	}

	return wf, nil
}

// buildJob converts extracted job data to a workflow.Job.
func (b *Builder) buildJob(ej *runner.ExtractedJob) (*workflow.Job, error) {
	job := &workflow.Job{}

	data := ej.Data

	// Set basic fields
	if name, ok := data["Name"].(string); ok {
		job.Name = name
	}

	if runsOn, ok := data["RunsOn"]; ok {
		job.RunsOn = runsOn
	}

	if ifCond, ok := data["If"]; ok {
		job.If = ifCond
	}

	if env, ok := data["Env"].(map[string]any); ok {
		job.Env = env
	}

	if timeoutMinutes, ok := data["TimeoutMinutes"].(int); ok {
		job.TimeoutMinutes = timeoutMinutes
	}

	if continueOnError, ok := data["ContinueOnError"].(bool); ok {
		job.ContinueOnError = continueOnError
	}

	// Handle Needs - convert to strings
	if needs, ok := data["Needs"].([]any); ok {
		job.Needs = needs
	}

	// Handle outputs
	if outputs, ok := data["Outputs"].(map[string]any); ok {
		job.Outputs = outputs
	}

	// Handle strategy
	if strategy, ok := data["Strategy"].(*workflow.Strategy); ok {
		job.Strategy = strategy
	}

	// Handle permissions
	if permissions, ok := data["Permissions"].(*workflow.Permissions); ok {
		job.Permissions = permissions
	}

	// Handle defaults
	if defaults, ok := data["Defaults"].(*workflow.JobDefaults); ok {
		job.Defaults = defaults
	}

	// Handle concurrency
	if concurrency, ok := data["Concurrency"].(*workflow.Concurrency); ok {
		job.Concurrency = concurrency
	}

	// Handle container
	if container, ok := data["Container"].(*workflow.Container); ok {
		job.Container = container
	}

	// Handle services
	if services, ok := data["Services"].(map[string]workflow.Service); ok {
		job.Services = services
	}

	// Handle steps - handle both direct type and slice reconstruction
	if steps, ok := data["Steps"].([]workflow.Step); ok {
		// Convert []workflow.Step to []any
		anySteps := make([]any, len(steps))
		for i, s := range steps {
			anySteps[i] = s
		}
		job.Steps = anySteps
	} else if stepsSlice, ok := data["Steps"].([]any); ok {
		job.Steps = b.reconstructStepsAsAny(stepsSlice)
	}

	// Handle environment
	if env, ok := data["Environment"].(*workflow.Environment); ok {
		job.Environment = env
	}

	return job, nil
}

// filterAndOrderJobs returns jobs in dependency order, filtered to only include specified jobs.
func (b *Builder) filterAndOrderJobs(jobNames []string, sortedJobs []string) []string {
	// Create a set of job names
	jobSet := make(map[string]bool)
	for _, name := range jobNames {
		jobSet[name] = true
	}

	// Filter sorted jobs to only include jobs in the set
	var result []string
	for _, name := range sortedJobs {
		if jobSet[name] {
			result = append(result, name)
		}
	}

	return result
}

// OrderJobs returns jobs in dependency order using Kahn's algorithm.
func OrderJobs(jobs []discover.DiscoveredJob) ([]string, error) {
	graph := NewGraph()

	for _, job := range jobs {
		graph.AddNode(job.Name)
		for _, dep := range job.Dependencies {
			graph.AddEdge(job.Name, dep)
		}
	}

	return graph.TopologicalSortKahn()
}

// ValidateJobDependencies checks that all job dependencies are valid.
func ValidateJobDependencies(jobs []discover.DiscoveredJob) []string {
	var errors []string

	// Build set of job names
	jobNames := make(map[string]bool)
	for _, job := range jobs {
		jobNames[job.Name] = true
	}

	// Check that all dependencies exist
	for _, job := range jobs {
		for _, dep := range job.Dependencies {
			if !jobNames[dep] {
				errors = append(errors, fmt.Sprintf("job %s: unknown dependency %q", job.Name, dep))
			}
		}
	}

	// Check for cycles
	graph := NewGraph()
	for _, job := range jobs {
		graph.AddNode(job.Name)
		for _, dep := range job.Dependencies {
			if jobNames[dep] {
				graph.AddEdge(job.Name, dep)
			}
		}
	}

	cycles := graph.DetectCycles()
	for _, cycle := range cycles {
		errors = append(errors, fmt.Sprintf("dependency cycle: %v", cycle))
	}

	sort.Strings(errors)
	return errors
}

// reconstructTriggers builds a Triggers struct from a generic map.
func (b *Builder) reconstructTriggers(data map[string]any) workflow.Triggers {
	triggers := workflow.Triggers{}

	if pushData, ok := data["Push"]; ok {
		if pushMap, ok := pushData.(map[string]any); ok {
			push := &workflow.PushTrigger{}
			if branches, ok := pushMap["Branches"].([]any); ok {
				push.Branches = anySliceToStrings(branches)
			}
			if tags, ok := pushMap["Tags"].([]any); ok {
				push.Tags = anySliceToStrings(tags)
			}
			if paths, ok := pushMap["Paths"].([]any); ok {
				push.Paths = anySliceToStrings(paths)
			}
			triggers.Push = push
		}
	}

	if prData, ok := data["PullRequest"]; ok {
		if prMap, ok := prData.(map[string]any); ok {
			pr := &workflow.PullRequestTrigger{}
			if branches, ok := prMap["Branches"].([]any); ok {
				pr.Branches = anySliceToStrings(branches)
			}
			if types, ok := prMap["Types"].([]any); ok {
				pr.Types = anySliceToStrings(types)
			}
			triggers.PullRequest = pr
		}
	}

	if wdData, ok := data["WorkflowDispatch"]; ok && wdData != nil {
		triggers.WorkflowDispatch = &workflow.WorkflowDispatchTrigger{}
	}

	if wcData, ok := data["WorkflowCall"]; ok && wcData != nil {
		triggers.WorkflowCall = &workflow.WorkflowCallTrigger{}
	}

	if schedData, ok := data["Schedule"].([]any); ok {
		for _, s := range schedData {
			if sched, ok := s.(map[string]any); ok {
				if cron, ok := sched["Cron"].(string); ok {
					triggers.Schedule = append(triggers.Schedule, workflow.ScheduleTrigger{Cron: cron})
				}
			}
		}
	}

	return triggers
}

// reconstructStepsAsAny builds a slice of Steps from a generic slice.
func (b *Builder) reconstructStepsAsAny(data []any) []any {
	var steps []any

	for _, item := range data {
		stepMap, ok := item.(map[string]any)
		if !ok {
			continue
		}

		step := workflow.Step{}

		if id, ok := stepMap["ID"].(string); ok {
			step.ID = id
		}
		if name, ok := stepMap["Name"].(string); ok {
			step.Name = name
		}
		if uses, ok := stepMap["Uses"].(string); ok {
			step.Uses = uses
		}
		if run, ok := stepMap["Run"].(string); ok {
			step.Run = run
		}
		if shell, ok := stepMap["Shell"].(string); ok {
			step.Shell = shell
		}
		if ifCond, ok := stepMap["If"].(string); ok {
			step.If = ifCond
		}
		if wd, ok := stepMap["WorkingDirectory"].(string); ok {
			step.WorkingDirectory = wd
		}
		if env, ok := stepMap["Env"].(map[string]any); ok {
			step.Env = env
		}
		if with, ok := stepMap["With"].(map[string]any); ok {
			step.With = with
		}
		if timeout, ok := stepMap["TimeoutMinutes"].(float64); ok {
			step.TimeoutMinutes = int(timeout)
		}

		steps = append(steps, step)
	}

	return steps
}

// anySliceToStrings converts []any to []string.
func anySliceToStrings(slice []any) []string {
	var result []string
	for _, v := range slice {
		if s, ok := v.(string); ok {
			result = append(result, s)
		}
	}
	return result
}
