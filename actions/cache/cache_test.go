package cache

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestCache_Action(t *testing.T) {
	c := Cache{}
	if got := c.Action(); got != "actions/cache@v4" {
		t.Errorf("Action() = %q, want %q", got, "actions/cache@v4")
	}
}

func TestCache_Inputs(t *testing.T) {
	c := Cache{
		Path:        "~/.cache/go-build",
		Key:         "go-cache-key",
		RestoreKeys: "go-cache-",
	}

	inputs := c.Inputs()

	if c.Action() != "actions/cache@v4" {
		t.Errorf("Action() = %q, want %q", c.Action(), "actions/cache@v4")
	}

	if inputs["path"] != "~/.cache/go-build" {
		t.Errorf("inputs[path] = %v, want %q", inputs["path"], "~/.cache/go-build")
	}

	if inputs["key"] != "go-cache-key" {
		t.Errorf("inputs[key] = %v, want %q", inputs["key"], "go-cache-key")
	}

	if inputs["restore-keys"] != "go-cache-" {
		t.Errorf("inputs[restore-keys] = %v, want %q", inputs["restore-keys"], "go-cache-")
	}
}

func TestCache_Inputs_Empty(t *testing.T) {
	c := Cache{}
	inputs := c.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty Cache.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestCache_Inputs_BoolFields(t *testing.T) {
	c := Cache{
		EnableCrossOsArchive: true,
		FailOnCacheMiss:      true,
		LookupOnly:           true,
		SaveAlways:           true,
	}

	inputs := c.Inputs()

	if inputs["enableCrossOsArchive"] != true {
		t.Errorf("inputs[enableCrossOsArchive] = %v, want true", inputs["enableCrossOsArchive"])
	}

	if inputs["fail-on-cache-miss"] != true {
		t.Errorf("inputs[fail-on-cache-miss] = %v, want true", inputs["fail-on-cache-miss"])
	}

	if inputs["lookup-only"] != true {
		t.Errorf("inputs[lookup-only] = %v, want true", inputs["lookup-only"])
	}

	if inputs["save-always"] != true {
		t.Errorf("inputs[save-always] = %v, want true", inputs["save-always"])
	}
}

func TestCache_Inputs_IntFields(t *testing.T) {
	c := Cache{
		UploadChunkSize: 1048576,
	}

	inputs := c.Inputs()

	if inputs["upload-chunk-size"] != 1048576 {
		t.Errorf("inputs[upload-chunk-size] = %v, want 1048576", inputs["upload-chunk-size"])
	}
}

func TestCache_ImplementsStepAction(t *testing.T) {
	var _ workflow.StepAction = Cache{}
}
