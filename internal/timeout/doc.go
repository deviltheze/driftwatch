// Package timeout provides deadline enforcement for SSH connect and command
// operations within driftwatch.
//
// An Enforcer is constructed from a Policy that specifies separate durations
// for the connection handshake and for individual remote commands. Both
// ConnectContext and CommandContext return a child context with the appropriate
// deadline attached, allowing callers to pass it directly to the SSH client.
//
// The Wrap helper executes an arbitrary function in a goroutine and returns a
// descriptive error if the command deadline is exceeded before the function
// completes, making it straightforward to add timeout semantics to any
// blocking SSH operation.
package timeout
