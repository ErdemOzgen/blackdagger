package cmd

import (
	"log"
	"path/filepath"

	"github.com/ErdemOzgen/blackdagger/internal/agent"
	"github.com/ErdemOzgen/blackdagger/internal/config"
	"github.com/ErdemOzgen/blackdagger/internal/dag"
	"github.com/ErdemOzgen/blackdagger/internal/logger"
	"github.com/spf13/cobra"
)

func dryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dry [flags] /path/to/spec.yaml",
		Short: "Dry-runs specified DAG",
		Long:  `blackdagger dry [--params="param1 param2"] /path/to/spec.yaml`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.Load()
			if err != nil {
				log.Fatalf("Failed to load config: %v", err)
			}
			initLogger := logger.NewLogger(logger.NewLoggerArgs{
				Debug:  cfg.Debug,
				Format: cfg.LogFormat,
			})

			params, err := cmd.Flags().GetString("params")
			if err != nil {
				initLogger.Fatal("Parameter retrieval failed", "error", err)
			}

			workflow, err := dag.Load(cfg.BaseConfig, args[0], removeQuotes(params))
			if err != nil {
				initLogger.Fatal("Workflow load failed", "error", err, "file", args[0])
			}

			requestID, err := generateRequestID()
			if err != nil {
				initLogger.Fatal("Request ID generation failed", "error", err)
			}

			logFile, err := logger.OpenLogFile(logger.LogFileConfig{
				Prefix:    "dry_",
				LogDir:    cfg.LogDir,
				DAGLogDir: workflow.LogDir,
				DAGName:   workflow.Name,
				RequestID: requestID,
			})

			if err != nil {
				initLogger.Fatal(
					"Log file creation failed",
					"error",
					err,
					"workflow",
					workflow.Name,
				)
			}
			defer logFile.Close()

			agentLogger := logger.NewLogger(logger.NewLoggerArgs{
				Debug:   cfg.Debug,
				Format:  cfg.LogFormat,
				LogFile: logFile,
			})

			dataStore := newDataStores(cfg)
			cli := newClient(cfg, dataStore, agentLogger)

			agt := agent.New(
				requestID,
				workflow,
				agentLogger,
				filepath.Dir(logFile.Name()),
				logFile.Name(),
				cli,
				dataStore,
				&agent.Options{Dry: true})

			ctx := cmd.Context()

			listenSignals(ctx, agt)

			if err := agt.Run(ctx); err != nil {
				agentLogger.Fatal("Workflow execution failed",
					"error", err,
					"workflow", workflow.Name,
					"requestID", requestID)
			}
		},
	}

	cmd.Flags().StringP("params", "p", "", "parameters")
	return cmd
}
