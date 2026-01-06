package serialize

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/lex00/wetwire-github-go/dependabot"
)

// DependabotToYAML serializes a Dependabot config to YAML bytes.
func DependabotToYAML(d *dependabot.Dependabot) ([]byte, error) {
	m, err := dependabotToMap(d)
	if err != nil {
		return nil, fmt.Errorf("converting dependabot to map: %w", err)
	}

	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(m); err != nil {
		return nil, fmt.Errorf("encoding YAML: %w", err)
	}
	if err := encoder.Close(); err != nil {
		return nil, fmt.Errorf("closing encoder: %w", err)
	}

	return buf.Bytes(), nil
}

// dependabotToMap converts a Dependabot config to a map for YAML serialization.
func dependabotToMap(d *dependabot.Dependabot) (map[string]any, error) {
	m := make(map[string]any)

	m["version"] = d.Version

	if d.EnableBetaEcosystems {
		m["enable-beta-ecosystems"] = true
	}

	if len(d.Updates) > 0 {
		updates := make([]any, len(d.Updates))
		for i, u := range d.Updates {
			updates[i] = updateToMap(&u)
		}
		m["updates"] = updates
	}

	if len(d.Registries) > 0 {
		registries := make(map[string]any)
		for name, r := range d.Registries {
			registries[name] = registryToMap(&r)
		}
		m["registries"] = registries
	}

	return m, nil
}

// updateToMap converts an Update to a map for YAML serialization.
func updateToMap(u *dependabot.Update) map[string]any {
	m := make(map[string]any)

	m["package-ecosystem"] = u.PackageEcosystem

	if u.Directory != "" {
		m["directory"] = u.Directory
	}
	if len(u.Directories) > 0 {
		m["directories"] = u.Directories
	}

	m["schedule"] = scheduleToMap(&u.Schedule)

	if len(u.Allow) > 0 {
		allow := make([]any, len(u.Allow))
		for i, a := range u.Allow {
			allow[i] = allowToMap(&a)
		}
		m["allow"] = allow
	}

	if len(u.Ignore) > 0 {
		ignore := make([]any, len(u.Ignore))
		for i, ig := range u.Ignore {
			ignore[i] = ignoreToMap(&ig)
		}
		m["ignore"] = ignore
	}

	if len(u.Labels) > 0 {
		m["labels"] = u.Labels
	}
	if len(u.Assignees) > 0 {
		m["assignees"] = u.Assignees
	}
	if len(u.Reviewers) > 0 {
		m["reviewers"] = u.Reviewers
	}
	if u.Milestone > 0 {
		m["milestone"] = u.Milestone
	}
	if u.OpenPullRequestsLimit > 0 {
		m["open-pull-requests-limit"] = u.OpenPullRequestsLimit
	}
	if u.RebaseStrategy != "" {
		m["rebase-strategy"] = u.RebaseStrategy
	}
	if u.VersioningStrategy != "" {
		m["versioning-strategy"] = u.VersioningStrategy
	}
	if u.Vendor {
		m["vendor"] = true
	}
	if u.TargetBranch != "" {
		m["target-branch"] = u.TargetBranch
	}
	if u.Registries != nil {
		m["registries"] = u.Registries
	}

	if len(u.Groups) > 0 {
		groups := make(map[string]any)
		for name, g := range u.Groups {
			groups[name] = groupToMap(&g)
		}
		m["groups"] = groups
	}

	if u.CommitMessage != nil {
		m["commit-message"] = commitMessageToMap(u.CommitMessage)
	}
	if u.PullRequestBranchName != nil {
		m["pull-request-branch-name"] = map[string]any{
			"separator": u.PullRequestBranchName.Separator,
		}
	}
	if u.InsecureExternalCodeExecution != "" {
		m["insecure-external-code-execution"] = u.InsecureExternalCodeExecution
	}

	return m
}

// scheduleToMap converts a Schedule to a map.
func scheduleToMap(s *dependabot.Schedule) map[string]any {
	m := make(map[string]any)

	m["interval"] = s.Interval

	if s.Day != "" {
		m["day"] = s.Day
	}
	if s.Time != "" {
		m["time"] = s.Time
	}
	if s.Timezone != "" {
		m["timezone"] = s.Timezone
	}

	return m
}

// allowToMap converts an Allow to a map.
func allowToMap(a *dependabot.Allow) map[string]any {
	m := make(map[string]any)

	if a.DependencyName != "" {
		m["dependency-name"] = a.DependencyName
	}
	if a.DependencyType != "" {
		m["dependency-type"] = a.DependencyType
	}

	return m
}

// ignoreToMap converts an Ignore to a map.
func ignoreToMap(i *dependabot.Ignore) map[string]any {
	m := make(map[string]any)

	if i.DependencyName != "" {
		m["dependency-name"] = i.DependencyName
	}
	if len(i.Versions) > 0 {
		m["versions"] = i.Versions
	}
	if len(i.UpdateTypes) > 0 {
		m["update-types"] = i.UpdateTypes
	}

	return m
}

// groupToMap converts a Group to a map.
func groupToMap(g *dependabot.Group) map[string]any {
	m := make(map[string]any)

	if len(g.Patterns) > 0 {
		m["patterns"] = g.Patterns
	}
	if g.DependencyType != "" {
		m["dependency-type"] = g.DependencyType
	}
	if len(g.UpdateTypes) > 0 {
		m["update-types"] = g.UpdateTypes
	}
	if len(g.ExcludePatterns) > 0 {
		m["exclude-patterns"] = g.ExcludePatterns
	}
	if g.AppliesTo != "" {
		m["applies-to"] = g.AppliesTo
	}

	return m
}

// registryToMap converts a Registry to a map.
func registryToMap(r *dependabot.Registry) map[string]any {
	m := make(map[string]any)

	m["type"] = r.Type

	if r.URL != "" {
		m["url"] = r.URL
	}
	if r.Username != "" {
		m["username"] = r.Username
	}
	if r.Password != "" {
		m["password"] = r.Password
	}
	if r.Token != "" {
		m["token"] = r.Token
	}
	if r.Key != "" {
		m["key"] = r.Key
	}
	if r.Organization != "" {
		m["organization"] = r.Organization
	}
	if r.ReplacesBase {
		m["replaces-base"] = true
	}

	return m
}

// commitMessageToMap converts a CommitMessage to a map.
func commitMessageToMap(c *dependabot.CommitMessage) map[string]any {
	m := make(map[string]any)

	if c.Prefix != "" {
		m["prefix"] = c.Prefix
	}
	if c.PrefixDevelopment != "" {
		m["prefix-development"] = c.PrefixDevelopment
	}
	if c.Include != "" {
		m["include"] = c.Include
	}

	return m
}
