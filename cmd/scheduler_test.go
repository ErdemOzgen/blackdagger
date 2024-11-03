package cmd

import (
	"testing"
	"time"

	"github.com/ErdemOzgen/blackdagger/internal/test"
)

func TestSchedulerCommand(t *testing.T) {
	t.Run("StartScheduler", func(t *testing.T) {
		setup := test.SetupTest(t)
		defer setup.Cleanup()

		go func() {
			testRunCommand(t, schedulerCmd(), cmdTest{
				args:        []string{"scheduler"},
				expectedOut: []string{"starting blackdagger scheduler"},
			})
		}()

		time.Sleep(time.Millisecond * 500)
	})
}
