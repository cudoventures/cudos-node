package main

import (
	"errors"
	"io"

	"github.com/CudoVentures/cudos-node/app"
	"github.com/CudoVentures/cudos-node/app/params"
	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
)

type AppCreator struct {
	encCfg *params.EncodingConfig
}

// newApp is an AppCreator
func (appCreator *AppCreator) newApp(logger log.Logger, db dbm.DB, traceStore io.Writer, appOpts servertypes.AppOptions) servertypes.Application {
	baseappOptions := server.DefaultBaseappOptions(appOpts)

	return app.NewCudosApp(
		logger,
		db,
		traceStore,
		true,
		*appCreator.encCfg,
		appOpts,
		baseappOptions...,
	)
}

// appExport creates a new simapp (optionally at a given height)
func (appCreator *AppCreator) appExport(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	height int64,
	forZeroHeight bool,
	jailAllowedAddrs []string,
	appOpts servertypes.AppOptions,
	modulesToExport []string,
) (servertypes.ExportedApp, error) {
	homePath, ok := appOpts.Get(flags.FlagHome).(string)
	if !ok || homePath == "" {
		return servertypes.ExportedApp{}, errors.New("application home not set")
	}

	cudosApp := app.NewCudosApp(logger, db, traceStore, height == -1, *appCreator.encCfg, appOpts)
	if height != -1 {
		if err := cudosApp.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	}

	return cudosApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs, modulesToExport)
}
