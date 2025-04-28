package cmd

import (
	"log"

	"github.com/ErdemOzgen/blackdagger/internal/config"
	"github.com/ErdemOzgen/blackdagger/internal/dag"
	"github.com/ErdemOzgen/blackdagger/internal/dag/scheduler"
	"github.com/ErdemOzgen/blackdagger/internal/logger"
	"github.com/ErdemOzgen/blackdagger/internal/util"
	"github.com/spf13/cobra"
)

func stopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop [/path/to/spec.yaml]",
		Short: "Stop a running workflow or all running workflows if no file is specified",
		Long:  `blackdagger stop [/path/to/spec.yaml]`,
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.Load()
			if err != nil {
				log.Fatalf("Configuration load failed: %v", err)
			}

			quiet, err := cmd.Flags().GetBool("quiet")
			if err != nil {
				log.Fatalf("Flag retrieval failed (quiet): %v", err)
			}

			logger := logger.NewLogger(logger.NewLoggerArgs{
				Debug:  cfg.Debug,
				Format: cfg.LogFormat,
				Quiet:  quiet,
			})

			dataStore := newDataStores(cfg)
			cli := newClient(cfg, dataStore, logger)

			var dagFiles []string
			if len(args) == 0 {
				logger.Info("No workflow spec provided. Stopping all running workflows...")
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

				if curStatus.Status != scheduler.Status(scheduler.NodeStatusRunning) {
					logger.Info("Workflow not running, skipping", "workflow", workflow.Name, "status", curStatus.Status)
					continue
				}

				logger.Info("Stopping workflow", "workflow", workflow.Name)

				if err := cli.Stop(workflow); err != nil {
					logger.Fatal("Workflow stop operation failed", "error", err, "workflow", workflow.Name)
					continue
				}
			}
		},
	}
	cmd.Flags().BoolP("quiet", "q", false, "suppress output")
	return cmd
}
