package main

import (
	"os"

	"github.com/CudoVentures/cudos-node/app"
	"github.com/CudoVentures/cudos-node/cmd/cudos-noded/cmd"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func main() {
	sdk.DefaultPowerReduction = sdk.NewIntFromUint64(1000000000000000000)
	app.SetConfig()
	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
