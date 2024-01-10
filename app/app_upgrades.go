package app

import (
	"strings"

	cudoMinttypes "github.com/CudoVentures/cudos-node/x/cudoMint/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/authz"

	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

const AddressBookModuleName = "addressbook"
const NftModuleName = "nft"
const MarketplaceModuleName = "marketplace"
const GroupModuleName = "group"

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
			delete(fromVM, GroupModuleName)
			delete(fromVM, AddressBookModuleName)
			delete(fromVM, MarketplaceModuleName)

			if _, ok := fromVM[NftModuleName]; ok {
				if fromVM[NftModuleName] == 2 {
					fromVM[NftModuleName] = 1
				}
			} else {
				fromVM[NftModuleName] = 1
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
			Added: []string{authz.ModuleName, GroupModuleName, AddressBookModuleName, MarketplaceModuleName},
		}

		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}

func setHandlerForVersion_1_2(app *App) {
	const upgradeVersion string = "v1.2"

	app.UpgradeKeeper.SetUpgradeHandler(upgradeVersion, func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		return app.mm.RunMigrations(ctx, app.configurator, fromVM)
	})

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	if upgradeInfo.Name == upgradeVersion && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{},
			Deleted: []string{
				AddressBookModuleName,
				MarketplaceModuleName,
				NftModuleName,
				GroupModuleName,
			},
		}

		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}
