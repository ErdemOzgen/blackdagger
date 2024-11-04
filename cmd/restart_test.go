package cmd

import (
	"testing"
	"time"

	"github.com/ErdemOzgen/blackdagger/internal/dag"
	"github.com/ErdemOzgen/blackdagger/internal/dag/scheduler"
	"github.com/ErdemOzgen/blackdagger/internal/logger"
	"github.com/ErdemOzgen/blackdagger/internal/test"
	"github.com/stretchr/testify/require"
)

const (
	waitForStatusUpdate = time.Millisecond * 100
)

func TestRestartCommand(t *testing.T) {
	t.Run("RestartDAG", func(t *testing.T) {
		setup := test.SetupTest(t)
		defer setup.Cleanup()

		dagFile := testDAGFile("restart.yaml")

		// Start the DAG.
		go func() {
			testRunCommand(
				t,
				startCmd(),
				cmdTest{args: []string{"start", `--params="foo"`, dagFile}},
			)
		}()

		time.Sleep(waitForStatusUpdate)
		cli := setup.Client()

		// Wait for the DAG running.
		testStatusEventual(t, cli, dagFile, scheduler.StatusRunning)

		// Restart the DAG.
		done := make(chan struct{})
		go func() {
			testRunCommand(t, restartCmd(), cmdTest{args: []string{"restart", dagFile}})
			close(done)
		}()

		time.Sleep(waitForStatusUpdate)

		// Wait for the DAG running again.
		testStatusEventual(t, cli, dagFile, scheduler.StatusRunning)

		// Stop the restarted DAG.
		testRunCommand(t, stopCmd(), cmdTest{args: []string{"stop", dagFile}})

		time.Sleep(waitForStatusUpdate)

		// Wait for the DAG is stopped.
		testStatusEventual(t, cli, dagFile, scheduler.StatusNone)

		// Check parameter was the same as the first execution
		workflow, err := dag.Load(setup.Config.BaseConfig, dagFile, "")
		require.NoError(t, err)

		dataStore := newDataStores(setup.Config)
		recentHistory := newClient(
			setup.Config,
			dataStore,
			logger.Default,
		).GetRecentHistory(workflow, 2)

		require.Len(t, recentHistory, 2)
		require.Equal(t, recentHistory[0].Status.Params, recentHistory[1].Status.Params)

		<-done
	})
}
