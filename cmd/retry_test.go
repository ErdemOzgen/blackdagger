package cmd

import (
	"fmt"
	"testing"

	"github.com/ErdemOzgen/blackdagger/internal/dag/scheduler"
	"github.com/ErdemOzgen/blackdagger/internal/test"
	"github.com/stretchr/testify/require"
)

func TestRetryCommand(t *testing.T) {
	t.Run("RetryDAG", func(t *testing.T) {
		setup := test.SetupTest(t)
		defer setup.Cleanup()

		dagFile := testDAGFile("retry.yaml")

		// Run a DAG.
		testRunCommand(t, startCmd(), cmdTest{args: []string{"start", `--params="foo"`, dagFile}})

		// Find the request ID.
		cli := setup.Client()
		status, err := cli.GetStatus(dagFile)
		require.NoError(t, err)
		require.Equal(t, status.Status.Status, scheduler.StatusSuccess)
		require.NotNil(t, status.Status)

		requestID := status.Status.RequestID

		// Retry with the request ID.
		testRunCommand(t, retryCmd(), cmdTest{
			args:        []string{"retry", fmt.Sprintf("--req=%s", requestID), dagFile},
			expectedOut: []string{"param is foo"},
		})
	})
}
