package app

import (
	"fmt"
	"strings"

	addressbookTypes "github.com/CudoVentures/cudos-node/x/addressbook/types"
	admintypes "github.com/CudoVentures/cudos-node/x/admin/types"
	cudoMinttypes "github.com/CudoVentures/cudos-node/x/cudoMint/types"
	marketplaceTypes "github.com/CudoVentures/cudos-node/x/marketplace/types"
	nfttypes "github.com/CudoVentures/cudos-node/x/nft/types"
	gravitytypes "github.com/althea-net/cosmos-gravity-bridge/module/x/gravity/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/cosmos-sdk/x/group"

	addressbooktypes "github.com/CudoVentures/cudos-node/x/addressbook/types"
	cudominttypes "github.com/CudoVentures/cudos-node/x/cudoMint/types"
	marketplacetypes "github.com/CudoVentures/cudos-node/x/marketplace/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ica "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icahosttypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	ibcfeetypes "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctmmigrations "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint/migrations"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

func (app *CudosApp) RegisterUpgradeHandlers() {
	setHandlerForVersion_1_0(app)
	setHandlerForVersion_1_1(app)
	setHandlerForVersion_1_2(app)
}

func setHandlerForVersion_1_0(app *CudosApp) {
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

func setHandlerForVersion_1_1(app *CudosApp) {
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

func setHandlerForVersion_1_2(app *CudosApp) {
	const upgradeVersion string = "v1.2"

	for _, subspace := range app.ParamsKeeper.GetSubspaces() {
		subspace := subspace

		var keyTable paramstypes.KeyTable
		switch subspace.Name() {
		case authtypes.ModuleName:
			keyTable = authtypes.ParamKeyTable() //nolint:staticcheck
		case banktypes.ModuleName:
			keyTable = banktypes.ParamKeyTable() //nolint:staticcheck
		case stakingtypes.ModuleName:
			keyTable = stakingtypes.ParamKeyTable()
		case distrtypes.ModuleName:
			keyTable = distrtypes.ParamKeyTable() //nolint:staticcheck
		case slashingtypes.ModuleName:
			keyTable = slashingtypes.ParamKeyTable() //nolint:staticcheck
		case govtypes.ModuleName:
			keyTable = govv1.ParamKeyTable() //nolint:staticcheck
		case crisistypes.ModuleName:
			keyTable = crisistypes.ParamKeyTable() //nolint:staticcheck
		case ibctransfertypes.ModuleName:
			keyTable = ibctransfertypes.ParamKeyTable()
		case icahosttypes.SubModuleName:
			keyTable = icahosttypes.ParamKeyTable()
		case icacontrollertypes.SubModuleName:
			keyTable = icacontrollertypes.ParamKeyTable()
		case wasmtypes.ModuleName:
			keyTable = wasmtypes.ParamKeyTable() //nolint:staticcheck
		case gravitytypes.ModuleName:
			keyTable = gravitytypes.ParamKeyTable() //nolint:staticcheck
		case addressbooktypes.ModuleName:
			keyTable = addressbooktypes.ParamKeyTable() //nolint:staticcheck
		case cudominttypes.ModuleName:
			keyTable = cudominttypes.ParamKeyTable() //nolint:staticcheck
		case marketplacetypes.ModuleName:
			keyTable = marketplacetypes.ParamKeyTable() //nolint:staticcheck
		default:
			continue
		}

		if !subspace.HasKeyTable() {
			subspace.WithKeyTable(keyTable)
		}
	}

	baseAppLegacySS := app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramstypes.ConsensusParamsKeyTable())

	app.UpgradeKeeper.SetUpgradeHandler(upgradeVersion, func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// IBC 2 -> 3
		fromVM[icatypes.ModuleName] = app.mm.Modules[icatypes.ModuleName].(module.HasConsensusVersion).ConsensusVersion()
		controllerParams := icacontrollertypes.Params{
			ControllerEnabled: true,
		}
		hostParams := icahosttypes.Params{
			HostEnabled:   true,
			AllowMessages: []string{"/cosmos.bank.v1beta1.MsgSend"},
		}
		app.mm.Modules[icatypes.ModuleName].(ica.AppModule).InitModule(ctx, controllerParams, hostParams)

		// IBC 6 -> 7
		// prune expired tendermint consensus states to save storage space
		_, err := ibctmmigrations.PruneExpiredConsensusStates(ctx, app.appCodec, app.IBCKeeper.ClientKeeper)
		if err != nil {
			return nil, err
		}

		// COSMOS-SDK 46 -> 47 AND IBC 6 -> 7
		baseapp.MigrateParams(ctx, baseAppLegacySS, &app.ConsensusParamsKeeper)

		// IBC 7 -> 7.1
		// explicitly update the IBC 02-client params, adding the localhost client type
		params := app.IBCKeeper.ClientKeeper.GetParams(ctx)
		params.AllowedClients = append(params.AllowedClients, ibcexported.Localhost)
		app.IBCKeeper.ClientKeeper.SetParams(ctx, params)

		// MIGRATE GRAVITY TOKEN CONTRACT ADDRESS (IF APPLICABLE) FROM RINKEBY TO SEPOLIA
		cudosDenom := "acudos"
		tokenContractEthAddress, found := app.GravityKeeper.GetCosmosOriginatedERC20(ctx, cudosDenom)
		if !found {
			return nil, fmt.Errorf("unable to find token to erc20 token mapping by searching for denom %s", cudosDenom)
		}
		// if it is a testnet then change the token contract address
		if tokenContractEthAddress.GetAddress() == strings.ToLower("0x28ea52f3ee46CaC5a72f72e8B3A387C0291d586d") || tokenContractEthAddress.GetAddress() == strings.ToLower("0x12d474723cb8c02bcbf46cd335a3bb4c75e9de44") {
			store := ctx.KVStore(app.keys[gravitytypes.StoreKey])

			// delete old values
			store.Delete(gravitytypes.GetDenomToERC20Key(cudosDenom))
			store.Delete(gravitytypes.GetERC20ToDenomKey(*tokenContractEthAddress))

			// set new token contract address
			tokenContractEthAddress, errorSetAddress := gravitytypes.NewEthAddress("0xE92f6A5b005B8f98F30313463Ada5cb35500a919")
			if errorSetAddress != nil {
				return nil, errorSetAddress
			}

			// set new values
			store.Set(gravitytypes.GetDenomToERC20Key(cudosDenom), []byte(tokenContractEthAddress.GetAddress()))
			store.Set(gravitytypes.GetERC20ToDenomKey(*tokenContractEthAddress), []byte(cudosDenom))
		}

		return app.mm.RunMigrations(ctx, app.configurator, fromVM)
	})

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	if upgradeInfo.Name == upgradeVersion && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{
				crisistypes.ModuleName,
				consensustypes.ModuleName,
				ibcfeetypes.ModuleName,
				icahosttypes.SubModuleName,
				icacontrollertypes.SubModuleName,
				admintypes.ModuleName,
			},
		}

		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}
