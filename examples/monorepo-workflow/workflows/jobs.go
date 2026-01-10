package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// DetectChanges determines which services have changed.
// Outputs are used by downstream jobs to conditionally execute.
var DetectChanges = workflow.Job{
	Name:   "Detect Changes",
	RunsOn: "ubuntu-latest",
	Outputs: map[string]any{
		"api":    workflow.Steps.Get("changes", "api").String(),
		"web":    workflow.Steps.Get("changes", "web").String(),
		"shared": workflow.Steps.Get("changes", "shared").String(),
	},
	Steps: DetectChangesSteps,
}

// APICondition checks if API or shared files changed.
// Runs when api == 'true' OR shared == 'true'.
var APICondition = workflow.Needs.Get("detect-changes", "api").
	Or(workflow.Needs.Get("detect-changes", "shared"))

// BuildAPI builds and tests the API service.
// Only runs if API files or shared library changed.
var BuildAPI = workflow.Job{
	Name:   "Build API",
	RunsOn: "ubuntu-latest",
	Needs:  []any{DetectChanges},
	If:     APICondition.String(),
	Steps:  APIBuildSteps,
}

// WebCondition checks if Web or shared files changed.
// Runs when web == 'true' OR shared == 'true'.
var WebCondition = workflow.Needs.Get("detect-changes", "web").
	Or(workflow.Needs.Get("detect-changes", "shared"))

// BuildWeb builds and tests the Web service.
// Only runs if Web files or shared library changed.
var BuildWeb = workflow.Job{
	Name:   "Build Web",
	RunsOn: "ubuntu-latest",
	Needs:  []any{DetectChanges},
	If:     WebCondition.String(),
	Steps:  WebBuildSteps,
}

// SharedCondition checks if shared files changed.
var SharedCondition = workflow.Needs.Get("detect-changes", "shared")

// BuildShared builds and tests the Shared library.
// Only runs if shared library files changed.
var BuildShared = workflow.Job{
	Name:   "Build Shared",
	RunsOn: "ubuntu-latest",
	Needs:  []any{DetectChanges},
	If:     SharedCondition.String(),
	Steps:  SharedBuildSteps,
}
