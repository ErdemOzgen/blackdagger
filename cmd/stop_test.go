package cmd

import (
	"testing"
	"time"

	"github.com/ErdemOzgen/blackdagger/internal/dag/scheduler"
	"github.com/ErdemOzgen/blackdagger/internal/test"
)

func TestStopCommand(t *testing.T) {
	t.Run("StopDAG", func(t *testing.T) {
		setup := test.SetupTest(t)
		defer setup.Cleanup()

		dagFile := testDAGFile("long2.yaml")

		// Start the DAG.
		done := make(chan struct{})
		go func() {
			testRunCommand(t, startCmd(), cmdTest{args: []string{"start", dagFile}})
			close(done)
		}()

		time.Sleep(time.Millisecond * 100)

		// Wait for the DAG running.
		testLastStatusEventual(
			t,
			setup.DataStore().HistoryStore(),
			dagFile,
			scheduler.StatusRunning,
		)

		// Stop the DAG.
		testRunCommand(t, stopCmd(), cmdTest{
			args:        []string{"stop", dagFile},
			expectedOut: []string{"Stopping..."}})

		// Check the last execution is cancelled.
		testLastStatusEventual(
			t, setup.DataStore().HistoryStore(), dagFile, scheduler.StatusCancel,
		)
		<-done
	})
}
