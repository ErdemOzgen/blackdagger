package cmd

import (
	"github.com/ErdemOzgen/blackdagger/internal/constants"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	constants.Version = "1.2.3"
	testRunCommand(t, versionCmd(), cmdTest{
		args:        []string{"version"},
		expectedOut: []string{"1.2.3"}})
}
