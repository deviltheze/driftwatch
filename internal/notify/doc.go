// Package notify provides event notification for config drift detections.
//
// A Notifier formats drift reports into structured log lines and writes them
// to a configured io.Writer (defaulting to os.Stdout). Each event includes
// a timestamp, severity level (INFO or ALERT), the target host, and a summary
// of drifted vs total checked files.
//
// Usage:
//
//	n := notify.New(os.Stdout)
//	if err := n.Notify("web-01", report); err != nil {
//		log.Printf("notification failed: %v", err)
//	}
package notify
