package setup_dotnet

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestSetupDotnet_Action(t *testing.T) {
	a := SetupDotnet{}
	if got := a.Action(); got != "actions/setup-dotnet@v4" {
		t.Errorf("Action() = %q, want %q", got, "actions/setup-dotnet@v4")
	}
}

func TestSetupDotnet_Inputs(t *testing.T) {
	a := SetupDotnet{
		DotnetVersion: "8.0.x",
	}

	inputs := a.Inputs()

	if a.Action() != "actions/setup-dotnet@v4" {
		t.Errorf("Action() = %q, want %q", a.Action(), "actions/setup-dotnet@v4")
	}

	if inputs["dotnet-version"] != "8.0.x" {
		t.Errorf("inputs[dotnet-version] = %v, want %q", inputs["dotnet-version"], "8.0.x")
	}
}

func TestSetupDotnet_Inputs_Empty(t *testing.T) {
	a := SetupDotnet{}
	inputs := a.Inputs()

	if a.Action() != "actions/setup-dotnet@v4" {
		t.Errorf("Action() = %q, want %q", a.Action(), "actions/setup-dotnet@v4")
	}

	if _, ok := inputs["dotnet-version"]; ok {
		t.Error("Empty dotnet-version should not be in inputs")
	}
}

func TestSetupDotnet_Inputs_AllFields(t *testing.T) {
	a := SetupDotnet{
		DotnetVersion:       "8.0.x",
		DotnetQuality:       "ga",
		GlobalJsonFile:      "./global.json",
		IncludePrerelease:   true,
		Cache:               true,
		CacheDependencyPath: "**/packages.lock.json",
	}

	inputs := a.Inputs()

	if inputs["dotnet-version"] != "8.0.x" {
		t.Errorf("dotnet-version = %v, want %q", inputs["dotnet-version"], "8.0.x")
	}
	if inputs["dotnet-quality"] != "ga" {
		t.Errorf("dotnet-quality = %v, want %q", inputs["dotnet-quality"], "ga")
	}
	if inputs["global-json-file"] != "./global.json" {
		t.Errorf("global-json-file = %v, want %q", inputs["global-json-file"], "./global.json")
	}
	if inputs["include-prerelease"] != true {
		t.Errorf("include-prerelease = %v, want %v", inputs["include-prerelease"], true)
	}
	if inputs["cache"] != true {
		t.Errorf("cache = %v, want %v", inputs["cache"], true)
	}
}

func TestSetupDotnet_Inputs_Versions(t *testing.T) {
	versions := []string{"6.0.x", "7.0.x", "8.0.x", "8.0.100", "9.0.x"}

	for _, ver := range versions {
		t.Run(ver, func(t *testing.T) {
			a := SetupDotnet{
				DotnetVersion: ver,
			}

			inputs := a.Inputs()
			if inputs["dotnet-version"] != ver {
				t.Errorf("dotnet-version = %v, want %q", inputs["dotnet-version"], ver)
			}
		})
	}
}

func TestSetupDotnet_Inputs_Quality(t *testing.T) {
	qualities := []string{"daily", "signed", "validated", "preview", "ga"}

	for _, q := range qualities {
		t.Run(q, func(t *testing.T) {
			a := SetupDotnet{
				DotnetVersion: "8.0.x",
				DotnetQuality: q,
			}

			inputs := a.Inputs()
			if inputs["dotnet-quality"] != q {
				t.Errorf("dotnet-quality = %v, want %q", inputs["dotnet-quality"], q)
			}
		})
	}
}

func TestSetupDotnet_ImplementsStepAction(t *testing.T) {
	a := SetupDotnet{}
	// Verify SetupDotnet implements StepAction interface
	var _ workflow.StepAction = a
}
