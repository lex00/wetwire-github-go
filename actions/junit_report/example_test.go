package junit_report_test

import (
	"fmt"

	"github.com/lex00/wetwire-github-go/actions/checkout"
	"github.com/lex00/wetwire-github-go/actions/junit_report"
	"github.com/lex00/wetwire-github-go/actions/setup_go"
	"github.com/lex00/wetwire-github-go/workflow"
)

func ExampleJUnitReport() {
	// Define a test job that runs tests and publishes JUnit reports
	testJob := workflow.Job{
		Name:   "test",
		RunsOn: "ubuntu-latest",
		Steps: []any{
			checkout.Checkout{},
			setup_go.SetupGo{
				GoVersion: "1.23",
			},
			workflow.Step{
				Name: "Run tests",
				Run:  "go test -v ./... -coverprofile=coverage.out 2>&1 | tee test-results.txt",
			},
			workflow.Step{
				Name: "Convert results to JUnit",
				Run:  "go install github.com/jstemmer/go-junit-report/v2@latest && cat test-results.txt | go-junit-report -set-exit-code > junit-report.xml",
			},
			junit_report.JUnitReport{
				ReportPaths:    "**/junit-report.xml",
				CheckName:      "Test Results",
				FailOnFailure:  true,
				DetailedSummary: true,
			},
		},
	}

	fmt.Println(testJob.Name)
	// Output: test
}

func ExampleJUnitReport_minimal() {
	// Minimal configuration - just provide the report paths
	step := junit_report.JUnitReport{
		ReportPaths: "**/test-results/*.xml",
	}

	fmt.Println(step.Action())
	// Output: mikepenz/action-junit-report@v4
}

func ExampleJUnitReport_withPRComments() {
	// Configure JUnit report with PR comments enabled
	step := junit_report.JUnitReport{
		ReportPaths:    "**/junit-reports/*.xml",
		CheckName:      "Unit Test Results",
		Comment:        true,
		UpdateComment:  true,
		DetailedSummary: true,
		FlakySummary:   true,
	}

	inputs := step.Inputs()
	fmt.Println(inputs["comment"])
	fmt.Println(inputs["detailed_summary"])
	// Output:
	// true
	// true
}
