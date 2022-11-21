package marketplace

import (
	"math/rand"

	"github.com/CudoVentures/cudos-node/testutil/sample"
	marketplacesimulation "github.com/CudoVentures/cudos-node/x/marketplace/simulation"
	"github.com/CudoVentures/cudos-node/x/marketplace/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = marketplacesimulation.FindAccount
	_ = simappparams.StakePerAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
	opWeightMsgPublishCollection = "op_weight_msg_publish_collection"
	// TODO: Determine the simulation weight value
	defaultWeightMsgPublishCollection int = 20

	opWeightMsgPublishNft = "op_weight_msg_publish_nft"
	// TODO: Determine the simulation weight value
	defaultWeightMsgPublishNft int = 20

	opWeightMsgBuyNft = "op_weight_msg_buy_nft"
	// TODO: Determine the simulation weight value
	defaultWeightMsgBuyNft int = 100

	opWeightMsgMintNft = "op_weight_msg_mint_nft"
	// TODO: Determine the simulation weight value
	defaultWeightMsgMintNft int = 100

	opWeightMsgRemoveNft = "op_weight_msg_remove_nft"
	// TODO: Determine the simulation weight value
	defaultWeightMsgRemoveNft int = 100

	opWeightMsgVerifyCollection = "op_weight_msg_verify_collection"
	// TODO: Determine the simulation weight value
	defaultWeightMsgVerifyCollection int = 100

	opWeightMsgUnverifyCollection = "op_weight_msg_unverify_collection"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUnverifyCollection int = 100

	opWeightMsgCreateCollection = "op_weight_msg_create_collection"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateCollection int = 100

	opWeightMsgUpdateRoyalties = "op_weight_msg_update_royalties"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateRoyalties int = 100

	opWeightMsgUpdatePrice = "op_weight_msg_update_price"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdatePrice int = 100

	opWeightMsgAddAdmin = "op_weight_msg_add_admin"
	// TODO: Determine the simulation weight value
	defaultWeightMsgAddAdmin int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	marketplaceGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&marketplaceGenesis)
}

// ProposalContents doesn't return any content functions for governance proposals
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized  param changes for the simulator
func (am AppModule) RandomizedParams(_ *rand.Rand) []simtypes.ParamChange {

	return []simtypes.ParamChange{}
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgPublishCollection int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgPublishCollection, &weightMsgPublishCollection, nil,
		func(_ *rand.Rand) {
			weightMsgPublishCollection = defaultWeightMsgPublishCollection
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgPublishCollection,
		marketplacesimulation.SimulateMsgPublishCollection(am.accountKeeper, am.bankKeeper, am.nftKeeper, am.keeper),
	))

	var weightMsgPublishNft int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgPublishNft, &weightMsgPublishNft, nil,
		func(_ *rand.Rand) {
			weightMsgPublishNft = defaultWeightMsgPublishNft
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgPublishNft,
		marketplacesimulation.SimulateMsgPublishNft(am.accountKeeper, am.bankKeeper, am.nftKeeper, am.keeper),
	))

	var weightMsgBuyNft int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgBuyNft, &weightMsgBuyNft, nil,
		func(_ *rand.Rand) {
			weightMsgBuyNft = defaultWeightMsgBuyNft
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgBuyNft,
		marketplacesimulation.SimulateMsgBuyNft(am.accountKeeper, am.bankKeeper, am.nftKeeper, am.keeper),
	))

	var weightMsgMintNft int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgMintNft, &weightMsgMintNft, nil,
		func(_ *rand.Rand) {
			weightMsgMintNft = defaultWeightMsgMintNft
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgMintNft,
		marketplacesimulation.SimulateMsgMintNft(am.accountKeeper, am.bankKeeper, am.nftKeeper, am.keeper),
	))

	var weightMsgRemoveNft int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgRemoveNft, &weightMsgRemoveNft, nil,
		func(_ *rand.Rand) {
			weightMsgRemoveNft = defaultWeightMsgRemoveNft
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRemoveNft,
		marketplacesimulation.SimulateMsgRemoveNft(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgVerifyCollection int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgVerifyCollection, &weightMsgVerifyCollection, nil,
		func(_ *rand.Rand) {
			weightMsgVerifyCollection = defaultWeightMsgVerifyCollection
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgVerifyCollection,
		marketplacesimulation.SimulateMsgVerifyCollection(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUnverifyCollection int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUnverifyCollection, &weightMsgUnverifyCollection, nil,
		func(_ *rand.Rand) {
			weightMsgUnverifyCollection = defaultWeightMsgUnverifyCollection
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUnverifyCollection,
		marketplacesimulation.SimulateMsgUnverifyCollection(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgCreateCollection int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreateCollection, &weightMsgCreateCollection, nil,
		func(_ *rand.Rand) {
			weightMsgCreateCollection = defaultWeightMsgCreateCollection
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateCollection,
		marketplacesimulation.SimulateMsgCreateCollection(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateRoyalties int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateRoyalties, &weightMsgUpdateRoyalties, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateRoyalties = defaultWeightMsgUpdateRoyalties
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateRoyalties,
		marketplacesimulation.SimulateMsgUpdateRoyalties(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdatePrice int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdatePrice, &weightMsgUpdatePrice, nil,
		func(_ *rand.Rand) {
			weightMsgUpdatePrice = defaultWeightMsgUpdatePrice
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdatePrice,
		marketplacesimulation.SimulateMsgUpdatePrice(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgAddAdmin int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgAddAdmin, &weightMsgAddAdmin, nil,
		func(_ *rand.Rand) {
			weightMsgAddAdmin = defaultWeightMsgAddAdmin
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddAdmin,
		marketplacesimulation.SimulateMsgAddAdmin(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
