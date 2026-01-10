package dependabot

import (
	"testing"
)

func TestDependabot_ResourceType(t *testing.T) {
	d := Dependabot{}
	if got := d.ResourceType(); got != "dependabot" {
		t.Errorf("ResourceType() = %q, want %q", got, "dependabot")
	}
}

func TestDependabot_Basic(t *testing.T) {
	d := Dependabot{
		Version: 2,
		Updates: []Update{
			{
				PackageEcosystem: "gomod",
				Directory:        "/",
				Schedule: Schedule{
					Interval: "weekly",
				},
			},
		},
	}

	if d.Version != 2 {
		t.Errorf("Version = %d, want 2", d.Version)
	}

	if len(d.Updates) != 1 {
		t.Errorf("len(Updates) = %d, want 1", len(d.Updates))
	}

	if d.Updates[0].PackageEcosystem != "gomod" {
		t.Errorf("Updates[0].PackageEcosystem = %q, want %q", d.Updates[0].PackageEcosystem, "gomod")
	}
}

func TestDependabot_WithRegistries(t *testing.T) {
	d := Dependabot{
		Version:              2,
		EnableBetaEcosystems: true,
		Registries: map[string]Registry{
			"npm-github": {
				Type:     "npm-registry",
				URL:      "https://npm.pkg.github.com",
				Token:    "${{ secrets.GITHUB_TOKEN }}",
				Username: "x-access-token",
			},
		},
		Updates: []Update{
			{
				PackageEcosystem: "npm",
				Directory:        "/",
				Schedule: Schedule{
					Interval: "daily",
				},
				Registries: "*",
			},
		},
	}

	if !d.EnableBetaEcosystems {
		t.Error("EnableBetaEcosystems should be true")
	}

	if len(d.Registries) != 1 {
		t.Errorf("len(Registries) = %d, want 1", len(d.Registries))
	}

	reg, ok := d.Registries["npm-github"]
	if !ok {
		t.Fatal("Registries[npm-github] not found")
	}

	if reg.Type != "npm-registry" {
		t.Errorf("Registry.Type = %q, want %q", reg.Type, "npm-registry")
	}
}

func TestUpdate_AllFields(t *testing.T) {
	u := Update{
		PackageEcosystem:      "npm",
		Directory:             "/frontend",
		Directories:           []string{"/frontend", "/backend"},
		Schedule:              Schedule{Interval: "weekly", Day: "monday"},
		Allow:                 []Allow{{DependencyName: "lodash"}},
		Ignore:                []Ignore{{DependencyName: "react", Versions: []string{">18.0.0"}}},
		Labels:                []string{"dependencies", "npm"},
		Assignees:             []string{"user1"},
		Reviewers:             []string{"team1"},
		Milestone:             1,
		OpenPullRequestsLimit: 10,
		RebaseStrategy:        "auto",
		VersioningStrategy:    "increase",
		Vendor:                true,
		TargetBranch:          "develop",
		Groups: map[string]Group{
			"dev-deps": {
				Patterns:       []string{"eslint*", "prettier*"},
				DependencyType: "development",
			},
		},
		CommitMessage: &CommitMessage{
			Prefix: "deps",
		},
		PullRequestBranchName: &PullRequestBranchName{
			Separator: "-",
		},
		InsecureExternalCodeExecution: "deny",
	}

	if u.PackageEcosystem != "npm" {
		t.Errorf("PackageEcosystem = %q, want %q", u.PackageEcosystem, "npm")
	}

	if u.OpenPullRequestsLimit != 10 {
		t.Errorf("OpenPullRequestsLimit = %d, want 10", u.OpenPullRequestsLimit)
	}

	if !u.Vendor {
		t.Error("Vendor should be true")
	}

	if u.CommitMessage.Prefix != "deps" {
		t.Errorf("CommitMessage.Prefix = %q, want %q", u.CommitMessage.Prefix, "deps")
	}

	if u.PullRequestBranchName.Separator != "-" {
		t.Errorf("PullRequestBranchName.Separator = %q, want %q", u.PullRequestBranchName.Separator, "-")
	}
}

func TestSchedule_Weekly(t *testing.T) {
	s := Schedule{
		Interval: "weekly",
		Day:      "monday",
		Time:     "09:00",
		Timezone: "America/New_York",
	}

	if s.Interval != "weekly" {
		t.Errorf("Interval = %q, want %q", s.Interval, "weekly")
	}

	if s.Day != "monday" {
		t.Errorf("Day = %q, want %q", s.Day, "monday")
	}

	if s.Time != "09:00" {
		t.Errorf("Time = %q, want %q", s.Time, "09:00")
	}

	if s.Timezone != "America/New_York" {
		t.Errorf("Timezone = %q, want %q", s.Timezone, "America/New_York")
	}
}

func TestSchedule_Daily(t *testing.T) {
	s := Schedule{
		Interval: "daily",
	}

	if s.Interval != "daily" {
		t.Errorf("Interval = %q, want %q", s.Interval, "daily")
	}

	if s.Day != "" {
		t.Errorf("Day should be empty for daily schedule, got %q", s.Day)
	}
}

func TestGroup_AllFields(t *testing.T) {
	g := Group{
		Patterns:        []string{"lodash*", "underscore*"},
		DependencyType:  "production",
		UpdateTypes:     []string{"minor", "patch"},
		ExcludePatterns: []string{"lodash-es"},
		AppliesTo:       "version-updates",
	}

	if len(g.Patterns) != 2 {
		t.Errorf("len(Patterns) = %d, want 2", len(g.Patterns))
	}

	if g.DependencyType != "production" {
		t.Errorf("DependencyType = %q, want %q", g.DependencyType, "production")
	}

	if len(g.UpdateTypes) != 2 {
		t.Errorf("len(UpdateTypes) = %d, want 2", len(g.UpdateTypes))
	}

	if g.AppliesTo != "version-updates" {
		t.Errorf("AppliesTo = %q, want %q", g.AppliesTo, "version-updates")
	}
}

func TestAllow(t *testing.T) {
	a := Allow{
		DependencyName: "lodash",
		DependencyType: "direct",
	}

	if a.DependencyName != "lodash" {
		t.Errorf("DependencyName = %q, want %q", a.DependencyName, "lodash")
	}

	if a.DependencyType != "direct" {
		t.Errorf("DependencyType = %q, want %q", a.DependencyType, "direct")
	}
}

func TestIgnore(t *testing.T) {
	i := Ignore{
		DependencyName: "react",
		Versions:       []string{">18.0.0", "<16.0.0"},
		UpdateTypes:    []string{"version-update:semver-major"},
	}

	if i.DependencyName != "react" {
		t.Errorf("DependencyName = %q, want %q", i.DependencyName, "react")
	}

	if len(i.Versions) != 2 {
		t.Errorf("len(Versions) = %d, want 2", len(i.Versions))
	}

	if len(i.UpdateTypes) != 1 {
		t.Errorf("len(UpdateTypes) = %d, want 1", len(i.UpdateTypes))
	}
}

func TestCommitMessage(t *testing.T) {
	cm := CommitMessage{
		Prefix:            "chore(deps)",
		PrefixDevelopment: "chore(dev-deps)",
		Include:           "scope",
	}

	if cm.Prefix != "chore(deps)" {
		t.Errorf("Prefix = %q, want %q", cm.Prefix, "chore(deps)")
	}

	if cm.PrefixDevelopment != "chore(dev-deps)" {
		t.Errorf("PrefixDevelopment = %q, want %q", cm.PrefixDevelopment, "chore(dev-deps)")
	}

	if cm.Include != "scope" {
		t.Errorf("Include = %q, want %q", cm.Include, "scope")
	}
}

func TestPullRequestBranchName(t *testing.T) {
	prbn := PullRequestBranchName{
		Separator: "_",
	}

	if prbn.Separator != "_" {
		t.Errorf("Separator = %q, want %q", prbn.Separator, "_")
	}
}
