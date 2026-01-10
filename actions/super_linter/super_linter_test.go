package super_linter

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestSuperLinter_Action(t *testing.T) {
	a := SuperLinter{}
	if got := a.Action(); got != "super-linter/super-linter@v7" {
		t.Errorf("Action() = %q, want %q", got, "super-linter/super-linter@v7")
	}
}

func TestSuperLinter_Inputs(t *testing.T) {
	a := SuperLinter{
		ValidateAllCodebase: true,
		DefaultBranch:       "main",
		GithubToken:         "${{ secrets.GITHUB_TOKEN }}",
	}

	inputs := a.Inputs()

	if inputs["validate_all_codebase"] != true {
		t.Errorf("inputs[validate_all_codebase] = %v, want true", inputs["validate_all_codebase"])
	}
	if inputs["default_branch"] != "main" {
		t.Errorf("inputs[default_branch] = %v, want main", inputs["default_branch"])
	}
	if inputs["github_token"] != "${{ secrets.GITHUB_TOKEN }}" {
		t.Errorf("inputs[github_token] = %v, want secret reference", inputs["github_token"])
	}
}

func TestSuperLinter_Inputs_Empty(t *testing.T) {
	a := SuperLinter{}
	inputs := a.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty SuperLinter.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestSuperLinter_Inputs_FilterRegex(t *testing.T) {
	a := SuperLinter{
		FilterRegexExclude: ".*_test\\.go$",
		FilterRegexInclude: "src/.*",
	}

	inputs := a.Inputs()

	if inputs["filter_regex_exclude"] != ".*_test\\.go$" {
		t.Errorf("inputs[filter_regex_exclude] = %v, want .*_test\\.go$", inputs["filter_regex_exclude"])
	}
	if inputs["filter_regex_include"] != "src/.*" {
		t.Errorf("inputs[filter_regex_include] = %v, want src/.*", inputs["filter_regex_include"])
	}
}

func TestSuperLinter_Inputs_LogLevel(t *testing.T) {
	a := SuperLinter{
		LogLevel: "DEBUG",
	}

	inputs := a.Inputs()

	if inputs["log_level"] != "DEBUG" {
		t.Errorf("inputs[log_level] = %v, want DEBUG", inputs["log_level"])
	}
}

func TestSuperLinter_Inputs_OutputFormat(t *testing.T) {
	a := SuperLinter{
		OutputFormat:  "tap",
		OutputDetails: "detailed",
	}

	inputs := a.Inputs()

	if inputs["output_format"] != "tap" {
		t.Errorf("inputs[output_format] = %v, want tap", inputs["output_format"])
	}
	if inputs["output_details"] != "detailed" {
		t.Errorf("inputs[output_details] = %v, want detailed", inputs["output_details"])
	}
}

func TestSuperLinter_Inputs_ValidateLanguages(t *testing.T) {
	a := SuperLinter{
		ValidateGo:         true,
		ValidateJavascript: true,
		ValidateTypescript: true,
		ValidatePython:     true,
	}

	inputs := a.Inputs()

	if inputs["validate_go"] != true {
		t.Errorf("inputs[validate_go] = %v, want true", inputs["validate_go"])
	}
	if inputs["validate_javascript"] != true {
		t.Errorf("inputs[validate_javascript] = %v, want true", inputs["validate_javascript"])
	}
	if inputs["validate_typescript"] != true {
		t.Errorf("inputs[validate_typescript] = %v, want true", inputs["validate_typescript"])
	}
	if inputs["validate_python"] != true {
		t.Errorf("inputs[validate_python] = %v, want true", inputs["validate_python"])
	}
}

func TestSuperLinter_Inputs_ValidateConfigFiles(t *testing.T) {
	a := SuperLinter{
		ValidateYaml:       true,
		ValidateJson:       true,
		ValidateMarkdown:   true,
		ValidateDockerfile: true,
		ValidateBash:       true,
	}

	inputs := a.Inputs()

	if inputs["validate_yaml"] != true {
		t.Errorf("inputs[validate_yaml] = %v, want true", inputs["validate_yaml"])
	}
	if inputs["validate_json"] != true {
		t.Errorf("inputs[validate_json] = %v, want true", inputs["validate_json"])
	}
	if inputs["validate_markdown"] != true {
		t.Errorf("inputs[validate_markdown] = %v, want true", inputs["validate_markdown"])
	}
	if inputs["validate_dockerfile"] != true {
		t.Errorf("inputs[validate_dockerfile] = %v, want true", inputs["validate_dockerfile"])
	}
	if inputs["validate_bash"] != true {
		t.Errorf("inputs[validate_bash] = %v, want true", inputs["validate_bash"])
	}
}

func TestSuperLinter_Inputs_Workspace(t *testing.T) {
	a := SuperLinter{
		DefaultWorkspace: "/github/workspace",
		LinterRulesPath:  ".github/linters",
	}

	inputs := a.Inputs()

	if inputs["default_workspace"] != "/github/workspace" {
		t.Errorf("inputs[default_workspace] = %v, want /github/workspace", inputs["default_workspace"])
	}
	if inputs["linter_rules_path"] != ".github/linters" {
		t.Errorf("inputs[linter_rules_path] = %v, want .github/linters", inputs["linter_rules_path"])
	}
}

func TestSuperLinter_Inputs_BooleanFalse(t *testing.T) {
	a := SuperLinter{
		ValidateAllCodebase: false,
		ValidateGo:          false,
		ValidateJavascript:  false,
		ValidateTypescript:  false,
		ValidatePython:      false,
		ValidateYaml:        false,
		ValidateJson:        false,
		ValidateMarkdown:    false,
		ValidateDockerfile:  false,
		ValidateBash:        false,
	}

	inputs := a.Inputs()

	boolFields := []string{
		"validate_all_codebase",
		"validate_go",
		"validate_javascript",
		"validate_typescript",
		"validate_python",
		"validate_yaml",
		"validate_json",
		"validate_markdown",
		"validate_dockerfile",
		"validate_bash",
	}

	for _, field := range boolFields {
		if _, ok := inputs[field]; ok {
			t.Errorf("%s=false should not be in inputs", field)
		}
	}
}

func TestSuperLinter_Inputs_AllFields(t *testing.T) {
	a := SuperLinter{
		ValidateAllCodebase: true,
		DefaultBranch:       "develop",
		GithubToken:         "${{ secrets.GITHUB_TOKEN }}",
		FilterRegexExclude:  "vendor/.*",
		FilterRegexInclude:  "src/.*",
		LogLevel:            "INFO",
		OutputFormat:        "tap",
		OutputDetails:       "simpler",
		ValidateGo:          true,
		ValidateJavascript:  true,
		ValidateTypescript:  true,
		ValidatePython:      true,
		ValidateYaml:        true,
		ValidateJson:        true,
		ValidateMarkdown:    true,
		ValidateDockerfile:  true,
		ValidateBash:        true,
		DefaultWorkspace:    "/workspace",
		LinterRulesPath:     ".linters",
	}

	inputs := a.Inputs()

	expected := map[string]any{
		"validate_all_codebase": true,
		"default_branch":        "develop",
		"github_token":          "${{ secrets.GITHUB_TOKEN }}",
		"filter_regex_exclude":  "vendor/.*",
		"filter_regex_include":  "src/.*",
		"log_level":             "INFO",
		"output_format":         "tap",
		"output_details":        "simpler",
		"validate_go":           true,
		"validate_javascript":   true,
		"validate_typescript":   true,
		"validate_python":       true,
		"validate_yaml":         true,
		"validate_json":         true,
		"validate_markdown":     true,
		"validate_dockerfile":   true,
		"validate_bash":         true,
		"default_workspace":     "/workspace",
		"linter_rules_path":     ".linters",
	}

	if len(inputs) != len(expected) {
		t.Errorf("inputs has %d entries, want %d", len(inputs), len(expected))
	}

	for key, want := range expected {
		if got := inputs[key]; got != want {
			t.Errorf("inputs[%q] = %v, want %v", key, got, want)
		}
	}
}

func TestSuperLinter_ImplementsStepAction(t *testing.T) {
	a := SuperLinter{}
	var _ workflow.StepAction = a
}
