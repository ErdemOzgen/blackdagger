package cmd

import (
	"log"

	"github.com/ErdemOzgen/blackdagger/internal/config"
	"github.com/ErdemOzgen/blackdagger/internal/dag"
	"github.com/ErdemOzgen/blackdagger/internal/logger"
	"github.com/spf13/cobra"
)

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status /path/to/spec.yaml",
		Short: "Display current status of the DAG",
		Long:  `blackdagger status /path/to/spec.yaml`,
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			cfg, err := config.Load()
			if err != nil {
				log.Fatalf("Configuration load failed: %v", err)
			}
			logger := logger.NewLogger(logger.NewLoggerArgs{
				Debug:  cfg.Debug,
				Format: cfg.LogFormat,
			})

			// Load the DAG file and get the current running status.
			workflow, err := dag.Load(cfg.BaseConfig, args[0], "")
			if err != nil {
				logger.Fatal("Workflow load failed", "error", err, "file", args[0])
			}

			dataStore := newDataStores(cfg)
			cli := newClient(cfg, dataStore, logger)

			curStatus, err := cli.GetCurrentStatus(workflow)

			if err != nil {
				logger.Fatal("Current status retrieval failed", "error", err)
			}

			logger.Info("Current status", "pid", curStatus.PID, "status", curStatus.Status)
		},
	}
}
