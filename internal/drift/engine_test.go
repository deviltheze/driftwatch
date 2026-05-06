package drift

import (
	"testing"

	"github.com/user/driftwatch/internal/config"
)

func TestNewEngine_NotNil(t *testing.T) {
	cfg := &config.Config{
		SSHKeyPath: "/tmp/key",
		SSHPort:    22,
		Hosts:      []config.Host{},
	}
	checker := NewChecker(nil)
	reporter := NewReporter(nil)

	engine := NewEngine(cfg, checker, reporter)
	if engine == nil {
		t.Fatal("expected non-nil engine")
	}
}

func TestNewEngine_StoresConfig(t *testing.T) {
	cfg := &config.Config{
		SSHKeyPath: "/tmp/id_rsa",
		SSHPort:    2222,
		Hosts: []config.Host{
			{Address: "10.0.0.1", User: "admin"},
		},
	}
	checker := NewChecker(nil)
	reporter := NewReporter(nil)

	engine := NewEngine(cfg, checker, reporter)

	if engine.cfg.SSHPort != 2222 {
		t.Errorf("expected SSHPort 2222, got %d", engine.cfg.SSHPort)
	}
	if len(engine.cfg.Hosts) != 1 {
		t.Errorf("expected 1 host, got %d", len(engine.cfg.Hosts))
	}
}

func TestEngine_Run_NoHosts(t *testing.T) {
	cfg := &config.Config{
		SSHKeyPath: "/tmp/key",
		SSHPort:    22,
		Hosts:      []config.Host{},
	}
	checker := NewChecker(nil)

	var written []Report
	reporter := NewReporter(func(r Report) error {
		written = append(written, r)
		return nil
	})

	engine := NewEngine(cfg, checker, reporter)
	err := engine.Run()
	if err != nil {
		t.Fatalf("expected no error with empty hosts, got %v", err)
	}
	if len(written) != 1 {
		t.Errorf("expected reporter to be called once, got %d", len(written))
	}
}
