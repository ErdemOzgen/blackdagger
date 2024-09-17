package cmd

import (
	"github.com/ErdemOzgen/blackdagger/internal/config"
	scheduler "github.com/ErdemOzgen/blackdagger/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func schedulerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scheduler",
		Short: "Start the scheduler",
		Long:  `blackdagger scheduler [--dags=<DAGs dir>]`,
		PreRun: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(config.LoadConfig())
		},
		Run: func(cmd *cobra.Command, args []string) {
			config.Get().DAGs = getFlagString(cmd, "dags", config.Get().DAGs)

			err := scheduler.New(topLevelModule).Start(cmd.Context())
			checkError(err)
		},
	}
	cmd.Flags().StringP("dags", "d", "", "location of DAG files (default is $HOME/.blackdagger/dags)")
	_ = viper.BindPFlag("dags", cmd.Flags().Lookup("dags"))

	return cmd
}
