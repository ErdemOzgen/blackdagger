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
			{
				args:        []string{"dry", testDAGFile("params.yaml"), "--", "p5", "p6"},
				expectedOut: []string{`[1=p5 2=p6]`},
			},
		}

		for _, tc := range tests {
			testRunCommand(t, dryCmd(), tc)
		}
	})
}
