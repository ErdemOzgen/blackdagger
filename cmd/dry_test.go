package cmd

import (
	"testing"

	"github.com/ErdemOzgen/blackdagger/internal/test"
)

func TestDryCommand(t *testing.T) {
	t.Run("DryRun", func(t *testing.T) {
		setup := test.SetupTest(t)
		defer setup.Cleanup()

		tests := []cmdTest{
			{
				args:        []string{"dry", testDAGFile("success.yaml")},
				expectedOut: []string{"Starting DRY-RUN"},
			},
			{
				args:        []string{"dry", testDAGFile("params.yaml"), "--params", "p3 p4"},
				expectedOut: []string{`[1=p3 2=p4]`},
			},
		}

		for _, tc := range tests {
			testRunCommand(t, dryCmd(), tc)
		}
	})
}
