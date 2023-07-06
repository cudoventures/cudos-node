package addressbook

import (
	"math/rand"

	// simappparams "cosmossdk.io/simapp/params"
	"github.com/CudoVentures/cudos-node/testutil/sample"
	addressbooksimulation "github.com/CudoVentures/cudos-node/x/addressbook/simulation"
	"github.com/CudoVentures/cudos-node/x/addressbook/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = addressbooksimulation.FindAccount
	// _ = simappparams.StakePerAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
	opWeightMsgCreateAddress = "op_weight_msg_address"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateAddress int = 100

	opWeightMsgUpdateAddress = "op_weight_msg_address"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateAddress int = 100

	opWeightMsgDeleteAddress = "op_weight_msg_address"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteAddress int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	addressbookGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		AddressList: []types.Address{
			{
				Creator: sample.AccAddress(),
				Network: "BTC",
				Label:   "1@testdenom",
			},
			{
				Creator: sample.AccAddress(),
				Network: "ETH",
				Label:   "2@newdenom",
			},
		},
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&addressbookGenesis)
}

// ProposalContents doesn't return any content functions for governance proposals
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized  param changes for the simulator
func (am AppModule) RandomizedParams(_ *rand.Rand) []simtypes.LegacyParamChange {
	return []simtypes.LegacyParamChange{}
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCreateAddress int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreateAddress, &weightMsgCreateAddress, nil,
		func(_ *rand.Rand) {
			weightMsgCreateAddress = defaultWeightMsgCreateAddress
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateAddress,
		addressbooksimulation.SimulateMsgCreateAddress(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateAddress int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateAddress, &weightMsgUpdateAddress, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateAddress = defaultWeightMsgUpdateAddress
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateAddress,
		addressbooksimulation.SimulateMsgUpdateAddress(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeleteAddress int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgDeleteAddress, &weightMsgDeleteAddress, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteAddress = defaultWeightMsgDeleteAddress
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteAddress,
		addressbooksimulation.SimulateMsgDeleteAddress(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
