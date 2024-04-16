package cmd

import (
	"os"
	"testing"
	"time"

	"github.com/ErdemOzgen/blackdagger/internal/scheduler"
)

func TestStopCommand(t *testing.T) {
	tmpDir, _, ds := setupTest(t)
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	dagFile := testDAGFile("stop.yaml")

	// Start the DAG.
	done := make(chan struct{})
	go func() {
		testRunCommand(t, startCmd(), cmdTest{args: []string{"start", dagFile}})
		close(done)
	}()

	time.Sleep(time.Millisecond * 50)

	// Wait for the DAG running.
	// TODO: Do not use history store.
	testLastStatusEventual(t, ds.NewHistoryStore(), dagFile, scheduler.Status_Running)

	// Stop the DAG.
	testRunCommand(t, stopCmd(), cmdTest{
		args:        []string{"stop", dagFile},
		expectedOut: []string{"Stopping..."}})

	// Check the last execution is cancelled.
	// TODO: Do not use history store.
	testLastStatusEventual(t, ds.NewHistoryStore(), dagFile, scheduler.Status_Cancel)
	<-done
}
