package app

import (
	nfttypes "github.com/CudoVentures/cudos-node/x/nft/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/cosmos-sdk/x/group"

	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func (app *App) SetUpgradeHandlers() {
	setV110Handler(app)
}

func setV110Handler(app *App) {
	const upgradeVersion string = "v1.1.0"

	app.UpgradeKeeper.SetUpgradeHandler(upgradeVersion, func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		if len(fromVM) == 0 {
			fromVM = app.mm.GetVersionMap()
			delete(fromVM, authz.ModuleName)
			delete(fromVM, group.ModuleName)

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
			Added: []string{authz.ModuleName, group.ModuleName},
		}

		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}
