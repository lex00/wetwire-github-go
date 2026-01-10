package setup_dotnet

import (
	"testing"
)

func TestSetupDotnet_Action(t *testing.T) {
	a := SetupDotnet{}
	if got := a.Action(); got != "actions/setup-dotnet@v4" {
		t.Errorf("Action() = %q, want %q", got, "actions/setup-dotnet@v4")
	}
}

func TestSetupDotnet_ToStep(t *testing.T) {
	a := SetupDotnet{
		DotnetVersion: "8.0.x",
	}

	step := a.ToStep()

	if step.Uses != "actions/setup-dotnet@v4" {
		t.Errorf("Uses = %q, want %q", step.Uses, "actions/setup-dotnet@v4")
	}

	if step.With["dotnet-version"] != "8.0.x" {
		t.Errorf("With[dotnet-version] = %v, want %q", step.With["dotnet-version"], "8.0.x")
	}
}

func TestSetupDotnet_ToStep_Empty(t *testing.T) {
	a := SetupDotnet{}
	step := a.ToStep()

	if step.Uses != "actions/setup-dotnet@v4" {
		t.Errorf("Uses = %q, want %q", step.Uses, "actions/setup-dotnet@v4")
	}

	if _, ok := step.With["dotnet-version"]; ok {
		t.Error("Empty dotnet-version should not be in With")
	}
}

func TestSetupDotnet_ToStep_AllFields(t *testing.T) {
	a := SetupDotnet{
		DotnetVersion:       "8.0.x",
		DotnetQuality:       "ga",
		GlobalJsonFile:      "./global.json",
		IncludePrerelease:   true,
		Cache:               true,
		CacheDependencyPath: "**/packages.lock.json",
	}

	step := a.ToStep()

	if step.With["dotnet-version"] != "8.0.x" {
		t.Errorf("dotnet-version = %v, want %q", step.With["dotnet-version"], "8.0.x")
	}
	if step.With["dotnet-quality"] != "ga" {
		t.Errorf("dotnet-quality = %v, want %q", step.With["dotnet-quality"], "ga")
	}
	if step.With["global-json-file"] != "./global.json" {
		t.Errorf("global-json-file = %v, want %q", step.With["global-json-file"], "./global.json")
	}
	if step.With["include-prerelease"] != true {
		t.Errorf("include-prerelease = %v, want %v", step.With["include-prerelease"], true)
	}
	if step.With["cache"] != true {
		t.Errorf("cache = %v, want %v", step.With["cache"], true)
	}
}

func TestSetupDotnet_ToStep_Versions(t *testing.T) {
	versions := []string{"6.0.x", "7.0.x", "8.0.x", "8.0.100", "9.0.x"}

	for _, ver := range versions {
		t.Run(ver, func(t *testing.T) {
			a := SetupDotnet{
				DotnetVersion: ver,
			}

			step := a.ToStep()
			if step.With["dotnet-version"] != ver {
				t.Errorf("dotnet-version = %v, want %q", step.With["dotnet-version"], ver)
			}
		})
	}
}

func TestSetupDotnet_ToStep_Quality(t *testing.T) {
	qualities := []string{"daily", "signed", "validated", "preview", "ga"}

	for _, q := range qualities {
		t.Run(q, func(t *testing.T) {
			a := SetupDotnet{
				DotnetVersion: "8.0.x",
				DotnetQuality: q,
			}

			step := a.ToStep()
			if step.With["dotnet-quality"] != q {
				t.Errorf("dotnet-quality = %v, want %q", step.With["dotnet-quality"], q)
			}
		})
	}
}
