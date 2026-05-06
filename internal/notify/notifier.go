package notify

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/driftwatch/internal/drift"
)

// Level represents the severity of a notification.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Event represents a single drift notification event.
type Event struct {
	Timestamp time.Time
	Level     Level
	Host      string
	Report    drift.Report
}

// Notifier sends drift events to a configured output.
type Notifier struct {
	writer io.Writer
}

// New creates a Notifier that writes to the given writer.
// If writer is nil, os.Stdout is used.
func New(writer io.Writer) *Notifier {
	if writer == nil {
		writer = os.Stdout
	}
	return &Notifier{writer: writer}
}

// Notify formats and writes a drift event based on the report.
func (n *Notifier) Notify(host string, report drift.Report) error {
	event := Event{
		Timestamp: time.Now().UTC(),
		Host:      host,
		Report:    report,
	}

	if report.HasDrift() {
		event.Level = LevelAlert
	} else {
		event.Level = LevelInfo
	}

	return n.write(event)
}

func (n *Notifier) write(e Event) error {
	_, err := fmt.Fprintf(
		n.writer,
		"[%s] %s host=%s drifted=%d total=%d\n",
		e.Timestamp.Format(time.RFC3339),
		e.Level,
		e.Host,
		e.Report.DriftedCount(),
		len(e.Report.Results),
	)
	return err
}
