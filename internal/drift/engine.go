package drift

import (
	"fmt"
	"log"

	"github.com/user/driftwatch/internal/config"
	"github.com/user/driftwatch/internal/ssh"
)

// Engine orchestrates drift checking across all configured hosts.
type Engine struct {
	cfg      *config.Config
	checker  *Checker
	reporter *Reporter
}

// NewEngine creates a new drift detection engine.
func NewEngine(cfg *config.Config, checker *Checker, reporter *Reporter) *Engine {
	return &Engine{
		cfg:      cfg,
		checker:  checker,
		reporter: reporter,
	}
}

// Run executes drift detection against all hosts defined in the config.
func (e *Engine) Run() error {
	var allResults []Result

	for _, host := range e.cfg.Hosts {
		results, err := e.checkHost(host)
		if err != nil {
			log.Printf("[warn] skipping host %s: %v", host.Address, err)
			continue
		}
		allResults = append(allResults, results...)
	}

	report := e.checker.BuildReport(allResults)
	return e.reporter.Write(report)
}

// checkHost connects to a single host and checks all configured files.
func (e *Engine) checkHost(host config.Host) ([]Result, error) {
	client, err := ssh.New(ssh.Config{
		Host:       host.Address,
		User:       host.User,
		KeyPath:    e.cfg.SSHKeyPath,
		Port:       e.cfg.SSHPort,
		TimeoutSec: e.cfg.TimeoutSec,
	})
	if err != nil {
		return nil, fmt.Errorf("ssh client init: %w", err)
	}

	runner := ssh.NewRunner(client)
	results := e.checker.CheckFiles(host, runner)
	return results, nil
}
