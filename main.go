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

var version = "1.0.6"

func init() {
	constants.Version = version
}
