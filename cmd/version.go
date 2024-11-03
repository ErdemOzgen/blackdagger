package cmd

import (
	"fmt"

	"github.com/ErdemOzgen/blackdagger/internal/config"
	"github.com/ErdemOzgen/blackdagger/internal/constants"
	"github.com/spf13/cobra"
)

var AsciiArt = `

██████╗░██╗░░░░░░█████╗░░█████╗░██╗░░██╗██████╗░░█████╗░░██████╗░░██████╗░███████╗██████╗░
██╔══██╗██║░░░░░██╔══██╗██╔══██╗██║░██╔╝██╔══██╗██╔══██╗██╔════╝░██╔════╝░██╔════╝██╔══██╗
██████╦╝██║░░░░░███████║██║░░╚═╝█████═╝░██║░░██║███████║██║░░██╗░██║░░██╗░█████╗░░██████╔╝
██╔══██╗██║░░░░░██╔══██║██║░░██╗██╔═██╗░██║░░██║██╔══██║██║░░╚██╗██║░░╚██╗██╔══╝░░██╔══██╗
██████╦╝███████╗██║░░██║╚█████╔╝██║░╚██╗██████╔╝██║░░██║╚██████╔╝╚██████╔╝███████╗██║░░██║
╚═════╝░╚══════╝╚═╝░░╚═╝░╚════╝░╚═╝░░╚═╝╚═════╝░╚═╝░░╚═╝░╚═════╝░░╚═════╝░╚══════╝╚═╝░░╚═╝              
`

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Display the binary version",
		Long:  `blackdagger version`,
		PreRun: func(cmd *cobra.Command, args []string) {
			_, err := config.Load()
			cobra.CheckErr(err)
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(AsciiArt)
			fmt.Println(constants.Version)
		},
	}
}
