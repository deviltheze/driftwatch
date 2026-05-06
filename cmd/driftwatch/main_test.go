package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestMain_VersionFlag(t *testing.T) {
	if os.Getenv("RUN_MAIN") == "1" {
		os.Args = []string{"driftwatch", "-version"}
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMain_VersionFlag")
	cmd.Env = append(os.Environ(), "RUN_MAIN=1")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("expected clean exit, got: %v", err)
	}
	got := string(out)
	if got == "" {
		t.Error("expected version output, got empty string")
	}
}

func TestMain_MissingConfig(t *testing.T) {
	if os.Getenv("RUN_MAIN") == "1" {
		os.Args = []string{"driftwatch", "-config", "/nonexistent/path.yaml"}
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMain_MissingConfig")
	cmd.Env = append(os.Environ(), "RUN_MAIN=1")
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit for missing config, got nil")
	}
}

func TestVersion_Default(t *testing.T) {
	if version != "dev" {
		t.Errorf("expected default version 'dev', got %q", version)
	}
}
