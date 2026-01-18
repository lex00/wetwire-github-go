package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var Build = workflow.Job{
	Name:  "Build for ${{ matrix.goos }}_${{ matrix.goarch }}",
	Needs: []any{"get-product-version", "get-go-version"},
}

var E2eTest = workflow.Job{
	Name:   "Run e2e test for ${{ matrix.goos }}_${{ matrix.goarch }}",
	RunsOn: "${{ matrix.runson }}",
	Needs:  []any{"get-product-version", "build", "e2etest-build"},
	Steps:  E2eTestSteps,
}

var E2eTestExec = workflow.Job{
	Name:   "Run terraform-exec test for linux amd64",
	RunsOn: "ubuntu-latest",
	Needs:  []any{"get-product-version", "get-go-version", "build"},
	Steps:  E2eTestExecSteps,
}

var E2etestBuild = workflow.Job{
	Name:   "Build e2etest for ${{ matrix.goos }}_${{ matrix.goarch }}",
	RunsOn: "ubuntu-latest",
	Needs:  []any{"get-product-version", "get-go-version"},
	Steps:  E2etestBuildSteps,
}

var GenerateMetadataFile = workflow.Job{
	Name:   "Generate release metadata",
	RunsOn: "ubuntu-latest",
	Needs:  []any{"get-product-version"},
	Steps:  GenerateMetadataFileSteps,
}

var GetGoVersion = workflow.Job{
	Name:   "Determine Go toolchain version",
	RunsOn: "ubuntu-latest",
	Steps:  GetGoVersionSteps,
}

var GetProductVersion = workflow.Job{
	Name:   "Determine intended Terraform version",
	RunsOn: "ubuntu-latest",
	Steps:  GetProductVersionSteps,
}

var PackageDocker = workflow.Job{
	Name:   "Build Docker image for linux_${{ matrix.arch }}",
	RunsOn: "ubuntu-latest",
	Needs:  []any{"get-product-version", "build"},
	Steps:  PackageDockerSteps,
}
