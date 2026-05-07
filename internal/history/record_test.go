package history_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/driftwatch/internal/history"
)

func tempHistoryPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "history.json")
}

func TestNewRecorder_EmptyFile(t *testing.T) {
	rec, err := history.NewRecorder(tempHistoryPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec == nil {
		t.Fatal("expected non-nil recorder")
	}
	if got := rec.Entries(); len(got) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(got))
	}
}

func TestRecord_AppendsEntry(t *testing.T) {
	path := tempHistoryPath(t)
	rec, _ := history.NewRecorder(path)

	e := history.Entry{Host: "web-01", Drifted: true, Details: []string{"file changed"}}
	if err := rec.Record(e); err != nil {
		t.Fatalf("Record failed: %v", err)
	}

	entries := rec.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Host != "web-01" {
		t.Errorf("expected host web-01, got %s", entries[0].Host)
	}
	if entries[0].Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestRecord_Persistence(t *testing.T) {
	path := tempHistoryPath(t)
	rec, _ := history.NewRecorder(path)
	_ = rec.Record(history.Entry{Host: "db-01", Drifted: false})
	_ = rec.Record(history.Entry{Host: "db-02", Drifted: true})

	// Reload from disk.
	rec2, err := history.NewRecorder(path)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	if got := len(rec2.Entries()); got != 2 {
		t.Fatalf("expected 2 entries after reload, got %d", got)
	}
}

func TestDriftedHosts_UniqueOnly(t *testing.T) {
	path := tempHistoryPath(t)
	rec, _ := history.NewRecorder(path)

	for i := 0; i < 3; i++ {
		_ = rec.Record(history.Entry{Host: "web-01", Drifted: true, Timestamp: time.Now().UTC()})
	}
	_ = rec.Record(history.Entry{Host: "web-02", Drifted: true})
	_ = rec.Record(history.Entry{Host: "db-01", Drifted: false})

	hosts := rec.DriftedHosts()
	if len(hosts) != 2 {
		t.Fatalf("expected 2 drifted hosts, got %d: %v", len(hosts), hosts)
	}
}

func TestNewRecorder_InvalidJSON(t *testing.T) {
	path := tempHistoryPath(t)
	_ = os.WriteFile(path, []byte("not-json"), 0o644)

	_, err := history.NewRecorder(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}
