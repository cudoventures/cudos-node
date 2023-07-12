package app

import (
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/runtime"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
)

const Name = "cudos-node"

var (
	DefaultNodeHome string
)

var (
	_ runtime.AppI            = (*CudosApp)(nil)
	_ servertypes.Application = (*CudosApp)(nil)
)

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
