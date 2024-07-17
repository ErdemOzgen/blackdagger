/*
Copyright © 2024 Mahmut Erdem Ozgen
Forked from Dagu-dev
*/
package main

import (
	"os"

	cmd "github.com/ErdemOzgen/blackdagger/cmd"
	"github.com/ErdemOzgen/blackdagger/internal/constants"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}

var version = "0.0.0"

func init() {
	constants.Version = version
}
