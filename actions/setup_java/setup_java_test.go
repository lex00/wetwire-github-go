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

func TestSetupJava_Inputs_JavaVersionFile(t *testing.T) {
	a := SetupJava{
		Distribution:    "temurin",
		JavaVersionFile: ".java-version",
	}

	inputs := a.Inputs()

	if inputs["java-version-file"] != ".java-version" {
		t.Errorf("inputs[java-version-file] = %v, want .java-version", inputs["java-version-file"])
	}
}

func TestSetupJava_Inputs_JdkFile(t *testing.T) {
	a := SetupJava{
		Distribution: "temurin",
		JdkFile:      "/path/to/jdk.tar.gz",
	}

	inputs := a.Inputs()

	if inputs["jdk-file"] != "/path/to/jdk.tar.gz" {
		t.Errorf("inputs[jdk-file] = %v, want /path/to/jdk.tar.gz", inputs["jdk-file"])
	}
}

func TestSetupJava_Inputs_SettingsPath(t *testing.T) {
	a := SetupJava{
		JavaVersion:  "17",
		Distribution: "temurin",
		SettingsPath: ".m2/settings.xml",
	}

	inputs := a.Inputs()

	if inputs["settings-path"] != ".m2/settings.xml" {
		t.Errorf("inputs[settings-path] = %v, want .m2/settings.xml", inputs["settings-path"])
	}
}

func TestSetupJava_Inputs_GPG(t *testing.T) {
	a := SetupJava{
		JavaVersion:   "17",
		Distribution:  "temurin",
		GPGPrivateKey: "${{ secrets.GPG_PRIVATE_KEY }}",
		GPGPassphrase: "GPG_PASSPHRASE",
	}

	inputs := a.Inputs()

	if inputs["gpg-private-key"] != "${{ secrets.GPG_PRIVATE_KEY }}" {
		t.Errorf("inputs[gpg-private-key] = %v, want secret reference", inputs["gpg-private-key"])
	}
	if inputs["gpg-passphrase"] != "GPG_PASSPHRASE" {
		t.Errorf("inputs[gpg-passphrase] = %v, want GPG_PASSPHRASE", inputs["gpg-passphrase"])
	}
}

func TestSetupJava_Inputs_CacheDependencyPath(t *testing.T) {
	a := SetupJava{
		JavaVersion:         "17",
		Distribution:        "temurin",
		Cache:               "gradle",
		CacheDependencyPath: "build.gradle",
	}

	inputs := a.Inputs()

	if inputs["cache-dependency-path"] != "build.gradle" {
		t.Errorf("inputs[cache-dependency-path] = %v, want build.gradle", inputs["cache-dependency-path"])
	}
}

func TestSetupJava_Inputs_Token(t *testing.T) {
	a := SetupJava{
		JavaVersion:  "17",
		Distribution: "temurin",
		Token:        "${{ secrets.GITHUB_TOKEN }}",
	}

	inputs := a.Inputs()

	if inputs["token"] != "${{ secrets.GITHUB_TOKEN }}" {
		t.Errorf("inputs[token] = %v, want secret reference", inputs["token"])
	}
}

func TestSetupJava_Inputs_MvnToolchain(t *testing.T) {
	a := SetupJava{
		JavaVersion:        "17",
		Distribution:       "temurin",
		MvnToolchainID:     "my-toolchain",
		MvnToolchainVendor: "openjdk",
	}

	inputs := a.Inputs()

	if inputs["mvn-toolchain-id"] != "my-toolchain" {
		t.Errorf("inputs[mvn-toolchain-id] = %v, want my-toolchain", inputs["mvn-toolchain-id"])
	}
	if inputs["mvn-toolchain-vendor"] != "openjdk" {
		t.Errorf("inputs[mvn-toolchain-vendor] = %v, want openjdk", inputs["mvn-toolchain-vendor"])
	}
}

func TestSetupJava_ImplementsStepAction(t *testing.T) {
	a := SetupJava{}
	// Verify SetupJava implements StepAction interface
	var _ workflow.StepAction = a
}
