package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var CI = workflow.Workflow{
	Name: "CI",
	On:   CITriggers,
	Jobs: map[string]workflow.Job{
		"build":                  Build,
		"build_all":              BuildAll,
		"build_all_status":       BuildAllStatus,
		"check_generated_parser": CheckGeneratedParser,
		"codeql":                 Codeql,
		"fuzzing":                Fuzzing,
		"golangci":               Golangci,
		"publish_main":           PublishMain,
		"publish_release":        PublishRelease,
		"publish_ui_release":     PublishUiRelease,
		"test_go":                TestGo,
		"test_go_more":           TestGoMore,
		"test_go_oldest":         TestGoOldest,
		"test_mixins":            TestMixins,
		"test_ui":                TestUi,
		"test_windows":           TestWindows,
	},
}
