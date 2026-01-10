package template

import (
	"strings"
	"testing"

	"github.com/lex00/wetwire-github-go/internal/discover"
	"github.com/lex00/wetwire-github-go/internal/runner"
)

func TestBuilder_BuildCodeowners_Empty(t *testing.T) {
	b := NewBuilder()

	discovered := &discover.CodeownersDiscoveryResult{
		Configs: []discover.DiscoveredCodeowners{},
	}
	extracted := &runner.CodeownersExtractionResult{
		Configs: []runner.ExtractedCodeowners{},
	}

	result, err := b.BuildCodeowners(discovered, extracted)
	if err != nil {
		t.Fatalf("BuildCodeowners() error = %v", err)
	}

	if len(result.Configs) != 0 {
		t.Errorf("Expected 0 configs, got %d", len(result.Configs))
	}
}

func TestBuilder_BuildCodeowners_SingleConfig(t *testing.T) {
	b := NewBuilder()

	discovered := &discover.CodeownersDiscoveryResult{
		Configs: []discover.DiscoveredCodeowners{
			{Name: "CodeOwners", File: "codeowners.go", Line: 10},
		},
	}
	extracted := &runner.CodeownersExtractionResult{
		Configs: []runner.ExtractedCodeowners{
			{
				Name: "CodeOwners",
				Rules: []runner.ExtractedCodeownersRule{
					{Pattern: "*", Owners: []string{"@default-team"}},
					{Pattern: "/docs/", Owners: []string{"@docs-team"}},
				},
			},
		},
	}

	result, err := b.BuildCodeowners(discovered, extracted)
	if err != nil {
		t.Fatalf("BuildCodeowners() error = %v", err)
	}

	if len(result.Configs) != 1 {
		t.Fatalf("Expected 1 config, got %d", len(result.Configs))
	}

	cfg := result.Configs[0]
	if cfg.Name != "CodeOwners" {
		t.Errorf("Config name = %q, want %q", cfg.Name, "CodeOwners")
	}

	if len(cfg.Content) == 0 {
		t.Error("Config content is empty")
	}

	// Check content contains expected patterns
	content := string(cfg.Content)
	if !strings.Contains(content, "* @default-team") {
		t.Errorf("Content missing expected rule, got: %s", content)
	}
	if !strings.Contains(content, "/docs/ @docs-team") {
		t.Errorf("Content missing expected rule, got: %s", content)
	}
}

func TestBuilder_BuildCodeowners_WithComments(t *testing.T) {
	b := NewBuilder()

	discovered := &discover.CodeownersDiscoveryResult{
		Configs: []discover.DiscoveredCodeowners{
			{Name: "Owners", File: "codeowners.go", Line: 10},
		},
	}
	extracted := &runner.CodeownersExtractionResult{
		Configs: []runner.ExtractedCodeowners{
			{
				Name: "Owners",
				Rules: []runner.ExtractedCodeownersRule{
					{Pattern: "*.go", Owners: []string{"@go-team"}, Comment: "Go source files"},
				},
			},
		},
	}

	result, err := b.BuildCodeowners(discovered, extracted)
	if err != nil {
		t.Fatalf("BuildCodeowners() error = %v", err)
	}

	if len(result.Configs) != 1 {
		t.Fatalf("Expected 1 config, got %d", len(result.Configs))
	}

	content := string(result.Configs[0].Content)
	if !strings.Contains(content, "# Go source files") {
		t.Errorf("Content missing comment, got: %s", content)
	}
}

func TestBuilder_BuildCodeowners_MultipleOwners(t *testing.T) {
	b := NewBuilder()

	discovered := &discover.CodeownersDiscoveryResult{
		Configs: []discover.DiscoveredCodeowners{
			{Name: "Owners", File: "codeowners.go", Line: 10},
		},
	}
	extracted := &runner.CodeownersExtractionResult{
		Configs: []runner.ExtractedCodeowners{
			{
				Name: "Owners",
				Rules: []runner.ExtractedCodeownersRule{
					{Pattern: "*.ts", Owners: []string{"@frontend", "@typescript-guild", "@code-review"}},
				},
			},
		},
	}

	result, err := b.BuildCodeowners(discovered, extracted)
	if err != nil {
		t.Fatalf("BuildCodeowners() error = %v", err)
	}

	content := string(result.Configs[0].Content)
	if !strings.Contains(content, "*.ts @frontend @typescript-guild @code-review") {
		t.Errorf("Content missing multiple owners, got: %s", content)
	}
}

func TestBuilder_BuildCodeowners_MissingExtraction(t *testing.T) {
	b := NewBuilder()

	discovered := &discover.CodeownersDiscoveryResult{
		Configs: []discover.DiscoveredCodeowners{
			{Name: "MissingConfig", File: "codeowners.go", Line: 10},
		},
	}
	// No extraction data
	extracted := &runner.CodeownersExtractionResult{
		Configs: []runner.ExtractedCodeowners{},
	}

	result, err := b.BuildCodeowners(discovered, extracted)
	if err != nil {
		t.Fatalf("BuildCodeowners() error = %v", err)
	}

	// Should have an error about missing extraction data
	if len(result.Errors) == 0 {
		t.Error("Expected error about missing extraction data")
	}

	if len(result.Configs) != 0 {
		t.Errorf("Expected 0 configs when extraction is missing, got %d", len(result.Configs))
	}
}

func TestBuilder_BuildCodeowners_HeaderComment(t *testing.T) {
	b := NewBuilder()

	discovered := &discover.CodeownersDiscoveryResult{
		Configs: []discover.DiscoveredCodeowners{
			{Name: "Owners", File: "codeowners.go", Line: 10},
		},
	}
	extracted := &runner.CodeownersExtractionResult{
		Configs: []runner.ExtractedCodeowners{
			{
				Name:  "Owners",
				Rules: []runner.ExtractedCodeownersRule{},
			},
		},
	}

	result, err := b.BuildCodeowners(discovered, extracted)
	if err != nil {
		t.Fatalf("BuildCodeowners() error = %v", err)
	}

	content := string(result.Configs[0].Content)
	if !strings.Contains(content, "# CODEOWNERS file generated by wetwire-github") {
		t.Errorf("Content missing header comment, got: %s", content)
	}
}
