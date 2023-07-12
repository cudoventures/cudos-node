package main

import (
	"os"

	"github.com/CudoVentures/cudos-node/app"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {
	app.InitializeSdk()
	app.InitializeGlobalAppVariables()

	rootCmd, _ := NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
