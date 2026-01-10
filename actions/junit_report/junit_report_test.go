package junit_report

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestJUnitReport_Action(t *testing.T) {
	jr := JUnitReport{}
	if got := jr.Action(); got != "mikepenz/action-junit-report@v4" {
		t.Errorf("Action() = %q, want %q", got, "mikepenz/action-junit-report@v4")
	}
}

func TestJUnitReport_Inputs(t *testing.T) {
	jr := JUnitReport{
		ReportPaths: "**/test-results/*.xml",
		CheckName:   "Test Results",
		Token:       "${{ secrets.GITHUB_TOKEN }}",
	}

	inputs := jr.Inputs()

	if jr.Action() != "mikepenz/action-junit-report@v4" {
		t.Errorf("Action() = %q, want %q", jr.Action(), "mikepenz/action-junit-report@v4")
	}

	if inputs["report_paths"] != "**/test-results/*.xml" {
		t.Errorf("inputs[report_paths] = %v, want %q", inputs["report_paths"], "**/test-results/*.xml")
	}

	if inputs["check_name"] != "Test Results" {
		t.Errorf("inputs[check_name] = %v, want %q", inputs["check_name"], "Test Results")
	}

	if inputs["token"] != "${{ secrets.GITHUB_TOKEN }}" {
		t.Errorf("inputs[token] = %v, want %q", inputs["token"], "${{ secrets.GITHUB_TOKEN }}")
	}
}

func TestJUnitReport_Inputs_EmptyWithMaps(t *testing.T) {
	jr := JUnitReport{}
	inputs := jr.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty JUnitReport.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestJUnitReport_Inputs_AllFields(t *testing.T) {
	jr := JUnitReport{
		ReportPaths:             "**/test-results/*.xml",
		Token:                   "${{ secrets.GITHUB_TOKEN }}",
		GroupReports:            true,
		TestFilesPrefix:         "src/",
		ExcludeSources:          "/build/,/node_modules/",
		CheckName:               "Custom Test Report",
		Commit:                  "${{ github.sha }}",
		FailOnFailure:           true,
		FailOnParseError:        true,
		RequireTests:            true,
		RequirePassedTests:      true,
		IncludePassed:           true,
		IncludeSkipped:          true,
		CheckRetries:            true,
		CheckTitleTemplate:      "{{SUITE_NAME}} - {{TEST_NAME}}",
		BreadCrumbDelimiter:     ".",
		Summary:                 "Test execution summary",
		CheckAnnotations:        true,
		UpdateCheck:             true,
		AnnotateOnly:            true,
		Transformers:            "[]",
		JobSummary:              true,
		JobSummaryText:          "Additional summary text",
		DetailedSummary:         true,
		FlakySummary:            true,
		VerboseSummary:          true,
		SkipSuccessSummary:      true,
		IncludeEmptyInSummary:   true,
		IncludeTimeInSummary:    true,
		SimplifiedSummary:       true,
		GroupSuite:              true,
		Comment:                 true,
		UpdateComment:           true,
		AnnotateNotice:          true,
		FollowSymlink:           true,
		JobName:                 "test-job",
		AnnotationsLimit:        10,
		SkipAnnotations:         true,
		TruncateStackTraces:     true,
		ResolveIgnoreClassname:  true,
		SkipCommentWithoutTests: true,
		PrID:                    123,
	}

	inputs := jr.Inputs()

	if inputs["report_paths"] != "**/test-results/*.xml" {
		t.Errorf("inputs[report_paths] = %v, want report_paths", inputs["report_paths"])
	}

	if inputs["token"] != "${{ secrets.GITHUB_TOKEN }}" {
		t.Errorf("inputs[token] = %v, want token", inputs["token"])
	}

	if inputs["group_reports"] != true {
		t.Errorf("inputs[group_reports] = %v, want true", inputs["group_reports"])
	}

	if inputs["test_files_prefix"] != "src/" {
		t.Errorf("inputs[test_files_prefix] = %v, want test_files_prefix", inputs["test_files_prefix"])
	}

	if inputs["exclude_sources"] != "/build/,/node_modules/" {
		t.Errorf("inputs[exclude_sources] = %v, want exclude_sources", inputs["exclude_sources"])
	}

	if inputs["check_name"] != "Custom Test Report" {
		t.Errorf("inputs[check_name] = %v, want check_name", inputs["check_name"])
	}

	if inputs["commit"] != "${{ github.sha }}" {
		t.Errorf("inputs[commit] = %v, want commit", inputs["commit"])
	}

	if inputs["fail_on_failure"] != true {
		t.Errorf("inputs[fail_on_failure] = %v, want true", inputs["fail_on_failure"])
	}

	if inputs["fail_on_parse_error"] != true {
		t.Errorf("inputs[fail_on_parse_error] = %v, want true", inputs["fail_on_parse_error"])
	}

	if inputs["require_tests"] != true {
		t.Errorf("inputs[require_tests] = %v, want true", inputs["require_tests"])
	}

	if inputs["require_passed_tests"] != true {
		t.Errorf("inputs[require_passed_tests] = %v, want true", inputs["require_passed_tests"])
	}

	if inputs["include_passed"] != true {
		t.Errorf("inputs[include_passed] = %v, want true", inputs["include_passed"])
	}

	if inputs["include_skipped"] != true {
		t.Errorf("inputs[include_skipped] = %v, want true", inputs["include_skipped"])
	}

	if inputs["check_retries"] != true {
		t.Errorf("inputs[check_retries] = %v, want true", inputs["check_retries"])
	}

	if inputs["check_title_template"] != "{{SUITE_NAME}} - {{TEST_NAME}}" {
		t.Errorf("inputs[check_title_template] = %v, want check_title_template", inputs["check_title_template"])
	}

	if inputs["bread_crumb_delimiter"] != "." {
		t.Errorf("inputs[bread_crumb_delimiter] = %v, want bread_crumb_delimiter", inputs["bread_crumb_delimiter"])
	}

	if inputs["summary"] != "Test execution summary" {
		t.Errorf("inputs[summary] = %v, want summary", inputs["summary"])
	}

	if inputs["check_annotations"] != true {
		t.Errorf("inputs[check_annotations] = %v, want true", inputs["check_annotations"])
	}

	if inputs["update_check"] != true {
		t.Errorf("inputs[update_check] = %v, want true", inputs["update_check"])
	}

	if inputs["annotate_only"] != true {
		t.Errorf("inputs[annotate_only] = %v, want true", inputs["annotate_only"])
	}

	if inputs["transformers"] != "[]" {
		t.Errorf("inputs[transformers] = %v, want transformers", inputs["transformers"])
	}

	if inputs["job_summary"] != true {
		t.Errorf("inputs[job_summary] = %v, want true", inputs["job_summary"])
	}

	if inputs["job_summary_text"] != "Additional summary text" {
		t.Errorf("inputs[job_summary_text] = %v, want job_summary_text", inputs["job_summary_text"])
	}

	if inputs["detailed_summary"] != true {
		t.Errorf("inputs[detailed_summary] = %v, want true", inputs["detailed_summary"])
	}

	if inputs["flaky_summary"] != true {
		t.Errorf("inputs[flaky_summary] = %v, want true", inputs["flaky_summary"])
	}

	if inputs["verbose_summary"] != true {
		t.Errorf("inputs[verbose_summary] = %v, want true", inputs["verbose_summary"])
	}

	if inputs["skip_success_summary"] != true {
		t.Errorf("inputs[skip_success_summary] = %v, want true", inputs["skip_success_summary"])
	}

	if inputs["include_empty_in_summary"] != true {
		t.Errorf("inputs[include_empty_in_summary] = %v, want true", inputs["include_empty_in_summary"])
	}

	if inputs["include_time_in_summary"] != true {
		t.Errorf("inputs[include_time_in_summary] = %v, want true", inputs["include_time_in_summary"])
	}

	if inputs["simplified_summary"] != true {
		t.Errorf("inputs[simplified_summary] = %v, want true", inputs["simplified_summary"])
	}

	if inputs["group_suite"] != true {
		t.Errorf("inputs[group_suite] = %v, want true", inputs["group_suite"])
	}

	if inputs["comment"] != true {
		t.Errorf("inputs[comment] = %v, want true", inputs["comment"])
	}

	if inputs["updateComment"] != true {
		t.Errorf("inputs[updateComment] = %v, want true", inputs["updateComment"])
	}

	if inputs["annotate_notice"] != true {
		t.Errorf("inputs[annotate_notice] = %v, want true", inputs["annotate_notice"])
	}

	if inputs["follow_symlink"] != true {
		t.Errorf("inputs[follow_symlink] = %v, want true", inputs["follow_symlink"])
	}

	if inputs["job_name"] != "test-job" {
		t.Errorf("inputs[job_name] = %v, want job_name", inputs["job_name"])
	}

	if inputs["annotations_limit"] != 10 {
		t.Errorf("inputs[annotations_limit] = %v, want 10", inputs["annotations_limit"])
	}

	if inputs["skip_annotations"] != true {
		t.Errorf("inputs[skip_annotations] = %v, want true", inputs["skip_annotations"])
	}

	if inputs["truncate_stack_traces"] != true {
		t.Errorf("inputs[truncate_stack_traces] = %v, want true", inputs["truncate_stack_traces"])
	}

	if inputs["resolve_ignore_classname"] != true {
		t.Errorf("inputs[resolve_ignore_classname] = %v, want true", inputs["resolve_ignore_classname"])
	}

	if inputs["skip_comment_without_tests"] != true {
		t.Errorf("inputs[skip_comment_without_tests] = %v, want true", inputs["skip_comment_without_tests"])
	}

	if inputs["pr_id"] != 123 {
		t.Errorf("inputs[pr_id] = %v, want 123", inputs["pr_id"])
	}
}

func TestJUnitReport_Inputs_MinimalConfig(t *testing.T) {
	jr := JUnitReport{
		ReportPaths: "**/test-results/*.xml",
	}

	inputs := jr.Inputs()

	if inputs["report_paths"] != "**/test-results/*.xml" {
		t.Errorf("inputs[report_paths] = %v, want report_paths", inputs["report_paths"])
	}

	if len(inputs) != 1 {
		t.Errorf("minimal JUnitReport should have 1 input entry, got %d", len(inputs))
	}
}

func TestJUnitReport_Inputs_BooleanFields(t *testing.T) {
	jr := JUnitReport{
		FailOnFailure:      true,
		RequireTests:       true,
		IncludePassed:      true,
		CheckRetries:       true,
		CheckAnnotations:   true,
		UpdateCheck:        true,
		AnnotateOnly:       true,
		JobSummary:         true,
		DetailedSummary:    true,
		FlakySummary:       true,
		VerboseSummary:     true,
		SkipSuccessSummary: true,
	}

	inputs := jr.Inputs()

	if inputs["fail_on_failure"] != true {
		t.Errorf("inputs[fail_on_failure] = %v, want true", inputs["fail_on_failure"])
	}

	if inputs["require_tests"] != true {
		t.Errorf("inputs[require_tests] = %v, want true", inputs["require_tests"])
	}

	if inputs["include_passed"] != true {
		t.Errorf("inputs[include_passed] = %v, want true", inputs["include_passed"])
	}

	if inputs["check_retries"] != true {
		t.Errorf("inputs[check_retries] = %v, want true", inputs["check_retries"])
	}

	if inputs["check_annotations"] != true {
		t.Errorf("inputs[check_annotations] = %v, want true", inputs["check_annotations"])
	}

	if inputs["update_check"] != true {
		t.Errorf("inputs[update_check] = %v, want true", inputs["update_check"])
	}

	if inputs["annotate_only"] != true {
		t.Errorf("inputs[annotate_only] = %v, want true", inputs["annotate_only"])
	}

	if inputs["job_summary"] != true {
		t.Errorf("inputs[job_summary] = %v, want true", inputs["job_summary"])
	}

	if inputs["detailed_summary"] != true {
		t.Errorf("inputs[detailed_summary] = %v, want true", inputs["detailed_summary"])
	}

	if inputs["flaky_summary"] != true {
		t.Errorf("inputs[flaky_summary] = %v, want true", inputs["flaky_summary"])
	}

	if inputs["verbose_summary"] != true {
		t.Errorf("inputs[verbose_summary] = %v, want true", inputs["verbose_summary"])
	}

	if inputs["skip_success_summary"] != true {
		t.Errorf("inputs[skip_success_summary] = %v, want true", inputs["skip_success_summary"])
	}
}

func TestJUnitReport_Inputs_NumericFields(t *testing.T) {
	jr := JUnitReport{
		AnnotationsLimit: 50,
		PrID:             456,
	}

	inputs := jr.Inputs()

	if inputs["annotations_limit"] != 50 {
		t.Errorf("inputs[annotations_limit] = %v, want 50", inputs["annotations_limit"])
	}

	if inputs["pr_id"] != 456 {
		t.Errorf("inputs[pr_id] = %v, want 456", inputs["pr_id"])
	}
}

func TestJUnitReport_ImplementsStepAction(t *testing.T) {
	var _ workflow.StepAction = JUnitReport{}
}
