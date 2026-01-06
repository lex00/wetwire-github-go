package template

import (
	"github.com/lex00/wetwire-github-go/dependabot"
	"github.com/lex00/wetwire-github-go/internal/discover"
	"github.com/lex00/wetwire-github-go/internal/runner"
	"github.com/lex00/wetwire-github-go/internal/serialize"
)

// DependabotBuildResult contains the result of building Dependabot templates.
type DependabotBuildResult struct {
	// Configs contains the assembled configs with YAML output
	Configs []BuiltDependabot

	// Errors contains any non-fatal errors encountered
	Errors []string
}

// BuiltDependabot represents a Dependabot config ready for output.
type BuiltDependabot struct {
	// Name is the config variable name
	Name string

	// Config is the assembled Dependabot config
	Config *dependabot.Dependabot

	// YAML is the serialized YAML output
	YAML []byte
}

// BuildDependabot assembles Dependabot templates from discovery and extraction results.
func (b *Builder) BuildDependabot(discovered *discover.DependabotDiscoveryResult, extracted *runner.DependabotExtractionResult) (*DependabotBuildResult, error) {
	result := &DependabotBuildResult{
		Configs: []BuiltDependabot{},
		Errors:  []string{},
	}

	// Process each config
	for _, dc := range discovered.Configs {
		// Find the extracted config data
		var configData map[string]any
		for _, ec := range extracted.Configs {
			if ec.Name == dc.Name {
				configData = ec.Data
				break
			}
		}

		if configData == nil {
			result.Errors = append(result.Errors, "config "+dc.Name+": extraction data not found")
			continue
		}

		// Reconstruct the Dependabot config from the map
		config := b.reconstructDependabot(configData)

		// Serialize to YAML
		yaml, err := serialize.DependabotToYAML(config)
		if err != nil {
			result.Errors = append(result.Errors, "config "+dc.Name+": "+err.Error())
			continue
		}

		result.Configs = append(result.Configs, BuiltDependabot{
			Name:   dc.Name,
			Config: config,
			YAML:   yaml,
		})
	}

	return result, nil
}

// reconstructDependabot reconstructs a Dependabot from a map.
func (b *Builder) reconstructDependabot(data map[string]any) *dependabot.Dependabot {
	config := &dependabot.Dependabot{}

	if v, ok := data["Version"].(int); ok {
		config.Version = v
	} else if v, ok := data["Version"].(float64); ok {
		config.Version = int(v)
	}

	if v, ok := data["EnableBetaEcosystems"].(bool); ok {
		config.EnableBetaEcosystems = v
	}

	if updates, ok := data["Updates"].([]any); ok {
		config.Updates = b.reconstructUpdates(updates)
	} else if updates, ok := data["Updates"].([]dependabot.Update); ok {
		config.Updates = updates
	}

	if registries, ok := data["Registries"].(map[string]any); ok {
		config.Registries = b.reconstructRegistries(registries)
	} else if registries, ok := data["Registries"].(map[string]dependabot.Registry); ok {
		config.Registries = registries
	}

	return config
}

// reconstructUpdates reconstructs Updates from a slice.
func (b *Builder) reconstructUpdates(data []any) []dependabot.Update {
	var updates []dependabot.Update

	for _, item := range data {
		if m, ok := item.(map[string]any); ok {
			updates = append(updates, b.reconstructUpdate(m))
		} else if u, ok := item.(dependabot.Update); ok {
			updates = append(updates, u)
		}
	}

	return updates
}

// reconstructUpdate reconstructs an Update from a map.
func (b *Builder) reconstructUpdate(data map[string]any) dependabot.Update {
	update := dependabot.Update{}

	if v, ok := data["PackageEcosystem"].(string); ok {
		update.PackageEcosystem = v
	}
	if v, ok := data["Directory"].(string); ok {
		update.Directory = v
	}
	if v, ok := data["Directories"].([]string); ok {
		update.Directories = v
	} else if v, ok := data["Directories"].([]any); ok {
		update.Directories = anySliceToStrings(v)
	}

	// Schedule
	if schedule, ok := data["Schedule"].(map[string]any); ok {
		update.Schedule = b.reconstructSchedule(schedule)
	} else if schedule, ok := data["Schedule"].(dependabot.Schedule); ok {
		update.Schedule = schedule
	}

	// Allow
	if allow, ok := data["Allow"].([]any); ok {
		update.Allow = b.reconstructAllowList(allow)
	}

	// Ignore
	if ignore, ok := data["Ignore"].([]any); ok {
		update.Ignore = b.reconstructIgnoreList(ignore)
	}

	// String slices
	if v, ok := data["Labels"].([]any); ok {
		update.Labels = anySliceToStrings(v)
	}
	if v, ok := data["Assignees"].([]any); ok {
		update.Assignees = anySliceToStrings(v)
	}
	if v, ok := data["Reviewers"].([]any); ok {
		update.Reviewers = anySliceToStrings(v)
	}

	// Integers
	if v, ok := data["Milestone"].(int); ok {
		update.Milestone = v
	} else if v, ok := data["Milestone"].(float64); ok {
		update.Milestone = int(v)
	}
	if v, ok := data["OpenPullRequestsLimit"].(int); ok {
		update.OpenPullRequestsLimit = v
	} else if v, ok := data["OpenPullRequestsLimit"].(float64); ok {
		update.OpenPullRequestsLimit = int(v)
	}

	// Strings
	if v, ok := data["RebaseStrategy"].(string); ok {
		update.RebaseStrategy = v
	}
	if v, ok := data["VersioningStrategy"].(string); ok {
		update.VersioningStrategy = v
	}
	if v, ok := data["TargetBranch"].(string); ok {
		update.TargetBranch = v
	}
	if v, ok := data["InsecureExternalCodeExecution"].(string); ok {
		update.InsecureExternalCodeExecution = v
	}

	// Booleans
	if v, ok := data["Vendor"].(bool); ok {
		update.Vendor = v
	}

	// Registries
	if v, ok := data["Registries"]; ok {
		update.Registries = v
	}

	// Groups
	if groups, ok := data["Groups"].(map[string]any); ok {
		update.Groups = b.reconstructGroups(groups)
	}

	// CommitMessage
	if cm, ok := data["CommitMessage"].(map[string]any); ok {
		update.CommitMessage = b.reconstructCommitMessage(cm)
	} else if cm, ok := data["CommitMessage"].(*dependabot.CommitMessage); ok {
		update.CommitMessage = cm
	}

	// PullRequestBranchName
	if prbn, ok := data["PullRequestBranchName"].(map[string]any); ok {
		update.PullRequestBranchName = b.reconstructPullRequestBranchName(prbn)
	} else if prbn, ok := data["PullRequestBranchName"].(*dependabot.PullRequestBranchName); ok {
		update.PullRequestBranchName = prbn
	}

	return update
}

// reconstructSchedule reconstructs a Schedule from a map.
func (b *Builder) reconstructSchedule(data map[string]any) dependabot.Schedule {
	schedule := dependabot.Schedule{}

	if v, ok := data["Interval"].(string); ok {
		schedule.Interval = v
	}
	if v, ok := data["Day"].(string); ok {
		schedule.Day = v
	}
	if v, ok := data["Time"].(string); ok {
		schedule.Time = v
	}
	if v, ok := data["Timezone"].(string); ok {
		schedule.Timezone = v
	}

	return schedule
}

// reconstructAllowList reconstructs Allow list from a slice.
func (b *Builder) reconstructAllowList(data []any) []dependabot.Allow {
	var result []dependabot.Allow

	for _, item := range data {
		if m, ok := item.(map[string]any); ok {
			allow := dependabot.Allow{}
			if v, ok := m["DependencyName"].(string); ok {
				allow.DependencyName = v
			}
			if v, ok := m["DependencyType"].(string); ok {
				allow.DependencyType = v
			}
			result = append(result, allow)
		}
	}

	return result
}

// reconstructIgnoreList reconstructs Ignore list from a slice.
func (b *Builder) reconstructIgnoreList(data []any) []dependabot.Ignore {
	var result []dependabot.Ignore

	for _, item := range data {
		if m, ok := item.(map[string]any); ok {
			ignore := dependabot.Ignore{}
			if v, ok := m["DependencyName"].(string); ok {
				ignore.DependencyName = v
			}
			if v, ok := m["Versions"].([]any); ok {
				ignore.Versions = anySliceToStrings(v)
			}
			if v, ok := m["UpdateTypes"].([]any); ok {
				ignore.UpdateTypes = anySliceToStrings(v)
			}
			result = append(result, ignore)
		}
	}

	return result
}

// reconstructGroups reconstructs Groups from a map.
func (b *Builder) reconstructGroups(data map[string]any) map[string]dependabot.Group {
	result := make(map[string]dependabot.Group)

	for name, item := range data {
		if m, ok := item.(map[string]any); ok {
			group := dependabot.Group{}
			if v, ok := m["Patterns"].([]any); ok {
				group.Patterns = anySliceToStrings(v)
			}
			if v, ok := m["DependencyType"].(string); ok {
				group.DependencyType = v
			}
			if v, ok := m["UpdateTypes"].([]any); ok {
				group.UpdateTypes = anySliceToStrings(v)
			}
			if v, ok := m["ExcludePatterns"].([]any); ok {
				group.ExcludePatterns = anySliceToStrings(v)
			}
			if v, ok := m["AppliesTo"].(string); ok {
				group.AppliesTo = v
			}
			result[name] = group
		}
	}

	return result
}

// reconstructRegistries reconstructs Registries from a map.
func (b *Builder) reconstructRegistries(data map[string]any) map[string]dependabot.Registry {
	result := make(map[string]dependabot.Registry)

	for name, item := range data {
		if m, ok := item.(map[string]any); ok {
			registry := dependabot.Registry{}
			if v, ok := m["Type"].(string); ok {
				registry.Type = v
			}
			if v, ok := m["URL"].(string); ok {
				registry.URL = v
			}
			if v, ok := m["Username"].(string); ok {
				registry.Username = v
			}
			if v, ok := m["Password"].(string); ok {
				registry.Password = v
			}
			if v, ok := m["Token"].(string); ok {
				registry.Token = v
			}
			if v, ok := m["Key"].(string); ok {
				registry.Key = v
			}
			if v, ok := m["Organization"].(string); ok {
				registry.Organization = v
			}
			if v, ok := m["ReplacesBase"].(bool); ok {
				registry.ReplacesBase = v
			}
			result[name] = registry
		}
	}

	return result
}

// reconstructCommitMessage reconstructs CommitMessage from a map.
func (b *Builder) reconstructCommitMessage(data map[string]any) *dependabot.CommitMessage {
	cm := &dependabot.CommitMessage{}

	if v, ok := data["Prefix"].(string); ok {
		cm.Prefix = v
	}
	if v, ok := data["PrefixDevelopment"].(string); ok {
		cm.PrefixDevelopment = v
	}
	if v, ok := data["Include"].(string); ok {
		cm.Include = v
	}

	return cm
}

// reconstructPullRequestBranchName reconstructs PullRequestBranchName from a map.
func (b *Builder) reconstructPullRequestBranchName(data map[string]any) *dependabot.PullRequestBranchName {
	prbn := &dependabot.PullRequestBranchName{}

	if v, ok := data["Separator"].(string); ok {
		prbn.Separator = v
	}

	return prbn
}
