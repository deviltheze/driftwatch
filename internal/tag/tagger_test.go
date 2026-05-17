package tag

import (
	"testing"
)

func TestNew_Empty(t *testing.T) {
	tr := New(nil)
	if tr == nil {
		t.Fatal("expected non-nil Tagger")
	}
	if got := tr.Resolve("host1"); len(got) != 0 {
		t.Fatalf("expected empty tags, got %v", got)
	}
}

func TestNew_GlobalTags_Normalised(t *testing.T) {
	tr := New([]string{"  Prod  ", "EU", ""})
	tags := tr.Resolve("any")
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d: %v", len(tags), tags)
	}
	if tags[0] != "prod" || tags[1] != "eu" {
		t.Fatalf("unexpected tags: %v", tags)
	}
}

func TestAddOverride_MergesWithGlobal(t *testing.T) {
	tr := New([]string{"global"})
	tr.AddOverride("web01", []string{"web", "frontend"})

	tags := tr.Resolve("web01")
	if len(tags) != 3 {
		t.Fatalf("expected 3 tags, got %d: %v", len(tags), tags)
	}
}

func TestAddOverride_DeduplicatesGlobal(t *testing.T) {
	tr := New([]string{"shared"})
	tr.AddOverride("db01", []string{"shared", "database"})

	tags := tr.Resolve("db01")
	if len(tags) != 2 {
		t.Fatalf("expected 2 unique tags, got %d: %v", len(tags), tags)
	}
}

func TestResolve_UnknownHost_ReturnsGlobalOnly(t *testing.T) {
	tr := New([]string{"global"})
	tr.AddOverride("known", []string{"extra"})

	tags := tr.Resolve("unknown")
	if len(tags) != 1 || tags[0] != "global" {
		t.Fatalf("expected only global tag, got %v", tags)
	}
}

func TestHasTag_Match(t *testing.T) {
	tr := New([]string{"prod"})
	if !tr.HasTag("any", "PROD") {
		t.Fatal("expected HasTag to return true for case-insensitive match")
	}
}

func TestHasTag_NoMatch(t *testing.T) {
	tr := New([]string{"prod"})
	if tr.HasTag("any", "staging") {
		t.Fatal("expected HasTag to return false")
	}
}

func TestHasTag_OverrideTag(t *testing.T) {
	tr := New(nil)
	tr.AddOverride("cache01", []string{"redis"})
	if !tr.HasTag("cache01", "redis") {
		t.Fatal("expected override tag to be found")
	}
	if tr.HasTag("other", "redis") {
		t.Fatal("override tag should not apply to other hosts")
	}
}
