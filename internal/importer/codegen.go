package importer

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// CodeGenerator generates Go code from parsed YAML.
type CodeGenerator struct {
	// PackageName is the Go package name for generated code
	PackageName string
	// SingleFile puts all code in one file when true
	SingleFile bool
}

// NewCodeGenerator creates a new CodeGenerator.
func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{
		PackageName: "workflows",
	}
}

// GeneratedCode represents generated Go source code.
type GeneratedCode struct {
	// Files maps filenames to their content
	Files map[string]string
	// Workflows is the count of generated workflows
	Workflows int
	// Jobs is the count of generated jobs
	Jobs int
	// Steps is the count of generated steps
	Steps int
}

// Generate generates Go code from a parsed workflow.
func (g *CodeGenerator) Generate(workflow *IRWorkflow, workflowName string) (*GeneratedCode, error) {
	result := &GeneratedCode{
		Files: make(map[string]string),
	}

	if g.SingleFile {
		code := g.generateSingleFile(workflow, workflowName)
		result.Files["workflows.go"] = code
	} else {
		// Generate separate files
		result.Files["workflows.go"] = g.generateWorkflowFile(workflow, workflowName)
		result.Files["triggers.go"] = g.generateTriggersFile(workflow, workflowName)
		result.Files["jobs.go"] = g.generateJobsFile(workflow)
		result.Files["steps.go"] = g.generateStepsFile(workflow)
	}

	// Count resources
	result.Workflows = 1
	result.Jobs = len(workflow.Jobs)
	for _, job := range workflow.Jobs {
		result.Steps += len(job.Steps)
	}

	return result, nil
}

// generateSingleFile generates all code in a single file.
func (g *CodeGenerator) generateSingleFile(workflow *IRWorkflow, workflowName string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("package %s\n\n", g.PackageName))
	sb.WriteString("import (\n")
	sb.WriteString("\t\"github.com/lex00/wetwire-github-go/workflow\"\n")
	sb.WriteString(")\n\n")

	// Generate workflow
	varName := toVarName(workflowName)
	sb.WriteString(g.generateWorkflow(workflow, varName))
	sb.WriteString("\n")

	// Generate triggers
	sb.WriteString(g.generateTriggers(workflow, varName))
	sb.WriteString("\n")

	// Generate jobs
	for jobID, job := range workflow.Jobs {
		sb.WriteString(g.generateJob(jobID, job))
		sb.WriteString("\n")
	}

	return sb.String()
}

// generateWorkflowFile generates the workflows.go file.
func (g *CodeGenerator) generateWorkflowFile(workflow *IRWorkflow, workflowName string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("package %s\n\n", g.PackageName))
	sb.WriteString("import (\n")
	sb.WriteString("\t\"github.com/lex00/wetwire-github-go/workflow\"\n")
	sb.WriteString(")\n\n")

	varName := toVarName(workflowName)
	sb.WriteString(g.generateWorkflow(workflow, varName))

	return sb.String()
}

// generateTriggersFile generates the triggers.go file.
func (g *CodeGenerator) generateTriggersFile(workflow *IRWorkflow, workflowName string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("package %s\n\n", g.PackageName))
	sb.WriteString("import (\n")
	sb.WriteString("\t\"github.com/lex00/wetwire-github-go/workflow\"\n")
	sb.WriteString(")\n\n")

	varName := toVarName(workflowName)
	sb.WriteString(g.generateTriggers(workflow, varName))

	return sb.String()
}

// generateJobsFile generates the jobs.go file.
func (g *CodeGenerator) generateJobsFile(workflow *IRWorkflow) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("package %s\n\n", g.PackageName))
	sb.WriteString("import (\n")
	sb.WriteString("\t\"github.com/lex00/wetwire-github-go/workflow\"\n")
	sb.WriteString(")\n\n")

	// Sort job IDs for deterministic output
	jobIDs := make([]string, 0, len(workflow.Jobs))
	for id := range workflow.Jobs {
		jobIDs = append(jobIDs, id)
	}
	sort.Strings(jobIDs)

	for _, jobID := range jobIDs {
		job := workflow.Jobs[jobID]
		sb.WriteString(g.generateJob(jobID, job))
		sb.WriteString("\n")
	}

	return sb.String()
}

// generateStepsFile generates the steps.go file.
func (g *CodeGenerator) generateStepsFile(workflow *IRWorkflow) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("package %s\n\n", g.PackageName))
	sb.WriteString("import (\n")
	sb.WriteString("\t\"github.com/lex00/wetwire-github-go/workflow\"\n")
	sb.WriteString(")\n\n")

	// Sort job IDs for deterministic output
	jobIDs := make([]string, 0, len(workflow.Jobs))
	for id := range workflow.Jobs {
		jobIDs = append(jobIDs, id)
	}
	sort.Strings(jobIDs)

	for _, jobID := range jobIDs {
		job := workflow.Jobs[jobID]
		if len(job.Steps) > 0 {
			varName := toVarName(jobID) + "Steps"
			sb.WriteString(g.generateSteps(varName, job.Steps))
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// generateWorkflow generates the workflow variable.
func (g *CodeGenerator) generateWorkflow(workflow *IRWorkflow, varName string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("var %s = workflow.Workflow{\n", varName))

	if workflow.Name != "" {
		sb.WriteString(fmt.Sprintf("\tName: %q,\n", workflow.Name))
	}

	sb.WriteString(fmt.Sprintf("\tOn:   %sTriggers,\n", varName))

	// Jobs map
	if len(workflow.Jobs) > 0 {
		sb.WriteString("\tJobs: map[string]workflow.Job{\n")
		jobIDs := make([]string, 0, len(workflow.Jobs))
		for id := range workflow.Jobs {
			jobIDs = append(jobIDs, id)
		}
		sort.Strings(jobIDs)
		for _, id := range jobIDs {
			sb.WriteString(fmt.Sprintf("\t\t%q: %s,\n", id, toVarName(id)))
		}
		sb.WriteString("\t},\n")
	}

	sb.WriteString("}\n")
	return sb.String()
}

// generateTriggers generates trigger variables.
func (g *CodeGenerator) generateTriggers(workflow *IRWorkflow, varName string) string {
	var sb strings.Builder

	// Generate individual trigger variables first
	if workflow.On.Push != nil {
		sb.WriteString(g.generatePushTriggerVar(workflow.On.Push, varName))
		sb.WriteString("\n")
	}
	if workflow.On.PullRequest != nil {
		sb.WriteString(g.generatePullRequestTriggerVar(workflow.On.PullRequest, varName))
		sb.WriteString("\n")
	}

	// Generate main triggers struct
	sb.WriteString(fmt.Sprintf("var %sTriggers = workflow.Triggers{\n", varName))

	if workflow.On.Push != nil {
		sb.WriteString(fmt.Sprintf("\tPush: &%sPush,\n", varName))
	}
	if workflow.On.PullRequest != nil {
		sb.WriteString(fmt.Sprintf("\tPullRequest: &%sPullRequest,\n", varName))
	}
	if workflow.On.WorkflowDispatch != nil {
		sb.WriteString("\tWorkflowDispatch: &workflow.WorkflowDispatchTrigger{},\n")
	}
	if workflow.On.WorkflowCall != nil {
		sb.WriteString("\tWorkflowCall: &workflow.WorkflowCallTrigger{},\n")
	}
	if len(workflow.On.Schedule) > 0 {
		sb.WriteString("\tSchedule: []workflow.ScheduleTrigger{\n")
		for _, s := range workflow.On.Schedule {
			sb.WriteString(fmt.Sprintf("\t\t{Cron: %q},\n", s.Cron))
		}
		sb.WriteString("\t},\n")
	}

	sb.WriteString("}\n")
	return sb.String()
}

// generatePushTriggerVar generates the push trigger variable definition.
func (g *CodeGenerator) generatePushTriggerVar(push *IRPushTrigger, varName string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("var %sPush = workflow.PushTrigger{\n", varName))

	if len(push.Branches) > 0 {
		sb.WriteString("\tBranches: []string{")
		for i, b := range push.Branches {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("%q", b))
		}
		sb.WriteString("},\n")
	}

	if len(push.BranchesIgnore) > 0 {
		sb.WriteString("\tBranchesIgnore: []string{")
		for i, b := range push.BranchesIgnore {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("%q", b))
		}
		sb.WriteString("},\n")
	}

	if len(push.Tags) > 0 {
		sb.WriteString("\tTags: []string{")
		for i, t := range push.Tags {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("%q", t))
		}
		sb.WriteString("},\n")
	}

	if len(push.Paths) > 0 {
		sb.WriteString("\tPaths: []string{")
		for i, p := range push.Paths {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("%q", p))
		}
		sb.WriteString("},\n")
	}

	sb.WriteString("}\n")
	return sb.String()
}

// generatePullRequestTriggerVar generates the pull request trigger variable definition.
func (g *CodeGenerator) generatePullRequestTriggerVar(pr *IRPullRequestTrigger, varName string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("var %sPullRequest = workflow.PullRequestTrigger{\n", varName))

	if len(pr.Branches) > 0 {
		sb.WriteString("\tBranches: []string{")
		for i, b := range pr.Branches {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("%q", b))
		}
		sb.WriteString("},\n")
	}

	if len(pr.BranchesIgnore) > 0 {
		sb.WriteString("\tBranchesIgnore: []string{")
		for i, b := range pr.BranchesIgnore {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("%q", b))
		}
		sb.WriteString("},\n")
	}

	if len(pr.Types) > 0 {
		sb.WriteString("\tTypes: []string{")
		for i, t := range pr.Types {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("%q", t))
		}
		sb.WriteString("},\n")
	}

	if len(pr.Paths) > 0 {
		sb.WriteString("\tPaths: []string{")
		for i, p := range pr.Paths {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("%q", p))
		}
		sb.WriteString("},\n")
	}

	sb.WriteString("}\n")
	return sb.String()
}

// generateJob generates a job variable.
func (g *CodeGenerator) generateJob(jobID string, job *IRJob) string {
	var sb strings.Builder
	varName := toVarName(jobID)

	sb.WriteString(fmt.Sprintf("var %s = workflow.Job{\n", varName))

	if job.Name != "" {
		sb.WriteString(fmt.Sprintf("\tName: %q,\n", job.Name))
	}

	// RunsOn
	if job.RunsOn != nil {
		runsOn := job.GetRunsOn()
		if runsOn != "" {
			sb.WriteString(fmt.Sprintf("\tRunsOn: %q,\n", runsOn))
		}
	}

	// Needs
	needs := job.GetNeeds()
	if len(needs) > 0 {
		sb.WriteString("\tNeeds: []any{")
		for i, n := range needs {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("%q", n))
		}
		sb.WriteString("},\n")
	}

	// If condition
	if job.If != "" {
		sb.WriteString(fmt.Sprintf("\tIf: %q,\n", job.If))
	}

	// TimeoutMinutes
	if job.TimeoutMinutes > 0 {
		sb.WriteString(fmt.Sprintf("\tTimeoutMinutes: %d,\n", job.TimeoutMinutes))
	}

	// Steps reference
	if len(job.Steps) > 0 {
		sb.WriteString(fmt.Sprintf("\tSteps: %sSteps,\n", varName))
	}

	sb.WriteString("}\n")
	return sb.String()
}

// generateSteps generates a steps slice variable.
func (g *CodeGenerator) generateSteps(varName string, steps []IRStep) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("var %s = []workflow.Step{\n", varName))

	for _, step := range steps {
		sb.WriteString("\t{\n")

		if step.ID != "" {
			sb.WriteString(fmt.Sprintf("\t\tID: %q,\n", step.ID))
		}
		if step.Name != "" {
			sb.WriteString(fmt.Sprintf("\t\tName: %q,\n", step.Name))
		}
		if step.Uses != "" {
			sb.WriteString(fmt.Sprintf("\t\tUses: %q,\n", step.Uses))
		}
		if step.Run != "" {
			// Handle multiline strings
			if strings.Contains(step.Run, "\n") {
				sb.WriteString(fmt.Sprintf("\t\tRun: `%s`,\n", step.Run))
			} else {
				sb.WriteString(fmt.Sprintf("\t\tRun: %q,\n", step.Run))
			}
		}
		if step.Shell != "" {
			sb.WriteString(fmt.Sprintf("\t\tShell: %q,\n", step.Shell))
		}
		if step.If != "" {
			sb.WriteString(fmt.Sprintf("\t\tIf: %q,\n", step.If))
		}
		if step.WorkingDirectory != "" {
			sb.WriteString(fmt.Sprintf("\t\tWorkingDirectory: %q,\n", step.WorkingDirectory))
		}
		if step.TimeoutMinutes > 0 {
			sb.WriteString(fmt.Sprintf("\t\tTimeoutMinutes: %d,\n", step.TimeoutMinutes))
		}

		// With map
		if len(step.With) > 0 {
			sb.WriteString("\t\tWith: map[string]any{\n")
			keys := make([]string, 0, len(step.With))
			for k := range step.With {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				v := step.With[k]
				sb.WriteString(fmt.Sprintf("\t\t\t%q: %s,\n", k, formatValue(v)))
			}
			sb.WriteString("\t\t},\n")
		}

		// Env map
		if len(step.Env) > 0 {
			sb.WriteString("\t\tEnv: map[string]any{\n")
			keys := make([]string, 0, len(step.Env))
			for k := range step.Env {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				v := step.Env[k]
				sb.WriteString(fmt.Sprintf("\t\t\t%q: %s,\n", k, formatValue(v)))
			}
			sb.WriteString("\t\t},\n")
		}

		sb.WriteString("\t},\n")
	}

	sb.WriteString("}\n")
	return sb.String()
}

// toVarName converts a string to a valid Go variable name.
func toVarName(s string) string {
	// Replace special characters with meaningful substitutes
	s = strings.ReplaceAll(s, "++", "pp")  // C++ -> Cpp
	s = strings.ReplaceAll(s, "+", "Plus") // single + -> Plus
	s = strings.ReplaceAll(s, "#", "Sharp") // C# -> CSharp
	s = strings.ReplaceAll(s, "&", "And")

	// Remove or replace non-alphanumeric characters
	s = strings.ReplaceAll(s, "(", "")
	s = strings.ReplaceAll(s, ")", "")
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.ReplaceAll(s, ".", "_")
	s = strings.ReplaceAll(s, "@", "_")
	s = strings.ReplaceAll(s, "!", "")
	s = strings.ReplaceAll(s, "?", "")
	s = strings.ReplaceAll(s, "'", "")
	s = strings.ReplaceAll(s, "\"", "")
	s = strings.ReplaceAll(s, ":", "_")
	s = strings.ReplaceAll(s, ";", "_")
	s = strings.ReplaceAll(s, ",", "_")

	// Convert spaces, kebab-case and snake_case to PascalCase
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "-", "_")
	parts := strings.Split(s, "_")
	var result []string
	for _, part := range parts {
		if len(part) > 0 {
			result = append(result, strings.ToUpper(part[:1])+part[1:])
		}
	}
	name := strings.Join(result, "")

	// Handle reserved words
	if isReserved(name) {
		name = name + "Job"
	}

	// Ensure the name starts with a letter (required for Go identifiers)
	if len(name) > 0 && !isLetter(rune(name[0])) {
		name = "X" + name
	}

	return name
}

// isLetter checks if a rune is a letter.
func isLetter(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
}

// isReserved checks if a name is a Go reserved word.
func isReserved(name string) bool {
	reserved := map[string]bool{
		"break": true, "case": true, "chan": true, "const": true,
		"continue": true, "default": true, "defer": true, "else": true,
		"fallthrough": true, "for": true, "func": true, "go": true,
		"goto": true, "if": true, "import": true, "interface": true,
		"map": true, "package": true, "range": true, "return": true,
		"select": true, "struct": true, "switch": true, "type": true,
		"var": true,
	}
	return reserved[strings.ToLower(name)]
}

// formatValue formats a value for Go source code.
func formatValue(v any) string {
	switch val := v.(type) {
	case string:
		return fmt.Sprintf("%q", val)
	case int, int64, float64:
		return fmt.Sprintf("%v", val)
	case bool:
		return fmt.Sprintf("%t", val)
	default:
		return fmt.Sprintf("%q", fmt.Sprintf("%v", v))
	}
}

// toFilename converts a workflow name to a valid filename.
var nonAlphanumericRE = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func toFilename(name string) string {
	name = nonAlphanumericRE.ReplaceAllString(name, "-")
	name = strings.Trim(name, "-")
	return strings.ToLower(name)
}
