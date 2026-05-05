package config

import (
	"os"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "driftwatch-*.yaml")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	content := `
interval: 30s
hosts:
  - name: web-01
    address: 192.168.1.10
    user: deploy
    key_path: ~/.ssh/id_rsa
checks:
  - name: nginx-config
    path: /etc/nginx/nginx.conf
    hosts: [web-01]
`
	path := writeTempConfig(t, content)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if cfg.Interval != 30*time.Second {
		t.Errorf("expected interval 30s, got %v", cfg.Interval)
	}
	if len(cfg.Hosts) != 1 || cfg.Hosts[0].Name != "web-01" {
		t.Errorf("unexpected hosts: %+v", cfg.Hosts)
	}
	if cfg.Hosts[0].Port != 22 {
		t.Errorf("expected default port 22, got %d", cfg.Hosts[0].Port)
	}
	if len(cfg.Checks) != 1 || cfg.Checks[0].Name != "nginx-config" {
		t.Errorf("unexpected checks: %+v", cfg.Checks)
	}
}

func TestLoad_DefaultInterval(t *testing.T) {
	content := `
hosts:
  - name: db-01
    address: 10.0.0.5
    user: admin
checks: []
`
	path := writeTempConfig(t, content)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if cfg.Interval != 60*time.Second {
		t.Errorf("expected default interval 60s, got %v", cfg.Interval)
	}
}

func TestLoad_MissingHostAddress(t *testing.T) {
	content := `
hosts:
  - name: bad-host
    user: root
checks: []
`
	path := writeTempConfig(t, content)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing address")
	}
}

func TestLoad_NoHosts(t *testing.T) {
	content := `checks: []`
	path := writeTempConfig(t, content)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error when no hosts defined")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
