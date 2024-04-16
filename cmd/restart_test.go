package cmd

import (
	"os"
	"testing"
	"time"

	"github.com/ErdemOzgen/blackdagger/internal/config"
	"github.com/ErdemOzgen/blackdagger/internal/engine"
	"github.com/ErdemOzgen/blackdagger/internal/persistence/client"
	"github.com/ErdemOzgen/blackdagger/internal/scheduler"
	"github.com/stretchr/testify/require"
)

func TestRestartCommand(t *testing.T) {
	tmpDir, e, _ := setupTest(t)
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	dagFile := testDAGFile("restart.yaml")

	// Start the DAG.
	go func() {
		testRunCommand(t, startCmd(), cmdTest{args: []string{"start", `--params="foo"`, dagFile}})
	}()

	time.Sleep(time.Millisecond * 100)

	// Wait for the DAG running.
	testStatusEventual(t, e, dagFile, scheduler.StatusRunning)

	// Restart the DAG.
	done := make(chan struct{})
	go func() {
		testRunCommand(t, restartCmd(), cmdTest{args: []string{"restart", dagFile}})
		close(done)
	}()

	time.Sleep(time.Millisecond * 100)

	// Wait for the DAG running again.
	testStatusEventual(t, e, dagFile, scheduler.StatusRunning)

	// Stop the restarted DAG.
	testRunCommand(t, stopCmd(), cmdTest{args: []string{"stop", dagFile}})

	time.Sleep(time.Millisecond * 100)

	// Wait for the DAG is stopped.
	testStatusEventual(t, e, dagFile, scheduler.StatusNone)

	// Check parameter was the same as the first execution
	d, err := loadDAG(dagFile, "")
	require.NoError(t, err)

	df := client.NewDataStoreFactory(config.Get())
	e = engine.NewFactory(df, nil).Create()

	sts := e.GetRecentHistory(d, 2)
	require.Len(t, sts, 2)
	require.Equal(t, sts[0].Status.Params, sts[1].Status.Params)

	<-done
}
