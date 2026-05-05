package ssh

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

// Client wraps an SSH client connection.
type Client struct {
	conn *ssh.Client
}

// Config holds SSH connection parameters.
type Config struct {
	Host       string
	Port       int
	User       string
	KeyPath    string
	Timeout    time.Duration
}

// New establishes an SSH connection using the provided config.
func New(cfg Config) (*Client, error) {
	key, err := os.ReadFile(cfg.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("reading key %q: %w", cfg.KeyPath, err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("parsing private key: %w", err)
	}

	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	port := cfg.Port
	if port == 0 {
		port = 22
	}

	clientCfg := &ssh.ClientConfig{
		User:            cfg.User,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: use known_hosts in production
		Timeout:         timeout,
	}

	addr := net.JoinHostPort(cfg.Host, fmt.Sprintf("%d", port))
	conn, err := ssh.Dial("tcp", addr, clientCfg)
	if err != nil {
		return nil, fmt.Errorf("dialing %s: %w", addr, err)
	}

	return &Client{conn: conn}, nil
}

// RunCommand executes a command on the remote host and returns its output.
func (c *Client) RunCommand(cmd string) (string, error) {
	session, err := c.conn.NewSession()
	if err != nil {
		return "", fmt.Errorf("creating session: %w", err)
	}
	defer session.Close()

	out, err := session.Output(cmd)
	if err != nil {
		return "", fmt.Errorf("running command %q: %w", cmd, err)
	}

	return string(out), nil
}

// Close terminates the underlying SSH connection.
func (c *Client) Close() error {
	return c.conn.Close()
}
