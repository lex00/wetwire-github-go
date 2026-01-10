package codeowners

import (
	"testing"
)

func TestOwners_ResourceType(t *testing.T) {
	o := Owners{}
	if got := o.ResourceType(); got != "codeowners" {
		t.Errorf("ResourceType() = %q, want %q", got, "codeowners")
	}
}

func TestOwners_Basic(t *testing.T) {
	o := Owners{
		Rules: []Rule{
			{Pattern: "*", Owners: []string{"@default-team"}},
		},
	}

	if len(o.Rules) != 1 {
		t.Errorf("len(Rules) = %d, want 1", len(o.Rules))
	}

	if o.Rules[0].Pattern != "*" {
		t.Errorf("Pattern = %q, want %q", o.Rules[0].Pattern, "*")
	}
}

func TestOwners_MultipleRules(t *testing.T) {
	o := Owners{
		Rules: []Rule{
			{Pattern: "*", Owners: []string{"@default-team"}},
			{Pattern: "/docs/", Owners: []string{"@docs-team"}},
			{Pattern: "*.go", Owners: []string{"@go-team", "@code-review"}},
			{Pattern: "/src/api/", Owners: []string{"@api-team"}, Comment: "API code owners"},
		},
	}

	if len(o.Rules) != 4 {
		t.Errorf("len(Rules) = %d, want 4", len(o.Rules))
	}

	// Check the last rule with comment
	lastRule := o.Rules[3]
	if lastRule.Pattern != "/src/api/" {
		t.Errorf("Pattern = %q, want %q", lastRule.Pattern, "/src/api/")
	}
	if lastRule.Comment != "API code owners" {
		t.Errorf("Comment = %q, want %q", lastRule.Comment, "API code owners")
	}
}

func TestRule_MultipleOwners(t *testing.T) {
	r := Rule{
		Pattern: "*.ts",
		Owners:  []string{"@frontend-team", "@typescript-guild", "@user1"},
	}

	if len(r.Owners) != 3 {
		t.Errorf("len(Owners) = %d, want 3", len(r.Owners))
	}

	if r.Owners[0] != "@frontend-team" {
		t.Errorf("Owners[0] = %q, want %q", r.Owners[0], "@frontend-team")
	}
}

func TestRule_DirectoryPattern(t *testing.T) {
	r := Rule{
		Pattern: "/src/components/",
		Owners:  []string{"@ui-team"},
	}

	if r.Pattern != "/src/components/" {
		t.Errorf("Pattern = %q, want %q", r.Pattern, "/src/components/")
	}
}

func TestRule_GlobPattern(t *testing.T) {
	r := Rule{
		Pattern: "src/**/*.tsx",
		Owners:  []string{"@react-team"},
	}

	if r.Pattern != "src/**/*.tsx" {
		t.Errorf("Pattern = %q, want %q", r.Pattern, "src/**/*.tsx")
	}
}

func TestOwners_EmptyRules(t *testing.T) {
	o := Owners{
		Rules: []Rule{},
	}

	if len(o.Rules) != 0 {
		t.Errorf("len(Rules) = %d, want 0", len(o.Rules))
	}
}

func TestRule_EmptyOwners(t *testing.T) {
	r := Rule{
		Pattern: "*.md",
		Owners:  []string{},
	}

	if len(r.Owners) != 0 {
		t.Errorf("len(Owners) = %d, want 0", len(r.Owners))
	}
}
