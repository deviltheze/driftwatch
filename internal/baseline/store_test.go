package baseline_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/driftwatch/internal/baseline"
)

func tempStorePath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "baseline.json")
}

func TestNewStore_CreatesEmpty(t *testing.T) {
	s, err := baseline.NewStore(tempStorePath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil store")
	}
}

func TestStore_SetAndGet(t *testing.T) {
	s, _ := baseline.NewStore(tempStorePath(t))
	if err := s.Set("host1", "/etc/hosts", "abc123"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	e, ok := s.Get("host1", "/etc/hosts")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Hash != "abc123" {
		t.Errorf("expected hash abc123, got %s", e.Hash)
	}
}

func TestStore_Get_Missing(t *testing.T) {
	s, _ := baseline.NewStore(tempStorePath(t))
	_, ok := s.Get("ghost", "/no/such/path")
	if ok {
		t.Error("expected missing entry")
	}
}

func TestStore_Persistence(t *testing.T) {
	path := tempStorePath(t)
	s1, _ := baseline.NewStore(path)
	_ = s1.Set("srv1", "/etc/nginx.conf", "deadbeef")

	s2, err := baseline.NewStore(path)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	e, ok := s2.Get("srv1", "/etc/nginx.conf")
	if !ok {
		t.Fatal("entry not persisted")
	}
	if e.Hash != "deadbeef" {
		t.Errorf("expected deadbeef, got %s", e.Hash)
	}
}

func TestStore_InvalidJSON(t *testing.T) {
	path := tempStorePath(t)
	if err := os.WriteFile(path, []byte("not-json"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := baseline.NewStore(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestStore_OverwriteEntry(t *testing.T) {
	s, _ := baseline.NewStore(tempStorePath(t))
	_ = s.Set("host1", "/etc/hosts", "first")
	_ = s.Set("host1", "/etc/hosts", "second")
	e, _ := s.Get("host1", "/etc/hosts")
	if e.Hash != "second" {
		t.Errorf("expected second, got %s", e.Hash)
	}
}
