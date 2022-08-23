package token

import (
	"math/rand"

	"github.com/CudoVentures/cudos-node/testutil/sample"
	tokensimulation "github.com/CudoVentures/cudos-node/x/token/simulation"
	"github.com/CudoVentures/cudos-node/x/token/types"
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
	_ = tokensimulation.FindAccount
	_ = simappparams.StakePerAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
	opWeightMsgCreateToken = "op_weight_msg_token"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateToken int = 100

	opWeightMsgUpdateToken = "op_weight_msg_token"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateToken int = 100

	opWeightMsgDeleteToken = "op_weight_msg_token"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteToken int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	tokenGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		TokenList: []types.Token{
			{
				Owner: sample.AccAddress(),
				Denom: "0",
			},
			{
				Owner: sample.AccAddress(),
				Denom: "1",
			},
		},
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&tokenGenesis)
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

	var weightMsgCreateToken int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreateToken, &weightMsgCreateToken, nil,
		func(_ *rand.Rand) {
			weightMsgCreateToken = defaultWeightMsgCreateToken
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateToken,
		tokensimulation.SimulateMsgCreateToken(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateToken int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateToken, &weightMsgUpdateToken, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateToken = defaultWeightMsgUpdateToken
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateToken,
		tokensimulation.SimulateMsgUpdateToken(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeleteToken int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgDeleteToken, &weightMsgDeleteToken, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteToken = defaultWeightMsgDeleteToken
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteToken,
		tokensimulation.SimulateMsgDeleteToken(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
