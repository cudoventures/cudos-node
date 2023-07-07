package app

import (
	"strings"

	addressbookTypes "github.com/CudoVentures/cudos-node/x/addressbook/types"
	cudoMinttypes "github.com/CudoVentures/cudos-node/x/cudoMint/types"
	marketplaceTypes "github.com/CudoVentures/cudos-node/x/marketplace/types"
	nfttypes "github.com/CudoVentures/cudos-node/x/nft/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/authz"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	"github.com/cosmos/cosmos-sdk/x/group"

	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func (app *App) SetUpgradeHandlers() {
	setHandlerForVersion_1_0(app)
	setHandlerForVersion_1_1(app)
	setHandlerForVersion_1_2(app)
}

func setHandlerForVersion_1_0(app *App) {
	app.UpgradeKeeper.SetUpgradeHandler("v1.0", func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ss, ok := app.ParamsKeeper.GetSubspace(cudoMinttypes.ModuleName)
		if ok {
			bpd := ss.GetRaw(ctx, []byte("BlocksPerDay"))

			bpdString := strings.ReplaceAll(string(bpd), "\"", "")

			bpdInt, parseOk := sdk.NewIntFromString(bpdString)
			if parseOk {
				ss.Set(ctx, []byte("IncrementModifier"), bpdInt)
			}
		}

		return fromVM, nil
	})
}

func setHandlerForVersion_1_1(app *App) {
	const upgradeVersion string = "v1.1"

	app.UpgradeKeeper.SetUpgradeHandler(upgradeVersion, func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		if len(fromVM) == 0 {
			fromVM = app.mm.GetVersionMap()
			delete(fromVM, authz.ModuleName)
			delete(fromVM, group.ModuleName)
			delete(fromVM, addressbookTypes.ModuleName)
			delete(fromVM, marketplaceTypes.ModuleName)

			if _, ok := fromVM[nfttypes.ModuleName]; ok {
				if fromVM[nfttypes.ModuleName] == 2 {
					fromVM[nfttypes.ModuleName] = 1
				}
			} else {
				fromVM[nfttypes.ModuleName] = 1
			}
		}

		return app.mm.RunMigrations(ctx, app.configurator, fromVM)
	})

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	if upgradeInfo.Name == upgradeVersion && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{authz.ModuleName, group.ModuleName, addressbookTypes.ModuleName, marketplaceTypes.ModuleName},
		}

		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}

func setHandlerForVersion_1_2(app *App) {
	const upgradeVersion string = "v1.2"

	baseAppLegacySS := app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramstypes.ConsensusParamsKeyTable())

	app.UpgradeKeeper.SetUpgradeHandler(upgradeVersion, func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		baseapp.MigrateParams(ctx, baseAppLegacySS, &app.ConsensusParamsKeeper)
		return app.mm.RunMigrations(ctx, app.configurator, fromVM)
	})

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	if upgradeInfo.Name == upgradeVersion && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{
				crisistypes.ModuleName, consensustypes.ModuleName,
			},
		}

		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}
