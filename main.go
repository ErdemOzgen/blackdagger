/*
Copyright © 2024 Mahmut Erdem Ozgen
Forked from Dagu-dev
*/
package main

import (
	cmd "github.com/ErdemOzgen/blackdagger/cmd"
	"github.com/ErdemOzgen/blackdagger/internal/constants"
)

func main() {
	cmd.Execute()
}

var version = "1.0.3"

func init() {
	constants.Version = version
}
