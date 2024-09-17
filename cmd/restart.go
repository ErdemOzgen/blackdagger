package cmd

import (
	"log"
	"time"

	"github.com/ErdemOzgen/blackdagger/internal/config"
	"github.com/ErdemOzgen/blackdagger/internal/dag"
	"github.com/ErdemOzgen/blackdagger/internal/engine"
	"github.com/ErdemOzgen/blackdagger/internal/persistence/client"
	"github.com/ErdemOzgen/blackdagger/internal/scheduler"
	"github.com/spf13/cobra"
)

func restartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "restart <DAG file>",
		Short: "Restart the DAG",
		Long:  `blackdagger restart <DAG file>`,
		Args:  cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(config.LoadConfig())
		},
		Run: func(cmd *cobra.Command, args []string) {
			dagFile := args[0]
			loadedDAG, err := loadDAG(dagFile, "")
			checkError(err)

			df := client.NewDataStoreFactory(config.Get())
			e := engine.NewFactory(df, config.Get()).Create()

			// Check the current status and stop the DAG if it is running.
			stopDAGIfRunning(e, loadedDAG)

			// Wait for the specified amount of time before restarting.
			waitForRestart(loadedDAG.RestartWait)

			// Retrieve the parameter of the previous execution.
			log.Printf("Restarting %s...", loadedDAG.Name)
			params := getPreviousExecutionParams(e, loadedDAG)

			// Start the DAG with the same parameter.
			loadedDAG, err = loadDAG(dagFile, params)
			checkError(err)
			cobra.CheckErr(start(cmd.Context(), e, loadedDAG, false))
		},
	}
}

func stopDAGIfRunning(e engine.Engine, d *dag.DAG) {
	st, err := e.GetCurrentStatus(d)
	checkError(err)

	// Stop the DAG if it is running.
	if st.Status == scheduler.StatusRunning {
		log.Printf("Stopping %s for restart...", d.Name)
		cobra.CheckErr(stopRunningDAG(e, d))
	}
}

func stopRunningDAG(e engine.Engine, d *dag.DAG) error {
	for {
		st, err := e.GetCurrentStatus(d)
		checkError(err)

		if st.Status != scheduler.StatusRunning {
			return nil
		}
		checkError(e.Stop(d))
		time.Sleep(time.Millisecond * 100)
	}
}

func waitForRestart(restartWait time.Duration) {
	if restartWait > 0 {
		log.Printf("Waiting for %s...", restartWait)
		time.Sleep(restartWait)
	}
}

func getPreviousExecutionParams(e engine.Engine, d *dag.DAG) string {
	st, err := e.GetLatestStatus(d)
	checkError(err)

	return st.Params
}
