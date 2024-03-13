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

const legacyPath = ".blackdagger"

func customHelpFunc(cmd *cobra.Command, strings []string) {
	fmt.Print(AsciiArt)
	fmt.Println(cmd.UsageString())
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	//fmt.Println(AsciiArt)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.blackdagger/admin.yaml)")

	cobra.OnInitialize(initConfig)
	rootCmd.SetHelpFunc(customHelpFunc)
	registerCommands(rootCmd)
}

var (
	homeDir string
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		cobra.CheckErr(err)
	}
	homeDir = home
}

func initConfig() {
	setConfigFile(homeDir)
}

func setConfigFile(home string) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		setDefaultConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("admin")
	}
}

func setDefaultConfigPath(home string) {
	viper.AddConfigPath(path.Join(home, legacyPath))
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
	rootCmd.AddCommand(createStatusCommand())
	rootCmd.AddCommand(versionCmd())
	rootCmd.AddCommand(serverCmd())
	rootCmd.AddCommand(createSchedulerCommand())
	rootCmd.AddCommand(retryCmd())
	rootCmd.AddCommand(startAllCmd())
}
