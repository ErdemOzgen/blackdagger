package cmd

import (
	"log"

	"github.com/ErdemOzgen/blackdagger/internal/config"
	"github.com/ErdemOzgen/blackdagger/internal/dag"
	"github.com/ErdemOzgen/blackdagger/internal/logger"
	"github.com/spf13/cobra"
)

func stopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop /path/to/spec.yaml",
		Short: "Stop the running workflow",
		Long:  `blackdagger stop /path/to/spec.yaml`,
		Args:  cobra.ExactArgs(1),
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

			workflow, err := dag.Load(cfg.BaseConfig, args[0], "")
			if err != nil {
				logger.Fatal("Workflow load failed", "error", err, "file", args[0])
			}

			logger.Info("Workflow stop initiated", "workflow", workflow.Name)

			dataStore := newDataStores(cfg)
			cli := newClient(cfg, dataStore, logger)

			if err := cli.Stop(workflow); err != nil {
				logger.Fatal(
					"Workflow stop operation failed",
					"error",
					err,
					"workflow",
					workflow.Name,
				)
			}
		},
	}
	cmd.Flags().BoolP("quiet", "q", false, "suppress output")
	return cmd
}
