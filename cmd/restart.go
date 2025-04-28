package cmd

import (
	"log"
	"path/filepath"
	"time"

	"github.com/ErdemOzgen/blackdagger/internal/agent"
	"github.com/ErdemOzgen/blackdagger/internal/client"
	"github.com/ErdemOzgen/blackdagger/internal/config"
	"github.com/ErdemOzgen/blackdagger/internal/dag"
	"github.com/ErdemOzgen/blackdagger/internal/dag/scheduler"
	"github.com/ErdemOzgen/blackdagger/internal/logger"
	"github.com/ErdemOzgen/blackdagger/internal/util"
	"github.com/spf13/cobra"
)

func restartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restart [/path/to/spec.yaml]",
		Short: "Stop the running DAG and restart it",
		Long:  `blackdagger restart [/path/to/spec.yaml]`,
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

			initLogger := logger.NewLogger(logger.NewLoggerArgs{
				Debug:  cfg.Debug,
				Format: cfg.LogFormat,
				Quiet:  quiet,
			})

			dataStore := newDataStores(cfg)
			cli := newClient(cfg, dataStore, initLogger)

			var dagFiles []string

			if len(args) == 0 {
				initLogger.Info("Restarting all DAGs...")

				dagFiles, err = util.GetAllDAGFiles(cfg.DAGs)
				if err != nil {
					initLogger.Fatal("Failed to list DAG files", "error", err)
				}
			} else {
				dagFiles = []string{args[0]}
			}

			for _, specFilePath := range dagFiles {
				workflow, err := dag.Load(cfg.BaseConfig, specFilePath, "")
				if err != nil {
					initLogger.Error("Workflow load failed", "error", err, "file", specFilePath)
					continue
				}

				curStatus, err := cli.GetCurrentStatus(workflow)
				if err != nil {
					initLogger.Error("Status retrieval failed", "error", err, "workflow", workflow.Name)
					continue
				}

				// Skip workflows that are in NodeStatusNone
				if curStatus.Status == scheduler.StatusNone {
					initLogger.Info("Skipping DAG with no status", "workflow", workflow.Name)
					continue
				}

				initLogger.Info("Restarting workflow", "workflow", workflow.Name)

				if err := stopDAGIfRunning(cli, workflow, initLogger); err != nil {
					initLogger.Fatal("Workflow stop operation failed",
						"error", err,
						"workflow", workflow.Name)
					continue
				}

				// Wait for the specified amount of time before restarting.
				waitForRestart(workflow.RestartWait, initLogger)

				// Retrieve the parameter of the previous execution.
				params, err := getPreviousExecutionParams(cli, workflow)
				if err != nil {
					initLogger.Fatal("Previous execution parameter retrieval failed",
						"error", err,
						"workflow", workflow.Name)
					continue
				}

				// Start the DAG with the same parameter.
				// Need to reload the DAG file with the parameter.
				workflow, err = dag.Load(cfg.BaseConfig, specFilePath, params)
				if err != nil {
					initLogger.Fatal("Workflow reload failed",
						"error", err,
						"file", specFilePath,
						"params", params)
					continue
				}

				requestID, err := generateRequestID()
				if err != nil {
					initLogger.Fatal("Request ID generation failed", "error", err)
					continue
				}

				logFile, err := logger.OpenLogFile(logger.LogFileConfig{
					Prefix:    "restart_",
					LogDir:    cfg.LogDir,
					DAGLogDir: workflow.LogDir,
					DAGName:   workflow.Name,
					RequestID: requestID,
				})
				if err != nil {
					initLogger.Fatal("Log file creation failed",
						"error", err,
						"workflow", workflow.Name)
					continue
				}
				defer logFile.Close()

				agentLogger := logger.NewLogger(logger.NewLoggerArgs{
					Debug:   cfg.Debug,
					Format:  cfg.LogFormat,
					LogFile: logFile,
					Quiet:   quiet,
				})

				agentLogger.Info("Workflow restart initiated",
					"workflow", workflow.Name,
					"requestID", requestID,
					"logFile", logFile.Name())

				agt := agent.New(
					requestID,
					workflow,
					agentLogger,
					filepath.Dir(logFile.Name()),
					logFile.Name(),
					newClient(cfg, dataStore, agentLogger),
					dataStore,
					&agent.Options{Dry: false},
				)

				listenSignals(cmd.Context(), agt)
				if err := agt.Run(cmd.Context()); err != nil {
					agentLogger.Fatal("Workflow restart failed",
						"error", err,
						"workflow", workflow.Name,
						"requestID", requestID)
					continue
				}
			}
		},
	}
	cmd.Flags().BoolP("quiet", "q", false, "suppress output")
	return cmd
}

// stopDAGIfRunning stops the DAG if it is running.
// Otherwise, it does nothing.
func stopDAGIfRunning(e client.Client, workflow *dag.DAG, lg logger.Logger) error {
	curStatus, err := e.GetCurrentStatus(workflow)
	if err != nil {
		return err
	}

	if curStatus.Status == scheduler.StatusRunning {
		lg.Infof("Stopping: %s", workflow.Name)
		cobra.CheckErr(stopRunningDAG(e, workflow))
	}
	return nil
}

// stopRunningDAG attempts to stop the running DAG
// by sending a stop signal to the agent.
func stopRunningDAG(e client.Client, workflow *dag.DAG) error {
	for {
		curStatus, err := e.GetCurrentStatus(workflow)
		if err != nil {
			return err
		}

		// If the DAG is not running, do nothing.
		if curStatus.Status != scheduler.StatusRunning {
			return nil
		}

		if err := e.Stop(workflow); err != nil {
			return err
		}

		time.Sleep(time.Millisecond * 100)
	}
}

// waitForRestart waits for the specified amount of time before restarting
// the DAG.
func waitForRestart(restartWait time.Duration, lg logger.Logger) {
	if restartWait > 0 {
		lg.Info("Waiting for restart", "duration", restartWait)
		time.Sleep(restartWait)
	}
}

func getPreviousExecutionParams(e client.Client, workflow *dag.DAG) (string, error) {
	latestStatus, err := e.GetLatestStatus(workflow)
	if err != nil {
		return "", err
	}

	return latestStatus.Params, nil
}
