package plugin_test

import (
	"testing"

	"github.com/yourusername/driftwatch/internal/plugin"
)

func noop(_, _ string) (string, error) { return "", nil }

func TestNewRegistry_Empty(t *testing.T) {
	reg := plugin.NewRegistry()
	if reg.Len() != 0 {
		t.Fatalf("expected 0 plugins, got %d", reg.Len())
	}
}

func TestRegister_Success(t *testing.T) {
	reg := plugin.NewRegistry()
	err := reg.Register(plugin.Plugin{Name: "p1", Check: noop})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reg.Len() != 1 {
		t.Fatalf("expected 1 plugin, got %d", reg.Len())
	}
}

func TestRegister_EmptyName_ReturnsError(t *testing.T) {
	reg := plugin.NewRegistry()
	err := reg.Register(plugin.Plugin{Name: "", Check: noop})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestRegister_NilCheck_ReturnsError(t *testing.T) {
	reg := plugin.NewRegistry()
	err := reg.Register(plugin.Plugin{Name: "p1", Check: nil})
	if err == nil {
		t.Fatal("expected error for nil CheckFn")
	}
}

func TestRegister_Duplicate_ReturnsError(t *testing.T) {
	reg := plugin.NewRegistry()
	_ = reg.Register(plugin.Plugin{Name: "dup", Check: noop})
	err := reg.Register(plugin.Plugin{Name: "dup", Check: noop})
	if err == nil {
		t.Fatal("expected error for duplicate plugin name")
	}
}

func TestGet_ExistingPlugin(t *testing.T) {
	reg := plugin.NewRegistry()
	_ = reg.Register(plugin.Plugin{Name: "found", Check: noop})
	p, ok := reg.Get("found")
	if !ok {
		t.Fatal("expected plugin to be found")
	}
	if p.Name != "found" {
		t.Fatalf("expected name 'found', got %q", p.Name)
	}
}

func TestGet_MissingPlugin(t *testing.T) {
	reg := plugin.NewRegistry()
	_, ok := reg.Get("missing")
	if ok {
		t.Fatal("expected plugin not to be found")
	}
}

func TestNames_ReturnsAllNames(t *testing.T) {
	reg := plugin.NewRegistry()
	_ = reg.Register(plugin.Plugin{Name: "a", Check: noop})
	_ = reg.Register(plugin.Plugin{Name: "b", Check: noop})
	names := reg.Names()
	if len(names) != 2 {
		t.Fatalf("expected 2 names, got %d", len(names))
	}
}

func TestCheckFn_Invoked(t *testing.T) {
	called := false
	reg := plugin.NewRegistry()
	_ = reg.Register(plugin.Plugin{
		Name: "spy",
		Check: func(host, out string) (string, error) {
			called = true
			return "annotated", nil
		},
	})
	p, _ := reg.Get("spy")
	annotation, err := p.Check("host1", "raw output")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected CheckFn to be called")
	}
	if annotation != "annotated" {
		t.Fatalf("unexpected annotation: %q", annotation)
	}
}
