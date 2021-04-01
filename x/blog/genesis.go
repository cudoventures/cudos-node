package blog

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"cudos.org/cudos-poc-01/x/blog/keeper"
	"cudos.org/cudos-poc-01/x/blog/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return types.DefaultGenesis()
}
