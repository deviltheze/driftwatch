package filter_test

import (
	"sort"
	"testing"

	"github.com/yourusername/driftwatch/internal/filter"
)

func TestMatch_NoIncludeNoExclude_AllPass(t *testing.T) {
	f := filter.New(nil, nil)
	if !f.Match([]string{"production", "web"}) {
		t.Error("expected match with empty include/exclude")
	}
}

func TestMatch_IncludeTag_Matches(t *testing.T) {
	f := filter.New([]string{"production"}, nil)
	if !f.Match([]string{"production", "web"}) {
		t.Error("expected host with 'production' tag to match")
	}
}

func TestMatch_IncludeTag_NoMatch(t *testing.T) {
	f := filter.New([]string{"production"}, nil)
	if f.Match([]string{"development"}) {
		t.Error("expected host without 'production' tag to not match")
	}
}

func TestMatch_ExcludeTag_Suppressed(t *testing.T) {
	f := filter.New(nil, []string{"maintenance"})
	if f.Match([]string{"production", "maintenance"}) {
		t.Error("expected host with 'maintenance' tag to be excluded")
	}
}

func TestMatch_IncludeAndExclude_ExcludeTakesPrecedence(t *testing.T) {
	f := filter.New([]string{"production"}, []string{"maintenance"})
	if f.Match([]string{"production", "maintenance"}) {
		t.Error("expected exclude to take precedence over include")
	}
}

func TestMatch_CaseInsensitive(t *testing.T) {
	f := filter.New([]string{"Production"}, nil)
	if !f.Match([]string{"production"}) {
		t.Error("expected case-insensitive match")
	}
}

func TestFilter_ReturnsMatchingHosts(t *testing.T) {
	f := filter.New([]string{"production"}, []string{"maintenance"})
	hosts := map[string][]string{
		"web-01": {"production", "web"},
		"db-01":  {"production", "maintenance"},
		"dev-01": {"development"},
	}

	got := f.Filter(hosts)
	sort.Strings(got)

	if len(got) != 1 || got[0] != "web-01" {
		t.Errorf("expected [web-01], got %v", got)
	}
}

func TestFilter_EmptyHosts_ReturnsNil(t *testing.T) {
	f := filter.New(nil, nil)
	got := f.Filter(map[string][]string{})
	if len(got) != 0 {
		t.Errorf("expected empty result, got %v", got)
	}
}

func TestFilter_NoInclude_AllNonExcluded(t *testing.T) {
	f := filter.New(nil, []string{"disabled"})
	hosts := map[string][]string{
		"alpha": {"production"},
		"beta":  {"disabled"},
	}
	got := f.Filter(hosts)
	if len(got) != 1 || got[0] != "alpha" {
		t.Errorf("expected [alpha], got %v", got)
	}
}
