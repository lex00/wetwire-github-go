package workflow_test

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestList(t *testing.T) {
	t.Run("strings", func(t *testing.T) {
		list := workflow.List("main", "develop", "release")

		if len(list) != 3 {
			t.Errorf("expected 3 items, got %d", len(list))
		}

		if list[0] != "main" {
			t.Errorf("expected first item to be 'main', got %q", list[0])
		}

		if list[2] != "release" {
			t.Errorf("expected third item to be 'release', got %q", list[2])
		}
	})

	t.Run("integers", func(t *testing.T) {
		list := workflow.List(1, 2, 3, 4, 5)

		if len(list) != 5 {
			t.Errorf("expected 5 items, got %d", len(list))
		}

		if list[0] != 1 {
			t.Errorf("expected first item to be 1, got %d", list[0])
		}

		if list[4] != 5 {
			t.Errorf("expected last item to be 5, got %d", list[4])
		}
	})

	t.Run("empty", func(t *testing.T) {
		list := workflow.List[string]()

		if len(list) != 0 {
			t.Errorf("expected 0 items, got %d", len(list))
		}
	})

	t.Run("single item", func(t *testing.T) {
		list := workflow.List("only")

		if len(list) != 1 {
			t.Errorf("expected 1 item, got %d", len(list))
		}

		if list[0] != "only" {
			t.Errorf("expected 'only', got %q", list[0])
		}
	})
}

func TestStrings(t *testing.T) {
	t.Run("multiple strings", func(t *testing.T) {
		list := workflow.Strings("a", "b", "c")

		if len(list) != 3 {
			t.Errorf("expected 3 items, got %d", len(list))
		}

		if list[1] != "b" {
			t.Errorf("expected second item to be 'b', got %q", list[1])
		}
	})

	t.Run("empty", func(t *testing.T) {
		list := workflow.Strings()

		if len(list) != 0 {
			t.Errorf("expected 0 items, got %d", len(list))
		}
	})
}

func TestPtr(t *testing.T) {
	t.Run("bool", func(t *testing.T) {
		val := workflow.Ptr(true)

		if val == nil {
			t.Fatal("expected non-nil pointer")
		}

		if !*val {
			t.Error("expected pointer to true")
		}
	})

	t.Run("string", func(t *testing.T) {
		val := workflow.Ptr("test")

		if val == nil {
			t.Fatal("expected non-nil pointer")
		}

		if *val != "test" {
			t.Errorf("expected 'test', got %q", *val)
		}
	})

	t.Run("int", func(t *testing.T) {
		val := workflow.Ptr(42)

		if val == nil {
			t.Fatal("expected non-nil pointer")
		}

		if *val != 42 {
			t.Errorf("expected 42, got %d", *val)
		}
	})

	t.Run("false bool", func(t *testing.T) {
		val := workflow.Ptr(false)

		if val == nil {
			t.Fatal("expected non-nil pointer")
		}

		if *val {
			t.Error("expected pointer to false")
		}
	})
}

func TestEnvType(t *testing.T) {
	env := workflow.Env{
		"KEY1": "value1",
		"KEY2": workflow.Secrets.Get("SECRET"),
		"KEY3": 123,
	}

	if env["KEY1"] != "value1" {
		t.Errorf("expected KEY1=value1, got %v", env["KEY1"])
	}

	expr, ok := env["KEY2"].(workflow.Expression)
	if !ok {
		t.Error("expected KEY2 to be an Expression")
	}

	if expr.Raw() != "secrets.SECRET" {
		t.Errorf("expected secrets.SECRET, got %s", expr.Raw())
	}

	if env["KEY3"] != 123 {
		t.Errorf("expected KEY3=123, got %v", env["KEY3"])
	}
}

func TestWithType(t *testing.T) {
	with := workflow.With{
		"param1": "value1",
		"param2": 42,
		"param3": true,
		"param4": workflow.GitHub.SHA(),
	}

	if with["param1"] != "value1" {
		t.Errorf("expected param1=value1, got %v", with["param1"])
	}

	if with["param2"] != 42 {
		t.Errorf("expected param2=42, got %v", with["param2"])
	}

	if with["param3"] != true {
		t.Errorf("expected param3=true, got %v", with["param3"])
	}

	expr, ok := with["param4"].(workflow.Expression)
	if !ok {
		t.Error("expected param4 to be an Expression")
	}

	if expr.Raw() != "github.sha" {
		t.Errorf("expected github.sha, got %s", expr.Raw())
	}
}
