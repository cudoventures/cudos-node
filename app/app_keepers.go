package app

import (
	"path/filepath"

	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/authz"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/spf13/cast"

	adminkeeper "github.com/CudoVentures/cudos-node/x/admin/keeper"
	admintypes "github.com/CudoVentures/cudos-node/x/admin/types"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/upgrade"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"

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

func (app *CudosApp) AddKeepers(appOpts servertypes.AppOptions) {
	// ParamsKeeper
	app.ParamsKeeper = initParamsKeeper(app.appCodec, app.legacyAminoCodec, app.keys[paramstypes.StoreKey], app.tkeys[paramstypes.TStoreKey])

	// ConsensusParamsKeeper
	app.ConsensusParamsKeeper = consensusparamkeeper.NewKeeper(app.appCodec, app.keys[consensusparamstypes.StoreKey], authtypes.NewModuleAddress(govtypes.ModuleName).String())
	app.BaseApp.SetParamStore(&app.ConsensusParamsKeeper)

	// CapabilityKeeper
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(app.appCodec, app.keys[capabilitytypes.StoreKey], app.memKeys[capabilitytypes.MemStoreKey])

	app.ScopedIBCKeeper = app.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	// app.ScopedICAHostKeeper = app.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)
	// app.ScopedICAControllerKeeper = app.CapabilityKeeper.ScopeToModule(icacontrollertypes.SubModuleName)
	app.ScopedTransferKeeper = app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	app.ScopedWasmKeeper = app.CapabilityKeeper.ScopeToModule(wasm.ModuleName)
	app.CapabilityKeeper.Seal()

	// AccountKeeper
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		app.appCodec,
		app.keys[authtypes.StoreKey],
		authtypes.ProtoBaseAccount,
		maccPerms,
		sdk.Bech32MainPrefix,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// BankKeeper
	bankKeeper := bankkeeper.NewBaseKeeper(
		app.appCodec,
		app.keys[banktypes.StoreKey],
		app.AccountKeeper,
		app.BlockedAddrs(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// StakingKeeper
	app.StakingKeeper = stakingkeeper.NewKeeper(
		app.appCodec,
		app.keys[stakingtypes.StoreKey],
		app.AccountKeeper,
		bankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// DistrKeeper
	app.DistrKeeper = distrkeeper.NewKeeper(
		app.appCodec,
		app.keys[distrtypes.StoreKey],
		app.AccountKeeper,
		bankKeeper,
		app.StakingKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Bank<->Distr
	bankKeeper.SetDistrKeeper(app.DistrKeeper)
	app.BankKeeper = bankKeeper

	// SlashingKeeper
	app.SlashingKeeper = slashingkeeper.NewKeeper(
		app.appCodec,
		app.legacyAminoCodec,
		app.keys[slashingtypes.StoreKey],
		app.StakingKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Staking<->Distr<->Slashing
	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.StakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	)

	// CrisisKeeper
	invCheckPeriod := cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod))
	app.CrisisKeeper = crisiskeeper.NewKeeper(
		app.appCodec,
		app.keys[crisistypes.StoreKey],
		invCheckPeriod,
		app.BankKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// FeeGrantKeeper
	app.FeeGrantKeeper = feegrantkeeper.NewKeeper(
		app.appCodec,
		app.keys[feegrant.StoreKey],
		app.AccountKeeper,
	)

	// AuthzKeeper
	app.AuthzKeeper = authzkeeper.NewKeeper(
		app.keys[authzkeeper.StoreKey],
		app.appCodec,
		app.BaseApp.MsgServiceRouter(),
		app.AccountKeeper,
	)

	// GroupKeeper
	groupConfig := group.DefaultConfig()
	app.GroupKeeper = groupkeeper.NewKeeper(
		app.keys[group.StoreKey],
		app.appCodec,
		app.MsgServiceRouter(),
		app.AccountKeeper,
		groupConfig,
	)

	// UpgradeKeeper
	// get skipUpgradeHeights from the app options
	skipUpgradeHeights := map[int64]bool{}
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}
	homePath := cast.ToString(appOpts.Get(flags.FlagHome))
	app.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		app.keys[upgradetypes.StoreKey],
		app.appCodec,
		homePath,
		app.BaseApp,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// IBCKeeper
	app.IBCKeeper = ibckeeper.NewKeeper(
		app.appCodec,
		app.keys[ibcexported.StoreKey],
		app.GetSubspace(ibcexported.ModuleName),
		app.StakingKeeper,
		app.UpgradeKeeper,
		app.ScopedIBCKeeper,
	)

	// GovKeeper
	govConfig := govtypes.DefaultConfig()
	govKeeper := govkeeper.NewKeeper(
		app.appCodec,
		app.keys[govtypes.StoreKey],
		app.AccountKeeper, app.BankKeeper,
		app.StakingKeeper,
		app.MsgServiceRouter(),
		govConfig,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	app.GovKeeper = *govKeeper.SetHooks(
		govtypes.NewMultiGovHooks(
		// register the governance hooks
		),
	)
	// Register the proposal types
	// Deprecated: Avoid adding new handlers, instead use the new proposal flow
	// by granting the governance module the right to execute the message.
	// See: https://docs.cosmos.network/main/modules/gov#proposal-messages
	govRouter := govv1beta1.NewRouter()
	if len(GetEnabledProposals()) != 0 {
		govRouter.AddRoute(wasm.RouterKey, wasm.NewWasmProposalHandler(app.WasmKeeper, GetEnabledProposals()))
	}
	govRouter.AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler)
	govRouter.AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper))
	govRouter.AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper))
	govRouter.AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(app.IBCKeeper.ClientKeeper))
	app.GovKeeper.SetLegacyRouter(govRouter)

	// NftKeeper
	app.NftKeeper = *nftmodulekeeper.NewKeeper(
		app.appCodec,
		app.keys[nftmoduletypes.StoreKey],
		app.keys[nftmoduletypes.MemStoreKey],
	)

	// MarketplaceKeeper
	app.MarketplaceKeeper = *marketplacekeeper.NewKeeper(
		app.appCodec,
		app.keys[marketplacetypes.StoreKey],
		app.keys[marketplacetypes.MemStoreKey],
		app.GetSubspace(marketplacetypes.ModuleName),
		app.BankKeeper,
		app.NftKeeper,
	)

	// AddressbookKeeper
	app.AddressbookKeeper = *addressbookkeeper.NewKeeper(
		app.appCodec,
		app.keys[addressbooktypes.StoreKey],
		app.keys[addressbooktypes.MemStoreKey],
		app.GetSubspace(addressbooktypes.ModuleName),
	)

	// Create evidence Keeper for to register the IBC light client misbehaviour evidence route
	// If evidence needs to be handled for the app, set routes in router here and seal
	app.EvidenceKeeper = *evidencekeeper.NewKeeper(
		app.appCodec,
		app.keys[evidencetypes.StoreKey],
		app.StakingKeeper,
		app.SlashingKeeper,
	)

	// IBCFeeKeeper
	app.IBCFeeKeeper = ibcfeekeeper.NewKeeper(
		app.appCodec, app.keys[ibcfeetypes.StoreKey],
		app.IBCKeeper.ChannelKeeper, // may be replaced with IBC middleware
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		app.BankKeeper,
	)

	// TransferKeeper
	app.TransferKeeper = ibctransferkeeper.NewKeeper(
		app.appCodec,
		app.keys[ibctransfertypes.StoreKey],
		app.GetSubspace(ibctransfertypes.ModuleName),
		app.IBCFeeKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		app.ScopedTransferKeeper,
	)

	// ICAControllerKeeper
	// app.ICAControllerKeeper = icacontrollerkeeper.NewKeeper(
	// 	app.appCodec,
	// 	app.keys[icacontrollertypes.StoreKey],
	// 	app.GetSubspace(icacontrollertypes.SubModuleName),
	// 	app.IBCFeeKeeper, // use ics29 fee as ics4Wrapper in middleware stack
	// 	app.IBCKeeper.ChannelKeeper,
	// 	&app.IBCKeeper.PortKeeper,
	// 	app.ScopedICAControllerKeeper,
	// 	app.MsgServiceRouter(),
	// )

	// GravityKeeper
	app.GravityKeeper = gravitykeeper.NewKeeper(
		app.appCodec,
		app.keys[gravitytypes.StoreKey],
		app.GetSubspace(gravitytypes.ModuleName),
		app.StakingKeeper,
		app.BankKeeper,
		app.SlashingKeeper,
		app.AccountKeeper,
	)

	// WasmKeeper
	wasmDir := filepath.Join(homePath, "wasm")
	wasmConfig, err := wasm.ReadWasmConfig(appOpts)
	if err != nil {
		panic("error while reading wasm config: " + err.Error())
	}
	supportedFeatures := "iterator,staking,stargate,cosmwasm_1_1,cosmwasm_1_2"

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
		app.ScopedWasmKeeper,
		app.TransferKeeper,
		app.MsgServiceRouter(),
		app.GRPCQueryRouter(),
		wasmDir,
		wasmConfig,
		supportedFeatures,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		GetCustomPlugins(app.NftKeeper, app.MarketplaceKeeper, app.AddressbookKeeper)...,
	)

	// AdminKeeper
	app.AdminKeeper = *adminkeeper.NewKeeper(
		app.appCodec,
		app.keys[admintypes.StoreKey],
		app.keys[admintypes.MemStoreKey],
		app.DistrKeeper,
		app.BankKeeper,
	)

	// CudoMintKeeper
	app.CudoMintKeeper = *cudoMintkeeper.NewKeeper(
		app.appCodec,
		app.keys[cudoMinttypes.StoreKey],
		app.keys[cudoMinttypes.MemStoreKey],
		app.BankKeeper,
		app.AccountKeeper,
		app.GetSubspace(cudoMinttypes.ModuleName),
		authtypes.FeeCollectorName,
	)
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName)
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(feegrant.ModuleName)
	paramsKeeper.Subspace(authz.ModuleName)
	paramsKeeper.Subspace(group.ModuleName)
	// ibc-related
	paramsKeeper.Subspace(ibcexported.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	// paramsKeeper.Subspace(icahosttypes.SubModuleName)
	// paramsKeeper.Subspace(icacontrollertypes.SubModuleName)
	// external
	paramsKeeper.Subspace(gravitytypes.ModuleName)
	paramsKeeper.Subspace(wasm.ModuleName)
	// cudos
	paramsKeeper.Subspace(cudoMinttypes.ModuleName)
	paramsKeeper.Subspace(nftmoduletypes.ModuleName)
	paramsKeeper.Subspace(addressbooktypes.ModuleName)
	paramsKeeper.Subspace(marketplacetypes.ModuleName)

	return paramsKeeper
}
