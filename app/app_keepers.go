package app

import (
	"path/filepath"

	"github.com/CosmWasm/wasmd/x/wasm"

	adminkeeper "github.com/CudoVentures/cudos-node/x/admin/keeper"
	admintypes "github.com/CudoVentures/cudos-node/x/admin/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
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

	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibctransferkeeper "github.com/cosmos/ibc-go/v2/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v2/modules/apps/transfer/types"
	ibchost "github.com/cosmos/ibc-go/v2/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v2/modules/core/keeper"

	ibcclient "github.com/cosmos/ibc-go/v2/modules/core/02-client"
	ibcclienttypes "github.com/cosmos/ibc-go/v2/modules/core/02-client/types"

	// this line is used by starport scaffolding # stargate/app/moduleImport
	cudoMintkeeper "github.com/CudoVentures/cudos-node/x/cudoMint/keeper"
	cudoMinttypes "github.com/CudoVentures/cudos-node/x/cudoMint/types"
	nftmodulekeeper "github.com/CudoVentures/cudos-node/x/nft/keeper"
	nftmoduletypes "github.com/CudoVentures/cudos-node/x/nft/types"

	gravitykeeper "github.com/althea-net/cosmos-gravity-bridge/module/x/gravity/keeper"
	gravitytypes "github.com/althea-net/cosmos-gravity-bridge/module/x/gravity/types"

	"github.com/cosmos/cosmos-sdk/x/group"
	groupkeeper "github.com/cosmos/cosmos-sdk/x/group/keeper"
)

func (app *App) AddKeepers(skipUpgradeHeights map[int64]bool, homePath string, appOpts servertypes.AppOptions) {
	app.ParamsKeeper = InitParamsKeeper(app.appCodec, app.cdc, app.keys[paramstypes.StoreKey], app.tkeys[paramstypes.TStoreKey])

	// set the BaseApp's parameter store
	app.BaseApp.SetParamStore(app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramskeeper.ConsensusParamsKeyTable()))

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(app.appCodec, app.keys[capabilitytypes.StoreKey], app.memKeys[capabilitytypes.MemStoreKey])

	// grant capabilities for the ibc and ibc-transfer modules
	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)
	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	scopedWasmKeeper := app.CapabilityKeeper.ScopeToModule(wasm.ModuleName)

	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedTransferKeeper = scopedTransferKeeper

	app.CapabilityKeeper.Seal()

	// add keepers
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		app.appCodec, app.keys[authtypes.StoreKey], app.GetSubspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, maccPerms,
	)

	app.AuthzKeeper = authzkeeper.NewKeeper(app.keys[authzkeeper.StoreKey], app.appCodec, app.BaseApp.MsgServiceRouter())

	bankKeeper := bankkeeper.NewBaseKeeper(
		app.appCodec, app.keys[banktypes.StoreKey], app.AccountKeeper, app.GetSubspace(banktypes.ModuleName), app.BlockedAddrs(),
	)
	stakingKeeper := stakingkeeper.NewKeeper(
		app.appCodec, app.keys[stakingtypes.StoreKey], app.AccountKeeper, bankKeeper, app.GetSubspace(stakingtypes.ModuleName),
	)

	app.DistrKeeper = distrkeeper.NewKeeper(
		app.appCodec, app.keys[distrtypes.StoreKey], app.GetSubspace(distrtypes.ModuleName), app.AccountKeeper, bankKeeper,
		&stakingKeeper, authtypes.FeeCollectorName, app.ModuleAccountAddrs(),
	)

	bankKeeper.SetDistrKeeper(app.DistrKeeper)

	app.BankKeeper = bankKeeper

	app.SlashingKeeper = slashingkeeper.NewKeeper(
		app.appCodec, app.keys[slashingtypes.StoreKey], &stakingKeeper, app.GetSubspace(slashingtypes.ModuleName),
	)
	app.CrisisKeeper = crisiskeeper.NewKeeper(
		app.GetSubspace(crisistypes.ModuleName), app.invCheckPeriod, app.BankKeeper, authtypes.FeeCollectorName,
	)
	app.UpgradeKeeper = upgradekeeper.NewKeeper(skipUpgradeHeights, app.keys[upgradetypes.StoreKey], app.appCodec, homePath, app.BaseApp)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.StakingKeeper = *stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	)

	app.feegrantKeeper = feegrantkeeper.NewKeeper(app.appCodec, app.keys[feegrant.StoreKey], app.AccountKeeper)

	// Create IBC Keeper
	app.IBCKeeper = ibckeeper.NewKeeper(
		app.appCodec, app.keys[ibchost.StoreKey], app.GetSubspace(ibchost.ModuleName), app.StakingKeeper, app.UpgradeKeeper, scopedIBCKeeper,
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

	supportedFeatures := "iterator,staking,stargate"
	customEncoderOptions := GetCustomMsgEncodersOptions()
	customQueryOptions := GetCustomMsgQueryOptions(app.NftKeeper)
	wasmOpts := append(customEncoderOptions, customQueryOptions...)

	// The last arguments can contain custom message handlers, and custom query handlers,
	// if we want to allow any custom callbacks
	app.wasmKeeper = wasm.NewKeeper(
		app.appCodec,
		app.keys[wasm.StoreKey],
		app.GetSubspace(wasm.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		app.DistrKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		scopedWasmKeeper,
		app.TransferKeeper,
		app.MsgServiceRouter(),
		app.GRPCQueryRouter(),
		wasmDir,
		wasmConfig,
		supportedFeatures,
		wasmOpts...,
	)

	app.adminKeeper = *adminkeeper.NewKeeper(
		app.appCodec, app.keys[admintypes.StoreKey], app.keys[admintypes.MemStoreKey],
		app.DistrKeeper, app.BankKeeper,
	)

	govKeeper := govtypes.NewRouter()

	// The gov proposal types can be individually enabled
	if len(GetEnabledProposals()) != 0 {
		govKeeper.AddRoute(wasm.RouterKey, wasm.NewWasmProposalHandler(app.wasmKeeper, GetEnabledProposals()))
	}

	// register the proposal types
	govKeeper.AddRoute(govtypes.RouterKey, govtypes.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
		AddRoute(distrtypes.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.DistrKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(app.IBCKeeper.ClientKeeper))

	app.GovKeeper = govkeeper.NewKeeper(
		app.appCodec, app.keys[govtypes.StoreKey], app.GetSubspace(govtypes.ModuleName), app.AccountKeeper, app.BankKeeper,
		&stakingKeeper, govKeeper,
	)

	// Create Transfer Keepers
	app.TransferKeeper = ibctransferkeeper.NewKeeper(
		app.appCodec, app.keys[ibctransfertypes.StoreKey], app.GetSubspace(ibctransfertypes.ModuleName),
		app.IBCKeeper.ChannelKeeper, &app.IBCKeeper.PortKeeper,
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
