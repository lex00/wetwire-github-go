package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var Build = workflow.Job{
	Name:   "Build Prometheus for common architectures",
	RunsOn: "ubuntu-latest",
	If:     "!(github.event_name == 'push' && startsWith(github.ref, 'refs/tags/v2.'))\n&&\n!(github.event_name == 'push' && startsWith(github.ref, 'refs/tags/v3.'))\n&&\n!(github.event_name == 'pull_request' && startsWith(github.event.pull_request.base.ref, 'release-'))\n&&\n!(github.event_name == 'push' && github.event.ref == 'refs/heads/main')\n",
	Steps:  BuildSteps,
}

var BuildAll = workflow.Job{
	Name:   "Build Prometheus for all architectures",
	RunsOn: "ubuntu-latest",
	If:     "(github.event_name == 'push' && startsWith(github.ref, 'refs/tags/v2.'))\n||\n(github.event_name == 'push' && startsWith(github.ref, 'refs/tags/v3.'))\n||\n(github.event_name == 'pull_request' && startsWith(github.event.pull_request.base.ref, 'release-'))\n||\n(github.event_name == 'push' && github.event.ref == 'refs/heads/main')\n",
	Steps:  BuildAllSteps,
}

var BuildAllStatus = workflow.Job{
	Name:   "Report status of build Prometheus for all architectures",
	RunsOn: "ubuntu-latest",
	Needs:  []any{"build_all"},
	If:     "always() && github.event_name == 'pull_request' && startsWith(github.event.pull_request.base.ref, 'release-')",
	Steps:  BuildAllStatusSteps,
}

var CheckGeneratedParser = workflow.Job{
	Name:   "Check generated parser",
	RunsOn: "ubuntu-latest",
	Steps:  CheckGeneratedParserSteps,
}

var Codeql = workflow.Job{}

var Fuzzing = workflow.Job{
	If: "github.event_name == 'pull_request'",
}

var Golangci = workflow.Job{
	Name:   "golangci-lint",
	RunsOn: "ubuntu-latest",
	Steps:  GolangciSteps,
}

var PublishMain = workflow.Job{
	Name:   "Publish main branch artifacts",
	RunsOn: "ubuntu-latest",
	Needs:  []any{"test_ui", "test_go", "test_go_more", "test_go_oldest", "test_windows", "golangci", "codeql", "build_all"},
	If:     "github.event_name == 'push' && github.event.ref == 'refs/heads/main'",
	Steps:  PublishMainSteps,
}

var PublishRelease = workflow.Job{
	Name:   "Publish release artefacts",
	RunsOn: "ubuntu-latest",
	Needs:  []any{"test_ui", "test_go", "test_go_more", "test_go_oldest", "test_windows", "golangci", "codeql", "build_all"},
	If:     "(github.event_name == 'push' && startsWith(github.ref, 'refs/tags/v2.'))\n||\n(github.event_name == 'push' && startsWith(github.ref, 'refs/tags/v3.'))\n",
	Steps:  PublishReleaseSteps,
}

var PublishUiRelease = workflow.Job{
	Name:   "Publish UI on npm Registry",
	RunsOn: "ubuntu-latest",
	Needs:  []any{"test_ui", "codeql"},
	Steps:  PublishUiReleaseSteps,
}

var TestGo = workflow.Job{
	Name:   "Go tests",
	RunsOn: "ubuntu-latest",
	Steps:  TestGoSteps,
}

var TestGoMore = workflow.Job{
	Name:   "More Go tests",
	RunsOn: "ubuntu-latest",
	Steps:  TestGoMoreSteps,
}

var TestGoOldest = workflow.Job{
	Name:   "Go tests with previous Go version",
	RunsOn: "ubuntu-latest",
	Steps:  TestGoOldestSteps,
}

var TestMixins = workflow.Job{
	Name:   "Mixins tests",
	RunsOn: "ubuntu-latest",
	Steps:  TestMixinsSteps,
}

var TestUi = workflow.Job{
	Name:   "UI tests",
	RunsOn: "ubuntu-latest",
	Steps:  TestUiSteps,
}

var TestWindows = workflow.Job{
	Name:   "Go tests on Windows",
	RunsOn: "windows-latest",
	Steps:  TestWindowsSteps,
}
