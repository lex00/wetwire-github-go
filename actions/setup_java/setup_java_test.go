package setup_java

import (
	"testing"
)

func TestSetupJava_Action(t *testing.T) {
	a := SetupJava{}
	if got := a.Action(); got != "actions/setup-java@v4" {
		t.Errorf("Action() = %q, want %q", got, "actions/setup-java@v4")
	}
}

func TestSetupJava_ToStep(t *testing.T) {
	a := SetupJava{
		JavaVersion:  "21",
		Distribution: "temurin",
	}

	step := a.ToStep()

	if step.Uses != "actions/setup-java@v4" {
		t.Errorf("Uses = %q, want %q", step.Uses, "actions/setup-java@v4")
	}

	if step.With["java-version"] != "21" {
		t.Errorf("With[java-version] = %v, want %q", step.With["java-version"], "21")
	}

	if step.With["distribution"] != "temurin" {
		t.Errorf("With[distribution] = %v, want %q", step.With["distribution"], "temurin")
	}
}

func TestSetupJava_ToStep_Empty(t *testing.T) {
	a := SetupJava{}
	step := a.ToStep()

	if step.Uses != "actions/setup-java@v4" {
		t.Errorf("Uses = %q, want %q", step.Uses, "actions/setup-java@v4")
	}

	// Empty values should not be in With
	if _, ok := step.With["java-version"]; ok {
		t.Error("Empty java-version should not be in With")
	}
}

func TestSetupJava_ToStep_AllFields(t *testing.T) {
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

	step := a.ToStep()

	if step.With["java-version"] != "17" {
		t.Errorf("java-version = %v, want %q", step.With["java-version"], "17")
	}
	if step.With["distribution"] != "zulu" {
		t.Errorf("distribution = %v, want %q", step.With["distribution"], "zulu")
	}
	if step.With["java-package"] != "jdk" {
		t.Errorf("java-package = %v, want %q", step.With["java-package"], "jdk")
	}
	if step.With["check-latest"] != true {
		t.Errorf("check-latest = %v, want %v", step.With["check-latest"], true)
	}
	if step.With["cache"] != "maven" {
		t.Errorf("cache = %v, want %q", step.With["cache"], "maven")
	}
}

func TestSetupJava_ToStep_Distributions(t *testing.T) {
	distributions := []string{"temurin", "zulu", "adopt", "liberica", "microsoft", "corretto", "semeru", "oracle"}

	for _, dist := range distributions {
		t.Run(dist, func(t *testing.T) {
			a := SetupJava{
				JavaVersion:  "21",
				Distribution: dist,
			}

			step := a.ToStep()
			if step.With["distribution"] != dist {
				t.Errorf("distribution = %v, want %q", step.With["distribution"], dist)
			}
		})
	}
}
