package personas

import (
	"testing"
)

func TestAll(t *testing.T) {
	personas := All()

	if len(personas) != 5 {
		t.Errorf("All() returned %d personas, want 5", len(personas))
	}

	expectedNames := []string{"beginner", "intermediate", "expert", "terse", "verbose"}
	for i, p := range personas {
		if p.Name != expectedNames[i] {
			t.Errorf("personas[%d].Name = %q, want %q", i, p.Name, expectedNames[i])
		}
	}
}

func TestGet_ValidPersonas(t *testing.T) {
	tests := []string{"beginner", "intermediate", "expert", "terse", "verbose"}

	for _, name := range tests {
		t.Run(name, func(t *testing.T) {
			p, err := Get(name)
			if err != nil {
				t.Fatalf("Get(%q) error = %v", name, err)
			}
			if p.Name != name {
				t.Errorf("Get(%q).Name = %q, want %q", name, p.Name, name)
			}
		})
	}
}

func TestGet_CaseInsensitive(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"BEGINNER", "beginner"},
		{"Beginner", "beginner"},
		{"EXPERT", "expert"},
		{"Expert", "expert"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p, err := Get(tt.input)
			if err != nil {
				t.Fatalf("Get(%q) error = %v", tt.input, err)
			}
			if p.Name != tt.want {
				t.Errorf("Get(%q).Name = %q, want %q", tt.input, p.Name, tt.want)
			}
		})
	}
}

func TestGet_Unknown(t *testing.T) {
	_, err := Get("unknown")
	if err == nil {
		t.Error("Get(\"unknown\") should return error")
	}
}

func TestNames(t *testing.T) {
	names := Names()

	if len(names) != 5 {
		t.Errorf("Names() returned %d names, want 5", len(names))
	}

	expected := []string{"beginner", "intermediate", "expert", "terse", "verbose"}
	for i, name := range names {
		if name != expected[i] {
			t.Errorf("Names()[%d] = %q, want %q", i, name, expected[i])
		}
	}
}

func TestPersona_Fields(t *testing.T) {
	tests := []struct {
		persona     Persona
		wantName    string
		wantDescLen int // minimum description length
	}{
		{Beginner, "beginner", 20},
		{Intermediate, "intermediate", 20},
		{Expert, "expert", 20},
		{Terse, "terse", 20},
		{Verbose, "verbose", 20},
	}

	for _, tt := range tests {
		t.Run(tt.wantName, func(t *testing.T) {
			if tt.persona.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", tt.persona.Name, tt.wantName)
			}

			if len(tt.persona.Description) < tt.wantDescLen {
				t.Errorf("Description too short: %d chars", len(tt.persona.Description))
			}

			if tt.persona.SystemPrompt == "" {
				t.Error("SystemPrompt is empty")
			}

			if tt.persona.ExpectedBehavior == "" {
				t.Error("ExpectedBehavior is empty")
			}
		})
	}
}

func TestBeginner_GitHubActionsContext(t *testing.T) {
	// Verify beginner persona has GitHub Actions-specific content
	if Beginner.SystemPrompt == "" {
		t.Error("Beginner.SystemPrompt is empty")
	}

	// Should mention GitHub Actions concepts
	containsGitHub := false
	keywords := []string{"GitHub Actions", "CI", "triggers", "workflow"}
	for _, kw := range keywords {
		if contains(Beginner.SystemPrompt, kw) || contains(Beginner.Description, kw) {
			containsGitHub = true
			break
		}
	}
	if !containsGitHub {
		t.Error("Beginner persona should reference GitHub Actions concepts")
	}
}

func TestExpert_GitHubActionsContext(t *testing.T) {
	// Verify expert persona has GitHub Actions-specific content
	if Expert.SystemPrompt == "" {
		t.Error("Expert.SystemPrompt is empty")
	}

	// Should mention advanced GitHub Actions concepts
	advancedKeywords := []string{"matrix", "runner", "secrets", "permissions"}
	foundAdvanced := 0
	for _, kw := range advancedKeywords {
		if contains(Expert.SystemPrompt, kw) {
			foundAdvanced++
		}
	}
	if foundAdvanced < 2 {
		t.Errorf("Expert persona should mention advanced concepts, found %d", foundAdvanced)
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 &&
		(s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || contains(s[1:], substr)))
}
