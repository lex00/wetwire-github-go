package checkout

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestCheckout_Action(t *testing.T) {
	c := Checkout{}
	if got := c.Action(); got != "actions/checkout@v4" {
		t.Errorf("Action() = %q, want %q", got, "actions/checkout@v4")
	}
}

func TestCheckout_Inputs(t *testing.T) {
	c := Checkout{
		Repository: "owner/repo",
		Ref:        "main",
		FetchDepth: 1,
		Submodules: "recursive",
	}

	inputs := c.Inputs()

	if inputs["repository"] != "owner/repo" {
		t.Errorf("inputs[repository] = %v, want %q", inputs["repository"], "owner/repo")
	}

	if inputs["ref"] != "main" {
		t.Errorf("inputs[ref] = %v, want %q", inputs["ref"], "main")
	}

	if inputs["fetch-depth"] != 1 {
		t.Errorf("inputs[fetch-depth] = %v, want 1", inputs["fetch-depth"])
	}

	if inputs["submodules"] != "recursive" {
		t.Errorf("inputs[submodules] = %v, want %q", inputs["submodules"], "recursive")
	}
}

func TestCheckout_Inputs_Empty(t *testing.T) {
	c := Checkout{}
	inputs := c.Inputs()

	// Empty checkout should have no inputs
	if len(inputs) != 0 {
		t.Errorf("empty Checkout.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestCheckout_Inputs_BoolFields(t *testing.T) {
	c := Checkout{
		Clean:              true,
		LFS:                true,
		PersistCredentials: true,
	}

	inputs := c.Inputs()

	if inputs["clean"] != true {
		t.Errorf("inputs[clean] = %v, want true", inputs["clean"])
	}

	if inputs["lfs"] != true {
		t.Errorf("inputs[lfs] = %v, want true", inputs["lfs"])
	}

	if inputs["persist-credentials"] != true {
		t.Errorf("inputs[persist-credentials] = %v, want true", inputs["persist-credentials"])
	}
}

func TestCheckout_ImplementsStepAction(t *testing.T) {
	c := Checkout{}
	// Verify Checkout implements StepAction interface
	var _ workflow.StepAction = c
}

func TestCheckout_Inputs_Token(t *testing.T) {
	c := Checkout{
		Token: "ghp_test123",
	}

	inputs := c.Inputs()

	if inputs["token"] != "ghp_test123" {
		t.Errorf("inputs[token] = %v, want %q", inputs["token"], "ghp_test123")
	}
}

func TestCheckout_Inputs_SSHKey(t *testing.T) {
	c := Checkout{
		SSHKey: "ssh-rsa AAAAB3NzaC1yc2E...",
	}

	inputs := c.Inputs()

	if inputs["ssh-key"] != "ssh-rsa AAAAB3NzaC1yc2E..." {
		t.Errorf("inputs[ssh-key] = %v, want %q", inputs["ssh-key"], "ssh-rsa AAAAB3NzaC1yc2E...")
	}
}

func TestCheckout_Inputs_SSHKnownHosts(t *testing.T) {
	c := Checkout{
		SSHKnownHosts: "github.com ssh-rsa ...",
	}

	inputs := c.Inputs()

	if inputs["ssh-known-hosts"] != "github.com ssh-rsa ..." {
		t.Errorf("inputs[ssh-known-hosts] = %v, want %q", inputs["ssh-known-hosts"], "github.com ssh-rsa ...")
	}
}

func TestCheckout_Inputs_SSHStrict(t *testing.T) {
	c := Checkout{
		SSHStrict: true,
	}

	inputs := c.Inputs()

	if inputs["ssh-strict"] != true {
		t.Errorf("inputs[ssh-strict] = %v, want true", inputs["ssh-strict"])
	}
}

func TestCheckout_Inputs_Path(t *testing.T) {
	c := Checkout{
		Path: "custom/path",
	}

	inputs := c.Inputs()

	if inputs["path"] != "custom/path" {
		t.Errorf("inputs[path] = %v, want %q", inputs["path"], "custom/path")
	}
}

func TestCheckout_Inputs_Filter(t *testing.T) {
	c := Checkout{
		Filter: "blob:none",
	}

	inputs := c.Inputs()

	if inputs["filter"] != "blob:none" {
		t.Errorf("inputs[filter] = %v, want %q", inputs["filter"], "blob:none")
	}
}

func TestCheckout_Inputs_SparseCheckout(t *testing.T) {
	c := Checkout{
		SparseCheckout: "src/\ndocs/",
	}

	inputs := c.Inputs()

	if inputs["sparse-checkout"] != "src/\ndocs/" {
		t.Errorf("inputs[sparse-checkout] = %v, want %q", inputs["sparse-checkout"], "src/\ndocs/")
	}
}

func TestCheckout_Inputs_SparseCheckoutConeMode(t *testing.T) {
	c := Checkout{
		SparseCheckoutConeMode: true,
	}

	inputs := c.Inputs()

	if inputs["sparse-checkout-cone-mode"] != true {
		t.Errorf("inputs[sparse-checkout-cone-mode] = %v, want true", inputs["sparse-checkout-cone-mode"])
	}
}

func TestCheckout_Inputs_FetchTags(t *testing.T) {
	c := Checkout{
		FetchTags: true,
	}

	inputs := c.Inputs()

	if inputs["fetch-tags"] != true {
		t.Errorf("inputs[fetch-tags] = %v, want true", inputs["fetch-tags"])
	}
}

func TestCheckout_Inputs_ShowProgress(t *testing.T) {
	c := Checkout{
		ShowProgress: true,
	}

	inputs := c.Inputs()

	if inputs["show-progress"] != true {
		t.Errorf("inputs[show-progress] = %v, want true", inputs["show-progress"])
	}
}

func TestCheckout_Inputs_SetSafeDirectory(t *testing.T) {
	c := Checkout{
		SetSafeDirectory: true,
	}

	inputs := c.Inputs()

	if inputs["set-safe-directory"] != true {
		t.Errorf("inputs[set-safe-directory] = %v, want true", inputs["set-safe-directory"])
	}
}

func TestCheckout_Inputs_GithubServerURL(t *testing.T) {
	c := Checkout{
		GithubServerURL: "https://github.enterprise.com",
	}

	inputs := c.Inputs()

	if inputs["github-server-url"] != "https://github.enterprise.com" {
		t.Errorf("inputs[github-server-url] = %v, want %q", inputs["github-server-url"], "https://github.enterprise.com")
	}
}

func TestCheckout_Inputs_AllFields(t *testing.T) {
	c := Checkout{
		Repository:             "owner/repo",
		Ref:                    "v1.0.0",
		Token:                  "token123",
		SSHKey:                 "ssh-key-value",
		SSHKnownHosts:          "known-hosts",
		SSHStrict:              true,
		PersistCredentials:     true,
		Path:                   "my-path",
		Clean:                  true,
		Filter:                 "tree:0",
		SparseCheckout:         "src/",
		SparseCheckoutConeMode: true,
		FetchDepth:             5,
		FetchTags:              true,
		ShowProgress:           true,
		LFS:                    true,
		Submodules:             "true",
		SetSafeDirectory:       true,
		GithubServerURL:        "https://custom.github.com",
	}

	inputs := c.Inputs()

	// Verify all fields are present
	expected := map[string]any{
		"repository":               "owner/repo",
		"ref":                      "v1.0.0",
		"token":                    "token123",
		"ssh-key":                  "ssh-key-value",
		"ssh-known-hosts":          "known-hosts",
		"ssh-strict":               true,
		"persist-credentials":      true,
		"path":                     "my-path",
		"clean":                    true,
		"filter":                   "tree:0",
		"sparse-checkout":          "src/",
		"sparse-checkout-cone-mode": true,
		"fetch-depth":              5,
		"fetch-tags":               true,
		"show-progress":            true,
		"lfs":                      true,
		"submodules":               "true",
		"set-safe-directory":       true,
		"github-server-url":        "https://custom.github.com",
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

func TestCheckout_Inputs_FalseBoolFields(t *testing.T) {
	// Test that false boolean values are not included in inputs
	c := Checkout{
		SSHStrict:              false,
		PersistCredentials:     false,
		Clean:                  false,
		SparseCheckoutConeMode: false,
		FetchTags:              false,
		ShowProgress:           false,
		LFS:                    false,
		SetSafeDirectory:       false,
	}

	inputs := c.Inputs()

	// None of these should be in the inputs map
	if len(inputs) != 0 {
		t.Errorf("inputs for false bools has %d entries, want 0. Got: %v", len(inputs), inputs)
	}
}

func TestCheckout_Inputs_ZeroFetchDepth(t *testing.T) {
	// Test that FetchDepth = 0 is not included (0 means all history)
	c := Checkout{
		FetchDepth: 0,
	}

	inputs := c.Inputs()

	if _, exists := inputs["fetch-depth"]; exists {
		t.Errorf("inputs[fetch-depth] should not exist for FetchDepth=0")
	}
}

func TestCheckout_Inputs_NegativeFetchDepth(t *testing.T) {
	// Edge case: negative fetch depth (should be included as non-zero)
	c := Checkout{
		FetchDepth: -1,
	}

	inputs := c.Inputs()

	if inputs["fetch-depth"] != -1 {
		t.Errorf("inputs[fetch-depth] = %v, want -1", inputs["fetch-depth"])
	}
}
