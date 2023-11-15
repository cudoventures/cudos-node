package cudoMint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/CudoVentures/cudos-node/x/cudoMint/keeper"
	"github.com/CudoVentures/cudos-node/x/cudoMint/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	k.SetMinter(ctx, data.Minter)
	k.SetParams(ctx, data.Params)
	// this line is used by starport scaffolding # genesis/module/init

	// this line is used by starport scaffolding # ibc/genesis/init
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	minter := k.GetMinter(ctx)
	params := k.GetParams(ctx)

	// this line is used by starport scaffolding # genesis/module/export

	// this line is used by starport scaffolding # ibc/genesis/export
	return types.NewGenesisState(minter, params)
}
