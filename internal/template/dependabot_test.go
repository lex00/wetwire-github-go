package template

import (
	"strings"
	"testing"

	"github.com/lex00/wetwire-github-go/dependabot"
	"github.com/lex00/wetwire-github-go/internal/discover"
	"github.com/lex00/wetwire-github-go/internal/runner"
)

func TestBuilder_BuildDependabot_Empty(t *testing.T) {
	b := NewBuilder()

	discovered := &discover.DependabotDiscoveryResult{
		Configs: []discover.DiscoveredDependabot{},
	}
	extracted := &runner.DependabotExtractionResult{
		Configs: []runner.ExtractedDependabot{},
	}

	result, err := b.BuildDependabot(discovered, extracted)
	if err != nil {
		t.Fatalf("BuildDependabot() error = %v", err)
	}

	if len(result.Configs) != 0 {
		t.Errorf("Expected 0 configs, got %d", len(result.Configs))
	}
}

func TestBuilder_BuildDependabot_SingleConfig(t *testing.T) {
	b := NewBuilder()

	discovered := &discover.DependabotDiscoveryResult{
		Configs: []discover.DiscoveredDependabot{
			{Name: "DependabotConfig", File: "dependabot.go", Line: 10},
		},
	}
	extracted := &runner.DependabotExtractionResult{
		Configs: []runner.ExtractedDependabot{
			{
				Name: "DependabotConfig",
				Data: map[string]any{
					"Version": 2,
					"Updates": []any{
						map[string]any{
							"PackageEcosystem": "gomod",
							"Directory":        "/",
							"Schedule": map[string]any{
								"Interval": "weekly",
							},
						},
					},
				},
			},
		},
	}

	result, err := b.BuildDependabot(discovered, extracted)
	if err != nil {
		t.Fatalf("BuildDependabot() error = %v", err)
	}

	if len(result.Configs) != 1 {
		t.Fatalf("Expected 1 config, got %d", len(result.Configs))
	}

	cfg := result.Configs[0]
	if cfg.Name != "DependabotConfig" {
		t.Errorf("Config name = %q, want %q", cfg.Name, "DependabotConfig")
	}

	if len(cfg.YAML) == 0 {
		t.Error("Config YAML is empty")
	}

	if cfg.Config.Version != 2 {
		t.Errorf("Config.Version = %d, want %d", cfg.Config.Version, 2)
	}

	if len(cfg.Config.Updates) != 1 {
		t.Fatalf("Expected 1 update, got %d", len(cfg.Config.Updates))
	}
}

func TestBuilder_BuildDependabot_MissingExtraction(t *testing.T) {
	b := NewBuilder()

	discovered := &discover.DependabotDiscoveryResult{
		Configs: []discover.DiscoveredDependabot{
			{Name: "MissingConfig", File: "dependabot.go", Line: 10},
		},
	}
	// No extraction data
	extracted := &runner.DependabotExtractionResult{
		Configs: []runner.ExtractedDependabot{},
	}

	result, err := b.BuildDependabot(discovered, extracted)
	if err != nil {
		t.Fatalf("BuildDependabot() error = %v", err)
	}

	// Should have an error about missing extraction data
	if len(result.Errors) == 0 {
		t.Error("Expected error about missing extraction data")
	}

	if len(result.Configs) != 0 {
		t.Errorf("Expected 0 configs when extraction is missing, got %d", len(result.Configs))
	}
}

func TestBuilder_reconstructDependabot(t *testing.T) {
	b := NewBuilder()

	tests := []struct {
		name string
		data map[string]any
		want func(*dependabot.Dependabot) bool
	}{
		{
			name: "version as int",
			data: map[string]any{
				"Version": 2,
			},
			want: func(d *dependabot.Dependabot) bool {
				return d.Version == 2
			},
		},
		{
			name: "version as float64",
			data: map[string]any{
				"Version": 2.0,
			},
			want: func(d *dependabot.Dependabot) bool {
				return d.Version == 2
			},
		},
		{
			name: "enable beta ecosystems",
			data: map[string]any{
				"Version":               2,
				"EnableBetaEcosystems": true,
			},
			want: func(d *dependabot.Dependabot) bool {
				return d.EnableBetaEcosystems
			},
		},
		{
			name: "with updates",
			data: map[string]any{
				"Version": 2,
				"Updates": []any{
					map[string]any{
						"PackageEcosystem": "npm",
						"Directory":        "/",
					},
				},
			},
			want: func(d *dependabot.Dependabot) bool {
				return len(d.Updates) == 1 && d.Updates[0].PackageEcosystem == "npm"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := b.reconstructDependabot(tt.data)
			if result == nil {
				t.Fatal("reconstructDependabot() returned nil")
			}
			if !tt.want(result) {
				t.Errorf("reconstructDependabot() validation failed")
			}
		})
	}
}

func TestBuilder_reconstructUpdate(t *testing.T) {
	b := NewBuilder()

	tests := []struct {
		name string
		data map[string]any
		want func(dependabot.Update) bool
	}{
		{
			name: "basic update",
			data: map[string]any{
				"PackageEcosystem": "gomod",
				"Directory":        "/",
			},
			want: func(u dependabot.Update) bool {
				return u.PackageEcosystem == "gomod" && u.Directory == "/"
			},
		},
		{
			name: "with schedule",
			data: map[string]any{
				"PackageEcosystem": "npm",
				"Directory":        "/frontend",
				"Schedule": map[string]any{
					"Interval": "daily",
					"Time":     "09:00",
				},
			},
			want: func(u dependabot.Update) bool {
				return u.Schedule.Interval == "daily" && u.Schedule.Time == "09:00"
			},
		},
		{
			name: "with labels",
			data: map[string]any{
				"PackageEcosystem": "docker",
				"Directory":        "/",
				"Labels":           []any{"dependencies", "docker"},
			},
			want: func(u dependabot.Update) bool {
				return len(u.Labels) == 2 && u.Labels[0] == "dependencies"
			},
		},
		{
			name: "with milestone int",
			data: map[string]any{
				"PackageEcosystem": "pip",
				"Directory":        "/",
				"Milestone":        5,
			},
			want: func(u dependabot.Update) bool {
				return u.Milestone == 5
			},
		},
		{
			name: "with milestone float64",
			data: map[string]any{
				"PackageEcosystem": "pip",
				"Directory":        "/",
				"Milestone":        5.0,
			},
			want: func(u dependabot.Update) bool {
				return u.Milestone == 5
			},
		},
		{
			name: "with open pull requests limit",
			data: map[string]any{
				"PackageEcosystem":      "maven",
				"Directory":             "/",
				"OpenPullRequestsLimit": 10,
			},
			want: func(u dependabot.Update) bool {
				return u.OpenPullRequestsLimit == 10
			},
		},
		{
			name: "with ignore list",
			data: map[string]any{
				"PackageEcosystem": "npm",
				"Directory":        "/",
				"Ignore": []any{
					map[string]any{
						"DependencyName": "react",
						"Versions":       []any{"17.0.0"},
					},
				},
			},
			want: func(u dependabot.Update) bool {
				return len(u.Ignore) == 1 && u.Ignore[0].DependencyName == "react"
			},
		},
		{
			name: "with allow list",
			data: map[string]any{
				"PackageEcosystem": "npm",
				"Directory":        "/",
				"Allow": []any{
					map[string]any{
						"DependencyName": "lodash",
						"DependencyType": "production",
					},
				},
			},
			want: func(u dependabot.Update) bool {
				return len(u.Allow) == 1 && u.Allow[0].DependencyName == "lodash"
			},
		},
		{
			name: "with vendor boolean",
			data: map[string]any{
				"PackageEcosystem": "gomod",
				"Directory":        "/",
				"Vendor":           true,
			},
			want: func(u dependabot.Update) bool {
				return u.Vendor
			},
		},
		{
			name: "with groups",
			data: map[string]any{
				"PackageEcosystem": "npm",
				"Directory":        "/",
				"Groups": map[string]any{
					"dev-dependencies": map[string]any{
						"DependencyType": "development",
					},
				},
			},
			want: func(u dependabot.Update) bool {
				return len(u.Groups) == 1
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := b.reconstructUpdate(tt.data)
			if !tt.want(result) {
				t.Errorf("reconstructUpdate() validation failed")
			}
		})
	}
}

func TestBuilder_reconstructSchedule(t *testing.T) {
	b := NewBuilder()

	data := map[string]any{
		"Interval": "weekly",
		"Day":      "monday",
		"Time":     "09:00",
		"Timezone": "America/New_York",
	}

	result := b.reconstructSchedule(data)

	if result.Interval != "weekly" {
		t.Errorf("Interval = %q, want %q", result.Interval, "weekly")
	}
	if result.Day != "monday" {
		t.Errorf("Day = %q, want %q", result.Day, "monday")
	}
	if result.Time != "09:00" {
		t.Errorf("Time = %q, want %q", result.Time, "09:00")
	}
	if result.Timezone != "America/New_York" {
		t.Errorf("Timezone = %q, want %q", result.Timezone, "America/New_York")
	}
}

func TestBuilder_reconstructAllowList(t *testing.T) {
	b := NewBuilder()

	data := []any{
		map[string]any{
			"DependencyName": "express",
			"DependencyType": "production",
		},
		map[string]any{
			"DependencyName": "lodash",
		},
	}

	result := b.reconstructAllowList(data)

	if len(result) != 2 {
		t.Fatalf("Expected 2 items, got %d", len(result))
	}

	if result[0].DependencyName != "express" {
		t.Errorf("First item DependencyName = %q, want %q", result[0].DependencyName, "express")
	}
	if result[0].DependencyType != "production" {
		t.Errorf("First item DependencyType = %q, want %q", result[0].DependencyType, "production")
	}
}

func TestBuilder_reconstructIgnoreList(t *testing.T) {
	b := NewBuilder()

	data := []any{
		map[string]any{
			"DependencyName": "react",
			"Versions":       []any{"17.0.0", "17.0.1"},
			"UpdateTypes":    []any{"version-update:semver-major"},
		},
	}

	result := b.reconstructIgnoreList(data)

	if len(result) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(result))
	}

	if result[0].DependencyName != "react" {
		t.Errorf("DependencyName = %q, want %q", result[0].DependencyName, "react")
	}
	if len(result[0].Versions) != 2 {
		t.Errorf("Expected 2 versions, got %d", len(result[0].Versions))
	}
	if len(result[0].UpdateTypes) != 1 {
		t.Errorf("Expected 1 update type, got %d", len(result[0].UpdateTypes))
	}
}

func TestBuilder_reconstructGroups(t *testing.T) {
	b := NewBuilder()

	data := map[string]any{
		"production": map[string]any{
			"Patterns":       []any{"*"},
			"DependencyType": "production",
		},
		"development": map[string]any{
			"Patterns":  []any{"eslint*", "prettier*"},
			"AppliesTo": "version-updates",
		},
	}

	result := b.reconstructGroups(data)

	if len(result) != 2 {
		t.Fatalf("Expected 2 groups, got %d", len(result))
	}

	prodGroup, ok := result["production"]
	if !ok {
		t.Fatal("production group not found")
	}
	if prodGroup.DependencyType != "production" {
		t.Errorf("production DependencyType = %q, want %q", prodGroup.DependencyType, "production")
	}

	devGroup, ok := result["development"]
	if !ok {
		t.Fatal("development group not found")
	}
	if len(devGroup.Patterns) != 2 {
		t.Errorf("Expected 2 patterns, got %d", len(devGroup.Patterns))
	}
}

func TestBuilder_reconstructRegistries(t *testing.T) {
	b := NewBuilder()

	data := map[string]any{
		"npm-registry": map[string]any{
			"Type":     "npm-registry",
			"URL":      "https://npm.pkg.github.com",
			"Token":    "${{secrets.NPM_TOKEN}}",
			"Username": "octocat",
		},
		"docker-registry": map[string]any{
			"Type":     "docker-registry",
			"URL":      "https://docker.pkg.github.com",
			"Username": "octocat",
			"Password": "${{secrets.DOCKER_PASSWORD}}",
		},
	}

	result := b.reconstructRegistries(data)

	if len(result) != 2 {
		t.Fatalf("Expected 2 registries, got %d", len(result))
	}

	npmReg, ok := result["npm-registry"]
	if !ok {
		t.Fatal("npm-registry not found")
	}
	if npmReg.Type != "npm-registry" {
		t.Errorf("npm-registry Type = %q, want %q", npmReg.Type, "npm-registry")
	}
	if npmReg.Token != "${{secrets.NPM_TOKEN}}" {
		t.Errorf("npm-registry Token = %q, want %q", npmReg.Token, "${{secrets.NPM_TOKEN}}")
	}

	dockerReg, ok := result["docker-registry"]
	if !ok {
		t.Fatal("docker-registry not found")
	}
	if dockerReg.Password != "${{secrets.DOCKER_PASSWORD}}" {
		t.Errorf("docker-registry Password = %q, want %q", dockerReg.Password, "${{secrets.DOCKER_PASSWORD}}")
	}
}

func TestBuilder_reconstructCommitMessage(t *testing.T) {
	b := NewBuilder()

	data := map[string]any{
		"Prefix":            "chore",
		"PrefixDevelopment": "build",
		"Include":           "scope",
	}

	result := b.reconstructCommitMessage(data)

	if result.Prefix != "chore" {
		t.Errorf("Prefix = %q, want %q", result.Prefix, "chore")
	}
	if result.PrefixDevelopment != "build" {
		t.Errorf("PrefixDevelopment = %q, want %q", result.PrefixDevelopment, "build")
	}
	if result.Include != "scope" {
		t.Errorf("Include = %q, want %q", result.Include, "scope")
	}
}

func TestBuilder_reconstructPullRequestBranchName(t *testing.T) {
	b := NewBuilder()

	data := map[string]any{
		"Separator": "/",
	}

	result := b.reconstructPullRequestBranchName(data)

	if result.Separator != "/" {
		t.Errorf("Separator = %q, want %q", result.Separator, "/")
	}
}

func TestBuilder_BuildDependabot_ComplexConfig(t *testing.T) {
	b := NewBuilder()

	discovered := &discover.DependabotDiscoveryResult{
		Configs: []discover.DiscoveredDependabot{
			{Name: "FullConfig", File: "dependabot.go", Line: 10},
		},
	}
	extracted := &runner.DependabotExtractionResult{
		Configs: []runner.ExtractedDependabot{
			{
				Name: "FullConfig",
				Data: map[string]any{
					"Version":               2,
					"EnableBetaEcosystems": true,
					"Updates": []any{
						map[string]any{
							"PackageEcosystem": "gomod",
							"Directory":        "/",
							"Schedule": map[string]any{
								"Interval": "weekly",
								"Day":      "monday",
								"Time":     "09:00",
								"Timezone": "America/New_York",
							},
							"Labels":                  []any{"dependencies", "go"},
							"Assignees":               []any{"@maintainer"},
							"Reviewers":               []any{"@team"},
							"OpenPullRequestsLimit":   5.0,
							"RebaseStrategy":          "auto",
							"VersioningStrategy":      "increase",
							"TargetBranch":            "main",
							"Vendor":                  true,
							"Allow": []any{
								map[string]any{
									"DependencyName": "github.com/stretchr/testify",
									"DependencyType": "direct",
								},
							},
							"Ignore": []any{
								map[string]any{
									"DependencyName": "github.com/old/package",
									"Versions":       []any{"1.0.0"},
								},
							},
							"CommitMessage": map[string]any{
								"Prefix":  "deps",
								"Include": "scope",
							},
							"Groups": map[string]any{
								"test-deps": map[string]any{
									"Patterns":       []any{"*test*"},
									"DependencyType": "development",
								},
							},
						},
					},
					"Registries": map[string]any{
						"github": map[string]any{
							"Type":     "git",
							"URL":      "https://github.com",
							"Username": "x-access-token",
							"Password": "${{secrets.GITHUB_TOKEN}}",
						},
					},
				},
			},
		},
	}

	result, err := b.BuildDependabot(discovered, extracted)
	if err != nil {
		t.Fatalf("BuildDependabot() error = %v", err)
	}

	if len(result.Configs) != 1 {
		t.Fatalf("Expected 1 config, got %d", len(result.Configs))
	}

	cfg := result.Configs[0]

	// Verify complex structure was reconstructed correctly
	if !cfg.Config.EnableBetaEcosystems {
		t.Error("EnableBetaEcosystems should be true")
	}

	if len(cfg.Config.Updates) != 1 {
		t.Fatalf("Expected 1 update, got %d", len(cfg.Config.Updates))
	}

	update := cfg.Config.Updates[0]
	if update.PackageEcosystem != "gomod" {
		t.Errorf("PackageEcosystem = %q, want %q", update.PackageEcosystem, "gomod")
	}
	if update.Schedule.Interval != "weekly" {
		t.Errorf("Schedule.Interval = %q, want %q", update.Schedule.Interval, "weekly")
	}
	if len(update.Labels) != 2 {
		t.Errorf("Expected 2 labels, got %d", len(update.Labels))
	}
	if len(update.Allow) != 1 {
		t.Errorf("Expected 1 allow entry, got %d", len(update.Allow))
	}
	if len(update.Ignore) != 1 {
		t.Errorf("Expected 1 ignore entry, got %d", len(update.Ignore))
	}
	if len(update.Groups) != 1 {
		t.Errorf("Expected 1 group, got %d", len(update.Groups))
	}

	if len(cfg.Config.Registries) != 1 {
		t.Errorf("Expected 1 registry, got %d", len(cfg.Config.Registries))
	}

	// Check YAML output
	yaml := string(cfg.YAML)
	if !strings.Contains(yaml, "version:") {
		t.Error("YAML should contain version field")
	}
	if !strings.Contains(yaml, "package-ecosystem:") {
		t.Error("YAML should contain package-ecosystem field")
	}
}
