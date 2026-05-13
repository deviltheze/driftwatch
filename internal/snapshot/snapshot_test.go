package snapshot

import (
	"strings"
	"testing"
	"time"
)

func makeEntry(host, path, hash string) Entry {
	return Entry{Host: host, Path: path, Hash: hash, CapturedAt: time.Now()}
}

func TestCompare_NoDiff(t *testing.T) {
	prev := makeEntry("host1", "/etc/hosts", "abc123")
	curr := makeEntry("host1", "/etc/hosts", "abc123")
	_, changed := Compare(prev, curr)
	if changed {
		t.Fatal("expected no change for identical hashes")
	}
}

func TestCompare_WithDiff(t *testing.T) {
	prev := makeEntry("host1", "/etc/hosts", "aabbccdd")
	curr := makeEntry("host1", "/etc/hosts", "11223344")
	diff, changed := Compare(prev, curr)
	if !changed {
		t.Fatal("expected change for different hashes")
	}
	if diff.Host != "host1" || diff.Path != "/etc/hosts" {
		t.Errorf("unexpected diff metadata: %+v", diff)
	}
}

func TestDiff_Changed(t *testing.T) {
	d := Diff{OldHash: "aaa", NewHash: "bbb"}
	if !d.Changed() {
		t.Fatal("expected Changed() true")
	}
	d2 := Diff{OldHash: "aaa", NewHash: "aaa"}
	if d2.Changed() {
		t.Fatal("expected Changed() false")
	}
}

func TestDiff_String(t *testing.T) {
	d := Diff{Host: "h", Path: "/p", OldHash: "aaaaaaaa1111", NewHash: "bbbbbbbb2222"}
	s := d.String()
	if !strings.Contains(s, "h") || !strings.Contains(s, "/p") {
		t.Errorf("String() missing host or path: %s", s)
	}
}

func TestStore_SetAndGet(t *testing.T) {
	s := NewStore()
	e := makeEntry("host1", "/etc/passwd", "deadbeef")
	s.Set(e)
	got, ok := s.Get("host1", "/etc/passwd")
	if !ok {
		t.Fatal("expected entry to be found")
	}
	if got.Hash != "deadbeef" {
		t.Errorf("unexpected hash: %s", got.Hash)
	}
}

func TestStore_Get_Missing(t *testing.T) {
	s := NewStore()
	_, ok := s.Get("ghost", "/no/such/file")
	if ok {
		t.Fatal("expected missing entry")
	}
}

func TestStore_Size(t *testing.T) {
	s := NewStore()
	if s.Size() != 0 {
		t.Fatal("expected empty store")
	}
	s.Set(makeEntry("h", "/a", "1"))
	s.Set(makeEntry("h", "/b", "2"))
	if s.Size() != 2 {
		t.Errorf("expected size 2, got %d", s.Size())
	}
}

func TestStore_All(t *testing.T) {
	s := NewStore()
	s.Set(makeEntry("h1", "/x", "h1"))
	s.Set(makeEntry("h2", "/x", "h2"))
	all := s.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}
