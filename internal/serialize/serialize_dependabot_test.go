package serialize_test

import (
	"strings"
	"testing"

	"github.com/lex00/wetwire-github-go/dependabot"
	"github.com/lex00/wetwire-github-go/internal/serialize"
)

// TestBasicDependabot tests basic dependabot configuration serialization.
func TestBasicDependabot(t *testing.T) {
	d := &dependabot.Dependabot{
		Version: 2,
		Updates: []dependabot.Update{
			{
				PackageEcosystem: "go",
				Directory:        "/",
				Schedule: dependabot.Schedule{
					Interval: "daily",
				},
			},
		},
	}

	yaml, err := serialize.DependabotToYAML(d)
	if err != nil {
		t.Fatalf("DependabotToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	if !strings.Contains(yamlStr, "version: 2") {
		t.Errorf("expected 'version: 2', got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "package-ecosystem: go") {
		t.Errorf("expected 'package-ecosystem: go', got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "interval: daily") {
		t.Errorf("expected 'interval: daily', got:\n%s", yamlStr)
	}
}

// TestDependabotWithAllFields tests dependabot with all optional fields.
func TestDependabotWithAllFields(t *testing.T) {
	d := &dependabot.Dependabot{
		Version:              2,
		EnableBetaEcosystems: true,
		Updates: []dependabot.Update{
			{
				PackageEcosystem: "npm",
				Directory:        "/frontend",
				Schedule: dependabot.Schedule{
					Interval: "weekly",
					Day:      "monday",
					Time:     "10:00",
					Timezone: "America/New_York",
				},
				Allow: []dependabot.Allow{
					{DependencyName: "lodash"},
					{DependencyType: "production"},
				},
				Ignore: []dependabot.Ignore{
					{
						DependencyName: "webpack",
						Versions:       []string{">= 5.0.0, < 6.0.0"},
					},
				},
				Labels:                []string{"dependencies", "npm"},
				Assignees:             []string{"@reviewer1"},
				Reviewers:             []string{"@team/reviewers"},
				Milestone:             5,
				OpenPullRequestsLimit: 10,
				RebaseStrategy:        "auto",
				VersioningStrategy:    "increase",
				Vendor:                true,
				TargetBranch:          "develop",
				CommitMessage: &dependabot.CommitMessage{
					Prefix:            "chore",
					PrefixDevelopment: "dev",
					Include:           "scope",
				},
				PullRequestBranchName: &dependabot.PullRequestBranchName{
					Separator: "/",
				},
				InsecureExternalCodeExecution: "allow",
				Groups: map[string]dependabot.Group{
					"production": {
						Patterns:        []string{"*"},
						DependencyType:  "production",
						UpdateTypes:     []string{"minor", "patch"},
						ExcludePatterns: []string{"test-*"},
						AppliesTo:       "version-updates",
					},
				},
			},
		},
		Registries: map[string]dependabot.Registry{
			"npm-private": {
				Type:         "npm-registry",
				URL:          "https://npm.example.com",
				Token:        "${{ secrets.NPM_TOKEN }}",
				ReplacesBase: true,
			},
		},
	}

	yaml, err := serialize.DependabotToYAML(d)
	if err != nil {
		t.Fatalf("DependabotToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	expectedFields := []string{
		"enable-beta-ecosystems: true",
		"package-ecosystem: npm",
		"schedule:",
		"day: monday",
		"time: \"10:00\"",
		"timezone: America/New_York",
		"allow:",
		"ignore:",
		"labels:",
		"assignees:",
		"reviewers:",
		"milestone: 5",
		"open-pull-requests-limit: 10",
		"rebase-strategy: auto",
		"versioning-strategy: increase",
		"vendor: true",
		"target-branch: develop",
		"commit-message:",
		"prefix: chore",
		"pull-request-branch-name:",
		"separator: /",
		"groups:",
		"registries:",
	}

	for _, field := range expectedFields {
		if !strings.Contains(yamlStr, field) {
			t.Errorf("expected field %q, got:\n%s", field, yamlStr)
		}
	}
}

// TestDependabotMultipleUpdates tests multiple update configurations.
func TestDependabotMultipleUpdates(t *testing.T) {
	d := &dependabot.Dependabot{
		Version: 2,
		Updates: []dependabot.Update{
			{
				PackageEcosystem: "go",
				Directory:        "/",
				Schedule: dependabot.Schedule{
					Interval: "daily",
				},
			},
			{
				PackageEcosystem: "docker",
				Directory:        "/",
				Schedule: dependabot.Schedule{
					Interval: "weekly",
				},
			},
			{
				PackageEcosystem: "github-actions",
				Directory:        "/",
				Schedule: dependabot.Schedule{
					Interval: "weekly",
				},
			},
		},
	}

	yaml, err := serialize.DependabotToYAML(d)
	if err != nil {
		t.Fatalf("DependabotToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	if !strings.Contains(yamlStr, "package-ecosystem: go") {
		t.Errorf("expected go ecosystem, got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "package-ecosystem: docker") {
		t.Errorf("expected docker ecosystem, got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "package-ecosystem: github-actions") {
		t.Errorf("expected github-actions ecosystem, got:\n%s", yamlStr)
	}
}

// TestDependabotRegistryAllFields tests all registry fields.
func TestDependabotRegistryAllFields(t *testing.T) {
	d := &dependabot.Dependabot{
		Version: 2,
		Updates: []dependabot.Update{
			{
				PackageEcosystem: "npm",
				Directory:        "/",
				Schedule: dependabot.Schedule{
					Interval: "daily",
				},
			},
		},
		Registries: map[string]dependabot.Registry{
			"npm-registry": {
				Type:         "npm-registry",
				URL:          "https://npm.example.com",
				Username:     "user",
				Password:     "${{ secrets.NPM_PASSWORD }}",
				Token:        "${{ secrets.NPM_TOKEN }}",
				Key:          "${{ secrets.NPM_KEY }}",
				Organization: "my-org",
				ReplacesBase: true,
			},
		},
	}

	yaml, err := serialize.DependabotToYAML(d)
	if err != nil {
		t.Fatalf("DependabotToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	expectedFields := []string{
		"type: npm-registry",
		"url: https://npm.example.com",
		"username: user",
		"password:",
		"token:",
		"key:",
		"organization: my-org",
		"replaces-base: true",
	}

	for _, field := range expectedFields {
		if !strings.Contains(yamlStr, field) {
			t.Errorf("expected field %q, got:\n%s", field, yamlStr)
		}
	}
}

// TestDependabotIgnoreAllFields tests all ignore fields.
func TestDependabotIgnoreAllFields(t *testing.T) {
	d := &dependabot.Dependabot{
		Version: 2,
		Updates: []dependabot.Update{
			{
				PackageEcosystem: "npm",
				Directory:        "/",
				Schedule: dependabot.Schedule{
					Interval: "daily",
				},
				Ignore: []dependabot.Ignore{
					{
						DependencyName: "webpack",
						Versions:       []string{"5.x", "6.x"},
						UpdateTypes:    []string{"version-update:semver-major"},
					},
				},
			},
		},
	}

	yaml, err := serialize.DependabotToYAML(d)
	if err != nil {
		t.Fatalf("DependabotToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	if !strings.Contains(yamlStr, "dependency-name: webpack") {
		t.Errorf("expected dependency-name, got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "versions:") {
		t.Errorf("expected versions, got:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "update-types:") {
		t.Errorf("expected update-types, got:\n%s", yamlStr)
	}
}
