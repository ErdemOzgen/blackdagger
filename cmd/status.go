package cmd

import (
	"log"

	"github.com/ErdemOzgen/blackdagger/internal/config"
	"github.com/ErdemOzgen/blackdagger/internal/dag"
	"github.com/ErdemOzgen/blackdagger/internal/logger"
	"github.com/ErdemOzgen/blackdagger/internal/util"
	"github.com/spf13/cobra"
)

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status [/path/to/spec.yaml]",
		Short: "Display current status of the DAG(s)",
		Long:  `blackdagger status [/path/to/spec.yaml]`,
		Args:  cobra.MaximumNArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			cfg, err := config.Load()
			if err != nil {
				log.Fatalf("Configuration load failed: %v", err)
			}
			logger := logger.NewLogger(logger.NewLoggerArgs{
				Debug:  cfg.Debug,
				Format: cfg.LogFormat,
			})

			dataStore := newDataStores(cfg)
			cli := newClient(cfg, dataStore, logger)

			var dagFiles []string
			if len(args) == 0 {
				logger.Info("Checking status of all DAGs...")
				dagFiles, err = util.GetAllDAGFiles(cfg.DAGs)
				if err != nil {
					logger.Fatal("DAG file list retrieval failed", "error", err)
				}
			} else {
				dagFiles = []string{args[0]}
			}

			for _, file := range dagFiles {
				workflow, err := dag.Load(cfg.BaseConfig, file, "")
				if err != nil {
					logger.Error("Workflow load failed", "error", err, "file", file)
					continue
				}

				curStatus, err := cli.GetCurrentStatus(workflow)
				if err != nil {
					logger.Error("Current status retrieval failed", "error", err, "workflow", workflow.Name)
					continue
				}

				logger.Info(
					"Workflow status",
					"workflow", workflow.Name,
					"pid", curStatus.PID,
					"status", curStatus.Status,
				)
			}
		},
	}
}
