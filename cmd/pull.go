/*
Copyright © 2024 ErdemOzgen m.erdemozgen@gmail.com
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ErdemOzgen/blackdagger/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Get YAML files from remote repository for collection of jobs (e.g. default,devops,cart,mlops,etc.)",
	Long: `Get YAML files from remote repository for collection of jobs (e.g. default,devops,cart,mlops,etc.)
		Example: blackdagger pull default
		Example: blackdagger pull devops 
		Example: blackdagger pull cart
		Example: blackdagger pull mlops
		Example: blackdagger pull all
		
		see
		`,
	PreRun: func(cmd *cobra.Command, args []string) {
		_ = viper.BindPFlag("dags", cmd.Flags().Lookup("dags"))       // /home/erdem/.blackdagger/dags
		_ = viper.BindPFlag("logDir", cmd.Flags().Lookup("logDir"))   // /home/erdem/.blackdagger/logs
		_ = viper.BindPFlag("dataDir", cmd.Flags().Lookup("dataDir")) // /home/erdem/.blackdagger/data
		cobra.CheckErr(config.LoadConfig(homeDir))
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting to pull the repository...")
		pulldags(args)
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pullCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pullCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func pulldags(args []string) {
	// TODO: move to config files
	jobCategories := []string{"mlops", "default", "devsecops", "devops"}
	if len(args) > 0 {
		if checkCategory(jobCategories, args[0]) {
			repoName := args[0]                 // Use the first argument as the repository name, e.g., "default"
			viper.AutomaticEnv()                // Ensure Viper is looking for env vars
			dagValue := viper.GetString("dags") // Get the DAG config value
			fmt.Printf("DAG configuration value: %s\n", dagValue)

			// Define the path for the folder based on the argument
			folderPath := filepath.Join(dagValue, repoName)

			// Check if the directory exists
			if _, err := os.Stat(folderPath); os.IsNotExist(err) {
				// Directory does not exist, clone the repository
				// TODO: move to config files
				gitCmd := exec.Command("git", "clone", fmt.Sprintf("https://github.com/ErdemOzgen/blackdagger-%s.git", repoName), folderPath)
				if output, err := gitCmd.CombinedOutput(); err != nil {
					fmt.Printf("Failed to clone the repository: %v, output: %s\n", err, string(output))
					fmt.Println("This category may not been public yet. Stay tuned for updates!")
				} else {
					fmt.Printf("Successfully cloned the repository into %s\n", folderPath)
					// Copy the YAML files to the DAGs folder
					CopyYAMLFiles(folderPath, dagValue)
				}
			} else {
				// Directory exists, pull the repository
				fmt.Printf("The directory %s already exists. Pulling updates...\n", folderPath)

				// Change the working directory to the repository's folder
				os.Chdir(folderPath)

				// Execute git pull to update the repository
				gitCmd := exec.Command("git", "pull")
				if output, err := gitCmd.CombinedOutput(); err != nil {
					fmt.Printf("Failed to pull the repository: %v, output: %s\n", err, string(output))
				} else {
					fmt.Printf("Successfully updated the repository in %s\n", folderPath)
					// Copy the YAML files to the DAGs folder
					CopyYAMLFiles(folderPath, dagValue)
				}
			}
		} else {
			fmt.Println("Please specify correct category name.")
			fmt.Println("Available categories are: ", jobCategories)
		}
	} else {
		fmt.Println("Please specify a category name.")
		fmt.Println("Available categories are: ", jobCategories)
	}

}

// checkCategory checks if a string is present in a category.
// It returns true if the string is found, otherwise false.
func checkCategory(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

// CopyYAMLFiles copies all YAML files from srcDir to destDir.
func CopyYAMLFiles(srcDir, destDir string) {
	files, err := os.ReadDir(srcDir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue // Skip directories
		}

		if strings.HasSuffix(file.Name(), ".yaml") || strings.HasSuffix(file.Name(), ".yml") {
			srcPath := filepath.Join(srcDir, file.Name())
			destPath := filepath.Join(destDir, file.Name())

			// Open the source file
			srcFile, err := os.Open(srcPath)
			if err != nil {
				panic(err)
			}
			defer srcFile.Close()

			// Create the destination file
			destFile, err := os.Create(destPath)
			if err != nil {
				panic(err)
			}
			defer destFile.Close()

			// Copy the contents of the source file to the destination file
			_, err = io.Copy(destFile, srcFile)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Copied %s to %s\n", srcPath, destPath)
		}
	}
}
