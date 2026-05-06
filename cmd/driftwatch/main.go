package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/example/driftwatch/internal/config"
	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/notify"
	"github.com/example/driftwatch/internal/schedule"
)

var version = "dev"

func main() {
	configPath := flag.String("config", "driftwatch.yaml", "path to config file")
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("driftwatch %s\n", version)
		os.Exit(0)
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger := log.New(os.Stdout, "[driftwatch] ", log.LstdFlags)

	engine := drift.NewEngine(cfg)
	notifier := notify.New(cfg, os.Stdout)

	runner := func() {
		report, err := engine.Run()
		if err != nil {
			logger.Printf("engine error: %v", err)
			return
		}
		if err := notifier.Notify(report); err != nil {
			logger.Printf("notify error: %v", err)
		}
	}

	sched := schedule.New(cfg, logger)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	logger.Printf("starting driftwatch (interval: %s)", cfg.Interval)
	sched.Start(runner)

	<-sigs
	logger.Println("shutting down")
	sched.Stop()
}
