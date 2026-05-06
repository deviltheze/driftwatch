package alert

import "time"

// Suppression tracks recently emitted alerts to prevent duplicate noise.
type Suppression struct {
	cooldown time.Duration
	seen     map[string]time.Time
}

// NewSuppression creates a Suppression with the given cooldown duration.
// Alerts for the same host within the cooldown window are suppressed.
func NewSuppression(cooldown time.Duration) *Suppression {
	return &Suppression{
		cooldown: cooldown,
		seen:     make(map[string]time.Time),
	}
}

// Allow returns true if the alert should be emitted (not suppressed).
// It records the alert host and timestamp when allowed.
func (s *Suppression) Allow(a *Alert) bool {
	if a == nil {
		return false
	}
	last, exists := s.seen[a.Host]
	if exists && time.Since(last) < s.cooldown {
		return false
	}
	s.seen[a.Host] = time.Now().UTC()
	return true
}

// Reset clears all suppression state.
func (s *Suppression) Reset() {
	s.seen = make(map[string]time.Time)
}
