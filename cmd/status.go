package cmd

import (
	"log"

	"github.com/ErdemOzgen/blackdagger/internal/config"
	"github.com/ErdemOzgen/blackdagger/internal/engine"
	"github.com/ErdemOzgen/blackdagger/internal/persistence/client"
	"github.com/ErdemOzgen/blackdagger/internal/persistence/model"
	"github.com/spf13/cobra"
)

func createStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status <DAG file>",
		Short: "Display current status of the DAG",
		Long:  `blackdagger status <DAG file>`,
		Args:  cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(config.LoadConfig(homeDir))
		},
		Run: func(cmd *cobra.Command, args []string) {
			loadedDAG, err := loadDAG(args[0], "")
			checkError(err)

			df := client.NewDataStoreFactory(config.Get())
			e := engine.NewFactory(df, config.Get()).Create()

			status, err := e.GetCurrentStatus(loadedDAG)
			checkError(err)

			res := &model.StatusResponse{Status: status}
			log.Printf("Pid=%d Status=%s", res.Status.Pid, res.Status.Status)
		},
	}
}
