package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var BuildWorkflow = workflow.Workflow{
	Name: "build",
	On:   BuildTriggers,
	Jobs: map[string]workflow.Job{
		"build":                  Build,
		"e2e-test":               E2eTest,
		"e2e-test-exec":          E2eTestExec,
		"e2etest-build":          E2etestBuild,
		"generate-metadata-file": GenerateMetadataFile,
		"get-go-version":         GetGoVersion,
		"get-product-version":    GetProductVersion,
		"package-docker":         PackageDocker,
	},
}
