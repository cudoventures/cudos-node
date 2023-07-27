package addressbook

import (
	"fmt"
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
	addressbookGenesis := types.GenesisState{
		Params:      types.DefaultParams(),
		AddressList: []types.Address{},
	}
	for i, acc := range simState.Accounts {
		addressbookGenesis.AddressList = append(addressbookGenesis.AddressList, types.Address{
			Creator: acc.Address.String(),
			Network: "BTC",
			Label:   fmt.Sprintf("%d@testdenom", i+1),
		})
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&addressbookGenesis)
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
