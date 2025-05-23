package cmd

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/ErdemOzgen/blackdagger/internal/config"
	"github.com/ErdemOzgen/blackdagger/internal/dag"
	"github.com/ErdemOzgen/blackdagger/internal/persistence"

	"github.com/ErdemOzgen/blackdagger/internal/client"
	"github.com/ErdemOzgen/blackdagger/internal/dag/scheduler"
	"github.com/ErdemOzgen/blackdagger/internal/util"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

// cmdTest is a helper struct to test commands.
// It contains the arguments to the command and the expected output.
type cmdTest struct {
	args        []string
	expectedOut []string
}

// testRunCommand is a helper function to test a command.
func testRunCommand(t *testing.T, cmd *cobra.Command, test cmdTest) {
	t.Helper()

	root := &cobra.Command{Use: "root"}
	root.AddCommand(cmd)

	// Set arguments.
	root.SetArgs(test.args)

	// Run the command

	// TODO: Fix thet test after update the logging code so that it can be
	err := root.Execute()
	require.NoError(t, err)

}

func testDAGFile(name string) string {
	return filepath.Join(
		filepath.Join(util.MustGetwd(), "testdata"),
		name,
	)
}

const (
	waitForStatusTimeout = time.Millisecond * 5000
	tick                 = time.Millisecond * 50
)

// testStatusEventual tests the status of a DAG to be the expected status.
func testStatusEventual(t *testing.T, e client.Client, dagFile string, expected scheduler.Status) {
	t.Helper()

	cfg, err := config.Load()
	require.NoError(t, err)

	workflow, err := dag.Load(cfg.BaseConfig, dagFile, "")
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		status, err := e.GetCurrentStatus(workflow)
		require.NoError(t, err)
		return expected == status.Status
	}, waitForStatusTimeout, tick)
}

// testLastStatusEventual tests the last status of a DAG to be the expected status.
func testLastStatusEventual(
	t *testing.T,
	hs persistence.HistoryStore,
	dg string,
	expected scheduler.Status,
) {
	t.Helper()

	require.Eventually(t, func() bool {
		status := hs.ReadStatusRecent(dg, 1)
		if len(status) < 1 {
			return false
		}
		return expected == status[0].Status.Status
	}, waitForStatusTimeout, tick)
}
