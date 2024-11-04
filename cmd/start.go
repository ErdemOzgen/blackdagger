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

func startCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start [flags] /path/to/spec.yaml",
		Short: "Runs the DAG",
		Long:  `blackdagger start [--params="param1 param2"] /path/to/spec.yaml`,
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

			initLogger := logger.NewLogger(logger.NewLoggerArgs{
				Debug:  cfg.Debug,
				Format: cfg.LogFormat,
				Quiet:  quiet,
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
				Prefix:    "start_",
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
				Quiet:   quiet,
			})

			dataStore := newDataStores(cfg)
			cli := newClient(cfg, dataStore, agentLogger)

			agentLogger.Info("Workflow execution initiated",
				"workflow", workflow.Name,
				"requestID", requestID,
				"logFile", logFile.Name())

			agt := agent.New(
				requestID,
				workflow,
				agentLogger,
				filepath.Dir(logFile.Name()),
				logFile.Name(),
				cli,
				dataStore,
				&agent.Options{})

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
	cmd.Flags().BoolP("quiet", "q", false, "suppress output")
	return cmd
}

// removeQuotes removes the surrounding quotes from the string.
func removeQuotes(s string) string {
	if len(s) > 1 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}
