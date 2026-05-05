package ssh

import (
	"fmt"
	"strings"
)

// CommandResult holds the output of a remote command.
type CommandResult struct {
	Host    string
	Command string
	Output  string
	Err     error
}

// Runner executes a set of commands against a remote host.
type Runner struct {
	client *Client
	host   string
}

// NewRunner creates a Runner using an established Client.
func NewRunner(client *Client, host string) *Runner {
	return &Runner{client: client, host: host}
}

// Run executes each command and returns a slice of results.
func (r *Runner) Run(commands []string) []CommandResult {
	results := make([]CommandResult, 0, len(commands))
	for _, cmd := range commands {
		out, err := r.client.RunCommand(cmd)
		results = append(results, CommandResult{
			Host:    r.host,
			Command: cmd,
			Output:  strings.TrimSpace(out),
			Err:     err,
		})
	}
	return results
}

// RunSingle executes a single command and returns the result.
func (r *Runner) RunSingle(cmd string) CommandResult {
	out, err := r.client.RunCommand(cmd)
	return CommandResult{
		Host:    r.host,
		Command: cmd,
		Output:  strings.TrimSpace(out),
		Err:     err,
	}
}

// String formats a CommandResult for human-readable output.
func (cr CommandResult) String() string {
	if cr.Err != nil {
		return fmt.Sprintf("[%s] $ %s\nERROR: %v", cr.Host, cr.Command, cr.Err)
	}
	return fmt.Sprintf("[%s] $ %s\n%s", cr.Host, cr.Command, cr.Output)
}
