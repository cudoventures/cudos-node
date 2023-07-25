package simapp

import (
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/runtime"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
)

const Name = "cudos-sim-node"

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string
)

var (
	_ runtime.AppI            = (*CudosSimApp)(nil)
	_ servertypes.Application = (*CudosSimApp)(nil)
)

func init() {
	InitializeGlobalAppVariables()
}

func InitializeGlobalAppVariables() {
	cudosHome, present := os.LookupEnv("CUDOS_HOME")
	if !present {
		userHomeDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		DefaultNodeHome = filepath.Join(userHomeDir, "cudos-data")
	} else {
		DefaultNodeHome = cudosHome
	}
}
