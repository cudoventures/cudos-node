package app

import (
	"io"
	"os"
	"path/filepath"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	reflectionv1 "cosmossdk.io/api/cosmos/reflection/v1"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	runtimeservices "github.com/cosmos/cosmos-sdk/runtime/services"
	"github.com/spf13/cast"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmos "github.com/cometbft/cometbft/libs/os"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	ibcfee "github.com/cosmos/ibc-go/v7/modules/apps/29-fee"

	appparams "github.com/CudoVentures/cudos-node/app/params"
	addressbooktypes "github.com/CudoVentures/cudos-node/x/addressbook/types"
	"github.com/CudoVentures/cudos-node/x/admin"
	admintypes "github.com/CudoVentures/cudos-node/x/admin/types"
	marketplacetypes "github.com/CudoVentures/cudos-node/x/marketplace/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/consensus"

	// Authz - Authorization for accounts to perform actions on behalf of other accounts.
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"

	// vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantmod "github.com/cosmos/cosmos-sdk/x/feegrant/module"

	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	transfer "github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"

	// this line is used by starport scaffolding # stargate/app/moduleImport
	"github.com/althea-net/cosmos-gravity-bridge/module/x/gravity"
	gravitytypes "github.com/althea-net/cosmos-gravity-bridge/module/x/gravity/types"

	"github.com/CudoVentures/cudos-node/x/cudoMint"
	cudominttypes "github.com/CudoVentures/cudos-node/x/cudoMint/types"
	nftmodule "github.com/CudoVentures/cudos-node/x/nft"
	nftmoduletypes "github.com/CudoVentures/cudos-node/x/nft/types"

	addressbook "github.com/CudoVentures/cudos-node/x/addressbook"
	marketplace "github.com/CudoVentures/cudos-node/x/marketplace"
	"github.com/cosmos/cosmos-sdk/x/group"
	groupmodule "github.com/cosmos/cosmos-sdk/x/group/module"
)

const Name = "cudos-node"

var (
	DefaultNodeHome string
)

var (
	_ CosmosApp               = (*App)(nil)
	_ servertypes.Application = (*App)(nil)
)

func init() {
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

// New returns a reference to an initialized Cudos.
// NewSimApp returns a reference to an initialized SimApp.
func New(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool, skipUpgradeHeights map[int64]bool,
	homePath string, invCheckPeriod uint, encodingConfig appparams.EncodingConfig,
	appOpts servertypes.AppOptions, baseAppOptions ...func(*baseapp.BaseApp),
) *App {

	appCodec := encodingConfig.Codec
	cdc := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bApp := baseapp.NewBaseApp(Name, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	// bApp.SetAppVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey,
		authzkeeper.StoreKey,
		banktypes.StoreKey,
		stakingtypes.StoreKey,
		crisistypes.StoreKey,
		distrtypes.StoreKey,
		slashingtypes.StoreKey,
		govtypes.StoreKey,
		paramstypes.StoreKey,
		consensusparamtypes.StoreKey,
		ibcexported.StoreKey,
		upgradetypes.StoreKey,
		evidencetypes.StoreKey,
		ibctransfertypes.StoreKey,
		capabilitytypes.StoreKey,
		// this line is used by starport scaffolding # stargate/app/storeKey
		nftmoduletypes.StoreKey,
		cudominttypes.StoreKey,
		wasm.StoreKey,
		gravitytypes.StoreKey,
		feegrant.StoreKey,
		group.StoreKey,
		addressbooktypes.StoreKey,
		marketplacetypes.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	app := &App{
		BaseApp:           bApp,
		cdc:               cdc,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())

	// add keepers
	app.AddKeepers(skipUpgradeHeights, homePath, appOpts)

	// set upgrades
	app.SetUpgradeHandlers()

	marketplaceModule := marketplace.NewAppModule(appCodec, app.MarketplaceKeeper, app.AccountKeeper, app.BankKeeper)

	addressbookModule := addressbook.NewAppModule(appCodec, app.AddressbookKeeper, app.AccountKeeper, app.BankKeeper)

	transferModule := transfer.NewAppModule(app.TransferKeeper)

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
		genutil.NewAppModule(app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx, encodingConfig.TxConfig),
		auth.NewAppModule(appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts, app.GetSubspace(authtypes.ModuleName)),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper, false),
		crisis.NewAppModule(app.CrisisKeeper, cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants)), app.GetSubspace(crisistypes.ModuleName)),
		gov.NewAppModule(appCodec, &app.GovKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(govtypes.ModuleName)),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(slashingtypes.ModuleName)),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(distrtypes.ModuleName)),
		staking.NewAppModule(appCodec, &app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		ibc.NewAppModule(app.IBCKeeper),
		params.NewAppModule(app.ParamsKeeper),
		transferModule,
		consensus.NewAppModule(appCodec, app.ConsensusParamsKeeper),
		wasm.NewAppModule(appCodec, &app.WasmKeeper, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.MsgServiceRouter(), app.GetSubspace(wasmtypes.ModuleName)),
		admin.NewAppModule(appCodec, app.adminKeeper),
		cudoMint.NewAppModule(appCodec, app.cudoMintKeeper),
		gravity.NewAppModule(app.GravityKeeper, app.BankKeeper),
		feegrantmod.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.feegrantKeeper, app.interfaceRegistry),
		// this line is used by starport scaffolding # stargate/app/appModule
		nftmodule.NewAppModule(appCodec, app.NftKeeper, app.AccountKeeper, app.BankKeeper),
		marketplaceModule,
		groupmodule.NewAppModule(appCodec, app.GroupKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		addressbookModule,
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	app.mm.SetOrderBeginBlockers(
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		cudominttypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		authtypes.ModuleName,
		authz.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		consensusparamtypes.ModuleName,
		gravitytypes.ModuleName,
		admintypes.ModuleName,
		nftmoduletypes.ModuleName,
		ibcexported.ModuleName,
		ibctransfertypes.ModuleName,
		wasmtypes.ModuleName,
		group.ModuleName,
		addressbooktypes.ModuleName,
		marketplacetypes.ModuleName,
	)

	app.mm.SetOrderEndBlockers(
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		authz.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		cudominttypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		consensusparamtypes.ModuleName,
		gravitytypes.ModuleName,
		admintypes.ModuleName,
		nftmoduletypes.ModuleName,
		ibcexported.ModuleName,
		ibctransfertypes.ModuleName,
		wasmtypes.ModuleName,
		group.ModuleName,
		addressbooktypes.ModuleName,
		marketplacetypes.ModuleName,
	)
	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	// NOTE: Gravity module must occur before genutils so that the pool are propertly initiallized
	// before gextxs
	app.mm.SetOrderInitGenesis(
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		authz.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		cudominttypes.ModuleName,
		crisistypes.ModuleName,
		gravitytypes.ModuleName, //MUST BE BEFORE GENUTIL!!!!
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		consensusparamtypes.ModuleName,
		ibcexported.ModuleName,
		ibctransfertypes.ModuleName,
		wasm.ModuleName,
		admintypes.ModuleName,
		nftmoduletypes.ModuleName,
		feegrant.ModuleName,
		upgradetypes.ModuleName,
		paramstypes.ModuleName,
		group.ModuleName,
		addressbooktypes.ModuleName,
		marketplacetypes.ModuleName,
	)

	app.mm.RegisterInvariants(app.CrisisKeeper)
	app.mm.RegisterServices(app.configurator)
	autocliv1.RegisterQueryServer(app.GRPCQueryRouter(), runtimeservices.NewAutoCLIQueryService(app.mm.Modules))

	reflectionSvc, err := runtimeservices.NewReflectionService()
	if err != nil {
		panic(err)
	}
	reflectionv1.RegisterReflectionServiceServer(app.GRPCQueryRouter(), reflectionSvc)

	anteHandler, err := ante.NewAnteHandler(
		ante.HandlerOptions{
			AccountKeeper:   app.AccountKeeper,
			BankKeeper:      app.BankKeeper,
			FeegrantKeeper:  app.feegrantKeeper,
			SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
			SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
		},
	)
	if err != nil {
		tmos.Exit(err.Error())
	}

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(anteHandler)
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}
	}

	return app
}
