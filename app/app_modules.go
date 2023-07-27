package app

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	addressbook "github.com/CudoVentures/cudos-node/x/addressbook"
	"github.com/CudoVentures/cudos-node/x/admin"
	"github.com/CudoVentures/cudos-node/x/cudoMint"
	marketplace "github.com/CudoVentures/cudos-node/x/marketplace"
	nftmodule "github.com/CudoVentures/cudos-node/x/nft"
	"github.com/althea-net/cosmos-gravity-bridge/module/x/gravity"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	feegrantmod "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	groupmodule "github.com/cosmos/cosmos-sdk/x/group/module"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	ibcfee "github.com/cosmos/ibc-go/v7/modules/apps/29-fee"
	transfer "github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	"github.com/spf13/cast"
)

func (app *CudosApp) addModules(appOpts servertypes.AppOptions) {
	appCodec := app.appCodec
	txConfig := app.txConfig

	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	var transferStack porttypes.IBCModule
	transferStack = transfer.NewIBCModule(app.TransferKeeper)
	transferStack = ibcfee.NewIBCMiddleware(transferStack, app.IBCFeeKeeper)

	var wasmStack porttypes.IBCModule
	wasmStack = wasm.NewIBCHandler(app.WasmKeeper, app.IBCKeeper.ChannelKeeper, app.IBCFeeKeeper)
	wasmStack = ibcfee.NewIBCMiddleware(wasmStack, app.IBCFeeKeeper)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferStack)
	ibcRouter.AddRoute(wasm.ModuleName, wasmStack)
	app.IBCKeeper.SetRouter(ibcRouter)

	app.mm = module.NewManager(
		genutil.NewAppModule(app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx, txConfig),
		auth.NewAppModule(appCodec, app.AccountKeeper, randomGenesisAccounts, app.GetSubspace(authtypes.ModuleName)),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper, false),
		feegrantmod.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		gov.NewAppModule(appCodec, &app.GovKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(govtypes.ModuleName)),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(slashingtypes.ModuleName)),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(distrtypes.ModuleName)),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		params.NewAppModule(app.ParamsKeeper),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		groupmodule.NewAppModule(appCodec, app.GroupKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		consensus.NewAppModule(appCodec, app.ConsensusParamsKeeper),
		crisis.NewAppModule(app.CrisisKeeper, skipGenesisInvariants, app.GetSubspace(crisistypes.ModuleName)),
		// ibc
		ibc.NewAppModule(app.IBCKeeper),
		ibcfee.NewAppModule(app.IBCFeeKeeper),
		transfer.NewAppModule(app.TransferKeeper),
		// ica.NewAppModule(&app.ICAControllerKeeper, &app.ICAHostKeeper),
		// external
		gravity.NewAppModule(app.GravityKeeper, app.BankKeeper),
		wasm.NewAppModule(appCodec, &app.WasmKeeper, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.MsgServiceRouter(), app.GetSubspace(wasmtypes.ModuleName)),
		// cudos
		admin.NewAppModule(appCodec, app.AdminKeeper),
		cudoMint.NewAppModule(appCodec, app.CudoMintKeeper),
		nftmodule.NewAppModule(appCodec, app.NftKeeper, app.AccountKeeper, app.BankKeeper),
		marketplace.NewAppModule(appCodec, app.MarketplaceKeeper, app.AccountKeeper, app.BankKeeper, app.NftKeeper),
		addressbook.NewAppModule(appCodec, app.AddressbookKeeper, app.AccountKeeper, app.BankKeeper),
	)
}

func (app *CudosApp) orderSimulationModules() {
	var gravityModuleIndex int = -1
	for i, module := range app.sm.Modules {
		if _, ok := module.(gravity.AppModule); ok {
			gravityModuleIndex = i
			break
		}
	}

	if gravityModuleIndex != -1 {
		modulesSize := len(app.sm.Modules)
		gravityModule := app.sm.Modules[gravityModuleIndex]
		for i := gravityModuleIndex; i < modulesSize-1; i++ {
			app.sm.Modules[i] = app.sm.Modules[i+1]
		}
		app.sm.Modules[modulesSize-1] = gravityModule
	}
}
