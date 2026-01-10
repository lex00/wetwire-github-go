// Package serialize provides YAML serialization for workflow types.
package serialize

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/lex00/wetwire-github-go/workflow"
)

// ToYAML serializes a workflow to YAML bytes.
func ToYAML(w *workflow.Workflow) ([]byte, error) {
	// Convert to a map for YAML serialization
	m, err := workflowToMap(w)
	if err != nil {
		return nil, fmt.Errorf("converting workflow to map: %w", err)
	}

	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(m); err != nil {
		return nil, fmt.Errorf("encoding YAML: %w", err)
	}
	if err := encoder.Close(); err != nil {
		return nil, fmt.Errorf("closing encoder: %w", err)
	}

	return buf.Bytes(), nil
}

// workflowToMap converts a Workflow to a map for YAML serialization.
func workflowToMap(w *workflow.Workflow) (map[string]any, error) {
	m := make(map[string]any)

	if w.Name != "" {
		m["name"] = w.Name
	}

	// Serialize triggers
	on, err := triggersToMap(&w.On)
	if err != nil {
		return nil, err
	}
	if len(on) > 0 {
		m["on"] = on
	}

	if len(w.Env) > 0 {
		m["env"] = serializeEnv(w.Env)
	}

	if w.Defaults != nil {
		m["defaults"] = structToMap(w.Defaults)
	}

	if w.Concurrency != nil {
		m["concurrency"] = structToMap(w.Concurrency)
	}

	if w.Permissions != nil {
		m["permissions"] = permissionsToMap(w.Permissions)
	}

	if len(w.Jobs) > 0 {
		jobs := make(map[string]any)
		for name, job := range w.Jobs {
			jobMap, err := jobToMap(&job)
			if err != nil {
				return nil, fmt.Errorf("serializing job %s: %w", name, err)
			}
			jobs[name] = jobMap
		}
		m["jobs"] = jobs
	}

	return m, nil
}

// jobToMap converts a Job to a map for YAML serialization.
func jobToMap(j *workflow.Job) (map[string]any, error) {
	m := make(map[string]any)

	if j.Name != "" {
		m["name"] = j.Name
	}

	if j.RunsOn != nil {
		m["runs-on"] = serializeValue(j.RunsOn)
	}

	if len(j.Needs) > 0 {
		m["needs"] = serializeNeeds(j.Needs)
	}

	if j.If != nil {
		m["if"] = serializeCondition(j.If)
	}

	if j.Permissions != nil {
		m["permissions"] = permissionsToMap(j.Permissions)
	}

	if j.Environment != nil {
		m["environment"] = structToMap(j.Environment)
	}

	if j.Concurrency != nil {
		m["concurrency"] = structToMap(j.Concurrency)
	}

	if len(j.Outputs) > 0 {
		m["outputs"] = serializeEnv(j.Outputs)
	}

	if len(j.Env) > 0 {
		m["env"] = serializeEnv(j.Env)
	}

	if j.Defaults != nil {
		m["defaults"] = structToMap(j.Defaults)
	}

	if j.Strategy != nil {
		strategy, err := strategyToMap(j.Strategy)
		if err != nil {
			return nil, err
		}
		m["strategy"] = strategy
	}

	if j.Container != nil {
		m["container"] = structToMap(j.Container)
	}

	if len(j.Services) > 0 {
		services := make(map[string]any)
		for name, svc := range j.Services {
			services[name] = structToMap(&svc)
		}
		m["services"] = services
	}

	if len(j.Steps) > 0 {
		steps := make([]any, len(j.Steps))
		for i, step := range j.Steps {
			stepMap, err := anyStepToMap(step)
			if err != nil {
				return nil, fmt.Errorf("serializing step %d: %w", i, err)
			}
			steps[i] = stepMap
		}
		m["steps"] = steps
	}

	if j.TimeoutMinutes > 0 {
		m["timeout-minutes"] = j.TimeoutMinutes
	}

	if j.ContinueOnError {
		m["continue-on-error"] = j.ContinueOnError
	}

	return m, nil
}

// anyStepToMap handles both workflow.Step and StepAction types.
func anyStepToMap(s any) (map[string]any, error) {
	switch v := s.(type) {
	case workflow.Step:
		return stepToMap(&v)
	case *workflow.Step:
		return stepToMap(v)
	case workflow.StepAction:
		// Convert action wrapper to Step and serialize
		step := workflow.ToStep(v)
		return stepToMap(&step)
	default:
		return nil, fmt.Errorf("unsupported step type: %T", s)
	}
}

// stepToMap converts a Step to a map for YAML serialization.
func stepToMap(s *workflow.Step) (map[string]any, error) {
	m := make(map[string]any)

	if s.ID != "" {
		m["id"] = s.ID
	}

	if s.Name != "" {
		m["name"] = s.Name
	}

	if s.If != nil {
		m["if"] = serializeCondition(s.If)
	}

	if s.Uses != "" {
		m["uses"] = s.Uses
	}

	if len(s.With) > 0 {
		m["with"] = serializeEnv(s.With)
	}

	if s.Run != "" {
		m["run"] = s.Run
	}

	if s.Shell != "" {
		m["shell"] = s.Shell
	}

	if len(s.Env) > 0 {
		m["env"] = serializeEnv(s.Env)
	}

	if s.WorkingDirectory != "" {
		m["working-directory"] = s.WorkingDirectory
	}

	if s.ContinueOnError {
		m["continue-on-error"] = s.ContinueOnError
	}

	if s.TimeoutMinutes > 0 {
		m["timeout-minutes"] = s.TimeoutMinutes
	}

	return m, nil
}

// strategyToMap converts a Strategy to a map for YAML serialization.
func strategyToMap(s *workflow.Strategy) (map[string]any, error) {
	m := make(map[string]any)

	if s.Matrix != nil {
		matrix := make(map[string]any)
		for k, v := range s.Matrix.Values {
			matrix[k] = v
		}
		if len(s.Matrix.Include) > 0 {
			matrix["include"] = s.Matrix.Include
		}
		if len(s.Matrix.Exclude) > 0 {
			matrix["exclude"] = s.Matrix.Exclude
		}
		m["matrix"] = matrix
	}

	if s.FailFast != nil {
		m["fail-fast"] = *s.FailFast
	}

	if s.MaxParallel > 0 {
		m["max-parallel"] = s.MaxParallel
	}

	return m, nil
}

// triggersToMap converts Triggers to a map for YAML serialization.
func triggersToMap(t *workflow.Triggers) (map[string]any, error) {
	m := make(map[string]any)

	if t.Push != nil {
		m["push"] = pushTriggerToMap(t.Push)
	}
	if t.PullRequest != nil {
		m["pull_request"] = pullRequestTriggerToMap(t.PullRequest)
	}
	if t.PullRequestTarget != nil {
		m["pull_request_target"] = pullRequestTargetTriggerToMap(t.PullRequestTarget)
	}
	if len(t.Schedule) > 0 {
		schedules := make([]map[string]any, len(t.Schedule))
		for i, s := range t.Schedule {
			schedules[i] = map[string]any{"cron": s.Cron}
		}
		m["schedule"] = schedules
	}
	if t.WorkflowDispatch != nil {
		m["workflow_dispatch"] = workflowDispatchToMap(t.WorkflowDispatch)
	}
	if t.WorkflowCall != nil {
		m["workflow_call"] = workflowCallToMap(t.WorkflowCall)
	}
	if t.WorkflowRun != nil {
		m["workflow_run"] = workflowRunToMap(t.WorkflowRun)
	}
	if t.RepositoryDispatch != nil {
		m["repository_dispatch"] = repositoryDispatchToMap(t.RepositoryDispatch)
	}

	// Simple triggers (empty structs)
	if t.Create != nil {
		m["create"] = nil
	}
	if t.Delete != nil {
		m["delete"] = nil
	}
	if t.Fork != nil {
		m["fork"] = nil
	}
	if t.Gollum != nil {
		m["gollum"] = nil
	}
	if t.Public != nil {
		m["public"] = nil
	}
	if t.PageBuild != nil {
		m["page_build"] = nil
	}
	if t.Status != nil {
		m["status"] = nil
	}

	// Triggers with types
	if t.IssueComment != nil {
		m["issue_comment"] = typesToMap(t.IssueComment.Types)
	}
	if t.Issues != nil {
		m["issues"] = typesToMap(t.Issues.Types)
	}
	if t.Label != nil {
		m["label"] = typesToMap(t.Label.Types)
	}
	if t.Milestone != nil {
		m["milestone"] = typesToMap(t.Milestone.Types)
	}
	if t.Project != nil {
		m["project"] = typesToMap(t.Project.Types)
	}
	if t.ProjectCard != nil {
		m["project_card"] = typesToMap(t.ProjectCard.Types)
	}
	if t.ProjectColumn != nil {
		m["project_column"] = typesToMap(t.ProjectColumn.Types)
	}
	if t.PullRequestReview != nil {
		m["pull_request_review"] = typesToMap(t.PullRequestReview.Types)
	}
	if t.PullRequestReviewComment != nil {
		m["pull_request_review_comment"] = typesToMap(t.PullRequestReviewComment.Types)
	}
	if t.Release != nil {
		m["release"] = typesToMap(t.Release.Types)
	}
	if t.Watch != nil {
		m["watch"] = typesToMap(t.Watch.Types)
	}
	if t.CheckRun != nil {
		m["check_run"] = typesToMap(t.CheckRun.Types)
	}
	if t.CheckSuite != nil {
		m["check_suite"] = typesToMap(t.CheckSuite.Types)
	}
	if t.Discussion != nil {
		m["discussion"] = typesToMap(t.Discussion.Types)
	}
	if t.DiscussionComment != nil {
		m["discussion_comment"] = typesToMap(t.DiscussionComment.Types)
	}
	if t.MergeGroup != nil {
		m["merge_group"] = typesToMap(t.MergeGroup.Types)
	}

	return m, nil
}

func pushTriggerToMap(t *workflow.PushTrigger) map[string]any {
	m := make(map[string]any)
	if len(t.Branches) > 0 {
		m["branches"] = t.Branches
	}
	if len(t.BranchesIgnore) > 0 {
		m["branches-ignore"] = t.BranchesIgnore
	}
	if len(t.Tags) > 0 {
		m["tags"] = t.Tags
	}
	if len(t.TagsIgnore) > 0 {
		m["tags-ignore"] = t.TagsIgnore
	}
	if len(t.Paths) > 0 {
		m["paths"] = t.Paths
	}
	if len(t.PathsIgnore) > 0 {
		m["paths-ignore"] = t.PathsIgnore
	}
	if len(m) == 0 {
		return nil
	}
	return m
}

func pullRequestTriggerToMap(t *workflow.PullRequestTrigger) map[string]any {
	m := make(map[string]any)
	if len(t.Types) > 0 {
		m["types"] = t.Types
	}
	if len(t.Branches) > 0 {
		m["branches"] = t.Branches
	}
	if len(t.BranchesIgnore) > 0 {
		m["branches-ignore"] = t.BranchesIgnore
	}
	if len(t.Paths) > 0 {
		m["paths"] = t.Paths
	}
	if len(t.PathsIgnore) > 0 {
		m["paths-ignore"] = t.PathsIgnore
	}
	if len(m) == 0 {
		return nil
	}
	return m
}

func pullRequestTargetTriggerToMap(t *workflow.PullRequestTargetTrigger) map[string]any {
	m := make(map[string]any)
	if len(t.Types) > 0 {
		m["types"] = t.Types
	}
	if len(t.Branches) > 0 {
		m["branches"] = t.Branches
	}
	if len(t.BranchesIgnore) > 0 {
		m["branches-ignore"] = t.BranchesIgnore
	}
	if len(t.Paths) > 0 {
		m["paths"] = t.Paths
	}
	if len(t.PathsIgnore) > 0 {
		m["paths-ignore"] = t.PathsIgnore
	}
	if len(m) == 0 {
		return nil
	}
	return m
}

func workflowDispatchToMap(t *workflow.WorkflowDispatchTrigger) map[string]any {
	if len(t.Inputs) == 0 {
		return nil
	}
	m := make(map[string]any)
	inputs := make(map[string]any)
	for name, input := range t.Inputs {
		inputMap := make(map[string]any)
		if input.Description != "" {
			inputMap["description"] = input.Description
		}
		if input.Required {
			inputMap["required"] = input.Required
		}
		if input.Default != nil {
			inputMap["default"] = input.Default
		}
		if input.Type != "" {
			inputMap["type"] = input.Type
		}
		if len(input.Options) > 0 {
			inputMap["options"] = input.Options
		}
		inputs[name] = inputMap
	}
	m["inputs"] = inputs
	return m
}

func workflowCallToMap(t *workflow.WorkflowCallTrigger) map[string]any {
	m := make(map[string]any)

	if len(t.Inputs) > 0 {
		inputs := make(map[string]any)
		for name, input := range t.Inputs {
			inputMap := make(map[string]any)
			if input.Description != "" {
				inputMap["description"] = input.Description
			}
			if input.Required {
				inputMap["required"] = input.Required
			}
			if input.Default != nil {
				inputMap["default"] = input.Default
			}
			if input.Type != "" {
				inputMap["type"] = input.Type
			}
			inputs[name] = inputMap
		}
		m["inputs"] = inputs
	}

	if len(t.Outputs) > 0 {
		outputs := make(map[string]any)
		for name, output := range t.Outputs {
			outputMap := make(map[string]any)
			if output.Description != "" {
				outputMap["description"] = output.Description
			}
			outputMap["value"] = output.Value.String()
			outputs[name] = outputMap
		}
		m["outputs"] = outputs
	}

	if len(t.Secrets) > 0 {
		secrets := make(map[string]any)
		for name, secret := range t.Secrets {
			secretMap := make(map[string]any)
			if secret.Description != "" {
				secretMap["description"] = secret.Description
			}
			if secret.Required {
				secretMap["required"] = secret.Required
			}
			secrets[name] = secretMap
		}
		m["secrets"] = secrets
	}

	if len(m) == 0 {
		return nil
	}
	return m
}

func workflowRunToMap(t *workflow.WorkflowRunTrigger) map[string]any {
	m := make(map[string]any)
	if len(t.Workflows) > 0 {
		m["workflows"] = t.Workflows
	}
	if len(t.Types) > 0 {
		m["types"] = t.Types
	}
	if len(t.Branches) > 0 {
		m["branches"] = t.Branches
	}
	if len(m) == 0 {
		return nil
	}
	return m
}

func repositoryDispatchToMap(t *workflow.RepositoryDispatchTrigger) map[string]any {
	if len(t.Types) == 0 {
		return nil
	}
	return map[string]any{"types": t.Types}
}

func typesToMap(types []string) map[string]any {
	if len(types) == 0 {
		return nil
	}
	return map[string]any{"types": types}
}

func permissionsToMap(p *workflow.Permissions) map[string]any {
	m := make(map[string]any)
	if p.Actions != "" {
		m["actions"] = p.Actions
	}
	if p.Checks != "" {
		m["checks"] = p.Checks
	}
	if p.Contents != "" {
		m["contents"] = p.Contents
	}
	if p.Deployments != "" {
		m["deployments"] = p.Deployments
	}
	if p.Discussions != "" {
		m["discussions"] = p.Discussions
	}
	if p.IDToken != "" {
		m["id-token"] = p.IDToken
	}
	if p.Issues != "" {
		m["issues"] = p.Issues
	}
	if p.Packages != "" {
		m["packages"] = p.Packages
	}
	if p.Pages != "" {
		m["pages"] = p.Pages
	}
	if p.PullRequests != "" {
		m["pull-requests"] = p.PullRequests
	}
	if p.RepositoryProjects != "" {
		m["repository-projects"] = p.RepositoryProjects
	}
	if p.SecurityEvents != "" {
		m["security-events"] = p.SecurityEvents
	}
	if p.Statuses != "" {
		m["statuses"] = p.Statuses
	}
	return m
}

// Helper functions

// serializeCondition converts a condition to a string.
func serializeCondition(c any) string {
	switch v := c.(type) {
	case workflow.Expression:
		return v.Raw()
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", c)
	}
}

// serializeNeeds converts job references to their names.
func serializeNeeds(needs []any) []string {
	result := make([]string, len(needs))
	for i, n := range needs {
		switch v := n.(type) {
		case string:
			result[i] = v
		case workflow.Job:
			result[i] = v.Name
		default:
			// Try to get Name field via reflection
			rv := reflect.ValueOf(n)
			if rv.Kind() == reflect.Struct {
				nameField := rv.FieldByName("Name")
				if nameField.IsValid() && nameField.Kind() == reflect.String {
					result[i] = nameField.String()
				} else {
					result[i] = fmt.Sprintf("%v", n)
				}
			} else {
				result[i] = fmt.Sprintf("%v", n)
			}
		}
	}
	return result
}

// serializeValue converts a value to YAML-safe format.
func serializeValue(v any) any {
	switch val := v.(type) {
	case workflow.Expression:
		return val.String()
	default:
		return v
	}
}

// serializeEnv converts an env map, handling Expression values.
func serializeEnv(env map[string]any) map[string]any {
	result := make(map[string]any)
	for k, v := range env {
		switch val := v.(type) {
		case workflow.Expression:
			result[k] = val.String()
		default:
			result[k] = v
		}
	}
	return result
}

// structToMap converts a struct to a map using YAML tags.
func structToMap(v any) map[string]any {
	if v == nil {
		return nil
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil
		}
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return nil
	}

	result := make(map[string]any)
	rt := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		field := rt.Field(i)
		value := rv.Field(i)

		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}

		// Get YAML tag
		tag := field.Tag.Get("yaml")
		if tag == "-" {
			continue
		}

		name := field.Name
		omitempty := false

		if tag != "" {
			parts := strings.Split(tag, ",")
			if parts[0] != "" {
				name = parts[0]
			}
			for _, p := range parts[1:] {
				if p == "omitempty" {
					omitempty = true
				}
			}
		} else {
			// Convert to kebab-case if no tag
			name = toKebabCase(name)
		}

		// Skip zero values if omitempty
		if omitempty && isZeroValue(value) {
			continue
		}

		result[name] = value.Interface()
	}

	return result
}

// isZeroValue checks if a reflect.Value is its zero value.
func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return v.IsNil()
	default:
		return v.IsZero()
	}
}

// toKebabCase converts PascalCase to kebab-case.
func toKebabCase(s string) string {
	re := regexp.MustCompile(`([a-z0-9])([A-Z])`)
	kebab := re.ReplaceAllString(s, `${1}-${2}`)
	return strings.ToLower(kebab)
}
