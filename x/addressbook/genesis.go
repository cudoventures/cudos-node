package addressbook

import (
	"github.com/CudoVentures/cudos-node/x/addressbook/keeper"
	"github.com/CudoVentures/cudos-node/x/addressbook/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the address
	for _, elem := range genState.AddressList {
		k.CreateNewAddress(ctx, elem)
	}
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.AddressList = k.GetAllAddress(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
