package cache

import (
	"testing"
)

func TestCache_Action(t *testing.T) {
	c := Cache{}
	if got := c.Action(); got != "actions/cache@v4" {
		t.Errorf("Action() = %q, want %q", got, "actions/cache@v4")
	}
}

func TestCache_ToStep(t *testing.T) {
	c := Cache{
		Path:        "~/.cache/go-build",
		Key:         "go-cache-key",
		RestoreKeys: "go-cache-",
	}

	step := c.ToStep()

	if step.Uses != "actions/cache@v4" {
		t.Errorf("step.Uses = %q, want %q", step.Uses, "actions/cache@v4")
	}

	if step.With["path"] != "~/.cache/go-build" {
		t.Errorf("step.With[path] = %v, want %q", step.With["path"], "~/.cache/go-build")
	}

	if step.With["key"] != "go-cache-key" {
		t.Errorf("step.With[key] = %v, want %q", step.With["key"], "go-cache-key")
	}

	if step.With["restore-keys"] != "go-cache-" {
		t.Errorf("step.With[restore-keys] = %v, want %q", step.With["restore-keys"], "go-cache-")
	}
}

func TestCache_ToStep_Empty(t *testing.T) {
	c := Cache{}
	step := c.ToStep()

	if len(step.With) != 0 {
		t.Errorf("empty Cache.ToStep() has %d with entries, want 0", len(step.With))
	}
}

func TestCache_ToStep_BoolFields(t *testing.T) {
	c := Cache{
		EnableCrossOsArchive: true,
		FailOnCacheMiss:      true,
		LookupOnly:           true,
		SaveAlways:           true,
	}

	step := c.ToStep()

	if step.With["enableCrossOsArchive"] != true {
		t.Errorf("step.With[enableCrossOsArchive] = %v, want true", step.With["enableCrossOsArchive"])
	}

	if step.With["fail-on-cache-miss"] != true {
		t.Errorf("step.With[fail-on-cache-miss] = %v, want true", step.With["fail-on-cache-miss"])
	}

	if step.With["lookup-only"] != true {
		t.Errorf("step.With[lookup-only] = %v, want true", step.With["lookup-only"])
	}

	if step.With["save-always"] != true {
		t.Errorf("step.With[save-always] = %v, want true", step.With["save-always"])
	}
}

func TestCache_ToStep_IntFields(t *testing.T) {
	c := Cache{
		UploadChunkSize: 1048576,
	}

	step := c.ToStep()

	if step.With["upload-chunk-size"] != 1048576 {
		t.Errorf("step.With[upload-chunk-size] = %v, want 1048576", step.With["upload-chunk-size"])
	}
}
