package setup_java

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestSetupJava_Action(t *testing.T) {
	a := SetupJava{}
	if got := a.Action(); got != "actions/setup-java@v4" {
		t.Errorf("Action() = %q, want %q", got, "actions/setup-java@v4")
	}
}

func TestSetupJava_Inputs(t *testing.T) {
	a := SetupJava{
		JavaVersion:  "21",
		Distribution: "temurin",
	}

	inputs := a.Inputs()

	if a.Action() != "actions/setup-java@v4" {
		t.Errorf("Action() = %q, want %q", a.Action(), "actions/setup-java@v4")
	}

	if inputs["java-version"] != "21" {
		t.Errorf("inputs[java-version] = %v, want %q", inputs["java-version"], "21")
	}

	if inputs["distribution"] != "temurin" {
		t.Errorf("inputs[distribution] = %v, want %q", inputs["distribution"], "temurin")
	}
}

func TestSetupJava_Inputs_Empty(t *testing.T) {
	a := SetupJava{}
	inputs := a.Inputs()

	if a.Action() != "actions/setup-java@v4" {
		t.Errorf("Action() = %q, want %q", a.Action(), "actions/setup-java@v4")
	}

	// Empty values should not be in inputs
	if _, ok := inputs["java-version"]; ok {
		t.Error("Empty java-version should not be in inputs")
	}
}

func TestSetupJava_Inputs_AllFields(t *testing.T) {
	a := SetupJava{
		JavaVersion:       "17",
		Distribution:      "zulu",
		JavaPackage:       "jdk",
		Architecture:      "x64",
		CheckLatest:       true,
		Cache:             "maven",
		ServerID:          "github",
		ServerUsername:    "MAVEN_USERNAME",
		ServerPassword:    "MAVEN_PASSWORD",
		OverwriteSettings: true,
	}

	inputs := a.Inputs()

	if inputs["java-version"] != "17" {
		t.Errorf("java-version = %v, want %q", inputs["java-version"], "17")
	}
	if inputs["distribution"] != "zulu" {
		t.Errorf("distribution = %v, want %q", inputs["distribution"], "zulu")
	}
	if inputs["java-package"] != "jdk" {
		t.Errorf("java-package = %v, want %q", inputs["java-package"], "jdk")
	}
	if inputs["check-latest"] != true {
		t.Errorf("check-latest = %v, want %v", inputs["check-latest"], true)
	}
	if inputs["cache"] != "maven" {
		t.Errorf("cache = %v, want %q", inputs["cache"], "maven")
	}
}

func TestSetupJava_Inputs_Distributions(t *testing.T) {
	distributions := []string{"temurin", "zulu", "adopt", "liberica", "microsoft", "corretto", "semeru", "oracle"}

	for _, dist := range distributions {
		t.Run(dist, func(t *testing.T) {
			a := SetupJava{
				JavaVersion:  "21",
				Distribution: dist,
			}

			inputs := a.Inputs()
			if inputs["distribution"] != dist {
				t.Errorf("distribution = %v, want %q", inputs["distribution"], dist)
			}
		})
	}
}

func TestSetupJava_ImplementsStepAction(t *testing.T) {
	a := SetupJava{}
	// Verify SetupJava implements StepAction interface
	var _ workflow.StepAction = a
}
