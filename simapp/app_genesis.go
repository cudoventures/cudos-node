package simapp

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// GenesisState The genesis state of the blockchain is represented here as a map of raw json
// messages key'd by a identifier string.
// The identifier is used to determine which module genesis information belongs
// to so it may be appropriately routed during init chain.
// Within this application default genesis information is retrieved from
// the ModuleBasicManager which populates json from each BasicModule
// object provided to it during init.
type GenesisState map[string]json.RawMessage

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState(cdc codec.JSONCodec) GenesisState {
	return ModuleBasics.DefaultGenesis(cdc)
}

// copy implementation from x/auth/simulation/genesis.go by removing vesting accounts
func randomGenesisAccounts(simState *module.SimulationState) types.GenesisAccounts {
	genesisAccs := make(types.GenesisAccounts, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		genesisAccs[i] = types.NewBaseAccountWithAddress(acc.Address)
	}

	return genesisAccs
}