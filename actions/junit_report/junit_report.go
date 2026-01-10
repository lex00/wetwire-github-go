// Package junit_report provides a typed wrapper for mikepenz/action-junit-report.
package junit_report

// JUnitReport wraps the mikepenz/action-junit-report@v4 action.
// Publish JUnit test results as GitHub checks and PR comments.
type JUnitReport struct {
	// Glob pattern for JUnit report file locations.
	ReportPaths string `yaml:"report_paths,omitempty"`

	// GitHub token for check creation.
	Token string `yaml:"token,omitempty"`

	// Group multiple reports together.
	GroupReports bool `yaml:"group_reports,omitempty"`

	// Prepend prefix to test file paths.
	TestFilesPrefix string `yaml:"test_files_prefix,omitempty"`

	// Comma-separated folders to ignore during source lookup.
	ExcludeSources string `yaml:"exclude_sources,omitempty"`

	// Name for the check run.
	CheckName string `yaml:"check_name,omitempty"`

	// Commit SHA for status updates.
	Commit string `yaml:"commit,omitempty"`

	// Fail build if tests fail.
	FailOnFailure bool `yaml:"fail_on_failure,omitempty"`

	// Fail if report cannot be parsed.
	FailOnParseError bool `yaml:"fail_on_parse_error,omitempty"`

	// Fail if no tests found.
	RequireTests bool `yaml:"require_tests,omitempty"`

	// Fail if no passed tests detected.
	RequirePassedTests bool `yaml:"require_passed_tests,omitempty"`

	// Include passing tests in annotations.
	IncludePassed bool `yaml:"include_passed,omitempty"`

	// Include skipped tests in summary.
	IncludeSkipped bool `yaml:"include_skipped,omitempty"`

	// Ignore original failures when retried.
	CheckRetries bool `yaml:"check_retries,omitempty"`

	// Custom format template for titles.
	CheckTitleTemplate string `yaml:"check_title_template,omitempty"`

	// Breadcrumb separator character.
	BreadCrumbDelimiter string `yaml:"bread_crumb_delimiter,omitempty"`

	// Additional text for summary output.
	Summary string `yaml:"summary,omitempty"`

	// Enable annotations in checks.
	CheckAnnotations bool `yaml:"check_annotations,omitempty"`

	// Use alternative API for 50+ annotations.
	UpdateCheck bool `yaml:"update_check,omitempty"`

	// Only annotate; skip check creation.
	AnnotateOnly bool `yaml:"annotate_only,omitempty"`

	// Custom filename transformers.
	Transformers string `yaml:"transformers,omitempty"`

	// Publish job summary results.
	JobSummary bool `yaml:"job_summary,omitempty"`

	// Additional job summary text.
	JobSummaryText string `yaml:"job_summary_text,omitempty"`

	// Include detailed test results table.
	DetailedSummary bool `yaml:"detailed_summary,omitempty"`

	// Include flaky results table.
	FlakySummary bool `yaml:"flaky_summary,omitempty"`

	// Note missing annotations in summary.
	VerboseSummary bool `yaml:"verbose_summary,omitempty"`

	// Skip summary if all tests pass.
	SkipSuccessSummary bool `yaml:"skip_success_summary,omitempty"`

	// Include zero-count entries.
	IncludeEmptyInSummary bool `yaml:"include_empty_in_summary,omitempty"`

	// Include test execution time.
	IncludeTimeInSummary bool `yaml:"include_time_in_summary,omitempty"`

	// Use icons instead of text.
	SimplifiedSummary bool `yaml:"simplified_summary,omitempty"`

	// Group test cases by suite.
	GroupSuite bool `yaml:"group_suite,omitempty"`

	// Add PR comment with summary.
	Comment bool `yaml:"comment,omitempty"`

	// Update existing comments.
	UpdateComment bool `yaml:"updateComment,omitempty"`

	// Annotate passed tests too.
	AnnotateNotice bool `yaml:"annotate_notice,omitempty"`

	// Follow symlinks in file search.
	FollowSymlink bool `yaml:"follow_symlink,omitempty"`

	// Check name to update.
	JobName string `yaml:"job_name,omitempty"`

	// Maximum annotation count.
	AnnotationsLimit int `yaml:"annotations_limit,omitempty"`

	// Disable all annotations.
	SkipAnnotations bool `yaml:"skip_annotations,omitempty"`

	// Limit stack trace to 2 lines.
	TruncateStackTraces bool `yaml:"truncate_stack_traces,omitempty"`

	// Ignore test case classname.
	ResolveIgnoreClassname bool `yaml:"resolve_ignore_classname,omitempty"`

	// Skip commenting if no tests.
	SkipCommentWithoutTests bool `yaml:"skip_comment_without_tests,omitempty"`

	// PR number for commenting.
	PrID int `yaml:"pr_id,omitempty"`
}

// Action returns the action reference.
func (a JUnitReport) Action() string {
	return "mikepenz/action-junit-report@v4"
}

// Inputs returns the action inputs as a map.
func (a JUnitReport) Inputs() map[string]any {
	with := make(map[string]any)

	if a.ReportPaths != "" {
		with["report_paths"] = a.ReportPaths
	}
	if a.Token != "" {
		with["token"] = a.Token
	}
	if a.GroupReports {
		with["group_reports"] = a.GroupReports
	}
	if a.TestFilesPrefix != "" {
		with["test_files_prefix"] = a.TestFilesPrefix
	}
	if a.ExcludeSources != "" {
		with["exclude_sources"] = a.ExcludeSources
	}
	if a.CheckName != "" {
		with["check_name"] = a.CheckName
	}
	if a.Commit != "" {
		with["commit"] = a.Commit
	}
	if a.FailOnFailure {
		with["fail_on_failure"] = a.FailOnFailure
	}
	if a.FailOnParseError {
		with["fail_on_parse_error"] = a.FailOnParseError
	}
	if a.RequireTests {
		with["require_tests"] = a.RequireTests
	}
	if a.RequirePassedTests {
		with["require_passed_tests"] = a.RequirePassedTests
	}
	if a.IncludePassed {
		with["include_passed"] = a.IncludePassed
	}
	if a.IncludeSkipped {
		with["include_skipped"] = a.IncludeSkipped
	}
	if a.CheckRetries {
		with["check_retries"] = a.CheckRetries
	}
	if a.CheckTitleTemplate != "" {
		with["check_title_template"] = a.CheckTitleTemplate
	}
	if a.BreadCrumbDelimiter != "" {
		with["bread_crumb_delimiter"] = a.BreadCrumbDelimiter
	}
	if a.Summary != "" {
		with["summary"] = a.Summary
	}
	if a.CheckAnnotations {
		with["check_annotations"] = a.CheckAnnotations
	}
	if a.UpdateCheck {
		with["update_check"] = a.UpdateCheck
	}
	if a.AnnotateOnly {
		with["annotate_only"] = a.AnnotateOnly
	}
	if a.Transformers != "" {
		with["transformers"] = a.Transformers
	}
	if a.JobSummary {
		with["job_summary"] = a.JobSummary
	}
	if a.JobSummaryText != "" {
		with["job_summary_text"] = a.JobSummaryText
	}
	if a.DetailedSummary {
		with["detailed_summary"] = a.DetailedSummary
	}
	if a.FlakySummary {
		with["flaky_summary"] = a.FlakySummary
	}
	if a.VerboseSummary {
		with["verbose_summary"] = a.VerboseSummary
	}
	if a.SkipSuccessSummary {
		with["skip_success_summary"] = a.SkipSuccessSummary
	}
	if a.IncludeEmptyInSummary {
		with["include_empty_in_summary"] = a.IncludeEmptyInSummary
	}
	if a.IncludeTimeInSummary {
		with["include_time_in_summary"] = a.IncludeTimeInSummary
	}
	if a.SimplifiedSummary {
		with["simplified_summary"] = a.SimplifiedSummary
	}
	if a.GroupSuite {
		with["group_suite"] = a.GroupSuite
	}
	if a.Comment {
		with["comment"] = a.Comment
	}
	if a.UpdateComment {
		with["updateComment"] = a.UpdateComment
	}
	if a.AnnotateNotice {
		with["annotate_notice"] = a.AnnotateNotice
	}
	if a.FollowSymlink {
		with["follow_symlink"] = a.FollowSymlink
	}
	if a.JobName != "" {
		with["job_name"] = a.JobName
	}
	if a.AnnotationsLimit != 0 {
		with["annotations_limit"] = a.AnnotationsLimit
	}
	if a.SkipAnnotations {
		with["skip_annotations"] = a.SkipAnnotations
	}
	if a.TruncateStackTraces {
		with["truncate_stack_traces"] = a.TruncateStackTraces
	}
	if a.ResolveIgnoreClassname {
		with["resolve_ignore_classname"] = a.ResolveIgnoreClassname
	}
	if a.SkipCommentWithoutTests {
		with["skip_comment_without_tests"] = a.SkipCommentWithoutTests
	}
	if a.PrID != 0 {
		with["pr_id"] = a.PrID
	}

	return with
}
