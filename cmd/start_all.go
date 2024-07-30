package cmd

import (
	"log"

	scheduler "github.com/ErdemOzgen/blackdagger/service"

	"github.com/ErdemOzgen/blackdagger/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func startAllCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-all",
		Short: "Launches both the Blackdagger web UI server and the scheduler process.",
		Long:  `blackdagger start-all [--dags=<DAGs dir>] [--host=<host>] [--port=<port>]`,
		PreRun: func(cmd *cobra.Command, args []string) {
			_ = viper.BindPFlag("port", cmd.Flags().Lookup("port"))
			_ = viper.BindPFlag("host", cmd.Flags().Lookup("host"))
			_ = viper.BindPFlag("dags", cmd.Flags().Lookup("dags"))
			cobra.CheckErr(config.LoadConfig())
		},
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			// TODO: move to config files
			pullDagList := []string{"default"}
			Pulldags(pullDagList)

			go func() {
				config.Get().DAGs = getFlagString(cmd, "dags", config.Get().DAGs)
				err := scheduler.New(topLevelModule).Start(cmd.Context())
				if err != nil {
					log.Fatal(err) // nolint // deep-exit
				}
			}()

			checkError(newFrontend().Start(ctx))
		},
	}
	bindStartAllCommandFlags(cmd)
	return cmd
}

func bindStartAllCommandFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("dags", "d", "", "location of DAG files (default is $HOME/.blackdagger/dags)")
	cmd.Flags().StringP("host", "s", "", "server host (default is localhost)")
	cmd.Flags().StringP("port", "p", "", "server port (default is 8080)")
}
