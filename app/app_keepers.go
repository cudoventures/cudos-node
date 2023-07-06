package app

import (
	"path/filepath"

	"github.com/CosmWasm/wasmd/x/wasm"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	adminkeeper "github.com/CudoVentures/cudos-node/x/admin/keeper"
	admintypes "github.com/CudoVentures/cudos-node/x/admin/types"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/upgrade"

	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"

	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"

	ibcfeekeeper "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/keeper"
	ibcfeetypes "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/types"

	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibctransferkeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"

	ibcclient "github.com/cosmos/ibc-go/v7/modules/core/02-client"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"

	// this line is used by starport scaffolding # stargate/app/moduleImport
	cudoMintkeeper "github.com/CudoVentures/cudos-node/x/cudoMint/keeper"
	cudoMinttypes "github.com/CudoVentures/cudos-node/x/cudoMint/types"
	nftmodulekeeper "github.com/CudoVentures/cudos-node/x/nft/keeper"
	nftmoduletypes "github.com/CudoVentures/cudos-node/x/nft/types"

	gravitykeeper "github.com/althea-net/cosmos-gravity-bridge/module/x/gravity/keeper"
	gravitytypes "github.com/althea-net/cosmos-gravity-bridge/module/x/gravity/types"

	"github.com/cosmos/cosmos-sdk/x/group"
	groupkeeper "github.com/cosmos/cosmos-sdk/x/group/keeper"

	marketplacekeeper "github.com/CudoVentures/cudos-node/x/marketplace/keeper"
	marketplacetypes "github.com/CudoVentures/cudos-node/x/marketplace/types"

	addressbookkeeper "github.com/CudoVentures/cudos-node/x/addressbook/keeper"
	addressbooktypes "github.com/CudoVentures/cudos-node/x/addressbook/types"

	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamstypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
)

func (app *App) AddKeepers(skipUpgradeHeights map[int64]bool, homePath string, appOpts servertypes.AppOptions) {
	app.ParamsKeeper = InitParamsKeeper(app.appCodec, app.cdc, app.keys[paramstypes.StoreKey], app.tkeys[paramstypes.TStoreKey])

	// set the BaseApp's parameter store
	app.ConsensusParamsKeeper = consensusparamkeeper.NewKeeper(app.appCodec, app.keys[consensusparamstypes.StoreKey], authtypes.NewModuleAddress(govtypes.ModuleName).String())
	app.BaseApp.SetParamStore(&app.ConsensusParamsKeeper)

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(app.appCodec, app.keys[capabilitytypes.StoreKey], app.memKeys[capabilitytypes.MemStoreKey])

	// grant capabilities for the ibc and ibc-transfer modules
	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	scopedWasmKeeper := app.CapabilityKeeper.ScopeToModule(wasm.ModuleName)

	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedTransferKeeper = scopedTransferKeeper

	app.CapabilityKeeper.Seal()

	// add keepers
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		app.appCodec, app.keys[authtypes.StoreKey], authtypes.ProtoBaseAccount, maccPerms, sdk.Bech32MainPrefix, authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.AuthzKeeper = authzkeeper.NewKeeper(app.keys[authzkeeper.StoreKey], app.appCodec, app.BaseApp.MsgServiceRouter(), app.AccountKeeper)

	bankKeeper := bankkeeper.NewBaseKeeper(
		app.appCodec, app.keys[banktypes.StoreKey], app.AccountKeeper, app.BlockedAddrs(), authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	stakingKeeper := stakingkeeper.NewKeeper(
		app.appCodec, app.keys[stakingtypes.StoreKey], app.AccountKeeper, bankKeeper, authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.DistrKeeper = distrkeeper.NewKeeper(
		app.appCodec, app.keys[distrtypes.StoreKey], app.AccountKeeper, bankKeeper,
		stakingKeeper, authtypes.FeeCollectorName, authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	bankKeeper.SetDistrKeeper(app.DistrKeeper)

	app.BankKeeper = bankKeeper

	app.SlashingKeeper = slashingkeeper.NewKeeper(
		app.appCodec, app.cdc, app.keys[slashingtypes.StoreKey], stakingKeeper, authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	app.CrisisKeeper = crisiskeeper.NewKeeper(
		app.appCodec, app.keys[crisistypes.StoreKey], app.invCheckPeriod, app.BankKeeper, authtypes.FeeCollectorName, authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	app.UpgradeKeeper = upgradekeeper.NewKeeper(skipUpgradeHeights, app.keys[upgradetypes.StoreKey], app.appCodec, homePath, app.BaseApp, authtypes.NewModuleAddress(govtypes.ModuleName).String())

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	)
	app.StakingKeeper = *stakingKeeper

	app.feegrantKeeper = feegrantkeeper.NewKeeper(app.appCodec, app.keys[feegrant.StoreKey], app.AccountKeeper)

	// Create IBC Keeper
	app.IBCKeeper = ibckeeper.NewKeeper(
		app.appCodec, app.keys[ibcexported.StoreKey], app.GetSubspace(ibcexported.ModuleName), app.StakingKeeper, app.UpgradeKeeper, scopedIBCKeeper,
	)

	wasmDir := filepath.Join(homePath, "wasm")
	wasmConfig, err := wasm.ReadWasmConfig(appOpts)
	if err != nil {
		panic("error while reading wasm config: " + err.Error())
	}

	app.NftKeeper = *nftmodulekeeper.NewKeeper(
		app.appCodec,
		app.keys[nftmoduletypes.StoreKey],
		app.keys[nftmoduletypes.MemStoreKey],
	)

	app.AddressbookKeeper = *addressbookkeeper.NewKeeper(
		app.appCodec,
		app.keys[addressbooktypes.StoreKey],
		app.keys[addressbooktypes.MemStoreKey],
		app.GetSubspace(addressbooktypes.ModuleName),
	)

	app.MarketplaceKeeper = *marketplacekeeper.NewKeeper(
		app.appCodec,
		app.keys[marketplacetypes.StoreKey],
		app.keys[marketplacetypes.MemStoreKey],
		app.GetSubspace(marketplacetypes.ModuleName),
		app.BankKeeper,
		app.NftKeeper,
	)

	app.IBCFeeKeeper = ibcfeekeeper.NewKeeper(
		app.appCodec, app.keys[ibcfeetypes.StoreKey],
		app.IBCKeeper.ChannelKeeper, // may be replaced with IBC middleware
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper, app.AccountKeeper, app.BankKeeper,
	)

	supportedFeatures := "iterator,staking,stargate"
	customEncoderOptions := GetCustomMsgEncodersOptions()
	customQueryOptions := GetCustomMsgQueryOptions(app.NftKeeper)
	wasmOpts := append(customEncoderOptions, customQueryOptions...)

	// The last arguments can contain custom message handlers, and custom query handlers,
	// if we want to allow any custom callbacks
	app.WasmKeeper = wasm.NewKeeper(
		app.appCodec,
		app.keys[wasm.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		distrkeeper.NewQuerier(app.DistrKeeper),
		app.IBCFeeKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		scopedWasmKeeper,
		app.TransferKeeper,
		app.MsgServiceRouter(),
		app.GRPCQueryRouter(),
		wasmDir,
		wasmConfig,
		supportedFeatures,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		wasmOpts...,
	)

	app.adminKeeper = *adminkeeper.NewKeeper(
		app.appCodec, app.keys[admintypes.StoreKey], app.keys[admintypes.MemStoreKey],
		app.DistrKeeper, app.BankKeeper,
	)

	govRouter := govv1beta1.NewRouter()
	// The gov proposal types can be individually enabled
	if len(GetEnabledProposals()) != 0 {
		govRouter.AddRoute(wasm.RouterKey, wasm.NewWasmProposalHandler(app.WasmKeeper, GetEnabledProposals()))
	}
	// register the proposal types
	govRouter.AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(app.IBCKeeper.ClientKeeper))

	govConfig := govtypes.DefaultConfig()

	govKeeper := govkeeper.NewKeeper(
		app.appCodec, app.keys[govtypes.StoreKey], app.AccountKeeper, app.BankKeeper, stakingKeeper, app.MsgServiceRouter(), govConfig, authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	govKeeper.SetLegacyRouter(govRouter)
	app.GovKeeper = *govKeeper.SetHooks(
		govtypes.NewMultiGovHooks(
		// register the governance hooks
		),
	)

	// Create Transfer Keepers
	app.TransferKeeper = ibctransferkeeper.NewKeeper(
		app.appCodec, app.keys[ibctransfertypes.StoreKey], app.GetSubspace(ibctransfertypes.ModuleName),
		app.IBCFeeKeeper, app.IBCKeeper.ChannelKeeper, &app.IBCKeeper.PortKeeper,
		app.AccountKeeper, app.BankKeeper, scopedTransferKeeper,
	)

	// Create evidence Keeper for to register the IBC light client misbehaviour evidence route
	// If evidence needs to be handled for the app, set routes in router here and seal
	app.EvidenceKeeper = *evidencekeeper.NewKeeper(
		app.appCodec, app.keys[evidencetypes.StoreKey], &app.StakingKeeper, app.SlashingKeeper,
	)

	app.cudoMintKeeper = *cudoMintkeeper.NewKeeper(
		app.appCodec,
		app.keys[cudoMinttypes.StoreKey],
		app.keys[cudoMinttypes.MemStoreKey],
		app.BankKeeper,
		app.AccountKeeper,
		app.GetSubspace(cudoMinttypes.ModuleName),
		authtypes.FeeCollectorName,
	)

	app.GravityKeeper = gravitykeeper.NewKeeper(
		app.appCodec, app.keys[gravitytypes.StoreKey], app.GetSubspace(gravitytypes.ModuleName), stakingKeeper, app.BankKeeper, app.SlashingKeeper, app.AccountKeeper,
	)

	groupConfig := group.DefaultConfig()
	app.GroupKeeper = groupkeeper.NewKeeper(
		app.keys[group.StoreKey],
		app.appCodec,
		app.MsgServiceRouter(),
		app.AccountKeeper,
		groupConfig,
	)
}
