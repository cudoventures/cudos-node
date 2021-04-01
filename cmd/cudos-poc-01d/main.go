package main

import (
	"os"

	"cudos.org/cudos-poc-01/app"
	"cudos.org/cudos-poc-01/cmd/cudos-poc-01d/cmd"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
