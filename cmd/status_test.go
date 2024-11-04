package cmd

import (
	"testing"

	"github.com/ErdemOzgen/blackdagger/internal/dag/scheduler"
	"github.com/ErdemOzgen/blackdagger/internal/test"
)

func TestStatusCommand(t *testing.T) {
	t.Run("StatusDAG", func(t *testing.T) {
		setup := test.SetupTest(t)
		defer setup.Cleanup()

		dagFile := testDAGFile("long.yaml")

		// Start the DAG.
		done := make(chan struct{})
		go func() {
			testRunCommand(t, startCmd(), cmdTest{args: []string{"start", dagFile}})
			close(done)
		}()

		testLastStatusEventual(
			t,
			setup.DataStore().HistoryStore(),
			dagFile,
			scheduler.StatusRunning,
		)

		// Check the current status.
		testRunCommand(t, statusCmd(), cmdTest{
			args:        []string{"status", dagFile},
			expectedOut: []string{"Status=running"},
		})

		// Stop the DAG.
		testRunCommand(t, stopCmd(), cmdTest{args: []string{"stop", dagFile}})
		<-done
	})
}
