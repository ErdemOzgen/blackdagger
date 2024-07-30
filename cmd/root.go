/*
Copyright © 2024 ErdemOzgen m.erdemozgen@gmail.com
*/
package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/ErdemOzgen/blackdagger/internal/config"
	"github.com/ErdemOzgen/blackdagger/internal/dag"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "blackdagger",
		Short: "YAML-based DAG scheduling tool for Red teaming,CART,DevOps,DevSecOps,MLOps,MLSecOps.",
		Long:  `YAML-based DAG scheduling tool for Red teaming,CART,DevOps,DevSecOps,MLOps,MLSecOps.`,
	}
)

func customHelpFunc(cmd *cobra.Command, strings []string) {
	fmt.Print(AsciiArt)
	fmt.Println(cmd.UsageString())
}

const configPath = ".blackdagger"

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.blackdagger/admin.yaml)")

	rootCmd.SetHelpFunc(customHelpFunc)
	cobra.OnInitialize(initialize)
	registerCommands(rootCmd)
}

func init() {
	_, err := os.UserHomeDir()
	if err != nil {
		cobra.CheckErr(err)
	}
}

func initialize() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		setDefaultConfigPath()
		viper.SetConfigType("yaml")
		viper.SetConfigName("admin")
	}
}

func setDefaultConfigPath() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("could not determine home directory")
	}
	viper.AddConfigPath(path.Join(homeDir, configPath))
}

func loadDAG(dagFile, params string) (d *dag.DAG, err error) {
	dagLoader := &dag.Loader{BaseConfig: config.Get().BaseConfig}
	return dagLoader.Load(dagFile, params)
}

func getFlagString(cmd *cobra.Command, name, fallback string) string {
	if s, _ := cmd.Flags().GetString(name); s != "" {
		return s
	}
	return fallback
}

func registerCommands(root *cobra.Command) {
	rootCmd.AddCommand(startCmd())
	rootCmd.AddCommand(stopCmd())
	rootCmd.AddCommand(restartCmd())
	rootCmd.AddCommand(dryCmd())
	rootCmd.AddCommand(statusCmd())
	rootCmd.AddCommand(versionCmd())
	rootCmd.AddCommand(serverCmd())
	rootCmd.AddCommand(schedulerCmd())
	rootCmd.AddCommand(retryCmd())
	rootCmd.AddCommand(startAllCmd())
}
