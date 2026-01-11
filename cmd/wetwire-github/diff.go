package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/lex00/wetwire-github-go/internal/discover"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var diffCmd = &cobra.Command{
	Use:   "diff <path1> <path2>",
	Short: "Compare two workflow configurations",
	Long: `Semantically compare two workflow configurations and show differences.

Shows added jobs, removed jobs, modified jobs, and dependency changes.

By default, compares two Go packages. Use --yaml to compare YAML files.

Supported output formats:
  - text (default): Human-readable output
  - json: Machine-readable JSON format
  - markdown: Markdown formatted diff report

Examples:
  # Compare two Go packages
  wetwire-github diff ./old-workflows ./new-workflows

  # Compare two YAML files
  wetwire-github diff old-ci.yml new-ci.yml --yaml

  # JSON output for automation
  wetwire-github diff ./v1 ./v2 --format json

  # Markdown report
  wetwire-github diff ./v1 ./v2 --format markdown`,
	Args: cobra.ExactArgs(2),
	RunE: runDiff,
}

func init() {
	diffCmd.Flags().Bool("yaml", false, "Compare YAML files instead of Go packages")
	diffCmd.Flags().String("format", "text", "Output format: text, json, markdown")
}

// jobInfo represents a discovered job with its dependencies.
type jobInfo struct {
	Name         string
	Dependencies []string
}

// jobDiff represents changes to a single job.
type jobDiff struct {
	Name    string   `json:"name"`
	Changes []string `json:"changes,omitempty"`
}

// dependencyChange represents dependency changes for a job.
type dependencyChange struct {
	Job         string   `json:"job"`
	AddedDeps   []string `json:"added_deps,omitempty"`
	RemovedDeps []string `json:"removed_deps,omitempty"`
}

// diffResult holds the comparison results.
type diffResult struct {
	Success           bool               `json:"success"`
	Message           string             `json:"message,omitempty"`
	AddedJobs         []string           `json:"added_jobs,omitempty"`
	RemovedJobs       []string           `json:"removed_jobs,omitempty"`
	ModifiedJobs      []jobDiff          `json:"modified_jobs,omitempty"`
	DependencyChanges []dependencyChange `json:"dependency_changes,omitempty"`
}

func runDiff(cmd *cobra.Command, args []string) error {
	yamlMode, _ := cmd.Flags().GetBool("yaml")
	outputFormat, _ := cmd.Flags().GetString("format")

	path1 := args[0]
	path2 := args[1]

	var oldJobs, newJobs []jobInfo
	var err error

	if yamlMode {
		oldJobs, err = parseYAMLWorkflowJobs(path1)
		if err != nil {
			return outputDiffError(cmd, outputFormat, fmt.Errorf("parse %s: %w", path1, err))
		}

		newJobs, err = parseYAMLWorkflowJobs(path2)
		if err != nil {
			return outputDiffError(cmd, outputFormat, fmt.Errorf("parse %s: %w", path2, err))
		}
	} else {
		disc := discover.NewDiscoverer()

		result1, err := disc.Discover(path1)
		if err != nil {
			return outputDiffError(cmd, outputFormat, fmt.Errorf("discover %s: %w", path1, err))
		}
		oldJobs = convertDiscoveredJobs(result1.Jobs)

		result2, err := disc.Discover(path2)
		if err != nil {
			return outputDiffError(cmd, outputFormat, fmt.Errorf("discover %s: %w", path2, err))
		}
		newJobs = convertDiscoveredJobs(result2.Jobs)
	}

	result := compareWorkflows(oldJobs, newJobs)
	result.Success = true

	switch outputFormat {
	case "json":
		return outputDiffJSON(cmd, result)
	case "markdown":
		outputDiffMarkdown(cmd, result)
	default:
		outputDiffText(cmd, result)
	}

	return nil
}

func convertDiscoveredJobs(jobs []discover.DiscoveredJob) []jobInfo {
	result := make([]jobInfo, len(jobs))
	for i, job := range jobs {
		result[i] = jobInfo{
			Name:         job.Name,
			Dependencies: job.Dependencies,
		}
	}
	return result
}

func compareWorkflows(oldJobs, newJobs []jobInfo) diffResult {
	result := diffResult{}

	oldJobMap := make(map[string]jobInfo)
	newJobMap := make(map[string]jobInfo)

	for _, job := range oldJobs {
		oldJobMap[job.Name] = job
	}
	for _, job := range newJobs {
		newJobMap[job.Name] = job
	}

	// Find added jobs
	for _, newJob := range newJobs {
		if _, exists := oldJobMap[newJob.Name]; !exists {
			result.AddedJobs = append(result.AddedJobs, newJob.Name)
		}
	}

	// Find removed jobs
	for _, oldJob := range oldJobs {
		if _, exists := newJobMap[oldJob.Name]; !exists {
			result.RemovedJobs = append(result.RemovedJobs, oldJob.Name)
		}
	}

	// Find modified jobs and dependency changes
	for _, newJob := range newJobs {
		oldJob, exists := oldJobMap[newJob.Name]
		if !exists {
			continue
		}

		var changes []string

		// Check dependency changes
		oldDepsSet := make(map[string]bool)
		newDepsSet := make(map[string]bool)
		for _, dep := range oldJob.Dependencies {
			oldDepsSet[dep] = true
		}
		for _, dep := range newJob.Dependencies {
			newDepsSet[dep] = true
		}

		var addedDeps, removedDeps []string
		for dep := range newDepsSet {
			if !oldDepsSet[dep] {
				addedDeps = append(addedDeps, dep)
			}
		}
		for dep := range oldDepsSet {
			if !newDepsSet[dep] {
				removedDeps = append(removedDeps, dep)
			}
		}

		if len(addedDeps) > 0 || len(removedDeps) > 0 {
			sort.Strings(addedDeps)
			sort.Strings(removedDeps)
			result.DependencyChanges = append(result.DependencyChanges, dependencyChange{
				Job:         newJob.Name,
				AddedDeps:   addedDeps,
				RemovedDeps: removedDeps,
			})

			if len(addedDeps) > 0 {
				changes = append(changes, fmt.Sprintf("added dependencies: %v", addedDeps))
			}
			if len(removedDeps) > 0 {
				changes = append(changes, fmt.Sprintf("removed dependencies: %v", removedDeps))
			}
		}

		if len(changes) > 0 {
			result.ModifiedJobs = append(result.ModifiedJobs, jobDiff{
				Name:    newJob.Name,
				Changes: changes,
			})
		}
	}

	sort.Strings(result.AddedJobs)
	sort.Strings(result.RemovedJobs)

	return result
}

func parseYAMLWorkflowJobs(path string) ([]jobInfo, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var workflow map[string]interface{}
	if err := yaml.Unmarshal(data, &workflow); err != nil {
		return nil, err
	}

	jobs, ok := workflow["jobs"].(map[string]interface{})
	if !ok {
		return []jobInfo{}, nil
	}

	var result []jobInfo
	for name, value := range jobs {
		jobMap, ok := value.(map[string]interface{})
		if !ok {
			continue
		}

		job := jobInfo{Name: name}

		// Extract dependencies from needs
		if needs, ok := jobMap["needs"]; ok {
			switch n := needs.(type) {
			case string:
				job.Dependencies = []string{n}
			case []interface{}:
				for _, need := range n {
					if s, ok := need.(string); ok {
						job.Dependencies = append(job.Dependencies, s)
					}
				}
			}
		}

		result = append(result, job)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result, nil
}

func outputDiffError(cmd *cobra.Command, outputFormat string, err error) error {
	result := diffResult{
		Success: false,
		Message: err.Error(),
	}

	if outputFormat == "json" {
		return outputDiffJSON(cmd, result)
	}

	fmt.Fprintf(cmd.ErrOrStderr(), "Error: %v\n", err)
	return nil
}

func outputDiffText(cmd *cobra.Command, result diffResult) {
	hasChanges := len(result.AddedJobs) > 0 ||
		len(result.RemovedJobs) > 0 ||
		len(result.ModifiedJobs) > 0 ||
		len(result.DependencyChanges) > 0

	if !hasChanges {
		fmt.Fprintln(cmd.OutOrStdout(), "No differences found.")
		return
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Workflow Diff")
	fmt.Fprintln(cmd.OutOrStdout(), "=============")
	fmt.Fprintln(cmd.OutOrStdout())

	if len(result.AddedJobs) > 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "Added Jobs:")
		for _, job := range result.AddedJobs {
			fmt.Fprintf(cmd.OutOrStdout(), "  + %s\n", job)
		}
		fmt.Fprintln(cmd.OutOrStdout())
	}

	if len(result.RemovedJobs) > 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "Removed Jobs:")
		for _, job := range result.RemovedJobs {
			fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", job)
		}
		fmt.Fprintln(cmd.OutOrStdout())
	}

	if len(result.ModifiedJobs) > 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "Modified Jobs:")
		for _, job := range result.ModifiedJobs {
			fmt.Fprintf(cmd.OutOrStdout(), "  ~ %s\n", job.Name)
			for _, change := range job.Changes {
				fmt.Fprintf(cmd.OutOrStdout(), "    - %s\n", change)
			}
		}
		fmt.Fprintln(cmd.OutOrStdout())
	}

	if len(result.DependencyChanges) > 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "Dependency Changes:")
		for _, dc := range result.DependencyChanges {
			fmt.Fprintf(cmd.OutOrStdout(), "  %s:\n", dc.Job)
			if len(dc.AddedDeps) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "    + %v\n", dc.AddedDeps)
			}
			if len(dc.RemovedDeps) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "    - %v\n", dc.RemovedDeps)
			}
		}
		fmt.Fprintln(cmd.OutOrStdout())
	}
}

func outputDiffJSON(cmd *cobra.Command, result diffResult) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "Error marshaling JSON: %v\n", err)
		return nil
	}
	fmt.Fprintln(cmd.OutOrStdout(), string(data))
	return nil
}

func outputDiffMarkdown(cmd *cobra.Command, result diffResult) {
	hasChanges := len(result.AddedJobs) > 0 ||
		len(result.RemovedJobs) > 0 ||
		len(result.ModifiedJobs) > 0 ||
		len(result.DependencyChanges) > 0

	fmt.Fprintln(cmd.OutOrStdout(), "## Workflow Diff")
	fmt.Fprintln(cmd.OutOrStdout())

	if !hasChanges {
		fmt.Fprintln(cmd.OutOrStdout(), "No differences found.")
		return
	}

	if len(result.AddedJobs) > 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "### Added Jobs")
		fmt.Fprintln(cmd.OutOrStdout())
		for _, job := range result.AddedJobs {
			fmt.Fprintf(cmd.OutOrStdout(), "- `%s`\n", job)
		}
		fmt.Fprintln(cmd.OutOrStdout())
	}

	if len(result.RemovedJobs) > 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "### Removed Jobs")
		fmt.Fprintln(cmd.OutOrStdout())
		for _, job := range result.RemovedJobs {
			fmt.Fprintf(cmd.OutOrStdout(), "- `%s`\n", job)
		}
		fmt.Fprintln(cmd.OutOrStdout())
	}

	if len(result.ModifiedJobs) > 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "### Modified Jobs")
		fmt.Fprintln(cmd.OutOrStdout())
		for _, job := range result.ModifiedJobs {
			fmt.Fprintf(cmd.OutOrStdout(), "#### `%s`\n\n", job.Name)
			for _, change := range job.Changes {
				fmt.Fprintf(cmd.OutOrStdout(), "- %s\n", change)
			}
			fmt.Fprintln(cmd.OutOrStdout())
		}
	}

	if len(result.DependencyChanges) > 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "### Dependency Changes")
		fmt.Fprintln(cmd.OutOrStdout())
		for _, dc := range result.DependencyChanges {
			fmt.Fprintf(cmd.OutOrStdout(), "#### `%s`\n\n", dc.Job)
			if len(dc.AddedDeps) > 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "Added dependencies:")
				for _, dep := range dc.AddedDeps {
					fmt.Fprintf(cmd.OutOrStdout(), "- `%s`\n", dep)
				}
			}
			if len(dc.RemovedDeps) > 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "Removed dependencies:")
				for _, dep := range dc.RemovedDeps {
					fmt.Fprintf(cmd.OutOrStdout(), "- `%s`\n", dep)
				}
			}
			fmt.Fprintln(cmd.OutOrStdout())
		}
	}
}
