package domain

import (
	"testing"

	coredomain "github.com/lex00/wetwire-core-go/domain"
)

func TestGitHubDomainImplementsInterface(t *testing.T) {
	// Compile-time check that GitHubDomain implements Domain
	var _ coredomain.Domain = (*GitHubDomain)(nil)
}

func TestGitHubDomainImplementsListerDomain(t *testing.T) {
	// Compile-time check that GitHubDomain implements ListerDomain
	var _ coredomain.ListerDomain = (*GitHubDomain)(nil)
}

func TestGitHubDomainImplementsGrapherDomain(t *testing.T) {
	// Compile-time check that GitHubDomain implements GrapherDomain
	var _ coredomain.GrapherDomain = (*GitHubDomain)(nil)
}

func TestGitHubDomainName(t *testing.T) {
	d := &GitHubDomain{}
	if d.Name() != "github" {
		t.Errorf("expected name 'github', got %q", d.Name())
	}
}

func TestGitHubDomainVersion(t *testing.T) {
	d := &GitHubDomain{}
	v := d.Version()
	if v == "" {
		t.Error("version should not be empty")
	}
}

func TestGitHubDomainBuilder(t *testing.T) {
	d := &GitHubDomain{}
	b := d.Builder()
	if b == nil {
		t.Error("builder should not be nil")
	}
}

func TestGitHubDomainLinter(t *testing.T) {
	d := &GitHubDomain{}
	l := d.Linter()
	if l == nil {
		t.Error("linter should not be nil")
	}
}

func TestGitHubDomainInitializer(t *testing.T) {
	d := &GitHubDomain{}
	i := d.Initializer()
	if i == nil {
		t.Error("initializer should not be nil")
	}
}

func TestGitHubDomainValidator(t *testing.T) {
	d := &GitHubDomain{}
	v := d.Validator()
	if v == nil {
		t.Error("validator should not be nil")
	}
}

func TestGitHubDomainLister(t *testing.T) {
	d := &GitHubDomain{}
	l := d.Lister()
	if l == nil {
		t.Error("lister should not be nil")
	}
}

func TestGitHubDomainGrapher(t *testing.T) {
	d := &GitHubDomain{}
	g := d.Grapher()
	if g == nil {
		t.Error("grapher should not be nil")
	}
}

func TestCreateRootCommand(t *testing.T) {
	cmd := CreateRootCommand(&GitHubDomain{})
	if cmd == nil {
		t.Fatal("root command should not be nil")
	}
	if cmd.Use != "wetwire-github" {
		t.Errorf("expected Use 'wetwire-github', got %q", cmd.Use)
	}
}
