/*
Copyright © 2023 Dagu Yota Hamada
*/
package main

import (
	cmd "github.com/ErdemOzgen/blackdagger/cmd"
	"github.com/ErdemOzgen/blackdagger/internal/constants"
)

func main() {
	cmd.Execute()
}

var version = "1.0.1"

func init() {
	constants.Version = version
}
