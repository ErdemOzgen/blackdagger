package cmd

import (
	"github.com/ErdemOzgen/blackdagger/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func serverCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start the server",
		Long:  `blackdagger server [--dags=<DAGs dir>] [--host=<host>] [--port=<port>]`,
		PreRun: func(cmd *cobra.Command, args []string) {
			_ = viper.BindPFlag("port", cmd.Flags().Lookup("port"))
			_ = viper.BindPFlag("host", cmd.Flags().Lookup("host"))
			_ = viper.BindPFlag("dags", cmd.Flags().Lookup("dags"))
			cobra.CheckErr(config.LoadConfig())
		},
		Run: func(cmd *cobra.Command, args []string) {
			checkError(newFrontend().Start(cmd.Context()))
		},
	}
	bindServerCommandFlags(cmd)
	return cmd
}

func bindServerCommandFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("dags", "d", "", "location of DAG files (default is $HOME/.blackdagger/dags)")
	cmd.Flags().StringP("host", "s", "", "server host (default is localhost)")
	cmd.Flags().StringP("port", "p", "", "server port (default is 8080)")
}
