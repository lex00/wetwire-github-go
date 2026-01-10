package configure_pages

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestConfigurePages_Action(t *testing.T) {
	c := ConfigurePages{}
	if got := c.Action(); got != "actions/configure-pages@v5" {
		t.Errorf("Action() = %q, want %q", got, "actions/configure-pages@v5")
	}
}

func TestConfigurePages_Inputs_Empty(t *testing.T) {
	c := ConfigurePages{}
	inputs := c.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty ConfigurePages.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestConfigurePages_Inputs_AllFields(t *testing.T) {
	c := ConfigurePages{
		StaticSiteGenerator: "next",
		GeneratorConfigFile: "next.config.js",
		Token:               "${{ secrets.GITHUB_TOKEN }}",
	}

	inputs := c.Inputs()

	expected := map[string]any{
		"static_site_generator": "next",
		"generator_config_file": "next.config.js",
		"token":                 "${{ secrets.GITHUB_TOKEN }}",
	}

	if len(inputs) != len(expected) {
		t.Errorf("inputs has %d entries, want %d", len(inputs), len(expected))
	}

	for key, want := range expected {
		if got := inputs[key]; got != want {
			t.Errorf("inputs[%q] = %v, want %v", key, got, want)
		}
	}
}

func TestConfigurePages_Inputs_StaticSiteGenerator(t *testing.T) {
	c := ConfigurePages{
		StaticSiteGenerator: "gatsby",
	}

	inputs := c.Inputs()

	if inputs["static_site_generator"] != "gatsby" {
		t.Errorf("inputs[static_site_generator] = %v, want %q", inputs["static_site_generator"], "gatsby")
	}
	if len(inputs) != 1 {
		t.Errorf("inputs has %d entries, want 1", len(inputs))
	}
}

func TestConfigurePages_Inputs_GeneratorConfigFile(t *testing.T) {
	c := ConfigurePages{
		GeneratorConfigFile: "gatsby-config.js",
	}

	inputs := c.Inputs()

	if inputs["generator_config_file"] != "gatsby-config.js" {
		t.Errorf("inputs[generator_config_file] = %v, want %q", inputs["generator_config_file"], "gatsby-config.js")
	}
	if len(inputs) != 1 {
		t.Errorf("inputs has %d entries, want 1", len(inputs))
	}
}

func TestConfigurePages_Inputs_Token(t *testing.T) {
	c := ConfigurePages{
		Token: "ghp_token123",
	}

	inputs := c.Inputs()

	if inputs["token"] != "ghp_token123" {
		t.Errorf("inputs[token] = %v, want %q", inputs["token"], "ghp_token123")
	}
	if len(inputs) != 1 {
		t.Errorf("inputs has %d entries, want 1", len(inputs))
	}
}

func TestConfigurePages_ImplementsStepAction(t *testing.T) {
	c := ConfigurePages{}
	// Verify ConfigurePages implements StepAction interface
	var _ workflow.StepAction = c
}

func TestConfigurePages_InSteps(t *testing.T) {
	// Test that ConfigurePages can be used in a steps slice
	steps := []any{
		ConfigurePages{
			StaticSiteGenerator: "jekyll",
		},
	}

	if len(steps) != 1 {
		t.Errorf("steps has %d entries, want 1", len(steps))
	}

	cp, ok := steps[0].(ConfigurePages)
	if !ok {
		t.Fatal("steps[0] is not ConfigurePages")
	}

	if cp.StaticSiteGenerator != "jekyll" {
		t.Errorf("StaticSiteGenerator = %q, want %q", cp.StaticSiteGenerator, "jekyll")
	}
}

func TestConfigurePages_Generators(t *testing.T) {
	// Test common static site generator values
	generators := []string{"next", "nuxt", "gatsby", "jekyll", "hugo", "sveltekit"}

	for _, gen := range generators {
		c := ConfigurePages{
			StaticSiteGenerator: gen,
		}
		inputs := c.Inputs()
		if inputs["static_site_generator"] != gen {
			t.Errorf("inputs[static_site_generator] = %v, want %q", inputs["static_site_generator"], gen)
		}
	}
}
