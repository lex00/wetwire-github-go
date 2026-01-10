package import_gpg

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestImportGPG_Action(t *testing.T) {
	a := ImportGPG{}
	if got := a.Action(); got != "crazy-max/ghaction-import-gpg@v6" {
		t.Errorf("Action() = %q, want %q", got, "crazy-max/ghaction-import-gpg@v6")
	}
}

func TestImportGPG_Inputs_Empty(t *testing.T) {
	a := ImportGPG{}
	inputs := a.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty ImportGPG.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestImportGPG_ImplementsStepAction(t *testing.T) {
	a := ImportGPG{}
	var _ workflow.StepAction = a
}

func TestImportGPG_Inputs_GPGPrivateKey(t *testing.T) {
	a := ImportGPG{
		GPGPrivateKey: "${{ secrets.GPG_PRIVATE_KEY }}",
	}

	inputs := a.Inputs()

	if inputs["gpg_private_key"] != "${{ secrets.GPG_PRIVATE_KEY }}" {
		t.Errorf("inputs[gpg_private_key] = %v, want %q", inputs["gpg_private_key"], "${{ secrets.GPG_PRIVATE_KEY }}")
	}
}

func TestImportGPG_Inputs_Passphrase(t *testing.T) {
	a := ImportGPG{
		Passphrase: "${{ secrets.PASSPHRASE }}",
	}

	inputs := a.Inputs()

	if inputs["passphrase"] != "${{ secrets.PASSPHRASE }}" {
		t.Errorf("inputs[passphrase] = %v, want %q", inputs["passphrase"], "${{ secrets.PASSPHRASE }}")
	}
}

func TestImportGPG_Inputs_GitUserSigningkey(t *testing.T) {
	a := ImportGPG{
		GitUserSigningkey: true,
	}

	inputs := a.Inputs()

	if inputs["git_user_signingkey"] != true {
		t.Errorf("inputs[git_user_signingkey] = %v, want true", inputs["git_user_signingkey"])
	}
}

func TestImportGPG_Inputs_GitCommitGpgsign(t *testing.T) {
	a := ImportGPG{
		GitCommitGpgsign: true,
	}

	inputs := a.Inputs()

	if inputs["git_commit_gpgsign"] != true {
		t.Errorf("inputs[git_commit_gpgsign] = %v, want true", inputs["git_commit_gpgsign"])
	}
}

func TestImportGPG_Inputs_GitTagGpgsign(t *testing.T) {
	a := ImportGPG{
		GitTagGpgsign: true,
	}

	inputs := a.Inputs()

	if inputs["git_tag_gpgsign"] != true {
		t.Errorf("inputs[git_tag_gpgsign] = %v, want true", inputs["git_tag_gpgsign"])
	}
}

func TestImportGPG_Inputs_GitPushGpgsign(t *testing.T) {
	a := ImportGPG{
		GitPushGpgsign: true,
	}

	inputs := a.Inputs()

	if inputs["git_push_gpgsign"] != true {
		t.Errorf("inputs[git_push_gpgsign] = %v, want true", inputs["git_push_gpgsign"])
	}
}

func TestImportGPG_Inputs_Fingerprint(t *testing.T) {
	a := ImportGPG{
		Fingerprint: "ABC123DEF456",
	}

	inputs := a.Inputs()

	if inputs["fingerprint"] != "ABC123DEF456" {
		t.Errorf("inputs[fingerprint] = %v, want %q", inputs["fingerprint"], "ABC123DEF456")
	}
}

func TestImportGPG_Inputs_TrustLevel(t *testing.T) {
	a := ImportGPG{
		TrustLevel: "5",
	}

	inputs := a.Inputs()

	if inputs["trust_level"] != "5" {
		t.Errorf("inputs[trust_level] = %v, want %q", inputs["trust_level"], "5")
	}
}

func TestImportGPG_Inputs_GitConfigGlobal(t *testing.T) {
	a := ImportGPG{
		GitConfigGlobal: true,
	}

	inputs := a.Inputs()

	if inputs["git_config_global"] != true {
		t.Errorf("inputs[git_config_global] = %v, want true", inputs["git_config_global"])
	}
}

func TestImportGPG_Inputs_Workdir(t *testing.T) {
	a := ImportGPG{
		Workdir: "/path/to/repo",
	}

	inputs := a.Inputs()

	if inputs["workdir"] != "/path/to/repo" {
		t.Errorf("inputs[workdir] = %v, want %q", inputs["workdir"], "/path/to/repo")
	}
}

func TestImportGPG_Inputs_AllFields(t *testing.T) {
	a := ImportGPG{
		GPGPrivateKey:     "${{ secrets.GPG_PRIVATE_KEY }}",
		Passphrase:        "${{ secrets.PASSPHRASE }}",
		GitUserSigningkey: true,
		GitCommitGpgsign:  true,
		GitTagGpgsign:     true,
		GitPushGpgsign:    true,
		Fingerprint:       "ABC123DEF456",
		TrustLevel:        "5",
		GitConfigGlobal:   true,
		Workdir:           "/custom/path",
	}

	inputs := a.Inputs()

	expected := map[string]any{
		"gpg_private_key":    "${{ secrets.GPG_PRIVATE_KEY }}",
		"passphrase":         "${{ secrets.PASSPHRASE }}",
		"git_user_signingkey": true,
		"git_commit_gpgsign": true,
		"git_tag_gpgsign":    true,
		"git_push_gpgsign":   true,
		"fingerprint":        "ABC123DEF456",
		"trust_level":        "5",
		"git_config_global":  true,
		"workdir":            "/custom/path",
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

func TestImportGPG_Inputs_FalseBoolFields(t *testing.T) {
	a := ImportGPG{
		GitUserSigningkey: false,
		GitCommitGpgsign:  false,
		GitTagGpgsign:     false,
		GitPushGpgsign:    false,
		GitConfigGlobal:   false,
	}

	inputs := a.Inputs()

	if len(inputs) != 0 {
		t.Errorf("inputs for false bools has %d entries, want 0. Got: %v", len(inputs), inputs)
	}
}

func TestImportGPG_Inputs_CommonUsage(t *testing.T) {
	// Test common usage pattern: signing commits
	a := ImportGPG{
		GPGPrivateKey:     "${{ secrets.GPG_PRIVATE_KEY }}",
		Passphrase:        "${{ secrets.PASSPHRASE }}",
		GitUserSigningkey: true,
		GitCommitGpgsign:  true,
	}

	inputs := a.Inputs()

	if len(inputs) != 4 {
		t.Errorf("common usage has %d entries, want 4", len(inputs))
	}

	if inputs["gpg_private_key"] != "${{ secrets.GPG_PRIVATE_KEY }}" {
		t.Errorf("inputs[gpg_private_key] = %v, want secret reference", inputs["gpg_private_key"])
	}
}
