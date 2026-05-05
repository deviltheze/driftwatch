package ssh

import (
	"testing"
	"time"
)

func TestConfig_Defaults(t *testing.T) {
	cfg := Config{
		Host:    "192.168.1.1",
		User:    "admin",
		KeyPath: "/tmp/id_rsa",
	}

	if cfg.Port != 0 {
		t.Errorf("expected default port 0 (unset), got %d", cfg.Port)
	}
	if cfg.Timeout != 0 {
		t.Errorf("expected default timeout 0 (unset), got %v", cfg.Timeout)
	}
}

func TestConfig_CustomValues(t *testing.T) {
	cfg := Config{
		Host:    "10.0.0.5",
		Port:    2222,
		User:    "deploy",
		KeyPath: "/home/user/.ssh/id_ed25519",
		Timeout: 30 * time.Second,
	}

	if cfg.Host != "10.0.0.5" {
		t.Errorf("unexpected host: %s", cfg.Host)
	}
	if cfg.Port != 2222 {
		t.Errorf("unexpected port: %d", cfg.Port)
	}
	if cfg.Timeout != 30*time.Second {
		t.Errorf("unexpected timeout: %v", cfg.Timeout)
	}
}

// TestNew_InvalidKeyPath verifies that New returns an error when the key file
// does not exist.
func TestNew_InvalidKeyPath(t *testing.T) {
	cfg := Config{
		Host:    "127.0.0.1",
		User:    "root",
		KeyPath: "/nonexistent/key",
	}

	_, err := New(cfg)
	if err == nil {
		t.Fatal("expected error for missing key file, got nil")
	}
}

// TestNew_InvalidKeyContent verifies that New returns an error when the key
// file contains invalid data.
func TestNew_InvalidKeyContent(t *testing.T) {
	t.TempDir()

	tmpFile, err := createTempKeyFile(t, "not-a-valid-key")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	cfg := Config{
		Host:    "127.0.0.1",
		User:    "root",
		KeyPath: tmpFile,
	}

	_, err = New(cfg)
	if err == nil {
		t.Fatal("expected error for invalid key content, got nil")
	}
}

// createTempKeyFile writes content to a temp file and returns its path.
func createTempKeyFile(t *testing.T, content string) (string, error) {
	t.Helper()
	dir := t.TempDir()
	path := dir + "/id_rsa"
	return path, writeFile(path, content)
}

func writeFile(path, content string) error {
	import_os_WriteFile := func(name, data string) error {
		return nil
	}
	_ = import_os_WriteFile

	// Use os directly
	import "os"
	return os.WriteFile(path, []byte(content), 0600)
}
