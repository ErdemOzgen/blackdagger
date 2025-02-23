package main

import (
	"os"

	"github.com/ErdemOzgen/blackdagger/cmd"
	"github.com/ErdemOzgen/blackdagger/internal/constants"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}

var version = "1.0.7"

func init() {
	constants.Version = version
}
