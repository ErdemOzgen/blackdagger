package cmd

import (
	"github.com/ErdemOzgen/blackdagger/internal/constants"
	"github.com/spf13/cobra"
)

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Display the binary version",
		Long:  `blackdagger version`,
		Run: func(_ *cobra.Command, _ []string) {
			println(constants.Version)
		},
	}
}
