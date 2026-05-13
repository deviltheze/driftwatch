// Package audit provides structured audit logging for drift detection events.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// EventType classifies the kind of audit event.
type EventType string

const (
	EventCheckStarted  EventType = "check.started"
	EventCheckFinished EventType = "check.finished"
	EventDriftDetected EventType = "drift.detected"
	EventAlertSent     EventType = "alert.sent"
)

// Event represents a single audit log entry.
type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Type      EventType `json:"type"`
	Host      string    `json:"host,omitempty"`
	Message   string    `json:"message"`
	Meta      map[string]string `json:"meta,omitempty"`
}

// Logger writes audit events as newline-delimited JSON.
type Logger struct {
	w io.Writer
}

// New creates a new audit Logger writing to w.
// If w is nil, os.Stdout is used.
func New(w io.Writer) *Logger {
	if w == nil {
		w = os.Stdout
	}
	return &Logger{w: w}
}

// Log writes a single audit event.
func (l *Logger) Log(eventType EventType, host, message string, meta map[string]string) error {
	e := Event{
		Timestamp: time.Now().UTC(),
		Type:      eventType,
		Host:      host,
		Message:   message,
		Meta:      meta,
	}
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal event: %w", err)
	}
	_, err = fmt.Fprintf(l.w, "%s\n", b)
	return err
}
