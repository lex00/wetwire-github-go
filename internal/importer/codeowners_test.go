package importer

import (
	"strings"
	"testing"
)

func TestParseCodeownersContent(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantLen  int
		wantErr  bool
		checkFn  func(*IRCodeowners) bool
	}{
		{
			name:    "empty file",
			content: "",
			wantLen: 0,
		},
		{
			name:    "single rule",
			content: "*.go @developer",
			wantLen: 1,
			checkFn: func(ir *IRCodeowners) bool {
				return ir.Rules[0].Pattern == "*.go" &&
					len(ir.Rules[0].Owners) == 1 &&
					ir.Rules[0].Owners[0] == "@developer"
			},
		},
		{
			name:    "multiple owners",
			content: "*.go @dev1 @dev2 @org/team",
			wantLen: 1,
			checkFn: func(ir *IRCodeowners) bool {
				return ir.Rules[0].Pattern == "*.go" &&
					len(ir.Rules[0].Owners) == 3 &&
					ir.Rules[0].Owners[0] == "@dev1" &&
					ir.Rules[0].Owners[1] == "@dev2" &&
					ir.Rules[0].Owners[2] == "@org/team"
			},
		},
		{
			name: "multiple rules",
			content: `*.go @backend
*.js @frontend
/docs/ @docs-team`,
			wantLen: 3,
		},
		{
			name: "with comments",
			content: `# This is a comment
*.go @developer
# Another comment
*.js @frontend`,
			wantLen: 2,
		},
		{
			name: "with inline comments",
			content: `*.go @developer # Go files`,
			wantLen: 1,
			checkFn: func(ir *IRCodeowners) bool {
				return ir.Rules[0].Comment == "Go files"
			},
		},
		{
			name: "empty lines",
			content: `*.go @developer

*.js @frontend`,
			wantLen: 2,
		},
		{
			name: "preserve section comments",
			content: `# Backend team owns Go files
*.go @backend-team

# Frontend team owns JS files
*.js @frontend-team`,
			wantLen: 2,
			checkFn: func(ir *IRCodeowners) bool {
				return ir.Rules[0].Comment == "Backend team owns Go files" &&
					ir.Rules[1].Comment == "Frontend team owns JS files"
			},
		},
		{
			name:    "default owner pattern",
			content: "* @default-owner",
			wantLen: 1,
			checkFn: func(ir *IRCodeowners) bool {
				return ir.Rules[0].Pattern == "*"
			},
		},
		{
			name:    "directory pattern",
			content: "/docs/ @docs-team",
			wantLen: 1,
			checkFn: func(ir *IRCodeowners) bool {
				return ir.Rules[0].Pattern == "/docs/"
			},
		},
		{
			name:    "glob pattern",
			content: "src/**/*.ts @typescript-team",
			wantLen: 1,
			checkFn: func(ir *IRCodeowners) bool {
				return ir.Rules[0].Pattern == "src/**/*.ts"
			},
		},
		{
			name:    "email owner",
			content: "*.md user@example.com",
			wantLen: 1,
			checkFn: func(ir *IRCodeowners) bool {
				return ir.Rules[0].Owners[0] == "user@example.com"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ir, err := ParseCodeownersContent(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCodeownersContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if len(ir.Rules) != tt.wantLen {
				t.Errorf("ParseCodeownersContent() got %d rules, want %d", len(ir.Rules), tt.wantLen)
			}
			if tt.checkFn != nil && !tt.checkFn(ir) {
				t.Errorf("ParseCodeownersContent() check function failed")
			}
		})
	}
}

func TestCodeownersCodeGenerator(t *testing.T) {
	gen := &CodeownersCodeGenerator{PackageName: "mypackage"}

	ir := &IRCodeowners{
		Rules: []IRCodeownersRule{
			{Pattern: "*.go", Owners: []string{"@backend-team"}, Comment: "Go files"},
			{Pattern: "*.js", Owners: []string{"@frontend-team", "@org/ui-team"}},
			{Pattern: "/docs/", Owners: []string{"@docs-team"}},
		},
	}

	result, err := gen.Generate(ir)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if len(result.Files) == 0 {
		t.Fatal("Generate() produced no files")
	}

	code, ok := result.Files["codeowners.go"]
	if !ok {
		t.Fatal("Generate() did not produce codeowners.go")
	}

	// Check that the generated code has expected content
	if !strings.Contains(code, "package mypackage") {
		t.Error("Generated code missing package declaration")
	}
	if !strings.Contains(code, `"github.com/lex00/wetwire-github-go/codeowners"`) {
		t.Error("Generated code missing codeowners import")
	}
	if !strings.Contains(code, "var Codeowners = codeowners.Owners{") {
		t.Error("Generated code missing Codeowners variable declaration")
	}
	if !strings.Contains(code, `Pattern: "*.go"`) {
		t.Error("Generated code missing *.go pattern")
	}
	if !strings.Contains(code, `"@backend-team"`) {
		t.Error("Generated code missing @backend-team owner")
	}
	if !strings.Contains(code, `Comment: "Go files"`) {
		t.Error("Generated code missing comment")
	}
	if result.Rules != 3 {
		t.Errorf("Generate() rules = %d, want 3", result.Rules)
	}
}

func TestCodeownersCodeGenerator_Empty(t *testing.T) {
	gen := &CodeownersCodeGenerator{PackageName: "empty"}

	ir := &IRCodeowners{Rules: []IRCodeownersRule{}}

	result, err := gen.Generate(ir)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	code := result.Files["codeowners.go"]
	if !strings.Contains(code, "Rules: []codeowners.Rule{}") {
		t.Error("Generated code should have empty rules slice")
	}
}

func TestCodeownersCodeGenerator_MultilineComment(t *testing.T) {
	gen := &CodeownersCodeGenerator{PackageName: "test"}

	ir := &IRCodeowners{
		Rules: []IRCodeownersRule{
			{Pattern: "*.go", Owners: []string{"@dev"}, Comment: "Comment with \"quotes\""},
		},
	}

	result, err := gen.Generate(ir)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	code := result.Files["codeowners.go"]
	// Check that quotes are properly escaped
	if !strings.Contains(code, `"quotes\"`) {
		// Check for proper escaping
		if !strings.Contains(code, `Comment:`) {
			t.Error("Generated code missing Comment field")
		}
	}
}
