/*
Copyright © 2024 ErdemOzgen m.erdemozgen@gmail.com
*/
package cmd

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
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
		Example: blackdagger pull mlsecops
		Example: blackdagger pull devsecops
		Example: blackdagger pull <category>
		
		see
		`,
	PreRun: func(cmd *cobra.Command, args []string) {
		_ = viper.BindPFlag("dags", cmd.Flags().Lookup("dags"))
		_ = viper.BindPFlag("logDir", cmd.Flags().Lookup("logDir"))
		_ = viper.BindPFlag("dataDir", cmd.Flags().Lookup("dataDir"))
		_ = viper.BindPFlag("force", rootCmd.PersistentFlags().Lookup("force"))
		_ = viper.BindPFlag("keep", rootCmd.PersistentFlags().Lookup("keep"))
		_, err := config.Load()
		cobra.CheckErr(err)

		force := viper.GetBool("force")
		keep := viper.GetBool("keep")
		if force && keep {
			fmt.Println("Error: --force and --keep cannot be used together.")
			os.Exit(1)
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting to pull the repository...")

		Pulldags(args)
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)

	pullCmd.Flags().String("origin", "", "Custom origin URL to pull from")
	_ = viper.BindPFlag("origin", pullCmd.Flags().Lookup("origin"))

	pullCmd.Flags().String("folder", "", "Folder to pull into")
	_ = viper.BindPFlag("folder", pullCmd.Flags().Lookup("folder"))

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pullCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pullCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func Pulldags(args []string) {
	// TODO: move to config files
	jobCategories := []string{"mlops", "default", "devsecops", "devops", "mlsecops", "cart"}

	if len(args) == 0 && viper.GetString("origin") == "" {
		fmt.Println("Please specify a category name or provide --origin.")
		fmt.Println("Available categories are:", jobCategories)
		return
	}

	viper.AutomaticEnv()
	dagValue := viper.GetString("dags")
	origin := viper.GetString("origin")
	folder := viper.GetString("folder")

	var repoName, repoURL, folderName string

	if origin != "" {
		repoURL = origin
		split := strings.Split(strings.TrimSuffix(filepath.Base(origin), ".git"), "/")
		folderName = split[len(split)-1]
	} else if CheckCategory(jobCategories, args[0]) {
		repoName = args[0]
		repoURL = fmt.Sprintf("%s%s.git", viper.GetString("dagRepo"), repoName)
		folderName = repoName
	} else {
		fmt.Println("Please specify a correct category name.")
		fmt.Println("Available categories are:", jobCategories)
		return
	}

	max := big.NewInt(1000000)
	randomInt, err := rand.Int(rand.Reader, max)
	if err != nil {
		fmt.Printf("Failed to generate a random number: %v\n", err)
		return
	}
	tempFolderName := fmt.Sprintf("blackdagger-%s-%s", folderName, randomInt.String())

	repoBase := filepath.Join(dagValue, "repos", tempFolderName)
	destBase := filepath.Join(dagValue, folder)

	if err := os.MkdirAll(repoBase, os.ModePerm); err != nil {
		fmt.Printf("Failed to create temp repo folder: %v\n", err)
		return
	}
	if err := os.MkdirAll(destBase, os.ModePerm); err != nil {
		fmt.Printf("Failed to create dag folder: %v\n", err)
		return
	}

	fmt.Printf("Cloning repository from %s into %s...\n", repoURL, repoBase)
	prefix := fmt.Sprintf("blackdagger-%s", folderName)

	matches, err := filepath.Glob(filepath.Join(filepath.Join(dagValue, "repos"), prefix+"-*"))
	if err != nil {
		fmt.Printf("Failed to find existing folders: %v\n", err)
		return
	}
	for _, match := range matches {
		if stat, err := os.Stat(match); err == nil && stat.IsDir() {
			if err := os.RemoveAll(match); err != nil {
				fmt.Printf("Failed to remove existing folder %s: %v\n", match, err)
				return
			}
		}
	}

	gitCmd := exec.Command("git", "clone", repoURL, repoBase)
	if output, err := gitCmd.CombinedOutput(); err != nil {
		fmt.Printf("Failed to clone the repository: %v, output: %s\n", err, string(output))
		fmt.Println("This category may not been public yet. Stay tuned for updates!")

		return
	}

	CopyWithChangeDetection(repoBase, destBase)

	_ = os.RemoveAll(repoBase)
}

func CopyWithChangeDetection(srcDir, destDir string) {
	files, err := os.ReadDir(srcDir)
	if err != nil {
		panic(err)
	}

	if viper.GetBool("force") {
		fmt.Println("Overwriting local changes due to --force flag.")
	} else if viper.GetBool("keep") {
		fmt.Println("Keeping all local changes due to --keep flag.")
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if strings.HasSuffix(file.Name(), ".yaml") || strings.HasSuffix(file.Name(), ".yml") {
			srcPath := filepath.Join(srcDir, file.Name())
			destPath := filepath.Join(destDir, file.Name())

			if viper.GetBool("force") {
				copyFile(srcPath, destPath)
				fmt.Printf("Overwritten %s\n", destPath)
				continue
			}

			// Check if destination exists
			if _, err := os.Stat(destPath); err == nil {
				// File exists: check if changed
				same, err := filesAreEqual(srcPath, destPath)
				if err != nil {
					fmt.Printf("Error comparing files: %v\n", err)
					continue
				}
				if !same {
					fmt.Printf("Detected change in %s\n", file.Name())
					fmt.Print("Overwrite? [y]es / [n]o / [d]iff: ")
					var choice string
					fmt.Scanln(&choice)
					switch strings.ToLower(choice) {
					case "y":
						copyFile(srcPath, destPath)
						fmt.Printf("Overwritten %s\n", destPath)
					case "d":
						showDiff(srcPath, destPath)
						fmt.Printf("-----  -----\n\n")
						fmt.Print("Revert after viewing? [y/N]: ")
						fmt.Scanln(&choice)
						if strings.ToLower(choice) == "y" {
							copyFile(srcPath, destPath)
							fmt.Printf("Copied new file: %s\n", destPath)
						} else {
							fmt.Println("Skipped.")
						}
					default:
						fmt.Println("Skipped.")
					}
				} else {
					fmt.Printf("No changes in %s\n", file.Name())
				}
			} else {
				fmt.Printf("New file detected: %s\n", file.Name())
				fmt.Print("Copy to destination? [y]es / [n]o / [v]iew: ")
				var choice string
				fmt.Scanln(&choice)
				switch strings.ToLower(choice) {
				case "y":
					copyFile(srcPath, destPath)
					fmt.Printf("Copied new file: %s\n", destPath)
				case "v":
					content, err := os.ReadFile(srcPath)
					if err != nil {
						fmt.Printf("Error reading file: %v\n", err)
						break
					}
					fmt.Printf("----- %s -----\n%s\n", file.Name(), string(content))
					fmt.Print("Copy after viewing? [y/N]: ")
					fmt.Scanln(&choice)
					if strings.ToLower(choice) == "y" {
						copyFile(srcPath, destPath)
						fmt.Printf("Copied new file: %s\n", destPath)
					} else {
						fmt.Println("Skipped.")
					}
				default:
					fmt.Println("Skipped.")
				}
			}

		}
	}
}
func filesAreEqual(path1, path2 string) (bool, error) {
	file1, err := os.ReadFile(path1)
	if err != nil {
		return false, err
	}
	file2, err := os.ReadFile(path2)
	if err != nil {
		return false, err
	}
	return string(file1) == string(file2), nil
}

func copyFile(src, dst string) {
	in, err := os.Open(src)
	if err != nil {
		panic(err)
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		panic(err)
	}
}
func showDiff(src, dst string) {
	cmd := exec.Command("diff", "-u", dst, src)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}

func DetectChanges(repoPath string) ([]string, []string, error) {
	cmd := exec.Command("git", "-C", repoPath, "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return nil, nil, err
	}

	var modified, deleted []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if len(line) < 4 {
			continue
		}
		status := strings.TrimSpace(line[:2])
		file := strings.TrimSpace(line[3:])
		if strings.HasSuffix(file, ".yaml") || strings.HasSuffix(file, ".yml") {
			switch status {
			case "M", " M", "MM", "AM", "MA":
				modified = append(modified, file)
			case "D", " D":
				deleted = append(deleted, file)
			}
		}
	}

	return modified, deleted, nil
}

func AskUserToRevertChanges(repoPath string, modified, deleted []string) bool {
	if len(modified) == 0 && len(deleted) == 0 {
		return true
	}

	fmt.Println("Detected changes in the following YAML files:")

	if len(modified) > 0 {
		fmt.Println("Modified files:")
		for _, f := range modified {
			fmt.Println(" -", f)
		}
	}

	if len(deleted) > 0 {
		fmt.Println("Deleted files:")
		for _, f := range deleted {
			fmt.Println(" -", f)
		}
	}

	fmt.Print("Do you want to discard local changes and pull fresh? (y/N): ")
	var answer string
	fmt.Scanln(&answer)
	answer = strings.ToLower(strings.TrimSpace(answer))
	return answer == "y" || answer == "yes"
}

// CheckCategory checks if a string is present in a category.
// It returns true if the string is found, otherwise false.
func CheckCategory(slice []string, str string) bool {
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
	err = os.RemoveAll(srcDir)
	if err != nil {
		panic(err)
	}
	fmt.Println("Pulled YAMLs are ready to use")
}
