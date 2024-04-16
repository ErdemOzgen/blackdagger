package cmd

import (
	"testing"

	"github.com/ErdemOzgen/blackdagger/internal/constants"
)

func TestVersionCommand(t *testing.T) {
	constants.Version = "1.0.5"
	testRunCommand(t, versionCmd(), cmdTest{
		args:        []string{"version"},
		expectedOut: []string{"1.0.5"}})
}
